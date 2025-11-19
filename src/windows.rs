use std::process::Command;

use anyhow::Context;
use clap::Subcommand;
use contiker_init::GitManager;

pub fn handle_up() -> anyhow::Result<()> {
    todo!()
}

pub fn handle_rm() -> anyhow::Result<()> {
    todo!()
}

pub fn handle_exec(args: ExecArgs) -> anyhow::Result<()> {
    todo!()
}

pub fn handle_reset() -> anyhow::Result<()> {
    todo!()
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
    println!("To make this repository the default, add the following to your $PATH:");
    println!();
    println!("CNG_PATH={:?}", git.volume());
    println!();
    println!("To temporarly make it the default, use th following command:");
    println!();
    println!("$env:PATH=$PATH:{}", git.volume());
    println!();

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
