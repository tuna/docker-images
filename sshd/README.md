This is for running a container of `sshd` to allow uploading files. To start, use: 

```
/usr/bin/docker run --rm \
  -p ${port}:22 \
  -v /path/config:/config:ro \
  -v /path/srv:/srv \
  -e USER=${username} \
  tunathu/sshd
```

, where `/path/config` contains `ssh_host_{dsa,ecdsa,rsa}_key{,.pub}` for the host key of the container and `authorized_keys` for the public keys allowed to login. All the files in `/path/config` shall be owned by `root`. `/path/srv` is an example for the location where the login users can store files.
