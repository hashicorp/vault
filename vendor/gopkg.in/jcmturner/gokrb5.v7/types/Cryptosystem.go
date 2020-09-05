package types

import (
	"github.com/jcmturner/gofork/encoding/asn1"
)

// Reference: https://www.ietf.org/rfc/rfc4120.txt
// Section: 5.2.9

// EncryptedData implements RFC 4120 type: https://tools.ietf.org/html/rfc4120#section-5.2.9
type EncryptedData struct {
	EType  int32  `asn1:"explicit,tag:0"`
	KVNO   int    `asn1:"explicit,optional,tag:1"`
	Cipher []byte `asn1:"explicit,tag:2"`
}

// EncryptionKey implements RFC 4120 type: https://tools.ietf.org/html/rfc4120#section-5.2.9
// AKA KeyBlock
type EncryptionKey struct {
	KeyType  int32  `asn1:"explicit,tag:0"`
	KeyValue []byte `asn1:"explicit,tag:1"`
}

// Checksum implements RFC 4120 type: https://tools.ietf.org/html/rfc4120#section-5.2.9
type Checksum struct {
	CksumType int32  `asn1:"explicit,tag:0"`
	Checksum  []byte `asn1:"explicit,tag:1"`
}

// Unmarshal bytes into the EncryptedData.
func (a *EncryptedData) Unmarshal(b []byte) error {
	_, err := asn1.Unmarshal(b, a)
	return err
}

// Marshal the EncryptedData.
func (a *EncryptedData) Marshal() ([]byte, error) {
	edb, err := asn1.Marshal(*a)
	if err != nil {
		return edb, err
	}
	return edb, nil
}

// Unmarshal bytes into the EncryptionKey.
func (a *EncryptionKey) Unmarshal(b []byte) error {
	_, err := asn1.Unmarshal(b, a)
	return err
}

// Unmarshal bytes into the Checksum.
func (a *Checksum) Unmarshal(b []byte) error {
	_, err := asn1.Unmarshal(b, a)
	return err
}
