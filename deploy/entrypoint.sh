#!/bin/sh
set -e

KEY_PATH="${HOST_KEY_DIR:-/app/.ssh}/id_ed25519"

# Generate host key on first run
if [ ! -f "$KEY_PATH" ]; then
    echo "Generating SSH host key..."
    ssh-keygen -t ed25519 -f "$KEY_PATH" -N "" -q
    echo "Host key generated at $KEY_PATH"
fi

exec "$@"
