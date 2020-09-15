FROM golang:1.15-alpine AS builder

RUN apk add curl jq git
RUN mkdir -p /tmp/build
WORKDIR /tmp/build
RUN \
  github="cyfdecyf/cow"; \
  release="$(curl -Ls "https://api.github.com/repos/$github/releases/latest" | jq -r .tarball_url)"; \
  curl -Ls "$release" | tar -xvz --strip-components 1 
RUN \
  go get github.com/shadowsocks/shadowsocks-go/shadowsocks && \
  go get github.com/cyfdecyf/bufio && \
  go get github.com/cyfdecyf/leakybuf && \
  go get github.com/cyfdecyf/color
RUN go build -o cow  

FROM alpine:3.12
COPY --chown=root:root --from=builder /tmp/build/cow /usr/bin/cow
