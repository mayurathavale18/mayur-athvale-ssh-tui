#!/bin/bash
set -e

# SSH Portfolio — One-command deploy to EC2
#
# Usage:
#   ./deploy/deploy.sh
#
# Config (set as env vars or edit defaults below):
#   EC2_HOST    — EC2 public IP or hostname
#   EC2_USER    — SSH user (default: ubuntu)
#   EC2_PORT    — Admin SSH port (default: 2222)
#   EC2_KEY     — Path to SSH key (default: ~/.ssh/id_ed25519)

EC2_HOST="${EC2_HOST:-}"
EC2_USER="${EC2_USER:-ubuntu}"
EC2_PORT="${EC2_PORT:-2222}"
EC2_KEY="${EC2_KEY:-$HOME/.ssh/id_ed25519}"

BINARY_NAME="ssh-portfolio-linux-amd64"
BUILD_DIR="build"
REMOTE_TMP="/tmp/$BINARY_NAME"
INSTALL_DIR="/opt/ssh-portfolio"

# --- Validate ---
if [ -z "$EC2_HOST" ]; then
    echo "Error: EC2_HOST is not set."
    echo "Usage: EC2_HOST=<your-ec2-ip> ./deploy/deploy.sh"
    exit 1
fi

SSH_CMD="ssh -p $EC2_PORT -i $EC2_KEY -o StrictHostKeyChecking=no $EC2_USER@$EC2_HOST"
SCP_CMD="scp -P $EC2_PORT -i $EC2_KEY -O -o StrictHostKeyChecking=no"

echo "==> Building binary for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $BUILD_DIR/$BINARY_NAME ./cmd/server
echo "    Built: $(du -sh $BUILD_DIR/$BINARY_NAME | cut -f1)"

echo "==> Copying binary to $EC2_USER@$EC2_HOST..."
$SCP_CMD $BUILD_DIR/$BINARY_NAME $EC2_USER@$EC2_HOST:$REMOTE_TMP

echo "==> Installing and restarting service..."
$SSH_CMD << EOF
    sudo systemctl stop ssh-portfolio
    sudo cp $REMOTE_TMP $INSTALL_DIR/ssh-portfolio
    sudo chmod +x $INSTALL_DIR/ssh-portfolio
    sudo chown portfolio:portfolio $INSTALL_DIR/ssh-portfolio
    sudo systemctl start ssh-portfolio
    sleep 1
    sudo systemctl is-active --quiet ssh-portfolio && echo "    Service is running" || echo "    ERROR: service failed to start"
    rm -f $REMOTE_TMP
EOF

echo ""
echo "==> Deploy complete!"
echo "    Portfolio:  ssh $EC2_HOST"
echo "    Logs:       ssh -p $EC2_PORT -i $EC2_KEY $EC2_USER@$EC2_HOST 'sudo journalctl -u ssh-portfolio -f'"
