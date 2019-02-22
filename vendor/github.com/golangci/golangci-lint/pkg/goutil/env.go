package goutil

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

type Env struct {
	vars   map[string]string
	log    logutils.Log
	debugf logutils.DebugFunc
}

func NewEnv(log logutils.Log) *Env {
	return &Env{
		vars:   map[string]string{},
		log:    log,
		debugf: logutils.Debug("env"),
	}
}

func (e *Env) Discover(ctx context.Context) error {
	out, err := exec.CommandContext(ctx, "go", "env", "-json").Output()
	if err != nil {
		return errors.Wrap(err, "failed to run 'go env'")
	}

	if err = json.Unmarshal(out, &e.vars); err != nil {
		return errors.Wrap(err, "failed to parse go env json")
	}

	e.debugf("Read go env: %#v", e.vars)
	return nil
}

func (e Env) Get(k string) string {
	envValue := os.Getenv(k)
	if envValue != "" {
		return envValue
	}

	return e.vars[k]
}
