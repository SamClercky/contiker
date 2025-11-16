
$wsl_distro_name = "contiker-wsl"
$wsl_export_name = "$wsl_distro_name.wsl"
$wsl_export_filename = Join-Path -Path $PSScriptRoot -ChildPath $wsl_export_name

Write-Output "[*] Check if WSL distro has been generated"

if (-not (Test-Path -Path $wsl_export_filename)) {
    Write-Output "[*] No distro found, so generating one"
    & "$PSScriptRoot/generate.ps1"
}

Write-Output "[*] Check if there is an older version of the distro installed"
$installed_distros = wsl --list --all
if ($installed_distros -contains $wsl_distro_name) {
    # If so, remove and delete distro
    Write-Output "[*] Cleaning up old version of distro"
    wsl --unregister $wsl_distro_name
}

Write-Output "[*] Import distro into WSL"
wsl --install --name $wsl_distro_name --from-file $wsl_export_filename
