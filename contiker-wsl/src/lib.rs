use std::path::PathBuf;

use anyhow::anyhow;
use wsl_api::Wsl2;

pub struct WslManager {
    client: Wsl2,
    volume_path: Option<PathBuf>,
}

impl WslManager {
    pub fn new(volume_path: Option<PathBuf>) -> anyhow::Result<Self> {
        Ok(Self {
            client: Wsl2::new()
                .map_err(|err| anyhow!("while creating a WSL2 connection: {err}"))?,
            volume_path,
        })
    }

    pub fn rm(&self) -> anyhow::Result<()> {
        todo!()
    }

    pub fn exec(&self, cmd: Vec<String>) -> anyhow::Result<()> {
        todo!()
    }

    pub fn is_up(&self) -> anyhow::Result<bool> {
        todo!()
    }

    pub fn reset(&self) -> anyhow::Result<()> {
        todo!()
    }
}
