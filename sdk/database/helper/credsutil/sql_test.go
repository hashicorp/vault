package credsutil

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

func TestGenerateUsername(t *testing.T) {
	// The length of random alphanumeric string is hardcoded
	// https://github.com/hashicorp/vault/blob/73b1eed7c9066b121f6d4c77aa2d8ace21b48fef/sdk/database/helper/credsutil/sql.go#L58
	RandomAlphaNumericLen := 20
	UnixTimestampLen := len(fmt.Sprint(time.Now().Unix()))

	credsProducer := &SQLCredentialsProducer{
		DisplayNameLen: 30,
		RoleNameLen:    8,
		UsernameLen:    55,
		Separator:      "-",
	}

	type testCase struct {
		displayName        string
		roleName           string
		expectedUserPrefix string
	}

	tests := map[string]testCase{
		"short display name and role name are fully shown": {
			displayName: "short",
			roleName: "ro",
			expectedUserPrefix: "v-short-ro-",
		},
		"long display name is truncated": {
			displayName: "longlonglonglonglonglonglonglong",
			roleName: "ro",
			expectedUserPrefix: "v-longlonglonglonglonglonglonglo-ro-",
		},
		"long role name is truncated": {
			displayName: "short",
			roleName: "readwrite",
			expectedUserPrefix: "v-short-readwrit-",
		},
		"long display name and role name are truncated": {
			displayName: "longlonglonglonglonglonglonglong",
			roleName: "readwrite",
			expectedUserPrefix: "v-longlonglonglonglonglonglonglo-readwrit-",
		},
		"empty display name": {
			displayName: "",
			roleName: "ro",
			expectedUserPrefix: "v-ro-",
		},
		"empty role name": {
			displayName: "short",
			roleName: "",
			expectedUserPrefix: "v-short-",
		},
		"both display name and role name are empty": {
			displayName: "",
			roleName: "",
			expectedUserPrefix: "v-",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			usernameConfig := dbplugin.UsernameConfig{
				DisplayName: test.displayName,
				RoleName:    test.roleName,
			}

			username, err := credsProducer.GenerateUsername(usernameConfig)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !strings.HasPrefix(username, test.expectedUserPrefix) {
				t.Fatalf("Expected username to begin with %s. But got %s", test.expectedUserPrefix, username)
			}

			var expectedLen int
			if len(test.expectedUserPrefix) + RandomAlphaNumericLen + 1 + UnixTimestampLen < credsProducer.UsernameLen {
				expectedLen = len(test.expectedUserPrefix) + RandomAlphaNumericLen + 1 + UnixTimestampLen
			} else {
				expectedLen = credsProducer.UsernameLen
			}

			if len(username) != expectedLen {
				t.Fatalf("Expected username's length to be %d. But got %d", expectedLen, len(username))
			}
		})
	}
}
