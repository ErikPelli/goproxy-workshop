package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
)

func startBackend() (*httptest.Server, error) {
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello from backend, you requested %s\n", r.URL.Path)
	}))

	listener, err := net.Listen("tcp", ":11000")
	if err != nil {
		backend.Close()
		return nil, err
	}

	backend.Listener = listener
	backend.Start()

	log.Printf("backend server listening on %s\n", backend.Listener.Addr().String())
	return backend, nil
}

func startTLSBackend() (*httptest.Server, error) {
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello from TLS backend, you requested %s\n", r.URL.Path)
	}))

	listener, err := net.Listen("tcp", ":11001")
	if err != nil {
		backend.Close()
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair("certs_localhost/localhost.crt", "certs_localhost/localhost.key")
	if err != nil {
		_ = listener.Close()
		backend.Close()
		return nil, err
	}

	backend.Listener = listener
	backend.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	backend.StartTLS()

	log.Printf("TLS backend server listening on %s\n", backend.Listener.Addr().String())
	return backend, nil
}
