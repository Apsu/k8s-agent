#!/usr/bin/env bash

set -euo pipefail

### Overridable defaults ###
API_IP=${API_IP:-"192.168.1.1"}
API_CIDR=${API_CIDR:-"32"}
INTERFACE=${INTERFACE:-"eno1"}
RKE2_RELEASE=${RKE2_RELEASE:-"v1.30.5+rke2r1"}
NODE_ROLE=${NODE_ROLE:-"worker"}
NODE_REGION=${NODE_REGION:-"us-south-2"}
GPU_COUNT=${GPU_COUNT:-"0"}
AGENT_VERSION=${AGENT_VERSION:-"v0.0.6"}

### Init ###
# Check if required data is provided
if [[ -z "${TOKEN:-}" || -z "${PUBLIC_IP:-}" || -z "${NODE_NAME:-}" ]]; then
    echo <<- EOF
    Usage: TOKEN=... PUBLIC_IP=... NODE_NAME=... $0

    Environment variables:
    TOKEN           Secure token for cluster joining (required)
    PUBLIC_IP       Public IP this node is reachable at (required)
    NODE_NAME       Hostname to use for this node (required)
    API_IP          IP used for API endpoint (default: $API_IP)
    API_CIDR        IP CIDR mask (default: $API_CIDR)
    INTERFACE       Interface to use for K8s networking (default: $INTERFACE)
    RKE2_RELEASE    RKE2 release to deploy (default: $RKE2_RELEASE)
    NODE_ROLE       Node role; [bootstrap | controller | worker] (default: $NODE_ROLE)
    NODE_REGION     Node region (default: $NODE_REGION)
    GPU_COUNT       Node GPU count (default: $GPU_COUNT)
    AGENT_VERSION   Version of this agent to deploy (default: $AGENT_VERSION)
EOF
    exit 1
fi

### Paths ###
TGZ_URL="https://github.com/Apsu/k8s-agent/releases/download/$AGENT_VERSION/package.tgz"
INSTALL_PATH="/opt/k8s-agent"
CONFIG_PATH="/etc/default/k8s-agent"
SERVICE_PATH="/etc/systemd/system/k8s-agent.service"

# Store config
cat <<EOF > $CONFIG_PATH
STATE=Initialized
TOKEN=$TOKEN
PUBLIC_IP=$PUBLIC_IP
API_IP=$API_IP
API_CIDR=$API_CIDR
INTERFACE=$INTERFACE
RKE2_RELEASE=$RKE2_RELEASE
NODE_NAME=$NODE_NAME
NODE_ROLE=$NODE_ROLE
NODE_REGION=$NODE_REGION
GPU_COUNT=$GPU_COUNT
AGENT_VERSION=$AGENT_VERSION
EOF

# Download and extract the release
mkdir -p $INSTALL_PATH
curl -sfL $TGZ_URL -o /tmp/package.tgz
tar -xzf /tmp/package.tgz -C $INSTALL_PATH

# Create systemd service
cp $INSTALL_PATH/configs/k8s-agent.service $SERVICE_PATH

# Start the agent service
systemctl daemon-reload
systemctl enable --now k8s-agent

echo "Agent installed and started."
