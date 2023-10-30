FROM debian:bullseye
MAINTAINER Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>

RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free" > /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free" >> /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free" >> /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian-security/ bullseye-security main contrib non-free" >> /etc/apt/sources.list

RUN apt-get update && \
	apt-get install -y spawn-fcgi fcgiwrap multiwatch git
RUN git config --system uploadpack.allowReachableSHA1InWant true
RUN git config --system uploadpack.allowFilter true

COPY cmd.sh /cmd.sh
RUN mkdir -p /data 

EXPOSE 5000
CMD /cmd.sh

USER nobody
