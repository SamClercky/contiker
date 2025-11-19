use contiker_docker::{DockerManager, User};

use crate::ExecArgs;

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
    let docker = DockerManager::new(None)?;
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
