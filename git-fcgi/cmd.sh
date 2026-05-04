#!/bin/bash

THREADS=${THREADS:-"64"}
QUEUE_SIZE=${QUEUE_SIZE:-"48"}
export HOME=/var/www
(
    while true; do
        /go-queue --queue-size "${QUEUE_SIZE}" --port-number 8888
    done
)&

spawn_fcgi_args=(-u www-data -g www-data -p 5000 -n)
if [[ -n "${FCGI_BIND_ADDRESS:-}" ]]; then
    spawn_fcgi_args+=(-a "${FCGI_BIND_ADDRESS}")
fi

exec /usr/bin/spawn-fcgi "${spawn_fcgi_args[@]}" -- /usr/bin/multiwatch -f "${THREADS}" -- /usr/sbin/fcgiwrap
