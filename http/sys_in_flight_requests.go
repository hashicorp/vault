package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/vault"
)

func handleUnAuthenticatedInFlightRequest(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
		now := time.Now()
		currentInFlightReqMap := make(map[string]interface{})
		syncMapRangeResult := true
		core.RangeInFlightReqData(func(key, value interface{}) bool {
			v, ok := value.(*vault.InFlightReqData)
			if !ok {
				syncMapRangeResult = false
				return false
			}
			// don't report the request to the in-flight-req path
			if v.ReqPath != "/v1/sys/in-flight-req" {
				v.Duration = fmt.Sprintf("%v microseconds", now.Sub(v.StartTime).Microseconds())
				currentInFlightReqMap[key.(string)] = v
			}

			return true
		})

		// TODO: should an error be returned here? and if so, what status code should be returned? 500 or 400?
		if !syncMapRangeResult {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("failed to read recorded in-flight requests"))
			return
		}

		respondOk(w, currentInFlightReqMap)
	})
}
