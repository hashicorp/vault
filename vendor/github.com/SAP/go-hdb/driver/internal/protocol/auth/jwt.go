package auth

import (
	"fmt"
)

// JWT implements JWT authentication.
type JWT struct {
	token     string
	logonname string
	_cookie   []byte
}

// NewJWT creates a new authJWT instance.
func NewJWT(token string) *JWT { return &JWT{token: token} }

func (a *JWT) String() string { return fmt.Sprintf("method type %s token %s", a.Typ(), a.token) }

// Cookie implements the AuthCookieGetter interface.
func (a *JWT) Cookie() (string, []byte) { return a.logonname, a._cookie }

// Typ implements the Method interface.
func (a *JWT) Typ() string { return MtJWT }

// Order implements the Method interface.
func (a *JWT) Order() byte { return MoJWT }

// PrepareInitReq implements the Method interface.
func (a *JWT) PrepareInitReq(prms *Prms) error {
	prms.addString(a.Typ())
	prms.addString(a.token)
	return nil
}

// InitRepDecode implements the Method interface.
func (a *JWT) InitRepDecode(d *Decoder) error {
	a.logonname = d.String()
	return nil
}

// PrepareFinalReq implements the Method interface.
func (a *JWT) PrepareFinalReq(prms *Prms) error {
	prms.AddCESU8String(a.logonname)
	prms.addString(a.Typ())
	prms.addEmpty() // empty parameter
	return nil
}

// FinalRepDecode implements the Method interface.
func (a *JWT) FinalRepDecode(d *Decoder) error {
	if err := d.NumPrm(2); err != nil {
		return err
	}
	mt := d.String()
	if err := checkAuthMethodType(mt, a.Typ()); err != nil {
		return err
	}
	a._cookie = d.bytes()
	return nil
}
