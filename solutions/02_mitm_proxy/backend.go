package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
)

func startHijackBackend() (*httptest.Server, error) {
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Hijacked backend should not be called")
	}))

	listener, err := net.Listen("tcp", ":11000")
	if err != nil {
		backend.Close()
		return nil, err
	}

	backend.Listener = listener
	backend.Start()

	log.Printf("hijack backend server listening on %s\n", backend.Listener.Addr().String())
	return backend, nil
}

func startMitmBackend() (*httptest.Server, *http.Transport, error) {
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "original data")
	}))

	listener, err := net.Listen("tcp", ":11001")
	if err != nil {
		backend.Close()
		return nil, nil, err
	}

	backend.Listener = listener
	backend.StartTLS()

	// httptest generates certificates valid only for example.com, not localhost,
	// so we configure the Go transport to accept it even if the request is for
	// the host "localhost"
	transport := backend.Client().Transport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		ServerName: "example.com",
		RootCAs:    transport.TLSClientConfig.RootCAs,
	}

	log.Printf("mitm backend server listening on %s\n", backend.Listener.Addr().String())
	return backend, transport, nil
}
