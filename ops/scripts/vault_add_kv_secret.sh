#!/usr/bin/env bash
# adds a secret to Vault kv engine
#
# Assumes that:
# Vault address is set in $VAULT_ADDR
# vault login is already performed (https://www.vaultproject.io/docs/commands/login.html)
#
# Arguments:
# 1st (required) - secret path
# 2nd (required) - secret key
# 3rd (required) - secret value

if [ "${#}" != "3" ]; then
    echo "usage ${0} <secret_path> <secret_key> <secret_value>"
    exit 1
fi

vault kv put "${1}" "${2}"="${3}" 1>/dev/null
