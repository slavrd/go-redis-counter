#! /usr/bin/env bash
# provisions the VM with the needed packages

# list of packages to install, separated by space
PKGS="git vim curl jq"

# version of Golang to install
GOVER="1.12"

which $PKGS || {
    sudo apt-get update
    sudo apt-get install -y $PKGS
}

# download and run golang install script
wget -q -P /tmp https://raw.githubusercontent.com/slavrd/bash-various-scripts/master/install-golang.sh && \
    bash /tmp/install-golang.sh 1.12 linux amd64
