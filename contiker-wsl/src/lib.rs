use std::{path::PathBuf, process::Command};

use anyhow::{Context, anyhow, bail};
use uuid::Uuid;
use wsl_api::Wsl2;

mod constants;
mod gh_release;

pub struct WslManager {
    client: Wsl2,
    _volume_path: Option<PathBuf>,
}

impl WslManager {
    pub fn new(volume_path: Option<PathBuf>) -> anyhow::Result<Self> {
        Ok(Self {
            client: Wsl2::new()
                .map_err(|err| anyhow!("while creating a WSL2 connection: {err}"))?,
            _volume_path: volume_path,
        })
    }

    pub fn rm(&self) -> anyhow::Result<()> {
        let command_status = Command::new("wsl")
            .arg("--unregister")
            .arg("contiker-wsl")
            .status()
            .context("while removing contiker")?;

        if !command_status.success() {
            bail!("The command to remove contiker did not finish succesfully");
        }

        Ok(())
    }

    pub fn exec(&self, cmd: Vec<String>) -> anyhow::Result<()> {
        self.ensure_contiker_installed(true)?;

        let status = Command::new("wsl")
            .arg("-d")
            .arg(constants::CONTIKER_WSL_NAME)
            .arg("--")
            .arg(if !cmd.is_empty() {
                cmd.join(" ")
            } else {
                "bash".to_string()
            })
            .status()
            .context("while executing command")?;

        if !status.success() {
            eprintln!("Command did not finish succesfully");
        }

        Ok(())
    }

    pub fn is_up(&self) -> anyhow::Result<Option<Uuid>> {
        let distros = self
            .client
            .enumerate_distributions()
            .map_err(|err| anyhow!("while enumerating installed distributions: {err}"))?;

        Ok(distros
            .into_iter()
            .find(|dis| dis.name == constants::CONTIKER_WSL_NAME)
            .map(|distro| distro.uuid))
    }

    pub fn reset(&self) -> anyhow::Result<()> {
        self.ensure_contiker_installed(false)?;

        Ok(())
    }

    fn ensure_contiker_installed(&self, use_cache: bool) -> anyhow::Result<()> {
        let is_up = self.is_up()?;
        if is_up.is_none() || !use_cache {
            // Install Contiker
            self.install_wsl_distro()
                .context("while installing the WSL Contiker distribution")?;
        }

        Ok(())
    }

    fn install_wsl_distro(&self) -> anyhow::Result<()> {
        let client = reqwest::blocking::Client::new();

        println!("[*] Querying for the latest Contiker WSL distribution release");

        let release = gh_release::get_latest_release(&client)
            .context("while retrieving the latest GH release")?;

        let Some(contiker_distro) = release
            .assets
            .into_iter()
            .find(|asset| asset.name == constants::CONTIKER_WSL_FILENAME)
        else {
            bail!(
                "Could not find the latest Contiker WSL distro. This may be a temporary issue, try again later."
            );
        };

        println!("[*] Downloading the latest Contiker WSL distribution release.");
        println!("This may take a while...");

        // Create a temporary directory to store distribution into
        let out_dir = tempdir::TempDir::new("contiki-wsl-distro")
            .context("could not create a temporary directory")?;
        let wsl_file = contiker_distro
            .download(out_dir.path(), &client)
            .context("while downloading WSL Contiker Distro")?;

        let is_up = self.is_up()?;
        if is_up.is_some() {
            // Unregister the distribution as you can only have one version up
            println!("[*] Detected an already existing contiker distribution, now removing");
            self.rm()?;
        }

        println!("[*] Registering the WSL distribution");
        let command_status = Command::new("wsl")
            .arg("--install")
            .arg("--name")
            .arg(constants::CONTIKER_WSL_NAME)
            .arg("--from-file")
            .arg(&wsl_file)
            .status()
            .context("while registering the WSL Contiker distribution")?;

        if !command_status.success() {
            bail!("WSL did not finish succesfully when installing the new distribution");
        }

        Ok(())
    }
}
