#!/bin/bash

THREADS=${THREADS:-"64"}
QUEUE_SIZE=${QUEUE_SIZE:-"48"}
MAX_BUFFER_SIZE=${MAX_BUFFER_SIZE:-"1G"}
BUFFER_PATH=${BUFFER_PATH:-"/tmp"}

export HOME=/var/www
(
    while true; do
        /go-queue --queue-size "${QUEUE_SIZE}" --port-number 8888
    done
)&
(
    while true; do
        /buffering-proxy --listen ":5000" --max-buffer-size "${MAX_BUFFER_SIZE}" --on-disk-buffer-path "${BUFFER_PATH}" --upstream "unix:/tmp/fcgi.sock"
    done
)&
exec /usr/bin/spawn-fcgi -u www-data -g www-data -s /tmp/fcgi.sock -n -- /usr/bin/multiwatch -f "${THREADS}" -- /usr/sbin/fcgiwrap
