package vault

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
	log "github.com/mgutz/logxi/v1"
)

var (
	clusterTestPausePeriod = 2 * time.Second
)

func TestClusterFetching(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	err := c.setupCluster()
	if err != nil {
		t.Fatal(err)
	}

	cluster, err := c.Cluster()
	if err != nil {
		t.Fatal(err)
	}
	// Test whether expected values are found
	if cluster == nil || cluster.Name == "" || cluster.ID == "" {
		t.Fatalf("cluster information missing: cluster: %#v", cluster)
	}
}

func TestClusterHAFetching(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	redirect := "http://127.0.0.1:8200"

	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirect,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, _ := TestCoreInit(t, c)
	for _, key := range keys {
		if _, err := TestCoreUnseal(c, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	sealed, err := c.Sealed()
	if err != nil {
		t.Fatalf("err checking seal status: %s", err)
	}
	if sealed {
		t.Fatal("should not be sealed")
	}

	// Wait for core to become active
	TestWaitActive(t, c)

	cluster, err := c.Cluster()
	if err != nil {
		t.Fatal(err)
	}
	// Test whether expected values are found
	if cluster == nil || cluster.Name == "" || cluster.ID == "" {
		t.Fatalf("cluster information missing: cluster:%#v", cluster)
	}
}

func TestCluster_ListenForRequests(t *testing.T) {
	// Make this nicer for tests
	manualStepDownSleepPeriod = 5 * time.Second

	cluster := NewTestCluster(t, nil, &TestClusterOptions{
		KeepStandbysSealed: true,
	})
	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores

	// Wait for core to become active
	TestWaitActive(t, cores[0].Core)

	// Use this to have a valid config after sealing since ClusterTLSConfig returns nil
	var lastTLSConfig *tls.Config
	checkListenersFunc := func(expectFail bool) {
		tlsConfig, err := cores[0].ClusterTLSConfig()
		if err != nil {
			if err.Error() != consts.ErrSealed.Error() {
				t.Fatal(err)
			}
			tlsConfig = lastTLSConfig
		} else {
			tlsConfig.NextProtos = []string{"h2"}
			lastTLSConfig = tlsConfig
		}

		for _, ln := range cores[0].Listeners {
			tcpAddr, ok := ln.Addr().(*net.TCPAddr)
			if !ok {
				t.Fatalf("%s not a TCP port", tcpAddr.String())
			}

			conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", tcpAddr.IP.String(), tcpAddr.Port+105), tlsConfig)
			if err != nil {
				if expectFail {
					t.Logf("testing %s:%d unsuccessful as expected", tcpAddr.IP.String(), tcpAddr.Port+105)
					continue
				}
				t.Fatalf("error: %v\nlisteners are\n%#v\n%#v\n", err, cores[0].Listeners[0], cores[0].Listeners[1])
			}
			if expectFail {
				t.Fatalf("testing %s:%d not unsuccessful as expected", tcpAddr.IP.String(), tcpAddr.Port+105)
			}
			err = conn.Handshake()
			if err != nil {
				t.Fatal(err)
			}
			connState := conn.ConnectionState()
			switch {
			case connState.Version != tls.VersionTLS12:
				t.Fatal("version mismatch")
			case connState.NegotiatedProtocol != "h2" || !connState.NegotiatedProtocolIsMutual:
				t.Fatal("bad protocol negotiation")
			}
			t.Logf("testing %s:%d successful", tcpAddr.IP.String(), tcpAddr.Port+105)
		}
	}

	time.Sleep(clusterTestPausePeriod)
	checkListenersFunc(false)

	err := cores[0].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: cluster.RootToken,
	})
	if err != nil {
		t.Fatal(err)
	}

	// StepDown doesn't wait during actual preSeal so give time for listeners
	// to close
	time.Sleep(clusterTestPausePeriod)
	checkListenersFunc(true)

	// After this period it should be active again
	time.Sleep(manualStepDownSleepPeriod)
	checkListenersFunc(false)

	err = cores[0].Seal(cluster.RootToken)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	// After sealing it should be inactive again
	checkListenersFunc(true)
}

func TestCluster_ForwardRequests(t *testing.T) {
	// Make this nicer for tests
	manualStepDownSleepPeriod = 5 * time.Second

	testCluster_ForwardRequestsCommon(t)
}

func testCluster_ForwardRequestsCommon(t *testing.T) {
	cluster := NewTestCluster(t, nil, nil)
	cores := cluster.Cores
	cores[0].Handler.(*http.ServeMux).HandleFunc("/core1", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte("core1"))
	})
	cores[1].Handler.(*http.ServeMux).HandleFunc("/core2", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(202)
		w.Write([]byte("core2"))
	})
	cores[2].Handler.(*http.ServeMux).HandleFunc("/core3", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(203)
		w.Write([]byte("core3"))
	})
	cluster.Start()
	defer cluster.Cleanup()

	root := cluster.RootToken

	// Wait for core to become active
	TestWaitActive(t, cores[0].Core)

	// Test forwarding a request. Since we're going directly from core to core
	// with no fallback we know that if it worked, request handling is working
	testCluster_ForwardRequests(t, cores[1], root, "core1")
	testCluster_ForwardRequests(t, cores[2], root, "core1")

	//
	// Now we do a bunch of round-robining. The point is to make sure that as
	// nodes come and go, we can always successfully forward to the active
	// node.
	//

	// Ensure active core is cores[1] and test
	err := cores[0].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[2].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[1].Core)
	testCluster_ForwardRequests(t, cores[0], root, "core2")
	testCluster_ForwardRequests(t, cores[2], root, "core2")

	// Ensure active core is cores[2] and test
	err = cores[1].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[0].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[2].Core)
	testCluster_ForwardRequests(t, cores[0], root, "core3")
	testCluster_ForwardRequests(t, cores[1], root, "core3")

	// Ensure active core is cores[0] and test
	err = cores[2].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[1].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[0].Core)
	testCluster_ForwardRequests(t, cores[1], root, "core1")
	testCluster_ForwardRequests(t, cores[2], root, "core1")

	// Ensure active core is cores[1] and test
	err = cores[0].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[2].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[1].Core)
	testCluster_ForwardRequests(t, cores[0], root, "core2")
	testCluster_ForwardRequests(t, cores[2], root, "core2")

	// Ensure active core is cores[2] and test
	err = cores[1].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[0].StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[2].Core)
	testCluster_ForwardRequests(t, cores[0], root, "core3")
	testCluster_ForwardRequests(t, cores[1], root, "core3")
}

func testCluster_ForwardRequests(t *testing.T, c *TestClusterCore, rootToken, remoteCoreID string) {
	standby, err := c.Standby()
	if err != nil {
		t.Fatal(err)
	}
	if !standby {
		t.Fatal("expected core to be standby")
	}

	// We need to call Leader as that refreshes the connection info
	isLeader, _, _, err := c.Leader()
	if err != nil {
		panic(err.Error())
		t.Fatal(err)
	}
	if isLeader {
		t.Fatal("core should not be leader")
	}

	bodBuf := bytes.NewReader([]byte(`{ "foo": "bar", "zip": "zap" }`))
	req, err := http.NewRequest("PUT", "https://pushit.real.good:9281/"+remoteCoreID, bodBuf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Vault-Token", rootToken)

	statusCode, header, respBytes, err := c.ForwardRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if header == nil {
		t.Fatal("err: expected at least a content-type header")
	}
	if header.Get("Content-Type") != "application/json" {
		t.Fatalf("bad content-type: %s", header.Get("Content-Type"))
	}

	body := string(respBytes)

	if body != remoteCoreID {
		t.Fatalf("expected %s, got %s", remoteCoreID, body)
	}
	switch body {
	case "core1":
		if statusCode != 201 {
			t.Fatal("bad response")
		}
	case "core2":
		if statusCode != 202 {
			t.Fatal("bad response")
		}
	case "core3":
		if statusCode != 203 {
			t.Fatal("bad response")
		}
	}
}

func TestCluster_CustomCipherSuites(t *testing.T) {
	cluster := NewTestCluster(t, &CoreConfig{
		ClusterCipherSuites: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	}, nil)
	cluster.Start()
	defer cluster.Cleanup()
	core := cluster.Cores[0]

	// Wait for core to become active
	TestWaitActive(t, core.Core)

	tlsConf, err := core.Core.ClusterTLSConfig()
	if err != nil {
		t.Fatal(err)
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", core.Listeners[0].Address.IP.String(), core.Listeners[0].Address.Port+105), tlsConf)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	err = conn.Handshake()
	if err != nil {
		t.Fatal(err)
	}
	if conn.ConnectionState().CipherSuite != tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 {
		var availCiphers string
		for _, cipher := range core.clusterCipherSuites {
			availCiphers += fmt.Sprintf("%x ", cipher)
		}
		t.Fatalf("got bad negotiated cipher %x, core-set suites are %s", conn.ConnectionState().CipherSuite, availCiphers)
	}
}
