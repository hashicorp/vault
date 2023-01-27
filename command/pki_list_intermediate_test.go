package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestPKIListIntermediate(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	// Relationship Map to Create
	//          pki-root			| pki-newroot  | pki-empty
	// RootX1    RootX2    RootX4     RootX3
	//   |								 |
	// ----------------------------------------------
	//   v								 v
	// IntX1					  	   IntX2       pki-int
	//   |								 |
	//   v								 v
	// IntX3 (-----------------------) IntX3(also)
	//
	// Here X1,X2 have the same name (same mount)
	// RootX4 uses the same key as RootX1 (but a different common_name/subject)
	// RootX3 has the same name, and is on a different mount
	// RootX1 has issued IntX1; RootX3 has issued IntX2
	createComplicatedIssuerSetUp(t, client)

	runPkiListIntermediateTests(t, client)
}

func runPkiListIntermediateTests(t *testing.T, client *api.Client) {
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
			"rootX1-match-everything-no-constraints",
			[]string{
				"pki", "list-intermediates", "-format=json", "-use_names=true",
				"-subject_match=false", "-key_id_match=false", "-direct_sign=false", "-indirect_sign=false", "-path_match=false",
				"pki-root/issuer/rootX1",
			},
			map[string]bool{
				"pki-root/issuer/rootX1":    true,
				"pki-root/issuer/rootX2":    true,
				"pki-newroot/issuer/rootX3": true,
				"pki-root/issuer/rootX4":    true,
				"pki-int/issuer/intX1":      true,
				"pki-int/issuer/intX2":      true,
				"pki-int/issuer/intX3":      true,
				"pki-int/issuer/intX3also":  true,
				"pki-int/issuer/rootX1":     true,
				"pki-int/issuer/rootX3":     true,
			},
			true,
			false,
			"",
			"",
			"",
		},
		{
			"rootX1-default-children",
			[]string{"pki", "list-intermediates", "-format=json", "-use_names=true", "pki-root/issuer/rootX1"},
			map[string]bool{
				"pki-root/issuer/rootX1":    true,
				"pki-root/issuer/rootX2":    false,
				"pki-newroot/issuer/rootX3": false,
				"pki-root/issuer/rootX4":    false,
				"pki-int/issuer/intX1":      true,
				"pki-int/issuer/intX2":      false,
				"pki-int/issuer/intX3":      false,
				"pki-int/issuer/intX3also":  false,
				"pki-int/issuer/rootX1":     true,
				"pki-int/issuer/rootX3":     false,
			},
			true,
			false,
			"",
			"",
			"",
		},
		{
			"rootX1-subject-match-only",
			[]string{
				"pki", "list-intermediates", "-format=json", "-use_names=true",
				"-key_id_match=false", "-direct_sign=false", "-indirect_sign=false",
				"pki-root/issuer/rootX1",
			},
			map[string]bool{
				"pki-root/issuer/rootX1":    true,
				"pki-root/issuer/rootX2":    true,
				"pki-newroot/issuer/rootX3": true,
				"pki-root/issuer/rootX4":    false,
				"pki-int/issuer/intX1":      true,
				"pki-int/issuer/intX2":      true,
				"pki-int/issuer/intX3":      false,
				"pki-int/issuer/intX3also":  false,
				"pki-int/issuer/rootX1":     true,
				"pki-int/issuer/rootX3":     true,
			},
			true,
			false,
			"",
			"",
			"",
		},
		{
			"rootX1-in-path",
			[]string{
				"pki", "list-intermediates", "-format=json", "-use_names=true",
				"-subject_match=false", "-key_id_match=false", "-direct_sign=false", "-indirect_sign=false", "-path_match=true",
				"pki-root/issuer/rootX1",
			},
			map[string]bool{
				"pki-root/issuer/rootX1":    true,
				"pki-root/issuer/rootX2":    false,
				"pki-newroot/issuer/rootX3": false,
				"pki-root/issuer/rootX4":    false,
				"pki-int/issuer/intX1":      true,
				"pki-int/issuer/intX2":      false,
				"pki-int/issuer/intX3":      true,
				"pki-int/issuer/intX3also":  false,
				"pki-int/issuer/rootX1":     true,
				"pki-int/issuer/rootX3":     false,
			},
			true,
			false,
			"",
			"",
			"",
		},
		{
			"rootX1-only-int-mount",
			[]string{
				"pki", "list-intermediates", "-format=json", "-use_names=true",
				"-subject_match=false", "-key_id_match=false", "-direct_sign=false", "-indirect_sign=false", "-path_match=true",
				"pki-root/issuer/rootX1", "pki-int/",
			},
			map[string]bool{
				"pki-int/issuer/intX1":     true,
				"pki-int/issuer/intX2":     false,
				"pki-int/issuer/intX3":     true,
				"pki-int/issuer/intX3also": false,
				"pki-int/issuer/rootX1":    true,
				"pki-int/issuer/rootX3":    false,
			},
			true,
			false,
			"",
			"",
			"",
		},
		{
			"rootX1-subject-match-root-mounts-only",
			[]string{
				"pki", "list-intermediates", "-format=json", "-use_names=true",
				"-key_id_match=false", "-direct_sign=false", "-indirect_sign=false",
				"pki-root/issuer/rootX1", "pki-root/", "pki-newroot", "pki-empty",
			},
			map[string]bool{
				"pki-root/issuer/rootX1":    true,
				"pki-root/issuer/rootX2":    true,
				"pki-newroot/issuer/rootX3": true,
				"pki-root/issuer/rootX4":    false,
			},
			true,
			false,
			"",
			"",
			"",
		},
		{
			"rootX1-subject-match-these-certs-only",
			[]string{
				"pki", "list-intermediates", "-format=json", "-use_names=true",
				"-key_id_match=false", "-direct_sign=false", "-indirect_sign=false",
				"pki-root/issuer/rootX1", "pki-root/issuer/rootX2", "pki-newroot/issuer/rootX3", "pki-root/issuer/rootX4",
			},
			map[string]bool{
				"pki-root/issuer/rootX2":    true,
				"pki-newroot/issuer/rootX3": true,
				"pki-root/issuer/rootX4":    false,
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
