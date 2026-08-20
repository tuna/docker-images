FROM golang:1.26-alpine AS builder
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY go-queue.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o go-queue go-queue.go

FROM debian:bookworm-slim
LABEL maintainer="Justin Wong <yuzhi.wang@tuna.tsinghua.edu.cn>"

RUN apt-get update && \
	apt-get install -y spawn-fcgi fcgiwrap multiwatch git --no-install-recommends && \
	rm -rf /var/lib/apt/lists/*

COPY cmd.sh /cmd.sh
COPY queue-wrapper.sh /queue-wrapper.sh

RUN git config --system uploadpack.allowReachableSHA1InWant true && \
	git config --system uploadpack.allowFilter true && \
	git config --system uploadpack.packObjectsHook /queue-wrapper.sh && \
	git config --system --add safe.directory '*' && \
	mkdir -p /data

COPY --from=builder --chown=0:0 /app/go-queue /go-queue

EXPOSE 5000
CMD ["/cmd.sh"]

USER nobody
