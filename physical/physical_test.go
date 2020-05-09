package physical

import (
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestReusableInmemStorage(t *testing.T) {

	logger := logging.NewVaultLogger(hclog.Debug).Named(t.Name())

	_, cleanup := teststorage.MakeReusableStorage(
		t, logger, teststorage.MakeInmemBackend(t, logger))
	defer cleanup()
}
