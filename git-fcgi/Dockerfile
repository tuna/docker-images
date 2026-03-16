FROM debian:bookworm as builder
RUN apt-get update && \
	apt-get install -y golang --no-install-recommends

RUN mkdir -p /app
WORKDIR /app
COPY go-queue.go /app/
COPY buffering-proxy.go /app/
RUN CGO_ENABLED=0 go build -o go-queue go-queue.go
RUN CGO_ENABLED=0 go build -o buffering-proxy buffering-proxy.go

FROM debian:bookworm
LABEL maintainer="Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>"

RUN apt-get update && \
	apt-get install -y spawn-fcgi fcgiwrap multiwatch git --no-install-recommends && \
	apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*
RUN git config --system uploadpack.allowReachableSHA1InWant true
RUN git config --system uploadpack.allowFilter true
RUN git config --system uploadpack.packObjectsHook /queue-wrapper.sh
RUN git config --system --add safe.directory '*'

COPY cmd.sh /cmd.sh
COPY queue-wrapper.sh /queue-wrapper.sh
RUN mkdir -p /data

COPY --from=builder --chown=0:0 /app/go-queue /go-queue
COPY --from=builder --chown=0:0 /app/buffering-proxy /buffering-proxy

EXPOSE 5000
CMD [ "/cmd.sh" ]

USER nobody
