// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldap

import (
	"context"
	"fmt"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
)

func PrepareTestContainer(t *testing.T, version string) (cleanup func(), cfg *ldaputil.ConfigEntry) {
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		// Currently set to "michelvocks" until https://github.com/rroemhild/docker-test-openldap/pull/14
		// has been merged.
		ImageRepo:     "docker.mirror.hashicorp.services/michelvocks/docker-test-openldap",
		ImageTag:      version,
		ContainerName: "ldap",
		Ports:         []string{"389/tcp"},
		// Env:        []string{"LDAP_DEBUG_LEVEL=384"},
	})
	if err != nil {
		t.Fatalf("could not start local LDAP docker container: %s", err)
	}

	cfg = new(ldaputil.ConfigEntry)
	cfg.UserDN = "ou=people,dc=planetexpress,dc=com"
	cfg.UserAttr = "cn"
	cfg.UserFilter = "({{.UserAttr}}={{.Username}})"
	cfg.BindDN = "cn=admin,dc=planetexpress,dc=com"
	cfg.BindPassword = "GoodNewsEveryone"
	cfg.GroupDN = "ou=people,dc=planetexpress,dc=com"
	cfg.GroupAttr = "cn"
	cfg.RequestTimeout = 60

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		connURL := fmt.Sprintf("ldap://%s:%d", host, port)
		cfg.Url = connURL
		logger := hclog.New(nil)
		client := ldaputil.Client{
			LDAP:   ldaputil.NewLDAP(),
			Logger: logger,
		}

		conn, err := client.DialLDAP(cfg)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		if _, err := client.GetUserBindDN(cfg, conn, "Philip J. Fry"); err != nil {
			return nil, err
		}

		return docker.NewServiceURLParse(connURL)
	})
	if err != nil {
		t.Fatalf("could not start local LDAP docker container: %s", err)
	}

	return svc.Cleanup, cfg
}
