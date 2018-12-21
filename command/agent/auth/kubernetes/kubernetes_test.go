package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/helper/logging"
)

func TestKubernetesAuth_basic(t *testing.T) {
	testCases := map[string]struct {
		tokenPath string
		data      *mockJWTFile
		e         error
	}{
		"normal": {
			data: newMockJWTFile(jwtData),
		},
		"projected": {
			tokenPath: "/some/other/path",
			data:      newMockJWTFile(jwtProjectedData),
		},
		"not_found": {
			e: errors.New("open /var/run/secrets/kubernetes.io/serviceaccount/token: no such file or directory"),
		},
		"projected_not_found": {
			tokenPath: "/some/other/path",
			e:         errors.New("open /some/other/path: no such file or directory"),
		},
	}

	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			authCfg := auth.AuthConfig{
				Logger:    logging.NewVaultLogger(hclog.Trace),
				MountPath: "kubernetes",
				Config: map[string]interface{}{
					"role": "plugin-test",
				},
			}

			if tc.tokenPath != "" {
				authCfg.Config["token_path"] = tc.tokenPath
			}

			a, err := NewKubernetesAuthMethod(&authCfg)
			if err != nil {
				t.Fatal(err)
			}

			// Type assert to set the kubernetesMethod jwtData, to mock out reading
			// files from the pod.
			k := a.(*kubernetesMethod)
			if tc.data != nil {
				k.jwtData = tc.data
			}

			_, data, err := k.Authenticate(context.Background(), nil)
			if err != nil && tc.e == nil {
				t.Fatal(err)
			}

			if err != nil && !errwrap.Contains(err, tc.e.Error()) {
				t.Fatalf("expected \"no such file\" error, got: (%s)", err)
			}

			if err == nil && tc.e != nil {
				t.Fatal("expected error, but got none")
			}

			if tc.e == nil {
				authJWTraw, ok := data["jwt"]
				if !ok {
					t.Fatal("expected to find jwt data")
				}

				authJWT := authJWTraw.(string)
				token := jwtData
				if tc.tokenPath != "" {
					token = jwtProjectedData
				}
				if authJWT != token {
					t.Fatalf("error with auth tokens, expected (%s) got (%s)", token, authJWT)
				}
			}
		})
	}

}

// jwt for default service account
var jwtData = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InZhdWx0LWF1dGgtdG9rZW4tdDVwY24iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoidmF1bHQtYXV0aCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImQ3N2Y4OWJjLTkwNTUtMTFlNy1hMDY4LTA4MDAyNzZkOTliZiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnZhdWx0LWF1dGgifQ.HKUcqgrvan5ZC_mnpaMEx4RW3KrhfyH_u8G_IA2vUfkLK8tH3T7fJuJaPr7W6K_BqCrbeM5y3owszOzb4NR0Lvw6GBt2cFcen2x1Ua4Wokr0bJjTT7xQOIOw7UvUDyVS17wAurlfUnmWMwMMMOebpqj5K1t6GnyqghH1wPdHYRGX-q5a6C323dBCgM5t6JY_zTTaBgM6EkFq0poBaifmSMiJRPrdUN_-IgyK8fgQRiFYYkgS6DMIU4k4nUOb_sUFf5xb8vMs3SMteKiuWFAIt4iszXTj5IyBUNqe0cXA3zSY3QiNCV6bJ2CWW0Qf9WDtniT79VAqcR4GYaTC_gxjNA"

// jwt for projected service account
var jwtProjectedData = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJhdWQiOlsia3ViZXJuZXRlcy5kZWZhdWx0LnN2YyJdLCJleHAiOjE2MDMwNTM1NjMsImlhdCI6MTUzOTk4MTU2MywiaXNzIjoia3ViZXJuZXRlcy9zZXJ2aWNlYWNjb3VudCIsImt1YmVybmV0ZXMuaW8iOnsibmFtZXNwYWNlIjoiZGVmYXVsdCIsInBvZCI6eyJuYW1lIjoidmF1bHQiLCJ1aWQiOiIxMDA2YTA2Yy1kM2RmLTExZTgtOGZlMi0wODAwMjdlNTVlYTgifSwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImRlZmF1bHQiLCJ1aWQiOiJiMzg5YjNiMi1kMzAyLTExZTgtYjE0Yy0wODAwMjdlNTVlYTgifX0sIm5iZiI6MTUzOTk4MTU2Mywic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGVmYXVsdCJ9.byu3BpCbs0tzQvEBCRTayXF3-kV1Ey7YvStBcCwovfSl6evBze43FFaDps78HtdDAMszjE_yn55_1BMN87EzOZYsF3GBoPLWxkofxhPIy88wmPTpurBsSx-nCKdjf4ayXhTpqGG9gy0xlkUc_xL4pM3Q8XZiqYqwq_T0PHXOpSfdzVy1oabFSZXr5QTZ377v8bvrMgAVWJF_4vZsSMG3XVCK8KBWNRw4_wt6yOelVKE5OGLPJvNu1CFjEKh4HBFBcQnB_Sgpe1nPlnm5utp-1-OVfd7zopOGDAp_Pk_Apu8OPDdPSafn6HpzIeuhMtWXcv1K8ZhZYDLC1wLywZPNyw"

// mockJWTFile provides a mock ReadCloser struct to inject into
// kubernetesMethod.jwtData
type mockJWTFile struct {
	b *bytes.Buffer
}

var _ io.ReadCloser = &mockJWTFile{}

func (j *mockJWTFile) Read(p []byte) (n int, err error) {
	return j.b.Read(p)
}

func (j *mockJWTFile) Close() error { return nil }

func newMockJWTFile(s string) *mockJWTFile {
	return &mockJWTFile{
		b: bytes.NewBufferString(s),
	}
}
