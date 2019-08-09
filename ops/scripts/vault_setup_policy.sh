#!/usr/bin/env bash
# sets up a policy in Vault using CLI.
# 
# usage: vault_setup_policy.sh <policy_name> <policy_file_path>

if [ "$#" != "2" ]; then
    echo "usage: vault_setup_policy.sh <policy_name> <policy_file_path>" >&2
fi

vault policy write "$1" "$2" 1>/dev/null
