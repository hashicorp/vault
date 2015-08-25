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
	for _, keyfile := range strings.Split(value, ",") {
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

		*p = append(*p, base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
	return nil
}
