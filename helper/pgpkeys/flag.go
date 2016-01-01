package pgpkeys

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

// PGPPubKeyFiles implements the flag.Value interface and allows
// parsing and reading a list of pgp public key files
type PubKeyFilesFlag []string

func (p *PubKeyFilesFlag) String() string {
	return fmt.Sprint(*p)
}

func (p *PubKeyFilesFlag) Set(value string) error {
	if len(*p) > 0 {
		return errors.New("pgp-keys can only be specified once")
	}

	// First, resolve all the keybase entries...do it in one go so we only
	// round trip to the API once, then store locally
	keybaseMap := map[string]string{}
	keybaseUsers := []string{}
	for _, keyfile := range strings.Split(value, ",") {
		if strings.HasPrefix(keyfile, "keybase:") {
			keybaseUsers = append(keybaseUsers, strings.TrimPrefix(keyfile, "keybase:"))
		}
	}
	keybaseKeys, err := FetchKeybasePubkeys(keybaseUsers)
	if err != nil {
		return err
	}
	for i, key := range keybaseKeys {
		keybaseMap[keybaseUsers[i]] = key
	}

	// Now go through the actual flag, and substitute in resolved keybase
	// entries where appropriate
	for _, keyfile := range strings.Split(value, ",") {
		if strings.HasPrefix(keyfile, "keybase:") {
			username := strings.TrimPrefix(keyfile, "keybase:")
			key := keybaseMap[username]
			if key == "" {
				return fmt.Errorf("key for keybase user %s was not found in the map")
			}
			*p = append(*p, key)
			continue
		}
		if keyfile[0] == '@' {
			keyfile = keyfile[1:]
		}
		f, err := os.Open(keyfile)
		if err != nil {
			return err
		}
		defer f.Close()
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(f)
		if err != nil {
			return err
		}

		_, err = base64.StdEncoding.DecodeString(buf.String())
		if err == nil {
			*p = append(*p, buf.String())
		} else {
			*p = append(*p, base64.StdEncoding.EncodeToString(buf.Bytes()))
		}
	}
	return nil
}
