package ldap

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/ory/dockertest"
	"testing"
)

func PrepareTestContainer(t *testing.T, version string) (cleanup func(), cfg *ldaputil.ConfigEntry) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "rroemhild/test-openldap",
		Tag:        version,
		Privileged: true,
		//Env:        []string{"LDAP_DEBUG_LEVEL=384"},
	}
	resource, err := pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local LDAP %s docker container: %s", version, err)
	}

	cleanup = func() {
		docker.CleanupResource(t, pool, resource)
	}

	//pool.MaxWait = time.Second
	// exponential backoff-retry
	if err = pool.Retry(func() error {
		logger := hclog.New(nil)
		client := ldaputil.Client{
			LDAP:   ldaputil.NewLDAP(),
			Logger: logger,
		}

		cfg = new(ldaputil.ConfigEntry)
		cfg.Url = fmt.Sprintf("ldap://localhost:%s", resource.GetPort("389/tcp"))
		cfg.UserDN = "ou=people,dc=planetexpress,dc=com"
		cfg.UserAttr = "cn"
		cfg.BindDN = "cn=admin,dc=planetexpress,dc=com"
		cfg.BindPassword = "GoodNewsEveryone"
		cfg.GroupDN = "ou=people,dc=planetexpress,dc=com"
		cfg.GroupAttr = "memberOf"
		conn, err := client.DialLDAP(cfg)
		if err != nil {
			return err
		}
		defer conn.Close()

		if _, err := client.GetUserBindDN(cfg, conn, "Philip J. Fry"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to docker: %s", err)
	}

	return cleanup, cfg
}
