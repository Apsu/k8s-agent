#!/usr/bin/env bash

set -euo pipefail

### Paths ###
AGENT_CONF=/etc/default/k8s-agent
AGENT_DIR=/opt/k8s-agent
RKE2_CONF_DIR=/etc/rancher/rke2
RKE2_DATA_DIR=/var/lib/rancher/rke2
RKE_MANIFEST_DIR=$RKE2_DATA_DIR/server/manifests
KUBECONFIG=$RKE2_CONF_DIR/rke2.yaml
KUBECTL=$RKE2_DATA_DIR/bin/kubectl
NETWORKD_DIR=/etc/systemd/network

### Logging ###
exec >> $AGENT_DIR/deploy.log 2>&1

### Load agent vars ###
set -a
. $AGENT_CONF
set +a

# Handle bootstrap case
if [[ $NODE_ROLE == "bootstrap" ]]; then
  RKE2_TYPE="server"
elif [[ $NODE_ROLE == "controller" ]]; then
  RKE2_TYPE="server"
else
  RKE2_TYPE="agent"
fi

### Summary ###
cat <<EOF
Deploying $NODE_NAME:
- Token: $TOKEN
- Public IP: $PUBLIC_IP
- API IP: $API_IP/$API_CIDR
- Interface: $INTERFACE
- Node role: $NODE_ROLE
- Node region: $NODE_REGION
- GPU count: $GPU_COUNT
- RKE2 type: $RKE2_TYPE
- RKE2 release: $RKE2_RELEASE
EOF

### Fixup networking ###

# Remove public IP from network config
sed -i "/^Address=$PUBLIC_IP.*/d" $NETWORKD_DIR/20-wired.network

# Insert private VIP route as on-link
if ! grep -q $API_IP $NETWORKD_DIR/20-wired.network; then
  envsubst < $AGENT_DIR/configs/vip-route.network >> $NETWORKD_DIR/20-wired.network
fi

# Apply the change
networkctl reload

### RKE2 ###
# Ensure directories
mkdir -m755 -p $RKE2_CONF_DIR $RKE_MANIFEST_DIR

# Install RKE2
curl -sfL https://get.rke2.io | INSTALL_RKE2_TYPE=$RKE2_TYPE INSTALL_RKE2_VERSION=$RKE2_RELEASE sh -

# Copy custom registry
# cp configs/registries.yaml $RKE2_CONF_DIR/registries.yaml

# Reload network config
networkctl reload

if [[ $RKE2_TYPE == "server" ]]; then
  if [[ $NODE_ROLE == "bootstrap" ]]; then
    # # Generate self-signed cert
    # openssl genrsa -out tls.key 2048
    # openssl req -new -key tls.key -out tls.csr -subj "/CN=$PUBLIC_IP"
    # openssl x509 -req -in tls.csr -sha256 -days 365 -signkey tls.key -out tls.crt

    # export TLS_KEY=$(base64 -w0 tls.key)
    # export TLS_CERT=$(base64 -w0 tls.crt)

    if ! ip -o address show dev $INTERFACE | grep -q $API_IP/$API_CIDR; then
      ip address add $API_IP/$API_CIDR dev $INTERFACE
    fi

    # Render bootstrap config
    envsubst < $AGENT_DIR/configs/bootstrap-config.yaml > $RKE2_CONF_DIR/config.yaml
  fi

  # Render manifests
  envsubst < $AGENT_DIR/manifests/kube-vip.yaml > $RKE_MANIFEST_DIR/kube-vip.yaml
  envsubst < $AGENT_DIR/manifests/kube-vip-rbac.yaml > $RKE_MANIFEST_DIR/kube-vip-rbac.yaml
  envsubst < $AGENT_DIR/manifests/rke2-cilium-values.yaml > $RKE_MANIFEST_DIR/rke2-cilium-values.yaml
  envsubst < $AGENT_DIR/manifests/rke2-coredns-values.yaml > $RKE_MANIFEST_DIR/rke2-coredns-values.yaml
  envsubst < $AGENT_DIR/manifests/nvidia-gpu-operator.yaml > $RKE_MANIFEST_DIR/nvidia-gpu-operator.yaml
  envsubst < $AGENT_DIR/manifests/kube-prometheus-stack.yaml > $RKE_MANIFEST_DIR/kube-prometheus-stack.yaml
  # envsubst < $AGENT_DIR/manifests/web-gateway.yaml > $RKE_MANIFEST_DIR/web-gateway.yaml
  # envsubst < $AGENT_DIR/manifests/grafana-routes.yaml > $RKE_MANIFEST_DIR/grafana-routes.yaml

  # Copy CRDs
  cp $AGENT_DIR/manifests/rke2-cilium-crds.yaml $RKE_MANIFEST_DIR/rke2-cilium-crds.yaml

  if [[ $NODE_ROLE == "bootstrap" ]]; then
    # Start and wait
    systemctl start rke2-server.service

    # Wait for node to become active
    while ! $KUBECTL --kubeconfig $KUBECONFIG get node $NODE_NAME | grep -q "Ready"; do
      sleep 5
    done
  fi

  # Render server config
  envsubst < $AGENT_DIR/configs/server-config.yaml > $RKE2_CONF_DIR/config.yaml

  # Run it
  systemctl enable rke2-server.service
  systemctl restart --no-block rke2-server.service
  systemctl mask rke2-agent.service
elif [[ $RKE2_TYPE == "agent" ]]; then
  # Render agent config
  envsubst < $AGENT_DIR/configs/agent-config.yaml > $RKE2_CONF_DIR/config.yaml

  # Run it
  systemctl enable --now --no-block rke2-agent.service
  systemctl mask rke2-server.service
fi

# Secure config
chown 0600 $RKE2_CONF_DIR/config.yaml

### Finish ###
touch $AGENT_DIR/.deployed
echo "Deployment time: $SECONDS seconds"
