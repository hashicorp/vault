.DEFAULT_GOAL := help
SHELL := /bin/bash

TEST?=$$(go list ./... | grep -v /vendor/)
VETARGS?=-asmdecl \
			-atomic \
			-bool \
			-buildtags \
			-copylocks \
			-methods \
			-nilfunc \
			-printf \
			-rangeloops \
			-shift \
			-structtags \
			-unsafeptr

VAULT_LOCAL_PATH := $(CURDIR)
VAULT_CTR_MOUNT := /go/src/github.com/hashicorp/vault/

ZK_IP = $(shell docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' vault-test-zk)
CONSUL_IP = $(shell docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' vault-test-consul)

BACKEND_ENVS := -e ZOOKEEPER_ADDR=$(ZK_IP):2181 \
				-e CONSUL_ADDR=$(CONSUL_IP):8500

DEVKIT_COMMON_DOCKER_OPTS := --name vault-devkit \
	-p 8200:8200 \
	-v $(VAULT_LOCAL_PATH):$(VAULT_CTR_MOUNT)

.PHONY: clean-aux-containers
clean-aux-containers:
	-docker rm -vf vault-test-consul > /dev/null 2>&1
	-docker rm -vf vault-test-zk > /dev/null 2>&1

.PHONY: clean-devkit-container
clean-devkit-container:
	-docker rm -vf vault-devkit > /dev/null 2>&1

.PHONY: clean-containers
clean-containers: clean-devkit-container clean-aux-containers

.PHONY: clean
clean: clean-containers
	@echo "+ Cleaning up binaries..."
	-sudo rm -vf $(VAULT_LOCAL_PATH)/bin/*

.PHONY: devkit
devkit:
	@# The problem here is that we need to compile SQLSurvivor using
	@# paths/mounts that are unavailabe from Dockerfile. So ATM the best
	@# solution seems to be having one aditionall "build step" in Makefiles
	@# devkit target
	if $$(docker images | grep mesosphereci/vault-devkit | grep -q latest); then \
		echo "+ Devkit image already build"; \
	else \
		echo "+ Building devkit image"; \
		docker rmi -f mesosphereci/vault-devkit:latest; \
		docker build --rm --force-rm -t mesosphereci/vault-devkit:latest -f Dockerfile ./ ||\
		exit 1 ; \
	fi

.PHONY: update-devkit
update-devkit: clean-devkit-container
	docker build -t mesosphereci/vault-devkit:latest -f Dockerfile ./

# ZK super creds: 'super:secret'
.PHONY: aux
aux:
	$(eval ZK_CID := $(shell docker ps -a -q -f name=vault-test-zk))
	$(eval CONSUL_CID := $(shell docker ps -a -q -f name=vault-test-consul))
	if [[ -z "$(ZK_CID)" ]]; then \
		docker run -d \
			-p 2181:2181 -p 2888:2888 -p 3888:3888 \
			-e ZOOKEEPER_TICK_TIME=100 \
			-e JVMFLAGS=-Dzookeeper.DigestAuthenticationProvider.superDigest=super:lK75jTNcA+U9vtVEw5vB51mj/w4= \
			--name=vault-test-zk \
			digitalwonderland/zookeeper:latest; \
	else \
		docker start vault-test-zk; \
	fi
	if [[ -z "$(CONSUL_CID)" ]]; then \
		docker run -d \
			--name=vault-test-consul \
			-p 8400:8400 \
			-p 8500:8500 \
			-p 8600:53/udp \
			-h node1 \
			progrium/consul \
				-server -bootstrap -ui-dir /ui; \
		grep -m 1 "joined, marking health alive" <(docker logs -f `docker ps -a -q -f name=vault-test-consul`); \
	else \
		docker start vault-test-consul; \
	fi

.PHONY: shell
shell: clean-devkit-container devkit
	@# Run privileged in case user wants to debug stuff inside container
	docker run --rm -it \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		--privileged \
		mesosphereci/vault-devkit:latest /bin/bash

.PHONY: shell-aux
shell-aux: clean-devkit-container devkit aux
	@# Run privileged in case user wants to debug stuff inside container
	docker run --rm -it \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		$(BACKEND_ENVS) \
		--privileged \
		mesosphereci/vault-devkit:latest /bin/bash

.PHONY: run
run: clean-devkit-container devkit build
	docker run -d \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		mesosphereci/vault-devkit:latest \
	@echo "+ vault running in background, issue 'docker logs -f vault-devkit' for logs"

.PHONY: stop
stop: clean-devkit-container

.PHONY: help
help:
	@echo "Please see README.md file."

# build creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin
.PHONY: build
build:  clean-devkit-container devkit generate
	docker run --rm \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		mesosphereci/vault-devkit:latest \
			/bin/bash -x -c ' \
				VAULT_DEV_BUILD=1 sh -c "./scripts/build.sh"'

# test runs the unit tests and vets the code
.PHONY: test
test: generate testplain testrace

.PHONY: testplain
testplain: clean-containers devkit generate aux
	docker run --rm \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		$(BACKEND_ENVS) \
		mesosphereci/vault-devkit:latest \
			/bin/bash -x -c ' \
				VAULT_TOKEN= TF_ACC= \
				go test $(TEST) $(TESTARGS) -timeout=120s -parallel=4'

# testacc runs acceptance tests
.PHONY: testacc
testacc: clean-containers devkit generate aux
	docker run --rm \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		$(BACKEND_ENVS) \
		mesosphereci/vault-devkit:latest \
			/bin/bash -x -c ' \
				if [ "$(TEST)" = "./..." ]; then \
					echo "ERROR: Set TEST to a specific package"; \
					exit 1; \
				fi; \
				TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 45m'

# testrace runs the race checker
.PHONY: testrace
testrace: clean-containers devkit generate aux
	docker run --rm \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		$(BACKEND_ENVS) \
		mesosphereci/vault-devkit:latest \
			/bin/bash -x -c ' \
				CGO_ENABLED=1 VAULT_TOKEN= \
				TF_ACC= go test -race $(TEST) $(TESTARGS)'

.PHONY: cover
cover: clean-containers devkit aux
	docker run --rm \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		$(BACKEND_ENVS) \
		mesosphereci/vault-devkit:latest \
			/bin/bash -x -c ' \
				./scripts/coverage.sh --html'

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
.PHONY: vet
vet: clean-devkit-container devkit
	docker run --rm \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		mesosphereci/vault-devkit:latest \
			/bin/bash -x -c ' \
				go list -f '{{.Dir}}' ./... | grep -v /vendor/ \
					| grep -v '.*github.com/hashicorp/vault$$' \
					| xargs go tool vet ; if [ $$? -eq 1 ]; then \
						echo ""; \
						echo "Vet found suspicious constructs. Please check the reported constructs"; \
						echo "and fix them if necessary before submitting the code for reviewal."; \
					fi'

# generate runs `go generate` to build the dynamically generated
# source files.
.PHONY: generate
generate: clean-devkit-container devkit
	docker run --rm \
		$(DEVKIT_COMMON_DOCKER_OPTS) \
		mesosphereci/vault-devkit:latest \
			/bin/bash -x -c ' \
				go generate $(go list ./... | grep -v /vendor/)'
