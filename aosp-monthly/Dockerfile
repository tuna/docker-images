FROM debian:bookworm-slim
LABEL maintainer="Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>"

RUN apt-get update && \
	apt-get install -y --no-install-recommends ca-certificates curl git python3 && \
	rm -rf /var/lib/apt/lists/*

RUN curl --fail --location --show-error \
		https://mirrors.tuna.tsinghua.edu.cn/git/git-repo \
		--output /usr/local/bin/repo && \
	chmod a+x /usr/local/bin/repo

RUN ln -sr /usr/bin/python3 /usr/bin/python

COPY aosp-monthly.sh /usr/local/bin/

RUN useradd -m aosp && install -d -o aosp -g aosp /data

USER aosp
CMD ["/usr/local/bin/aosp-monthly.sh"]
