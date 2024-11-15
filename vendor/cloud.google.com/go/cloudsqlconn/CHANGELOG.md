# Changelog

## [1.13.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.12.1...v1.13.0) (2024-10-23)


### Features

* Automatically reset connection when the DNS record changes. ([#868](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/868)) ([4d7abd8](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4d7abd877edf5fba3173b69e14181b6ddf911b24))


### Bug Fixes

* update bytes_sent and bytes_received to use Sum ([#874](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/874)) ([73b6f38](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/73b6f3860ef28dedd995a41b74b5f12168d3ff06))

## [1.12.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.12.0...v1.12.1) (2024-09-19)


### Bug Fixes

* update dependencies to latest versions ([#872](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/872)) ([4eed622](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4eed622e482a1fbcaecaf16124c445a0f7509e0c))

## [1.12.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.11.1...v1.12.0) (2024-08-13)


### Features

* add `bytes_sent` and `bytes_received` as metrics ([#856](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/856)) ([d0e493f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/d0e493fc3859debd625e56874c4df32aeca02403))
* add support for Go 1.23 and drop Go 1.20 ([#860](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/860)) ([8ce98e8](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/8ce98e858c236efccd5eb21a84f24c4b20f4a2cb))
* Configure connections using DNS domain names ([#843](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/843)) ([ec6b3a0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/ec6b3a09bdfd1e13df30786e973ccecd48e9b3a6))
* support Cloud SQL CAS instances. ([#850](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/850)) ([511fae4](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/511fae491ed1101c2ce0998120291e0cb8180d40))

## [1.11.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.11.0...v1.11.1) (2024-07-10)


### Bug Fixes

* bump dependencies to latest versions ([#839](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/839)) ([ce7f28f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/ce7f28ff56481d9cfd4031d940fcd4fcd61219ee))

## [1.11.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.10.1...v1.11.0) (2024-06-12)


### Features

* generate RSA key lazily for lazy refresh ([#826](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/826)) ([bf293e2](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/bf293e25e2d52f395734c597c86dfe85ede5f4cd)), closes [#823](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/823)
* invalidate cache on failed IP lookup ([#812](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/812)) ([4b68de3](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4b68de3693e25642acd847d0c8ac393982d00c9b)), closes [#780](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/780)


### Bug Fixes

* ensure connection count is correctly reported ([#824](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/824)) ([b286049](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/b286049a7ade2a9e3cf44ea36f56946cfa58f60a))
* invalidate cache on failed `Warmup` and `EngineVersion` ([#827](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/827)) ([c3915a6](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/c3915a6790f3d4e3cff266a0d8c506a09ecf9634))

## [1.10.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.10.0...v1.10.1) (2024-05-22)


### Bug Fixes

* remove duplicate refresh operations ([#806](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/806)) ([beb3605](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/beb36052af2221d7ff238edc4c98c733cac2999d)), closes [#771](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/771)

## [1.10.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.9.0...v1.10.0) (2024-05-14)


### Features

* expose context to debug logger ([#797](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/797)) ([847f7c1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/847f7c10cc796761e81a86e0551f00832a5056d5))


### Bug Fixes

* retry 50x errors with exponential backoff ([#781](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/781)) ([40dc789](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/40dc789baabbe40cebabee7a287222940b120e6a))

## [1.9.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.8.1...v1.9.0) (2024-04-16)


### Features

* add support for a lazy refresh ([#772](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/772)) ([931150f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/931150f492cb461cf623a9bbafae6f704b9c5a36)), closes [#770](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/770)


### Bug Fixes

* return a friendly error if the dialer is closed ([#766](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/766)) ([d1c13e0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/d1c13e039a29ccbc085e2d3ca8451f83825e8d32))

## [1.8.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.8.0...v1.8.1) (2024-03-12)


### Bug Fixes

* strip monotonic clock reading in cert check ([#750](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/750)) ([6ae33b0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/6ae33b0a6e281293823e75ff97a51575c053bf9f)), closes [#749](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/749)

## [1.8.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.7.0...v1.8.0) (2024-03-08)


### Features

* add support for TPC ([#732](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/732)) ([b7364d9](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/b7364d93cc93893b2af8eeda6cdf9cf36aaf9d67))

## [1.7.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.6.0...v1.7.0) (2024-02-13)


### Features

* add support for debug logging ([#726](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/726)) ([d8ca89e](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/d8ca89e4403e2e3cf6ac278a19b4d93b77797ec6))
* add support for Go 1.22 ([#723](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/723)) ([ebe31dc](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/ebe31dcaf2ec215470ce3b224732f4ff6282ba22))


### Bug Fixes

* ensure background refresh is closed cleanly ([#715](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/715)) ([0b4c342](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/0b4c3420bb5158cab63c51158e109b3bea926b59))

## [1.6.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.5.2...v1.6.0) (2024-01-17)


### Features

* add connection name to public API ([#698](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/698)) ([84f3b6e](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/84f3b6eedcf13402bcbf7da720924cf242893beb))

## [1.5.2](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.5.1...v1.5.2) (2023-12-12)


### Bug Fixes

* ensure cert refresh recovers from sleep ([#686](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/686)) ([95671ad](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/95671ada40905cf14209b5c54058463689ce6b20))

## [1.5.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.5.0...v1.5.1) (2023-11-14)


### Bug Fixes

* bump dependencies to latest ([#667](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/667)) ([86544f5](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/86544f5a477f694c8ceb862b13c3b83d19d72d5d))

## [1.5.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.4.5...v1.5.0) (2023-10-24)


### Features

* add pgx v5 support ([#639](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/639)) ([#642](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/642)) ([8d86d92](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/8d86d92147d06ca10d754439638d6fd1b2154182))


### Bug Fixes

* use different driver names for v4 and v5 testing ([#639](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/639)) ([#654](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/654)) ([fa73c41](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/fa73c4184a9887e6e9217e5b50db97aa3fdc0d28))
* use HandshakeContext by default ([#656](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/656)) ([49aad1f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/49aad1f30bf560e6cf1e2ff52da46f3ff2cd2312))

## [1.4.5](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.4.4...v1.4.5) (2023-10-11)


### Bug Fixes

* bump dependencies to latest ([#649](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/649)) ([0ddac9f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/0ddac9fa7de17f740021408ed25ffbb0b0133d9e))
* bump minimum supported Go version to 1.19 ([#637](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/637)) ([4a28a78](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4a28a788a94d64e1ce6ddd76fa3a041c82c8f2b1))

## [1.4.4](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.4.3...v1.4.4) (2023-09-12)


### Bug Fixes

* update dependencies to latest versions ([#621](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/621)) ([32f1e27](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/32f1e2762b8ced0a3332e4928fdc61ad5d731530))

## [1.4.3](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.4.2...v1.4.3) (2023-08-18)


### Bug Fixes

* update ForceRefresh to block if invalid ([#605](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/605)) ([61c72e3](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/61c72e3e76d04863b6971aeb86726c3b1252e5ed))

## [1.4.2](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.4.1...v1.4.2) (2023-08-15)


### Bug Fixes

* re-use existing connection info on force refresh ([#602](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/602)) ([d049851](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/d049851361fc48bb339232c6609a2f2932d2d684))

## [1.4.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.4.0...v1.4.1) (2023-08-07)


### Bug Fixes

* avoid holding lock over IO ([#576](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/576)) ([1e4560f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/1e4560f7b41547882a2e9f7ef3ece94bb1bb48be))

## [1.4.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.3.0...v1.4.0) (2023-07-06)


### Features

* add support for PSC connections ([#565](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/565)) ([10a46b0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/10a46b0a36440d6b84498468346833729c21bbb4))

## [1.3.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.2.4...v1.3.0) (2023-06-13)


### Features

* add support for WithOneOffDialFunc ([#558](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/558)) ([14592f3](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/14592f3d21e58fbd038cffdb6c4f67d7e3526302))


### Bug Fixes

* close background refresh for bad instances ([#550](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/550)) ([31f06fc](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/31f06fc078f097b6cef4f7c19228a724a00c3408))

## [1.2.4](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.2.3...v1.2.4) (2023-05-09)


### Bug Fixes

* update dependencies to latest versions ([#539](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/539)) ([f1a4008](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/f1a40083289ef0051b757f7a12921cfefc65a249))

## [1.2.3](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.2.2...v1.2.3) (2023-04-11)


### Bug Fixes

* update dependencies to latest versions ([#517](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/517)) ([55bad80](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/55bad80b3ae64b4b9c7135db2c12dd49e0ad230e))

## [1.2.2](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.2.1...v1.2.2) (2023-03-09)


### Bug Fixes

* strip monotonic clock readings for refresh calculations ([#471](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/471)) ([94048af](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/94048afd001fd960f316e961501b871ab648296e))

## [1.2.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.2.0...v1.2.1) (2023-02-15)


### Bug Fixes

* don't initialize default creds when using a token ([#460](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/460)) ([fc5c435](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/fc5c435b92ddfe6be5bbe77264486c0b712ba4d1))

## [1.2.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.1.1...v1.2.0) (2023-02-14)


### Features

* add support for Go 1.20 ([#445](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/445)) ([4df53ef](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4df53ef4e742d6cd4c80bb79ed90d7ecd2110868))


### Bug Fixes

* error when dialer is misconfigured with token source ([#453](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/453)) ([7b45a7e](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/7b45a7e27c164dbf1f7903ed7792e4d81dd467b7))
* improve reliability of certificate refresh ([#448](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/448)) ([47bd3f3](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/47bd3f385ad0cc7bbd057f3273ed03d2587e9ac8))
* prevent repeated context expired errors ([#458](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/458)) ([7ffeafe](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/7ffeafea9729d08ad04c403c07b70d4f184664a0))

## [1.1.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.1.0...v1.1.1) (2023-01-10)


### Bug Fixes

* move MySQL liveness check into driver code ([#417](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/417)) ([0de68fb](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/0de68fbc32d87e4cabab301be8a11f9eba50e13d))
* use handshake context when possible ([#427](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/427)) ([37c4e70](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/37c4e70aa7082c49b84aaedb2066ddb67e1d920f))

## [1.1.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.0.1...v1.1.0) (2022-12-06)


### Features

* add support for MySQL Auto IAM AuthN ([#309](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/309)) ([6c4f20e](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/6c4f20eae857c215098b7b991fffc7d15bbead5b))
* improve refresh duration calculation ([#364](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/364)) ([10b0bf7](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/10b0bf7d9d3c69238df3d0a88ffab54f03f7d7a6))


### Bug Fixes

* handle context cancellations during instance refresh ([#372](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/372)) ([cdb59c7](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/cdb59c797968f46419673378c96e79d40da453dc)), closes [#370](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/370)
* remove leading slash from metric names ([#393](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/393)) ([ac5ca26](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/ac5ca264e17adf0c5780ea2317f4df03c6e1923d))

## [1.0.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v1.0.0...v1.0.1) (2022-11-01)


### Bug Fixes

* update dependencies to latest versions ([#365](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/365)) ([5479502](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/547950268712f48d8613aac3d7e2a1e494b6a680))

## [1.0.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.5.2...v1.0.0) (2022-10-18)


### Features

* add WithAutoIP option ([#346](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/346)) ([bd20b6b](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/bd20b6bfe746cfea778b9e1a9702de28047e5950))
* Downscope OAuth2 token included in ephemeral certificate ([#332](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/332)) ([d13dd6f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/d13dd6f3e7db0179511539315dec1c2dc96f0e3e))


### Bug Fixes

* throw error when Auto IAM AuthN is unsupported ([#310](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/310)) ([652e196](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/652e196b427ce9673676e214c6ad3905b21a68b0))


### Miscellaneous Chores

* set next version to v1.0.0 ([#349](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/349)) ([a76d2db](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/a76d2db0b31447dc96707679973ff87b3c755bf5))

## [0.5.2](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.5.1...v0.5.2) (2022-09-07)


### Bug Fixes

* update dependencies to latest versions ([#300](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/300)) ([5504df6](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/5504df6e03bda7b56e01146e63b715f775443d85))

## [0.5.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.5.0...v0.5.1) (2022-08-01)


### Bug Fixes

* remove unnecessary import path restrictions ([#258](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/258)) ([bc57877](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/bc57877f16a61e42c603d4dc50ff4d01fc01d9d9))

## [0.5.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.4.0...v0.5.0) (2022-07-12)


### Features

* expose the WithQuotaProject dialer option ([#237](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/237)) ([bda8917](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/bda891776d5d44d49ed3e4a268f27bd10a23427e))


### Bug Fixes

* support MySQL driver’s conn check. ([#226](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/226)) ([4b48e3b](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4b48e3bfe7a5bd8c398592f21eb25ac43644e123))

## [0.4.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.3.1...v0.4.0) (2022-06-07)


### Features

* add DialOption for IAM DB Authentication ([#171](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/171)) ([c103acc](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/c103acc6b49f94a1a733dc0e5c8b41890172dd8b))
* Add Warmup function for starting background refresh ([#163](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/163)) ([2459f92](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/2459f92911eeca46102f56966c8cefa7cee8a0ae))


### Bug Fixes

* adjust alignment for 32-bit arch ([#197](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/197)) ([86e96ad](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/86e96adf30cbc82ba170dc70ce4d0694a3b595ce))

### [0.3.1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.3.0...v0.3.1) (2022-05-03)


### Bug Fixes

* update dependencies to latest versions ([#185](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/185)) ([702a380](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/702a3802d0383c0d71277779d80d62a5e5c23157))

## [0.3.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.2.0...v0.3.0) (2022-04-04)


### Features

* add option to configure SQL Admin API URL ([#148](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/148)) ([c791369](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/c79136972083480d16f65a4696a7747bae942afe))
* add WithUserAgent opt ([#156](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/156)) ([bd89dc5](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/bd89dc50bb50d1d6ff9cf5a146071b307a54683a))
* drop support for Go 1.15 ([#145](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/145)) ([791641b](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/791641bb2d0ab93955b218b9bc6f5335b8ead243))
* use connect API for instance metadata ([#150](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/150)) ([1086ad0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/1086ad01cc7907051147d572f4f27ab1ba538027))


### Bug Fixes

* memory leak in database/sql integration ([#162](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/162)) ([47cdf2d](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/47cdf2da2230801b591bf4f459bfcbe7e9432cd1))
* prevent unnecessary allocation of conn config ([#164](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/164)) ([49c7828](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/49c782809aff84b6141027f1a2634b0a0db2b18a))

## [0.2.0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/compare/v0.1.0...v0.2.0) (2022-03-01)


### ⚠ BREAKING CHANGES

* use singular name for package (#101)

### Features

* add dial_failure_count metric ([#127](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/127)) ([34cdbb9](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/34cdbb92efa6f186bd8afdde3c8dcc810e77911e))
* add metrics for refresh success and failure ([#133](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/133)) ([a36a212](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/a36a212dbd30474721669f10fbfda1e76a22d325))
* drop support and testing for Go 1.14 ([#128](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/128)) ([aceadcc](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/aceadcc4835b6fe18639a696755302bb00f82bc2))


### Bug Fixes

* custom drivers report error on cleanup ([#102](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/102)) ([648b75a](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/648b75a4d8e43b3641d827086047a9c6783c1306))
* use singular name for package ([#101](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/101)) ([5e5589d](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/5e5589db3bb0a86d9c167cd6b85358535238176a))


## 0.1.0 (2022-02-08)


### ⚠ BREAKING CHANGES

* remove singleton Dial (#92)
* return cleanup func to close dialer (#75)
* dialer is a io.Closer (#76)
* initialize dialer in register func (#73)
* rename DialerOption to Option (#64)

### Features

* Add Close method to Dialer ([#34](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/34)) ([91ee305](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/91ee305b6af83d48ba5fc445ad1191fd99785079))
* add concrete errors to public API ([#36](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/36)) ([7441b71](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/7441b7176d8bce5d2e054aa7e53f1509aece9898))
* add custom driver for MySQL ([#70](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/70)) ([755c334](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/755c3344f28e33d18a1d7acc414352ee73e39d8a))
* add custom driver for SQL Server ([#71](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/71)) ([14eb60a](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/14eb60a88532dd81cda4d602d044c98013ee0af6))
* add default useragent ([#17](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/17)) ([57d7ed9](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/57d7ed9da73c731196bdc5120134b6dec72d9c68))
* Add DialerOption for specifying a refresh timeout ([#12](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/12)) ([94df7cf](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/94df7cfa21dc60463afb1ad3519455d507d610f3))
* add DialOptions for configuring Dial  ([#8](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/8)) ([e2d53ee](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/e2d53ee6c66ba58114d8a49ca86f0eb3a56ce481))
* Add EngineVersion method to Dialer ([#59](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/59)) ([6a78bfd](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/6a78bfd4a73807e4fce455ae0d6cd4f531710edd))
* Add initial dialer ([#1](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/1)) ([7e89552](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/7e8955216cc91999e3d8d17ed9eced8f63564ca7))
* add initial support for metrics ([#40](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/40)) ([ee396ff](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/ee396fffb10ea52af9072d0fdd09a8b4e9d4b736))
* add support for configuring the HTTP client ([#55](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/55)) ([de9e72e](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/de9e72e1dc6961f6b6ed3fe9cf4381344dd5fa37))
* add support for IAM DB Authn ([#44](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/44)) ([92e28cf](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/92e28cfccd573c0908588ad3594ef9de403e5e51))
* add support for tracing ([#32](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/32)) ([4d2acbc](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4d2acbcecb11acbbc58f95c711051a02fb31e82f))
* allow for configuring the Dial func ([#57](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/57)) ([4cb523e](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4cb523e80b4a388b37c8ce251a533a3b8d370029))
* expose Dialer and add DialerOptions ([#7](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/7)) ([1235a9f](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/1235a9f62beb678f18695afc6d22d0b8e6b7b506))
* force early refresh of instance info if connect fails ([#19](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/19)) ([eb06ae2](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/eb06ae26609cbc46fa65e50c080508d53ec0b9c2))
* improve reliablity of refresh operations ([#49](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/49)) ([3a52440](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/3a5244075f68f3c95f26218f9008bb7451934f80))
* improve RSA keypair generation ([#10](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/10)) ([e2a5238](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/e2a52388ff047144272089db60cb0b1fce7c16bf))
* initialize dialer in register func ([#73](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/73)) ([7633cfd](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/7633cfd2eaadeef065686f85ae9f2faa5087e917))
* **postgres/pgxv4:** add support for postgres driver ([#61](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/61)) ([295a5dc](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/295a5dcfbdaeb12884333e678f8b9f7f44de2b46))
* remove singleton Dial ([#92](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/92)) ([0a1966c](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/0a1966c4fe0400e8dcd14b2531db20ad7bc10855))
* return cleanup func to close dialer ([#75](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/75)) ([fa9b845](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/fa9b84576a7adcf8f0ad4296723685d681ada89e))
* use cloud.google.com/go/cloudsqlconn ([#30](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/30)) ([a251fd7](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/a251fd727813223dc08f40bc5060add3235564e6))


### Bug Fixes

* dialer is a io.Closer ([#76](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/76)) ([89de96c](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/89de96c2a4d636cc3dfe44aa1b47ab3492d5cf0c))
* perform refresh operations asynchronously ([#11](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/11)) ([925d6c2](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/925d6c2686d519d182dc196c752ed0c7edb0e28c))
* rate limit refresh attempts per instance ([#18](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/18)) ([1092ccc](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/1092ccc04361293f6ea07fdc97cde30cf1cb1866))
* rename DialerOption to Option ([#64](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/64)) ([016a821](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/016a821ba191b7b2117c7d240507e32c289e3f0e))
* schedule refreshes based on result expiration instead of fixed interval ([#21](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/21)) ([65073d0](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/65073d0ea9582abbe01c7ca0698681624e3c7834))
* **trace:** use LastValue for open connections ([#58](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/58)) ([4ee6bea](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/4ee6bea069c196454dd48034457a16ba416b725c))
* use ctx for NewService ([#24](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/issues/24)) ([77fd677](https://github.com/GoogleCloudPlatform/cloud-sql-go-connector/commit/77fd677ccb827feb89e6bb41eb45c22f3a2b1861))
