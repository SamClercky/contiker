use std::path::PathBuf;

use clap::{Args, Parser};

#[cfg(target_os = "windows")]
mod windows;
#[cfg(target_os = "windows")]
use windows::*;
#[cfg(target_os = "linux")]
mod linux;
#[cfg(target_os = "linux")]
use linux::*;

#[derive(Parser, Debug, Clone)]
/// Manage the environment for Contiki-NG
///
/// On Linux, Docker is used
///
/// On Windows, WSL is used
enum Cli {
    /// Startup the enviroment
    Up,
    /// Remove the environment
    Rm,
    /// Initialize Contiki-NG
    Init(InitArgs),
    /// Auto-fix some common issues
    Fix,
    /// Execute a command directly
    Exec(ExecArgs),
    /// Reset the environment
    Reset,
    /// Start Cooja (shorthand for `contiker exec cooja`)
    Cooja(ExecArgs),
}

#[derive(Args, Debug, Clone)]
struct InitArgs {
    /// Optional custom git url
    #[arg(long)]
    pub git: Option<String>,
    /// Place to put Contiki folder
    #[arg(short, long)]
    pub volume: Option<PathBuf>,
    /// Do not perform a shallow clone (worse performance)
    #[arg(long)]
    pub no_shallow: bool,
    #[arg(short, long)]
    /// Optionally specify a branch
    pub branch: Option<String>,
}

#[derive(Args, Debug, Clone)]
struct ExecArgs {
    /// The command to pass into the Contiki environment
    pub command: Vec<String>,
    #[arg(long)]
    /// alias for `uid 0`
    pub root: bool,
    #[arg(short, long)]
    /// Set the uid
    pub uid: Option<u32>,
    #[arg(short, long)]
    /// Set the gid
    pub gid: Option<u32>,
    #[arg(short, long)]
    /// Set the mount volume explicitly and asks if this is different from the last time
    pub volume: Option<PathBuf>,
}

fn main() {
    let cli = Cli::parse();

    let result = match cli {
        Cli::Up => handle_up(),
        Cli::Rm => handle_rm(),
        Cli::Init(init_args) => handle_init(init_args),
        Cli::Fix => todo!(),
        Cli::Exec(exec_args) => handle_exec(exec_args),
        Cli::Reset => handle_reset(),
        Cli::Cooja(exec_args) => handle_exec(ExecArgs {
            command: vec!["cooja".to_string()],
            ..exec_args
        }),
    };

    if let Err(err) = result {
        eprintln!("[ERROR] {err:?}");
    }
}
