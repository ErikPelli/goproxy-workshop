package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/httptest"
)

type randomResponse struct {
	Number int `json:"number"`
}

func startRandomBackend() (*httptest.Server, error) {
	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Random number between 0 and 9
		n := rand.IntN(10)
		_ = json.NewEncoder(w).Encode(randomResponse{
			Number: n,
		})
	}))

	listener, err := net.Listen("tcp", ":11000")
	if err != nil {
		backend.Close()
		return nil, err
	}

	backend.Listener = listener
	backend.Start()

	log.Printf("random backend server listening on %s\n", backend.Listener.Addr().String())
	return backend, nil
}
