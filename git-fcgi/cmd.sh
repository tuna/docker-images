THREADS=${THREADS:-"64"}
export HOME=/var/www
/usr/bin/spawn-fcgi -u www-data -g www-data -p 5000 -n -- /usr/bin/multiwatch -f ${THREADS} -- /usr/sbin/fcgiwrap
