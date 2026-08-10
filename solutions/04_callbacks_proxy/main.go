package main

import (
	"crypto/sha512"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

const (
	_proxyUser = "test"

	// SHA-512 hash of the accepted password ("secret").
	// Generate it with: printf '%s' 'secret' | sha512sum
	_proxyPasswordHashHex = "bd2b1aaf7ef4f09be9f52ce2d8d599674d81aa9d6a4421696dc4d93dd0619d682ce56b4d64a9ef097761ced99e0f67265b5f76085e5b0ee7ca4696b2ad6fe2b2"
)

func main() {
	proxyPasswordHash, err := hex.DecodeString(_proxyPasswordHashHex)
	if err != nil {
		log.Fatalf("invalid password hash: %v", err)
	}

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

	proxy := goproxy.NewProxyHttpServer()

	// Verify Source IP before authenticating user
	// Allow only private IP addresses
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if !isPrivateIP(r) {
			return r, goproxy.NewResponse(r, "application/json", http.StatusForbidden, `{"message":"Invalid IP address"}`)
		}
		return r, nil
	})

	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		if !isPrivateIP(ctx.Req) {
			return goproxy.RejectConnect, host
		}
		return nil, host
	})

	// Verify Proxy-Authorization header
	auth.ProxyBasic(proxy, "Test", func(user, passwd string) bool {
		if user != _proxyUser {
			return false
		}

		hash := sha512.Sum512([]byte(passwd))
		return subtle.ConstantTimeCompare(proxyPasswordHash, hash[:]) == 1
	})

	// mTLS configuration: the proxy serves HTTPS and verifies the
	// client certificate against the custom CA.
	tlsConfig, err := createTLSConfig(
		"certs/server-cert.pem",
		"certs/server-key.pem",
		"certs/ca-cert.pem",
	)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:      ":10999",
		Handler:   proxy,
		TLSConfig: tlsConfig,
	}

	log.Printf("proxy server listening on %s (mTLS)", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func isPrivateIP(r *http.Request) bool {
	// Allow only private IP addresses
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(clientIP)
	return ip.IsPrivate() || ip.IsLoopback()
}

func createTLSConfig(cert, key, ca string) (*tls.Config, error) {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("unable to append CA certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{c},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
