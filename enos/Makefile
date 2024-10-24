VAULT_VERSION=$$(cat $(CURDIR)/../version/VERSION)

.PHONY: default
default: check-fmt shellcheck

.PHONY: check-fmt
check-fmt: check-fmt-enos check-fmt-modules check-shfmt

.PHONY: fmt
fmt: fmt-enos fmt-modules shfmt

.PHONY: check-fmt-enos
check-fmt-enos:
	enos fmt --check --diff .
	enos fmt --check --diff ./k8s

.PHONY: fmt-enos
fmt-enos:
	enos fmt .
	enos fmt ./k8s

.PHONY: gen-enos
gen-enos:
	pushd ../tools/pipeline &> /dev/null && go run ./... generate enos-dynamic-config -d ../../enos -f enos-dynamic-config.hcl -e ce -v $(VAULT_VERSION) -n 3 --log info && popd &> /dev/null

.PHONY: check-fmt-modules
check-fmt-modules:
	terraform fmt -check -diff -recursive ./modules

.PHONY: fmt-modules
fmt-modules:
	terraform fmt -diff -recursive ./modules

.PHONY: validate-enos
validate-enos:
	enos scenario validate --timeout 30m0s --chdir ./k8s
	enos scenario validate --timeout 30m0s

.PHONY: lint
lint: check-fmt check-fmt-modules check-shfmt shellcheck validate-enos

.PHONY: shellcheck
shellcheck:
	find ./modules/ -type f -name '*.sh' | xargs shellcheck

.PHONY: shfmt
shfmt:
	find ./modules/ -type f -name '*.sh' | xargs shfmt -l -w -i 2 -bn -ci -kp -sr

.PHONY: check-shfmt
check-shfmt:
	find ./modules/ -type f -name '*.sh' | xargs shfmt -l -d -i 2 -bn -ci -kp -sr
