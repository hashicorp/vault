package kubernetes

import (
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestKubernetesBackend(t *testing.T) {
	namespace := os.Getenv("KUBERNETES_NAMESPACE")
	if namespace == "" {
		namespace = "vault"
	}

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewKubernetesBackend(map[string]string{
		"namespace": namespace,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}
