package file

import (
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
)

func Factory(conf map[string]string) (audit.Backend, error) {
	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("path is required")
	}

	return &Backend{Path: path}, nil
}

// Backend is the audit backend for the file-based audit store.
//
// NOTE: This audit backend is currently very simple: it appends to a file.
// It doesn't do anything more at the moment to assist with rotation
// or reset the write cursor, this should be done in the future.
type Backend struct {
	Path string

	once sync.Once
	f    *os.File
}

func (b *Backend) LogRequest(auth *logical.Auth, req *logical.Request) error {
	if err := b.open(); err != nil {
		return err
	}

	// TODO
	return nil
}

func (b *Backend) LogResponse(
	auth *logical.Auth,
	req *logical.Request,
	resp *logical.Response,
	err error) error {
	if err := b.open(); err != nil {
		return err
	}

	// TODO
	return nil
}

func (b *Backend) open() error {
	if b.f != nil {
		return nil
	}

	var err error
	b.f, err = os.Create(b.Path)
	if err != nil {
		return err
	}

	return nil
}
