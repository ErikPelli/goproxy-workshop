# Certification Authority
Generate the CA private key
```
$ openssl genpkey \
  -algorithm ed25519 \
  -out ca-key.pem
```

Generate the CA certificate (valid for 10 years)
```
$ openssl req \
  -key ca-key.pem \
  -new -x509 \
  -days 3650 \
  -out ca-cert.pem \
  -subj "/C=IT/ST=Padua/L=Padua/O=GoProxy/OU=workshop/CN=ca.goproxy.com"
```

# Server
Generate the server private key
```
$ openssl genpkey \
  -algorithm ed25519 \
  -out server-key.pem
```

Generate the server certificate request
```
$ openssl req \
  -new \
  -key server-key.pem \
  -out server-req.pem \
  -subj "/C=IT/ST=Padua/L=Padua/O=GoProxy/OU=workshop/CN=workshop.goproxy.com" \
  -addext "subjectAltName=DNS:workshop.goproxy.com,DNS:localhost"
```

Sign the certificate request with the CA
```
$ openssl x509 \
  -req \
  -in server-req.pem \
  -CA ca-cert.pem \
  -CAkey ca-key.pem \
  -CAcreateserial \
  -days 3650 \
  -out server-cert.pem \
  -extfile <(printf "subjectAltName=DNS:workshop.goproxy.com,DNS:localhost")
```

Remove the server certificate request
```
$ rm server-req.pem ca-cert.srl
```

# Client
Generate the node private key
```
$ openssl genpkey \
  -algorithm ed25519 \
  -out client-key.pem
```

Generate the client certificate request
```
$ openssl req \
  -new \
  -key client-key.pem \
  -out client-req.pem \
  -subj "/C=IT/ST=Padua/L=Padua/O=GoProxy/OU=workshop/CN=client"
```

Sign the certificate request with the CA
```
$ openssl x509 \
  -req \
  -in client-req.pem \
  -CA ca-cert.pem \
  -CAkey ca-key.pem \
  -CAcreateserial \
  -days 3650 \
  -out client-cert.pem
```

Remove the server certificate request
```
$ rm client-req.pem ca-cert.srl
```
