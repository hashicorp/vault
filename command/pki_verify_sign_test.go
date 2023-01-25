package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
	"testing"
)

func TestPKIVerifySign(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	// Relationship Map to Create
	//          pki-root			| pki-newroot
	// RootX1    RootX2    RootX4     RootX3
	//   |								 |
	// ----------------------------------------------
	//   v								 v
	// IntX1					  	   IntX2       pki-int
	//   |								 |
	//   v								 v
	// IntX3 (-----------------------) IntX3
	//
	// Here X1,X2 have the same name (same mount)
	// RootX4 uses the same key as RootX1 (but a different common_name/subject)
	// RootX3 has the same name, and is on a different mount
	// RootX1 has issued IntX1; RootX3 has issued IntX2
	createComplicatedIssuerSetUp(t, client)

	cases := []struct {
		name               string
		args               []string
		expectedMatches    map[string]bool
		jsonOut            bool
		shouldError        bool
		expectErrorCont    string
		expectErrorNotCont string
		nonJsonOutputCont  string
	}{
		{
			"rootX1-matches-rootX1",
			[]string{"pki", "verify-sign", "-format=json", "pki-root/issuer/rootX1", "pki-root/issuer/rootX1"},
			map[string]bool{
				"key_id_match":    true,
				"path_match":      true,
				"signature_match": true,
				"subject_match":   true,
				"trust_match":     true,
			},
			true,
			false,
			"",
			"",
			"",
		},
		{
			"rootX1-on-rootX2-onlySameName",
			[]string{"pki", "verify-sign", "-format=json", "pki-root/issuer/rootX1", "pki-root/issuer/rootX2"},
			map[string]bool{
				"key_id_match":    false,
				"path_match":      false,
				"signature_match": false,
				"subject_match":   true,
				"trust_match":     false,
			},
			true,
			false,
			"",
			"",
			"",
		},
	}
	for _, testCase := range cases {
		var errString string
		var results map[string]interface{}
		var stdOut string

		if testCase.jsonOut {
			results, errString = execPKIVerifyJson(t, client, false, testCase.shouldError, testCase.args)
		} else {
			stdOut, errString = execPKIVerifyNonJson(t, client, testCase.shouldError, testCase.args)
		}

		// Verify Error Behavior
		if testCase.shouldError {
			if errString == "" {
				t.Fatalf("Expected error in Testcase %s : no error produced, got results %s", testCase.name, results)
			}
			if testCase.expectErrorCont != "" && !strings.Contains(errString, testCase.expectErrorCont) {
				t.Fatalf("Expected error in Testcase %s to contain %s, but got error %s", testCase.name, testCase.expectErrorCont, errString)
			}
			if testCase.expectErrorNotCont != "" && strings.Contains(errString, testCase.expectErrorNotCont) {
				t.Fatalf("Expected error in Testcase %s to not contain %s, but got error %s", testCase.name, testCase.expectErrorNotCont, errString)
			}
		} else {
			if errString != "" {
				t.Fatalf("Error in Testcase %s : no error expected, but got error: %s", testCase.name, errString)
			}
		}

		// Verify Output
		if testCase.jsonOut {
			isMatch, errString := verifyExpectedJson(testCase.expectedMatches, results)
			if !isMatch {
				t.Fatalf("Expected Results for Testcase %s, do not match returned results %s", testCase.name, errString)
			}
		} else {
			if !strings.Contains(stdOut, testCase.nonJsonOutputCont) {
				t.Fatalf("Expected standard output for Testcase %s to contain %s, but got %s", testCase.name, testCase.nonJsonOutputCont, stdOut)
			}
		}

	}

}

func execPKIVerifyJson(t *testing.T, client *api.Client, expectErrorUnmarshalling bool, expectErrorOut bool, callArgs []string) (map[string]interface{}, string) {
	stdout, stderr := execPKIVerifyNonJson(t, client, expectErrorOut, callArgs)

	var results map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &results); err != nil && !expectErrorUnmarshalling {
		t.Fatalf("failed to decode json response : %v \n json: \n%v", err, stdout)
	}

	return results, stderr
}

func execPKIVerifyNonJson(t *testing.T, client *api.Client, expectErrorOut bool, callArgs []string) (string, string) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	runOpts := &RunOptions{
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	}

	code := RunCustom(callArgs, runOpts)
	if !expectErrorOut && code != 0 {
		t.Fatalf("running command `%v` unsuccessful (ret %v)\nerr: %v", strings.Join(callArgs, " "), code, stderr.String())
	}

	t.Log(stdout.String() + stderr.String())

	return stdout.String(), stderr.String()
}

func convertListOfInterfaceToString(list []interface{}, sep string) string {
	newList := make([]string, len(list))
	for i, interfa := range list {
		newList[i] = interfa.(string)
	}
	return strings.Join(newList, sep)
}

func createComplicatedIssuerSetUp(t *testing.T, client *api.Client) {
	// Relationship Map to Create
	//          pki-root			| pki-newroot
	// RootX1    RootX2    RootX4     RootX3
	//   |								 |
	// ----------------------------------------------
	//   v								 v
	// IntX1					  	   IntX2       pki-int
	//   |								 |
	//   v								 v
	// IntX3 (-----------------------) IntX3
	//
	// Here X1,X2 have the same name (same mount)
	// RootX4 uses the same key as RootX1 (but a different common_name/subject)
	// RootX3 has the same name, and is on a different mount
	// RootX1 has issued IntX1; RootX3 has issued IntX2

	if err := client.Sys().Mount("pki-root", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			MaxLeaseTTL: "36500d",
		},
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	if err := client.Sys().Mount("pki-newroot", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			MaxLeaseTTL: "36500d",
		},
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	if err := client.Sys().Mount("pki-int", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			MaxLeaseTTL: "36500d",
		},
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	resp, err := client.Logical().Write("pki-root/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X",
		"ttl":         "3650d",
		"issuer_name": "rootX1",
		"key_name":    "rootX1",
	})
	t.Logf("%s", resp.Data)
	if err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	resp, err = client.Logical().Write("pki-root/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X",
		"ttl":         "3650d",
		"issuer_name": "rootX2",
	})
	t.Logf("%s", resp.Data)
	if err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	if resp, err := client.Logical().Write("pki-newroot/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X",
		"ttl":         "3650d",
		"issuer_name": "rootX3",
	}); err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	if resp, err := client.Logical().Write("pki-root/root/generate/existing", map[string]interface{}{
		"common_name": "Root X4",
		"ttl":         "3650d",
		"issuer_name": "rootX4",
		"key_ref":     "rootX1",
	}); err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	// Intermediate X1
	int1CsrResp, err := client.Logical().Write("pki-int/intermediate/generate/internal", map[string]interface{}{
		"key_type":    "rsa",
		"common_name": "Int X1",
		"ttl":         "3650d",
	})
	if err != nil || int1CsrResp == nil {
		t.Fatalf("failed to generate CSR: %v", err)
	}
	int1CsrRaw, ok := int1CsrResp.Data["csr"]
	if !ok {
		t.Fatalf("no csr produced when generating intermediate, resp: %v", int1CsrResp)
	}
	int1Csr := int1CsrRaw.(string)
	int1CertResp, err := client.Logical().Write("pki-root/issuer/rootX1/sign-intermediate", map[string]interface{}{
		"csr": int1Csr,
	})
	if err != nil || int1CertResp == nil {
		t.Fatalf("failed to sign CSR: %v", err)
	}
	int1CertChainRaw, ok := int1CertResp.Data["ca_chain"]
	if !ok {
		t.Fatalf("no ca_chain produced when signing intermediate, resp: %v", int1CertResp)
	}
	int1CertChain := convertListOfInterfaceToString(int1CertChainRaw.([]interface{}), "\n")
	importInt1Resp, err := client.Logical().Write("pki-int/issuers/import/cert", map[string]interface{}{
		"pem_bundle": int1CertChain,
	})
	if err != nil || importInt1Resp == nil {
		t.Fatalf("failed to import certificate: %v", err)
	}
	importIssuerIdMap, ok := importInt1Resp.Data["mapping"]
	if !ok {
		t.Fatalf("no mapping data returned on issuer import: %v", importInt1Resp)
	}
	importIssuerId := ""
	for key, value := range importIssuerIdMap.(map[string]interface{}) {
		if value != nil && len(value.(string)) > 0 {
			importIssuerId = key
			break
		}
	}
	if resp, err := client.Logical().Write("pki-int/issuer/"+importIssuerId, map[string]interface{}{
		"issuer_name": "intX1",
	}); err != nil || resp == nil {
		t.Fatalf("error naming issuer %v", err)
	}

	// Intermediate X2
	int2CsrResp, err := client.Logical().Write("pki-int/intermediate/generate/internal", map[string]interface{}{
		"key_type":    "ed25519",
		"common_name": "Int X2",
		"ttl":         "3650d",
	})
	if err != nil || int2CsrResp == nil {
		t.Fatalf("failed to generate CSR: %v", err)
	}
	int2CsrRaw, ok := int2CsrResp.Data["csr"]
	if !ok {
		t.Fatalf("no csr produced when generating intermediate, resp: %v", int2CsrResp)
	}
	int2Csr := int2CsrRaw.(string)
	int2CertResp, err := client.Logical().Write("pki-newroot/issuer/rootX3/sign-intermediate", map[string]interface{}{
		"csr": int2Csr,
	})
	if err != nil || int2CertResp == nil {
		t.Fatalf("failed to sign CSR: %v", err)
	}
	int2CertChainRaw, ok := int2CertResp.Data["ca_chain"]
	if !ok {
		t.Fatalf("no ca_chain produced when signing intermediate, resp: %v", int2CertResp)
	}
	int2CertChain := convertListOfInterfaceToString(int2CertChainRaw.([]interface{}), "\n")
	importInt2Resp, err := client.Logical().Write("pki-int/issuers/import/cert", map[string]interface{}{
		"pem_bundle":  int2CertChain,
		"issuer_name": "intX2",
	})
	if err != nil || importInt2Resp == nil {
		t.Fatalf("failed to import certificate: %v", err)
	}
	importIssuer2IdMap, ok := importInt2Resp.Data["mapping"]
	if !ok {
		t.Fatalf("no mapping data returned on issuer import: %v", importInt2Resp)
	}
	importIssuer2Id := ""
	for key, value := range importIssuer2IdMap.(map[string]interface{}) {
		if value != nil && len(value.(string)) > 0 {
			importIssuer2Id = key
			break
		}
	}
	if resp, err := client.Logical().Write("pki-int/issuer/"+importIssuer2Id, map[string]interface{}{
		"issuer_name": "intX2",
	}); err != nil || resp == nil {
		t.Fatalf("error naming issuer %v", err)
	}

	// Intermediate X3
	int3CsrResp, err := client.Logical().Write("pki-int/intermediate/generate/internal", map[string]interface{}{
		"key_type":    "rsa",
		"common_name": "Int X3",
		"ttl":         "3650d",
	})
	if err != nil || int3CsrResp == nil {
		t.Fatalf("failed to generate CSR: %v", err)
	}
	int3CsrRaw, ok := int3CsrResp.Data["csr"]
	if !ok {
		t.Fatalf("no csr produced when generating intermediate, resp: %v", int3CsrResp)
	}
	int3Csr := int3CsrRaw.(string)
	// sign by intX1 and import
	int3CertResp1, err := client.Logical().Write("pki-int/issuer/intX1/sign-intermediate", map[string]interface{}{
		"csr": int3Csr,
	})
	if err != nil || int3CertResp1 == nil {
		t.Fatalf("failed to sign CSR: %v", err)
	}
	int3CertChainRaw1, ok := int3CertResp1.Data["ca_chain"]
	if !ok {
		t.Fatalf("no ca_chain produced when signing intermediate, resp: %v", int3CertResp1)
	}
	int3CertChain1 := convertListOfInterfaceToString(int3CertChainRaw1.([]interface{}), "\n")
	importInt3Resp1, err := client.Logical().Write("pki-int/issuers/import/cert", map[string]interface{}{
		"pem_bundle": int3CertChain1,
	})
	if err != nil || importInt3Resp1 == nil {
		t.Fatalf("failed to import certificate: %v", err)
	}
	importIssuer3IdMap1, ok := importInt3Resp1.Data["mapping"]
	if !ok {
		t.Fatalf("no mapping data returned on issuer import: %v", importInt2Resp)
	}
	importIssuer3Id1 := ""
	for key, value := range importIssuer3IdMap1.(map[string]interface{}) {
		if value != nil && len(value.(string)) > 0 {
			importIssuer3Id1 = key
			break
		}
	}
	if resp, err := client.Logical().Write("pki-int/issuer/"+importIssuer3Id1, map[string]interface{}{
		"issuer_name": "intX3",
	}); err != nil || resp == nil {
		t.Fatalf("error naming issuer %v", err)
	}
	// sign by intX2 and import
	int3CertResp2, err := client.Logical().Write("pki-int/issuer/intX2/sign-intermediate", map[string]interface{}{
		"csr": int3Csr,
	})
	if err != nil || int3CertResp2 == nil {
		t.Fatalf("failed to sign CSR: %v", err)
	}
	int3CertChainRaw2, ok := int3CertResp2.Data["ca_chain"]
	if !ok {
		t.Fatalf("no ca_chain produced when signing intermediate, resp: %v", int3CertResp2)
	}
	int3CertChain2 := convertListOfInterfaceToString(int3CertChainRaw2.([]interface{}), "\n")
	importInt3Resp2, err := client.Logical().Write("pki-int/issuers/import/cert", map[string]interface{}{
		"pem_bundle": int3CertChain2,
	})
	if err != nil || importInt3Resp2 == nil {
		t.Fatalf("failed to import certificate: %v", err)
	}
	importIssuer3IdMap2, ok := importInt3Resp2.Data["mapping"]
	if !ok {
		t.Fatalf("no mapping data returned on issuer import: %v", importInt2Resp)
	}
	importIssuer3Id2 := ""
	for key, value := range importIssuer3IdMap2.(map[string]interface{}) {
		if value != nil && len(value.(string)) > 0 {
			importIssuer3Id2 = key
			break
		}
	}
	if resp, err := client.Logical().Write("pki-int/issuer/"+importIssuer3Id2, map[string]interface{}{
		"issuer_name": "intX3also",
	}); err != nil || resp == nil {
		t.Fatalf("error naming issuer %v", err)
	}

}

func verifyExpectedJson(expectedResults map[string]bool, results map[string]interface{}) (isMatch bool, error string) {
	if len(expectedResults) != len(results) {
		return false, fmt.Sprintf("Different Number of Keys in Expected Results (%d), than results (%d)",
			len(expectedResults), len(results))
	}
	for key, value := range expectedResults {
		if results[key].(bool) != value {
			return false, fmt.Sprintf("Different value for key %s : expected %t got %s", key, value, results[key])
		}
	}
	return true, ""
}
