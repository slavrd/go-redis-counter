#!/usr/bin/env bash
# runs the console app and verifies the output

# usage test_rc_console.sh <console_app_path> <reids_addr> <redis_pass>

if [ "$#" != 3 ]; then
    echo "usage test_rc_console.sh <console_app_path> <reids_addr> <redis_pass>"
fi

VAL1="$($1 -a $2 -pass $3)" || {
    echo "failed running $1"
    exit 1
}

VAL2="$($1 -a $2 -pass $3)" || {
    echo "failed running $1"
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
