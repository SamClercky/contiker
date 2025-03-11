# Contiker

New and improved contiker script to work with variable workspace directories and docker.

## Prerequisites

- Go
- Make

## Install

```bash
make install
```

To allow X11 to work with Cooja on Linux and Docker, use the following command:
```bash
xhost +local:docker
```

## Run

```bash
contiker -v . # Run with . as mounted dir
contiker # Run with contents of $CNG_PATH as mounted dir
contiker -sh # Open a shell in already running container
contiker -rm # Remove current container
```
