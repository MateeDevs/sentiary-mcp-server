$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$repository = "MateeDevs/sentiary-mcp-server"
$binary = "sentiary-cli"

$processorArchitecture = if ($env:PROCESSOR_ARCHITEW6432) {
    $env:PROCESSOR_ARCHITEW6432
} else {
    $env:PROCESSOR_ARCHITECTURE
}

$architecture = switch ($processorArchitecture.ToUpperInvariant()) {
    "AMD64" { "amd64" }
    "X86_64" { "amd64" }
    "ARM64" { "arm64" }
    default { throw "Unsupported architecture: $processorArchitecture" }
}

$asset = "$binary-windows-$architecture.exe"
$url = "https://github.com/$repository/releases/latest/download/$asset"
$installDirectory = if ($env:INSTALL_DIR) {
    $env:INSTALL_DIR
} else {
    Join-Path $env:LOCALAPPDATA "Programs\Sentiary\bin"
}
$installPath = Join-Path $installDirectory "$binary.exe"
$temporaryPath = Join-Path ([IO.Path]::GetTempPath()) "$([Guid]::NewGuid()).exe"

try {
    New-Item -ItemType Directory -Force -Path $installDirectory | Out-Null
    Invoke-WebRequest -UseBasicParsing -Uri $url -OutFile $temporaryPath
    Move-Item -Force -Path $temporaryPath -Destination $installPath
} finally {
    if (Test-Path $temporaryPath) {
        Remove-Item -Force $temporaryPath
    }
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$pathEntries = @($userPath -split ";" | Where-Object { $_ })
$hasInstallDirectory = $pathEntries | Where-Object {
    $_.TrimEnd("\") -ieq $installDirectory.TrimEnd("\")
}

if (-not $hasInstallDirectory) {
    $newUserPath = (@($pathEntries) + $installDirectory) -join ";"
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
}

$currentPathEntries = @($env:Path -split ";")
if ($currentPathEntries -notcontains $installDirectory) {
    $env:Path = "$installDirectory;$env:Path"
}

$version = & $installPath --version
Write-Output "Installed $binary to $installPath"
if ($version) {
    Write-Output "Version: $version"
}
