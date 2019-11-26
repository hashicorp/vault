package connutil

import (
	"context"
	"net/url"
	"testing"
)

func TestSQLPasswordChars(t *testing.T) {
	testCases := []struct {
		Username string
		Password string
	}{
		{"postgres", "password{0}"},
		{"postgres", "pass:word"},
		{"postgres", "pass/word"},
		{"postgres", "p@ssword"},
		{"postgres", "pass\"word\""},
		// Much to my surprise, CREATE USER "{{password}}" PASSWORD 'foo' worked.
		{"{{password}}", "foo"},
		{"user", "{{username}}"},
	}
	for _, tc := range testCases {
		t.Logf("username %q password %q", tc.Username, tc.Password)

		sql := &SQLConnectionProducer{}
		ctx := context.Background()
		conf := map[string]interface{}{
			"connection_url": "postgres://{{username}}:{{password}}@localhost:5432/mydb",
			"username":       tc.Username,
			"password":       tc.Password,
		}
		_, err := sql.Init(ctx, conf, false)
		if err != nil {
			t.Errorf("Init error on %q %q: %+v", tc.Username, tc.Password, err)
		} else {
			// This jumps down a few layers...
			// Connection() uses sql.Open uses lib/pq uses net/url.Parse
			u, err := url.Parse(sql.ConnectionURL)
			if err != nil {
				t.Errorf("URL parse error on %q %q: %+v", tc.Username, tc.Password, err)
			} else {
				username := u.User.Username()
				password, pPresent := u.User.Password()
				if username != tc.Username {
					t.Errorf("Parsed username %q != original username %q", username, tc.Username)
				}
				if !pPresent {
					t.Errorf("Password %q not present", tc.Password)
				} else if password != tc.Password {
					t.Errorf("Parsed password %q != original password %q", password, tc.Password)
				}
			}
		}
	}
}
