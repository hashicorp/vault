package pgpkeys

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/keybase/go-crypto/openpgp"
)

// PubKeyFileFlag implements flag.Value and command.Example to receive exactly
// one PGP or keybase key via a flag.
type PubKeyFileFlag string

func (p *PubKeyFileFlag) String() string { return string(*p) }

func (p *PubKeyFileFlag) Set(val string) error {
	if p != nil && *p != "" {
		return errors.New("can only be specified once")
	}

	keys, err := ParsePGPKeys(strings.Split(val, ","))
	if err != nil {
		return err
	}

	if len(keys) > 1 {
		return errors.New("can only specify one pgp key")
	}

	*p = PubKeyFileFlag(keys[0])
	return nil
}

func (p *PubKeyFileFlag) Example() string { return "keybase:user" }

// PGPPubKeyFiles implements the flag.Value interface and allows parsing and
// reading a list of PGP public key files.
type PubKeyFilesFlag []string

func (p *PubKeyFilesFlag) String() string {
	return fmt.Sprint(*p)
}

func (p *PubKeyFilesFlag) Set(val string) error {
	if len(*p) > 0 {
		return errors.New("can only be specified once")
	}

	keys, err := ParsePGPKeys(strings.Split(val, ","))
	if err != nil {
		return err
	}

	*p = PubKeyFilesFlag(keys)
	return nil
}

func (p *PubKeyFilesFlag) Example() string { return "keybase:user1, keybase:user2, ..." }

// ParsePGPKeys takes a list of PGP keys and parses them either using keybase
// or reading them from disk and returns the "expanded" list of pgp keys in
// the same order.
func ParsePGPKeys(keyfiles []string) ([]string, error) {
	keys := make([]string, len(keyfiles))

	keybaseMap, err := FetchKeybasePubkeys(keyfiles)
	if err != nil {
		return nil, err
	}

	for i, keyfile := range keyfiles {
		keyfile = strings.TrimSpace(keyfile)

		if strings.HasPrefix(keyfile, kbPrefix) {
			key, ok := keybaseMap[keyfile]
			if !ok || key == "" {
				return nil, fmt.Errorf("keybase user %q not found", strings.TrimPrefix(keyfile, kbPrefix))
			}
			keys[i] = key
			continue
		}

		pgpStr, err := ReadPGPFile(keyfile)
		if err != nil {
			return nil, err
		}
		keys[i] = pgpStr
	}

	return keys, nil
}

// ReadPGPFile reads the given PGP file from disk.
func ReadPGPFile(path string) (string, error) {
	if len(path) <= 0 {
		return "", errors.New("empty path")
	}
	if path[0] == '@' {
		path = path[1:]
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return "", err
	}

	// First parse as an armored keyring file, if that doesn't work, treat it as a straight binary/b64 string
	keyReader := bytes.NewReader(buf.Bytes())
	entityList, err := openpgp.ReadArmoredKeyRing(keyReader)
	if err == nil {
		if len(entityList) != 1 {
			return "", fmt.Errorf("more than one key found in file %q", path)
		}
		if entityList[0] == nil {
			return "", fmt.Errorf("primary key was nil for file %q", path)
		}

		serializedEntity := bytes.NewBuffer(nil)
		err = entityList[0].Serialize(serializedEntity)
		if err != nil {
			return "", errwrap.Wrapf(fmt.Sprintf("error serializing entity for file %q: {{err}}", path), err)
		}

		return base64.StdEncoding.EncodeToString(serializedEntity.Bytes()), nil
	}

	_, err = base64.StdEncoding.DecodeString(buf.String())
	if err == nil {
		return buf.String(), nil
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil

}

// SharePathsFlag implements the flag.Value interface and allows
// writing the encrypted shares to the specified paths.
type SharePathsFlag []string

func (p *SharePathsFlag) String() string {
	return fmt.Sprint(*p)
}

func (p *SharePathsFlag) Set(val string) error {
	if len(*p) > 0 {
		return errors.New("can only be specified once")
	}

	paths := strings.Split(val, ",")
	err := ParseSharePaths(paths); if err != nil {
		return err
	}

	*p = SharePathsFlag(paths)
	return nil
}

func (p *SharePathsFlag) Example() string { return "user1.gpg, user2.gpg, ..." }

// ParseSharePaths takes a list of paths to output encrypted shares and checks
// if paths don't exist and they are writeable before attempting rekeying
func ParseSharePaths(SharePaths []string) error {
	for _, path := range SharePaths {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("file %s already exists", path))
		}
		file, err := os.Create(path)
		file.Close()
		if err != nil {
			return err
		}
		os.Remove(path)
	}
	return nil
}

// WriteSharePaths writes the given encrypted shares to specified paths in filesystem.
func WriteSharePaths(StoredKeys map[string][]string, SharePaths []string) error {
	if len(StoredKeys) != len(SharePaths) {
		return errors.New("number of supplied pgp-shares paths does not match number of entries in backup")
	}
	i := 0
	for _, v := range StoredKeys {
		keyBytes, err := hex.DecodeString(strings.Join(v, "")); if err != nil {
			return err
		}
		err = ioutil.WriteFile(SharePaths[i], keyBytes, 0600); if err != nil {
			return err
		}
		i++
	}
	return nil
}
