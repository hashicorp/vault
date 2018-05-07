package cryptoutil

import "golang.org/x/crypto/blake2b"

func Blake2b256Hash(key string) ([]byte, error) {
	hf, err := blake2b.New256(nil)
	if err != nil {
		return nil, err
	}

	hf.Write([]byte(key))

	return hf.Sum(nil), nil
}
