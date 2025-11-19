use std::io::{Read, Write};
use std::path::Path;
use std::time::Duration;
use std::{collections::HashMap, path::PathBuf};

use anyhow::{Context, bail};
use bollard::exec::StartExecOptions;
use bollard::query_parameters::{ListImagesOptions, StartContainerOptions};
use bollard::{
    Docker,
    exec::StartExecResults,
    query_parameters::{CreateContainerOptions, ListContainersOptions, RemoveContainerOptions},
    secret::{
        ContainerCreateBody, ContainerSummary, ContainerSummaryStateEnum, ExecConfig, HostConfig,
        MountPoint,
    },
};
use futures_util::StreamExt;
use termion::raw::IntoRawMode;
use tokio::io::AsyncWriteExt;
use tokio::runtime::{self, Runtime};

mod constants;

#[derive(Debug, Clone)]
pub enum DockerCommand {
    Rm,
    Up,
    Exec { command: Vec<String>, user: User },
}

pub struct DockerManager {
    client: Docker,
    volume_path: Option<PathBuf>,
    term_size: Option<(u16, u16)>,
    rt: Runtime,
}

impl DockerManager {
    pub fn new(volume_path: Option<PathBuf>) -> anyhow::Result<Self> {
        Ok(Self {
            client: Docker::connect_with_local_defaults().context("while connecting to Docker")?,
            volume_path,
            term_size: termion::terminal_size().ok(),
            rt: runtime::Builder::new_current_thread()
                .enable_all()
                .build()?,
        })
    }

    pub fn rm(&self) -> anyhow::Result<()> {
        self.rt.block_on(async { self.rm_contiker().await })
    }

    pub fn exec(&self, cmd: Vec<String>, user: User) -> anyhow::Result<()> {
        self.rt.block_on(async {
            // Check if we have a volume set
            if let Some(volume_path) = self.volume() {
                let docker_status = self.query_docker().await?;
                if let Some(docker_status) = docker_status {
                    // We already have a container, check if we still have the same volume mounted
                    // If not, we ask if we need to change
                    if !docker_status.has_mount(&volume_path) {
                        print!("Detected that the current running instance of contiker is using a different mount point. Do you want to restart and switch to {:?}? [y|n] ", volume_path);
                        let _ = std::io::stdout().flush();

                        let stdin = std::io::stdin();
                        let mut answer = String::new();
                        let mut valid_answer = false;
                        let mut answer_result = false;
                        while !valid_answer {
                            let _ = stdin.read_line(&mut answer).context("could not read from stdin")?;
                            let answer = answer.trim();
                            if "yes".starts_with(answer) {
                                answer_result = true;
                                valid_answer = true;
                            } else if "no".starts_with(answer) {
                                answer_result = false;
                                valid_answer = true;
                            }
                        }

                        if answer_result {
                            // We answered yes, so we restart by shutting down now, and later
                            // restarting
                            self.rm_contiker().await?;
                        }
                    }
                }
            }

            self.exec_docker(cmd, user).await
        })
    }

    pub fn is_up(&self) -> anyhow::Result<bool> {
        self.rt
            .block_on(async { Ok(self.query_docker().await?.is_some()) })
    }

    pub fn reset(&self) -> anyhow::Result<()> {
        self.rt.block_on(async {
            self.rm_contiker().await?;
            self.pull_contiker_container(false).await?;
            self.ensure_contiker_up(User::infer()).await?;

            Ok(())
        })
    }

    fn volume(&self) -> Option<PathBuf> {
        self.volume_path
            .as_ref()
            .map(|path| std::path::absolute(path).unwrap())
    }

    fn volume_path_or_cng(&self) -> Option<PathBuf> {
        self.volume_path
            .clone()
            .or_else(|| {
                std::env::var("CNG_PATH")
                    .ok()
                    .map(|path| Path::new(&path).to_path_buf())
            })
            .map(|path| std::path::absolute(&path).unwrap())
    }

    async fn pull_contiker_container(&self, use_cache: bool) -> anyhow::Result<()> {
        let mut needs_pulling = true;

        if use_cache {
            // Check if we have already the image in cache
            let has_contiker_image = !self
                .client
                .list_images(Some(ListImagesOptions {
                    filters: Some(HashMap::from([(
                        "reference".to_string(),
                        vec!["contiker/contiki-ng".to_string()],
                    )])),
                    ..Default::default()
                }))
                .await
                .context("while listing available images")?
                .is_empty();

            needs_pulling = !has_contiker_image;
        }

        if needs_pulling {
            // TODO: Refactor this to also use the API and not the CLI
            println!("Pulling image from Docker Hub");
            let command_status = std::process::Command::new("docker")
                .arg("pull")
                .arg("contiker/contiki-ng")
                .status()
                .context("while pulling the contiker image")?;

            if !command_status.success() {
                bail!("pull command did not finish succesfully");
            }
        }
        Ok(())
    }

    async fn query_docker(&self) -> anyhow::Result<Option<ContainerInfo>> {
        let filters = HashMap::from([(
            "name".to_string(),
            vec![constants::CONTIKER_CONTAINER.to_string()],
        )]);

        let options = ListContainersOptions {
            all: true,
            filters: Some(filters),
            ..Default::default()
        };

        let containers = self
            .client
            .list_containers(Some(options))
            .await
            .context("while querying docker for the contiker container")?;

        match containers.len() {
            0 => Ok(None),
            1 => Ok(Some(containers.into_iter().next().unwrap().try_into()?)),
            _ => {
                println!(
                    "[WARN] Multiple contiker containers found, only using the first one found"
                );

                Ok(Some(containers.into_iter().next().unwrap().try_into()?))
            }
        }
    }

    async fn ensure_contiker_up(&self, user: User) -> anyhow::Result<ContainerInfo> {
        let display = std::env::var("DISPLAY").unwrap_or_else(|_| {
            println!("[WARN] No DISPLAY environment found. This will result in you not being able to use cooja.");
            ":0".to_string()
        });

        // Pull if necessary
        self.pull_contiker_container(true).await?;

        let container = self.query_docker().await?;
        let container = match container {
            Some(container) => container,
            None => {
                // No container ready, so start one up
                let volume_path = self.volume_path_or_cng();
                let volume_path = volume_path
                    .as_ref()
                    .context("no volume mount point provided")?
                    .to_str()
                    .context("volume mount path cannot be converted to utf8")?;

                let response = self.client
                    .create_container(
                        Some(CreateContainerOptions {
                            name: Some(constants::CONTIKER_CONTAINER.to_string()),
                            ..Default::default()
                        }),
                        ContainerCreateBody {
                            image: Some("contiker/contiki-ng".to_string()),
                            env: Some(vec![
                                format!("DISPLAY={display}"),
                                "_JAVA_AWT_WM_NONREPARENTING=1".to_string(),
                                format!("LOCAL_UID={}", user.uid),
                                format!("LOCAL_GID={}", user.gid),
                                "JDK_JAVA_OPTIONS='-Dawt.useSystemAAFontSettings=on -Dswing.aatext=true -Dswing.defaultlaf=com.sun.java.swing.plaf.gtk.GTKLookAndFeel -Dsun.java2d.opengl=true'".to_string(),
                            ]),
                            host_config: Some(HostConfig {
                                privileged: Some(true),
                                ipc_mode: Some("host".to_string()),
                                binds: Some(vec![
                                    "/dev/:/dev/".to_string(),
                                    "/tmp/.X11-unix:/tmp/.X11-unix".to_string(),
                                    format!("{}:/home/user/contiki-ng", volume_path),
                                ]),
                                network_mode: Some("host".to_string()),
                                ..Default::default()
                            }),
                            entrypoint: Some(vec!["sleep".to_string(), "infinity".to_string()]),
                            open_stdin: Some(true),
                            ..Default::default()
                        },
                    )
                    .await.context("while creating the container")?;

                // Print warnings if any
                for warning in response.warnings {
                    println!("[WARN] {warning}");
                }

                self.client
                    .start_container(constants::CONTIKER_CONTAINER, None::<StartContainerOptions>)
                    .await
                    .context("while starting the contiker container")?;

                // Query again to verify previous step
                self.query_docker()
                    .await
                    .context(
                        "we just created a contiker container but cannot connect to Docker anymore",
                    )?
                    .context(
                        "we just created a contiker container but cannot find it, please try again",
                    )?
            }
        };

        // Make sure that container is running
        if !matches!(container.state, ContainerSummaryStateEnum::RUNNING) {
            // If it is not running, start the container
            self.client
                .start_container(constants::CONTIKER_CONTAINER, None::<StartContainerOptions>)
                .await
                .context("while starting the contiker container")?;
        }

        Ok(container)
    }

    async fn exec_docker(&self, cmd: Vec<String>, user: User) -> anyhow::Result<()> {
        // Make sure container is up, and get the information
        let _ = self.ensure_contiker_up(user).await?;

        // Start command in container
        let response = self
            .client
            .create_exec(
                constants::CONTIKER_CONTAINER,
                ExecConfig {
                    attach_stdin: Some(true),
                    attach_stdout: Some(true),
                    attach_stderr: Some(true),
                    console_size: self.term_size.map(|(w, h)| vec![w as usize, h as usize]),
                    tty: Some(true),
                    cmd: Some(if cmd.is_empty() {
                        vec!["bash".to_string()]
                    } else {
                        cmd
                    }),
                    privileged: Some(true),
                    user: Some(format!("{}:{}", user.uid, user.gid)),
                    ..Default::default()
                },
            )
            .await
            .context("while preparing the shell command")?;

        let StartExecResults::Attached {
            mut output,
            mut input,
        } = self
            .client
            .start_exec(
                &response.id,
                Some(StartExecOptions {
                    detach: false,
                    tty: true,
                    ..Default::default()
                }),
            )
            .await
            .context("while starting the shell command")?
        else {
            panic!(
                "The resulting exec should always be attached as this is how it has been configured. This is a programming mistake"
            );
        };

        // Connect the command to std io
        tokio::spawn(async move {
            #[allow(clippy::unbuffered_bytes)]
            let mut stdin = termion::async_stdin().bytes();
            loop {
                if let Some(Ok(byte)) = stdin.next() {
                    input.write_all(&[byte]).await.ok();
                } else {
                    tokio::time::sleep(Duration::from_nanos(10)).await;
                }
            }
        });

        // Set `stdout` in raw mode so we can do `TTY` stuff
        let stdout = std::io::stdout();
        let mut stdout = stdout.lock().into_raw_mode()?;

        // Pipe Docker attach output into `stdout`
        while let Some(Ok(output)) = output.next().await {
            stdout.write_all(output.into_bytes().as_ref())?;
            stdout.flush()?;
        }

        Ok(())
    }

    async fn rm_contiker(&self) -> anyhow::Result<()> {
        let options = RemoveContainerOptions {
            force: true,
            ..Default::default()
        };

        self.client
            .remove_container(constants::CONTIKER_CONTAINER, Some(options))
            .await
            .context("while removing contiker container")
    }
}

#[derive(Debug, Clone)]
pub struct ContainerInfo {
    pub id: String,
    pub state: ContainerSummaryStateEnum,
    pub mounts: Vec<MountPoint>,
}

impl ContainerInfo {
    pub fn has_mount(&self, mount: &Path) -> bool {
        // When we add a volume, it will be a bind volume with the given mount point
        self.mounts.iter().any(|mount_info| {
            mount_info.typ == Some(bollard::secret::MountPointTypeEnum::BIND)
                && mount_info.source.as_ref().map(Path::new) == Some(mount)
        })
    }
}

impl TryFrom<ContainerSummary> for ContainerInfo {
    type Error = anyhow::Error;

    fn try_from(value: ContainerSummary) -> Result<Self, Self::Error> {
        Ok(Self {
            id: value.id.context("no id found")?,
            state: value.state.context("no status found")?,
            mounts: value.mounts.context("no mountpoints found")?,
        })
    }
}

#[derive(Debug, Clone, Copy)]
pub struct User {
    pub uid: u32,
    pub gid: u32,
}

impl User {
    pub fn new(uid: u32, gid: u32) -> Self {
        Self { uid, gid }
    }

    pub fn root() -> Self {
        Self { uid: 0, gid: 0 }
    }

    pub fn infer() -> Self {
        let uid = rustix::process::getuid();
        let gid = rustix::process::getgid();

        Self {
            uid: uid.as_raw(),
            gid: gid.as_raw(),
        }
    }
}

impl Default for User {
    fn default() -> Self {
        Self {
            uid: 1000,
            gid: 1000,
        }
    }
}
