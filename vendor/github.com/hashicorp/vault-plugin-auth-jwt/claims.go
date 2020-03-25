package jwtauth

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/mitchellh/pointerstructure"
	"github.com/ryanuber/go-glob"
)

// getClaim returns a claim value from allClaims given a provided claim string.
// If this string is a valid JSONPointer, it will be interpreted as such to locate
// the claim. Otherwise, the claim string will be used directly.
func getClaim(logger log.Logger, allClaims map[string]interface{}, claim string) interface{} {
	var val interface{}
	var err error

	if !strings.HasPrefix(claim, "/") {
		val = allClaims[claim]
	} else {
		val, err = pointerstructure.Get(allClaims, claim)
		if err != nil {
			logger.Warn(fmt.Sprintf("unable to locate %s in claims: %s", claim, err.Error()))
			return nil
		}
	}

	// The claims unmarshalled by go-oidc don't use UseNumber, so there will
	// be mismatches if they're coming in as float64 since Vault's config will
	// be represented as json.Number. If the operator can coerce claims data to
	// be in string form, there is no problem. Alternatively, we could try to
	// intelligently convert float64 to json.Number, e.g.:
	//
	// switch v := val.(type) {
	// case float64:
	// 	val = json.Number(strconv.Itoa(int(v)))
	// }
	//
	// Or we fork and/or PR go-oidc.

	return val
}

// extractMetadata builds a metadata map from a set of claims and claims mappings.
// The referenced claims must be strings and the claims mappings must be of the structure:
//
//   {
//       "/some/claim/pointer": "metadata_key1",
//       "another_claim": "metadata_key2",
//        ...
//   }
func extractMetadata(logger log.Logger, allClaims map[string]interface{}, claimMappings map[string]string) (map[string]string, error) {
	metadata := make(map[string]string)
	for source, target := range claimMappings {
		if value := getClaim(logger, allClaims, source); value != nil {
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("error converting claim '%s' to string", source)
			}

			metadata[target] = strValue
		}
	}
	return metadata, nil
}

// validateAudience checks whether any of the audiences in audClaim match those
// in boundAudiences. If strict is true and there are no bound audiences, then the
// presence of any audience in the received claim is considered an error.
func validateAudience(boundAudiences, audClaim []string, strict bool) error {
	if strict && len(boundAudiences) == 0 && len(audClaim) > 0 {
		return errors.New("audience claim found in JWT but no audiences bound to the role")
	}

	if len(boundAudiences) > 0 {
		for _, v := range boundAudiences {
			if strutil.StrListContains(audClaim, v) {
				return nil
			}
		}
		return errors.New("aud claim does not match any bound audience")
	}

	return nil
}

// validateBoundClaims checks that all of the claim:value requirements in boundClaims are
// met in allClaims.
func validateBoundClaims(logger log.Logger, boundClaimsType string, boundClaims, allClaims map[string]interface{}) error {
	useGlobs := boundClaimsType == boundClaimsTypeGlob

	for claim, expValue := range boundClaims {
		actValue := getClaim(logger, allClaims, claim)
		if actValue == nil {
			return fmt.Errorf("claim %q is missing", claim)
		}

		actVals, ok := normalizeList(actValue)
		if !ok {
			return fmt.Errorf("received claim is not a string or list: %v", actValue)
		}

		expVals, ok := normalizeList(expValue)
		if !ok {
			return fmt.Errorf("bound claim is not a string or list: %v", expValue)
		}

		found, err := matchFound(expVals, actVals, useGlobs)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("claim %q does not match any associated bound claim values", claim)
		}
	}
	return nil
}

func matchFound(expVals, actVals []interface{}, useGlobs bool) (bool, error) {
	for _, expVal := range expVals {
		for _, actVal := range actVals {
			if useGlobs {
				// Only string globbing is supported.
				expValStr, ok := expVal.(string)
				if !ok {
					return false, fmt.Errorf("received claim is not a glob string: %expVal", expVal)
				}
				actValStr, ok := actVal.(string)
				if !ok {
					continue
				}
				if !glob.Glob(expValStr, actValStr) {
					continue
				}
			} else {
				if actVal != expVal {
					continue
				}
			}
			return true, nil
		}
	}
	return false, nil
}

// normalizeList takes a string, bool or list and returns a list. This is useful when
// providers are expected to return a list (typically of strings) but reduce it
// to a string type when the list count is 1.
func normalizeList(raw interface{}) ([]interface{}, bool) {
	var normalized []interface{}

	switch v := raw.(type) {
	case []interface{}:
		normalized = v
	case string, bool:
		normalized = []interface{}{v}
	default:
		return nil, false
	}

	return normalized, true
}
