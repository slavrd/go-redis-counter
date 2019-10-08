#!/usr/bin/env bash
# creates a vault token accosiated with a vault policy 
# and writes it out to a file.

# usage vault_create_token.sh <policy_name>
# assumes the VAULT_TOKEN and VAULT_ADDR variables are set.

if [ "$#" != 1 ]; then
    echo "usage: ${0} <policy_name>" >&2
    exit 1
fi

V_POICY_NAME=${1}

# check/install prerequisite packages
PKGS="curl jq" 
which ${PKGS} >>/dev/null || {
    sudo apt-get update
    sudo apt-get install -y ${PKGS}
}

# create the token via Vault's API
RESP=$(curl -sSf \
            --header "X-Vault-Token: $VAULT_TOKEN" \
            --request "POST" \
            --data "{\"policies\":[\"$V_POICY_NAME\"]}" \
            "$VAULT_ADDR/v1/auth/token/create")

if [ "$?" != 0 ]; then
    echo "failed craeting vault token" >&2
    exit 1
fi

VT=$(echo $RESP | jq -r '.auth.client_token')

# wirte out vault token to a variable.
[ -d "/vagrant/.vagrant_cache" ] || mkdir -p "/vagrant/.vagrant_cache"
echo -n ${VT} > "/vagrant/.vagrant_cache/.vault-policybased-token"
