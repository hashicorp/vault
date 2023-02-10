package command

import (
	"bytes"
	"testing"

	"github.com/hashicorp/vault/api"
)

// TestPKIReIssueIntermediate tests that the pki reissue command line tool accurately copies information from the
// template certificate to the newly issued certificate, by issuing and reissuing several certificates and seeing how
// they related to each other.
func TestPKIReIssueIntermediate(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	// Relationship Map to Create
	//          pki-root			| pki-newroot | pki-empty
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
	createComplicatedIssuerSetUpWithReIssueIntermediate(t, client)

	runPkiVerifySignTests(t, client)

	runPkiListIntermediateTests(t, client)
}

func createComplicatedIssuerSetUpWithReIssueIntermediate(t *testing.T, client *api.Client) {
	// Relationship Map to Create
	//          pki-root			| pki-newroot  | pki-empty
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

	// Used to check handling empty list responses: Not Used for Any Issuers / Certificates
	if err := client.Sys().Mount("pki-empty", &api.MountInput{
		Type:   "pki",
		Config: api.MountConfigInput{},
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
	if err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	resp, err = client.Logical().Write("pki-root/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X",
		"ttl":         "3650d",
		"issuer_name": "rootX2",
	})
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

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	runOpts := &RunOptions{
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	}

	// Intermediate X1
	intX1CallArgs := []string{
		"pki", "issue", "-format=json", "-issuer_name=intX1",
		"pki-root/issuer/rootX1",
		"pki-int/",
		"key_type=rsa",
		"common_name=Int X1",
		"ou=thing",
		"ttl=3650d",
	}
	codeOut := RunCustom(intX1CallArgs, runOpts)
	if codeOut != 0 {
		t.Fatalf("error issuing intermediate X1, code: %d \n stdout: %v \n stderr: %v", codeOut, stdout, stderr)
	}

	// Intermediate X2 - using ReIssue
	intX2CallArgs := []string{
		"pki", "reissue", "-format=json", "-issuer_name=intX2",
		"pki-newroot/issuer/rootX3",
		"pki-int/issuer/intX1",
		"pki-int/",
		"key_type=ec",
		"common_name=Int X2",
	}
	codeOut = RunCustom(intX2CallArgs, runOpts)
	if codeOut != 0 {
		t.Fatalf("error issuing intermediate X2, code: %d \n stdout: %v \n stderr: %v", codeOut, stdout, stderr)
	}

	// Intermediate X3
	intX3OriginalCallArgs := []string{
		"pki", "issue", "-format=json", "-issuer_name=intX3",
		"pki-int/issuer/intX1",
		"pki-int/",
		"key_type=ec",
		"use_pss=true", // This is meaningful because rootX1 is an RSA key
		"signature_bits=512",
		"common_name=Int X3",
		"ttl=3650d",
	}
	codeOut = RunCustom(intX3OriginalCallArgs, runOpts)
	if codeOut != 0 {
		t.Fatalf("error issuing intermediate X3, code: %d \n stdout: %v \n stderr: %v", codeOut, stdout, stderr)
	}

	intX3AdaptedCallArgs := []string{
		"pki", "reissue", "-format=json", "-issuer_name=intX3also", "-type=existing",
		"pki-int/issuer/intX2", // This is a EC key
		"pki-int/issuer/intX3", // This template includes use_pss = true which can't be accomodated
		"pki-int/",
	}
	codeOut = RunCustom(intX3AdaptedCallArgs, runOpts)
	if codeOut != 0 {
		t.Fatalf("error issuing intermediate X3also, code: %d \n stdout: %v \n stderr: %v", codeOut, stdout, stderr)
	}
}
