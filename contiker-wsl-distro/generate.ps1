Write-Output "[*] Building Contiker WSL distro"

$prevPwd = $PWD; Set-Location -ErrorAction Stop -LiteralPath $PSScriptRoot
$container_export_name = "contiker_wsl_export"
$running_container_name = "contiker_wsl_export"
$wsl_export_name = "contiker-wsl"

try {
    Write-Output "[*] Prepare extraction"

    docker build -t $container_export_name .

    # Check if container is already up
    $up_containers = docker ps -a --format "{{.Names}}"
    if ($up_containers -contains $running_container_name) {
        # If so, remove and delete container
        docker container rm -f $running_container_name
    }

    # Make sure that the container exists
    docker run -t --name $running_container_name $container_export_name echo "From the contiker container: Hello world"

    # Export into a tar
    Write-Output "[*] Extract Contiker distro"
    Write-Output "This may take a while..."
    docker export "$running_container_name" | gzip > "$wsl_export_name.wsl"

    # Remove container
    Write-Output "[*] Cleaning up resources"
    docker container rm -f $running_container_name
}
finally {
    Set-Location $prevPwd # Return to old PWD
}
