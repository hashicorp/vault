To rebuild the cert.pem within this folder run the following commands

```shell
$ openssl x509 -in cert.pem -signkey key.pem -x509toreq -out cert.csr
$ openssl x509 -req -in cert.csr -CA ../root/rootcacert.pem -CAkey ../root/rootcakey.pem -CAcreateserial -out cert.pem -days 9132 -sha256 -extensions v3_req -extfile <(echo "[v3_req]\nsubjectAltName=DNS:cert.example.com,IP:127.0.0.1")
```
