use clap::{Args, Parser};

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
}

#[derive(Args, Debug, Clone)]
struct InitArgs {
    /// Optional custom git url
    #[arg(long)]
    pub git: Option<String>,
    /// Placo to put Contiki folder
    #[arg(short, long)]
    pub volume: Option<String>,
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
    pub uid: Option<u8>,
    #[arg(short, long)]
    /// Set the gid
    pub gid: Option<u8>,
}

fn main() {
    let cli = Cli::parse();
    println!("{cli:?}");
}
