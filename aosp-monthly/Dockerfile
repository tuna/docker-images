FROM debian:bullseye
MAINTAINER Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>

RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free" > /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free" >> /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free" >> /etc/apt/sources.list && \
	echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian-security/ bullseye-security main contrib non-free" >> /etc/apt/sources.list

RUN apt-get update && \
	apt-get install -y wget curl git python3-dev python3-pip

RUN pip3 install --upgrade pip setuptools

RUN mkdir -p /data

RUN curl https://mirrors.tuna.tsinghua.edu.cn/git/git-repo > /usr/local/bin/repo && \
	chmod a+x /usr/local/bin/repo

RUN ln -sr /usr/bin/python3 /usr/bin/python

COPY aosp-monthly.sh /usr/local/bin/
CMD /usr/local/bin/aosp-monthly.sh
RUN useradd -m aosp
USER aosp
