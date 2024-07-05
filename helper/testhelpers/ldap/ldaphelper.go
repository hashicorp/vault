// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/cap/ldap"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
)

func PrepareTestContainer(t *testing.T, version string) (cleanup func(), cfg *ldaputil.ConfigEntry) {
	// note: this image isn't supported on arm64 architecture in CI.
	// but if you're running on Apple Silicon, feel free to comment out the code below locally.
	if strings.Contains(runtime.GOARCH, "arm") {
		t.Skip("Skipping, as this image is not supported on ARM architectures")
	}

	logsWriter := bytes.NewBuffer([]byte{})

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "ghcr.io/rroemhild/docker-test-openldap",
		ImageTag:      version,
		ContainerName: "ldap",
		Ports:         []string{"10389/tcp"},
		// Env:        []string{"LDAP_DEBUG_LEVEL=384"},
		LogStderr: logsWriter,
		LogStdout: logsWriter,
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
	cfg.MaximumPageSize = 1000

	var started bool

	for i := 0; i < 3; i++ {
		svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
			connURL := fmt.Sprintf("ldap://%s:%d", host, port)
			cfg.Url = connURL

			client, err := ldap.NewClient(ctx, ldaputil.ConvertConfig(cfg))
			if err != nil {
				return nil, err
			}

			defer client.Close(ctx)

			_, err = client.Authenticate(ctx, "Philip J. Fry", "fry")
			if err != nil {
				return nil, err
			}

			return docker.NewServiceURLParse(connURL)
		})
		if err != nil {
			t.Logf("could not start local LDAP docker container: %s", err)
			t.Log("Docker container logs: ")
			t.Log(logsWriter.String())
			continue
		}

		started = true
		cleanup = func() {
			if t.Failed() {
				t.Log(logsWriter.String())
			}
			svc.Cleanup()
		}
		break
	}

	if !started {
		t.FailNow()
	}

	return cleanup, cfg
}
