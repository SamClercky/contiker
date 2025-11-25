# Contiker

New and improved contiker script to work with variable workspace directories and docker.

## Prerequisites

### All supported operating systems

- Git
- Rust
- Wireshark (not necessary, but necessary if want to work with `tunslip`)

### Additionally for Linux

- Docker
- xhost (used to fix permissions)

### Additionally for Windows

- usbipd (allow WSL to connect to an USB device)
- WSL2 (requires Windows 11 with gWSL support, should be standard for at least the past 4 years)
- Docker (optional, if you want to build the distribution yourself)

## Install

```bash
cargo install --git https://github.com/SamClercky/contiker.git --tag v0.2.0 --locked
```

### Additionally for Linux

To allow X11 to work with Cooja on Linux and Docker, use the following command:
```bash
xhost +local:docker # or
contiker fix xhost # To fix through contiker
```

## Run

```bash
# Init Contiki
contiker init # Install Contiki in the current directory
contiker init --git https://github.com/contiki-ng/contiki-ng.git # With specific git (Linux only)
contiker init -v contiki-ng # In specific folder (Linux only)

# Start Contiki container
contiker exec -v . # Run with . as mounted dir
contiker exec bash # Run specific command
contiker cooja # Alias for `contiker exec cooja`
contiker exec --root # Create a root shell
contiker reset # Stop already running Contiker instances and redownload container/WSL distro
contiker code # Run code in the current directory

# Management of Contiker containers
contiker rm # Remove current container/WSL
contiker up # Check if currently a container/WSL is up

# Common fixes (Linux)
contiker fix all # Run all available fixes
contiker fix docker # Add current user the `docker` group
contiker fix xhost # Fix common xhost issue with cooja
contiker fix fileperm # Fix common file permission errors in current Contiki instance
```

## How it works 

The contiker application is mostly an automation to get students quickly up and running with Contiki-NG.
When there are some issues, it is important to know what is actually happening.

### Linux

On Linux, contiker uses Docker to orchestrate a `contiker` container.
This container gets configured to work on a specific installation of contiki-ng.

This makes it easy to provide some of the shortcuts like the `contiker cooja` command.
Editing is done the same way as you would when modifying any project on Linux.
Only compiling and executing some of the commands is done through the container.

The container is also put entirely open such that it becomes possible to use tools like `tunslip` in a container and Wireshark outside.
In the case of `cooja`, GUI applications get passed through as long as X11 passthrough is a possibility (most if not all Linux environments).

### Windows

On Windows, it is technically possible to use Docker and make `cooja` work.
The issue is that this will incur some high memory/CPU usage as you have certain context switches that need to be made (Windows <-> WSL <-> Docker).

Since WSL can forward GUI applications, it thus becomes interesting to use WSL directly.
And this is what contiker does.

It provides a new `contiker-wsl` distribution that is derived from the same Docker container as is used on the Linux side, but with some alterations to make it compatible with Windows.
The main task that contiker does in this case is downloading the latest version from GitHub and installing it on the user's computer.

Due to the context switch when interacting with NTFS <-> EXT4, it is not recommended to install `contiki-ng` on a Windows-accessible drive.
This is why the default installation will clone the necessary git repositories in WSL only.

When you want to develop code in WSL, you can use VS Code to make that connection through the remote development plugin.

Due to some limitations of Windows and PowerShell, it is currently not possible to have complex commands executed directly in WSL through contiker.
So, if you want to have complex commands, first enter a shell and then type your command.

## How to replicate this tool?

Cloning the contiki-ng repository:
```bash
git clone --recurse-submodules --shallow-submodules --depth 1 https://github.com/contiki-ng/contiki-ng.git ~/contiki-ng
```

### Linux

Starting the container in the current directory
```bash
docker run -it \ # Make it interactive
    -d \ # Detach immediatly
    -v $PATH:/home/user/contiki-ng \ # Connect the current working directory
    -v /dev/:/dev/ \ # Bind USB devices
    -v /tmp/.X11-unix:/tmp/.X11-unix \ # GUI on Wayland
    -e DISPLAY=$DISPLAY \ # GUI for X11
    -e _JAVA_AWT_WM_NONREPARENTING=1 \ # X11 and tiling window managers
    -e "JDK_JAVA_OPTIONS='-Dawt.useSystemAAFontSettings=on -Dswing.aatext=true -Dswing.defaultlaf=com.sun.java.swing.plaf.gtk.GTKLookAndFeel -Dsun.java2d.opengl=true'" \ # Better font rendering
    -e LOCAL_UID=1000 \ # set uid
    -e LOCAL_GID=1000 \ # set gid
    --network host \ # Enable WireShark and tunslip
    --privileged \ # Allow access to USB devices
    --name contiker \ # Set a name for easy reference later
    contiker/contiki-ng \ # Use this image
    sleep infinity # Keep the container alive
```

Then to run a command:
```bash
docker exec -it \ # Make it interactive
    contiker \ # Name of the container
    bash # The command we want to run
```

Remove the container:
```bash
docker rm -f contiker
```

Perform the `xhost` fix:
```bash
xhost +local:docker
```

Perform the Docker permissions fix:
```bash
sudo usermod -aG $(whoami) docker
```
Reboot after executing this command for it to take effect.

Some students have issues with accidentally building with root, the `fileperm` fix fixes this.
```bash
docker exec -it contiker chown -R 1000:1000 /home/user/contiki-ng
```

### Windows

Building the WSL distribution:
```bash
. ".\\contiker-wsl-distro\\generate.ps1"
```

Installing the distribution:
1. Download the latest `contiker-wsl.wsl` artifact from this repository
2. Run
```bash
wsl --install --name contiker-wsl --from-file <FILE>
```

Execute a command inside WSL:
```bash
wsl -d contiker-wsl -- bash
```

Remove the `contiker-wsl` distribution:
```bash
wsl --unregister contiker-wsl
```
**Note:** All files in this instance will also be deleted!

Run the default installation inside the container:
```bash
wsl -d contiker-wsl -- /etc/oobe.sh
```

Run `cooja` after installing contiki-ng through `/etc/oobe.sh`:
```bash
contiker exec
cd ~/contiki-ng/tools/cooja
./gradlew run
```
