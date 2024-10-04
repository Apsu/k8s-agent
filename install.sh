#!/usr/bin/env bash

set -euo pipefail

TOKEN=${TOKEN:-$1}
TGZ_URL="https://your.short.url/package.tgz"
INSTALL_PATH="/opt/k8s-agent"
CONFIG_PATH="/etc/default/k8s-agent"
SERVICE_PATH="/etc/systemd/system/k8s-agent.service"

# Check if token is provided
if [[ -z "$TOKEN" ]]; then
    echo "TOKEN is required."
    exit 1
fi

# Download and extract the release
mkdir -p $INSTALL_PATH
curl -sfL $TGZ_URL -o /tmp/package.tgz
tar -xzf /tmp/package.tgz -C $INSTALL_PATH

# Store token in config
echo "TOKEN=$TOKEN" > $CONFIG_PATH

# Create systemd service
cp $INSTALL_PATH/k8s-agent.service /etc/systemd/system/k8s-agent.service

# Start the agent service
systemctl daemon-reload
systemctl enable --now k8s-agent

echo "Agent installed and started."
