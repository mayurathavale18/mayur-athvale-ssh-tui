#!/bin/bash
set -e

# SSH Portfolio Server — EC2 Setup Script
#
# Usage:
#   Step 1: sudo bash setup.sh --move-ssh-port
#   Step 2: sudo bash setup.sh /path/to/binary

INSTALL_DIR="/opt/ssh-portfolio"
SERVICE_USER="portfolio"

if [ "$1" = "--move-ssh-port" ]; then
    echo "==> Moving system SSH to port 2222..."
    CURRENT_PORT=$(grep -E "^Port [0-9]+" /etc/ssh/sshd_config 2>/dev/null | head -1 | awk '{print $2}')
    CURRENT_PORT=${CURRENT_PORT:-22}

    if [ "$CURRENT_PORT" = "2222" ]; then
        echo "    Already on port 2222. You're good."
        exit 0
    fi

    # SAFETY: Add port 2222 while KEEPING the old port
    # This way you can't get locked out
    echo "    Adding port 2222 alongside current port $CURRENT_PORT..."
    grep -q "^Port 2222" /etc/ssh/sshd_config || echo "Port 2222" >> /etc/ssh/sshd_config
    systemctl restart sshd

    PUBLIC_IP=$(curl -s ifconfig.me 2>/dev/null || echo "<your-ip>")
    echo ""
    echo "    Done! SSH now listens on BOTH port $CURRENT_PORT and 2222."
    echo ""
    echo "    Test from another terminal:"
    echo "      ssh -p 2222 -i <your-key> $(whoami)@$PUBLIC_IP"
    echo ""
    echo "    Once confirmed, run step 2:"
    echo "      sudo bash setup.sh /tmp/ssh-portfolio-linux-amd64"
    echo ""
    echo "    (The old port $CURRENT_PORT will be removed in step 2"
    echo "     when the portfolio takes over port 22)"
    exit 0
fi

BINARY_PATH=${1:-"/tmp/ssh-portfolio-linux-amd64"}

echo "==> SSH Portfolio Server Setup"
echo ""

if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: Binary not found at $BINARY_PATH"
    echo "Usage: sudo bash setup.sh /path/to/ssh-portfolio-linux-amd64"
    exit 1
fi

# 1. Verify port 2222 is available
if ! ss -tlnp | grep -q ':2222 '; then
    echo "Error: SSH is not listening on port 2222."
    echo "Run first:  sudo bash setup.sh --move-ssh-port"
    exit 1
fi
echo "==> System SSH confirmed on port 2222"

# 2. Now remove old port and keep only 2222
echo "==> Cleaning up SSH ports..."
# Remove all Port lines, then add only 2222
sed -i '/^Port /d' /etc/ssh/sshd_config
echo "Port 2222" >> /etc/ssh/sshd_config
systemctl restart sshd
echo "    System SSH now only on port 2222"

# 3. Stop any existing portfolio service
systemctl stop ssh-portfolio 2>/dev/null || true

# 4. Create service user and directories
echo "==> Creating service user and directories..."
id -u $SERVICE_USER &>/dev/null || useradd -r -s /bin/false -d $INSTALL_DIR $SERVICE_USER
mkdir -p $INSTALL_DIR/{.ssh,data}

# 5. Install binary
echo "==> Installing binary..."
cp "$BINARY_PATH" $INSTALL_DIR/ssh-portfolio
chmod +x $INSTALL_DIR/ssh-portfolio
chown -R $SERVICE_USER:$SERVICE_USER $INSTALL_DIR

# 6. Generate SSH host key
if [ ! -f "$INSTALL_DIR/.ssh/id_ed25519" ]; then
    echo "==> Generating SSH host key..."
    sudo -u $SERVICE_USER ssh-keygen -t ed25519 -f $INSTALL_DIR/.ssh/id_ed25519 -N "" -q
fi

# 7. Install systemd service
echo "==> Installing systemd service..."
cat > /etc/systemd/system/ssh-portfolio.service << 'EOF'
[Unit]
Description=SSH Portfolio TUI Server
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=portfolio
Group=portfolio
WorkingDirectory=/opt/ssh-portfolio

ExecStart=/opt/ssh-portfolio/ssh-portfolio

Environment=SSH_HOST=0.0.0.0
Environment=SSH_PORT=22
Environment=HOST_KEY_DIR=/opt/ssh-portfolio/.ssh
Environment=DB_PATH=/opt/ssh-portfolio/data/analytics.db

Restart=always
RestartSec=5

NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/ssh-portfolio/data /opt/ssh-portfolio/.ssh
PrivateTmp=true
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
EOF

# 8. Enable and start
echo "==> Starting service..."
systemctl daemon-reload
systemctl enable ssh-portfolio
systemctl start ssh-portfolio

# 9. Verify
sleep 2
if systemctl is-active --quiet ssh-portfolio; then
    PUBLIC_IP=$(curl -s ifconfig.me 2>/dev/null || echo "<your-ip>")
    echo ""
    echo "==> Setup complete!"
    echo ""
    echo "    Portfolio:  ssh $PUBLIC_IP"
    echo "    Admin SSH:  ssh -p 2222 $(whoami)@$PUBLIC_IP"
    echo "    Logs:       sudo journalctl -u ssh-portfolio -f"
    echo "    Status:     sudo systemctl status ssh-portfolio"
else
    echo ""
    echo "==> Service failed to start. Check logs:"
    echo "    sudo journalctl -u ssh-portfolio -e"
    exit 1
fi
