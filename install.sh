#!/usr/bin/env bash

set -euo pipefail

### Init ###
# Check if required data is provided
if [[ -z "$TOKEN" || -z "$PUBLIC_IP" ]]; then
    echo <<- EOF
    Usage: TOKEN=... PUBLIC_IP=... $0

    Environment variables:
    TOKEN           Secure token for cluster joining (required)
    PUBLIC_IP       Public IP this node is reachable at (required)
    VIRTUAL_IP      VIP used for API endpoint (default: 192.168.1.1)
    VIRTUAL_CIDR    VIP CIDR mask (default: 32)
    INTERFACE       Interface to use for K8s networking (default: eno1)
    NODE_TYPE       RKE2 node type; [server | agent] (default: agent)
    FIRST_NODE      First controller in the cluster; [true | false] (default: false)
    RKE2_VERSION    RKE2 version to deploy (default: v1.30.5+rke2r1)
EOF
    exit 1
fi

### Overridable defaults ###
VIRTUAL_IP=${VIRTUAL_IP:-192.168.1.1}
VIRTUAL_CIDR=${VIRTUAL_CIDR:-32}
INTERFACE=${INTERFACE:-eno1}
NODE_TYPE=${NODE_TYPE:-agent}                  # Default to worker node
FIRST_NODE=${FIRST_NODE:-false}                # Default to secondary controller when type is 'server'
RKE2_VERSION=${RKE2_VERSION:-v1.30.5+rke2r1}

### Paths ###
TGZ_URL="https://github.com/Apsu/k8s-agent/releases/download/v0.0.3/package.tgz"
INSTALL_PATH="/opt/k8s-agent"
CONFIG_PATH="/etc/default/k8s-agent"
SERVICE_PATH="/etc/systemd/system/k8s-agent.service"

# Store config
cat <<EOF > $CONFIG_PATH
TOKEN=$TOKEN
PUBLIC_IP=$PUBLIC_IP
VIRTUAL_IP=$VIRTUAL_IP
VIRTUAL_CIDR=$VIRTUAL_CIDR
INTERFACE=$INTERFACE
NODE_TYPE=$NODE_TYPE
FIRST_NODE=$FIRST_NODE
RKE2_VERSION=$RKE2_VERSION
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
