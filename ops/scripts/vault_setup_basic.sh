#!/usr/bin/env bash
# Disables Vault default kv engine at secret/ and start a new v1 one on kv/
# Sets up Vault VAULT_TOKEN and VAULT_ADDR for vagrant user
#
# environment variables are set in ~/.profile for the user executing the script
# if ~/.bash_profile is present ~/.profile will not be loaded so need to take appropriate measures
#
# vault root token is assumed in /etc/vault.d/.vault-token
# this is the default location used by vault_init.sh included in the slavrd/vault vagrant box

# set environment variables
ROOT_TOKEN_PATH="/etc/vault.d/.vault-token"
PROFILE="${HOME}/.profile"

grep 'VAULT_ADDR=' ${PROFILE} || {
    echo "export VAULT_ADDR=http://127.0.0.1:8200" | tee -a ${PROFILE}
} 

grep 'VAULT_TOKEN=' ${PROFILE} || {
    echo "export VAULT_TOKEN=$(cat ${ROOT_TOKEN_PATH})" | tee -a ${PROFILE}
} && {
    sed -i "s/VAULT_TOKEN=.*/VAULT_TOKEN=$(cat ${ROOT_TOKEN_PATH})/" ${PROFILE}
}

source ${PROFILE}

# disable default kv secrets engine (it's v2 for vault in dev)
vault secrets disable secret/ 1>/dev/null

# enable a v1 kv secrets engine
vault secrets list | grep kv/ || vault secrets enable -version=1 kv 1>/dev/null
