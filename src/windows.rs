use anyhow::Context;
use clap::Subcommand;
use contiker_init::GitManager;
use contiker_wsl::WslManager;

use crate::*;

pub fn handle_up() -> anyhow::Result<()> {
    let wsl = WslManager::new(None)?;
    let up = wsl
        .is_up()
        .context("while querying if the distribution is installed")?;

    if up.is_some() {
        println!("The Contiker WSL distribution is succesfully installed");
    } else {
        println!("The Contiker WSL distribution is not installed");
    }

    Ok(())
}

pub fn handle_rm() -> anyhow::Result<()> {
    let wsl = WslManager::new(None)?;
    wsl.rm()
        .context("while removing the Contiker WSL distribution")?;

    Ok(())
}

pub fn handle_exec(args: ExecArgs) -> anyhow::Result<()> {
    let wsl = WslManager::new(args.volume)?;
    wsl.exec(args.command)?;

    Ok(())
}

pub fn handle_reset() -> anyhow::Result<()> {
    let wsl = WslManager::new(None)?;
    wsl.reset()?;

    Ok(())
}

pub fn handle_init(args: InitArgs) -> anyhow::Result<()> {
    if args.windows {
        let git = GitManager::new(
            args.volume.as_deref(),
            None,
            args.branch.as_deref(),
            !args.no_shallow,
        );
        git.install()?;

        println!("Succesfully cloned repository");
        println!("To make this repository the default, add the following to your $PATH:");
        println!();
        println!("CNG_PATH={}", git.volume().display());
        println!();
        println!("To temporarly make it the default, use th following command:");
        println!();
        println!("$env:PATH=$PATH:{}", git.volume().display());
        println!();
    } else {
        let wsl = WslManager::new(args.volume)?;
        wsl.exec(vec!["/etc/oobe.sh".to_string()])?;
    }

    Ok(())
}

pub fn handle_code(args: ExecArgs) -> anyhow::Result<()> {
    let wsl = WslManager::new(args.volume)?;
    wsl.exec(vec!["code".to_string()])?;
    Ok(())
}

#[derive(Subcommand, Default, Debug, Clone)]
pub enum Fixes {
    #[default]
    All,
}

impl Fixes {
    pub fn apply(&self) -> anyhow::Result<()> {
        match self {
            Fixes::All => Ok(()),
        }
    }
}
