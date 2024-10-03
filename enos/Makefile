.PHONY: default
default: check-fmt shellcheck

.PHONY: check-fmt
check-fmt: check-fmt-enos check-fmt-modules check-shfmt

.PHONY: fmt
fmt: fmt-enos fmt-modules shfmt

.PHONY: check-fmt-enos
check-fmt-enos:
	enos fmt --check --diff .

.PHONY: fmt-enos
fmt-enos:
	enos fmt .

.PHONY: check-fmt-modules
check-fmt-modules:
	terraform fmt -check -diff -recursive ./modules

.PHONY: fmt-modules
fmt-modules:
	terraform fmt -diff -recursive ./modules

# Generate dynamic global variables. We execute this in CI so that each branch has the correct
# config.
.PHONY: gen-dynamic-globals
gen-dynamic-globals:
	pushd ../tools/ &> /dev/null &&\
	upgrade_initial_versions=$$(go run ./pipeline/... releases list-versions -n 3 -u $$(cat ../version/VERSION) -e $(edition) | jq -Mec .versions) &&\
	popd &> /dev/null &&\
	echo "globals {\n  upgrade_initial_versions = $$upgrade_initial_versions\n}" > enos-dynamic-globals.hcl

.PHONY: validate-enos
validate-enos:
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
