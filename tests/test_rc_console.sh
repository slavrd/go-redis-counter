#!/usr/bin/env bash
# runs the console app and verifies the output

# usage test_rc_console.sh <console_app_path> <reids_addr> <redis_pass>

# Vault address and toke must be set in $VAULT_ADDR and $VAULT_TOKEN

if [ "$#" != 3 ]; then
    echo "usage test_rc_console.sh <console_app_path> <reids_addr> <redis_pass>"
    exit 1
fi

APP="$1"
IFS=':' RADDR=($2)
RPASS="$3"

# run app once and record current counter value
VAL1="$($APP -a ${RADDR[0]} -pass $RPASS)" || {
    echo "failed running $APP"
    exit 1
}

# check app again and confirm value increases by 1 
VAL2="$($APP -a ${RADDR[0]} -pass $RPASS)" || {
    echo "failed running $APP"
    exit 1
}

WANT=$(expr $VAL1 + 1) || {
    echo "failed calculating expected result"
    exit 1
}

if [ "$VAL2" != "$WANT" ]; then
    echo "received wrong value got:$VAL2, want:$WANT, initial value:$VAL1"
    exit 1
fi

# run app again, taking redis password from vault
VAL3="$($APP -a ${RADDR[0]} -v-use)" || {
    echo "failed running $APP"
    exit 1
}

WANT2=$(expr $VAL2 + 1) || {
    echo "failed calculating expected result"
    exit 1
}

if [ "$VAL3" != "$WANT2" ]; then
    echo "received wrong value got:$VAL3, want:$WANT2, initial value:$VAL2"
    exit 1
fi
