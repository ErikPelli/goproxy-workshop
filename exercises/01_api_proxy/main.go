package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/go-chi/httprate"
)

func main() {
	randomSrv, err := startRandomBackend()
	if err != nil {
		log.Fatal(err)
	}
	defer randomSrv.Close()

	proxy := goproxy.NewProxyHttpServer()

	// The certificate cache avoids signing a new certificate on every
	// HTTPS request. `NewCertStorage` is implemented in cache.go.
	proxy.CertStore = NewCertStorage()

	// TODO(2): enable MITM
	// To inspect HTTPS traffic the proxy must perform a Man-in-the-Middle.
	// Register a handler that responds to any CONNECT request with
	// `goproxy.AlwaysMitm`:
	//   proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	// See: https://pkg.go.dev/github.com/elazarl/goproxy#AlwaysMitm

	// The condition selects only the API requests (header X-Random).
	var apiCondition goproxy.ReqConditionFunc = func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return req.Header.Get("X-Random") != ""
	}

	// TODO(3): rate limiting
	// Create a rate limiter with a limit of 1 request per minute:
	//   rateLimiter := httprate.NewRateLimiter(limit, window)
	// where `limit` is an int and `window` is a time.Duration.
	// See: https://pkg.go.dev/github.com/go-chi/httprate
	//
	// In the handler below, if `rateLimiter.OnLimit(...)` returns true,
	// respond with a 429 ("Too many requests"), otherwise with a 200
	// ("Correct response"). Build the response with goproxy.NewResponse:
	// https://pkg.go.dev/github.com/elazarl/goproxy#NewResponse
	proxy.OnRequest(apiCondition).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ip, _, _ := strings.Cut(req.RemoteAddr, ":")
		if rateLimiter.OnLimit(mockResponseWriter{}, req, ip) {
			return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusTooManyRequests, "Too many requests")
		}
		return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "Correct response")
	})

	// Default handler: requests that do not match apiCondition are rejected.
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusUnauthorized, "wrong request, retry")
	})

	// TODO(4): enrich the response
	// When an API response is 200, query the backend at http://localhost:11000,
	// decode the JSON `{"number": N}` and set the `X-Random` header on the
	// proxy response.
	// Suggestions:
	//   - build the request with http.NewRequestWithContext
	//   - execute it with http.DefaultClient.Do
	//   - read the body with io.ReadAll and decode it with json.Unmarshal
	// See: https://pkg.go.dev/net/http#NewRequestWithContext
	//      https://pkg.go.dev/encoding/json#Unmarshal
	proxy.OnResponse(apiCondition).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp.StatusCode != http.StatusOK {
			return resp
		}

		// TODO: fetch the random number from the backend and set the
		// X-Random header (use strconv.Itoa for the conversion).

		return resp
	})

	log.Fatal(http.ListenAndServe(":10999", proxy))
}

// mockResponseWriter mimics an http.ResponseWriter for the rate limiter.
type mockResponseWriter struct{}

func (m mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m mockResponseWriter) WriteHeader(_ int) {
	return
}
