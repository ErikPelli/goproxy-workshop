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
	// credentials in the Proxy-Authorization header
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

	// Edge proxy: the last proxy in the chain, it authenticates every request
	// before forwarding it to the target
	endProxy := goproxy.NewProxyHttpServer()
	auth.ProxyBasic(endProxy, "my_realm", func(user, passwd string) bool {
		return subtle.ConstantTimeCompare([]byte(user), []byte(_proxyUser)) == 1 &&
			subtle.ConstantTimeCompare([]byte(passwd), []byte(_proxyPassword)) == 1
	})

	log.Println("edge proxy listening on localhost:10999")
	go func() {
		log.Fatal(http.ListenAndServe("localhost:10999", endProxy))
	}()

	// Middle proxy: the first proxy in the chain, it forwards every request to
	// the edge proxy, adding the credentials
	middleProxy := goproxy.NewProxyHttpServer()

	// Keep the Proxy-Authorization header added below, so that it reaches the
	// edge proxy instead of being stripped as a single-hop header
	middleProxy.KeepHeader = true

	// Forward plain HTTP requests to the edge proxy
	middleProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return &url.URL{
			Scheme: "http",
			Host:   "localhost:10999",
		}, nil
	}

	// Forward CONNECT (https://) requests through the edge proxy, adding the
	// Proxy-Authorization header to the CONNECT request
	middleProxy.ConnectDial = middleProxy.NewConnectDialToProxyWithHandler("http://localhost:10999", func(req *http.Request) {
		setProxyAuth(req)
	})

	// Add the Proxy-Authorization header to every plain HTTP request
	middleProxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		setProxyAuth(req)
		return req, nil
	}))

	log.Println("middle proxy listening on localhost:10998")
	log.Fatal(http.ListenAndServe("localhost:10998", middleProxy))
}

func setProxyAuth(req *http.Request) {
	req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(_proxyUser+":"+_proxyPassword)))
}
