$ErrorActionPreference = "Stop"

$Repo = "KingHanzala/gcpsec"
$BinaryName = "gcpsec.exe"
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { Join-Path $env:USERPROFILE "AppData\Local\Microsoft\WindowsApps" }

$arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
switch ($arch) {
    "x64" { $GoArch = "amd64" }
    "arm64" { $GoArch = "arm64" }
    default { throw "Unsupported Windows architecture: $arch" }
}

$latestUrl = "https://github.com/$Repo/releases/latest/download"
$archiveName = "gcpsec_windows_${GoArch}.zip"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("gcpsec-install-" + [System.Guid]::NewGuid().ToString("N"))
$archivePath = Join-Path $tempDir $archiveName

New-Item -ItemType Directory -Force -Path $tempDir | Out-Null
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

try {
    Write-Host "Downloading $archiveName..."
    Invoke-WebRequest -Uri "$latestUrl/$archiveName" -OutFile $archivePath

    Expand-Archive -Path $archivePath -DestinationPath $tempDir -Force

    $targetPath = Join-Path $InstallDir $BinaryName
    Copy-Item -Path (Join-Path $tempDir $BinaryName) -Destination $targetPath -Force

    Write-Host "Installed $BinaryName to $targetPath"

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (($userPath -split ";") -notcontains $InstallDir) {
        Write-Host ""
        Write-Host "Add this folder to your User PATH if needed:"
        Write-Host $InstallDir
    }

    Write-Host "Run: gcpsec version"
}
finally {
    if (Test-Path $tempDir) {
        Remove-Item -Recurse -Force $tempDir
    }
}
