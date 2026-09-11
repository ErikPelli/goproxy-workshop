# Exercise 04 - Callbacks Proxy (mTLS)

> Estimated time: ~30 minutes.

## Goal
Build a proxy that, using goproxy's callbacks, enforces three controls:
1. **Private IPs only**: requests from non-private IPs are rejected.
2. **Basic authentication**: a valid `Proxy-Authorization` header is
   required.
3. **mTLS**: the proxy serves HTTPS and verifies the client certificate
   against a custom CA.

## Learning outcomes
- Register callbacks for plain requests (`OnRequest().DoFunc`) and for
  CONNECT requests (`OnRequest().HandleConnectFunc`).
- Configure Basic proxy authentication with `auth.ProxyBasic`.
- Set up a mutual-TLS (mTLS) HTTP server.
- Store passwords as hashes and compare them in constant time.

## Background (read before coding, ~5 min)
- **Callbacks**: goproxy exposes several hook points. In this exercise you
  will use:
  - `proxy.OnRequest().DoFunc(...)` - runs before a plain HTTP request is
    forwarded.
  - `proxy.OnRequest().HandleConnectFunc(...)` - runs when a CONNECT arrives.
  - `auth.ProxyBasic(proxy, realm, check)` - registers both handlers above,
    returning `407` when the check fails.
- **mTLS**: the server presents its certificate AND asks the client for one
  (`tls.RequireAndVerifyClientCert`). The client certificate is verified
  against the `ClientCAs` pool.
  See: https://pkg.go.dev/crypto/tls#Config
- **Password storage**: never store plain passwords. Store an SHA-512 hash
  and compare with `subtle.ConstantTimeCompare` to avoid timing attacks.
  See: https://pkg.go.dev/crypto/subtle

## Prerequisites
- Go 1.24+ (https://go.dev/doc/install)

### Certificates
The proxy uses the `certs/` folder (CA, server certificate, client
certificate).

You have two options:
- **Generate your own** with `openssl` (you'll likely already have it):
  follow the instructions in `certs/generation.md`.
- **Copy the provided ones** from the main project folder:
  ```sh
  cp -r ../../certs .
  ```

## Project layout
```
exercises/04_callbacks_proxy/
├── main.go        <- to complete
├── backend.go     <- ready: HTTP backend (:11000) and HTTPS backend (:11001)
├── TEST.md        <- expected output of the curl tests
└── certs/         <- generate or copy from the main folder
```

> Note: `main.go` compiles only once all `TODO`s are completed; the Go
> compiler tells you exactly what is missing.

## Steps

### Step 1 - Check the source IP (plain requests)
Register a callback that rejects requests from non-private IPs:
```go
proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if !isPrivateIP(r) {
		return r, goproxy.NewResponse(r, "application/json", http.StatusForbidden, `{"message":"Invalid IP address"}`)
	}
	return r, nil
})
```
`isPrivateIP` is already provided (uses `net.IP.IsPrivate()` and
`net.IP.IsLoopback()`).
See: https://pkg.go.dev/net#IP.IsPrivate

### Step 2 - Check the source IP (CONNECT)
CONNECT goes through a different callback. If the IP is not private, reject
with `goproxy.RejectConnect`, otherwise return `nil` so the default action
applies.
See: https://pkg.go.dev/github.com/elazarl/goproxy#HttpsHandler

### Step 3 - Authentication
Require a valid `Proxy-Authorization` header:
```go
auth.ProxyBasic(proxy, "Test", func(user, passwd string) bool {
	if user != _proxyUser {
		return false
	}
	hash := sha512.Sum512([]byte(passwd))
	return subtle.ConstantTimeCompare(proxyPasswordHash, hash[:]) == 1
})
```
- Credentials: user `test`, password `secret` (stored only as a SHA-512 hash).
- `auth.ProxyBasic`: https://pkg.go.dev/github.com/elazarl/goproxy/ext/auth#ProxyBasic
- `crypto/sha512`: https://pkg.go.dev/crypto/sha512

### Step 4 - mTLS server
Create an `http.Server` and start it with TLS:
```go
server := &http.Server{
	Addr:      ":10999",
	Handler:   proxy,
	TLSConfig: tlsConfig,
}
log.Fatal(server.ListenAndServeTLS("", ""))
```
`createTLSConfig` is already provided: it loads the server certificate, the
CA, and sets `ClientAuth: tls.RequireAndVerifyClientCert`.

## Check your understanding
1. Why does `OnRequest().DoFunc` return a `*http.Response`? What does a
   non-nil response mean for the request flow?
2. Why is the source-IP check duplicated for plain requests and CONNECT?
   Could a CONNECT bypass the `OnRequest().DoFunc` callback?
3. `auth.ProxyBasic` registers two handlers. Which ones, and for which kinds
   of traffic?
4. Why store only the hash of the password instead of the plain text?
5. What does `tls.RequireAndVerifyClientCert` change compared to
   `tls.VerifyClientCertIfGiven`?
6. When curl connects to `https://localhost:10999`, which certificate is it
   verifying with `--proxy-cacert`? And which one does it present with
   `--proxy-cert`?

## Bonus (optional)
- Generate a new hash for a different password and verify that the old
  password stops working.
- Log the client certificate subject (`ctx.Req.TLS.PeerCertificates`) on each
  request.
- Add a callback that rejects requests coming from the `fe80::/10` range too.

## Verify
Run the proxy:
```sh
go run .   # from exercises/04_callbacks_proxy
```

> Tip: open `TEST.md` and run the curl commands it contains: the expected
> output is recorded there, so you can check whether your proxy matches it.

In another terminal:
```sh
# mTLS: without a client certificate the handshake fails
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem \
  http://localhost:11000/

# Client certificate but no credentials -> 407
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem \
  --proxy-cert certs/client-cert.pem --proxy-key certs/client-key.pem \
  http://localhost:11000/

# Client certificate and credentials -> 200
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem \
  --proxy-cert certs/client-cert.pem --proxy-key certs/client-key.pem \
  --proxy-user test:secret http://localhost:11000/

# CONNECT (HTTPS) -> 200
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem \
  --proxy-cert certs/client-cert.pem --proxy-key certs/client-key.pem \
  --proxy-user test:secret --cacert certs_localhost/localhost.crt \
  https://localhost:11001/
```

> Note: the CONNECT test also needs `certs_localhost/` (copy it from the main
> folder as in exercise 02, or use `-k`).

## References
- goproxy: https://pkg.go.dev/github.com/elazarl/goproxy
- ext/auth: https://pkg.go.dev/github.com/elazarl/goproxy/ext/auth
- crypto/tls: https://pkg.go.dev/crypto/tls
- crypto/sha512: https://pkg.go.dev/crypto/sha512
- crypto/subtle: https://pkg.go.dev/crypto/subtle
