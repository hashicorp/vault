.PHONY: default
default: check-fmt

.PHONY: check-fmt
check-fmt: check-fmt-enos check-fmt-modules

.PHONY: fmt
fmt: fmt-enos fmt-modules

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
