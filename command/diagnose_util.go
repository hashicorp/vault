package command

import (
	"github.com/hashicorp/vault/command/server"
)

type DiagnoseObserver interface {
	Success(key string)
	Error(key string, err error)
	ConfigCreated(config *server.Config)
	IsEnabled() bool
}

type NullDiagnoseObserver struct {
}

func (n *NullDiagnoseObserver) Success(key string) {
}

func (n *NullDiagnoseObserver) Error(key string, err error) {
}

func (n *NullDiagnoseObserver) ConfigCreated(config *server.Config) {
}

func (n *NullDiagnoseObserver) IsEnabled() bool {
	return false
}
