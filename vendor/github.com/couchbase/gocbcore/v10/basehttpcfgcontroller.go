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

type baseHTTPConfigController struct {
	cfgMgr               *configManagementComponent
	confHTTPRetryDelay   time.Duration
	confHTTPRedialPeriod time.Duration
	confHTTPMaxWait      time.Duration
	httpComponent        *httpComponent
	bucketName           string
	endpointCallback     func(uint64) string

	looperStopSig chan struct{}

	fetchErr error
	errLock  sync.Mutex
}

type httpPollerProperties struct {
	confHTTPRetryDelay   time.Duration
	confHTTPRedialPeriod time.Duration
	confHTTPMaxWait      time.Duration
	httpComponent        *httpComponent
}

func newBaseHTTPConfigController(bucketName string, props httpPollerProperties, cfgMgr *configManagementComponent,
	endpointCallback func(uint64) string) *baseHTTPConfigController {
	return &baseHTTPConfigController{
		cfgMgr:               cfgMgr,
		confHTTPRedialPeriod: props.confHTTPRedialPeriod,
		confHTTPRetryDelay:   props.confHTTPRetryDelay,
		confHTTPMaxWait:      props.confHTTPMaxWait,
		httpComponent:        props.httpComponent,
		bucketName:           bucketName,

		looperStopSig: make(chan struct{}),

		endpointCallback: endpointCallback,
	}
}
func (hcc *baseHTTPConfigController) Error() error {
	hcc.errLock.Lock()
	defer hcc.errLock.Unlock()
	return hcc.fetchErr
}

func (hcc *baseHTTPConfigController) setError(err error) {
	hcc.errLock.Lock()
	hcc.fetchErr = err
	hcc.errLock.Unlock()
}

func (hcc *baseHTTPConfigController) Stop() {
	logDebugf("HTTP Looper stopping.")
	close(hcc.looperStopSig)
}

// Reset must never be called concurrently with the Stop or whilst the poll loop is running.
func (hcc *baseHTTPConfigController) Reset() {
	hcc.looperStopSig = make(chan struct{})
}

func (hcc *baseHTTPConfigController) DoLoop() {
	hcc.doLoop()
	logDebugf("HTTP Looper stopped.")
}

func (hcc *baseHTTPConfigController) doLoop() {
	waitPeriod := hcc.confHTTPRetryDelay
	maxConnPeriod := hcc.confHTTPRedialPeriod

	var iterNum uint64 = 1
	iterSawConfig := false

	logDebugf("HTTP Looper starting.")

	for {
		select {
		case <-hcc.looperStopSig:
			return
		default:
		}

		pickedSrv := hcc.endpointCallback(iterNum)

		if pickedSrv == "" {
			logDebugf("Pick Failed.")
			// All servers have been visited during this iteration

			if !iterSawConfig {
				logDebugf("Looper waiting...")
				// Wait for a period before trying again if there was a problem...
				// We also watch for the client being shut down.
				select {
				case <-hcc.looperStopSig:
					return
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
			uri := fmt.Sprintf("/pools/default/%s/%s", streamPath, url.PathEscape(hcc.bucketName))
			logDebugf("Requesting config from: %s/%s.", pickedSrv, uri)

			req := &httpRequest{
				Service:  MgmtService,
				Method:   "GET",
				Path:     uri,
				Endpoint: pickedSrv,
				UniqueID: uuid.New().String(),
				Deadline: time.Now().Add(hcc.confHTTPMaxWait),
			}

			var err error
			resp, err = hcc.httpComponent.DoInternalHTTPRequest(req, true)
			if err != nil {
				logWarnf("Failed to connect to host. %v", err)
				hcc.setError(err)
				return 0
			}

			if resp.StatusCode != 200 {
				err := resp.Body.Close()
				if err != nil {
					logErrorf("Socket close failed handling status code != 200 (%s)", err)
				}
				if resp.StatusCode == 401 {
					logWarnf("Failed to connect to host, bad auth.")
					hcc.setError(errAuthenticationFailure)
					return -1
				} else if resp.StatusCode == 404 {
					if is2x {
						logWarnf("Failed to connect to host, bad bucket.")
						hcc.setError(errAuthenticationFailure)
						return -1
					}

					return doConfigRequest(true)
				}
				logWarnf("Failed to connect to host, unexpected status code: %v.", resp.StatusCode)
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
}
