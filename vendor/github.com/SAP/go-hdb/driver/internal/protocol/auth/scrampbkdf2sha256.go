package auth

// Salted Challenge Response Authentication Mechanism (SCRAM)

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/SAP/go-hdb/driver/internal/protocol/cache"
	"golang.org/x/crypto/pbkdf2"
)

func scrampbkdf2sha256Key(password, salt []byte, rounds int) []byte {
	return _sha256(pbkdf2.Key(password, salt, rounds, clientProofSize, sha256.New))
}

// use cache as key calculation is expensive.
var scrampbkdf2KeyCache = cache.NewList(3, func(k *SCRAMPBKDF2SHA256) []byte {
	return scrampbkdf2sha256Key([]byte(k.password), k.salt, int(k.rounds))
})

// SCRAMPBKDF2SHA256 implements SCRAMPBKDF2SHA256 authentication.
type SCRAMPBKDF2SHA256 struct {
	username, password    string
	clientChallenge       []byte
	salt, serverChallenge []byte
	serverProof           []byte
	rounds                uint32
}

// NewSCRAMPBKDF2SHA256 creates a new authSCRAMPBKDF2SHA256 instance.
func NewSCRAMPBKDF2SHA256(username, password string) *SCRAMPBKDF2SHA256 {
	return &SCRAMPBKDF2SHA256{username: username, password: password, clientChallenge: clientChallenge()}
}

func (a *SCRAMPBKDF2SHA256) String() string {
	return fmt.Sprintf("method type %s clientChallenge %v", a.Typ(), a.clientChallenge)
}

// Compare implements cache.Compare interface.
func (a *SCRAMPBKDF2SHA256) Compare(a1 *SCRAMPBKDF2SHA256) bool {
	return a.password == a1.password && bytes.Equal(a.salt, a1.salt) && a.rounds == a1.rounds
}

// Typ implements the Method interface.
func (a *SCRAMPBKDF2SHA256) Typ() string { return MtSCRAMPBKDF2SHA256 }

// Order implements the Method interface.
func (a *SCRAMPBKDF2SHA256) Order() byte { return MoSCRAMPBKDF2SHA256 }

// PrepareInitReq implements the Method interface.
func (a *SCRAMPBKDF2SHA256) PrepareInitReq(prms *Prms) error {
	prms.addString(a.Typ())
	prms.addBytes(a.clientChallenge)
	return nil
}

// InitRepDecode implements the Method interface.
func (a *SCRAMPBKDF2SHA256) InitRepDecode(d *Decoder) error {
	d.subSize() // sub parameters
	if err := d.NumPrm(3); err != nil {
		return err
	}
	a.salt = d.bytes()
	a.serverChallenge = d.bytes()
	if err := checkSalt(a.salt); err != nil {
		return err
	}
	if err := checkServerChallenge(a.serverChallenge); err != nil {
		return err
	}
	var err error
	if a.rounds, err = d.bigUint32(); err != nil {
		return err
	}
	return nil
}

// PrepareFinalReq implements the Method interface.
func (a *SCRAMPBKDF2SHA256) PrepareFinalReq(prms *Prms) error {
	key := scrampbkdf2KeyCache.Get(a)
	clientProof, err := clientProof(key, a.salt, a.serverChallenge, a.clientChallenge)
	if err != nil {
		return err
	}

	prms.AddCESU8String(a.username)
	prms.addString(a.Typ())
	subPrms := prms.addPrms()
	subPrms.addBytes(clientProof)

	return nil
}

// FinalRepDecode implements the Method interface.
func (a *SCRAMPBKDF2SHA256) FinalRepDecode(d *Decoder) error {
	if err := d.NumPrm(2); err != nil {
		return err
	}
	mt := d.String()
	if err := checkAuthMethodType(mt, a.Typ()); err != nil {
		return err
	}
	d.subSize()
	if err := d.NumPrm(1); err != nil {
		return err
	}
	a.serverProof = d.bytes()
	return nil
}
