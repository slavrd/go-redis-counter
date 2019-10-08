#!/usr/bin/env bash
# sets up a vault token for the current user

VT=$(cat /home/vagrant/go/src/github.com/slavrd/go-redis-counter/.vagrant_cache/.vault-policybased-token)

# set token to VAULT_TOKEN env variable for current user
PROFILE="$HOME/.profile"

grep 'VAULT_TOKEN=' ${PROFILE} >>/dev/null || {
    echo "export VAULT_TOKEN=${VT}" | tee -a ${PROFILE}
} && {
    sed -i "s/VAULT_TOKEN=.*/VAULT_TOKEN=${VT}/" ${PROFILE}
}
