package credsutil

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

func TestGenerateUsername(t *testing.T) {
	credsProducer := &SQLCredentialsProducer{
		DisplayNameLen: 30,
		RoleNameLen:    8,
		UsernameLen:    63,
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

			if len(username) > credsProducer.UsernameLen {
				t.Fatalf("Maximum username's length is %d. But got %d", credsProducer.UsernameLen, len(username))
			}
		})
	}
}
