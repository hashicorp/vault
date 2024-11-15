package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/driver/internal/protocol/auth"
	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

// AuthHnd holds the client authentication methods dependent on the driver.Connector attributes and handles the authentication hdb protocol.
type AuthHnd struct {
	logonname string
	methods   auth.Methods
	selected  auth.Method // selected method
}

// NewAuthHnd creates a new AuthHnd instance.
func NewAuthHnd(logonname string) *AuthHnd {
	return &AuthHnd{logonname: logonname, methods: auth.Methods{}}
}

func (a *AuthHnd) String() string { return "logonname " + a.logonname }

// AddSessionCookie adds session cookie authentication method.
func (a *AuthHnd) AddSessionCookie(cookie []byte, logonname, clientID string) {
	a.methods[auth.MtSessionCookie] = auth.NewSessionCookie(cookie, logonname, clientID)
}

// AddBasic adds basic authentication methods.
func (a *AuthHnd) AddBasic(username, password string) {
	a.methods[auth.MtSCRAMPBKDF2SHA256] = auth.NewSCRAMPBKDF2SHA256(username, password)
	a.methods[auth.MtSCRAMSHA256] = auth.NewSCRAMSHA256(username, password)
}

// AddJWT adds JWT authentication method.
func (a *AuthHnd) AddJWT(token string) { a.methods[auth.MtJWT] = auth.NewJWT(token) }

// AddX509 adds X509 authentication method.
func (a *AuthHnd) AddX509(certKey *auth.CertKey) { a.methods[auth.MtX509] = auth.NewX509(certKey) }

// Selected returns the selected authentication method.
func (a *AuthHnd) Selected() auth.Method { return a.selected }

func (a *AuthHnd) setMethod(mt string) error {
	var ok bool
	if a.selected, ok = a.methods[mt]; !ok {
		return fmt.Errorf("invalid method type: %s", mt)
	}
	return nil
}

// InitRequest returns the init request part.
func (a *AuthHnd) InitRequest() (*AuthInitRequest, error) {
	prms := &auth.Prms{}
	prms.AddCESU8String(a.logonname)
	for _, m := range a.methods.Order() {
		if err := m.PrepareInitReq(prms); err != nil {
			return nil, err
		}
	}
	return &AuthInitRequest{prms: prms}, nil
}

// InitReply returns the init reply part.
func (a *AuthHnd) InitReply() (*AuthInitReply, error) { return &AuthInitReply{authHnd: a}, nil }

// FinalRequest returns the final request part.
func (a *AuthHnd) FinalRequest() (*AuthFinalRequest, error) {
	prms := &auth.Prms{}
	if err := a.selected.PrepareFinalReq(prms); err != nil {
		return nil, err
	}
	return &AuthFinalRequest{prms}, nil
}

// FinalReply returns the final reply part.
func (a *AuthHnd) FinalReply() (*AuthFinalReply, error) {
	return &AuthFinalReply{method: a.selected}, nil
}

// AuthInitRequest represents an authentication initial request.
type AuthInitRequest struct {
	prms *auth.Prms
}

func (r *AuthInitRequest) String() string                     { return r.prms.String() }
func (r *AuthInitRequest) size() int                          { return r.prms.Size() }
func (r *AuthInitRequest) decode(dec *encoding.Decoder) error { return r.prms.Decode(dec) }
func (r *AuthInitRequest) encode(enc *encoding.Encoder) error { return r.prms.Encode(enc) }

// AuthInitReply represents an authentication initial reply.
type AuthInitReply struct {
	authHnd *AuthHnd
}

func (r *AuthInitReply) String() string { return r.authHnd.String() }
func (r *AuthInitReply) decode(dec *encoding.Decoder) error {
	if r.authHnd == nil {
		return nil
	}

	d := auth.NewDecoder(dec)

	if err := d.NumPrm(2); err != nil {
		return err
	}
	mt := d.String()

	if err := r.authHnd.setMethod(mt); err != nil {
		return err
	}
	if err := r.authHnd.selected.InitRepDecode(d); err != nil {
		return err
	}
	return dec.Error()
}

// AuthFinalRequest represents an authentication final request.
type AuthFinalRequest struct {
	prms *auth.Prms
}

func (r *AuthFinalRequest) String() string { return r.prms.String() }
func (r *AuthFinalRequest) size() int      { return r.prms.Size() }
func (r *AuthFinalRequest) decode(dec *encoding.Decoder) error {
	return nil
	// panic("not implemented yet")
}
func (r *AuthFinalRequest) encode(enc *encoding.Encoder) error { return r.prms.Encode(enc) }

// AuthFinalReply represents an authentication final reply.
type AuthFinalReply struct {
	method auth.Method
}

func (r *AuthFinalReply) String() string { return r.method.String() }
func (r *AuthFinalReply) decode(dec *encoding.Decoder) error {
	if r.method == nil {
		return nil
	}

	if err := r.method.FinalRepDecode(auth.NewDecoder(dec)); err != nil {
		return err
	}
	return dec.Error()
}
