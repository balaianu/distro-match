#!/bin/sh
# DistroMatch installer — detects OS/arch and downloads the correct binary.
#
# Usage:
#   curl -fsSL https://distromatch.dev/install.sh | sh
#   curl -fsSL https://distromatch.dev/install.sh | sh -s -- --version v0.1.0
#   curl -fsSL https://distromatch.dev/install.sh | sh -s -- --dir /usr/local/bin
#   curl -fsSL https://distromatch.dev/install.sh | sh -s -- --dir ~/.local/bin
#
# Flags:
#   --version <tag>   Version to download (default: latest)
#   --dir <path>      Install directory (default: ~/.local/bin)
#   --dry-run         Print what would happen without downloading

set -eu

# ── Defaults ─────────────────────────────────────────────────────────────────

REPO="balaianu/distro-match"
VERSION="latest"
INSTALL_DIR="${DISTROMATCH_INSTALL_DIR:-$HOME/.local/bin}"
BASE_URL="${DISTROMATCH_BASE_URL:-https://github.com/${REPO}/releases/download}"
DRY_RUN=false
BINARY_NAME="distromatch"

# ── Argument parsing ─────────────────────────────────────────────────────────

while [ $# -gt 0 ]; do
	case "$1" in
		--version)
			VERSION="$2"
			shift 2
			;;
		--dir)
			INSTALL_DIR="$2"
			shift 2
			;;
		--dry-run)
			DRY_RUN=true
			shift
			;;
		--help|-h)
			cat <<EOF
DistroMatch installer

Usage: install.sh [flags]

Flags:
  --version <tag>   Version to download (default: latest)
  --dir <path>      Install directory (default: ~/.local/bin)
  --dry-run         Print what would happen without downloading
  --help            Show this help
EOF
			exit 0
			;;
		*)
			echo "Unknown flag: $1" >&2
			exit 1
			;;
	esac
done

# ── OS/arch detection ────────────────────────────────────────────────────────

OS="$(uname -s 2>/dev/null || echo unknown)"
ARCH="$(uname -m 2>/dev/null || echo unknown)"

case "$OS" in
	Linux*)  GOOS="linux"   ;;
	Darwin*) GOOS="darwin"  ;;
	*) echo "Unsupported OS: $OS (expected Linux or macOS)" >&2; exit 1 ;;
esac

case "$ARCH" in
	x86_64|amd64)  GOARCH="amd64" ;;
	aarch64|arm64) GOARCH="arm64" ;;
	*) echo "Unsupported architecture: $ARCH (expected amd64 or arm64)" >&2; exit 1 ;;
esac

# ── Resolve version ──────────────────────────────────────────────────────────

if [ "$VERSION" = "latest" ]; then
	if [ "$DRY_RUN" = true ]; then
		VERSION="(latest — would resolve from GitHub)"
	else
		# Query the GitHub API for the latest release tag.
		API_URL="https://api.github.com/repos/${REPO}/releases/latest"
		if command -v curl >/dev/null 2>&1; then
			VERSION="$(curl -fsSL "$API_URL" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
		elif command -v wget >/dev/null 2>&1; then
			VERSION="$(wget -qO- "$API_URL" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
		else
			echo "Error: need curl or wget to resolve latest version" >&2
			exit 1
		fi

		if [ -z "$VERSION" ]; then
			echo "Error: could not determine latest version from GitHub" >&2
			exit 1
		fi
	fi
fi

# ── Construct download URL ───────────────────────────────────────────────────

ASSET="distro-match-${GOOS}-${GOARCH}"
URL="${BASE_URL}/${VERSION}/${ASSET}"

# ── Print plan ───────────────────────────────────────────────────────────────

cat <<EOF
DistroMatch installer
  OS:       ${GOOS} (${ARCH})
  Version:  ${VERSION}
  Binary:   ${ASSET}
  URL:      ${URL}
  Install:  ${INSTALL_DIR}/${BINARY_NAME}
EOF

if [ "$DRY_RUN" = true ]; then
	echo "Dry run — not downloading."
	exit 0
fi

# ── Download ─────────────────────────────────────────────────────────────────

TMPFILE="$(mktemp -t distromatch.XXXXXX)"
trap 'rm -f "$TMPFILE"' EXIT

echo "Downloading..."
if command -v curl >/dev/null 2>&1; then
	curl -fsSL "$URL" -o "$TMPFILE"
elif command -v wget >/dev/null 2>&1; then
	wget -qO "$TMPFILE" "$URL"
else
	echo "Error: need curl or wget to download" >&2
	exit 1
fi

chmod +x "$TMPFILE"

# ── Install ──────────────────────────────────────────────────────────────────

mkdir -p "$INSTALL_DIR"

# Try sudo only if the target dir isn't writable.
if [ -w "$INSTALL_DIR" ]; then
	mv "$TMPFILE" "${INSTALL_DIR}/${BINARY_NAME}"
else
	echo "Install directory requires sudo: ${INSTALL_DIR}"
	sudo mv "$TMPFILE" "${INSTALL_DIR}/${BINARY_NAME}"
fi

chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

# ── Verify ───────────────────────────────────────────────────────────────────

echo ""
echo "Installed to: ${INSTALL_DIR}/${BINARY_NAME}"

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
	echo ""
	echo "Warning: ${INSTALL_DIR} is not in your PATH."
	echo "Add it by running:"
	echo ""
	case "$(basename "$SHELL")" in
		fish)
			echo "  fish_add_path ${INSTALL_DIR}"
			;;
		*)
			echo "  echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.${SHELL##*/}rc"
			;;
	esac
fi

echo ""
echo "Run 'distromatch' to start, or 'distromatch --help' for options."
