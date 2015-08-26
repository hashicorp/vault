package pgpkeys

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

// EncryptShares takes an ordered set of Shamir key share fragments and
// PGP public keys and encrypts each Shamir key fragment with the corresponding
// public key
//
// Note: There is no corresponding test function; this functionality is
// thoroughly tested in the init and rekey command unit tests
func EncryptShares(secretShares [][]byte, pgpKeys []string) ([][]byte, error) {
	if len(secretShares) != len(pgpKeys) {
		return nil, fmt.Errorf("Mismatch between number of generated shares and number of PGP keys")
	}
	encryptedShares := [][]byte{}
	for i, keystring := range pgpKeys {
		data, err := base64.StdEncoding.DecodeString(keystring)
		if err != nil {
			return nil, fmt.Errorf("Error decoding given PGP key: %s", err)
		}
		entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(data)))
		if err != nil {
			return nil, fmt.Errorf("Error parsing given PGP key: %s", err)
		}
		ctBuf := bytes.NewBuffer(nil)
		pt, err := openpgp.Encrypt(ctBuf, []*openpgp.Entity{entity}, nil, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("Error setting up encryption for PGP message: %s", err)
		}
		_, err = pt.Write([]byte(hex.EncodeToString(secretShares[i])))
		if err != nil {
			return nil, fmt.Errorf("Error encrypting PGP message: %s", err)
		}
		pt.Close()
		encryptedShares = append(encryptedShares, ctBuf.Bytes())
	}
	return encryptedShares, nil
}
