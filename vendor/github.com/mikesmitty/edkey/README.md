# edkey
edkey allows you to marshal/write ED25519 private keys in the OpenSSH private key format

## Example
```go
package main

import (
	"crypto/rand"
	"encoding/pem"
	"io/ioutil"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

func main() {
	// Generate a new private/public keypair for OpenSSH
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	publicKey, _ := ssh.NewPublicKey(pubKey)

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	}
	privateKey := pem.EncodeToMemory(pemKey)
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)

	_ = ioutil.WriteFile("id_ed25519", privateKey, 0600)
	_ = ioutil.WriteFile("id_ed25519.pub", authorizedKey, 0644)
}
```
