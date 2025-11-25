use anyhow::{Context, bail};
use std::path::{Path, PathBuf};

pub struct GitManager {
    volume: PathBuf,
    url: String,
    branch: Option<String>,
    shallow: bool,
}

impl GitManager {
    pub fn new(
        volume: Option<&Path>,
        url: Option<&str>,
        branch: Option<&str>,
        shallow: bool,
    ) -> Self {
        Self {
            volume: volume.unwrap_or(Path::new(".")).to_path_buf(),
            url: url
                .unwrap_or("https://github.com/contiki-ng/contiki-ng.git")
                .to_string(),
            branch: branch.map(|b| b.to_string()),
            shallow,
        }
    }

    pub fn volume(&self) -> PathBuf {
        std::path::absolute(self.volume.as_path()).unwrap()
    }

    pub fn install(&self) -> anyhow::Result<()> {
        let mut command = std::process::Command::new("git");
        command
            .arg("clone")
            .arg("--recurse-submodules")
            .arg(format!("-j{}", num_cpus::get()));

        if self.shallow {
            command.arg("--shallow-submodules");
            command.arg("--depth");
            command.arg("1");
        }

        if let Some(branch) = &self.branch {
            command.arg("--branch");
            command.arg(branch);
        }

        command.arg(&self.url).arg(&self.volume);

        let command_status = command.status().context("while cloning repository")?;

        if !command_status.success() {
            bail!("[ERROR] Unsuccesful status code after clone");
        }

        Ok(())
    }
}
