# Test

## Hijack
```sh
curl -v --proxy http://localhost:10999 --cacert certs_localhost/localhost.crt --http1.1 https://localhost:11000/
```

Expected Response:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* CONNECT: no ALPN negotiated
* allocate connect buffer
* Establish HTTP proxy tunnel to localhost:11000
> CONNECT localhost:11000 HTTP/1.1
> Host: localhost:11000
> User-Agent: curl/8.18.0
> Proxy-Connection: Keep-Alive
> 
< HTTP/1.1 200 Ok
< 
* CONNECT phase completed
* CONNECT tunnel established, response 200
* ALPN: curl offers http/1.1
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* SSL Trust Anchors:
*   CAfile: certs_localhost/localhost.crt
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519MLKEM768 / RSASSA-PSS
* ALPN: server did not agree on a protocol. Uses default.
* Server certificate:
*   subject: CN=localhost
*   start date: Jul  7 21:03:28 2026 GMT
*   expire date: Aug  6 21:03:28 2026 GMT
*   issuer: CN=localhost
*   Certificate level 0: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 33554 
* using HTTP/1.x
> GET / HTTP/1.1
> Host: localhost:11000
> User-Agent: curl/8.18.0
> Accept: */*
> 
* Request completely sent off
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
< HTTP/1.1 200 OK
< Connection: close
< test: 1234
< 
* TLSv1.3 (IN), TLS alert, close notify (256):
* shutting down connection #0
hijacked response
```

## MITM
```sh
curl -v --proxy http://localhost:10999 --cacert certs/ca-cert.pem https://localhost:11001/
```

Expected Response:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* CONNECT: no ALPN negotiated
* allocate connect buffer
* Establish HTTP proxy tunnel to localhost:11001
> CONNECT localhost:11001 HTTP/1.1
> Host: localhost:11001
> User-Agent: curl/8.18.0
> Proxy-Connection: Keep-Alive
> 
< HTTP/1.0 200 OK
< 
* CONNECT phase completed
* CONNECT tunnel established, response 200
* ALPN: curl offers h2,http/1.1
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* SSL Trust Anchors:
*   CAfile: certs/ca-cert.pem
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519MLKEM768 / ED25519
* ALPN: server did not agree on a protocol. Uses default.
* Server certificate:
*   subject: O=GoProxy untrusted MITM proxy Inc; CN=localhost
*   start date: Jun  7 20:44:23 2026 GMT
*   expire date: Jul  7 20:44:23 2027 GMT
*   issuer: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=ca.goproxy.com
*   Certificate level 0: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   Certificate level 1: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 44690 
* using HTTP/1.x
> GET / HTTP/1.1
> Host: localhost:11001
> User-Agent: curl/8.18.0
> Accept: */*
> 
* Request completely sent off
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
< HTTP/1.1 200 OK
< Transfer-Encoding: chunked
< Content-Type: text/plain; charset=utf-8
< Date: Tue, 07 Jul 2026 20:44:23 GMT
< 
* Connection #0 to host localhost:10999 left intact
mitm body
```
