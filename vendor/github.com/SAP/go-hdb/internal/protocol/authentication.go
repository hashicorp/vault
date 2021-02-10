// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

// Salted Challenge Response Authentication Mechanism (SCRAM)

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"

	"golang.org/x/crypto/pbkdf2"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
	"github.com/SAP/go-hdb/internal/unicode/cesu8"
)

const (
	mnSCRAMSHA256       = "SCRAMSHA256"       // password
	mnSCRAMPBKDF2SHA256 = "SCRAMPBKDF2SHA256" // pbkdf2
)

const (
	clientChallengeSize = 64
	serverChallengeSize = 48
	saltSize            = 16
	clientProofSize     = 32
)

const (
	int16Size  = 2
	uint32Size = 4
)

type authStepper interface {
	next() (partReadWriter, error)
}

type _authShortCESU8String struct{}
type _authShortBytes struct{}

var authShortCESU8String = _authShortCESU8String{}
var authShortBytes = _authShortBytes{}

func (_authShortCESU8String) decode(dec *encoding.Decoder) string {
	size := dec.Byte()
	return string(dec.CESU8Bytes(int(size)))
}

func (_authShortCESU8String) encode(enc *encoding.Encoder, s string) error {
	size := cesu8.StringSize(s)
	if size > math.MaxUint8 {
		return fmt.Errorf("invalid auth parameter length %d", size)
	}
	enc.Byte(byte(size))
	enc.CESU8String(s)
	return nil
}

func (_authShortBytes) decode(dec *encoding.Decoder) []byte {
	size := dec.Byte()
	b := make([]byte, size)
	dec.Bytes(b)
	return b
}

func (_authShortBytes) encode(enc *encoding.Encoder, b []byte) error {
	size := len(b)
	if size > math.MaxUint8 {
		return fmt.Errorf("invalid auth parameter length %d", size)
	}
	enc.Byte(byte(size))
	enc.Bytes(b)
	return nil
}

type authMethod struct {
	method          string
	clientChallenge []byte
}

func (m *authMethod) String() string {
	return fmt.Sprintf("method %s clientChallenge %v", m.method, m.clientChallenge)
}

func (m *authMethod) size() int {
	size := 2 // number of parameters
	size += len(m.method)
	size += len(m.clientChallenge)
	return size
}

func (m *authMethod) decode(dec *encoding.Decoder, ph *partHeader) error {
	m.method = string(authShortBytes.decode(dec))
	m.clientChallenge = authShortBytes.decode(dec)
	return nil
}

func (m *authMethod) encode(enc *encoding.Encoder) error {
	if err := authShortBytes.encode(enc, []byte(m.method)); err != nil {
		return err
	}
	if err := authShortBytes.encode(enc, m.clientChallenge); err != nil {
		return err
	}
	return nil
}

type authInitReq struct {
	username string
	methods  []*authMethod
}

func (r *authInitReq) String() string {
	return fmt.Sprintf("username %s methods %v", r.username, r.methods)
}

func (r *authInitReq) size() int {
	size := int16Size // no of parameters
	size++            // len byte username
	size += cesu8.StringSize(r.username)
	for _, m := range r.methods {
		size += m.size()
	}
	return size
}

func (r *authInitReq) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	r.username = authShortCESU8String.decode(dec)
	numMethod := (numPrm - 1) / 2
	r.methods = make([]*authMethod, numMethod)
	for i := 0; i < len(r.methods); i++ {
		authMethod := &authMethod{}
		r.methods[i] = authMethod
		if err := authMethod.decode(dec, ph); err != nil {
			return err
		}
	}
	return nil
}

func (r *authInitReq) encode(enc *encoding.Encoder) error {
	enc.Int16(int16(1 + len(r.methods)*2)) // username + methods á each two fields
	if err := authShortCESU8String.encode(enc, r.username); err != nil {
		return err
	}
	for _, m := range r.methods {
		m.encode(enc)
	}
	return nil
}

type authInitSCRAMSHA256Rep struct {
	salt, serverChallenge []byte
}

func (r *authInitSCRAMSHA256Rep) String() string {
	return fmt.Sprintf("salt %v serverChallenge %v", r.salt, r.serverChallenge)
}

func (r *authInitSCRAMSHA256Rep) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	if numPrm != 2 {
		return fmt.Errorf("invalid number of parameters %d - expected %d", numPrm, 2)
	}
	r.salt = authShortBytes.decode(dec)
	r.serverChallenge = authShortBytes.decode(dec)
	return nil
}

type authInitSCRAMPBKDF2SHA256Rep struct {
	salt, serverChallenge []byte
	rounds                uint32
}

func (r *authInitSCRAMPBKDF2SHA256Rep) String() string {
	return fmt.Sprintf("salt %v serverChallenge %v rounds %d", r.salt, r.serverChallenge, r.rounds)
}

func (r *authInitSCRAMPBKDF2SHA256Rep) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	if numPrm != 3 {
		return fmt.Errorf("invalid number of parameters %d - expected %d", numPrm, 3)
	}
	r.salt = authShortBytes.decode(dec)
	r.serverChallenge = authShortBytes.decode(dec)
	size := dec.Byte()
	if size != uint32Size {
		return fmt.Errorf("invalid auth uint32 size %d - expected %d", size, uint32Size)
	}
	r.rounds = dec.Uint32ByteOrder(binary.BigEndian) // big endian coded (e.g. rounds param)
	return nil
}

type authInitRep struct {
	method string
	prms   partDecoder
}

func (r *authInitRep) String() string {
	return fmt.Sprintf("method %s parameters %v", r.method, r.prms)
}
func (r *authInitRep) size() int                          { panic("not implemented") }
func (r *authInitRep) encode(enc *encoding.Encoder) error { panic("not implemented") }

func (r *authInitRep) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	if numPrm != 2 {
		return fmt.Errorf("invalid number of parameters %d - expected %d", numPrm, 2)
	}
	r.method = string(authShortBytes.decode(dec))

	dec.Byte() // sub parameter length

	switch r.method {
	case mnSCRAMSHA256:
		r.prms = &authInitSCRAMSHA256Rep{}
		return r.prms.decode(dec, ph)
	case mnSCRAMPBKDF2SHA256:
		r.prms = &authInitSCRAMPBKDF2SHA256Rep{}
		return r.prms.decode(dec, ph)
	default:
		return fmt.Errorf("invalid or not supported authentication method %s", r.method)
	}
}

type authClientProofReq struct {
	clientProof []byte
}

func (r *authClientProofReq) String() string { return fmt.Sprintf("clientProof %v", r.clientProof) }

func (r *authClientProofReq) size() int {
	size := int16Size // no of parameters
	size += len(r.clientProof) + 1
	return size
}

func (r *authClientProofReq) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	if numPrm != 1 {
		return fmt.Errorf("invalid number of parameters %d - expected %d", numPrm, 1)
	}
	r.clientProof = authShortBytes.decode(dec)
	return nil
}

func (r *authClientProofReq) encode(enc *encoding.Encoder) error {
	enc.Int16(1)
	if err := authShortBytes.encode(enc, r.clientProof); err != nil {
		return err
	}
	return nil
}

type authFinalReq struct {
	username, method string
	prms             partDecodeEncoder
}

func (r *authFinalReq) String() string {
	return fmt.Sprintf("username %s methods %s parameter %v", r.username, r.method, r.prms)
}

func (r *authFinalReq) size() int {
	size := int16Size // no of parameters
	size += cesu8.StringSize(r.username) + 1
	size += len(r.method) + 1
	size++ // len sub parameters
	size += r.prms.size()
	return size
}

func (r *authFinalReq) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	if numPrm != 3 {
		return fmt.Errorf("invalid number of parameters %d - expected %d", numPrm, 3)
	}
	r.username = authShortCESU8String.decode(dec)
	r.method = string(authShortBytes.decode(dec))
	dec.Byte() // sub parameters
	r.prms = &authClientProofReq{}
	return r.prms.decode(dec, ph)
}

func (r *authFinalReq) encode(enc *encoding.Encoder) error {
	enc.Int16(3)
	if err := authShortCESU8String.encode(enc, r.username); err != nil {
		return err
	}
	if err := authShortBytes.encode(enc, []byte(r.method)); err != nil {
		return err
	}
	enc.Byte(byte(r.prms.size()))
	return r.prms.encode(enc)
}

type authServerProofRep struct {
	serverProof []byte
}

func (r *authServerProofRep) String() string { return fmt.Sprintf("serverProof %v", r.serverProof) }

func (r *authServerProofRep) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	if numPrm != 1 {
		return fmt.Errorf("invalid number of server proof parameters %d - expected %d", numPrm, 1)
	}
	r.serverProof = authShortBytes.decode(dec)
	return nil
}

func (r *authServerProofRep) encode(enc *encoding.Encoder) error {
	enc.Int16(1)
	if err := authShortBytes.encode(enc, r.serverProof); err != nil {
		return err
	}
	return nil
}

type authFinalRep struct {
	method string
	prms   partDecoder
}

func (r *authFinalRep) String() string {
	return fmt.Sprintf("method %s parameter %v", r.method, r.prms)
}
func (r *authFinalRep) size() int                          { panic("not implemented") }
func (r *authFinalRep) encode(enc *encoding.Encoder) error { panic("not implemented") }

func (r *authFinalRep) decode(dec *encoding.Decoder, ph *partHeader) error {
	numPrm := int(dec.Int16())
	if numPrm != 2 {
		return fmt.Errorf("invalid number of parameters %d - expected %d", numPrm, 2)
	}
	r.method = string(authShortBytes.decode(dec))
	if size := dec.Byte(); size == 0 { // sub parameter length
		// mnSCRAMSHA256: server does not return server proof parameter
		return nil
	}
	r.prms = &authServerProofRep{}
	return r.prms.decode(dec, ph)
}

type auth struct {
	step               int
	username, password string
	methods            []*authMethod
	initRep            *authInitRep
}

func newAuth(username, password string) *auth {
	return &auth{
		username: username,
		password: password,
		methods: []*authMethod{
			{method: mnSCRAMPBKDF2SHA256, clientChallenge: clientChallenge()},
			{method: mnSCRAMSHA256, clientChallenge: clientChallenge()},
		},
		initRep: &authInitRep{},
	}
}

func (a *auth) clientChallenge(method string) []byte {
	for _, m := range a.methods {
		if m.method == method {
			return m.clientChallenge
		}
	}
	panic("should never happen")
}

func (a *auth) next() (partReadWriter, error) {
	defer func() { a.step++ }()

	switch a.step {
	case 0:
		for _, m := range a.methods {
			if len(m.clientChallenge) != clientChallengeSize {
				return nil, fmt.Errorf("invalid client challenge size %d - expected %d", len(m.clientChallenge), clientChallengeSize)
			}
		}
		return &authInitReq{username: a.username, methods: a.methods}, nil
	case 1:
		return a.initRep, nil
	case 2:
		var salt, serverChallenge, key []byte

		switch a.initRep.method {
		case mnSCRAMSHA256:
			prms := a.initRep.prms.(*authInitSCRAMSHA256Rep)
			salt, serverChallenge = prms.salt, prms.serverChallenge
			key = scramsha256Key([]byte(a.password), prms.salt)
		case mnSCRAMPBKDF2SHA256:
			prms := a.initRep.prms.(*authInitSCRAMPBKDF2SHA256Rep)
			salt, serverChallenge = prms.salt, prms.serverChallenge
			key = scrampbkdf2sha256Key([]byte(a.password), prms.salt, int(prms.rounds))
		default:
			panic("should never happen")
		}
		if len(salt) != saltSize {
			return nil, fmt.Errorf("invalid salt size %d - expected %d", len(salt), saltSize)
		}
		if len(serverChallenge) != serverChallengeSize {
			return nil, fmt.Errorf("invalid server challenge size %d - expected %d", len(serverChallenge), serverChallengeSize)
		}
		clientProof := clientProof(key, salt, serverChallenge, a.clientChallenge(a.initRep.method))
		if len(clientProof) != clientProofSize {
			return nil, fmt.Errorf("invalid client proof size %d - expected %d", len(clientProof), clientProofSize)
		}
		return &authFinalReq{username: a.username, method: a.initRep.method, prms: &authClientProofReq{clientProof: clientProof}}, nil
	case 3:
		return &authFinalRep{}, nil
	}
	panic("should never happen")
}

func clientChallenge() []byte {
	r := make([]byte, clientChallengeSize)
	if _, err := rand.Read(r); err != nil {
		plog.Fatalf("client challenge fatal error")
	}
	return r
}

func scramsha256Key(password, salt []byte) []byte {
	return _sha256(_hmac(password, salt))
}

func scrampbkdf2sha256Key(password, salt []byte, rounds int) []byte {
	return _sha256(pbkdf2.Key(password, salt, rounds, clientProofSize, sha256.New))
}

func clientProof(key, salt, serverChallenge, clientChallenge []byte) []byte {
	sig := _hmac(_sha256(key), salt, serverChallenge, clientChallenge)
	proof := xor(sig, key)
	return proof
}

func _sha256(p []byte) []byte {
	hash := sha256.New()
	hash.Write(p)
	s := hash.Sum(nil)
	return s
}

func _hmac(key []byte, prms ...[]byte) []byte {
	hash := hmac.New(sha256.New, key)
	for _, p := range prms {
		hash.Write(p)
	}
	s := hash.Sum(nil)
	return s
}

func xor(sig, key []byte) []byte {
	r := make([]byte, len(sig))

	for i, v := range sig {
		r[i] = v ^ key[i]
	}
	return r
}
