# Test

## Run
The program is a self-contained demo: it starts the WebSocket echo backend and
the proxy, then connects a client through the proxy to send a couple of
messages and read the echoes.

```sh
go run .
```

Expected output:
```
2026/08/10 15:43:35 websocket backend listening on [::]:11000
2026/08/10 15:43:35 proxy listening on :10999
2026/08/10 15:43:35 WebSocket handshake request dump:
GET / HTTP/1.1
Host: localhost:11000
Accept-Encoding: gzip
Connection: Upgrade
Sec-Websocket-Key: DwQt5H5/Jy0aPQwY9qHsRw==
Sec-Websocket-Version: 13
Upgrade: websocket
User-Agent: Go-http-client/1.1

2026/08/10 15:43:35 client send: Hello
2026/08/10 15:43:35 backend received: Hello
2026/08/10 15:43:35 client received: ECHO: Hello
2026/08/10 15:43:35 client send: World
2026/08/10 15:43:35 backend received: World
2026/08/10 15:43:35 client received: ECHO: World
```

## What happens
1. The client dials `wss://localhost:11000` using the proxy at `localhost:10999`.
2. The proxy MITMs the `CONNECT` and dumps the WebSocket handshake request
   (`GET /` with `Upgrade: websocket`), logging it.
3. goproxy forwards the handshake to the WebSocket server, then proxies the
   frames in both directions for the client.
