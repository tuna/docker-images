THREADS=${THREADS:-"64"}
QUEUE_SIZE=${QUEUE_SIZE:-"48"}
export HOME=/var/www
(
    while true; do
        /go-queue --queue-size "${QUEUE_SIZE}" --port-number 8888
    done
)&
exec /usr/bin/spawn-fcgi -u www-data -g www-data -p 5000 -n -- /usr/bin/multiwatch -f ${THREADS} -- /usr/sbin/fcgiwrap
