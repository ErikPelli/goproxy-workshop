# Test

## Random
```sh
curl -v --proxy http://localhost:10999 --cacert goproxy_ca.pem --header "X-Random: true" https://fake-url.com/
```

Expected Response #1:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* CONNECT: no ALPN negotiated
* allocate connect buffer
* Establish HTTP proxy tunnel to fake-url.com:443
> CONNECT fake-url.com:443 HTTP/1.1
> Host: fake-url.com:443
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
*   CAfile: goproxy_ca.pem
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
*   subject: O=GoProxy untrusted MITM proxy Inc; CN=fake-url.com
*   start date: Jun  9 16:32:09 2026 GMT
*   expire date: Jul  9 16:32:09 2027 GMT
*   issuer: C=IL; ST=Center; L=Lod; O=GoProxy; OU=GoProxy; CN=goproxy.github.io; emailAddress=elazarl@gmail.com
*   Certificate level 0: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
*   Certificate level 1: Public key type RSA (4096/152 Bits/secBits), signed using sha256WithRSAEncryption
*   subjectAltName: "fake-url.com" matches cert's "fake-url.com"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 36290 
* using HTTP/1.x
> GET / HTTP/1.1
> Host: fake-url.com
> User-Agent: curl/8.18.0
> Accept: */*
> X-Random: true
> 
* Request completely sent off
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
< HTTP/1.1 200 OK
< Transfer-Encoding: chunked
< Content-Type: text/plain
< X-Random: 8
< 
* Connection #0 to host localhost:10999 left intact
Correct response
```

Expected Response #2:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* CONNECT: no ALPN negotiated
* allocate connect buffer
* Establish HTTP proxy tunnel to fake-url.com:443
> CONNECT fake-url.com:443 HTTP/1.1
> Host: fake-url.com:443
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
*   CAfile: goproxy_ca.pem
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
*   subject: O=GoProxy untrusted MITM proxy Inc; CN=fake-url.com
*   start date: Jun  9 16:32:14 2026 GMT
*   expire date: Jul  9 16:32:14 2027 GMT
*   issuer: C=IL; ST=Center; L=Lod; O=GoProxy; OU=GoProxy; CN=goproxy.github.io; emailAddress=elazarl@gmail.com
*   Certificate level 0: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
*   Certificate level 1: Public key type RSA (4096/152 Bits/secBits), signed using sha256WithRSAEncryption
*   subjectAltName: "fake-url.com" matches cert's "fake-url.com"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 55006 
* using HTTP/1.x
> GET / HTTP/1.1
> Host: fake-url.com
> User-Agent: curl/8.18.0
> Accept: */*
> X-Random: true
> 
* Request completely sent off
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
< HTTP/1.1 429 Too Many Requests
< Content-Length: 17
< Content-Type: text/plain
< 
* Connection #0 to host localhost:10999 left intact
Too many requests
```
