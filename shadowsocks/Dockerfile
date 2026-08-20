FROM golang:1.21-alpine AS builder

RUN apk add --no-cache ca-certificates curl git jq

WORKDIR /tmp/build
ENV GO111MODULE=off

ADD https://api.github.com/repos/cyfdecyf/cow/releases/latest /tmp/cow-release.json

RUN \
  release="$(jq -r .tarball_url /tmp/cow-release.json)"; \
  curl --fail --location --show-error "$release" | tar -xz --strip-components 1
RUN \
  go get github.com/shadowsocks/shadowsocks-go/shadowsocks && \
  go get github.com/cyfdecyf/bufio && \
  go get github.com/cyfdecyf/leakybuf && \
  go get github.com/cyfdecyf/color
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o cow

FROM alpine:3.23
COPY --chown=root:root --from=builder /tmp/build/cow /usr/bin/cow
