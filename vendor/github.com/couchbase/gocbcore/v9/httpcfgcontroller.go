package gocbcore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type configStreamBlock struct {
	Bytes []byte
}

func (i *configStreamBlock) UnmarshalJSON(data []byte) error {
	i.Bytes = make([]byte, len(data))
	copy(i.Bytes, data)
	return nil
}

func hostnameFromURI(uri string) string {
	uriInfo, err := url.Parse(uri)
	if err != nil {
		return uri
	}

	hostname, err := hostFromHostPort(uriInfo.Host)
	if err != nil {
		return uri
	}

	return hostname
}

type httpConfigController struct {
	muxer                *httpMux
	cfgMgr               *configManagementComponent
	confHTTPRetryDelay   time.Duration
	confHTTPRedialPeriod time.Duration
	httpComponent        *httpComponent
	bucketName           string

	looperStopSig chan struct{}
	looperDoneSig chan struct{}

	fetchErr error
	errLock  sync.Mutex
}

type httpPollerProperties struct {
	confHTTPRetryDelay   time.Duration
	confHTTPRedialPeriod time.Duration
	httpComponent        *httpComponent
}

func newHTTPConfigController(bucketName string, props httpPollerProperties, muxer *httpMux,
	cfgMgr *configManagementComponent) *httpConfigController {
	return &httpConfigController{
		muxer:                muxer,
		cfgMgr:               cfgMgr,
		confHTTPRedialPeriod: props.confHTTPRedialPeriod,
		confHTTPRetryDelay:   props.confHTTPRetryDelay,
		httpComponent:        props.httpComponent,
		bucketName:           bucketName,

		looperStopSig: make(chan struct{}),
		looperDoneSig: make(chan struct{}),
	}
}

func (hcc *httpConfigController) Error() error {
	hcc.errLock.Lock()
	defer hcc.errLock.Unlock()
	return hcc.fetchErr
}

func (hcc *httpConfigController) setError(err error) {
	hcc.errLock.Lock()
	hcc.fetchErr = err
	hcc.errLock.Unlock()
}

func (hcc *httpConfigController) Pause(paused bool) {
}

func (hcc *httpConfigController) Done() chan struct{} {
	return hcc.looperDoneSig
}

func (hcc *httpConfigController) Stop() {
	close(hcc.looperStopSig)
}

func (hcc *httpConfigController) Reset() {
	hcc.looperStopSig = make(chan struct{})
	hcc.looperDoneSig = make(chan struct{})
}

func (hcc *httpConfigController) DoLoop() {
	waitPeriod := hcc.confHTTPRetryDelay
	maxConnPeriod := hcc.confHTTPRedialPeriod

	var iterNum uint64 = 1
	iterSawConfig := false
	seenNodes := make(map[string]uint64)

	logDebugf("HTTP Looper starting.")

Looper:
	for {
		select {
		case <-hcc.looperStopSig:
			break Looper
		default:
		}

		var pickedSrv string
		for _, srv := range hcc.muxer.MgmtEps() {
			if seenNodes[srv] >= iterNum {
				continue
			}
			pickedSrv = srv
			break
		}

		if pickedSrv == "" {
			logDebugf("Pick Failed.")
			// All servers have been visited during this iteration

			if !iterSawConfig {
				logDebugf("Looper waiting...")
				// Wait for a period before trying again if there was a problem...
				// We also watch for the client being shut down.
				select {
				case <-hcc.looperStopSig:
					break Looper
				case <-time.After(waitPeriod):
				}
			}
			logDebugf("Looping again.")
			// Go to next iteration and try all servers again
			iterNum++
			iterSawConfig = false
			continue
		}

		logDebugf("Http Picked: %s.", pickedSrv)

		seenNodes[pickedSrv] = iterNum

		hostname := hostnameFromURI(pickedSrv)
		logDebugf("HTTP Hostname: %s.", hostname)

		var resp *HTTPResponse
		// 1 on success, 0 on failure for node, -1 for generic failure
		var doConfigRequest func(bool) int

		doConfigRequest = func(is2x bool) int {
			streamPath := "bs"
			if is2x {
				streamPath = "bucketsStreaming"
			}
			// HTTP request time!
			uri := fmt.Sprintf("/pools/default/%s/%s", streamPath, hcc.bucketName)
			logDebugf("Requesting config from: %s/%s.", pickedSrv, uri)

			req := &httpRequest{
				Service:  MgmtService,
				Method:   "GET",
				Path:     uri,
				Endpoint: pickedSrv,
				UniqueID: uuid.New().String(),
			}

			var err error
			resp, err = hcc.httpComponent.DoInternalHTTPRequest(req, true)
			if err != nil {
				logDebugf("Failed to connect to host. %v", err)
				hcc.setError(err)
				return 0
			}

			if resp.StatusCode != 200 {
				err := resp.Body.Close()
				if err != nil {
					logErrorf("Socket close failed handling status code != 200 (%s)", err)
				}
				if resp.StatusCode == 401 {
					logDebugf("Failed to connect to host, bad auth.")
					hcc.setError(errAuthenticationFailure)
					return -1
				} else if resp.StatusCode == 404 {
					if is2x {
						logDebugf("Failed to connect to host, bad bucket.")
						hcc.setError(errAuthenticationFailure)
						return -1
					}

					return doConfigRequest(true)
				}
				logDebugf("Failed to connect to host, unexpected status code: %v.", resp.StatusCode)
				hcc.setError(errCliInternalError)
				return 0
			}
			hcc.setError(nil)
			return 1
		}

		switch doConfigRequest(false) {
		case 0:
			continue
		case -1:
			continue
		}

		logDebugf("Connected.")

		var autoDisconnected int32

		// Autodisconnect eventually
		go func() {
			select {
			case <-time.After(maxConnPeriod):
			case <-hcc.looperStopSig:
			}

			logDebugf("Automatically resetting our HTTP connection")

			atomic.StoreInt32(&autoDisconnected, 1)

			err := resp.Body.Close()
			if err != nil {
				logErrorf("Socket close failed during auto-dc (%s)", err)
			}
		}()

		dec := json.NewDecoder(resp.Body)
		configBlock := new(configStreamBlock)
		for {
			err := dec.Decode(configBlock)
			if err != nil {
				if atomic.LoadInt32(&autoDisconnected) == 1 {
					// If we know we intentionally disconnected, we know we do not
					// need to close the client, nor log an error, since this was
					// expected behaviour
					break
				}

				logWarnf("Config block decode failure (%s)", err)

				if err != io.EOF {
					err = resp.Body.Close()
					if err != nil {
						logErrorf("Socket close failed after decode fail (%s)", err)
					}
				}

				break
			}

			logDebugf("Got Block: %v", string(configBlock.Bytes))

			bkCfg, err := parseConfig(configBlock.Bytes, hostname)
			if err != nil {
				logDebugf("Got error while parsing config: %v", err)

				err = resp.Body.Close()
				if err != nil {
					logErrorf("Socket close failed after parsing fail (%s)", err)
				}

				break
			}

			logDebugf("Got Config.")

			iterSawConfig = true
			logDebugf("HTTP Config Update")
			hcc.cfgMgr.OnNewConfig(bkCfg)
		}

		logDebugf("HTTP, Setting %s to iter %d", pickedSrv, iterNum)
	}

	close(hcc.looperDoneSig)
}
