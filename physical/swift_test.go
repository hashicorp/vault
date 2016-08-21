package physical

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/logformat"
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
	tenant := os.Getenv("OS_TENANT_NAME")

	ts := time.Now().UnixNano()
	container := fmt.Sprintf("vault-test-%d", ts)

	cleaner := swift.Connection{
		UserName:  username,
		ApiKey:    password,
		AuthUrl:   authUrl,
		Tenant:    tenant,
		Transport: cleanhttp.DefaultPooledTransport(),
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

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("swift", logger, map[string]string{
		"username":  username,
		"password":  password,
		"container": container,
		"auth_url":  authUrl,
		"tenant":    tenant,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}
