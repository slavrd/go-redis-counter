#!/usr/bin/env bash
# creates a vault token accosiated with a vault policy 
# and sets it up in VAULT_TOKEN env variable

# usage vault_create_token.sh <vault_addr> <vault_token> <policy_name>
# vault address must be in the form http[s]://<HOST>[:PORT]

if [ "$#" != 3 ]; then
    echo "usage: ${0} <vault_addr> <vault_token> [policy_name]" >&2
    exit 1
fi

VADDR=${1}
VTOKEN=${2}
VPNAME=${3}

# check/install prerequisite packages
PKGS="curl jq" 
which ${PKGS} >>/dev/null || {
    sudo apt-get update
    sudo apt-get install -y ${PKGS}
}

# create the token via Vault's API
RESP=$(curl -sSf \
            --header "X-Vault-Token: $VTOKEN" \
            --request "POST" \
            --data "{\"policies\":[\"$VPNAME\"]}" \
            "$VADDR/v1/auth/token/create")

if [ "$?" != 0 ]; then
    echo "failed craeting vault token" >&2
    exit 1
fi

VT=$(echo $RESP | jq -r '.auth.client_token')

# set token to VAULT_TOKEN env variable for current user
PROFILE="$HOME/.profile"

grep 'VAULT_TOKEN=' ${PROFILE} || {
    echo "export VAULT_TOKEN=${VT}" | tee -a ${PROFILE}
} && {
    sed -i "s/VAULT_TOKEN=.*/VAULT_TOKEN=${VT}/" ${PROFILE}
}
