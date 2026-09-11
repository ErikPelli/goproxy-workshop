package main

import (
	"crypto/tls"
	"sync"
)

// CertStorage is a simple certificate cache that keeps
// everything in memory.
type CertStorage struct {
	certs map[string]*tls.Certificate
	mtx   sync.RWMutex
}

// TODO(1): implement Fetch
// Return the cached certificate for `hostname` if present, otherwise
// generate a new one with `gen()` and store it for future use.
// Remember to protect the map access with the mutex:
//   - read path: mtx.RLock / mtx.RUnlock
//   - write path: mtx.Lock / mtx.Unlock
func (cs *CertStorage) Fetch(hostname string, gen func() (*tls.Certificate, error)) (*tls.Certificate, error) {
	return nil, nil // TODO
}

func NewCertStorage() *CertStorage {
	return &CertStorage{
		certs: make(map[string]*tls.Certificate),
	}
}
