#!/bin/bash

set -e

SCRIPT_DIR="$(dirname -- "$0")"
echo "SCRIPT_DIR: $SCRIPT_DIR"

echo "[*] Building Contiker WSL distro"

container_export_name="contiker_wsl_export"
running_container_name="contiker_wsl_export"
wsl_export_name="contiker-wsl"

echo "[*] Prepare extraction"

docker build -t $container_export_name "$SCRIPT_DIR"

# Check if container is already up
up_containers=$(docker ps -a --format "{{.Names}}")
if [[ $up_containers == *$running_container_name* ]]; then
  # If so, remove and delete container
  docker container rm -f $running_container_name
fi

# Make sure that the container exists
docker run -t --name $running_container_name $container_export_name echo "From the contiker container: Hello world"

# Export into a tar
echo "[*] Extract Contiker distro"
echo "This may take a while..."
docker export "$running_container_name" | gzip >"$SCRIPT_DIR/$wsl_export_name.wsl"

# Remove container
echo "[*] Cleaning up resources"
docker container rm -f $running_container_name
