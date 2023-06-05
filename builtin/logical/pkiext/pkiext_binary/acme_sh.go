package pkiext_binary

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/docker"
)

const ACMEshContainerfile = `
FROM docker.mirror.hashicorp.services/ubuntu:latest

RUN apt update && DEBIAN_FRONTEND="noninteractive" apt install -y openssl socat curl coreutils dnsutils tzdata sed tar jq libidn2-0 openssh-client git cron

RUN git clone https://github.com/acmesh-official/acme.sh.git /acme.sh
WORKDIR /acme.sh

RUN ./acme.sh --install --home /etc/acmesh --config-home /etc/ssl/data --cert-home /etc/ssl/certs --accountemail "webmaster@dadgarcorp.com"
`

var (
	acmeShRunner             *docker.Runner
	buildAcmeShContainerOnce sync.Once
)

func buildACMEshContainer(t *testing.T, network string) {
	bCtx := docker.NewBuildContext()

	imageName := "vault_pki_acme_acme_sh_integration"
	imageTag := "latest"
	containerName := "vault_acme_sh_container"

	var err error
	acmeShRunner, err = docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     imageName,
		ImageTag:      imageTag,
		ContainerName: containerName,
		NetworkName:   network,
		// We want to run sleep in the background so we're not stuck waiting
		// for the default ubuntu container's shell to prompt for input.
		//
		// We choose a slightly longer default sleep, 600s = 10min, as we
		// potentially want to allow the test to fail validating, which
		// might take a while.
		Entrypoint: []string{"sleep", "600"},
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
	output, err := acmeShRunner.BuildImage(ctx, ACMEshContainerfile, bCtx,
		docker.BuildRemove(true), docker.BuildForceRemove(true),
		docker.BuildPullParent(true),
		docker.BuildTags([]string{imageName + ":" + imageTag}))
	if err != nil {
		t.Fatalf("Could not build new image: %v", err)
	}

	t.Logf("Image build output: %v", string(output))
}

func GetACMEshContainer(t *testing.T, network string) (*docker.Runner, *docker.StartResult) {
	buildAcmeShContainerOnce.Do(func() {
		buildACMEshContainer(t, network)
	})

	ctx := context.Background()
	result, err := acmeShRunner.Start(ctx, true, false)
	if err != nil {
		t.Fatalf("Could not start acme.sh container: %s", err)
	}

	return acmeShRunner, result
}
