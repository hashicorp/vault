package vault

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
)

func TestCluster(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	cluster, err := c.Cluster()
	if err != nil {
		t.Fatal(err)
	}
	// Test whether expected values are found
	if cluster == nil || cluster.Name == "" || cluster.ID == "" {
		t.Fatalf("cluster information missing: cluster: %#v", cluster)
	}

	// Test whether a private key has been generated
	entry, err := c.barrier.Get(coreLocalClusterKeyPath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("missing local cluster private key")
	}

	var params privKeyParams
	if err = jsonutil.DecodeJSON(entry.Value, &params); err != nil {
		t.Fatal(err)
	}
	switch {
	case params.X == nil, params.Y == nil, params.D == nil:
		t.Fatalf("x or y or d are nil: %#v", params)
	case params.Type == corePrivateKeyTypeP521:
	default:
		t.Fatal("parameter error: %#v", params)
	}
}

func TestClusterHA(t *testing.T) {
	logger = log.New(os.Stderr, "", log.LstdFlags)
	advertise := "http://127.0.0.1:8200"

	c, err := NewCore(&CoreConfig{
		Physical:      physical.NewInmemHA(logger),
		HAPhysical:    physical.NewInmemHA(logger),
		AdvertiseAddr: advertise,
		DisableMlock:  true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	key, _ := TestCoreInit(t, c)
	if _, err := c.Unseal(TestKeyCopy(key)); err != nil {
		t.Fatalf("unseal err: %s", err)
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
	testWaitActive(t, c)

	cluster, err := c.Cluster()
	if err != nil {
		t.Fatal(err)
	}
	// Test whether expected values are found
	if cluster == nil || cluster.Name == "" || cluster.ID == "" || cluster.Certificate == nil || len(cluster.Certificate) == 0 {
		t.Fatalf("cluster information missing: cluster:%#v", cluster)
	}

	// Test whether a private key has been generated
	entry, err := c.barrier.Get(coreLocalClusterKeyPath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("missing local cluster private key")
	}

	var params privKeyParams
	if err = jsonutil.DecodeJSON(entry.Value, &params); err != nil {
		t.Fatal(err)
	}
	switch {
	case params.X == nil, params.Y == nil, params.D == nil:
		t.Fatalf("x or y or d are nil: %#v", params)
	case params.Type == corePrivateKeyTypeP521:
	default:
		t.Fatal("parameter error: %#v", params)
	}

	// Make sure the certificate meets expectations
	cert, err := x509.ParseCertificate(cluster.Certificate)
	if err != nil {
		t.Fatal("error parsing local cluster certificate: %v", err)
	}
	if cert.Subject.CommonName != "127.0.0.1" {
		t.Fatalf("bad common name: %#v", cert.Subject.CommonName)
	}
	if len(cert.DNSNames) != 1 || cert.DNSNames[0] != "127.0.0.1" {
		t.Fatalf("bad dns names: %#v", cert.DNSNames)
	}
	if len(cert.IPAddresses) != 1 || cert.IPAddresses[0].String() != "127.0.0.1" {
		t.Fatalf("bad ip sans: %#v", cert.IPAddresses)
	}

	// Make sure the cert pool is as expected
	if len(c.localClusterCertPool.Subjects()) != 1 {
		t.Fatal("unexpected local cluster cert pool length")
	}
	if !reflect.DeepEqual(cert.RawSubject, c.localClusterCertPool.Subjects()[0]) {
		t.Fatal("cert pool subject does not match expected")
	}
}

func TestCluster_ListenForRequests(t *testing.T) {
	logger = log.New(os.Stderr, "", log.LstdFlags)
	advertise := "http://127.0.0.1:8200"

	c, err := NewCore(&CoreConfig{
		Physical:      physical.NewInmemHA(logger),
		HAPhysical:    physical.NewInmemHA(logger),
		AdvertiseAddr: advertise,
		DisableMlock:  true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:8202")
	if err != nil {
		t.Fatal(err)
	}
	lns := []net.Listener{ln}
	ln, err = net.Listen("tcp", "127.0.0.1:8204")
	if err != nil {
		t.Fatal(err)
	}
	lns = append(lns, ln)

	defer func() {
		for _, ln := range lns {
			ln.Close()
		}
	}()

	clusterListenerSetupFunc := func() ([]net.Listener, http.Handler, error) {
		ret := make([]net.Listener, 0, len(lns))
		// Loop over the existing listeners and start listeners on appropriate ports
		for _, ln := range lns {
			tcpAddr, ok := ln.Addr().(*net.TCPAddr)
			if !ok {
				c.logger.Printf("[TRACE] command/server: %s not a candidate for cluster request handling", ln.Addr().String())
				continue
			}
			c.logger.Printf("[TRACE] command/server: %s is a candidate for cluster request handling at addr %s and port %d", tcpAddr.String(), tcpAddr.IP.String(), tcpAddr.Port+1)

			ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", tcpAddr.IP.String(), tcpAddr.Port+1))
			if err != nil {
				return nil, nil, err
			}
			ret = append(ret, ln)
		}

		return ret, nil, nil
	}

	c.SetClusterListenerSetupFunc(clusterListenerSetupFunc)

	key, root := TestCoreInit(t, c)
	if _, err := c.Unseal(TestKeyCopy(key)); err != nil {
		t.Fatalf("unseal err: %s", err)
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
	testWaitActive(t, c)

	tlsConfig, err := c.ClusterTLSConfig()
	if err != nil {
		t.Fatal(err)
	}

	checkListenersFunc := func(expectFail bool) {
		for _, ln := range lns {
			tcpAddr, ok := ln.Addr().(*net.TCPAddr)
			if !ok {
				t.Fatal("%s not a TCP port", tcpAddr.String())
			}

			conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", tcpAddr.IP.String(), tcpAddr.Port+1), tlsConfig)
			if err != nil {
				if expectFail {
					t.Logf("testing %s:%d unsuccessful as expected", tcpAddr.IP.String(), tcpAddr.Port+1)
					continue
				}
				t.Fatal(err)
			}
			if expectFail {
				t.Fatalf("testing %s:%d not unsuccessful as expected", tcpAddr.IP.String(), tcpAddr.Port+1)
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
			t.Logf("testing %s:%d successful", tcpAddr.IP.String(), tcpAddr.Port+1)
		}
	}

	checkListenersFunc(false)

	// Make this nicer for tests
	oldManualStepDownSleepPeriod := manualStepDownSleepPeriod
	manualStepDownSleepPeriod = 3 * time.Second
	// Restore this value for other tests
	defer func() { manualStepDownSleepPeriod = oldManualStepDownSleepPeriod }()

	err = c.StepDown(&logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/seal",
		ClientToken: root,
	})
	if err != nil {
		t.Fatal(err)
	}

	// StepDown doesn't wait during actual preSeal so give time for listeners
	// to close
	time.Sleep(1 * time.Second)
	checkListenersFunc(true)

	time.Sleep(manualStepDownSleepPeriod)
	checkListenersFunc(false)

	err = c.Seal(root)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	checkListenersFunc(true)
}
