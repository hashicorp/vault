package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestHTTP_Fallback(t *testing.T) {
	handler1 := http.NewServeMux()
	handler2 := http.NewServeMux()
	handler3 := http.NewServeMux()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		ClusterAddr: "https://127.0.0.1:8382",
	}

	// Chicken-and-egg: Handler needs a core. So we create handlers first, then
	// add routes chained to a Handler-created handler.
	cores := vault.TestCluster(t, []http.Handler{handler1, handler2, handler3}, coreConfig, true)
	for _, core := range cores {
		defer core.CloseListeners()
	}
	handler1.Handle("/", Handler(cores[0].Core))
	handler2.Handle("/", Handler(cores[1].Core))
	handler3.Handle("/", Handler(cores[2].Core))

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	root := cores[0].Root

	addrs := []string{"http://127.0.0.1:8206", "http://127.0.0.1:8208"}

	for _, addr := range addrs {
		config := api.DefaultConfig()
		config.Address = addr
		client, err := api.NewClient(config)
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(root)

		secret, err := client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			t.Fatal("secret is nil")
		}
		if secret.Data["id"].(string) != root {
			t.Fatal("token mismatch")
		}
	}
}

// This function recreates the fuzzy testing from transit to pipe a large
// number of requests from the standbys to the active node.
func TestHTTP_Forwarding_Stress(t *testing.T) {
	testPlaintext := "the quick brown fox"
	testPlaintextB64 := "dGhlIHF1aWNrIGJyb3duIGZveA=="

	handler1 := http.NewServeMux()
	handler2 := http.NewServeMux()
	handler3 := http.NewServeMux()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	}

	// Chicken-and-egg: Handler needs a core. So we create handlers first, then
	// add routes chained to a Handler-created handler.
	cores := vault.TestCluster(t, []http.Handler{handler1, handler2, handler3}, coreConfig, true)
	for _, core := range cores {
		defer core.CloseListeners()
	}
	handler1.Handle("/", Handler(cores[0].Core))
	handler2.Handle("/", Handler(cores[1].Core))
	handler3.Handle("/", Handler(cores[2].Core))

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	root := cores[0].Root

	wg := sync.WaitGroup{}

	funcs := []string{"encrypt", "decrypt", "rotate", "change_min_version"}
	keys := []string{"test1", "test2", "test3"}
	hosts := []string{"http://127.0.0.1:8206/v1/transit/", "http://127.0.0.1:8208/v1/transit/"}

	transport := cleanhttp.DefaultPooledTransport()

	//core.Logger().Printf("[TRACE] mounting transit")
	req, err := http.NewRequest("POST", "http://127.0.0.1:8202/v1/sys/mounts/transit",
		bytes.NewBuffer([]byte("{\"type\": \"transit\"}")))
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return fmt.Errorf("redirects not allowed in this test")
		},
	}
	req.Header.Set(AuthHeaderName, root)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	//core.Logger().Printf("[TRACE] done mounting transit")

	var totalOps int64
	var successfulOps int64
	var key1ver int64 = 1
	var key2ver int64 = 1
	var key3ver int64 = 1

	// This is the goroutine loop
	doFuzzy := func(id int) {
		// Check for panics, otherwise notify we're done
		defer func() {
			if err := recover(); err != nil {
				core.Logger().Printf("[ERR] got a panic: %v", err)
				t.Fail()
			}
			wg.Done()
		}()

		// Holds the latest encrypted value for each key
		latestEncryptedText := map[string]string{}

		startTime := time.Now()
		client := &http.Client{
			Transport: transport,
		}

		var chosenFunc, chosenKey, chosenHost string

		doReq := func(method, url string, body io.Reader) (*http.Response, error) {
			req, err := http.NewRequest(method, url, body)
			if err != nil {
				return nil, err
			}
			req.Header.Set(AuthHeaderName, root)
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		doResp := func(resp *http.Response) (*api.Secret, error) {
			if resp == nil {
				return nil, fmt.Errorf("nil response")
			}
			defer resp.Body.Close()

			// Make sure we weren't redirected
			if resp.StatusCode > 300 && resp.StatusCode < 400 {
				return nil, fmt.Errorf("got status code %d, resp was %#v", resp.StatusCode, *resp)
			}

			result := &api.Response{Response: resp}
			err = result.Error()
			if err != nil {
				return nil, err
			}

			secret, err := api.ParseSecret(result.Body)
			if err != nil {
				return nil, err
			}

			return secret, nil
		}

		for _, chosenHost := range hosts {
			for _, chosenKey := range keys {
				// Try to write the key to make sure it exists
				_, err := doReq("POST", chosenHost+"keys/"+chosenKey, bytes.NewBuffer([]byte("{}")))
				if err != nil {
					panic(err)
				}
			}
		}

		//core.Logger().Printf("[TRACE] Starting %d", id)
		for {
			// Stop after 10 seconds
			if time.Now().Sub(startTime) > 10*time.Second {
				return
			}

			atomic.AddInt64(&totalOps, 1)

			// Pick a function and a key
			chosenFunc = funcs[rand.Int()%len(funcs)]
			chosenKey = keys[rand.Int()%len(keys)]
			chosenHost = hosts[rand.Int()%len(hosts)]

			switch chosenFunc {
			// Encrypt our plaintext and store the result
			case "encrypt":
				//core.Logger().Printf("[TRACE] %s, %s, %d", chosenFunc, chosenKey, id)
				resp, err := doReq("POST", chosenHost+"encrypt/"+chosenKey, bytes.NewBuffer([]byte(fmt.Sprintf("{\"plaintext\": \"%s\"}", testPlaintextB64))))
				if err != nil {
					panic(err)
				}

				secret, err := doResp(resp)
				if err != nil {
					panic(err)
				}

				latest := secret.Data["ciphertext"].(string)
				if latest == "" {
					panic(fmt.Errorf("bad ciphertext"))
				}
				latestEncryptedText[chosenKey] = secret.Data["ciphertext"].(string)

				atomic.AddInt64(&successfulOps, 1)

			// Decrypt the ciphertext and compare the result
			case "decrypt":
				ct := latestEncryptedText[chosenKey]
				if ct == "" {
					atomic.AddInt64(&successfulOps, 1)
					continue
				}

				//core.Logger().Printf("[TRACE] %s, %s, %d", chosenFunc, chosenKey, id)
				resp, err := doReq("POST", chosenHost+"decrypt/"+chosenKey, bytes.NewBuffer([]byte(fmt.Sprintf("{\"ciphertext\": \"%s\"}", ct))))
				if err != nil {
					panic(err)
				}

				secret, err := doResp(resp)
				if err != nil {
					// This could well happen since the min version is jumping around
					if strings.Contains(err.Error(), transit.ErrTooOld) {
						atomic.AddInt64(&successfulOps, 1)
						continue
					}
					panic(err)
				}

				ptb64 := secret.Data["plaintext"].(string)
				pt, err := base64.StdEncoding.DecodeString(ptb64)
				if err != nil {
					panic(fmt.Errorf("got an error decoding base64 plaintext: %v", err))
				}
				if string(pt) != testPlaintext {
					panic(fmt.Errorf("got bad plaintext back: %s", pt))
				}

				atomic.AddInt64(&successfulOps, 1)

			// Rotate to a new key version
			case "rotate":
				//core.Logger().Printf("[TRACE] %s, %s, %d", chosenFunc, chosenKey, id)
				_, err := doReq("POST", chosenHost+"keys/"+chosenKey+"/rotate", bytes.NewBuffer([]byte("{}")))
				if err != nil {
					panic(err)
				}
				switch chosenKey {
				case "test1":
					atomic.AddInt64(&key1ver, 1)
				case "test2":
					atomic.AddInt64(&key2ver, 1)
				case "test3":
					atomic.AddInt64(&key3ver, 1)
				}
				atomic.AddInt64(&successfulOps, 1)

			// Change the min version, which also tests the archive functionality
			case "change_min_version":
				var latestVersion int64
				switch chosenKey {
				case "test1":
					latestVersion = atomic.LoadInt64(&key1ver)
				case "test2":
					latestVersion = atomic.LoadInt64(&key2ver)
				case "test3":
					latestVersion = atomic.LoadInt64(&key3ver)
				}

				setVersion := (rand.Int63() % latestVersion) + 1

				//core.Logger().Printf("[TRACE] %s, %s, %d, new min version %d", chosenFunc, chosenKey, id, setVersion)

				_, err := doReq("POST", chosenHost+"keys/"+chosenKey+"/config", bytes.NewBuffer([]byte(fmt.Sprintf("{\"min_decryption_version\": %d}", setVersion))))
				if err != nil {
					panic(err)
				}

				atomic.AddInt64(&successfulOps, 1)
			}
		}
	}

	// Spawn 20 of these workers for 10 seconds
	for i := 0; i < 20; i++ {
		wg.Add(1)
		//core.Logger().Printf("[TRACE] spawning %d", i)
		go doFuzzy(i)
	}

	// Wait for them all to finish
	wg.Wait()

	core.Logger().Printf("[TRACE] total operations tried: %d, total successful: %d", totalOps, successfulOps)
}
