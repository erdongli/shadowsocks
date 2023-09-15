# Shadowsocks

A Shadowsocks implementation in Go.

> [!NOTE]
> This only supports TCP CONNECT and AEAD_CHACHA20_POLY1305.

## Usage

### Docker

#### Building ss-local and ss-remote

```
make docker app=local
make docker app=remote
```

#### Running ss-remote

```
docker run shadowsocks/remote:latest -p <port> -k <access key> -l <log level>
```

#### Running ss-local
```
docker run shadowsocks/local:latest -p <port> -r <remote host>:<remote port> -k <access key> -l <log level>
```

### Firefox
1. Go to `about:preferences`
2. Scroll down to `Network Settings`
3. Select `Settings...`
4. Select `Manual proxy configuration`
5. Type in ss-local's hostname in `SOCKS Host`
6. Type in ss-local's port number in `Port`
7. Enable `Proxy DNS when using SOCKS v5`
