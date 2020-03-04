package http

import (
	"fmt"
	"net/http"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/monitor"
	"github.com/hashicorp/vault/vault"
)

func handleSysMonitor(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ll := r.URL.Query().Get("log_level")
		if ll == "" {
			ll = "INFO"
		}
		logLevel := log.LevelFromString(ll)

		if logLevel == log.NoLevel {
			respondError(w, http.StatusBadRequest, fmt.Errorf("unknown log level"))
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			respondError(w, http.StatusBadRequest, fmt.Errorf("streaming not supported"))
			return
		}

		isJson := core.SanitizedConfig()["log_format"] == "json"
		monitor := monitor.New(512, core.Logger(), &log.LoggerOptions{
			Level:      logLevel,
			JSONFormat: isJson,
		})

		logCh := monitor.Start()
		w.WriteHeader(http.StatusOK)

		// 0 byte write is needed before the Flush call so that if we are using
		// a gzip stream it will go ahead and write out the HTTP response header
		w.Write([]byte(""))
		flusher.Flush()

		// Stream logs until the connection is closed.
		for {
			select {
			case <-r.Context().Done():
				monitor.Stop()
				return
			case log := <-logCh:
				fmt.Fprint(w, string(log))
				flusher.Flush()
			}
		}
	})
}
