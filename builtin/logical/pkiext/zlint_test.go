// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pkiext

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

var (
	zRunner        *docker.Runner
	buildZLintOnce sync.Once
	releaseRegex   = regexp.MustCompile(`^go\d+\.\d+\.\d+$`)
)

func buildZLintContainer(t *testing.T) {
	// Leverage the Go version running the test to pull a version tagged image
	// to avoid the issues we sometimes encounter pulling images with the latest tag
	runtimeVer := runtime.Version()
	goVersion := "latest"
	// The version returned from Go might not be a release tag such as go1.23.2, if it
	// isn't fallback to latest
	if releaseRegex.MatchString(runtimeVer) {
		goVersion = strings.TrimPrefix(runtime.Version(), "go")
	}
	containerfile := fmt.Sprintf(`
FROM docker.mirror.hashicorp.services/library/golang:%s
RUN go install github.com/zmap/zlint/v3/cmd/zlint@v3.6.2
`, goVersion)

	bCtx := docker.NewBuildContext()

	imageName := "vault_pki_zlint_validator"
	imageTag := "latest"

	var err error
	zRunner, err = docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     imageName,
		ImageTag:      imageTag,
		ContainerName: "pki_zlint",
		// We want to run sleep in the background so we're not stuck waiting
		// for the default golang container's shell to prompt for input.
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

	// Sometimes we see timeouts and issues pulling the zlint code from GitHub
	testhelpers.RetryUntil(t, 30*time.Second, func() error {
		output, err := zRunner.BuildImage(ctx, containerfile, bCtx,
			docker.BuildRemove(true),
			docker.BuildForceRemove(true),
			docker.BuildPullParent(true),
			docker.BuildTags([]string{imageName + ":" + imageTag}))
		if err != nil {
			return fmt.Errorf("could not build new image with zlint: %w", err)
		}

		t.Logf("Image build output: %v", string(output))
		return nil
	})
}

func RunZLintContainer(t *testing.T, certificate string) []byte {
	buildZLintOnce.Do(func() {
		buildZLintContainer(t)
	})

	ctx := context.Background()
	// We don't actually care about the address, we just want to start the
	// container so we can run commands in it. We'd ideally like to skip this
	// step and only build a new image, but the zlint output would be
	// intermingled with container build stages, so its not that useful.
	result, err := zRunner.Start(ctx, true, false)
	if err != nil {
		t.Fatalf("Could not start golang container for zlint: %s", err)
	}

	// Copy the cert into the newly running container.
	certCtx := docker.NewBuildContext()
	certCtx["cert.pem"] = docker.PathContentsFromBytes([]byte(certificate))
	if err := zRunner.CopyTo(result.Container.ID, "/go/", certCtx); err != nil {
		t.Fatalf("Could not copy certificate into container: %v", err)
	}

	// Run the zlint command and save the output.
	cmd := []string{"/go/bin/zlint", "/go/cert.pem"}
	stdout, stderr, retcode, err := zRunner.RunCmdWithOutput(ctx, result.Container.ID, cmd)
	if err != nil {
		t.Fatalf("Could not run command in container: %v", err)
	}

	if len(stderr) != 0 {
		t.Logf("Got stderr from command:\n%v\n", string(stderr))
	}

	if retcode != 0 {
		t.Logf("Got stdout from command:\n%v\n", string(stdout))
		t.Fatalf("Got unexpected non-zero retcode from zlint: %v\n", retcode)
	}

	// Clean up after ourselves.
	if err := zRunner.Stop(context.Background(), result.Container.ID); err != nil {
		t.Fatalf("failed to stop container: %v", err)
	}

	return stdout
}

func RunZLintRootTest(t *testing.T, keyType string, keyBits int, usePSS bool, ignored []string) {
	b, s := pki.CreateBackendWithStorage(t)

	resp, err := pki.CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name":  "Root X1",
		"country":      "US",
		"organization": "Dadgarcorp",
		"ou":           "QA",
		"key_type":     keyType,
		"key_bits":     keyBits,
		"use_pss":      usePSS,
	})
	require.NoError(t, err)
	rootCert := resp.Data["certificate"].(string)

	var parsed map[string]interface{}
	output := RunZLintContainer(t, rootCert)

	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("failed to parse zlint output as JSON: %v\nOutput:\n%v\n\n", err, string(output))
	}

	for key, rawValue := range parsed {
		value := rawValue.(map[string]interface{})
		result, ok := value["result"]
		if !ok || result == "NA" {
			continue
		}

		if result == "error" {
			skip := false
			for _, allowedFailures := range ignored {
				if allowedFailures == key {
					skip = true
					break
				}
			}

			if !skip {
				t.Fatalf("got unexpected error from test %v: %v", key, value)
			}
		}
	}
}

func Test_ZLintRSA2048(t *testing.T) {
	t.Parallel()
	RunZLintRootTest(t, "rsa", 2048, false, nil)
}

func Test_ZLintRSA2048PSS(t *testing.T) {
	t.Parallel()
	RunZLintRootTest(t, "rsa", 2048, true, nil)
}

func Test_ZLintRSA3072(t *testing.T) {
	t.Parallel()
	RunZLintRootTest(t, "rsa", 3072, false, nil)
}

func Test_ZLintRSA3072PSS(t *testing.T) {
	t.Parallel()
	RunZLintRootTest(t, "rsa", 3072, true, nil)
}

func Test_ZLintECDSA256(t *testing.T) {
	t.Parallel()
	RunZLintRootTest(t, "ec", 256, false, nil)
}

func Test_ZLintECDSA384(t *testing.T) {
	t.Parallel()
	RunZLintRootTest(t, "ec", 384, false, nil)
}

func Test_ZLintECDSA521(t *testing.T) {
	t.Parallel()
	// Mozilla doesn't allow P-521 ECDSA keys.
	RunZLintRootTest(t, "ec", 521, false, []string{
		"e_mp_ecdsa_pub_key_encoding_correct",
		"e_mp_ecdsa_signature_encoding_correct",
	})
}
