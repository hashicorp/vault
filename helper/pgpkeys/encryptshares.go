package pgpkeys

import (
	"bytes"
	"encoding/base64"
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
func EncryptShares(secretShares *[][]byte, pgpKeys *[]string) error {
	if len(*secretShares) != len(*pgpKeys) {
		return fmt.Errorf("Mismatch between number of generated shares and number of PGP keys")
	}
	for i, keystring := range *pgpKeys {
		data, err := base64.StdEncoding.DecodeString(keystring)
		if err != nil {
			return fmt.Errorf("Error decoding given PGP key: %s", err)
		}
		entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(data)))
		if err != nil {
			return fmt.Errorf("Error parsing given PGP key: %s", err)
		}
		ctBuf := bytes.NewBuffer(nil)
		pt, err := openpgp.Encrypt(ctBuf, []*openpgp.Entity{entity}, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Error setting up encryption for PGP message: %s", err)
		}
		_, err = pt.Write((*secretShares)[i])
		if err != nil {
			return fmt.Errorf("Error encrypting PGP message: %s", err)
		}
		pt.Close()
		(*secretShares)[i] = ctBuf.Bytes()
	}
	return nil
}
