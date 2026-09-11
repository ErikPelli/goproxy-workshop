package main

import (
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

const (
	// The "edge" proxy (the last one before the target) accepts these
	// credentials in the Proxy-Authorization header.
	_proxyUser     = "foo"
	_proxyPassword = "bar"
)

func main() {
	backendSrv, err := startBackend()
	if err != nil {
		log.Fatal(err)
	}
	defer backendSrv.Close()

	tlsBackendSrv, err := startTLSBackend()
	if err != nil {
		log.Fatal(err)
	}
	defer tlsBackendSrv.Close()

	// TODO(1): edge proxy
	// The edge proxy is the last proxy in the chain (locally :10999): it
	// authenticates every request before forwarding it to the target.
	// Register the authentication with auth.ProxyBasic using the
	// _proxyUser/_proxyPassword credentials:
	//   auth.ProxyBasic(endProxy, realm, func(user, passwd string) bool {...})
	// See: https://pkg.go.dev/github.com/elazarl/goproxy/ext/auth#ProxyBasic
	//
	// Tips:
	//   - compare user/passwd with _proxyUser/_proxyPassword using
	//     subtle.ConstantTimeCompare to avoid timing attacks
	//     (https://pkg.go.dev/crypto/subtle)
	endProxy := goproxy.NewProxyHttpServer()
	// TODO(1): auth.ProxyBasic(endProxy, ...)

	log.Println("edge proxy listening on localhost:10999")
	go func() {
		log.Fatal(http.ListenAndServe("localhost:10999", endProxy))
	}()

	// TODO(2): middle proxy
	// The middle proxy (locally :10998) is the first in the chain: it
	// receives the client requests and forwards them to the edge proxy,
	// adding the credentials. Complete each sub-step:
	//
	// 2a. KeepHeader
	// Set `middleProxy.KeepHeader = true` so that the Proxy-Authorization
	// header added below is not stripped before the forwarding (it is a
	// single-hop header by default).
	//
	// 2b. Forward plain HTTP requests
	// Set `middleProxy.Tr.Proxy` to return a *url.URL pointing to
	// http://localhost:10999. This is the standard http.Transport hook for
	// routing through an upstream proxy.
	// See: https://pkg.go.dev/net/http#Transport.Proxy
	//
	// 2c. Forward CONNECT (https://) requests
	//   middleProxy.ConnectDial = middleProxy.NewConnectDialToProxyWithHandler(
	//       "http://localhost:10999", func(req *http.Request) { setProxyAuth(req) })
	// See: https://pkg.go.dev/github.com/elazarl/goproxy#ProxyHttpServer.NewConnectDialToProxyWithHandler
	//
	// 2d. Inject the credentials on every plain HTTP request
	// Register an `OnRequest().Do(...)` handler that calls setProxyAuth(req).
	middleProxy := goproxy.NewProxyHttpServer()
	// TODO(2): complete 2a..2d

	log.Println("middle proxy listening on localhost:10998")
	log.Fatal(http.ListenAndServe("localhost:10998", middleProxy))
}

func setProxyAuth(req *http.Request) {
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(_proxyUser+":"+_proxyPassword)))
}
