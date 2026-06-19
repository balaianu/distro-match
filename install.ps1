# DistroMatch installer for Windows — detects arch and downloads the correct binary.
#
# Usage (PowerShell):
#   irm https://distromatch.dev/install.ps1 | iex
#   & ([scriptblock]::Create((irm https://distromatch.dev/install.ps1))) -Version v0.1.0
#
# Flags:
#   -Version <tag>   Version to download (default: latest)
#   -InstallDir <path>  Install directory (default: $env:LOCALAPPDATA\Programs\distromatch)
#   -DryRun          Print what would happen without downloading

param(
	[string]$Version = "latest",
	[string]$InstallDir = "$env:LOCALAPPDATA\Programs\distromatch",
	[switch]$DryRun
)

$ErrorActionPreference = "Stop"

# ── Defaults ─────────────────────────────────────────────────────────────────

$Repo = "balaianu/distro-match"
$BinaryName = "distromatch.exe"

# ── Arch detection ───────────────────────────────────────────────────────────

$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
	"AMD64"     { $GOARCH = "amd64" }
	"ARM64"     { $GOARCH = "arm64" }
	default     { Write-Error "Unsupported architecture: $Arch (expected AMD64 or ARM64)"; exit 1 }
}

# ── Resolve version ──────────────────────────────────────────────────────────

if ($Version -eq "latest") {
	$ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
	try {
		$Release = Invoke-RestMethod -Uri $ApiUrl -Headers @{ "User-Agent" = "distromatch-installer" }
		$Version = $Release.tag_name
	} catch {
		Write-Error "Could not determine latest version from GitHub: $_"
		exit 1
	}
}

if ([string]::IsNullOrEmpty($Version)) {
	Write-Error "Could not determine version to download"
	exit 1
}

# ── Construct download URL ───────────────────────────────────────────────────

$Asset = "distro-match-windows-$GOARCH.exe"
$Url = "https://github.com/$Repo/releases/download/$Version/$Asset"
$DestPath = Join-Path $InstallDir $BinaryName

# ── Print plan ───────────────────────────────────────────────────────────────

Write-Host "DistroMatch installer"
Write-Host "  OS:       windows ($GOARCH)"
Write-Host "  Version:  $Version"
Write-Host "  Binary:   $Asset"
Write-Host "  URL:      $Url"
Write-Host "  Install:  $DestPath"

if ($DryRun) {
	Write-Host "Dry run — not downloading."
	exit 0
}

# ── Download ─────────────────────────────────────────────────────────────────

Write-Host "Downloading..."

if (-not (Test-Path $InstallDir)) {
	New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

try {
	Invoke-WebRequest -Uri $Url -OutFile $DestPath -UseBasicParsing
} catch {
	Write-Error "Download failed: $_"
	exit 1
}

# ── Verify ───────────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "Installed to: $DestPath"

# Check if install dir is in PATH
$PathDirs = $env:PATH -split ";"
if ($PathDirs -notcontains $InstallDir) {
	Write-Host ""
	Write-Host "Warning: $InstallDir is not in your PATH."
	Write-Host "Add it by running (as admin):"
	Write-Host ""
	Write-Host "  [Environment]::SetEnvironmentVariable('PATH', `"$InstallDir;`$([Environment]::GetEnvironmentVariable('PATH', 'User'))`", 'User')"
}

Write-Host ""
Write-Host "Run 'distromatch' to start, or 'distromatch --help' for options."
