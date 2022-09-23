package pki

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/helper/testhelpers/docker"
)

func RunZLintContainer(t *testing.T, certificate string) []byte {
	containerfile := `
FROM golang:latest

RUN go install github.com/zmap/zlint/v3/cmd/zlint@latest

COPY cert.pem /go/cert.pem
`

	bCtx := docker.NewBuildContext()
	bCtx["cert.pem"] = docker.PathContentsFromBytes([]byte(certificate))

	imageName := "vault_pki_zlint_validator"
	imageTag := "latest"

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     imageName,
		ImageTag:      imageTag,
		ContainerName: "pki_zlint",
		// We want to run sleep in the background so we're not stuck waiting
		// for the default golang container's shell to prompt for input.
		Entrypoint: []string{"sleep", "1200"},
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

	// We don't actually care about the address, we just want to start the
	// container so we can run commands in it. We'd ideally like to skip this
	// step and only build a new image, but the zlint output would be
	// intermingled with container build stages, so its not that useful.
	ctr, _, _, err := runner.Start(ctx)
	if err != nil {
		t.Fatalf("Could not start golang container for zlint: %s", err)
	}

	// Run the zlint command and save the output.
	cmd := []string{"/go/bin/zlint", "/go/cert.pem"}
	stdout, stderr, retcode, err := runner.RunCmdWithOutput(ctx, ctr.ID, cmd)
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

	return stdout
}

func RunZLintRootTest(t *testing.T, keyType string, keyBits int, usePSS bool, ignored []string) {
	b, s := createBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
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

func Test_ZLint(t *testing.T) {
	RunZLintRootTest(t, "rsa", 2048, false, nil)
	RunZLintRootTest(t, "rsa", 2048, true, nil)
	RunZLintRootTest(t, "rsa", 3072, false, nil)
	RunZLintRootTest(t, "rsa", 3072, true, nil)
	RunZLintRootTest(t, "ec", 256, false, nil)
	RunZLintRootTest(t, "ec", 384, false, nil)

	// Mozilla doesn't allow P-521 ECDSA keys.
	RunZLintRootTest(t, "ec", 521, false, []string{
		"e_mp_ecdsa_pub_key_encoding_correct",
		"e_mp_ecdsa_signature_encoding_correct"})
}
