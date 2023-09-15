# Shadowsocks

A Shadowsocks implementation in Go.

> [!NOTE]
> Only TCP CONNECT and AEAD_CHACHA20_POLY1305 are supported.

## Usage

### Docker

#### Build ss-remote

```
make docker app=remote
```

#### Build ss-local

```
make docker app=local
```

#### Run ss-remote

```
docker run -p <port>:<port> -d shadowsocks/remote:latest \
  -p <port> \
  -k <access key> \
  -l <log level>
```

#### Run ss-local
```
docker run -p <port>:<port> -d shadowsocks/local:latest \
  -p <port> \
  -r <remote host>:<remote port> \
  -k <access key> \
  -l <log level>
```

### Firefox Configuration
1. Go to `about:preferences`
2. Scroll down to `Network Settings`
3. Select `Settings...`
4. Select `Manual proxy configuration`
5. Type in ss-local's hostname in `SOCKS Host`
6. Type in ss-local's port number in `Port`
7. Enable `Proxy DNS when using SOCKS v5`
