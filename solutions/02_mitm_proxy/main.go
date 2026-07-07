package main

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
)

func main() {
	// Hijack certificate
	localhostCert, err := os.ReadFile("certs_localhost/localhost.crt")
	if err != nil {
		log.Fatal(err)
	}

	localhostKey, err := os.ReadFile("certs_localhost/localhost.key")
	if err != nil {
		log.Fatal(err)
	}

	localhostCertificate, err := parseCertificate(localhostCert, localhostKey)
	if err != nil {
		log.Fatal(err)
	}

	// MITM CA certificate
	caCert, err := os.ReadFile("certs/ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	caKey, err := os.ReadFile("certs/ca-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	caCertificate, err := parseCertificate(caCert, caKey)
	if err != nil {
		log.Fatal(err)
	}

	hijackSrv, err := startHijackBackend()
	if err != nil {
		log.Fatal(err)
	}
	defer hijackSrv.Close()

	mitmSrv, mitmTr, err := startMitmBackend()
	if err != nil {
		log.Fatal(err)
	}
	defer mitmSrv.Close()

	proxy := goproxy.NewProxyHttpServer()

	// Instead of returning the Internet response, send custom data from
	// our proxy server, using connection hijack
	proxy.OnRequest(goproxy.ReqHostIs("localhost:11000")).HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
		log.Println("connection hijacked")

		if req.ProtoMajor != 1 {
			log.Println("I'm too lazy to implement HTTP/2 by hand in hijack")
			_ = client.Close()
			return
		}

		_, _ = client.Write([]byte("HTTP/1.1 200 Ok\r\n\r\n"))

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{*localhostCertificate},
		}
		tlsClient := tls.Server(client, tlsConfig)

		if err := tlsClient.Handshake(); err != nil {
			log.Println("TLS Handshake error:", err)
			_ = client.Close()
			return
		}

		w := bufio.NewWriter(tlsClient)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"test": {"1234"},
			},
			Body: io.NopCloser(strings.NewReader("hijacked response")),
		}

		_ = resp.Write(w)
		_ = w.Flush()
		_ = tlsClient.Close()
	})

	customCaMitm := &goproxy.ConnectAction{
		Action:    goproxy.ConnectMitm,
		TLSConfig: goproxy.TLSConfigFromCA(caCertificate),
	}
	var customAlwaysMitm goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return customCaMitm, host
	}

	proxy.OnRequest(goproxy.ReqHostIs("localhost:11001")).HandleConnect(customAlwaysMitm)
	proxy.OnRequest(goproxy.ReqHostIs("localhost:11001")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.RoundTripper = &roundTripper{tr: mitmTr}
		return req, nil
	})
	proxy.OnResponse(goproxy.ReqHostIs("localhost:11001")).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		log.Println("connection mitm'ed")
		resp.Body = io.NopCloser(strings.NewReader("mitm body"))
		return resp
	})

	log.Fatal(http.ListenAndServe(":10999", proxy))
}

func parseCertificate(caCert, caKey []byte) (*tls.Certificate, error) {
	parsedCert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return nil, err
	}
	return &parsedCert, nil
}

type roundTripper struct {
	tr http.RoundTripper
}

func (r *roundTripper) RoundTrip(req *http.Request, _ *goproxy.ProxyCtx) (*http.Response, error) {
	return r.tr.RoundTrip(req)
}
