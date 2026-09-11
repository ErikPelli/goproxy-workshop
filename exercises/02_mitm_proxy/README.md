# Exercise 02 - MITM Proxy

> Estimated time: ~30 minutes.

## Goal
Build a proxy that does two different things on two ports:
1. **Hijack** (`localhost:11000`): steal the CONNECT connection and respond
   by itself with custom data, without contacting the real server.
2. **MITM with a custom CA** (`localhost:11001`): intercept the TLS using
   certificates generated on the fly and signed by our CA, modifying the
   response coming from the backend.

## Learning outcomes
- Understand TLS certificates, CAs and how goproxy signs per-host
  certificates.
- Use `HijackConnect` to take full control of a CONNECT connection.
- Configure MITM with a custom CA instead of the default one.
- Swap the upstream transport and rewrite a response body.

## Background (read before coding, ~5 min)
- **TLS handshake**: the server presents a certificate signed by a CA; the
  client verifies it against its trust store. goproxy can generate a
  certificate for any hostname, signed by a CA we control, so that a client
  trusting that CA accepts the proxy as the "real" server.
- **Hijack**: `HijackConnect` gives the proxy the raw TCP connection of the
  CONNECT request. The proxy can answer manually (even fake a TLS server)
  instead of tunnelling to the real target.
- **`goproxy.ConnectAction`**: tells goproxy what to do with a CONNECT:
  `ConnectAccept` (tunnel), `ConnectMitm` (intercept), `ConnectReject`.
  See: https://pkg.go.dev/github.com/elazarl/goproxy#ConnectAction

## Prerequisites
- Go 1.24+ (https://go.dev/doc/install)

### Certificates
The proxy uses two sets of certificates: `certs/` (the CA for the MITM
certificates) and `certs_localhost/` (the certificate for the hijack).

You have two options:
- **Generate your own** with `openssl` (you'll likely already have it):
  follow the instructions in `certs/generation.md` and
  `certs_localhost/generation.md`.
- **Copy the provided ones** from the main project folder:
  ```sh
  cp -r ../../certs .
  cp -r ../../certs_localhost .
  ```

## Project layout
```
exercises/02_mitm_proxy/
├── main.go            <- to complete
├── backend.go         <- ready: hijack backend (:11000) and MITM backend (:11001)
├── TEST.md            <- expected output of the curl tests
├── certs/             <- generate or copy from the main folder (CA for MITM)
└── certs_localhost/   <- generate or copy from the main folder (localhost certificate)
```

> Note: `main.go` compiles only once all `TODO`s are completed; the Go
> compiler tells you exactly what is missing.

## Steps

### Step 1 - Understand the backends (`backend.go`)
Read `backend.go`:
- `startHijackBackend`: binds `:11000` but panics if called. The proxy must
  never forward to it - it must hijack instead.
- `startMitmBackend`: binds `:11001`, serves "original data" over TLS, and
  returns a custom `*http.Transport` that trusts the httptest certificate
  (`ServerName: "example.com"`). Why is that transport needed? Because
  `httptest` generates certificates valid only for `example.com`, not
  `localhost`.
See: https://pkg.go.dev/net/http/httptest

### Step 2 - Connection hijack
Register the hijack handler for `localhost:11000`:
```go
proxy.OnRequest(goproxy.ReqHostIs("localhost:11000")).HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) { ... })
```
Inside:
- Write `HTTP/1.1 200 Ok\r\n\r\n` to the client.
- Wrap `client` in `tls.Server(client, tlsConfig)` using
  `localhostCertificate`, and call `Handshake()`.
- Write an `http.Response` with body `hijacked response` and header
  `test: 1234` using `bufio.NewWriter` and `resp.Write(w)`.

References:
- https://pkg.go.dev/crypto/tls#Server
- https://pkg.go.dev/net/http#Response.Write
- https://pkg.go.dev/net/http#Hijacker

### Step 3 - MITM with a custom CA
For `localhost:11001`, MITM the connection with certificates signed by the
workshop CA:
```go
customCaMitm := &goproxy.ConnectAction{
	Action:    goproxy.ConnectMitm,
	TLSConfig: goproxy.TLSConfigFromCA(caCertificate),
}
proxy.OnRequest(goproxy.ReqHostIs("localhost:11001")).HandleConnect(customAlwaysMitm)
```
- `goproxy.TLSConfigFromCA`: https://pkg.go.dev/github.com/elazarl/goproxy#TLSConfigFromCA

### Step 4 - Custom upstream transport
Once the connection is MITM'd, goproxy forwards the request to the backend
with its default transport, which would reject the `example.com` certificate.
Set `ctx.RoundTripper = &roundTripper{tr: mitmTr}` in an `OnRequest` handler
for `localhost:11001`. The `roundTripper` type is already at the bottom of
`main.go`.

### Step 5 - Modify the response
In an `OnResponse` handler for `localhost:11001`, replace the body:
```go
resp.Body = io.NopCloser(strings.NewReader("mitm body"))
```
- https://pkg.go.dev/io#NopCloser

## Check your understanding
1. Why can't the client talk to `localhost:11000` directly without the proxy
   seeing the data, even though it is local?
2. In the hijack, what happens if the client uses HTTP/2? Why does the code
   close the connection instead?
3. What is the difference between `goproxy.AlwaysMitm` and a custom
   `ConnectAction` with `TLSConfigFromCA`?
4. Why must the proxy use `mitmTr` (ServerName example.com) to talk to the
   `httptest` backend, while the client trusts a different CA?
5. `OnResponse` runs on the MITM'd response. In what order do `OnRequest`,
   the upstream round-trip and `OnResponse` run?
6. Could you reuse the same hijack technique for the MITM port? What would
   be missing compared to `ConnectMitm`?

## Bonus (optional)
- Change the hijacked body and header, and re-run curl to see them.
- Add logging of the MITM'd request headers (hint: `req.Header` in an
  `OnRequest` handler).
- Make the hijack also answer to a second path instead of closing.

## Verify
Run the proxy:
```sh
go run .   # from exercises/02_mitm_proxy
```

> Tip: open `TEST.md` and run the curl commands it contains: the expected
> output is recorded there, so you can check whether your proxy matches it.

In another terminal:
```sh
# Hijack
curl -v --proxy http://localhost:10999 --cacert certs_localhost/localhost.crt \
  --http1.1 https://localhost:11000/
# Expected: body "hijacked response", header "test: 1234"

# MITM
curl -v --proxy http://localhost:10999 --cacert certs/ca-cert.pem https://localhost:11001/
# Expected: body "mitm body"
```

## References
- goproxy: https://pkg.go.dev/github.com/elazarl/goproxy
- crypto/tls: https://pkg.go.dev/crypto/tls
- net/http: https://pkg.go.dev/net/http
