module github.com/hashicorp/vault/api/auth/gcp

go 1.16

replace github.com/hashicorp/vault/api => ../../../api

require (
	cloud.google.com/go v0.97.0
	github.com/hashicorp/vault/api v1.2.0
	google.golang.org/genproto v0.0.0-20210924002016-3dee208752a0
)
