//go:build !enterprise

package http

import (
	"net/http"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func adjustRequest(c *vault.Core, listener *configutil.Listener, r *http.Request) (*http.Request, int, error) {
	return r, 0, nil
}
