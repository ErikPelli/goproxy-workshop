package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/coder/websocket"
)

func startEchoBackend() (*httptest.Server, *http.Transport, error) {
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}
		defer func() {
			_ = c.Close(websocket.StatusNormalClosure, "")
		}()

		ctx := r.Context()
		for {
			mt, message, err := c.Read(ctx)
			if err != nil {
				break
			}
			log.Printf("backend received: %s", message)
			if err := c.Write(ctx, mt, append([]byte("ECHO: "), message...)); err != nil {
				break
			}
		}
	}))

	listener, err := net.Listen("tcp", ":11000")
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

	log.Printf("Websocket backend listening on %s", backend.Listener.Addr().String())
	return backend, transport, nil
}
