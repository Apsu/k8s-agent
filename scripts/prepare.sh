#!/usr/bin/env bash

set -euo pipefail

env | sort

### Paths ###
AGENT_CONF=/etc/default/k8s-agent
AGENT_DIR=/opt/k8s-agent

### Logging ###
exec >> $AGENT_DIR/prepare.log 2>&1

### Load agent vars ###
set -a
. $AGENT_CONF
set +a

### Env flags ###
export DEBIAN_FRONTEND=noninteractive

### Packages, kernel, systemd ###
echo "Preparing node..."

# Set hostname
hostnamectl set-hostname $NODE_NAME

# Stop lambda units
systemctl disable --now cloudflared.service cloudflared-update.{service,timer} lambda-jupyter.service lambda-first-boot.service || true

# Stop some system units
systemctl disable --now irqbalance multipathd unattended-upgrades || true

# Purge all 3rd-party packages
apt autoremove --purge -y '~i !~OUbuntu'

# Purge unneeded packages and deps
apt autoremove --purge -y apparmor snapd nvidia* mesa* ruby* golang* javascript* docker* containerd* ffmpeg* emacs* font* gfortran* packagekit* polkit* thermald* power* multipath* irqbalance*

# Purge old kernels if present
apt autoremove --purge -y linux-{image,headers,modules,virtual,generic}-*

# Remove 3rd party repos
rm -f /etc/apt/sources.list.d/* /etc/apt/cloud-init.gpg.d/*

# Remove custom pieces
rm -rf /etc/systemd/system/cloudflared*
rm -rf /etc/systemd/system/lambda-jupyter*
rm -rf /lib/systemd/system/lambda-first-boot*
rm -f /usr/bin/disable-mig.sh

# Update repos
apt update

# Upgrade rest of system
apt upgrade -y

# Disable apparmor in kernel
sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="apparmor=0"/' /etc/default/grub
update-grub

# Refresh units
systemctl daemon-reload

### Finish ###

# Cleanup ubuntu user
rm -rf /home/ubuntu/.{cache,config,ipython,jupyter,lambda,local}

# Mark completion
touch $AGENT_DIR/.prepared
echo "Preparation time: $SECONDS seconds"

# Apparmor/driver unload/etc
reboot
