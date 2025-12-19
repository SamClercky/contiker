use std::{path::Path, process::Command};

use anyhow::{Context, bail};
use clap::Subcommand;
use contiker_docker::{DockerManager, User};
use contiker_init::GitManager;

use crate::{ExecArgs, InitArgs};

pub fn handle_up() -> anyhow::Result<()> {
    let docker = DockerManager::new(None)?;
    if docker.is_up()? {
        println!("Contiker is up");
    } else {
        println!("Contiker is down");
    }

    Ok(())
}

pub fn handle_rm() -> anyhow::Result<()> {
    let docker = DockerManager::new(None)?;
    docker.rm()?;

    Ok(())
}

pub fn handle_exec(args: ExecArgs) -> anyhow::Result<()> {
    let docker = DockerManager::new(args.volume)?;
    docker.exec(args.command, {
        let mut user = User::infer();
        if let Some(uid) = args.uid {
            user.uid = uid;
        }
        if let Some(gid) = args.gid {
            user.gid = gid;
        }

        if args.root {
            user = User::root();
        }

        user
    })?;

    Ok(())
}

pub fn handle_reset() -> anyhow::Result<()> {
    let docker = DockerManager::new(None)?;
    docker.reset()
}

pub fn handle_init(args: InitArgs) -> anyhow::Result<()> {
    let git = GitManager::new(
        args.volume.as_deref(),
        None,
        args.branch.as_deref(),
        !args.no_shallow,
    );
    git.install()?;

    println!("Succesfully cloned repository");
    println!("To make this repository the default, put the following in your .bashrc:");
    println!();
    println!("export CNG_PATH={}", git.volume().display());
    println!();

    Ok(())
}

pub fn handle_code(args: ExecArgs) -> anyhow::Result<()> {
    let command_status = std::process::Command::new("code")
        .arg(args.volume.as_deref().unwrap_or(Path::new(".")))
        .status()
        .context("while launching vscode")?;

    if !command_status.success() {
        bail!("Starting VSCode was unsuccesful");
    }

    Ok(())
}

#[derive(Subcommand, Default, Debug, Clone)]
pub enum Fixes {
    #[default]
    All,
    /// Fix issue related with xhost permissions
    XHost,
    /// Fix user not in the docker user group
    DockerPerm,
    /// Fix all files are owned by root
    FilePerm,
}

impl Fixes {
    pub fn apply(&self) -> anyhow::Result<()> {
        match self {
            Fixes::All => {
                self.fix_xhost()?;
                self.fix_dockerperm()?;
                self.fix_fileperm()?;

                Ok(())
            }
            Fixes::XHost => self.fix_xhost(),
            Fixes::DockerPerm => self.fix_dockerperm(),
            Fixes::FilePerm => self.fix_fileperm(),
        }
    }

    fn fix_xhost(&self) -> anyhow::Result<()> {
        println!("Applying xhost fix");

        let command_status = Command::new("xhost")
            .arg("+local:docker")
            .status()
            .context("while running the xhost command")?;

        if !command_status.success() {
            bail!("non-successful error code by xhost command");
        }

        Ok(())
    }

    fn fix_dockerperm(&self) -> anyhow::Result<()> {
        println!("Applying Docker permissions fix");

        let command_status = Command::new("sudo")
            .arg("usermod")
            .arg("-aG")
            .arg("docker")
            .arg(User::infer().name()?)
            .status()
            .context("while running the usermod command")?;

        if !command_status.success() {
            bail!("non-successful error code by xhost command");
        }

        Ok(())
    }

    fn fix_fileperm(&self) -> anyhow::Result<()> {
        println!("Applying file permissions fix");

        let user = User::infer();
        handle_exec(ExecArgs {
            command: vec![
                "chown".to_string(),
                "-R".to_string(),
                format!("{}:{}", user.uid, user.gid),
                "/home/user/contiki-ng".to_string(),
            ],
            root: false,
            uid: Some(user.uid),
            gid: Some(user.gid),
            volume: None,
        })
        .context("while applying file permissions fix")?;

        Ok(())
    }
}
