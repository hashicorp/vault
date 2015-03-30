package disk

import (
	"testing"

	"github.com/hashicorp/vault/command/token"
)

func TestCommand(t *testing.T) {
	token.TestProcess(t)
}

func TestHelperProcess(t *testing.T) {
	token.TestHelperProcessCLI(t, new(Command))
}
