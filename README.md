# Contiker

New and improved contiker script to work with variable workspace directories and docker.

## Prerequisites

- Git
- Go
- Make
- Docker
- xhost (used to fix permissions)

## Install

```bash
git clone https://github.com/SamClercky/contiker.git
cd contiker
make install
```

To allow X11 to work with Cooja on Linux and Docker, use the following command:
```bash
xhost +local:docker # or
contiker fix -xhost # To fix through contiker
```

## Run

```bash
# Init Contiki
contiker init # Install Contiki in the current directory
contiker init -git https://github.com/contiki-ng/contiki-ng.git # With specific git
contiker init -v contiki-ng # In specific folder

# Start Contiki container
contiker # Run with contents of $CNG_PATH as mounted dir
contiker -v . # Run with . as mounted dir
contiker -e bash # Run specific command
contiker cooja # Alias for `contiker -e cooja`
contiker -root # Create a root shell

# Management of Contiker containers
contiker rm # Remove current container
contiker up # Check if currently a container is up

# Common fixes
contiker fix -docker # Add current user the `docker` group
contiker fix -xhost # Fix common xhost issue with cooja
```
