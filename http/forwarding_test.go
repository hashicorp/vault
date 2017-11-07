package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/http2"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestHTTP_Fallback_Bad_Address(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		ClusterAddr: "https://127.3.4.1:8382",
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	addrs := []string{
		fmt.Sprintf("https://127.0.0.1:%d", cores[1].Listeners[0].Address.Port),
		fmt.Sprintf("https://127.0.0.1:%d", cores[2].Listeners[0].Address.Port),
	}

	for _, addr := range addrs {
		config := api.DefaultConfig()
		config.Address = addr
		config.HttpClient.Transport.(*http.Transport).TLSClientConfig = cores[0].TLSConfig

		client, err := api.NewClient(config)
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(cluster.RootToken)

		secret, err := client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			t.Fatal("secret is nil")
		}
		if secret.Data["id"].(string) != cluster.RootToken {
			t.Fatal("token mismatch")
		}
	}
}

func TestHTTP_Fallback_Disabled(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		ClusterAddr: "empty",
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	addrs := []string{
		fmt.Sprintf("https://127.0.0.1:%d", cores[1].Listeners[0].Address.Port),
		fmt.Sprintf("https://127.0.0.1:%d", cores[2].Listeners[0].Address.Port),
	}

	for _, addr := range addrs {
		config := api.DefaultConfig()
		config.Address = addr
		config.HttpClient.Transport.(*http.Transport).TLSClientConfig = cores[0].TLSConfig

		client, err := api.NewClient(config)
		if err != nil {
			t.Fatal(err)
		}
		client.SetToken(cluster.RootToken)

		secret, err := client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			t.Fatal("secret is nil")
		}
		if secret.Data["id"].(string) != cluster.RootToken {
			t.Fatal("token mismatch")
		}
	}
}

// This function recreates the fuzzy testing from transit to pipe a large
// number of requests from the standbys to the active node.
func TestHTTP_Forwarding_Stress(t *testing.T) {
	testHTTP_Forwarding_Stress_Common(t, false, 50)
	testHTTP_Forwarding_Stress_Common(t, true, 50)
}

func testHTTP_Forwarding_Stress_Common(t *testing.T, parallel bool, num uint64) {
	testPlaintext := "the quick brown fox"
	testPlaintextB64 := "dGhlIHF1aWNrIGJyb3duIGZveA=="

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	wg := sync.WaitGroup{}

	funcs := []string{"encrypt", "decrypt", "rotate", "change_min_version"}
	keys := []string{"test1", "test2", "test3"}

	hosts := []string{
		fmt.Sprintf("https://127.0.0.1:%d/v1/transit/", cores[1].Listeners[0].Address.Port),
		fmt.Sprintf("https://127.0.0.1:%d/v1/transit/", cores[2].Listeners[0].Address.Port),
	}

	transport := &http.Transport{
		TLSClientConfig: cores[0].TLSConfig,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}

	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return fmt.Errorf("redirects not allowed in this test")
		},
	}

	//core.Logger().Printf("[TRACE] mounting transit")
	req, err := http.NewRequest("POST", fmt.Sprintf("https://127.0.0.1:%d/v1/sys/mounts/transit", cores[0].Listeners[0].Address.Port),
		bytes.NewBuffer([]byte("{\"type\": \"transit\"}")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(AuthHeaderName, cluster.RootToken)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	//core.Logger().Printf("[TRACE] done mounting transit")

	var totalOps uint64
	var successfulOps uint64
	var key1ver int64 = 1
	var key2ver int64 = 1
	var key3ver int64 = 1
	var numWorkers uint64 = 50
	var numWorkersStarted uint64
	var waitLock sync.Mutex
	waitCond := sync.NewCond(&waitLock)

	// This is the goroutine loop
	doFuzzy := func(id int, parallel bool) {
		var myTotalOps uint64
		var mySuccessfulOps uint64
		var keyVer int64 = 1
		// Check for panics, otherwise notify we're done
		defer func() {
			if err := recover(); err != nil {
				core.Logger().Error("got a panic: %v", err)
				t.Fail()
			}
			atomic.AddUint64(&totalOps, myTotalOps)
			atomic.AddUint64(&successfulOps, mySuccessfulOps)
			wg.Done()
		}()

		// Holds the latest encrypted value for each key
		latestEncryptedText := map[string]string{}

		client := &http.Client{
			Transport: transport,
		}

		var chosenFunc, chosenKey, chosenHost string

		myRand := rand.New(rand.NewSource(int64(id) * 400))

		doReq := func(method, url string, body io.Reader) (*http.Response, error) {
			req, err := http.NewRequest(method, url, body)
			if err != nil {
				return nil, err
			}
			req.Header.Set(AuthHeaderName, cluster.RootToken)
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
			err := result.Error()
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
				_, err := doReq("POST", chosenHost+"keys/"+fmt.Sprintf("%s-%t", chosenKey, parallel), bytes.NewBuffer([]byte("{}")))
				if err != nil {
					panic(err)
				}
			}
		}

		if !parallel {
			chosenHost = hosts[id%len(hosts)]
			chosenKey = fmt.Sprintf("key-%t-%d", parallel, id)

			_, err := doReq("POST", chosenHost+"keys/"+chosenKey, bytes.NewBuffer([]byte("{}")))
			if err != nil {
				panic(err)
			}
		}

		atomic.AddUint64(&numWorkersStarted, 1)

		waitCond.L.Lock()
		for atomic.LoadUint64(&numWorkersStarted) != numWorkers {
			waitCond.Wait()
		}
		waitCond.L.Unlock()
		waitCond.Broadcast()

		core.Logger().Trace("Starting goroutine", "id", id)

		startTime := time.Now()
		for {
			// Stop after 10 seconds
			if time.Now().Sub(startTime) > 10*time.Second {
				return
			}

			myTotalOps++

			// Pick a function and a key
			chosenFunc = funcs[myRand.Int()%len(funcs)]
			if parallel {
				chosenKey = fmt.Sprintf("%s-%t", keys[myRand.Int()%len(keys)], parallel)
				chosenHost = hosts[myRand.Int()%len(hosts)]
			}

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

				mySuccessfulOps++

			// Decrypt the ciphertext and compare the result
			case "decrypt":
				ct := latestEncryptedText[chosenKey]
				if ct == "" {
					mySuccessfulOps++
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
					if strings.Contains(err.Error(), keysutil.ErrTooOld) {
						mySuccessfulOps++
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

				mySuccessfulOps++

			// Rotate to a new key version
			case "rotate":
				//core.Logger().Printf("[TRACE] %s, %s, %d", chosenFunc, chosenKey, id)
				_, err := doReq("POST", chosenHost+"keys/"+chosenKey+"/rotate", bytes.NewBuffer([]byte("{}")))
				if err != nil {
					panic(err)
				}
				if parallel {
					switch chosenKey {
					case "test1":
						atomic.AddInt64(&key1ver, 1)
					case "test2":
						atomic.AddInt64(&key2ver, 1)
					case "test3":
						atomic.AddInt64(&key3ver, 1)
					}
				} else {
					keyVer++
				}

				mySuccessfulOps++

			// Change the min version, which also tests the archive functionality
			case "change_min_version":
				var latestVersion int64 = keyVer
				if parallel {
					switch chosenKey {
					case "test1":
						latestVersion = atomic.LoadInt64(&key1ver)
					case "test2":
						latestVersion = atomic.LoadInt64(&key2ver)
					case "test3":
						latestVersion = atomic.LoadInt64(&key3ver)
					}
				}

				setVersion := (myRand.Int63() % latestVersion) + 1

				//core.Logger().Printf("[TRACE] %s, %s, %d, new min version %d", chosenFunc, chosenKey, id, setVersion)

				_, err := doReq("POST", chosenHost+"keys/"+chosenKey+"/config", bytes.NewBuffer([]byte(fmt.Sprintf("{\"min_decryption_version\": %d}", setVersion))))
				if err != nil {
					panic(err)
				}

				mySuccessfulOps++
			}
		}
	}

	atomic.StoreUint64(&numWorkers, num)

	// Spawn some of these workers for 10 seconds
	for i := 0; i < int(atomic.LoadUint64(&numWorkers)); i++ {
		wg.Add(1)
		//core.Logger().Printf("[TRACE] spawning %d", i)
		go doFuzzy(i+1, parallel)
	}

	// Wait for them all to finish
	wg.Wait()

	if totalOps == 0 || totalOps != successfulOps {
		t.Fatalf("total/successful ops zero or mismatch: %d/%d; parallel: %t, num %d", totalOps, successfulOps, parallel, num)
	}
	t.Logf("total operations tried: %d, total successful: %d; parallel: %t, num %d", totalOps, successfulOps, parallel, num)
}

// This tests TLS connection state forwarding by ensuring that we can use a
// client TLS to authenticate against the cert backend
func TestHTTP_Forwarding_ClientTLS(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"cert": credCert.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(t, core)

	transport := cleanhttp.DefaultTransport()
	transport.TLSClientConfig = cores[0].TLSConfig
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}

	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://127.0.0.1:%d/v1/sys/auth/cert", cores[0].Listeners[0].Address.Port),
		bytes.NewBuffer([]byte("{\"type\": \"cert\"}")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(AuthHeaderName, cluster.RootToken)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	type certConfig struct {
		Certificate string `json:"certificate"`
		Policies    string `json:"policies"`
	}
	encodedCertConfig, err := json.Marshal(&certConfig{
		Certificate: string(cluster.CACertPEM),
		Policies:    "default",
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", fmt.Sprintf("https://127.0.0.1:%d/v1/auth/cert/certs/test", cores[0].Listeners[0].Address.Port),
		bytes.NewBuffer(encodedCertConfig))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(AuthHeaderName, cluster.RootToken)
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	addrs := []string{
		fmt.Sprintf("https://127.0.0.1:%d", cores[1].Listeners[0].Address.Port),
		fmt.Sprintf("https://127.0.0.1:%d", cores[2].Listeners[0].Address.Port),
	}

	// Ensure we can't possibly use lingering connections even though it should be to a different address

	transport = cleanhttp.DefaultTransport()
	transport.TLSClientConfig = cores[0].TLSConfig
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}

	client = &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return fmt.Errorf("redirects not allowed in this test")
		},
	}

	//cores[0].Logger().Printf("cluster.RootToken token is %s", cluster.RootToken)
	//time.Sleep(4 * time.Hour)

	for _, addr := range addrs {
		client := cores[0].Client
		client.SetAddress(addr)

		secret, err := client.Logical().Write("auth/cert/login", nil)
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			t.Fatal("secret is nil")
		}
		if secret.Auth == nil {
			t.Fatal("auth is nil")
		}
		if secret.Auth.Policies == nil || len(secret.Auth.Policies) == 0 || secret.Auth.Policies[0] != "default" {
			t.Fatalf("bad policies: %#v", secret.Auth.Policies)
		}
		if secret.Auth.ClientToken == "" {
			t.Fatalf("bad client token: %#v", *secret.Auth)
		}
		client.SetToken(secret.Auth.ClientToken)
		secret, err = client.Auth().Token().LookupSelf()
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil {
			t.Fatal("secret is nil")
		}
		if secret.Data == nil || len(secret.Data) == 0 {
			t.Fatal("secret data was empty")
		}
	}
}

func TestHTTP_Forwarding_HelpOperation(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	testHelp := func(client *api.Client) {
		help, err := client.Help("auth/token")
		if err != nil {
			t.Fatal(err)
		}
		if help == nil {
			t.Fatal("help was nil")
		}
	}

	testHelp(cores[0].Client)
	testHelp(cores[1].Client)
}
