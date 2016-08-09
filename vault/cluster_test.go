package vault

import (
	"crypto/x509"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/jsonutil"
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
