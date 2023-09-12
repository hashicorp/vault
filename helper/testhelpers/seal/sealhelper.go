// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package sealhelper

import (
	"path"
	"strconv"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/mitchellh/go-testing-interface"
)

type TransitSealServer struct {
	*vault.TestCluster
}

func NewTransitSealServer(t testing.T, idx int) *TransitSealServer {
	conf := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	}
	opts := &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: http.Handler,
		Logger:      corehelpers.NewTestLogger(t).Named("transit-seal" + strconv.Itoa(idx)),
	}
	teststorage.InmemBackendSetup(conf, opts)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()

	if err := cluster.Cores[0].Client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	}); err != nil {
		t.Fatal(err)
	}

	return &TransitSealServer{cluster}
}

func (tss *TransitSealServer) MakeKey(t testing.T, key string) {
	client := tss.Cores[0].Client
	if _, err := client.Logical().Write(path.Join("transit", "keys", key), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write(path.Join("transit", "keys", key, "config"), map[string]interface{}{
		"deletion_allowed": true,
	}); err != nil {
		t.Fatal(err)
	}
}

func (tss *TransitSealServer) MakeSeal(t testing.T, key string) (vault.Seal, error) {
	client := tss.Cores[0].Client
	wrapperConfig := map[string]string{
		"address":     client.Address(),
		"token":       client.Token(),
		"mount_path":  "transit",
		"key_name":    key,
		"tls_ca_cert": tss.CACertPEMFile,
	}
	transitSealWrapper, _, err := configutil.GetTransitKMSFunc(&configutil.KMS{Config: wrapperConfig})
	if err != nil {
		t.Fatalf("error setting wrapper config: %v", err)
	}

	access, err := seal.NewAccessFromSealWrappers(tss.Logger, 1, true, []seal.SealWrapper{
		{
			Wrapper:  transitSealWrapper,
			Priority: 1,
			Name:     "transit",
		},
	})
	if err != nil {
		return nil, err
	}
	return vault.NewAutoSeal(access), nil
}
