#!/usr/bin/env bash
# Sets up the environment variables for the current user
#
# usage: provision_client_env.sh <redis_addr> <redis_pass> <vault_addr>

if [ "$#" != "3" ]; then
    echo "usage: $0 <redis_addr> <redis_pass> <vault_addr>" >&2
    exit 1
fi

PROFILE="${HOME}/.profile"

grep 'REDIS_ADDR=' ${PROFILE} || {
    echo "export REDIS_ADDR='${1}'" | tee -a ${PROFILE}
} && {
    sed -i "s/REDIS_ADDR=.*/REDIS_ADDR=\'${1}\'/" ${PROFILE}
}

grep 'REDIS_PASS=' ${PROFILE} || {
    echo "export REDIS_PASS='${2}'" | tee -a ${PROFILE}
} && {
    sed -i "s/REDIS_PASS=.*/REDIS_PASS=\'${2}\'/" ${PROFILE}
}

grep 'VAULT_ADDR=' ${PROFILE} || {
    echo "export VAULT_ADDR='${3}'" | tee -a ${PROFILE}
} && {
    sed -i "s|VAULT_ADDR=.*|VAULT_ADDR=\'${3}\'|" ${PROFILE}
}
