# kind cluster name
KIND_CLUSTER_NAME?=vault-plugin-secrets-kubernetes

# kind k8s version
KIND_K8S_VERSION?=v1.26.2

PKG=github.com/hashicorp/vault-plugin-secrets-kubernetes
LDFLAGS?="-X '$(PKG).WALRollbackMinAge=10s'"

RUNNER_TEMP ?= $(TMPDIR)

.PHONY: default
default: dev

# dev target sets WALRollbackMinAge to 10s instead of the default 10 minutes to speed up integration tests
.PHONY: dev
dev:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o bin/vault-plugin-secrets-kubernetes cmd/vault-plugin-secrets-kubernetes/main.go

.PHONY: test
test: fmtcheck
	CGO_ENABLED=0 go test ./... $(TESTARGS) -timeout=20m

.PHONY: integration-test
integration-test:
	INTEGRATION_TESTS=true KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME) CGO_ENABLED=0 go test github.com/hashicorp/vault-plugin-secrets-kubernetes/integrationtest/... $(TESTARGS) -count=1 -timeout=40m

.PHONY: fmtcheck
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: setup-kind
# create a kind cluster for running the acceptance tests locally
setup-kind:
	kind get clusters | grep --silent "^${KIND_CLUSTER_NAME}$$" || \
	kind create cluster \
		--image kindest/node:${KIND_K8S_VERSION} \
		--name ${KIND_CLUSTER_NAME}  \
		--config $(CURDIR)/integrationtest/kind/config.yaml
	kubectl config use-context kind-${KIND_CLUSTER_NAME}

.PHONY: delete-kind
# delete the kind cluster
delete-kind:
	kind delete cluster --name ${KIND_CLUSTER_NAME} || true

.PHONY: vault-image
vault-image:
	GOOS=linux make dev
	docker build -f integrationtest/vault/Dockerfile bin/ --tag=hashicorp/vault:dev

.PHONY: vault-image-ent
vault-image-ent:
	GOOS=linux make dev
	docker build -f integrationtest/vault/Dockerfile --target enterprise bin/ --tag=hashicorp/vault:dev

# Create Vault inside the cluster with a locally-built version of kubernetes secrets.
.PHONY: setup-integration-test-common
setup-integration-test-common: SET_LICENSE=$(if $(VAULT_LICENSE_CI),--set server.enterpriseLicense.secretName=vault-license)
setup-integration-test-common: teardown-integration-test
	kind --name ${KIND_CLUSTER_NAME} load docker-image hashicorp/vault:dev
	kubectl create namespace test
	kubectl label namespaces test target=integration-test other=label

	# don't log the license
	printenv VAULT_LICENSE_CI > $(RUNNER_TEMP)/vault-license.txt || true
	if [ -s $(RUNNER_TEMP)/vault-license.txt ]; then \
		kubectl -n test create secret generic vault-license --from-file license=$(RUNNER_TEMP)/vault-license.txt; \
		rm -rf $(RUNNER_TEMP)/vault-license.txt; \
	fi

	helm install vault vault --repo https://helm.releases.hashicorp.com --version=0.24.1 \
		--wait --timeout=5m \
		--namespace=test \
		--set server.logLevel=debug \
		--set server.dev.enabled=true \
		--set server.image.tag=dev \
		--set server.image.pullPolicy=Never \
		--set injector.enabled=false \
		$(SET_LICENSE) \
		--set server.extraArgs="-dev-plugin-dir=/vault/plugin_directory"
	kubectl patch --namespace=test statefulset vault --patch-file integrationtest/vault/hostPortPatch.yaml
	kubectl apply --namespace=test -f integrationtest/vault/testRoles.yaml
	kubectl apply --namespace=test -f integrationtest/vault/testServiceAccounts.yaml
	kubectl apply --namespace=test -f integrationtest/vault/testBindings.yaml

	kubectl delete --namespace=test pod vault-0
	kubectl wait --namespace=test --for=condition=Ready --timeout=5m pod -l app.kubernetes.io/name=vault

.PHONY: setup-integration-test
setup-integration-test: vault-image setup-integration-test-common

.PHONY: setup-integration-test-ent
setup-integration-test-ent: check-license vault-image-ent setup-integration-test-common

.PHONY: check-license
check-license:
	(printenv VAULT_LICENSE_CI > /dev/null) || (echo "VAULT_LICENSE_CI must be set"; exit 1)

.PHONY: teardown-integration-test
teardown-integration-test:
	helm uninstall vault --namespace=test || true
	kubectl delete --ignore-not-found namespace test
	# kubectl delete --ignore-not-found clusterrolebinding vault-crb
	# kubectl delete --ignore-not-found clusterrole k8s-clusterrole
	kubectl delete --ignore-not-found --namespace=test -f integrationtest/vault/testBindings.yaml
	kubectl delete --ignore-not-found --namespace=test -f integrationtest/vault/testServiceAccounts.yaml
	kubectl delete --ignore-not-found --namespace=test -f integrationtest/vault/testRoles.yaml
