package certutil

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

// GetOctalFormatted returns the byte buffer formatted in octal with
// the specified separator between bytes
func GetOctalFormatted(buf []byte, sep string) string {
	var ret bytes.Buffer
	for _, cur := range buf {
		if ret.Len() > 0 {
			fmt.Fprintf(&ret, sep)
		}
		fmt.Fprintf(&ret, "%02x", cur)
	}
	return ret.String()
}

// ParsePKISecret takes a Secret returned from the PKI backend
// and returns a CertBundle for further processing (espcially
// by converting to a RawCertBundle)
func ParsePKISecret(secret *api.Secret) (*CertBundle, error) {
	var result CertBundle
	err := mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
