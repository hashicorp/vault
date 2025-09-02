// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestGroupsList(t *testing.T) {
	b, storage := getBackend(t)

	groups := []string{
		"%20\\",
		"foo",
		"zfoo",
		"🙂",
		"foo/nested",
		"foo/even/more/nested",
	}

	for _, group := range groups {
		req := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "groups/" + group,
			Storage:   storage,
			Data: map[string]interface{}{
				"policies": []string{group + "_a", group + "_b"},
			},
		}

		resp, err := b.HandleRequest(context.Background(), req)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("err:%s resp:%#v\n", err, resp)
		}

	}

	for _, group := range groups {
		for _, upper := range []bool{false, true} {
			groupPath := group
			if upper {
				groupPath = strings.ToUpper(group)
			}
			req := &logical.Request{
				Operation: logical.ReadOperation,
				Path:      "groups/" + groupPath,
				Storage:   storage,
			}

			resp, err := b.HandleRequest(context.Background(), req)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("err:%s resp:%#v\n", err, resp)
			}
			if resp == nil {
				t.Fatal("unexpected nil response")
			}

			expected := []string{group + "_a", group + "_b"}

			if diff := deep.Equal(resp.Data["policies"].([]string), expected); diff != nil {
				t.Fatal(diff)
			}
		}
	}

	req := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "groups",
		Storage:   storage,
	}

	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("err:%s resp:%#v\n", err, resp)
	}

	if diff := deep.Equal(resp.Data["keys"].([]string), groups); diff != nil {
		t.Fatal(diff)
	}
}

func getBackend(t *testing.T) (logical.Backend, logical.Storage) {
	defaultLeaseTTLVal := time.Hour * 12
	maxLeaseTTLVal := time.Hour * 24

	config := &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Trace),

		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
		StorageView: &logical.InmemStorage{},
	}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("unable to create backend: %v", err)
	}

	return b, config.StorageView
}
