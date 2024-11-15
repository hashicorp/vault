## Unreleased

## 0.19.0
### Sept 3, 2024

IMPROVEMENTS:
* Updated dependencies: (https://github.com/hashicorp/vault-plugin-auth-alicloud/pull/129)

## 0.18.0
### May 20, 2024

IMPROVEMENTS:
* Updated dependencies: (https://github.com/hashicorp/vault-plugin-auth-alicloud/pull/122)
* bump github.com/hashicorp/go-plugin v1.5.2 to => v1.6.0 to enable running the plugin in containers

## 0.17.0
### Feb 2, 2024

IMPROVEMENTS:
* Updated dependencies:
  * github.com/aliyun/alibaba-cloud-sdk-go v1.62.665 => v1.62.676
  * github.com/hashicorp/vault/api v1.10.0 => v1.11.0
  * github.com/stretchr/testify v1.8.3 => v1.8.4


## 0.16.1
### Jan 23, 2024

IMPROVEMENTS:
* Updated dependencies:
  * github.com/aliyun/alibaba-cloud-sdk-go v1.62.479 -> v1.62.665
  * github.com/hashicorp/go-hclog v1.5.0 -> v1.6.2
  * github.com/hashicorp/go-sockaddr v1.0.2 -> v1.0.6
  * github.com/hashicorp/vault/api v1.9.2 -> v1.10.0
  * github.com/hashicorp/vault/sdk v0.9.2 -> v0.10.2
  * golang.org/x/crypto 0.6.0 -> 0.17.0
  * golang.org/x/net 0.8.0 -> 0.17.0
  * github.com/go-jose/go-jose/v3 3.0.0 -> 3.0.1
  * google.golang.org/grpc 1.57.0 -> 1.57.1
  * github.com/docker/docker 24.0.5+incompatible -> 24.0.7+incompatible


## 0.16.0
### July 27, 2023

IMPROVEMENTS:
* Updated dependencies [GH-87](https://github.com/hashicorp/vault-plugin-auth-alicloud/pull/87):
  * `github.com/aliyun/alibaba-cloud-sdk-go` v1.62.301 -> v1.62.479
  * `github.com/hashicorp/vault/api` v1.9.1 -> v1.9.2
  * `github.com/hashicorp/vault/sdk` v0.9.0 -> v0.9.2

## 0.15.0
### May 24, 2023

IMPROVEMENTS:
* enable plugin multiplexing [GH-61](https://github.com/hashicorp/vault-plugin-auth-alicloud/pull/61)
* update dependencies
  * `github.com/hashicorp/vault/api` v1.9.1
  * `github.com/hashicorp/vault/sdk` v0.9.0
  * `github.com/aliyun/alibaba-cloud-sdk-go` v1.62.301

## 0.14.0
### February 6, 2023

Change:
* Require the `role` field on login

Bug Fixes:
* fix regression in vault login command that caused login to fail

Improvements:
* Update dependencies [GH-43](https://github.com/hashicorp/vault-plugin-auth-alicloud/pull/43)
  * github.com/aliyun/alibaba-cloud-sdk-go v1.61.1842
  * github.com/hashicorp/go-hclog v1.3.1
  * github.com/hashicorp/go-uuid v1.0.3
  * github.com/hashicorp/vault/api v1.8.2
  * github.com/hashicorp/vault/sdk v0.6.1

## 0.12.0
### May 25, 2022

* dep: update golang/x/sys to 9388b58f7150 [[GH-36](https://github.com/hashicorp/vault-plugin-auth-alicloud/pull/36)]
