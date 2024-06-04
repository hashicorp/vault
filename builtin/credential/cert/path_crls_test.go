// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestCRLFetch(t *testing.T) {
	storage := &logical.InmemStorage{}

	lb, err := Factory(context.Background(), &logical.BackendConfig{
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 300 * time.Second,
			MaxLeaseTTLVal:     1800 * time.Second,
		},
		StorageView: storage,
	})

	require.NoError(t, err)
	b := lb.(*backend)
	closeChan := make(chan bool)
	go func() {
		t := time.NewTicker(50 * time.Millisecond)
		for {
			select {
			case <-t.C:
				b.PeriodicFunc(context.Background(), &logical.Request{Storage: storage})
			case <-closeChan:
				break
			}
		}
	}()
	defer close(closeChan)

	if err != nil {
		t.Fatalf("error: %s", err)
	}
	connState, err := testConnState("test-fixtures/keys/cert.pem",
		"test-fixtures/keys/key.pem", "test-fixtures/root/rootcacert.pem")
	require.NoError(t, err)
	caPEM, err := ioutil.ReadFile("test-fixtures/root/rootcacert.pem")
	require.NoError(t, err)
	caKeyPEM, err := ioutil.ReadFile("test-fixtures/keys/key.pem")
	require.NoError(t, err)
	certPEM, err := ioutil.ReadFile("test-fixtures/keys/cert.pem")

	caBundle, err := certutil.ParsePEMBundle(string(caPEM))
	require.NoError(t, err)
	bundle, err := certutil.ParsePEMBundle(string(certPEM) + "\n" + string(caKeyPEM))
	require.NoError(t, err)
	//  Entry with one cert first

	revocationListTemplate := &x509.RevocationList{
		RevokedCertificates: []pkix.RevokedCertificate{
			{
				SerialNumber:   big.NewInt(1),
				RevocationTime: time.Now(),
			},
		},
		Number:             big.NewInt(1),
		ThisUpdate:         time.Now(),
		NextUpdate:         time.Now().Add(50 * time.Millisecond),
		SignatureAlgorithm: x509.SHA1WithRSA,
	}

	var crlBytesLock sync.Mutex
	crlBytes, err := x509.CreateRevocationList(rand.Reader, revocationListTemplate, caBundle.Certificate, bundle.PrivateKey)
	require.NoError(t, err)

	var serverURL *url.URL
	crlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == serverURL.Host {
			crlBytesLock.Lock()
			w.Write(crlBytes)
			crlBytesLock.Unlock()
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	serverURL, _ = url.Parse(crlServer.URL)

	req := &logical.Request{
		Connection: &logical.Connection{
			ConnState: &connState,
		},
		Storage: storage,
		Auth:    &logical.Auth{},
	}

	fd := &framework.FieldData{
		Raw: map[string]interface{}{
			"name":        "test",
			"certificate": string(caPEM),
			"policies":    "foo,bar",
		},
		Schema: pathCerts(b).Fields,
	}

	resp, err := b.pathCertWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}

	empty_login_fd := &framework.FieldData{
		Raw:    map[string]interface{}{},
		Schema: pathLogin(b).Fields,
	}
	resp, err = b.pathLogin(context.Background(), req, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	// Set a bad CRL
	fd = &framework.FieldData{
		Raw: map[string]interface{}{
			"name": "testcrl",
			"url":  "http://wrongserver.com",
		},
		Schema: pathCRLs(b).Fields,
	}
	resp, err = b.pathCRLWrite(context.Background(), req, fd)
	if err == nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	// Set good CRL
	fd = &framework.FieldData{
		Raw: map[string]interface{}{
			"name": "testcrl",
			"url":  crlServer.URL,
		},
		Schema: pathCRLs(b).Fields,
	}
	resp, err = b.pathCRLWrite(context.Background(), req, fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}

	b.crlUpdateMutex.Lock()
	if len(b.crls["testcrl"].Serials) != 1 {
		t.Fatalf("wrong number of certs in CRL got %d, expected 1", len(b.crls["testcrl"].Serials))
	}
	b.crlUpdateMutex.Unlock()

	// Add a cert to the CRL, then wait to see if it gets automatically picked up
	revocationListTemplate.RevokedCertificates = []pkix.RevokedCertificate{
		{
			SerialNumber:   big.NewInt(1),
			RevocationTime: revocationListTemplate.RevokedCertificates[0].RevocationTime,
		},
		{
			SerialNumber:   big.NewInt(2),
			RevocationTime: time.Now(),
		},
	}
	revocationListTemplate.ThisUpdate = time.Now()
	revocationListTemplate.NextUpdate = time.Now().Add(1 * time.Minute)
	revocationListTemplate.Number = big.NewInt(2)

	crlBytesLock.Lock()
	crlBytes, err = x509.CreateRevocationList(rand.Reader, revocationListTemplate, caBundle.Certificate, bundle.PrivateKey)
	crlBytesLock.Unlock()
	require.NoError(t, err)

	// Give ourselves a little extra room on slower CI systems to ensure we
	// can fetch the new CRL.
	corehelpers.RetryUntil(t, 2*time.Second, func() error {
		b.crlUpdateMutex.Lock()
		defer b.crlUpdateMutex.Unlock()

		serialCount := len(b.crls["testcrl"].Serials)
		if serialCount != 2 {
			return fmt.Errorf("CRL refresh did not occur serial count %d", serialCount)
		}
		return nil
	})
}
