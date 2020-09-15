To start a cow-proxy, run:

```
docker run -i --rm --name=cow-proxy \
  -p 127.0.0.1:8123:8123 \
  -v /path/to/cow.conf:/cow.conf:ro \
  tunathu/shadowsocks /usr/bin/cow -request -rc /cow.conf
```

This Dockerfile automatically finds the latest version of cow-proxy and builds it.
