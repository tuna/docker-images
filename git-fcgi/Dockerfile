FROM debian:buster
MAINTAINER Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>

RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free" > /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free" >> /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free" >> /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian-security/ buster/updates main contrib non-free" >> /etc/apt/sources.list

RUN apt-get update && \
	apt-get install -y spawn-fcgi fcgiwrap multiwatch git
RUN git config --system uploadpack.allowReachableSHA1InWant true

COPY cmd.sh /cmd.sh
RUN mkdir -p /data 

EXPOSE 5000
CMD /cmd.sh

USER nobody
