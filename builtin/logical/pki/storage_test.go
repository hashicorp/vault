package pki

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func Test_ConfigsRoundTrip(t *testing.T) {
	_, s := createBackendWithStorage(t)

	// Verify we handle nothing stored properly
	keyConfigEmpty, err := getKeysConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &keyConfig{}, keyConfigEmpty)

	issuerConfigEmpty, err := getIssuersConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, &issuerConfig{}, issuerConfigEmpty)

	// Now attempt to store and reload properly
	origKeyConfig := &keyConfig{
		DefaultKeyId: genKeyId(),
	}
	origIssuerConfig := &issuerConfig{
		DefaultIssuerId: genIssuerId(),
	}

	err = setKeysConfig(ctx, s, origKeyConfig)
	require.NoError(t, err)
	err = setIssuersConfig(ctx, s, origIssuerConfig)
	require.NoError(t, err)

	keyConfig, err := getKeysConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, origKeyConfig, keyConfig)

	issuerConfig, err := getIssuersConfig(ctx, s)
	require.NoError(t, err)
	require.Equal(t, origIssuerConfig, issuerConfig)
}

func Test_IssuerRoundTrip(t *testing.T) {
	b, s := createBackendWithStorage(t)
	issuer1, key1 := genIssuerAndKey(t, b)
	issuer2, key2 := genIssuerAndKey(t, b)

	// We get an error when issuer id not found
	_, err := fetchIssuerById(ctx, s, issuer1.ID)
	require.Error(t, err)

	// We get an error when key id not found
	_, err = fetchKeyById(ctx, s, key1.ID)
	require.Error(t, err)

	// Now write out our issuers and keys
	err = writeKey(ctx, s, key1)
	require.NoError(t, err)
	err = writeIssuer(ctx, s, &issuer1)
	require.NoError(t, err)

	err = writeKey(ctx, s, key2)
	require.NoError(t, err)
	err = writeIssuer(ctx, s, &issuer2)
	require.NoError(t, err)

	fetchedKey1, err := fetchKeyById(ctx, s, key1.ID)
	require.NoError(t, err)

	fetchedIssuer1, err := fetchIssuerById(ctx, s, issuer1.ID)
	require.NoError(t, err)

	require.Equal(t, &key1, fetchedKey1)
	require.Equal(t, &issuer1, fetchedIssuer1)

	keys, err := listKeys(ctx, s)
	require.NoError(t, err)

	require.ElementsMatch(t, []keyId{key1.ID, key2.ID}, keys)

	issuers, err := listIssuers(ctx, s)
	require.NoError(t, err)

	require.ElementsMatch(t, []issuerId{issuer1.ID, issuer2.ID}, issuers)
}

func Test_KeysIssuerImport(t *testing.T) {
	b, s := createBackendWithStorage(t)
	issuer1, key1 := genIssuerAndKey(t, b)
	issuer2, key2 := genIssuerAndKey(t, b)

	// Key 1 before Issuer 1; Issuer 2 before Key 2.
	// Remove KeyIDs from non-written entities before beginning.
	key1.ID = ""
	issuer1.ID = ""
	issuer1.KeyID = ""

	key1_ref1, existing, err := importKey(ctx, s, key1.PrivateKey, "key1")
	require.NoError(t, err)
	require.False(t, existing)
	require.Equal(t, key1.PrivateKey, key1_ref1.PrivateKey)

	// Make sure if we attempt to re-import the same private key, no import/updates occur.
	// So the existing flag should be set to true and we do not update the existing Name field.
	key1_ref2, existing, err := importKey(ctx, s, key1.PrivateKey, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, key1.PrivateKey, key1_ref1.PrivateKey)
	require.Equal(t, key1_ref1.ID, key1_ref2.ID)
	require.Equal(t, key1_ref1.Name, key1_ref2.Name)

	issuer1_ref1, existing, err := importIssuer(ctx, s, issuer1.Certificate, "issuer1")
	require.NoError(t, err)
	require.False(t, existing)
	require.Equal(t, issuer1.Certificate, issuer1_ref1.Certificate)
	require.Equal(t, key1_ref1.ID, issuer1_ref1.KeyID)
	require.Equal(t, "issuer1", issuer1_ref1.Name)

	// Make sure if we attempt to re-import the same issuer, no import/updates occur.
	// So the existing flag should be set to true and we do not update the existing Name field.
	issuer1_ref2, existing, err := importIssuer(ctx, s, issuer1.Certificate, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, issuer1.Certificate, issuer1_ref1.Certificate)
	require.Equal(t, issuer1_ref1.ID, issuer1_ref2.ID)
	require.Equal(t, key1_ref1.ID, issuer1_ref2.KeyID)
	require.Equal(t, issuer1_ref1.Name, issuer1_ref2.Name)

	err = writeIssuer(ctx, s, &issuer2)
	require.NoError(t, err)

	err = writeKey(ctx, s, key2)
	require.NoError(t, err)

	// Same double import tests as above, but make sure if the previous was created through writeIssuer not importIssuer.
	issuer2_ref, existing, err := importIssuer(ctx, s, issuer2.Certificate, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, issuer2.Certificate, issuer2_ref.Certificate)
	require.Equal(t, issuer2.ID, issuer2_ref.ID)
	require.Equal(t, "", issuer2_ref.Name)
	require.Equal(t, issuer2.KeyID, issuer2_ref.KeyID)

	// Same double import tests as above, but make sure if the previous was created through writeKey not importKey.
	key2_ref, existing, err := importKey(ctx, s, key2.PrivateKey, "ignore-me")
	require.NoError(t, err)
	require.True(t, existing)
	require.Equal(t, key2.PrivateKey, key2_ref.PrivateKey)
	require.Equal(t, key2.ID, key2_ref.ID)
	require.Equal(t, "", key2_ref.Name)
}

func Test_CAChainBuilding(t *testing.T) {
	// Each step of the process we import a CA and (its key) and validate
	// the CA Chain at that stage. Each Cert import should be unique, but
	// duplicating keys is expected. We map PEM->PEM to ensure things line up
	// nicely. If the import errors, the chain isn't checked.
	//
	// For ease of reading, these test cases are defined in a separate file,
	// constants_for_test.go.
	for testCaseIndex, test := range chainBuildingTestCases {
		// We want to guarantee that the results are stable. This is costly
		// to do at each step, so instead do it at the end. But, use the step
		// iteration to build the validation data.
		var mapIssuersChain map[issuerId]string

		_, s := createBackendWithStorage(t)
		for testStepIndex, step := range test.Steps {
			logPrefix := fmt.Sprintf("[test case %v / test step %v] ", testCaseIndex, testStepIndex)

			var err error
			var existing bool
			var k *key
			var i *issuer

			// We've gotta import keys before certs sometimes for testing,
			// so let's do it deterministically on the index of the cert.
			if testStepIndex%2 == 0 {
				if step.Key != "" {
					k, existing, err = importKey(ctx, s, step.Key, "")
					require.NoError(t, err, logPrefix+"expected importing key to always work without error")
					require.Equal(t, existing, step.KeyIsExisting, logPrefix+"expecting equal results from import, key existing values")
				}

				_, _, err = importIssuer(ctx, s, step.Cert, "")
				if err != nil != step.CertImportErrors {
					t.Fatalf(logPrefix+"expected cert import to error: %v -- actual err: %v", step.CertImportErrors, err)
				}

				// Can validate the key/issuer link when the key is imported
				// first.
				if i != nil && k != nil {
					require.Equal(t, k.ID, i.KeyID, logPrefix+"expecting imported key, issuer to match k.ID == i.KeyID")
				}
			} else {
				// As above, in opposite order.
				i, existing, err = importIssuer(ctx, s, step.Cert, "")
				if (err != nil) != step.CertImportErrors {
					t.Fatalf(logPrefix+"expected cert import to error: %v -- actual err: %v", step.CertImportErrors, err)
				}

				if step.Key != "" {
					k, existing, err = importKey(ctx, s, step.Key, "")
					require.NoError(t, err, logPrefix+"expected importing key to always work without error")
					require.Equal(t, existing, step.KeyIsExisting, logPrefix+"expecting equal results from import, key existing values")
				}

				// Can validate the key/issuer link when the key is existing.
				if existing && i != nil && k != nil {
					require.Equal(t, k.ID, i.KeyID, logPrefix+"expecting imported key, issuer to match k.ID == i.KeyID")
				}
			}

			if step.CertImportErrors {
				// Skip chain validation.
				continue
			}

			// Now validate all certs' CAChain fields. Set ourselves up to
			// guarantee a stable ordering too.
			issuers, err := listIssuers(ctx, s)
			require.NoError(t, err, logPrefix+"unable to list issuers")
			mapIssuersChain = make(map[issuerId]string)

			for _, identifier := range issuers {
				issuer, err := fetchIssuerById(ctx, s, identifier)
				require.NoError(t, err)

				for cert, chain := range step.CAChain {
					if issuer.Certificate != cert {
						continue
					}

					if len(issuer.CAChain) != len(chain) {
						t.Fatalf(logPrefix+"validating certificate %v / issuer %v: different length of chains: got %v / expected %v", cert, issuer, len(issuer.CAChain), len(chain))
					}

					for index, chainCert := range issuer.CAChain {
						mapIssuersChain[identifier] += chainCert

						if strings.Count(chainCert, "BEGIN CERTIFICATE") != 1 {
							t.Fatalf(logPrefix+"validating certificate %v / issuer %v: concat'd cert in chain field at index %v: got %v", cert, issuer, index, chainCert)
						}

						if !strings.Contains(chain[index], chainCert) {
							t.Fatalf(logPrefix+"validating certificate %v / issuer %v: different entry in chain at index %v: got %v / expected %v", cert, issuer, index, chainCert, chain[index])
						}
					}
				}
			}
		}

		// Finally, rebuild the full chain a few more times and ensure
		// the order is the same.
		for count := 0; count < 100; count++ {
			logPrefixStable := fmt.Sprintf("[test case %v / stability iteration: %v] ", testCaseIndex, count)
			err := rebuildIssuersChains(ctx, s, nil)
			require.NoError(t, err, logPrefixStable+"error building chain")

			if len(mapIssuersChain) == 0 {
				t.Fatal(logPrefixStable + "expected non-empty mapIssuersChain")
			}

			for identifier, originalChain := range mapIssuersChain {
				if len(originalChain) == 0 {
					t.Fatalf(logPrefixStable+"expected non-empty chain for issuer: %v", identifier)
				}

				issuer, err := fetchIssuerById(ctx, s, identifier)
				require.NoError(t, err, logPrefixStable+"unable to fetch issuer")

				var newChain string
				for _, chainCert := range issuer.CAChain {
					newChain += chainCert
				}

				require.Equal(t, originalChain, newChain, logPrefixStable+"expected stable sort order")
			}
		}
	}
}

func genIssuerAndKey(t *testing.T, b *backend) (issuer, key) {
	certBundle := genCertBundle(t, b)

	keyId := genKeyId()

	pkiKey := key{
		ID:             keyId,
		PrivateKeyType: certBundle.PrivateKeyType,
		PrivateKey:     certBundle.PrivateKey,
	}

	issuerId := genIssuerId()

	pkiIssuer := issuer{
		ID:           issuerId,
		KeyID:        keyId,
		Certificate:  certBundle.Certificate,
		CAChain:      certBundle.CAChain,
		SerialNumber: certBundle.SerialNumber,
	}

	return pkiIssuer, pkiKey
}

func genCertBundle(t *testing.T, b *backend) *certutil.CertBundle {
	// Pretty gross just to generate a cert bundle, but
	fields := addCACommonFields(map[string]*framework.FieldSchema{})
	fields = addCAKeyGenerationFields(fields)
	fields = addCAIssueFields(fields)
	apiData := &framework.FieldData{
		Schema: fields,
		Raw: map[string]interface{}{
			"exported": "internal",
			"cn":       "example.com",
			"ttl":      3600,
		},
	}
	_, _, role, respErr := b.getGenerationParams(ctx, apiData, "/pki")
	require.Nil(t, respErr)

	input := &inputBundle{
		req: &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "issue/testrole",
			Storage:   b.storage,
		},
		apiData: apiData,
		role:    role,
	}
	parsedCertBundle, err := generateCert(ctx, b, input, nil, true, rand.Reader)

	require.NoError(t, err)
	certBundle, err := parsedCertBundle.ToCertBundle()
	require.NoError(t, err)
	return certBundle
}
