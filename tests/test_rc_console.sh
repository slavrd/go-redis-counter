#!/usr/bin/env bash
# runs the console app and verifies the output

# usage test_rc_console.sh <console_app_path> <reids_addr> <redis_pass>

if [ "$#" != 3 ]; then
    echo "usage test_rc_console.sh <console_app_path> <reids_addr> <redis_pass>"
fi

APP="$1"
IFS=':' RADDR=($2)
RPASS="$3"

VAL1="$($APP -a ${RADDR[0]} -pass $RPASS)" || {
    echo "failed running $APP"
    exit 1
}

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
