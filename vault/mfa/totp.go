package mfa

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/pquerna/otp/totp"
)

// createTOTPKey uses the given method and identifier to create a TOTP secret.
// It returns an error on error, but may also add warnings to the response
// object.
func (b *backend) createTOTPKey(method *mfaMethodEntry, identifierEntry *mfaIdentifierEntry, resp *logical.Response) error {
	if identifierEntry.TOTPAccountName == "" {
		newResp := logical.ErrorResponse("\"totp_account_name\" must be set on the identifier")
		*resp = *newResp
		return nil
	}

	alg, err := method.totpAlgorithm()
	if err != nil {
		return err
	}

	opts := totp.GenerateOpts{
		Algorithm:   alg,
		AccountName: identifierEntry.TOTPAccountName,
		Issuer:      fmt.Sprintf("Vault MFA: %s/%s", method.Name, identifierEntry.Identifier),
	}

	if opts.AccountName != "" {
		opts.Issuer = fmt.Sprintf("%s (%s)", opts.Issuer, opts.AccountName)
	}

	key, err := totp.Generate(opts)
	if err != nil {
		return err
	}
	if key == nil {
		return fmt.Errorf("generated TOTP key is nil")
	}

	// We store the original URL so we have the full set of parameters in the
	// key, but also the secret itself as a shortcut to avoid having to parse
	// the URL over and over
	identifierEntry.TOTPURL = key.String()
	identifierEntry.TOTPSecret = key.Secret()

	if resp.Data == nil {
		resp.Data = map[string]interface{}{}
	}

	resp.Data["totp_secret"] = key.Secret()
	resp.Data["totp_url"] = key.String()

	// Generate a PNG. This isn't super huge and can be quite helpful to the
	// end user.
	keyImage, err := key.Image(1024, 1024)
	if err != nil {
		return err
	}
	pngBuf := bytes.NewBuffer(nil)
	err = png.Encode(pngBuf, keyImage)
	if err != nil {
		return err
	}
	resp.Data["totp_qrcode_png_b64"] = base64.StdEncoding.EncodeToString(pngBuf.Bytes())

	// Force the response into a cubbyhole
	if resp.WrapInfo == nil {
		resp.WrapInfo = &logical.WrapInfo{}
	}
	resp.WrapInfo.TTL = 5 * time.Minute

	return nil
}

// validateTOTP ensures that the given token is valid for the given method and identifier
func (b *backend) validateTOTP(methodName string, mfaInfo *logical.MFAInfo) (bool, error, error) {
	if methodName == "" {
		return false, fmt.Errorf("no method name supplied"), nil
	}

	if mfaInfo == nil {
		return false, fmt.Errorf("no MFA information found"), nil
	}

	if mfaInfo.Parameters == nil {
		return false, fmt.Errorf("no MFA parameters supplied"), nil
	}

	if mfaInfo.Parameters["token"] == "" {
		return false, fmt.Errorf("no token supplied"), nil
	}

	if mfaInfo.Parameters["identifier"] == "" {
		return false, fmt.Errorf("no identifier supplied"), nil
	}

	// Look up the identifier
	method, identifierEntry, err := b.mfaBackendMethodIdentifiers(methodName, mfaInfo.Parameters["identifier"])
	if err != nil {
		return false, nil, err
	}
	if method == nil {
		return false, fmt.Errorf("method not found"), nil
	}
	if identifierEntry == nil {
		return false, fmt.Errorf("identifier not found"), nil
	}

	// Get the algorithm used
	alg, err := method.totpAlgorithm()
	if err != nil {
		return false, nil, err
	}

	opts := totp.ValidateOpts{
		Algorithm: alg,
	}

	// Validate!
	valid, err := totp.ValidateCustom(mfaInfo.Parameters["token"], identifierEntry.TOTPSecret, time.Now().UTC(), opts)
	if err != nil {
		return false, err, nil
	}

	return valid, nil, nil
}
