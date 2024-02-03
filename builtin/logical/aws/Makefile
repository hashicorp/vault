GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt build start

build:
	GO111MODULE=on GOOS=$(OS) GOARCH="$(GOARCH)" go build -o builds/awsnew cmd/aws/main.go

start:
	vault server -dev -dev-root-token-id=root -log-level=trace -dev-listen-address=:8201 -dev-plugin-dir=./builds/ &

enable:
	vault secrets enable -path=awsupdated awsnew

config: 
	vault write awsupdated/roles/network \
    role_arns=arn:aws:iam::578479370890:role/network-admin-role \
    credential_type=assumed_role console_login=true \
	console_duration=1800

deploy:
	vault write awsupdated/sts/network ttl=60m

clean:
	rm -f ./vault/plugins/aws

fmt:
	go fmt $$(go list ./...)

register:
	export SHA256=$(sha256sum ./vault/plugins/awsnew | cut -d ' ' -f1)
	echo $SHA256
	vault plugin register -sha256=$SHA256 secret new

.PHONY: build clean fmt start enable
