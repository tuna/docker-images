#!/bin/bash
USER=${USER:-"user"}

mkdir -p "/home/${USER}/.ssh"
chown -R 1000:1000 "/home/${USER}/.ssh"
chmod 755 "/home/${USER}/.ssh" "/home/${USER}"
ln -sf /config/authorized_keys "/home/${USER}/.ssh"

echo "${USER}:!:1000:1000::/home/${USER}:/bin/bash" >> /etc/passwd
echo "${USER}:x:1000:" >> /etc/group

exec /usr/sbin/sshd -D -e
