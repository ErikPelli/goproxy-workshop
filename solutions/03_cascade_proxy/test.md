# Test

## HTTP through the cascade
```sh
curl -v --proxy http://localhost:10998 http://localhost:11000/
```

Expected Response:
```
* Host localhost:10998 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10998...
* connect to ::1 port 10998 from ::1 port 53260 failed: Connection refused
*   Trying 127.0.0.1:10998...
* Established connection to localhost (127.0.0.1 port 10998) from 127.0.0.1 port 57474 
* using HTTP/1.x
> GET http://localhost:11000/ HTTP/1.1
> Host: localhost:11000
> User-Agent: curl/8.18.0
> Accept: */*
> Proxy-Connection: Keep-Alive
> 
* Request completely sent off
< HTTP/1.1 200 OK
< Content-Length: 36
< Content-Type: text/plain; charset=utf-8
< Date: Mon, 10 Aug 2026 13:32:17 GMT
< 
{ [36 bytes data]
* Connection #0 to host localhost:10998 left intact
Hello from backend, you requested /
```

## HTTPS (CONNECT) through the cascade
```sh
curl -v --proxy http://localhost:10998 --cacert certs_localhost/localhost.crt https://localhost:11001/
```

Expected Response:
```
* Host localhost:10998 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10998...
* connect to ::1 port 10998 from ::1 port 53274 failed: Connection refused
*   Trying 127.0.0.1:10998...
* CONNECT: no ALPN negotiated
* allocate connect buffer
* Establish HTTP proxy tunnel to localhost:11001
> CONNECT localhost:11001 HTTP/1.1
> Host: localhost:11001
> User-Agent: curl/8.18.0
> Proxy-Connection: Keep-Alive
> 
< HTTP/1.0 200 Connection established
< 
* CONNECT phase completed
* CONNECT tunnel established, response 200
* ALPN: curl offers h2,http/1.1
} [5 bytes data]
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
} [1568 bytes data]
* SSL Trust Anchors:
*   CAfile: certs_localhost/localhost.crt
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Server hello (2):
{ [1210 bytes data]
* TLSv1.3 (IN), TLS change cipher, Change cipher spec (1):
{ [1 bytes data]
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
{ [25 bytes data]
* TLSv1.3 (IN), TLS handshake, Certificate (11):
{ [800 bytes data]
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
{ [264 bytes data]
* TLSv1.3 (IN), TLS handshake, Finished (20):
{ [36 bytes data]
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
} [1 bytes data]
* TLSv1.3 (OUT), TLS handshake, Finished (20):
} [36 bytes data]
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519MLKEM768 / RSASSA-PSS
* ALPN: server accepted http/1.1
* Server certificate:
*   subject: CN=localhost
*   start date: Aug 10 13:18:42 2026 GMT
*   expire date: Aug  7 13:18:42 2036 GMT
*   issuer: CN=localhost
*   Certificate level 0: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (127.0.0.1 port 10998) from 127.0.0.1 port 57478 
* using HTTP/1.x
} [5 bytes data]
> GET / HTTP/1.1
> Host: localhost:11001
> User-Agent: curl/8.18.0
> Accept: */*
> 
* Request completely sent off
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
{ [122 bytes data]
< HTTP/1.1 200 OK
< Date: Mon, 10 Aug 2026 13:32:17 GMT
< Content-Length: 40
< Content-Type: text/plain; charset=utf-8
< 
{ [40 bytes data]
* Connection #0 to host localhost:10998 left intact
Hello from TLS backend, you requested /
```

## Edge proxy authentication
Directly against the edge proxy, without credentials, the request is rejected:
```sh
curl -v --proxy http://localhost:10999 http://localhost:11000/
```

Expected Response:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* connect to ::1 port 10999 from ::1 port 42890 failed: Connection refused
*   Trying 127.0.0.1:10999...
* Established connection to localhost (127.0.0.1 port 10999) from 127.0.0.1 port 38086 
* using HTTP/1.x
> GET http://localhost:11000/ HTTP/1.1
> Host: localhost:11000
> User-Agent: curl/8.18.0
> Accept: */*
> Proxy-Connection: Keep-Alive
> 
* Request completely sent off
< HTTP/1.1 407 Proxy Authentication Required
< Proxy-Authenticate: Basic realm=my_realm
* HTTP/1.1 proxy connection set close
< Proxy-Connection: close
< Date: Mon, 10 Aug 2026 13:32:17 GMT
< Content-Length: 33
< Content-Type: text/plain; charset=utf-8
< 
{ [33 bytes data]
* shutting down connection #0
407 Proxy Authentication Required
```

With the credentials injected by the middle proxy, the cascade works (see the tests above).
