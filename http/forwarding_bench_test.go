package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/forwarding"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"golang.org/x/net/http2"
)

func BenchmarkHTTP_Forwarding_Stress(b *testing.B) {
	testPlaintextB64 := "dGhlIHF1aWNrIGJyb3duIGZveA=="

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	}

	cluster := vault.NewTestCluster(b, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
		Logger:      logging.NewVaultLoggerWithWriter(ioutil.Discard, log.Error),
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	// make it easy to get access to the active
	core := cores[0].Core
	vault.TestWaitActive(b, core)

	handler := cores[0].Handler
	host := fmt.Sprintf("https://127.0.0.1:%d/v1/transit/", cores[0].Listeners[0].Address.Port)

	transport := &http.Transport{
		TLSClientConfig: cores[0].TLSConfig,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		b.Fatal(err)
	}

	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://127.0.0.1:%d/v1/sys/mounts/transit", cores[0].Listeners[0].Address.Port),
		bytes.NewBuffer([]byte("{\"type\": \"transit\"}")))
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set(consts.AuthHeaderName, cluster.RootToken)
	_, err = client.Do(req)
	if err != nil {
		b.Fatal(err)
	}

	var numOps uint32

	doReq := func(b *testing.B, method, url string, body io.Reader) {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			b.Fatal(err)
		}
		req.Header.Set(consts.AuthHeaderName, cluster.RootToken)
		w := forwarding.NewRPCResponseWriter()
		handler.ServeHTTP(w, req)
		switch w.StatusCode() {
		case 200:
		case 204:
			if !strings.Contains(url, "keys") {
				b.Fatal("got 204")
			}
		default:
			b.Fatalf("bad status code: %d, resp: %s", w.StatusCode(), w.Body().String())
		}
		//b.Log(w.Body().String())
		numOps++
	}

	doReq(b, "POST", host+"keys/test1", bytes.NewBuffer([]byte("{}")))
	keyUrl := host + "encrypt/test1"
	reqBuf := []byte(fmt.Sprintf("{\"plaintext\": \"%s\"}", testPlaintextB64))

	b.Run("doreq", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			doReq(b, "POST", keyUrl, bytes.NewReader(reqBuf))
		}
	})

	b.Logf("total ops: %d", numOps)
}
