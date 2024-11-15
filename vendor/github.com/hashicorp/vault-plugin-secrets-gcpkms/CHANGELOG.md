## Unreleased

## v0.19.0 (Sep 9th, 2024)
* Updated dependencies:
  * `github.com/googleapis/enterprise-certificate-proxy` v0.3.3 -> v0.3.4 


## v0.18.0 (Sep 5th, 2024)
* Updated dependencies:
  * `cloud.google.com/go/kms` v1.17.0 -> v1.19.0
  * `github.com/docker/docker` v25.0.5 -> v25.0.6
  * `github.com/hashicorp/vault/api` v1.13.0 -> v1.14.0
  * `github.com/hashicorp/vault/sdk` v0.12.0 -> v0.13.0
  * `golang.org/x/oauth2` v0.20.0 -> v0.23.0
  * `google.golang.org/api` v0.181.0 -> v0.196.0
  * `google.golang.org/genproto` v0.0.0-20240520151616-dc85e6b867a5 -> v0.0.0-20240903143218-8af14fe29dc1
  * `google.golang.org/grpc` v1.64.0 -> v1.66.0
  * `github.com/hashicorp/go-retryablehttp` v0.7.6 -> v0.7.7


## v0.17.0 (May 21st, 2024)
IMPROVEMENTS:
* Updated dependencies:
  * `cloud.google.com/go/kms` v1.15.6 -> v1.17.0
  * `github.com/docker/docker` v25.0.2 -> v25.0.5
  * `github.com/golang/protobuf` v1.5.3 -> v1.5.4
  * `github.com/hashicorp/go-hclog` v1.6.2 -> v1.6.3
  * `github.com/hashicorp/vault/api` v1.11.0 -> v1.13.0
  * `github.com/hashicorp/vault/sdk` v0.10.2 -> v0.12.0
  * `golang.org/x/oauth2` v0.16.0 -> v0.20.0
  * `google.golang.org/api` v0.161.0 -> v0.181.0
  * `google.golang.org/genproto` v0.0.0-20240125205218-1f4bbc51befe -> v0.0.0-20240520151616-dc85e6b867a5
  * `google.golang.org/grpc` v1.61.0 -> v1.64.0


## v0.16.0 (February 2nd, 2024)

IMPROVEMENTS:

* Updated dependencies:
  * `cloud.google.com/go/kms` v1.15.1 -> v1.15.6
  * `github.com/hashicorp/go-hclog` v1.5.0 -> v1.6.2
  * `github.com/hashicorp/vault/api` v1.9.2 -> v1.11.0
  * `github.com/hashicorp/vault/sdk` v0.9.2 -> v0.10.2
  * `golang.org/x/oauth2` v0.11.0 -> v0.16.0
  * `google.golang.org/api` v0.138.0 -> v0.161.0
  * `google.golang.org/genproto` v0.0.0-20230822172742-b8732ec3820d -> v0.0.0-20240125205218-1f4bbc51befe
  * `google.golang.org/grpc` v1.57.0 -> v1.61.0

FIXES:
* Fixed an issue where the newly introduced key metadata cache would cause a nil dereference error [GH-37](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/37)

## v0.15.2 (September 22th, 2023)

IMPROVEMENTS:

* Added cache on key metadata lookup to avoid quota issues [GH-33](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/33)

## v0.15.1 (September 5th, 2023)

IMPROVEMENTS:

* Updated dependencies:
   * `cloud.google.com/go/kms` v1.6.0 -> v1.15.1
   * `github.com/gammazero/workerpool` v0.0.0-20190406235159-88d534f22b56 -> v1.1.3
   * `github.com/golang/protobuf` v1.5.2 -> v1.5.3
   * `github.com/hashicorp/go-hclog` v0.16.2 -> v1.5.0
   * `github.com/hashicorp/vault/api` v1.9.0 -> v1.9.2
   * `github.com/hashicorp/vault/sdk` v0.9.0 -> v0.9.2
   * `golang.org/x/oauth2` v0.4.0 -> v0.11.0
   * `google.golang.org/api` v0.103.0 -> v0.138.0
   * `google.golang.org/genproto` v0.0.0-20230110181048-76db0878b65f -> v0.0.0-20230822172742-b8732ec3820d
   * `google.golang.org/grpc` v1.47.0 -> v1.57.0

## v0.15.0

IMPROVEMENTS:

* enable plugin multiplexing [GH-26](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/26)
* update dependencies
  * `github.com/hashicorp/vault/api` v1.9.0
  * `github.com/hashicorp/vault/sdk` v0.8.1

## v0.14.0

CHANGES:

* Changes user-agent header value to use correct Vault version information and include
  the plugin type and name in the comment section. [[GH-21](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/21)]
* CreateOperation should only be implemented alongside ExistenceCheck [[GH-20](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/20)]

IMPROVEMENTS:

* Dependency updates
  * google.golang.org/api v0.5.0 => v0.83.0
  * github.com/hashicorp/vault/api v1.0.5-0.20200215224050-f6547fa8e820 => v1.8.3
  * github.com/hashicorp/vault/sdk v0.1.14-0.20200215224050-f6547fa8e820 => v0.7.0
  * golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 => v0.5.0
  * golang.org/x/net 0.0.0-20220722155237-a158d28d115b => v0.5.0
