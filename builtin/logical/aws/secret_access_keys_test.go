package aws

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestNormalizeDisplayName_NormRequired(t *testing.T) {
	invalidNames := map[string]string{
		"^#$test name\nshould be normalized)(*": "___test_name_should_be_normalized___",
		"^#$test name1 should be normalized)(*": "___test_name1_should_be_normalized___",
		"^#$test name  should be normalized)(*": "___test_name__should_be_normalized___",
		"^#$test name__should be normalized)(*": "___test_name__should_be_normalized___",
	}

	for k, v := range invalidNames {
		normalizedName := normalizeDisplayName(k)
		if normalizedName != v {
			t.Fatalf(
				"normalizeDisplayName does not normalize AWS name correctly: %s should resolve to %s",
				k,
				normalizedName)
		}
	}
}

func TestNormalizeDisplayName_NormNotRequired(t *testing.T) {
	validNames := []string{
		"test_name_should_normalize_to_itself@example.com",
		"test1_name_should_normalize_to_itself@example.com",
		"UPPERlower0123456789-_,.@example.com",
	}

	for _, n := range validNames {
		normalizedName := normalizeDisplayName(n)
		if normalizedName != n {
			t.Fatalf(
				"normalizeDisplayName erroneously normalizes valid names: expected %s but normalized to %s",
				n,
				normalizedName)
		}
	}
}

func TestGenUsername(t *testing.T) {

	testUsername, warning, err := genUsername("name1", "policy1", "iam_user", `{{ printf "vault-%s-%s-%s-%s" (.DisplayName) (.PolicyName) (unix_time) (random 20) | truncate 64 }}`)
	if err != nil {
		t.Fatalf(
			"expected no err; got %s",
			err,
		)
	}

	expectedUsernameRegex := `^vault-name1-policy1-[0-9]+-[a-zA-Z0-9]+`
	require.Regexp(t, expectedUsernameRegex, testUsername)
	// IAM usernames are capped at 64 characters
	if len(testUsername) > 64 {
		t.Fatalf(
			"expected IAM username to be of length 64, got %d",
			len(testUsername),
		)
	}

	testUsername, warning, err = genUsername(
		"this---is---a---very---long---name",
		"long------policy------name",
		"iam_user",
		`{{ printf "%s-%s-%s-%s" (.DisplayName) (.PolicyName) (unix_time) (random 20) }}`,
	)

	if warning == "" || !strings.Contains(warning, "calling token display name/IAM policy name were truncated to 64 characters") {
		t.Fatalf("expected a truncate warning; received empty string")
	}
	if len(testUsername) != 64 {
		t.Fatalf("expected a username cap at 64 chars; got length: %d", len(testUsername))
	}

	testUsername, warning, err = genUsername("name1", "policy1", "sts", defaultSTSTemplate)
	if strings.Contains(testUsername, "name1") || strings.Contains(testUsername, "policy1") {
		t.Fatalf(
			"expected sts username to not contain display name or policy name; got %s",
			testUsername,
		)
	}
	// STS usernames are capped at 64 characters
	if len(testUsername) > 32 {
		t.Fatalf(
			"expected sts username to be under 32 chars; got %s of length %d",
			testUsername,
			len(testUsername),
		)
	}
}

func TestReadConfig_DefaultTemplate(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b := Backend()
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	testTemplate := ""
	configData := map[string]interface{}{
		"connection_uri":    "test_uri",
		"username":          "guest",
		"password":          "guest",
		"username_template": testTemplate,
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/root",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err := b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	configResult, err := readConfig(context.Background(), config.StorageView)

	if err != nil {
		t.Fatalf("expected err to be nil; got %s", err)
	}

	// No template provided, config set to defaultUsernameTemplate
	if configResult.UsernameTemplate != defaultUserNameTemplate {
		t.Fatalf(
			"expected template %s; got %s",
			defaultUserNameTemplate,
			configResult.UsernameTemplate,
		)
	}
}

func TestReadConfig_CustomTemplate(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b := Backend()
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	testTemplate := "`foo-{{ .DisplayName }}`"
	configData := map[string]interface{}{
		"connection_uri":    "test_uri",
		"username":          "guest",
		"password":          "guest",
		"username_template": testTemplate,
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/root",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err := b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	configResult, err := readConfig(context.Background(), config.StorageView)

	if err != nil {
		t.Fatalf("expected err to be nil; got %s", err)
	}

	if configResult.UsernameTemplate != testTemplate {
		t.Fatalf(
			"expected template %s; got %s",
			testTemplate,
			configResult.UsernameTemplate,
		)
	}
}
