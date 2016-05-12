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

func (b *backend) createTOTPKey(method *mfaMethodEntry, identifierEntry *mfaIdentifierEntry, resp *logical.Response) error {
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

	identifierEntry.TOTPURL = key.String()
	identifierEntry.TOTPSecret = key.Secret()

	if resp.Data == nil {
		resp.Data = map[string]interface{}{}
	}
	resp.Data["totp_secret"] = key.Secret()

	keyImage, err := key.Image(1024, 1024)
	if err != nil {
		return err
	}

	pngBuf := bytes.NewBuffer(nil)
	err = png.Encode(pngBuf, keyImage)
	if err != nil {
		return err
	}
	resp.Data["totp_qrcode_png"] = base64.StdEncoding.EncodeToString(pngBuf.Bytes())

	// Force the response into a cubbyhole
	resp.WrapInfo.TTL = 5 * time.Minute

	return nil
}

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

	alg, err := method.totpAlgorithm()
	if err != nil {
		return false, nil, err
	}

	opts := totp.ValidateOpts{
		Algorithm: alg,
	}

	valid, err := totp.ValidateCustom(mfaInfo.Parameters["token"], identifierEntry.TOTPSecret, time.Now().UTC(), opts)
	if err != nil {
		return false, err, nil
	}

	return valid, nil, nil
}
