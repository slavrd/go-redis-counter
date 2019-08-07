#!/usr/bin/env bash
# starts Vault in dev mode bound to all interfaces and sets up needed variables
#
# environment variables are set in ~/.profile for the user executing the script
# if ~/.bash_profile is present ~/.profile will not be loaded so need to take appropriate measures
#
# Arguments:
# 1st (optional) - Vault Root token to use

# Check provided argument
if [ "${1}" != "" ]; then
    VT=$1
else
    VT='defDevV@ultRootT0ken'
fi

# set environment variables
PROFILE="${HOME}/.profile"

grep 'VAULT_ADDR=' ${PROFILE} || {
    echo "export VAULT_ADDR=http://127.0.0.1:8200" | tee -a ${PROFILE}
} 

grep 'VAULT_DEV_ROOT_TOKEN_ID=' ${PROFILE} || {
    echo "export VAULT_DEV_ROOT_TOKEN_ID=${VT}" | tee -a ${PROFILE}
} && {
    sed -i "s/VAULT_DEV_ROOT_TOKEN_ID=.*/VAULT_DEV_ROOT_TOKEN_ID=${VT}/" ${PROFILE}
}

source ${PROFILE}

# start Vault in dev mode
sudo killall vault
vault server -dev -dev-root-token-id="${VT}" --dev-listen-address='0.0.0.0:8200' &>vault.log &
sleep 5

# disable default kv secrets engine (it's v2 for vault in dev)
vault secrets disable secret/ 1>/dev/null

# enable a v1 kv secrets engine
vault secrets enable -version=1 kv 1>/dev/null