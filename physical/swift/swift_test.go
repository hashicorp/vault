package swift

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
	"github.com/ncw/swift"
)

func TestSwiftBackend(t *testing.T) {
	if os.Getenv("OS_USERNAME") == "" || os.Getenv("OS_PASSWORD") == "" ||
		os.Getenv("OS_AUTH_URL") == "" {
		t.SkipNow()
	}
	username := os.Getenv("OS_USERNAME")
	password := os.Getenv("OS_PASSWORD")
	authUrl := os.Getenv("OS_AUTH_URL")
	project := os.Getenv("OS_PROJECT_NAME")
	domain := os.Getenv("OS_USER_DOMAIN_NAME")
	projectDomain := os.Getenv("OS_PROJECT_DOMAIN_NAME")
	region := os.Getenv("OS_REGION_NAME")
	tenantID := os.Getenv("OS_TENANT_ID")

	ts := time.Now().UnixNano()
	container := fmt.Sprintf("vault-test-%d", ts)

	cleaner := swift.Connection{
		Domain:       domain,
		UserName:     username,
		ApiKey:       password,
		AuthUrl:      authUrl,
		Tenant:       project,
		TenantDomain: projectDomain,
		Region:       region,
		TenantId:     tenantID,
		Transport:    cleanhttp.DefaultPooledTransport(),
	}

	err := cleaner.Authenticate()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = cleaner.ContainerCreate(container, nil)
	if nil != err {
		t.Fatalf("Unable to create test container '%s': %v", container, err)
	}
	defer func() {
		newObjects, err := cleaner.ObjectNamesAll(container, nil)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		for _, o := range newObjects {
			err := cleaner.ObjectDelete(container, o)
			if err != nil {
				t.Fatalf("err: %s", err)
			}
		}
		err = cleaner.ContainerDelete(container)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewSwiftBackend(map[string]string{
		"username":       username,
		"password":       password,
		"container":      container,
		"auth_url":       authUrl,
		"project":        project,
		"domain":         domain,
		"project-domain": projectDomain,
		"tenant_id":      tenantID,
		"region":         region,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}
