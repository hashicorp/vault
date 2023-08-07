// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"context"
	"net/url"
	"strings"
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
			"connection_url":   "postgres://{{username}}:{{password}}@localhost:5432/mydb",
			"username":         tc.Username,
			"password":         tc.Password,
			"disable_escaping": false,
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

func TestSQLDisableEscaping(t *testing.T) {
	testCases := []struct {
		Username        string
		Password        string
		DisableEscaping bool
	}{
		{"mssql{0}", "password{0}", true},
		{"mssql{0}", "password{0}", false},
		{"ms\"sql\"", "pass\"word\"", true},
		{"ms\"sql\"", "pass\"word\"", false},
		{"ms'sq;l", "pass'wor;d", true},
		{"ms'sq;l", "pass'wor;d", false},
	}
	for _, tc := range testCases {
		t.Logf("username %q password %q disable_escaping %t", tc.Username, tc.Password, tc.DisableEscaping)

		sql := &SQLConnectionProducer{}
		ctx := context.Background()
		conf := map[string]interface{}{
			"connection_url":   "server=localhost;port=1433;user id={{username}};password={{password}};database=mydb;",
			"username":         tc.Username,
			"password":         tc.Password,
			"disable_escaping": tc.DisableEscaping,
		}
		_, err := sql.Init(ctx, conf, false)
		if err != nil {
			t.Errorf("Init error on %q %q: %+v", tc.Username, tc.Password, err)
		} else {
			if tc.DisableEscaping {
				if !strings.Contains(sql.ConnectionURL, tc.Username) || !strings.Contains(sql.ConnectionURL, tc.Password) {
					t.Errorf("Raw username and/or password missing from ConnectionURL")
				}
			} else {
				if strings.Contains(sql.ConnectionURL, tc.Username) || strings.Contains(sql.ConnectionURL, tc.Password) {
					t.Errorf("Raw username and/or password was present in ConnectionURL")
				}
			}
		}
	}
}
