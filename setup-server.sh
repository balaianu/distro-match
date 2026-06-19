#!/bin/sh
# DistroMatch SSH server setup script for Debian/Ubuntu ARM64 VMs.
#
# Run as root on a fresh VM:
#   curl -fsSL https://github.com/balaianu/distro-match/raw/master/setup-server.sh | sh
#
# Or with a specific version:
#   curl -fsSL https://github.com/balaianu/distro-match/raw/master/setup-server.sh | sh -s -- v0.1.0
#
# What it does:
#   1. Creates a dedicated 'distromatch' user
#   2. Downloads the linux-arm64 binary
#   3. Generates an SSH host key
#   4. Installs a systemd service that serves the TUI on port 22
#   5. Opens the firewall for SSH

set -eu

VERSION="${1:-latest}"
INSTALL_DIR="/opt/distromatch"
SERVICE_FILE="/etc/systemd/system/distromatch.service"
REPO="balaianu/distro-match"

# ── Must be root ──────────────────────────────────────────────────────────────

if [ "$(id -u)" -ne 0 ]; then
	echo "This script must be run as root." >&2
	exit 1
fi

# ── Resolve version ──────────────────────────────────────────────────────────

if [ "$VERSION" = "latest" ]; then
	API_URL="https://api.github.com/repos/${REPO}/releases/latest"
	if command -v curl >/dev/null 2>&1; then
		VERSION="$(curl -fsSL "$API_URL" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
	elif command -v wget >/dev/null 2>&1; then
		VERSION="$(wget -qO- "$API_URL" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')"
	else
		echo "Error: need curl or wget" >&2
		exit 1
	fi
	if [ -z "$VERSION" ]; then
		echo "Error: could not determine latest version" >&2
		exit 1
	fi
fi

echo "DistroMatch SSH server setup"
echo "  Version:     $VERSION"
echo "  Install dir: $INSTALL_DIR"
echo "  Port:        22"
echo ""

# ── Create user ──────────────────────────────────────────────────────────────

if ! id distromatch >/dev/null 2>&1; then
	echo "Creating user 'distromatch'..."
	useradd --system --no-create-home --shell /usr/sbin/nologin distromatch
else
	echo "User 'distromatch' already exists."
fi

# ── Download binary ──────────────────────────────────────────────────────────

echo "Downloading binary..."
mkdir -p "$INSTALL_DIR/bin"
mkdir -p "$INSTALL_DIR/ssh"

URL="https://github.com/${REPO}/releases/download/${VERSION}/distro-match-linux-arm64"

if command -v curl >/dev/null 2>&1; then
	curl -fsSL "$URL" -o "$INSTALL_DIR/bin/distro-match"
elif command -v wget >/dev/null 2>&1; then
	wget -qO "$INSTALL_DIR/bin/distro-match" "$URL"
fi

chmod +x "$INSTALL_DIR/bin/distro-match"
chown -R distromatch:distromatch "$INSTALL_DIR"

# ── Generate SSH host key ────────────────────────────────────────────────────

echo "Generating SSH host key..."
ssh-keygen -t ed25519 -f "$INSTALL_DIR/ssh/host_ed25519" -N "" -C "distromatch"
chown -R distromatch:distromatch "$INSTALL_DIR/ssh"

# ── Install systemd service ──────────────────────────────────────────────────

echo "Installing systemd service..."
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=DistroMatch SSH Server
After=network.target

[Service]
Type=simple
User=distromatch
Group=distromatch
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/bin/distro-match --serve :22 --host-key $INSTALL_DIR/ssh/host_ed25519
Restart=on-failure
RestartSec=5

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
PrivateTmp=yes
ReadWritePaths=$INSTALL_DIR/ssh
AmbientCapabilities=CAP_NET_BIND_SERVICE
CapabilityBoundingSet=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
EOF

# ── Stop OpenSSH if it's using port 22 ───────────────────────────────────────

if systemctl is-active --quiet sshd 2>/dev/null || systemctl is-active --quiet ssh 2>/dev/null; then
	echo "Stopping OpenSSH to free port 22..."
	systemctl stop sshd 2>/dev/null || true
	systemctl stop ssh 2>/dev/null || true
	systemctl disable sshd 2>/dev/null || true
	systemctl disable ssh 2>/dev/null || true
	echo "  OpenSSH stopped and disabled."
	echo "  WARNING: If you're connected via SSH, you may lose your session."
	echo "  Consider running the distromatch service on a different port first."
fi

# ── Open firewall ────────────────────────────────────────────────────────────

if command -v ufw >/dev/null 2>&1; then
	echo "Opening firewall for port 22..."
	ufw allow 22/tcp
elif command -v firewall-cmd >/dev/null 2>&1; then
	echo "Opening firewall for port 22..."
	firewall-cmd --permanent --add-port=22/tcp
	firewall-cmd --reload
fi

# ── Enable and start ─────────────────────────────────────────────────────────

systemctl daemon-reload
systemctl enable distromatch
systemctl start distromatch

echo ""
echo "DistroMatch SSH server is running!"
echo ""
echo "  Test:  ssh anyone@$(hostname -I 2>/dev/null | awk '{print $1}' || echo 'your-server-ip')"
echo ""
echo "  Service:  systemctl status distromatch"
echo "  Logs:     journalctl -u distromatch -f"
echo "  Stop:     systemctl stop distromatch"
