package vault

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
)

var (
	clusterTestPausePeriod = 2 * time.Second
)

func TestClusterFetching(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	err := c.setupCluster(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	cluster, err := c.Cluster(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// Test whether expected values are found
	if cluster == nil || cluster.Name == "" || cluster.ID == "" {
		t.Fatalf("cluster information missing: cluster: %#v", cluster)
	}
}

func TestClusterHAFetching(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)

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
	if c.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Wait for core to become active
	TestWaitActive(t, c)

	cluster, err := c.Cluster(context.Background())
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
	checkListenersFunc := func(expectFail bool) {
		cores[0].clusterListener.AddClient(requestForwardingALPN, &requestForwardingClusterClient{cores[0].Core})

		parsedCert := cores[0].localClusterParsedCert.Load().(*x509.Certificate)
		dialer := cores[0].getGRPCDialer(context.Background(), requestForwardingALPN, parsedCert.Subject.CommonName, parsedCert)
		for i := range cores[0].Listeners {

			clnAddr := cores[0].clusterListener.clusterListenerAddrs[i]
			netConn, err := dialer(clnAddr.String(), 0)
			conn := netConn.(*tls.Conn)
			if err != nil {
				if expectFail {
					t.Logf("testing %s unsuccessful as expected", clnAddr)
					continue
				}
				t.Fatalf("error: %v\ncluster listener is %s", err, clnAddr)
			}
			if expectFail {
				t.Fatalf("testing %s not unsuccessful as expected", clnAddr)
			}
			err = conn.Handshake()
			if err != nil {
				t.Fatal(err)
			}
			connState := conn.ConnectionState()
			switch {
			case connState.Version != tls.VersionTLS12:
				t.Fatal("version mismatch")
			case connState.NegotiatedProtocol != requestForwardingALPN || !connState.NegotiatedProtocolIsMutual:
				t.Fatal("bad protocol negotiation")
			}
			t.Logf("testing %s successful", clnAddr)
		}
	}

	time.Sleep(clusterTestPausePeriod)
	checkListenersFunc(false)

	err := cores[0].StepDown(context.Background(), &logical.Request{
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

	err = cores[0].Core.Seal(cluster.RootToken)
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
	err := cores[0].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[2].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[1].Core)
	testCluster_ForwardRequests(t, cores[0], root, "core2")
	testCluster_ForwardRequests(t, cores[2], root, "core2")

	// Ensure active core is cores[2] and test
	err = cores[1].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[0].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[2].Core)
	testCluster_ForwardRequests(t, cores[0], root, "core3")
	testCluster_ForwardRequests(t, cores[1], root, "core3")

	// Ensure active core is cores[0] and test
	err = cores[2].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[1].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[0].Core)
	testCluster_ForwardRequests(t, cores[1], root, "core1")
	testCluster_ForwardRequests(t, cores[2], root, "core1")

	// Ensure active core is cores[1] and test
	err = cores[0].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[2].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	time.Sleep(clusterTestPausePeriod)
	TestWaitActive(t, cores[1].Core)
	testCluster_ForwardRequests(t, cores[0], root, "core2")
	testCluster_ForwardRequests(t, cores[2], root, "core2")

	// Ensure active core is cores[2] and test
	err = cores[1].StepDown(context.Background(), &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/step-down",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(clusterTestPausePeriod)
	_ = cores[0].StepDown(context.Background(), &logical.Request{
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
	req.Header.Add(consts.AuthHeaderName, rootToken)
	req = req.WithContext(context.WithValue(req.Context(), "original_request_path", req.URL.Path))

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

	core.clusterListener.AddClient(requestForwardingALPN, &requestForwardingClusterClient{core.Core})

	parsedCert := core.localClusterParsedCert.Load().(*x509.Certificate)
	dialer := core.getGRPCDialer(context.Background(), requestForwardingALPN, parsedCert.Subject.CommonName, parsedCert)

	netConn, err := dialer(core.clusterListener.clusterListenerAddrs[0].String(), 0)
	conn := netConn.(*tls.Conn)
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
