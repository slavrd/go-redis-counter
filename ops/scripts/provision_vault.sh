#!/usr/bin/env bash
# basic tools and vault installation

# list of packages to install, separated by space
PKGS="curl"

# version of Vault to install
VAULT_VERSION="1.2.0"

which $PKGS || {
    sudo apt-get update
    sudo apt-get install -y $PKGS
}

# install Vault
curl -sSf -o /tmp/hc_install.sh https://raw.githubusercontent.com/slavrd/bash-various-scripts/master/install_hc_product.sh \
    && bash /tmp/hc_install.sh vault $VAULT_VERSION linux amd64
