package pkiext

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/testhelpers/docker"

	"github.com/hashicorp/go-uuid"

	"github.com/stretchr/testify/require"
)

var (
	cwRunner                 *docker.Runner
	builtNetwork             string
	buildClientContainerOnce sync.Once
)

const (
	protectedFile    = `dadgarcorp-internal-protected`
	unprotectedFile  = `hello-world`
	failureIndicator = `THIS-TEST-SHOULD-FAIL`
	uniqueHostname   = `dadgarcorpvaultpkitestingnginxwgetcurlcontainersexample.com`
	containerName    = `vault_pki_nginx_integration`
)

func buildNginxContainer(t *testing.T, root string, crl string, chain string, private string) (func(), string, int, string, string, int) {
	containerfile := `
FROM nginx:latest

RUN mkdir /www /etc/nginx/ssl && rm /etc/nginx/conf.d/*.conf

COPY testing.conf /etc/nginx/conf.d/
COPY root.pem /etc/nginx/ssl/root.pem
COPY fullchain.pem /etc/nginx/ssl/fullchain.pem
COPY privkey.pem /etc/nginx/ssl/privkey.pem
COPY crl.pem /etc/nginx/ssl/crl.pem
COPY /data /www/data
`

	siteConfig := `
server {
	listen 80;
	listen [::]:80;

    location / {
        return 301 $request_uri;
    }
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;

	ssl_certificate /etc/nginx/ssl/fullchain.pem;
	ssl_certificate_key /etc/nginx/ssl/privkey.pem;

	ssl_client_certificate /etc/nginx/ssl/root.pem;
	ssl_crl /etc/nginx/ssl/crl.pem;
	ssl_verify_client optional;

	# Magic per: https://serverfault.com/questions/891603/nginx-reverse-proxy-with-optional-ssl-client-authentication
	# Only necessary since we're too lazy to setup two different subdomains.
	set $ssl_status 'open';
	if ($request_uri ~ protected) {
		set $ssl_status 'closed';
	}

	if ($ssl_client_verify != SUCCESS) {
		set $ssl_status "$ssl_status-fail";
	}

	if ($ssl_status = "closed-fail") {
		return 403;
	}

	location / {
		root /www/data;
    }
}
`

	bCtx := docker.NewBuildContext()
	bCtx["testing.conf"] = docker.PathContentsFromString(siteConfig)
	bCtx["root.pem"] = docker.PathContentsFromString(root)
	bCtx["fullchain.pem"] = docker.PathContentsFromString(chain)
	bCtx["privkey.pem"] = docker.PathContentsFromString(private)
	bCtx["crl.pem"] = docker.PathContentsFromString(crl)
	bCtx["/data/index.html"] = docker.PathContentsFromString(unprotectedFile)
	bCtx["/data/protected.html"] = docker.PathContentsFromString(protectedFile)

	imageName := "vault_pki_nginx_integration"
	suffix, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("error generating unique suffix: %v", err)
	}
	imageTag := suffix

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     imageName,
		ImageTag:      imageTag,
		ContainerName: containerName,
		Ports:         []string{"443/tcp"},
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	})
	if err != nil {
		t.Fatalf("Could not provision docker service runner: %s", err)
	}

	ctx := context.Background()
	output, err := runner.BuildImage(ctx, containerfile, bCtx,
		docker.BuildRemove(true), docker.BuildForceRemove(true),
		docker.BuildPullParent(true),
		docker.BuildTags([]string{imageName + ":" + imageTag}))
	if err != nil {
		t.Fatalf("Could not build new image: %v", err)
	}

	t.Logf("Image build output: %v", string(output))

	svc, err := runner.StartService(ctx, func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		// Nginx loads fast, we're too lazy to validate this properly.
		time.Sleep(5 * time.Second)
		return docker.NewServiceHostPort(host, port), nil
	})
	if err != nil {
		t.Fatalf("Could not start nginx container: %v", err)
	}

	// We also need to find the network address of this node, and return
	// the non-local address associated with it so that we can spawn the
	// client command on the correct network/port.
	networks, err := runner.GetNetworkAndAddresses(svc.Container.ID)
	if err != nil {
		t.Fatalf("Could not interrogate container for addresses: %v", err)
	}

	var networkName string
	var networkAddr string
	for name, addr := range networks {
		if addr == "" {
			continue
		}

		networkName = name
		networkAddr = addr
		break
	}

	if networkName == "" || networkAddr == "" {
		t.Fatalf("failed to get network info for containers: empty network address: %v", networks)
	}

	pieces := strings.Split(svc.Config.Address(), ":")
	port, _ := strconv.Atoi(pieces[1])
	return svc.Cleanup, pieces[0], port, networkName, networkAddr, 443
}

func buildWgetCurlContainer(t *testing.T, network string) {
	containerfile := `
FROM ubuntu:latest

RUN apt update && DEBIAN_FRONTEND="noninteractive" apt install -y curl wget wget2
`

	bCtx := docker.NewBuildContext()

	imageName := "vault_pki_wget_curl_integration"
	imageTag := "latest"

	var err error
	cwRunner, err = docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     imageName,
		ImageTag:      imageTag,
		ContainerName: "vault_pki_wget_curl",
		NetworkID:     network,
		// We want to run sleep in the background so we're not stuck waiting
		// for the default ubuntu container's shell to prompt for input.
		Entrypoint: []string{"sleep", "45"},
		LogConsumer: func(s string) {
			if t.Failed() {
				t.Logf("container logs: %s", s)
			}
		},
	})
	if err != nil {
		t.Fatalf("Could not provision docker service runner: %s", err)
	}

	ctx := context.Background()
	output, err := cwRunner.BuildImage(ctx, containerfile, bCtx,
		docker.BuildRemove(true), docker.BuildForceRemove(true),
		docker.BuildPullParent(true),
		docker.BuildTags([]string{imageName + ":" + imageTag}))
	if err != nil {
		t.Fatalf("Could not build new image: %v", err)
	}

	t.Logf("Image build output: %v", string(output))
}

func CheckWithClients(t *testing.T, network string, address string, url string, rootCert string, certificate string, privatekey string) {
	// We assume the network doesn't change once assigned.
	buildClientContainerOnce.Do(func() {
		buildWgetCurlContainer(t, network)
		builtNetwork = network
	})

	if builtNetwork != network {
		t.Fatalf("failed assumption check: different built network (%v) vs run network (%v); must've changed while running tests", builtNetwork, network)
	}

	// Start our service with a random name to not conflict with other
	// threads.
	ctx := context.Background()
	ctr, _, _, err := cwRunner.Start(ctx, true, false)
	if err != nil {
		t.Fatalf("Could not start golang container for wget/curl checks: %s", err)
	}

	// Commands to run after potentially writing the certificate. We
	// might augment these if the certificate exists.
	//
	// We manually add the expected hostname to the local hosts file
	// to avoid resolving it over the network and instead resolving it
	// to this other container we just started (potentially in parallel
	// with other containers).
	hostPrimeCmd := []string{"sh", "-c", "echo '" + address + "	" + uniqueHostname + "' >> /etc/hosts"}
	wgetCmd := []string{"wget", "--verbose", "--ca-certificate=/root.pem", url}
	curlCmd := []string{"curl", "--verbose", "--cacert", "/root.pem", url}

	certCtx := docker.NewBuildContext()
	certCtx["root.pem"] = docker.PathContentsFromString(rootCert)
	if certificate != "" {
		// Copy the cert into the newly running container.
		certCtx["client-cert.pem"] = docker.PathContentsFromString(certificate)
		certCtx["client-privkey.pem"] = docker.PathContentsFromString(privatekey)

		wgetCmd = []string{"wget", "--verbose", "--ca-certificate=/root.pem", "--certificate=/client-cert.pem", "--private-key=/client-privkey.pem", url}
		curlCmd = []string{"curl", "--verbose", "--cacert", "/root.pem", "--cert", "/client-cert.pem", "--key", "/client-privkey.pem", url}
	}
	if err := cwRunner.CopyTo(ctr.ID, "/", certCtx); err != nil {
		t.Fatalf("Could not copy certificate and key into container: %v", err)
	}

	for _, cmd := range [][]string{hostPrimeCmd, wgetCmd, curlCmd} {
		t.Logf("Running client connection command: %v", cmd)

		stdout, stderr, retcode, err := cwRunner.RunCmdWithOutput(ctx, ctr.ID, cmd)
		if err != nil {
			t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
		}

		if len(stderr) != 0 {
			t.Logf("Got stderr from command (%v):\n%v\n", cmd, string(stderr))
		}

		if retcode != 0 {
			t.Logf("Got stdout from command (%v):\n%v\n", cmd, string(stdout))
			t.Fatalf("Got unexpected non-zero retcode from command (%v): %v\n", cmd, retcode)
		}
	}
}

func CheckDeltaCRL(t *testing.T, network string, address string, url string, rootCert string, crls string) {
	// We assume the network doesn't change once assigned.
	buildClientContainerOnce.Do(func() {
		buildWgetCurlContainer(t, network)
		builtNetwork = network
	})

	if builtNetwork != network {
		t.Fatalf("failed assumption check: different built network (%v) vs run network (%v); must've changed while running tests", builtNetwork, network)
	}

	// Start our service with a random name to not conflict with other
	// threads.
	ctx := context.Background()
	ctr, _, _, err := cwRunner.Start(ctx, true, false)
	if err != nil {
		t.Fatalf("Could not start golang container for wget2 delta CRL checks: %s", err)
	}

	// Commands to run after potentially writing the certificate. We
	// might augment these if the certificate exists.
	//
	// We manually add the expected hostname to the local hosts file
	// to avoid resolving it over the network and instead resolving it
	// to this other container we just started (potentially in parallel
	// with other containers).
	hostPrimeCmd := []string{"sh", "-c", "echo '" + address + "	" + uniqueHostname + "' >> /etc/hosts"}
	wgetCmd := []string{"wget2", "--verbose", "--ca-certificate=/root.pem", "--crl-file=/crls.pem", url}

	certCtx := docker.NewBuildContext()
	certCtx["root.pem"] = docker.PathContentsFromString(rootCert)
	certCtx["crls.pem"] = docker.PathContentsFromString(crls)
	if err := cwRunner.CopyTo(ctr.ID, "/", certCtx); err != nil {
		t.Fatalf("Could not copy certificate and key into container: %v", err)
	}

	for index, cmd := range [][]string{hostPrimeCmd, wgetCmd} {
		t.Logf("Running client connection command: %v", cmd)

		stdout, stderr, retcode, err := cwRunner.RunCmdWithOutput(ctx, ctr.ID, cmd)
		if err != nil {
			t.Fatalf("Could not run command (%v) in container: %v", cmd, err)
		}

		if len(stderr) != 0 {
			t.Logf("Got stderr from command (%v):\n%v\n", cmd, string(stderr))
		}

		if retcode != 0 && index == 0 {
			t.Logf("Got stdout from command (%v):\n%v\n", cmd, string(stdout))
			t.Fatalf("Got unexpected non-zero retcode from command (%v): %v\n", cmd, retcode)
		}

		if retcode == 0 && index == 1 {
			t.Logf("Got stdout from command (%v):\n%v\n", cmd, string(stdout))
			t.Fatalf("Got unexpected zero retcode from command; wanted this to fail (%v): %v\n", cmd, retcode)
		}
	}
}

func CheckWithGo(t *testing.T, rootCert string, clientCert string, clientChain []string, clientKey string, host string, port int, networkAddr string, networkPort int, url string, expected string, shouldFail bool) {
	// Ensure we can connect with Go.
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(rootCert))
	tlsConfig := &tls.Config{
		RootCAs: pool,
	}

	if clientCert != "" {
		var clientTLSCert tls.Certificate
		clientTLSCert.Certificate = append(clientTLSCert.Certificate, parseCert(t, clientCert).Raw)
		clientTLSCert.PrivateKey = parseKey(t, clientKey)
		for _, cert := range clientChain {
			clientTLSCert.Certificate = append(clientTLSCert.Certificate, parseCert(t, cert).Raw)
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, clientTLSCert)
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == host+":"+strconv.Itoa(port) {
				// If we can't resolve our hostname, try
				// accessing it via the docker protocol
				// instead of via the returned service
				// address.
				if _, err := net.LookupHost(host); err != nil && strings.Contains(err.Error(), "no such host") {
					addr = networkAddr + ":" + strconv.Itoa(networkPort)
				}
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}

	client := &http.Client{Transport: transport}
	clientResp, err := client.Get(url)
	if err != nil {
		if shouldFail {
			return
		}

		t.Fatalf("failed to fetch url (%v): %v", url, err)
	} else if shouldFail {
		if clientResp.StatusCode == 200 {
			t.Fatalf("expected failure to fetch url (%v): got response: %v", url, clientResp)
		}

		return
	}

	defer clientResp.Body.Close()
	body, err := io.ReadAll(clientResp.Body)
	if err != nil {
		t.Fatalf("failed to get read response body: %v", err)
	}
	if !strings.Contains(string(body), expected) {
		t.Fatalf("expected body to contain (%v) but was:\n%v", expected, string(body))
	}
}

func RunNginxRootTest(t *testing.T, caKeyType string, caKeyBits int, caUsePSS bool, roleKeyType string, roleKeyBits int, roleUsePSS bool) {
	t.Skipf("flaky in CI")

	b, s := pki.CreateBackendWithStorage(t)

	testSuffix := fmt.Sprintf(" - %v %v %v - %v %v %v", caKeyType, caKeyType, caUsePSS, roleKeyType, roleKeyBits, roleUsePSS)

	// Configure our mount to use auto-rotate, even though we don't have
	// a periodic func.
	_, err := pki.CBWrite(b, s, "config/crl", map[string]interface{}{
		"auto_rebuild": true,
		"enable_delta": true,
	})

	// Create a root and intermediate, setting the intermediate as default.
	resp, err := pki.CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name":  "Root X1" + testSuffix,
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"key_bits":     caKeyBits,
		"use_pss":      caUsePSS,
		"issuer_name":  "root",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create root cert")
	rootCert := resp.Data["certificate"].(string)
	resp, err = pki.CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
		"common_name":  "Intermediate I1" + testSuffix,
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"key_bits":     caKeyBits,
		"use_pss":      caUsePSS,
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create intermediate csr")
	resp, err = pki.CBWrite(b, s, "issuer/default/sign-intermediate", map[string]interface{}{
		"common_name":  "Intermediate I1",
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"csr":          resp.Data["csr"],
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to sign intermediate csr")
	intCert := resp.Data["certificate"].(string)
	resp, err = pki.CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": intCert,
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to sign intermediate csr")
	_, err = pki.CBWrite(b, s, "config/issuers", map[string]interface{}{
		"default": resp.Data["imported_issuers"].([]string)[0],
	})

	// Create a role+certificate valid for localhost only.
	_, err = pki.CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allow_any_name": true,
		"key_type":       roleKeyType,
		"key_bits":       roleKeyBits,
		"use_pss":        roleUsePSS,
		"ttl":            "60m",
	})
	require.NoError(t, err)
	resp, err = pki.CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": uniqueHostname,
		"ip_sans":     "127.0.0.1,::1",
		"sans":        uniqueHostname + ",localhost,localhost4,localhost6,localhost.localdomain",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create server leaf cert")
	leafCert := resp.Data["certificate"].(string)
	leafPrivateKey := resp.Data["private_key"].(string) + "\n"
	fullChain := leafCert + "\n"
	for _, cert := range resp.Data["ca_chain"].([]string) {
		fullChain += cert + "\n"
	}

	// Issue a client leaf certificate.
	resp, err = pki.CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "testing.client.dadgarcorp.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create client leaf cert")
	clientCert := resp.Data["certificate"].(string)
	clientKey := resp.Data["private_key"].(string) + "\n"
	clientWireChain := clientCert + "\n" + resp.Data["issuing_ca"].(string) + "\n"
	clientTrustChain := resp.Data["issuing_ca"].(string) + "\n" + rootCert + "\n"
	clientCAChain := resp.Data["ca_chain"].([]string)

	// Issue a client leaf cert and revoke it, placing it on the main CRL
	// via rotation.
	resp, err = pki.CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "revoked-crl.client.dadgarcorp.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create revoked client leaf cert")
	revokedCert := resp.Data["certificate"].(string)
	revokedKey := resp.Data["private_key"].(string) + "\n"
	// revokedFullChain := revokedCert + "\n" + resp.Data["issuing_ca"].(string) + "\n"
	// revokedTrustChain := resp.Data["issuing_ca"].(string) + "\n" + rootCert + "\n"
	revokedCAChain := resp.Data["ca_chain"].([]string)
	_, err = pki.CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": revokedCert,
	})
	require.NoError(t, err)
	_, err = pki.CBRead(b, s, "crl/rotate")
	require.NoError(t, err)

	// Issue a client leaf cert and revoke it, placing it on the delta CRL
	// via rotation.
	/*resp, err = pki.CBWrite(b, s, "issue/testing", map[string]interface{}{
	      "common_name": "revoked-delta-crl.client.dadgarcorp.com",
	  })
	  requireSuccessNonNilResponse(t, resp, err, "failed to create delta CRL revoked client leaf cert")
	  deltaCert := resp.Data["certificate"].(string)
	  deltaKey := resp.Data["private_key"].(string) + "\n"
	  //deltaFullChain := deltaCert + "\n" + resp.Data["issuing_ca"].(string) + "\n"
	  //deltaTrustChain := resp.Data["issuing_ca"].(string) + "\n" + rootCert + "\n"
	  deltaCAChain := resp.Data["ca_chain"].([]string)
	  _, err = pki.CBWrite(b, s, "revoke", map[string]interface{}{
	      "certificate": deltaCert,
	  })
	  require.NoError(t, err)
	  _, err = pki.CBRead(b, s, "crl/rotate-delta")
	  require.NoError(t, err)*/

	// Get the CRL and Delta CRLs.
	resp, err = pki.CBRead(b, s, "issuer/root/crl")
	require.NoError(t, err)
	rootCRL := resp.Data["crl"].(string) + "\n"
	resp, err = pki.CBRead(b, s, "issuer/default/crl")
	require.NoError(t, err)
	intCRL := resp.Data["crl"].(string) + "\n"

	// No need to fetch root Delta CRL as we've not revoked anything on it.
	resp, err = pki.CBRead(b, s, "issuer/default/crl/delta")
	require.NoError(t, err)
	deltaCRL := resp.Data["crl"].(string) + "\n"

	crls := rootCRL + intCRL + deltaCRL

	cleanup, host, port, networkName, networkAddr, networkPort := buildNginxContainer(t, rootCert, crls, fullChain, leafPrivateKey)
	defer cleanup()

	if host != "127.0.0.1" && host != "::1" && strings.HasPrefix(host, containerName) {
		t.Logf("Assuming %v:%v is a container name rather than localhost reference.", host, port)
		host = uniqueHostname
		port = networkPort
	}

	localBase := "https://" + host + ":" + strconv.Itoa(port)
	localURL := localBase + "/index.html"
	localProtectedURL := localBase + "/protected.html"
	containerBase := "https://" + uniqueHostname + ":" + strconv.Itoa(networkPort)
	containerURL := containerBase + "/index.html"
	containerProtectedURL := containerBase + "/protected.html"

	t.Logf("Spawned nginx container:\nhost: %v\nport: %v\nnetworkName: %v\nnetworkAddr: %v\nnetworkPort: %v\nlocalURL: %v\ncontainerURL: %v\n", host, port, networkName, networkAddr, networkPort, localBase, containerBase)

	// Ensure we can connect with Go. We do our checks for revocation here,
	// as this behavior is server-controlled and shouldn't matter based on
	// client type.
	CheckWithGo(t, rootCert, "", nil, "", host, port, networkAddr, networkPort, localURL, unprotectedFile, false)
	CheckWithGo(t, rootCert, "", nil, "", host, port, networkAddr, networkPort, localProtectedURL, failureIndicator, true)
	CheckWithGo(t, rootCert, clientCert, clientCAChain, clientKey, host, port, networkAddr, networkPort, localProtectedURL, protectedFile, false)
	CheckWithGo(t, rootCert, revokedCert, revokedCAChain, revokedKey, host, port, networkAddr, networkPort, localProtectedURL, protectedFile, true)
	// CheckWithGo(t, rootCert, deltaCert, deltaCAChain, deltaKey, host, port, networkAddr, networkPort, localProtectedURL, protectedFile, true)

	// Ensure we can connect with wget/curl.
	CheckWithClients(t, networkName, networkAddr, containerURL, rootCert, "", "")
	CheckWithClients(t, networkName, networkAddr, containerProtectedURL, clientTrustChain, clientWireChain, clientKey)

	// Ensure OpenSSL will validate the delta CRL by revoking our server leaf
	// and then using it with wget2. This will land on the intermediate's
	// Delta CRL.
	_, err = pki.CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": leafCert,
	})
	require.NoError(t, err)
	_, err = pki.CBRead(b, s, "crl/rotate-delta")
	require.NoError(t, err)
	resp, err = pki.CBRead(b, s, "issuer/default/crl/delta")
	require.NoError(t, err)
	deltaCRL = resp.Data["crl"].(string) + "\n"
	crls = rootCRL + intCRL + deltaCRL

	CheckDeltaCRL(t, networkName, networkAddr, containerURL, rootCert, crls)
}

func Test_NginxRSAPure(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "rsa", 2048, false, "rsa", 2048, false)
}

func Test_NginxRSAPurePSS(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "rsa", 2048, false, "rsa", 2048, true)
}

func Test_NginxRSAPSSPure(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "rsa", 2048, true, "rsa", 2048, false)
}

func Test_NginxRSAPSSPurePSS(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "rsa", 2048, true, "rsa", 2048, true)
}

func Test_NginxECDSA256Pure(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "ec", 256, false, "ec", 256, false)
}

func Test_NginxECDSAHybrid(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "ec", 256, false, "rsa", 2048, false)
}

func Test_NginxECDSAHybridPSS(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "ec", 256, false, "rsa", 2048, true)
}

func Test_NginxRSAHybrid(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "rsa", 2048, false, "ec", 256, false)
}

func Test_NginxRSAPSSHybrid(t *testing.T) {
	t.Parallel()
	RunNginxRootTest(t, "rsa", 2048, true, "ec", 256, false)
}
