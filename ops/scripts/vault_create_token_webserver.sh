#!/usr/bin/env bash
# creates a vault token accosiated with a vault policy and places it in enviroment.conf

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

# setup token
CF='/tmp/environment.conf'

[ -f ${CF} ] || {
    echo "File ${CF} does not exist"
    exit 1
}

sed -i "s/VAULT_TOKEN_VALUE/${VT}/" ${CF}