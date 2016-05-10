package pgpkeys

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"

	"golang.org/x/crypto/openpgp"
)

const (
	// MITKeyserver is the MIT PGP Keyserver url.
	MITKeyserver = "https://pgp.mit.edu"
)

// FetchKeyserverPubkeys fetches public keys from a given keyserver, a set of
// emails or ids, which are derived from correctly formatted input entries.
// The keys are returned as base64-encoded strings.
func FetchKeyserverPubkeys(keyserver string, input []string) (map[string]string, error) {
	client := cleanhttp.DefaultClient()

	if len(input) == 0 {
		return nil, nil
	}

	if keyserver == "" {
		return nil, errors.New("empty keyserver passed")
	}

	ret := make(map[string]string, len(input))
	serializedEntity := bytes.NewBuffer(nil)
	for i, lookup := range input {
		url := fmt.Sprintf("%s/pks/lookup?op=get&search=%s&options=mr", keyserver, lookup)
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status from server %s was %d", url, resp.StatusCode)
		}

		entityList, err := openpgp.ReadArmoredKeyRing(resp.Body)
		if err != nil {
			return nil, err
		}

		if len(entityList) != 1 {
			return nil, fmt.Errorf("primary key could not be parsed for user %s", input[i])
		}

		if entityList[0] == nil {
			return nil, fmt.Errorf("primary key was nil for user %s", input[i])
		}

		serializedEntity.Reset()
		err = entityList[0].Serialize(serializedEntity)
		if err != nil {
			return nil, fmt.Errorf("error serializing entity for user %s: %v", input[i], err)
		}

		// The API returns values in the same ordering requested, so this should properly match
		ret[input[i]] = base64.StdEncoding.EncodeToString(serializedEntity.Bytes())
	}

	return ret, nil
}
