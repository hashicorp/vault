package azureauth

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/version"
)

// Using the same time parsing logic from https://github.com/coreos/go-oidc
// This code is licensed under the Apache 2.0 license
type jsonTime time.Time

func (j *jsonTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	var unix int64

	if t, err := n.Int64(); err == nil {
		unix = t
	} else {
		f, err := n.Float64()
		if err != nil {
			return err
		}
		unix = int64(f)
	}
	*j = jsonTime(time.Unix(unix, 0))
	return nil
}

// strListContains does a case-insensitive search of the string
// list for the value
func strListContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if strings.ToLower(item) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}

const ossVaultGUID = `15cd22ce-24af-43a4-aa83-4c1a36a4b177`
const entVaultGUID = `b2c13ec1-60e8-4733-9a76-88dbb2ce2471`

// userAgent determines the User Agent to send on HTTP requests. This is mostly copied
// from the useragent helper in vault and may get replaced with something more general
// for plugins
func userAgent() string {
	ua := useragent.String()

	// ent has many version variations, so if it's not "dev" or "" we'll assume
	// it's an enterprise variation
	guid := ossVaultGUID
	ver := version.GetVersion()
	if ver.VersionMetadata != "" && ver.VersionMetadata != "dev" {
		guid = entVaultGUID
	}

	vaultIDString := fmt.Sprintf("; %s)", guid)

	return strings.Replace(ua, ")", vaultIDString, 1)
}
