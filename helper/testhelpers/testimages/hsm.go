// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package testimages

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/constants"
	dockhelper "github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/stretchr/testify/require"
)

// GetImageRepoAndTag returns an image repo and tag that can be used to start a vault
// node via docker.  Env vars are used as inputs: either VAULT_BINARY and VAULT_IMAGE
// if hsm is false, or VAULT_HSM_BINARY and VAULT_HSM_IMAGE if hsm is true.
//
// If a matching image var is set, we split that on ":" and return the two pieces
// as the repo and tag.  If instead a matching binary var is set, we create an image
// using a vault-enterprise docker image as a starting point, then add softhsm
// (when hsm is true) and the specified binary to it.  If neither the image or binary var
// are set, we fail the test.
//
// For devs on their workstations, they can either create an image or a binary and
// set the env vars appropriately.  Creating an hsm linux binary is more challenging and
// time-consuming than creating a regular binary, so we don't want to impose that
// on people running tests that don't require one.
//
// See also tools/testimagemaker for a way to build an image for this purpose
// from the CLI.
func GetImageRepoAndTag(t *testing.T, hsm bool) (string, string) {
	t.Helper()
	repo, tag, output, err := CreateOrReturnDockerImage(hsm)
	if err != nil && output != nil {
		t.Logf("docker image create output: %s", output)
	}

	require.NoError(t, err)

	// t.Logf("used bin=%s (%q) and img=%s (%q) to create %s:%s", bin, os.Getenv(bin), img, os.Getenv(img), repo, tag)
	t.Cleanup(func() {
		// When image build fails, it doesn't always return an error, but typically the error
		// is visible in the output
		if t.Failed() && output != nil {
			t.Logf("docker image create output: %s", output)
		}
	})

	return repo, tag
}

// CreateOrReturnDockerImage looks at the vaultBinary and vaultImage params.
// If vaultImage is populated, it is split by ":" and the two pieces are returned
// as the repo and tag.  If vault_binary is populated, an image is created based on
// the latest hsm image.
// (TODO: currently hardcoded as "docker.mirror.hashicorp.services/vault-enterprise:2.0.4-ent.hsm")
// This is done by installing SoftHSM and the vaultBinary on top of that image.
// If neither is populated an error is returned.
func CreateOrReturnDockerImage(hsm bool) (repo string, tag string, output []byte, err error) {
	binVar, imgVar := "VAULT_BINARY", "VAULT_IMAGE"
	if hsm {
		binVar, imgVar = "VAULT_HSM_BINARY", "VAULT_HSM_IMAGE"
	}
	bin, img := os.Getenv(binVar), os.Getenv(imgVar)
	switch {
	case bin == "" && img == "":
		return "", "", nil, fmt.Errorf("no docker image or binary provided")
	case img != "":
		// Ignore the binary if an image is specified
		pieces := strings.Split(img, ":")
		if len(pieces) != 2 {
			return "", "", nil, fmt.Errorf("bad input image format %q", img)
		}
		return pieces[0], pieces[1], nil, nil
	default:
		base := "hashicorp/vault"
		if constants.IsEnterprise {
			base += "-enterprise"
		}
		repo := base + "-ci"
		tag := "latest"
		source := "docker.mirror.hashicorp.services/" + base + ":latest"
		if hsm {
			source = "docker.mirror.hashicorp.services/hashicorp/vault-enterprise:2.0.4-ent.hsm"
			tag = "latest-hsm"
		}
		target := fmt.Sprintf("%s:%s", repo, tag)
		var output []byte
		var err error
		if hsm {
			output, err = CreateHSMDockerImage(source, target, bin)
		} else {
			output, err = CreateNonHSMDockerImage(source, target, bin)
		}
		return repo, tag, output, err
	}
}

func createBuildContextWithBinary(vaultBinary string) (dockhelper.BuildContext, error) {
	f, err := os.Open(vaultBinary)
	if err != nil {
		return nil, fmt.Errorf("error opening vault binary file: %w", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading vault binary file: %w", err)
	}

	bCtx := dockhelper.NewBuildContext()
	bCtx["vault"] = &dockhelper.FileContents{
		Data: data,
		Mode: 0o755,
	}

	return bCtx, nil
}

// createDockerImage creates an image named toImage from the given context and Dockerfile.
// Extra BuildOpt values (e.g. BuildArgs) may be passed to customise the build.
func createDockerImage(toImage, containerFile string, bCtx dockhelper.BuildContext, extraOpts ...dockhelper.BuildOpt) ([]byte, error) {
	client, err := dockhelper.NewDockerAPI()
	if err != nil {
		return nil, err
	}

	buildOpts := []dockhelper.BuildOpt{
		dockhelper.BuildRemove(true),
		dockhelper.BuildForceRemove(true),
		dockhelper.BuildPullParent(true),
		dockhelper.BuildTags([]string{toImage}),
	}
	buildOpts = append(buildOpts, extraOpts...)

	output, err := dockhelper.BuildImage(context.Background(), client, containerFile, bCtx, buildOpts...)
	if err != nil {
		return nil, fmt.Errorf("error building docker image: %w (output: %s)", err, output)
	}

	return output, nil
}

func CreateNonHSMDockerImage(fromImage, toImage, vaultBinary string) ([]byte, error) {
	bCtx := dockhelper.NewBuildContext()
	var err error
	bCtx, err = createBuildContextWithBinary(vaultBinary)
	if err != nil {
		return nil, err
	}

	containerFile := fmt.Sprintf(`
FROM %s
USER root

COPY vault /bin/vault

USER vault
CMD ["server", "-dev"]
`, fromImage)
	return createDockerImage(toImage, containerFile, bCtx)
}

// CreateHSMDockerImage creates a new vault-enterprise HSM docker image from an
// existing HSM base image. SoftHSM 2.7.0 is built from source using only
// build deps from the UBI 10 repos (no external mirrors required at runtime).
//
// SoftHSM 2.6.1 — the newest version packaged by any RHEL 10-family distro —
// crashes with a C++ bounds-check assertion when AES-GCM is called without
// additional authenticated data (https://github.com/softhsm/SoftHSMv2/pull/664,
// fixed in 2.7.0).  Until a distro ships 2.7.0 we must build from source.
//
// If GH_TOKEN or GITHUB_TOKEN is set in the environment the token is forwarded
// to the build as a build-arg so that the GitHub release download is
// authenticated and not subject to anonymous rate-limiting.
func CreateHSMDockerImage(fromImage, toImage, vaultBinary string) ([]byte, error) {
	bCtx := dockhelper.NewBuildContext()
	if vaultBinary != "" {
		var err error
		bCtx, err = createBuildContextWithBinary(vaultBinary)
		if err != nil {
			return nil, err
		}
	}
	bCtx["setup-softhsm.sh"] = &dockhelper.FileContents{
		Data: []byte(`#!/bin/bash

mkdir -p /vault/file/softhsm/tokens

# only create a new slot if there isn't an existing one
if [ ! -e /vault/file/hsm-slot ]; then
	softhsm2-util --init-token --slot 0 --so-pin=12345 --pin=12345 --label "vault" | grep -oE '[0-9]+$' > /vault/file/hsm-slot
fi

exec docker-entrypoint.sh "$@"
`),
		Mode: 0o755,
	}

	// softhsm270SHA256 is the SHA-256 of the SoftHSMv2 2.7.0 source archive from
	// https://github.com/softhsm/SoftHSMv2/archive/refs/tags/2.7.0.tar.gz
	softhsm270SHA256 := "be14a5820ec457eac5154462ffae51ba5d8a643f6760514d4b4b83a77be91573"

	// All build deps (gcc, openssl-devel, etc.) are available in the UBI 10
	// baseos/appstream/CRB repos that are already enabled in the base image.
	containerFile := fmt.Sprintf(`FROM %s AS builder

USER root

RUN microdnf install -y wget make gcc gcc-c++ openssl-devel automake autoconf libtool pkg-config

# GITHUB_TOKEN is an optional build-arg; when set it is forwarded as a Bearer
# token so the download is authenticated and not subject to anonymous rate-limiting.
ARG GITHUB_TOKEN
RUN if [ -n "${GITHUB_TOKEN}" ]; then \
        wget -q --header="Authorization: Bearer ${GITHUB_TOKEN}" \
            https://github.com/softhsm/SoftHSMv2/archive/refs/tags/2.7.0.tar.gz; \
    else \
        wget -q \
            https://github.com/softhsm/SoftHSMv2/archive/refs/tags/2.7.0.tar.gz; \
    fi && \
    echo "%s  2.7.0.tar.gz" | sha256sum -c && \
    tar -xzf 2.7.0.tar.gz

RUN cd SoftHSMv2-2.7.0 && sh autogen.sh && ./configure --disable-gost && make -j$(nproc)

FROM %s

USER root

COPY --from=builder /SoftHSMv2-2.7.0/src/lib/.libs/libsofthsm2.so /usr/lib64/libsofthsm2.so
COPY --from=builder /SoftHSMv2-2.7.0/src/bin/util/softhsm2-util /usr/bin/softhsm2-util
COPY --from=builder /SoftHSMv2-2.7.0/src/lib/common/softhsm2.conf /etc/softhsm2.conf

RUN mkdir -p /usr/local/lib/softhsm && \
    ln -sf /usr/lib64/libsofthsm2.so /usr/local/lib/softhsm/libsofthsm2.so

# Put the tokens under /vault/file since that's the data volume, and if we want
# to start a cluster using a pre-existing volume (i.e. resuming from a previous
# cluster) we need the tokens in order to unseal/unsealwrap.
RUN sed -i 's|directories.tokendir = .*|directories.tokendir = /vault/file/softhsm/tokens|g' /etc/softhsm2.conf

RUN sed -i 's/log.level = ERROR/log.level = DEBUG/' /etc/softhsm2.conf

COPY setup-softhsm.sh /usr/local/bin/setup-softhsm.sh

COPY vault /bin/vault

USER vault
CMD ["server", "-dev"]
ENTRYPOINT ["setup-softhsm.sh"]
`, fromImage, softhsm270SHA256, fromImage)

	ghToken := os.Getenv("GH_TOKEN")
	if ghToken == "" {
		ghToken = os.Getenv("GITHUB_TOKEN")
	}

	// Forward a GitHub token as a build-arg if one is available in the environment.
	var extraOpts []dockhelper.BuildOpt
	if ghToken != "" {
		extraOpts = append(extraOpts, dockhelper.BuildArgs{"GITHUB_TOKEN": &ghToken})
	}

	return createDockerImage(toImage, containerFile, bCtx, extraOpts...)
}

const (
	PKCS11Library    = "/usr/lib64/libsofthsm2.so"
	PKCS11Pin        = "12345"
	PKCS11TokenLabel = "vault"
)
