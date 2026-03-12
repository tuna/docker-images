#!/bin/bash

set -euo pipefail

QUEUE_SERVER_PORT=8888
exec 3<"/dev/tcp/127.0.0.1/${QUEUE_SERVER_PORT}"
exec 4>&1
exec >&2

all_input=$(cat)
waited=0
while read -u 3 -r line; do
    if [ "$line" = "0" ]; then
        break
    fi
    printf "Waiting in queue... (Position: %s)\r" "$line"
    waited=1
done
if [ "$waited" -eq 1 ]; then
    printf "\n"
fi

"$@" <<< "$all_input" 3<&- >&4 4>&-
