#!/usr/bin/env bash
# gets the application vault token from the Vault box
# and writes it out to the application servcie configuration file.
# usage: webserver_get_vault_token.sh <vault_ip_addr>

[ -z "${VAULT_IP_ADDR}" ] && {
    echo 'VAULT_IP_ADDR is not set' >&2
    exit 1
} 

CF='/tmp/environment.conf'
[ -f ${CF} ] || {
    echo "File ${CF} does not exist"
    exit 1
}

grep 'VAULT_TOKEN_VALUE' >>/dev/null ${CF} && {

    [ -f "/vagrant/.vagrant_cache/.vault-policybased-token" ] || {
        echo "/vagrant/.vagrant_cache/.vault-policybased-token does not exist" >&2
        exit 1
    }

    VT=$(cat /vagrant/.vagrant_cache/.vault-policybased-token)

    sed -i "s/VAULT_TOKEN_VALUE/${VT}/" ${CF}
}