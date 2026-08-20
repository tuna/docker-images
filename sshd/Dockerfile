FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends openssh-server less rsync && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /var/run/sshd && \
    echo 'PasswordAuthentication no' >> /etc/ssh/sshd_config && \
    rm -f /etc/ssh/ssh_host* && \
    for i in "dsa" "ecdsa" "rsa"; do \
      for j in "" ".pub"; do \
        ln -sf "/config/ssh_host_${i}_key${j}" /etc/ssh/; \
      done; \
    done

COPY start_ssh.sh /

EXPOSE 22

CMD ["/start_ssh.sh"]
