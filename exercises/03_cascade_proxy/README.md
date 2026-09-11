# Exercise 03 - Cascade Proxy

> Estimated time: ~30 minutes.

## Goal
Build a **cascade proxy** (proxy aggregator): a chain of two proxies where
the first (middle, `:10998`) forwards every request to the second
(edge, `:10999`), which authenticates and forwards to the final server.

```
client → middle proxy (:10998) → edge proxy (:10999) → backend
```

## Learning outcomes
- Chain two goproxy instances.
- Understand single-hop headers (`Proxy-Authorization`, `Proxy-Connection`).
- Route plain HTTP and CONNECT requests through an upstream proxy.
- Enforce authentication on the edge proxy and inject credentials on the
  middle proxy.

## Background (read before coding, ~5 min)
- **Proxy chains**: a proxy can use another proxy as its "transport": for
  plain HTTP it sets `Transport.Proxy`; for CONNECT it sets `ConnectDial`.
- **Single-hop headers**: `Proxy-Authorization` is meant for one proxy hop
  only. By default goproxy strips it before forwarding, so the middle proxy
  must re-add it for the edge proxy (or keep it with `KeepHeader`).
- **`auth.ProxyBasic`** registers both an `OnRequest` handler (for plain
  requests) and a `HandleConnect` handler (for CONNECT), returning `407` on
  missing/invalid credentials.
  See: https://pkg.go.dev/github.com/elazarl/goproxy/ext/auth#ProxyBasic

## Prerequisites
- Go 1.24+ (https://go.dev/doc/install)

### Certificates
The HTTPS backend (`:11001`) serves `certs_localhost/localhost.crt`.

You have two options:
- **Generate your own** with `openssl` (you'll likely already have it):
  follow the instructions in `certs_localhost/generation.md`.
- **Copy the provided one** from the main project folder:
  ```sh
  cp -r ../../certs_localhost .
  ```

## Project layout
```
exercises/03_cascade_proxy/
├── main.go            <- to complete
├── backend.go         <- ready: HTTP backend (:11000) and HTTPS backend (:11001)
├── TEST.md            <- expected output of the curl tests
└── certs_localhost/   <- generate or copy from the main folder
```

> Note: `main.go` compiles only once all `TODO`s are completed; the Go
> compiler tells you exactly what is missing.

## Steps

### Step 1 - Edge proxy
The edge proxy must authenticate every request:
```go
auth.ProxyBasic(endProxy, "my_realm", func(user, passwd string) bool {
	return subtle.ConstantTimeCompare([]byte(user), []byte(_proxyUser)) == 1 &&
		subtle.ConstantTimeCompare([]byte(passwd), []byte(_proxyPassword)) == 1
})
```
- `auth.ProxyBasic`: https://pkg.go.dev/github.com/elazarl/goproxy/ext/auth#ProxyBasic
- `subtle.ConstantTimeCompare`: https://pkg.go.dev/crypto/subtle#ConstantTimeCompare
- Start the listener on `localhost:10999` in a goroutine (the main goroutine
  will block on the middle proxy at the end of `main`).

### Step 2 - Middle proxy
The middle proxy must forward everything to the edge proxy, adding the
credentials. Complete sub-steps 2a–2d in `main.go`:
1. **KeepHeader**: `middleProxy.KeepHeader = true` so `Proxy-Authorization`
   is not stripped.
2. **`Tr.Proxy`**: route plain HTTP requests to `http://localhost:10999`.
   See: https://pkg.go.dev/net/http#Transport.Proxy
3. **`ConnectDial`**: tunnel CONNECT requests through the edge proxy and add
   the credentials on the CONNECT request via
   `NewConnectDialToProxyWithHandler`.
   See: https://pkg.go.dev/github.com/elazarl/goproxy#ProxyHttpServer.NewConnectDialToProxyWithHandler
4. **`OnRequest().Do`**: add `Proxy-Authorization` to every plain HTTP
   request using the provided `setProxyAuth` helper.

## Check your understanding
1. What is a single-hop header, and why does `Proxy-Authorization` need to be
   re-added by the middle proxy?
2. What does `KeepHeader = true` change exactly? What other headers does it
   end up preserving?
3. For plain HTTP the middle proxy uses `Tr.Proxy`; for HTTPS it uses
   `ConnectDial`. Why are two different mechanisms needed?
4. Why is `ConnectDial` implemented with `NewConnectDialToProxyWithHandler`
   instead of `NewConnectDialToProxy`? (Hint: think about authentication.)
5. Where does the 407 come from when you hit the edge proxy without
   credentials? Is it produced by the middle proxy or the edge proxy?
6. What would happen if the middle proxy forwarded to the edge proxy but
   forgot to add the credentials?

## Bonus (optional)
- Change the credentials on both proxies and verify that the old ones stop
  working.
- Add a third proxy in the middle of the chain (extra hop) and verify the
  request still reaches the backend.
- Log on the edge proxy the requests it receives, to confirm the client IP
  is the middle proxy's, not the client's.

## Verify
Run the program:
```sh
go run .   # from exercises/03_cascade_proxy
```

> Tip: open `TEST.md` and run the curl commands it contains: the expected
> output is recorded there, so you can check whether your cascade matches it.

In another terminal:
```sh
# HTTP through the cascade
curl -v --proxy http://localhost:10998 http://localhost:11000/
# Expected: "Hello from backend, you requested /"

# HTTPS (CONNECT) through the cascade
curl -v --proxy http://localhost:10998 --cacert certs_localhost/localhost.crt \
  https://localhost:11001/
# Expected: "Hello from TLS backend, you requested /"

# Directly against the edge proxy, without credentials
curl -v --proxy http://localhost:10999 http://localhost:11000/
# Expected: 407 Proxy Authentication Required
```

## References
- goproxy: https://pkg.go.dev/github.com/elazarl/goproxy
- ext/auth: https://pkg.go.dev/github.com/elazarl/goproxy/ext/auth
- net/http Transport: https://pkg.go.dev/net/http#Transport
- crypto/subtle: https://pkg.go.dev/crypto/subtle
