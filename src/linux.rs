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
    println!("export CNG_PATH={:?}", git.volume());
    println!();

    Ok(())
}
