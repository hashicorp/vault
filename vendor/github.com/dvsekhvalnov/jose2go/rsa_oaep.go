package jose

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"

	"github.com/dvsekhvalnov/jose2go/arrays"
)

// RS-AES using OAEP key management algorithm implementation
func init() {
	RegisterJwa(&RsaOaep{shaSizeBits: 1})
	RegisterJwa(&RsaOaep{shaSizeBits: 256})
	RegisterJwa(&RsaOaep{shaSizeBits: 384})
	RegisterJwa(&RsaOaep{shaSizeBits: 512})
}

type RsaOaep struct {
	shaSizeBits int
}

func (alg *RsaOaep) Name() string {
	switch alg.shaSizeBits {
	case 1:
		return RSA_OAEP
	case 256:
		return RSA_OAEP_256
	case 384:
		return RSA_OAEP_384
	default:
		return RSA_OAEP_512
	}
}

func (alg *RsaOaep) WrapNewKey(cekSizeBits int, key interface{}, header map[string]interface{}) (cek []byte, encryptedCek []byte, err error) {
	if pubKey, ok := key.(*rsa.PublicKey); ok {
		if cek, err = arrays.Random(cekSizeBits >> 3); err == nil {
			encryptedCek, err = rsa.EncryptOAEP(alg.sha(), rand.Reader, pubKey, cek, nil)
			return
		}

		return nil, nil, err
	}

	return nil, nil, errors.New("RsaOaep.WrapNewKey(): expected key to be '*rsa.PublicKey'")
}

func (alg *RsaOaep) Unwrap(encryptedCek []byte, key interface{}, cekSizeBits int, header map[string]interface{}) (cek []byte, err error) {
	if privKey, ok := key.(*rsa.PrivateKey); ok {
		return rsa.DecryptOAEP(alg.sha(), rand.Reader, privKey, encryptedCek, nil)
	}

	return nil, errors.New("RsaOaep.Unwrap(): expected key to be '*rsa.PrivateKey'")
}

func (alg *RsaOaep) sha() hash.Hash {
	switch alg.shaSizeBits {
	case 1:
		return sha1.New()
	case 256:
		return sha256.New()
	case 384:
		return sha512.New384()
	default:
		return sha512.New()
	}
}
