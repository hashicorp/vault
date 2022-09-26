package pki

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/docker"

	"github.com/hashicorp/go-uuid"

	"github.com/stretchr/testify/require"
)

var cwRunner *docker.Runner
var builtNetwork string
var buildClientContainerOnce sync.Once

const (
	protectedFile   = `dadgarcorp-internal-protected`
	unprotectedFile = `hello-world`
	uniqueHostname  = `dadgarcorpvaultpkitestingnginxwgetcurlcontainersexample.com`
)

func buildNginxContainer(t *testing.T, chain string, private string) (func(), string, int, string, string, int) {
	containerfile := `
FROM nginx:latest

RUN mkdir /www /etc/nginx/ssl && rm /etc/nginx/conf.d/*.conf

COPY testing.conf /etc/nginx/conf.d/
COPY fullchain.pem /etc/nginx/ssl/fullchain.pem
COPY privkey.pem /etc/nginx/ssl/privkey.pem
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

    location / {
		root /www/data;
    }
}
`

	bCtx := docker.NewBuildContext()
	bCtx["testing.conf"] = docker.PathContentsFromString(siteConfig)
	bCtx["fullchain.pem"] = docker.PathContentsFromString(chain)
	bCtx["privkey.pem"] = docker.PathContentsFromString(private)
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
		ContainerName: "vault_pki_nginx_integration",
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
		time.Sleep(2 * time.Second)
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

RUN apt update && DEBIAN_FRONTEND="noninteractive" apt install -y curl wget
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
		t.Fatalf("Could not start golang container for zlint: %s", err)
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

func RunNginxRootTest(t *testing.T, caKeyType string, caKeyBits int, caUsePSS bool, roleKeyType string, roleKeyBits int, roleUsePSS bool) {
	b, s := createBackendWithStorage(t)

	testSuffix := fmt.Sprintf(" - %v %v %v - %v %v %v", caKeyType, caKeyType, caUsePSS, roleKeyType, roleKeyBits, roleUsePSS)

	// Create a root and intermediate, setting the intermediate as default.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name":  "Root X1" + testSuffix,
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"key_bits":     caKeyBits,
		"use_pss":      caUsePSS,
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create root cert")
	rootCert := resp.Data["certificate"].(string)
	resp, err = CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
		"common_name":  "Intermediate I1" + testSuffix,
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"key_bits":     caKeyBits,
		"use_pss":      caUsePSS,
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create intermediate csr")
	resp, err = CBWrite(b, s, "issuer/default/sign-intermediate", map[string]interface{}{
		"common_name":  "Intermediate I1",
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     caKeyType,
		"csr":          resp.Data["csr"],
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to sign intermediate csr")
	resp, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": resp.Data["certificate"],
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to sign intermediate csr")
	_, err = CBWrite(b, s, "config/issuers", map[string]interface{}{
		"default": resp.Data["imported_issuers"].([]string)[0],
	})

	// Create a role+certificate valid for localhost only.
	_, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allow_any_name": true,
		"key_type":       roleKeyType,
		"key_bits":       roleKeyBits,
		"use_pss":        roleUsePSS,
		"ttl":            "60m",
	})
	require.NoError(t, err)
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": uniqueHostname,
		"ip_sans":     "127.0.0.1,::1",
		"sans":        uniqueHostname + ",localhost,localhost4,localhost6,localhost.localdomain",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to create leaf cert")
	leafCert := resp.Data["certificate"].(string)
	leafPrivateKey := resp.Data["private_key"].(string) + "\n"
	fullChain := leafCert + "\n"
	for _, cert := range resp.Data["ca_chain"].([]string) {
		fullChain += cert + "\n"
	}

	cleanup, host, port, networkName, networkAddr, networkPort := buildNginxContainer(t, fullChain, leafPrivateKey)
	defer cleanup()

	localURL := "https://" + host + ":" + strconv.Itoa(port) + "/index.html"
	containerURL := "https://" + uniqueHostname + ":" + strconv.Itoa(networkPort) + "/index.html"

	// Ensure we can connect with Go.
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(rootCert))
	tlsConfig := &tls.Config{
		RootCAs: pool,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	clientResp, err := client.Get(localURL)
	if err != nil {
		t.Fatalf("failed to fetch url (%v): %v", localURL, err)
	}
	defer clientResp.Body.Close()
	body, err := io.ReadAll(clientResp.Body)
	if err != nil {
		t.Fatalf("failed to get read response body: %v", err)
	}
	if !strings.Contains(string(body), unprotectedFile) {
		t.Fatalf("expected body to contain (%v) but was:\n%v", unprotectedFile, string(body))
	}

	// Ensure we can connect with wget/curl.
	CheckWithClients(t, networkName, networkAddr, containerURL, rootCert, "", "")
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
