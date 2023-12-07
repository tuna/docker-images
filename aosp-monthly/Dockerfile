FROM debian:bookworm
LABEL maintainer="Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>"

RUN apt-get update && \
	apt-get install -y wget curl git python3-pip python3-setuptools

RUN mkdir -p /data

RUN curl https://mirrors.tuna.tsinghua.edu.cn/git/git-repo > /usr/local/bin/repo && \
	chmod a+x /usr/local/bin/repo

RUN ln -sr /usr/bin/python3 /usr/bin/python

COPY aosp-monthly.sh /usr/local/bin/
CMD /usr/local/bin/aosp-monthly.sh
RUN useradd -m aosp
USER aosp
