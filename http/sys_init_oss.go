//go:build !enterprise

package http

import (
	"net/http"

	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/vault"
)

func init() {
	handleSysEnableUnsealRecovery = func(core *vault.Core) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := core.InitializeUnsealRecovery(r.Context())
			if err != nil {
				if _, ok := err.(errutil.UserError); ok {
					respondError(w, http.StatusBadRequest, err)
				} else {
					respondError(w, http.StatusInternalServerError, err)
				}
			}
			respondOk(w, nil)
		})
	}
}
