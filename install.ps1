# Seshmark Windows Installer

$ErrorActionPreference = "Stop"

$installDir = "$env:LOCALAPPDATA\seshmark"
$binDir = "$installDir\bin"

# Detect architecture
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "arm64" }

$binary = "seshmark-windows-$arch.exe"
$url = "https://github.com/seshmark/seshmark/releases/latest/download/$binary"

Write-Host "Installing seshmark..."
New-Item -ItemType Directory -Force -Path $binDir | Out-Null

# Download
Invoke-WebRequest -Uri $url -OutFile "$binDir\seshmark.exe"

# Create aliases
$aliases = @("agentblame.exe", "aiblame.exe", "git-agentblame.exe", "git-aiblame.exe")
foreach ($alias in $aliases) {
    $target = "$binDir\seshmark.exe"
    $link = "$binDir\$alias"
    if (Test-Path $link) { Remove-Item $link }
    cmd /c mklink $link $target | Out-Null
}

# Add to PATH
$path = [Environment]::GetEnvironmentVariable("Path", "User")
if ($path -notlike "*$binDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$path;$binDir", "User")
    Write-Host "Added $binDir to PATH. Restart your terminal."
}

# Install global hook
& "$binDir\seshmark.exe" hook install --global

# Install in current repo if inside one
if (git rev-parse --git-dir 2>$null) {
    & "$binDir\seshmark.exe" hook install
}

Write-Host "seshmark installed successfully!"
Write-Host ""
Write-Host "Try it now:"
Write-Host "  git agentblame <file>"
Write-Host "  seshmark query --agent cursor"
