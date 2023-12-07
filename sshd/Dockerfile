FROM debian:bookworm

RUN apt-get update && apt-get install -y openssh-server less rsync
RUN mkdir /var/run/sshd
RUN echo 'PasswordAuthentication no' >> /etc/ssh/sshd_config
RUN rm -rf /etc/ssh/ssh_host*
RUN for i in "dsa" "ecdsa" "rsa"; do \
      for j in "" ".pub"; do \
        ln -sf "/config/ssh_host_${i}_key${j}" /etc/ssh/; \
      done; \
    done

COPY start_ssh.sh /

EXPOSE 22

CMD exec /start_ssh.sh
