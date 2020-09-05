// Package types provides Kerberos 5 data types.
package types

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/jcmturner/gofork/encoding/asn1"
	"gopkg.in/jcmturner/gokrb5.v7/asn1tools"
	"gopkg.in/jcmturner/gokrb5.v7/iana"
	"gopkg.in/jcmturner/gokrb5.v7/iana/asnAppTag"
)

/*Authenticator   ::= [APPLICATION 2] SEQUENCE  {
authenticator-vno       [0] INTEGER (5),
crealm                  [1] Realm,
cname                   [2] PrincipalName,
cksum                   [3] Checksum OPTIONAL,
cusec                   [4] Microseconds,
ctime                   [5] KerberosTime,
subkey                  [6] EncryptionKey OPTIONAL,
seq-number              [7] UInt32 OPTIONAL,
authorization-data      [8] AuthorizationData OPTIONAL
}

   cksum
      This field contains a checksum of the application data that
      accompanies the KRB_AP_REQ, computed using a key usage value of 10
      in normal application exchanges, or 6 when used in the TGS-REQ
      PA-TGS-REQ AP-DATA field.

*/

// Authenticator - A record containing information that can be shown to have been recently generated using the session key known only by the client and server.
// https://tools.ietf.org/html/rfc4120#section-5.5.1
type Authenticator struct {
	AVNO              int               `asn1:"explicit,tag:0"`
	CRealm            string            `asn1:"generalstring,explicit,tag:1"`
	CName             PrincipalName     `asn1:"explicit,tag:2"`
	Cksum             Checksum          `asn1:"explicit,optional,tag:3"`
	Cusec             int               `asn1:"explicit,tag:4"`
	CTime             time.Time         `asn1:"generalized,explicit,tag:5"`
	SubKey            EncryptionKey     `asn1:"explicit,optional,tag:6"`
	SeqNumber         int64             `asn1:"explicit,optional,tag:7"`
	AuthorizationData AuthorizationData `asn1:"explicit,optional,tag:8"`
}

// NewAuthenticator creates a new Authenticator.
func NewAuthenticator(realm string, cname PrincipalName) (Authenticator, error) {
	seq, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		return Authenticator{}, err
	}
	t := time.Now().UTC()
	return Authenticator{
		AVNO:      iana.PVNO,
		CRealm:    realm,
		CName:     cname,
		Cksum:     Checksum{},
		Cusec:     int((t.UnixNano() / int64(time.Microsecond)) - (t.Unix() * 1e6)),
		CTime:     t,
		SeqNumber: seq.Int64(),
	}, nil
}

// GenerateSeqNumberAndSubKey sets the Authenticator's sequence number and subkey.
func (a *Authenticator) GenerateSeqNumberAndSubKey(keyType int32, keySize int) error {
	seq, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		return err
	}
	a.SeqNumber = seq.Int64()
	//Generate subkey value
	sk := make([]byte, keySize, keySize)
	rand.Read(sk)
	a.SubKey = EncryptionKey{
		KeyType:  keyType,
		KeyValue: sk,
	}
	return nil
}

// Unmarshal bytes into the Authenticator.
func (a *Authenticator) Unmarshal(b []byte) error {
	_, err := asn1.UnmarshalWithParams(b, a, fmt.Sprintf("application,explicit,tag:%v", asnAppTag.Authenticator))
	return err
}

// Marshal the Authenticator.
func (a *Authenticator) Marshal() ([]byte, error) {
	b, err := asn1.Marshal(*a)
	if err != nil {
		return nil, err
	}
	b = asn1tools.AddASNAppTag(b, asnAppTag.Authenticator)
	return b, nil
}
