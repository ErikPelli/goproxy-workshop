# Test

## mTLS: client certificate required
```sh
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem http://localhost:11000/
```

Expected Response:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* ALPN: curl offers http/1.1
} [5 bytes data]
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
} [1565 bytes data]
* SSL Trust Anchors:
*   CAfile: certs/ca-cert.pem
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Server hello (2):
{ [1210 bytes data]
* TLSv1.3 (IN), TLS change cipher, Change cipher spec (1):
{ [1 bytes data]
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
{ [25 bytes data]
* TLSv1.3 (IN), TLS handshake, Request CERT (13):
{ [182 bytes data]
* TLSv1.3 (IN), TLS handshake, Certificate (11):
{ [541 bytes data]
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
{ [72 bytes data]
* TLSv1.3 (IN), TLS handshake, Finished (20):
{ [36 bytes data]
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
} [1 bytes data]
* TLSv1.3 (OUT), TLS handshake, Certificate (11):
} [8 bytes data]
* TLSv1.3 (OUT), TLS handshake, Finished (20):
} [36 bytes data]
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519MLKEM768 / ED25519
* ALPN: server accepted http/1.1
* Proxy certificate:
*   subject: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=workshop.goproxy.com
*   start date: Jul  7 13:45:00 2026 GMT
*   expire date: Jul  4 13:45:00 2036 GMT
*   issuer: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=ca.goproxy.com
*   Certificate level 0: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   Certificate level 1: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 45852 
* using HTTP/1.x
} [5 bytes data]
> GET http://localhost:11000/ HTTP/1.1
> Host: localhost:11000
> User-Agent: curl/8.18.0
> Accept: */*
> Proxy-Connection: Keep-Alive
> 
* Request completely sent off
{ [5 bytes data]
* TLSv1.3 (IN), TLS alert, unknown (628):
{ [2 bytes data]
* OpenSSL SSL_read: OpenSSL/3.5.7: error:0A00045C:SSL routines::tlsv13 alert certificate required, errno 0
* closing connection #0
```

## Authentication: missing credentials
```sh
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem \
  --proxy-cert certs/client-cert.pem --proxy-key certs/client-key.pem http://localhost:11000/
```

Expected Response:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* ALPN: curl offers http/1.1
} [5 bytes data]
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
} [1565 bytes data]
* SSL Trust Anchors:
*   CAfile: certs/ca-cert.pem
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Server hello (2):
{ [1210 bytes data]
* TLSv1.3 (IN), TLS change cipher, Change cipher spec (1):
{ [1 bytes data]
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
{ [25 bytes data]
* TLSv1.3 (IN), TLS handshake, Request CERT (13):
{ [182 bytes data]
* TLSv1.3 (IN), TLS handshake, Certificate (11):
{ [541 bytes data]
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
{ [72 bytes data]
* TLSv1.3 (IN), TLS handshake, Finished (20):
{ [36 bytes data]
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
} [1 bytes data]
* TLSv1.3 (OUT), TLS handshake, Certificate (11):
} [983 bytes data]
* TLSv1.3 (OUT), TLS handshake, CERT verify (15):
} [72 bytes data]
* TLSv1.3 (OUT), TLS handshake, Finished (20):
} [36 bytes data]
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519MLKEM768 / ED25519
* ALPN: server accepted http/1.1
* Proxy certificate:
*   subject: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=workshop.goproxy.com
*   start date: Jul  7 13:45:00 2026 GMT
*   expire date: Jul  4 13:45:00 2036 GMT
*   issuer: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=ca.goproxy.com
*   Certificate level 0: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   Certificate level 1: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 45866 
* using HTTP/1.x
} [5 bytes data]
> GET http://localhost:11000/ HTTP/1.1
> Host: localhost:11000
> User-Agent: curl/8.18.0
> Accept: */*
> Proxy-Connection: Keep-Alive
> 
* Request completely sent off
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
{ [1598 bytes data]
< HTTP/1.1 407 Proxy Authentication Required
< Proxy-Authenticate: Basic realm=Test
* HTTP/1.1 proxy connection set close
< Proxy-Connection: close
< Date: Mon, 10 Aug 2026 12:49:43 GMT
< Content-Length: 33
< Content-Type: text/plain; charset=utf-8
< 
{ [33 bytes data]
* shutting down connection #0
407 Proxy Authentication Required
```

## Successful request
```sh
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem \
  --proxy-cert certs/client-cert.pem --proxy-key certs/client-key.pem \
  --proxy-user test:secret http://localhost:11000/
```

Expected Response:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* ALPN: curl offers http/1.1
} [5 bytes data]
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
} [1565 bytes data]
* SSL Trust Anchors:
*   CAfile: certs/ca-cert.pem
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Server hello (2):
{ [1210 bytes data]
* TLSv1.3 (IN), TLS change cipher, Change cipher spec (1):
{ [1 bytes data]
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
{ [25 bytes data]
* TLSv1.3 (IN), TLS handshake, Request CERT (13):
{ [182 bytes data]
* TLSv1.3 (IN), TLS handshake, Certificate (11):
{ [541 bytes data]
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
{ [72 bytes data]
* TLSv1.3 (IN), TLS handshake, Finished (20):
{ [36 bytes data]
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
} [1 bytes data]
* TLSv1.3 (OUT), TLS handshake, Certificate (11):
} [983 bytes data]
* TLSv1.3 (OUT), TLS handshake, CERT verify (15):
} [72 bytes data]
* TLSv1.3 (OUT), TLS handshake, Finished (20):
} [36 bytes data]
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519MLKEM768 / ED25519
* ALPN: server accepted http/1.1
* Proxy certificate:
*   subject: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=workshop.goproxy.com
*   start date: Jul  7 13:45:00 2026 GMT
*   expire date: Jul  4 13:45:00 2036 GMT
*   issuer: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=ca.goproxy.com
*   Certificate level 0: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   Certificate level 1: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 45874 
* using HTTP/1.x
* Proxy auth using Basic with user 'test'
} [5 bytes data]
> GET http://localhost:11000/ HTTP/1.1
> Host: localhost:11000
> Proxy-Authorization: Basic dGVzdDpzZWNyZXQ=
> User-Agent: curl/8.18.0
> Accept: */*
> Proxy-Connection: Keep-Alive
> 
* Request completely sent off
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
{ [1598 bytes data]
< HTTP/1.1 200 OK
< Content-Length: 36
< Content-Type: text/plain; charset=utf-8
< Date: Mon, 10 Aug 2026 12:49:43 GMT
< 
{ [36 bytes data]
* Connection #0 to host localhost:10999 left intact
Hello from backend, you requested /
```

## CONNECT tunnel: HTTPS backend
```sh
curl -v --proxy https://localhost:10999 --proxy-cacert certs/ca-cert.pem \
  --proxy-cert certs/client-cert.pem --proxy-key certs/client-key.pem \
  --proxy-user test:secret --cacert certs_localhost/localhost.crt https://localhost:11001/
```

Expected Response:
```
* Host localhost:10999 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:10999...
* ALPN: curl offers http/1.1
} [5 bytes data]
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
} [1565 bytes data]
* SSL Trust Anchors:
*   CAfile: certs/ca-cert.pem
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Server hello (2):
{ [1210 bytes data]
* TLSv1.3 (IN), TLS change cipher, Change cipher spec (1):
{ [1 bytes data]
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
{ [25 bytes data]
* TLSv1.3 (IN), TLS handshake, Request CERT (13):
{ [182 bytes data]
* TLSv1.3 (IN), TLS handshake, Certificate (11):
{ [541 bytes data]
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
{ [72 bytes data]
* TLSv1.3 (IN), TLS handshake, Finished (20):
{ [36 bytes data]
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
} [1 bytes data]
* TLSv1.3 (OUT), TLS handshake, Certificate (11):
} [983 bytes data]
* TLSv1.3 (OUT), TLS handshake, CERT verify (15):
} [72 bytes data]
* TLSv1.3 (OUT), TLS handshake, Finished (20):
} [36 bytes data]
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256 / X25519MLKEM768 / ED25519
* ALPN: server accepted http/1.1
* Proxy certificate:
*   subject: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=workshop.goproxy.com
*   start date: Jul  7 13:45:00 2026 GMT
*   expire date: Jul  4 13:45:00 2036 GMT
*   issuer: C=IT; ST=Padua; L=Padua; O=GoProxy; OU=workshop; CN=ca.goproxy.com
*   Certificate level 0: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   Certificate level 1: Public key type ED25519 (256/128 Bits/secBits), signed using ED25519
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* CONNECT: 'http/1.1' negotiated
* allocate connect buffer
* Proxy auth using Basic with user 'test'
* Establish HTTP proxy tunnel to localhost:11001
} [5 bytes data]
> CONNECT localhost:11001 HTTP/1.1
> Host: localhost:11001
> Proxy-Authorization: Basic dGVzdDpzZWNyZXQ=
> User-Agent: curl/8.18.0
> Proxy-Connection: Keep-Alive
> 
{ [5 bytes data]
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
{ [1598 bytes data]
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
*   start date: Aug 10 13:12:48 2026 GMT
*   expire date: Aug  7 13:12:48 2036 GMT
*   issuer: CN=localhost
*   Certificate level 0: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
*   subjectAltName: "localhost" matches cert's "localhost"
* SSL certificate verified via OpenSSL.
* Established connection to localhost (::1 port 10999) from ::1 port 40460 
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
< Date: Mon, 10 Aug 2026 13:13:13 GMT
< Content-Length: 40
< Content-Type: text/plain; charset=utf-8
< 
{ [40 bytes data]
* Connection #0 to host localhost:10999 left intact
Hello from TLS backend, you requested /
```
