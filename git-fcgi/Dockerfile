FROM debian:bookworm
LABEL maintainer="Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>"

RUN apt-get update && \
	apt-get install -y spawn-fcgi fcgiwrap multiwatch git
RUN git config --system uploadpack.allowReachableSHA1InWant true
RUN git config --system uploadpack.allowFilter true

COPY cmd.sh /cmd.sh
RUN mkdir -p /data 

EXPOSE 5000
CMD /cmd.sh

USER nobody
