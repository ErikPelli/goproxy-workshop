# goproxy-workshop

A workshop on HTTP proxies in Go, based on [elazarl/goproxy](https://github.com/elazarl/goproxy).

## Layout
- `exercises/`: exercises to complete (scaffolding + README with instructions).
- `solutions/`: complete reference solutions.
- `certs/` and `certs_localhost/`: shared certificates; copy them into an
  exercise folder when the exercise README asks for it.

## Exercises
1. [API Proxy](exercises/01_api_proxy/README.md): MITM, rate limiting and response enrichment.
2. [MITM Proxy](exercises/02_mitm_proxy/README.md): connection hijack and MITM with a custom CA.
3. [Cascade Proxy](exercises/03_cascade_proxy/README.md): proxy chain (middle → edge) with authentication.
4. [Callbacks Proxy](exercises/04_callbacks_proxy/README.md): goproxy callbacks, mTLS and source IP restriction.
5. [WebSocket Proxy](exercises/05_websocket_proxy/README.md): WebSocket proxy that logs the handshake.

## How to do an exercise
Each folder in `exercises/` contains a `README.md` with the steps to complete
and the related scaffolding, plus a `TEST.md` with the expected output of the
test commands: run them against your solution and check that the output
matches.

When an exercise needs certificates, you can either **generate your own** with
`openssl` (instructions in the `generation.md` files) or copy the provided
ones from the main project folder (see the exercise README).

## Requirements
- Go 1.24+
- Dependencies are already in the root `go.mod`.
