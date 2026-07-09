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

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	var apiCondition goproxy.ReqConditionFunc = func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		return req.Header.Get("X-Random") != ""
	}

	rateLimiter := httprate.NewRateLimiter(1, time.Minute)
	proxy.OnRequest(apiCondition).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ip, _, _ := strings.Cut(req.RemoteAddr, ":")
		if rateLimiter.OnLimit(mockResponseWriter{}, req, ip) {
			return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusTooManyRequests, "Too many requests")
		}
		return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "Correct response")
	})

	// Default handler when we dont match apiCondition filter
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusUnauthorized, "wrong request, retry")
	})

	proxy.OnResponse(apiCondition).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp.StatusCode != http.StatusOK {
			return resp
		}

		req, err := http.NewRequestWithContext(resp.Request.Context(), http.MethodGet, "http://localhost:11000", nil)
		if err != nil {
			return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusInternalServerError, err.Error())
		}

		respRand, err := http.DefaultClient.Do(req)
		if err != nil {
			return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusInternalServerError, err.Error())
		}
		defer func() {
			_ = respRand.Body.Close()
		}()

		respRandBody, err := io.ReadAll(respRand.Body)
		if err != nil {
			return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusInternalServerError, err.Error())
		}

		var result randomResponse
		if err := json.Unmarshal(respRandBody, &result); err != nil {
			return goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusInternalServerError, err.Error())
		}

		resp.Header.Set("X-Random", strconv.Itoa(result.Number))

		return resp
	})

	log.Fatal(http.ListenAndServe(":10999", proxy))
}

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
