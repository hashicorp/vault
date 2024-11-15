# Changes

## [1.73.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.72.0...spanner/v1.73.0) (2024-11-14)


### Features

* **spanner:** Add ResetForRetry method for stmt-based transactions ([#10956](https://github.com/googleapis/google-cloud-go/issues/10956)) ([02c191c](https://github.com/googleapis/google-cloud-go/commit/02c191c5dc13023857812217f63be2395bfcb382))


### Bug Fixes

* **spanner:** Add safecheck to avoid deadlock when creating multiplex session ([#11131](https://github.com/googleapis/google-cloud-go/issues/11131)) ([8ee5d05](https://github.com/googleapis/google-cloud-go/commit/8ee5d05e288c7105ddb1722071d6719933effea4))
* **spanner:** Allow non default service account only when direct path is enabled ([#11046](https://github.com/googleapis/google-cloud-go/issues/11046)) ([4250788](https://github.com/googleapis/google-cloud-go/commit/42507887523f41d0507ca8b1772235846947c3e0))
* **spanner:** Use spanner options when initializing monitoring exporter ([#11109](https://github.com/googleapis/google-cloud-go/issues/11109)) ([81413f3](https://github.com/googleapis/google-cloud-go/commit/81413f3647a0ea406f25d4159db19b7ad9f0682b))

## [1.72.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.71.0...spanner/v1.72.0) (2024-11-07)


### Features

* **spanner/spansql:** Add support for protobuf column types & Proto bundles ([#10945](https://github.com/googleapis/google-cloud-go/issues/10945)) ([91c6f0f](https://github.com/googleapis/google-cloud-go/commit/91c6f0fcaadfb7bd983e070e6ceffc8aeba7d5a2)), refs [#10944](https://github.com/googleapis/google-cloud-go/issues/10944)


### Bug Fixes

* **spanner:** Skip exporting metrics if attempt or operation is not captured. ([#11095](https://github.com/googleapis/google-cloud-go/issues/11095)) ([1d074b5](https://github.com/googleapis/google-cloud-go/commit/1d074b520c7a368fb8a7a27574ef56a120665c64))

## [1.71.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.70.0...spanner/v1.71.0) (2024-11-01)


### Features

* **spanner/admin/instance:** Add support for Cloud Spanner Default Backup Schedules ([706ecb2](https://github.com/googleapis/google-cloud-go/commit/706ecb2c813da3109035b986a642ca891a33847f))
* **spanner:** Client built in metrics ([#10998](https://github.com/googleapis/google-cloud-go/issues/10998)) ([d81a1a7](https://github.com/googleapis/google-cloud-go/commit/d81a1a75b9efbf7104bb077300364f8d63da89b5))


### Bug Fixes

* **spanner/test/opentelemetry/test:** Update google.golang.org/api to v0.203.0 ([8bb87d5](https://github.com/googleapis/google-cloud-go/commit/8bb87d56af1cba736e0fe243979723e747e5e11e))
* **spanner/test/opentelemetry/test:** WARNING: On approximately Dec 1, 2024, an update to Protobuf will change service registration function signatures to use an interface instead of a concrete type in generated .pb.go files. This change is expected to affect very few if any users of this client library. For more information, see https://togithub.com/googleapis/google-cloud-go/issues/11020. ([2b8ca4b](https://github.com/googleapis/google-cloud-go/commit/2b8ca4b4127ce3025c7a21cc7247510e07cc5625))
* **spanner:** Attempt latency for streaming call should capture the total latency till decoding of protos ([#11039](https://github.com/googleapis/google-cloud-go/issues/11039)) ([255c6bf](https://github.com/googleapis/google-cloud-go/commit/255c6bfcdd3e844dcf602a829bfa2ce495bcd72e))
* **spanner:** Decode PROTO to custom type variant of base type ([#11007](https://github.com/googleapis/google-cloud-go/issues/11007)) ([5e363a3](https://github.com/googleapis/google-cloud-go/commit/5e363a31cc9f2616832540ca82aa5cb998a3938c))
* **spanner:** Update google.golang.org/api to v0.203.0 ([8bb87d5](https://github.com/googleapis/google-cloud-go/commit/8bb87d56af1cba736e0fe243979723e747e5e11e))
* **spanner:** WARNING: On approximately Dec 1, 2024, an update to Protobuf will change service registration function signatures to use an interface instead of a concrete type in generated .pb.go files. This change is expected to affect very few if any users of this client library. For more information, see https://togithub.com/googleapis/google-cloud-go/issues/11020. ([2b8ca4b](https://github.com/googleapis/google-cloud-go/commit/2b8ca4b4127ce3025c7a21cc7247510e07cc5625))

## [1.70.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.69.0...spanner/v1.70.0) (2024-10-14)


### Features

* **spanner/admin/instance:** Define ReplicaComputeCapacity and AsymmetricAutoscalingOption ([78d8513](https://github.com/googleapis/google-cloud-go/commit/78d8513f7e31c6ef118bdfc784049b8c7f1e3249))
* **spanner:** Add INTERVAL API ([78d8513](https://github.com/googleapis/google-cloud-go/commit/78d8513f7e31c6ef118bdfc784049b8c7f1e3249))
* **spanner:** Add new QueryMode enum values (WITH_STATS, WITH_PLAN_AND_STATS) ([78d8513](https://github.com/googleapis/google-cloud-go/commit/78d8513f7e31c6ef118bdfc784049b8c7f1e3249))


### Documentation

* **spanner/admin/instance:** A comment for field `node_count` in message `spanner.admin.instance.v1.Instance` is changed ([78d8513](https://github.com/googleapis/google-cloud-go/commit/78d8513f7e31c6ef118bdfc784049b8c7f1e3249))
* **spanner/admin/instance:** A comment for field `processing_units` in message `spanner.admin.instance.v1.Instance` is changed ([78d8513](https://github.com/googleapis/google-cloud-go/commit/78d8513f7e31c6ef118bdfc784049b8c7f1e3249))
* **spanner:** Update comment for PROFILE QueryMode ([78d8513](https://github.com/googleapis/google-cloud-go/commit/78d8513f7e31c6ef118bdfc784049b8c7f1e3249))

## [1.69.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.68.0...spanner/v1.69.0) (2024-10-03)


### Features

* **spanner:** Add x-goog-spanner-end-to-end-tracing header for requests to Spanner ([#10241](https://github.com/googleapis/google-cloud-go/issues/10241)) ([7f61cd5](https://github.com/googleapis/google-cloud-go/commit/7f61cd579f7e4ed4f1ac161f2c2a28e931406f16))


### Bug Fixes

* **spanner:** Handle errors ([#10943](https://github.com/googleapis/google-cloud-go/issues/10943)) ([c67f964](https://github.com/googleapis/google-cloud-go/commit/c67f964de364808c02085dda61fa53e2b2fda850))


### Performance Improvements

* **spanner:** Use passthrough with emulator endpoint ([#10947](https://github.com/googleapis/google-cloud-go/issues/10947)) ([9e964dd](https://github.com/googleapis/google-cloud-go/commit/9e964ddc01a54819f25435cfcc9d5b37c91f5a1d))

## [1.68.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.67.0...spanner/v1.68.0) (2024-09-25)


### Features

* **spanner:** Add support for Go 1.23 iterators ([84461c0](https://github.com/googleapis/google-cloud-go/commit/84461c0ba464ec2f951987ba60030e37c8a8fc18))


### Bug Fixes

* **spanner/test:** Bump dependencies ([2ddeb15](https://github.com/googleapis/google-cloud-go/commit/2ddeb1544a53188a7592046b98913982f1b0cf04))
* **spanner:** Bump dependencies ([2ddeb15](https://github.com/googleapis/google-cloud-go/commit/2ddeb1544a53188a7592046b98913982f1b0cf04))
* **spanner:** Check errors in tests ([#10738](https://github.com/googleapis/google-cloud-go/issues/10738)) ([971bfb8](https://github.com/googleapis/google-cloud-go/commit/971bfb85ee7bf8c636117a6424280a4323b5fb3c))
* **spanner:** Enable toStruct support for structs with proto message pointer fields ([#10704](https://github.com/googleapis/google-cloud-go/issues/10704)) ([42cdde6](https://github.com/googleapis/google-cloud-go/commit/42cdde6ee34fc9058dc47c9c9ab39ba91b6b9c58))
* **spanner:** Ensure defers run at the right time in tests ([#9759](https://github.com/googleapis/google-cloud-go/issues/9759)) ([7ef0ded](https://github.com/googleapis/google-cloud-go/commit/7ef0ded2502dbb37f07bc93bc2e868e29f7121c4))
* **spanner:** Increase spanner ping timeout to give backend more time to process executeSQL requests ([#10874](https://github.com/googleapis/google-cloud-go/issues/10874)) ([6997991](https://github.com/googleapis/google-cloud-go/commit/6997991e2325e7a66d3ffa60c27622a1a13041a8))
* **spanner:** Json null handling ([#10660](https://github.com/googleapis/google-cloud-go/issues/10660)) ([4c519e3](https://github.com/googleapis/google-cloud-go/commit/4c519e37a124defc3451adfdbd0883a5e081eb2f))
* **spanner:** Support custom encoding and decoding of protos ([#10799](https://github.com/googleapis/google-cloud-go/issues/10799)) ([d410907](https://github.com/googleapis/google-cloud-go/commit/d410907f3e52bcc64bd92e0a341777c1277a6418))
* **spanner:** Unnecessary string formatting fixes ([#10736](https://github.com/googleapis/google-cloud-go/issues/10736)) ([1efe5c4](https://github.com/googleapis/google-cloud-go/commit/1efe5c4275dca6d739691e89b8d460b97160d953))
* **spanner:** Wait for things to complete ([#10095](https://github.com/googleapis/google-cloud-go/issues/10095)) ([7785cad](https://github.com/googleapis/google-cloud-go/commit/7785cad89effbc8c4e67043368f96d4768cdb40f))


### Performance Improvements

* **spanner:** Better error handling ([#10734](https://github.com/googleapis/google-cloud-go/issues/10734)) ([c342f65](https://github.com/googleapis/google-cloud-go/commit/c342f6550c24e3a16e32d1cd61c6fcfeaed77c7b)), refs [#9749](https://github.com/googleapis/google-cloud-go/issues/9749)


### Documentation

* **spanner:** Fix Key related document code to add package name ([#10711](https://github.com/googleapis/google-cloud-go/issues/10711)) ([bbe7b9c](https://github.com/googleapis/google-cloud-go/commit/bbe7b9ceed1deb85a4f40ea95572595ce63ff002))

## [1.67.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.66.0...spanner/v1.67.0) (2024-08-15)


### Features

* **spanner/admin/database:** Add resource reference annotation to backup schedules ([#10677](https://github.com/googleapis/google-cloud-go/issues/10677)) ([6593c0d](https://github.com/googleapis/google-cloud-go/commit/6593c0d62d48751c857bce3d3f858127467a4489))
* **spanner/admin/instance:** Add edition field to the instance proto ([6593c0d](https://github.com/googleapis/google-cloud-go/commit/6593c0d62d48751c857bce3d3f858127467a4489))
* **spanner:** Support commit options in mutation operations. ([#10668](https://github.com/googleapis/google-cloud-go/issues/10668)) ([62a56f9](https://github.com/googleapis/google-cloud-go/commit/62a56f953d3b8fe82083c42926831c2728312b9c))


### Bug Fixes

* **spanner/test/opentelemetry/test:** Update google.golang.org/api to v0.191.0 ([5b32644](https://github.com/googleapis/google-cloud-go/commit/5b32644eb82eb6bd6021f80b4fad471c60fb9d73))
* **spanner:** Update google.golang.org/api to v0.191.0 ([5b32644](https://github.com/googleapis/google-cloud-go/commit/5b32644eb82eb6bd6021f80b4fad471c60fb9d73))


### Documentation

* **spanner/admin/database:** Add an example to filter backups based on schedule name ([6593c0d](https://github.com/googleapis/google-cloud-go/commit/6593c0d62d48751c857bce3d3f858127467a4489))

## [1.66.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.65.0...spanner/v1.66.0) (2024-08-07)


### Features

* **spanner/admin/database:** Add support for Cloud Spanner Incremental Backups ([d949cc0](https://github.com/googleapis/google-cloud-go/commit/d949cc0e5d44af62154d9d5fd393f25a852f93ed))
* **spanner:** Add support of multiplexed session support in writeAtleastOnce mutations ([#10646](https://github.com/googleapis/google-cloud-go/issues/10646)) ([54009ea](https://github.com/googleapis/google-cloud-go/commit/54009eab1c3b11a28531ad9e621917d01c9e5339))
* **spanner:** Add support of using multiplexed session with ReadOnlyTransactions ([#10269](https://github.com/googleapis/google-cloud-go/issues/10269)) ([7797022](https://github.com/googleapis/google-cloud-go/commit/7797022e51d1ac07b8d919c421a8bfdf34a1d53c))

## [1.65.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.64.0...spanner/v1.65.0) (2024-07-29)


### Features

* **spanner/admin/database:** Add support for Cloud Spanner Scheduled Backups ([3b15f9d](https://github.com/googleapis/google-cloud-go/commit/3b15f9db9e0ee3bff3d8d5aafc82cdc2a31d60fc))
* **spanner:** Add RESOURCE_EXHAUSTED to retryable transaction codes ([#10412](https://github.com/googleapis/google-cloud-go/issues/10412)) ([29b52dc](https://github.com/googleapis/google-cloud-go/commit/29b52dc40f3d1a6ffe7fa40e6142d8035c0d95ee))


### Bug Fixes

* **spanner/test:** Bump google.golang.org/api@v0.187.0 ([8fa9e39](https://github.com/googleapis/google-cloud-go/commit/8fa9e398e512fd8533fd49060371e61b5725a85b))
* **spanner/test:** Bump google.golang.org/grpc@v1.64.1 ([8ecc4e9](https://github.com/googleapis/google-cloud-go/commit/8ecc4e9622e5bbe9b90384d5848ab816027226c5))
* **spanner/test:** Update dependencies ([257c40b](https://github.com/googleapis/google-cloud-go/commit/257c40bd6d7e59730017cf32bda8823d7a232758))
* **spanner:** Bump google.golang.org/api@v0.187.0 ([8fa9e39](https://github.com/googleapis/google-cloud-go/commit/8fa9e398e512fd8533fd49060371e61b5725a85b))
* **spanner:** Bump google.golang.org/grpc@v1.64.1 ([8ecc4e9](https://github.com/googleapis/google-cloud-go/commit/8ecc4e9622e5bbe9b90384d5848ab816027226c5))
* **spanner:** Fix negative values for max_in_use_sessions metrics [#10449](https://github.com/googleapis/google-cloud-go/issues/10449) ([#10508](https://github.com/googleapis/google-cloud-go/issues/10508)) ([4e180f4](https://github.com/googleapis/google-cloud-go/commit/4e180f4539012eb6e3d1d2788e68b291ef7230c3))
* **spanner:** HealthCheck should not decrement num_in_use sessions ([#10480](https://github.com/googleapis/google-cloud-go/issues/10480)) ([9b2b47f](https://github.com/googleapis/google-cloud-go/commit/9b2b47f107153d624d56709d9a8e6a6b72c39447))
* **spanner:** Update dependencies ([257c40b](https://github.com/googleapis/google-cloud-go/commit/257c40bd6d7e59730017cf32bda8823d7a232758))

## [1.64.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.63.0...spanner/v1.64.0) (2024-06-29)


### Features

* **spanner:** Add field lock_hint in spanner.proto ([3df3c04](https://github.com/googleapis/google-cloud-go/commit/3df3c04f0dffad3fa2fe272eb7b2c263801b9ada))
* **spanner:** Add field order_by in spanner.proto ([3df3c04](https://github.com/googleapis/google-cloud-go/commit/3df3c04f0dffad3fa2fe272eb7b2c263801b9ada))
* **spanner:** Add LockHint feature ([#10382](https://github.com/googleapis/google-cloud-go/issues/10382)) ([64bdcb1](https://github.com/googleapis/google-cloud-go/commit/64bdcb1a6a462d41a62d3badea6814425e271f22))
* **spanner:** Add OrderBy feature ([#10289](https://github.com/googleapis/google-cloud-go/issues/10289)) ([07b8bd2](https://github.com/googleapis/google-cloud-go/commit/07b8bd2f5dc738e0293305dfc459c13632d5ea65))
* **spanner:** Add support of checking row not found errors from ReadRow and ReadRowUsingIndex ([#10405](https://github.com/googleapis/google-cloud-go/issues/10405)) ([5cb0c26](https://github.com/googleapis/google-cloud-go/commit/5cb0c26013eeb3bbe51174bee628a20c2ec775e0))


### Bug Fixes

* **spanner:** Fix data-race caused by TrackSessionHandle ([#10321](https://github.com/googleapis/google-cloud-go/issues/10321)) ([23c5fff](https://github.com/googleapis/google-cloud-go/commit/23c5fffd06bcde408db50a981c015921cd4ecf0e)), refs [#10320](https://github.com/googleapis/google-cloud-go/issues/10320)
* **spanner:** Fix negative values for max_in_use_sessions metrics ([#10449](https://github.com/googleapis/google-cloud-go/issues/10449)) ([a1e198a](https://github.com/googleapis/google-cloud-go/commit/a1e198a9b18bd2f92c3438e4f609412047f8ccf4))
* **spanner:** Prevent possible panic for Session not found errors ([#10386](https://github.com/googleapis/google-cloud-go/issues/10386)) ([ba9711f](https://github.com/googleapis/google-cloud-go/commit/ba9711f87ec871153ae00cfd0827bce17c31ee9c)), refs [#10385](https://github.com/googleapis/google-cloud-go/issues/10385)

## [1.63.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.62.0...spanner/v1.63.0) (2024-05-24)


### Features

* **spanner:** Fix schema naming ([#10194](https://github.com/googleapis/google-cloud-go/issues/10194)) ([215e0c8](https://github.com/googleapis/google-cloud-go/commit/215e0c8125ea05246c834984bde1ca698c7dde4c))
* **spanner:** Update go mod to use latest grpc lib ([#10218](https://github.com/googleapis/google-cloud-go/issues/10218)) ([adf91f9](https://github.com/googleapis/google-cloud-go/commit/adf91f9fd37faa39ec7c6f9200273220f65d2a82))

## [1.62.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.61.0...spanner/v1.62.0) (2024-05-15)


### Features

* **spanner/admin/database:** Add support for multi region encryption config ([3e25053](https://github.com/googleapis/google-cloud-go/commit/3e250530567ee81ed4f51a3856c5940dbec35289))
* **spanner/executor:** Add QueryCancellationAction message in executor protos ([292e812](https://github.com/googleapis/google-cloud-go/commit/292e81231b957ae7ac243b47b8926564cee35920))
* **spanner:** Add `RESOURCE_EXHAUSTED` to the list of retryable error codes ([1d757c6](https://github.com/googleapis/google-cloud-go/commit/1d757c66478963d6cbbef13fee939632c742759c))
* **spanner:** Add support for Proto Columns ([#9315](https://github.com/googleapis/google-cloud-go/issues/9315)) ([3ffbbbe](https://github.com/googleapis/google-cloud-go/commit/3ffbbbe50225684f4211c6dbe3ca25acb3d02b8e))


### Bug Fixes

* **spanner:** Add ARRAY keywords to keywords ([#10079](https://github.com/googleapis/google-cloud-go/issues/10079)) ([8e675cd](https://github.com/googleapis/google-cloud-go/commit/8e675cd0ccf12c6912209aa5c56092db3716c40d))
* **spanner:** Handle unused errors ([#10067](https://github.com/googleapis/google-cloud-go/issues/10067)) ([a0c097c](https://github.com/googleapis/google-cloud-go/commit/a0c097c724b609cfa428e69f89075f02a3782a7b))
* **spanner:** Remove json-iterator dependency ([#10099](https://github.com/googleapis/google-cloud-go/issues/10099)) ([3917cca](https://github.com/googleapis/google-cloud-go/commit/3917ccac57c403b3b4d07514ac10a66a86e298c0)), refs [#9380](https://github.com/googleapis/google-cloud-go/issues/9380)
* **spanner:** Update staleness bound ([#10118](https://github.com/googleapis/google-cloud-go/issues/10118)) ([c07f1e4](https://github.com/googleapis/google-cloud-go/commit/c07f1e47c06387b696abb1edbfa339b391ec1fd5))

## [1.61.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.60.0...spanner/v1.61.0) (2024-04-30)


### Features

* **spanner/admin/instance:** Adding `EXPECTED_FULFILLMENT_PERIOD` to the indicate instance creation times (with `FULFILLMENT_PERIOD_NORMAL` or `FULFILLMENT_PERIOD_EXTENDED` ENUM) with the extended instance creation time triggered by On-Demand Capacity... ([#9693](https://github.com/googleapis/google-cloud-go/issues/9693)) ([aa93790](https://github.com/googleapis/google-cloud-go/commit/aa93790132ba830b4c97d217ef02764e2fb1b8ea))
* **spanner/executor:** Add SessionPoolOptions, SpannerOptions protos in executor protos ([2cdc40a](https://github.com/googleapis/google-cloud-go/commit/2cdc40a0b4288f5ab5f2b2b8f5c1d6453a9c81ec))
* **spanner:** Add support for change streams transaction exclusion option ([#9779](https://github.com/googleapis/google-cloud-go/issues/9779)) ([979ce94](https://github.com/googleapis/google-cloud-go/commit/979ce94758442b1224a78a4f3b1f5d592ab51660))
* **spanner:** Support MultiEndpoint ([#9565](https://github.com/googleapis/google-cloud-go/issues/9565)) ([0ac0d26](https://github.com/googleapis/google-cloud-go/commit/0ac0d265abedf946b05294ef874a892b2c5d6067))


### Bug Fixes

* **spanner/test/opentelemetry/test:** Bump x/net to v0.24.0 ([ba31ed5](https://github.com/googleapis/google-cloud-go/commit/ba31ed5fda2c9664f2e1cf972469295e63deb5b4))
* **spanner:** Bump x/net to v0.24.0 ([ba31ed5](https://github.com/googleapis/google-cloud-go/commit/ba31ed5fda2c9664f2e1cf972469295e63deb5b4))
* **spanner:** Fix uint8 conversion ([9221c7f](https://github.com/googleapis/google-cloud-go/commit/9221c7fa12cef9d5fb7ddc92f41f1d6204971c7b))

## [1.60.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.59.0...spanner/v1.60.0) (2024-03-19)


### Features

* **spanner:** Allow attempt direct path xds via env var ([e4b663c](https://github.com/googleapis/google-cloud-go/commit/e4b663cdcb6e010c5a8ac791e5624407aaa191b3))

## [1.59.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.58.0...spanner/v1.59.0) (2024-03-13)


### Features

* **spanner/spansql:** Support Table rename & Table synonym ([#9275](https://github.com/googleapis/google-cloud-go/issues/9275)) ([9b97ce7](https://github.com/googleapis/google-cloud-go/commit/9b97ce75d36980fdaa06f15b0398b7b65e0d6082))
* **spanner:** Add support of float32 type ([#9525](https://github.com/googleapis/google-cloud-go/issues/9525)) ([87d7ea9](https://github.com/googleapis/google-cloud-go/commit/87d7ea97787a56b18506b53e9b26d037f92759ca))


### Bug Fixes

* **spanner:** Add JSON_PARSE_ARRAY to funcNames slice ([#9557](https://github.com/googleapis/google-cloud-go/issues/9557)) ([f799597](https://github.com/googleapis/google-cloud-go/commit/f79959722352ead48bfb3efb3001fddd3a56db65))

## [1.58.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.57.0...spanner/v1.58.0) (2024-03-06)


### Features

* **spanner/admin/instance:** Add instance partition support to spanner instance proto ([ae1f547](https://github.com/googleapis/google-cloud-go/commit/ae1f5472bff1b476c3fd58e590ec135185446daf))
* **spanner:** Add field for multiplexed session in spanner.proto ([a86aa8e](https://github.com/googleapis/google-cloud-go/commit/a86aa8e962b77d152ee6cdd433ad94967150ef21))
* **spanner:** SelectAll struct spanner tag annotation match should be case-insensitive ([#9460](https://github.com/googleapis/google-cloud-go/issues/9460)) ([6cd6a73](https://github.com/googleapis/google-cloud-go/commit/6cd6a73be87a261729d3b6b45f3d28be93c3fdb3))
* **spanner:** Update TransactionOptions to include new option exclude_txn_from_change_streams ([0195fe9](https://github.com/googleapis/google-cloud-go/commit/0195fe9292274ff9d86c71079a8e96ed2e5f9331))

## [1.57.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.56.0...spanner/v1.57.0) (2024-02-13)


### Features

* **spanner:** Add OpenTelemetry implementation ([#9254](https://github.com/googleapis/google-cloud-go/issues/9254)) ([fc51cc2](https://github.com/googleapis/google-cloud-go/commit/fc51cc2ac71e8fb0b3e381379dc343630ed441e7))
* **spanner:** Support max_commit_delay in Spanner transactions ([#9299](https://github.com/googleapis/google-cloud-go/issues/9299)) ([a8078f0](https://github.com/googleapis/google-cloud-go/commit/a8078f0b841281bd439c548db9d303f6b5ce54e6))


### Bug Fixes

* **spanner:** Enable universe domain resolution options ([fd1d569](https://github.com/googleapis/google-cloud-go/commit/fd1d56930fa8a747be35a224611f4797b8aeb698))
* **spanner:** Internal test package should import local version ([#9416](https://github.com/googleapis/google-cloud-go/issues/9416)) ([f377281](https://github.com/googleapis/google-cloud-go/commit/f377281a73553af9a9a2bee2181efe2e354e1c68))
* **spanner:** SelectAll struct fields match should be case-insensitive ([#9417](https://github.com/googleapis/google-cloud-go/issues/9417)) ([7ff5356](https://github.com/googleapis/google-cloud-go/commit/7ff535672b868e6cba54abdf5dd92b9199e4d1d4))
* **spanner:** Support time.Time and other custom types using SelectAll ([#9382](https://github.com/googleapis/google-cloud-go/issues/9382)) ([dc21234](https://github.com/googleapis/google-cloud-go/commit/dc21234268b08a4a21b2b3a1ed9ed74d65a289f0))


### Documentation

* **spanner:** Update the comment regarding eligible SQL shapes for PartitionQuery ([e60a6ba](https://github.com/googleapis/google-cloud-go/commit/e60a6ba01acf2ef2e8d12e23ed5c6e876edeb1b7))

## [1.56.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.55.0...spanner/v1.56.0) (2024-01-30)


### Features

* **spanner/admin/database:** Add proto descriptors for proto and enum types in create/update/get database ddl requests ([97d62c7](https://github.com/googleapis/google-cloud-go/commit/97d62c7a6a305c47670ea9c147edc444f4bf8620))
* **spanner/spansql:** Add support for CREATE VIEW with SQL SECURITY DEFINER ([#8754](https://github.com/googleapis/google-cloud-go/issues/8754)) ([5f156e8](https://github.com/googleapis/google-cloud-go/commit/5f156e8c88f4729f569ee5b4ac9378dda3907997))
* **spanner:** Add FLOAT32 enum to TypeCode ([97d62c7](https://github.com/googleapis/google-cloud-go/commit/97d62c7a6a305c47670ea9c147edc444f4bf8620))
* **spanner:** Add max_commit_delay API ([af2f8b4](https://github.com/googleapis/google-cloud-go/commit/af2f8b4f3401c0b12dadb2c504aa0f902aee76de))
* **spanner:** Add proto and enum types ([00b9900](https://github.com/googleapis/google-cloud-go/commit/00b990061592a20a181e61faa6964b45205b76a7))
* **spanner:** Add SelectAll method to decode from Spanner iterator.Rows to golang struct ([#9206](https://github.com/googleapis/google-cloud-go/issues/9206)) ([802088f](https://github.com/googleapis/google-cloud-go/commit/802088f1322752bb9ce9bab1315c3fed6b3a99aa))

## [1.55.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.54.0...spanner/v1.55.0) (2024-01-08)


### Features

* **spanner:** Add directed reads feature ([#7668](https://github.com/googleapis/google-cloud-go/issues/7668)) ([a42604a](https://github.com/googleapis/google-cloud-go/commit/a42604a3a6ea90c38a2ff90d036a79fd070174fd))

## [1.54.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.53.1...spanner/v1.54.0) (2023-12-14)


### Features

* **spanner/executor:** Add autoscaling config in the instance to support autoscaling in systests ([29effe6](https://github.com/googleapis/google-cloud-go/commit/29effe600e16f24a127a1422ec04263c4f7a600a))
* **spanner:** New clients ([#9127](https://github.com/googleapis/google-cloud-go/issues/9127)) ([2c97389](https://github.com/googleapis/google-cloud-go/commit/2c97389ddacdfc140a06f74498cc2753bb040a4d))


### Bug Fixes

* **spanner:** Use json.Number for decoding unknown values from spanner ([#9054](https://github.com/googleapis/google-cloud-go/issues/9054)) ([40d1392](https://github.com/googleapis/google-cloud-go/commit/40d139297bd484408c63c9d6ad1d7035d9673c1c))

## [1.53.1](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.53.0...spanner/v1.53.1) (2023-12-01)


### Bug Fixes

* **spanner:** Handle nil error when cleaning up long running session ([#9052](https://github.com/googleapis/google-cloud-go/issues/9052)) ([a93bc26](https://github.com/googleapis/google-cloud-go/commit/a93bc2696bf9ae60aae93af0e8c4911b58514d31))
* **spanner:** MarshalJSON function caused errors for certain values ([#9063](https://github.com/googleapis/google-cloud-go/issues/9063)) ([afe7c98](https://github.com/googleapis/google-cloud-go/commit/afe7c98036c198995075530d4228f1f4ae3f1222))

## [1.53.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.52.0...spanner/v1.53.0) (2023-11-15)


### Features

* **spanner:** Enable long running transaction clean up ([#8969](https://github.com/googleapis/google-cloud-go/issues/8969)) ([5d181bb](https://github.com/googleapis/google-cloud-go/commit/5d181bb3a6fea55b8d9d596213516129006bdae2))

## [1.52.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.51.0...spanner/v1.52.0) (2023-11-14)


### Features

* **spanner:** Add directed_read_option in spanner.proto ([#8950](https://github.com/googleapis/google-cloud-go/issues/8950)) ([24e410e](https://github.com/googleapis/google-cloud-go/commit/24e410efbb6add2d33ecfb6ad98b67dc8894e578))
* **spanner:** Add DML, DQL, Mutation, Txn Actions and Utility methods for executor framework ([#8976](https://github.com/googleapis/google-cloud-go/issues/8976)) ([ca76671](https://github.com/googleapis/google-cloud-go/commit/ca7667194007394bdcade8058fa84c1fe19c06b1))
* **spanner:** Add lastUseTime property to session ([#8942](https://github.com/googleapis/google-cloud-go/issues/8942)) ([b560cfc](https://github.com/googleapis/google-cloud-go/commit/b560cfcf967ff6dec0cd6ac4b13045470945f30b))
* **spanner:** Add method ([#8945](https://github.com/googleapis/google-cloud-go/issues/8945)) ([411a51e](https://github.com/googleapis/google-cloud-go/commit/411a51e320fe21ffe830cdaa6bb4e4d77f7a996b))
* **spanner:** Add methods to return Row fields ([#8953](https://github.com/googleapis/google-cloud-go/issues/8953)) ([e22e70f](https://github.com/googleapis/google-cloud-go/commit/e22e70f44f83aab4f8b89af28fcd24216d2e740e))
* **spanner:** Add PG.OID type cod annotation ([#8749](https://github.com/googleapis/google-cloud-go/issues/8749)) ([ffb0dda](https://github.com/googleapis/google-cloud-go/commit/ffb0ddabf3d9822ba8120cabaf25515fd32e9615))
* **spanner:** Admin, Batch, Partition actions for executor framework ([#8932](https://github.com/googleapis/google-cloud-go/issues/8932)) ([b2db89e](https://github.com/googleapis/google-cloud-go/commit/b2db89e03a125cde31a7ea86eecc3fbb08ebd281))
* **spanner:** Auto-generated executor framework proto changes ([#8713](https://github.com/googleapis/google-cloud-go/issues/8713)) ([2ca939c](https://github.com/googleapis/google-cloud-go/commit/2ca939cba4bc240f2bfca7d5683708fd3a94fd74))
* **spanner:** BatchWrite ([#8652](https://github.com/googleapis/google-cloud-go/issues/8652)) ([507d232](https://github.com/googleapis/google-cloud-go/commit/507d232cdb09bd941ebfe800bdd4bfc020346f5d))
* **spanner:** Executor framework server and worker proxy ([#8714](https://github.com/googleapis/google-cloud-go/issues/8714)) ([6b931ee](https://github.com/googleapis/google-cloud-go/commit/6b931eefb9aa4a18758788167bdcf9e2fad1d7b9))
* **spanner:** Fix falkiness ([#8977](https://github.com/googleapis/google-cloud-go/issues/8977)) ([ca8d3cb](https://github.com/googleapis/google-cloud-go/commit/ca8d3cbf80f7fc2f47beb53b95138040c83097db))
* **spanner:** Long running transaction clean up - disabled ([#8177](https://github.com/googleapis/google-cloud-go/issues/8177)) ([461d11e](https://github.com/googleapis/google-cloud-go/commit/461d11e913414e9de822e5f1acdf19c8f3f953d5))
* **spanner:** Update code for session leaks cleanup ([#8978](https://github.com/googleapis/google-cloud-go/issues/8978)) ([cc83515](https://github.com/googleapis/google-cloud-go/commit/cc83515d0c837c8b1596a97b6f09d519a0f75f72))


### Bug Fixes

* **spanner:** Bump google.golang.org/api to v0.149.0 ([8d2ab9f](https://github.com/googleapis/google-cloud-go/commit/8d2ab9f320a86c1c0fab90513fc05861561d0880))
* **spanner:** Expose Mutations field in MutationGroup ([#8923](https://github.com/googleapis/google-cloud-go/issues/8923)) ([42180cf](https://github.com/googleapis/google-cloud-go/commit/42180cf1134885188270f75126a65fa71b03c033))
* **spanner:** Update grpc-go to v1.56.3 ([343cea8](https://github.com/googleapis/google-cloud-go/commit/343cea8c43b1e31ae21ad50ad31d3b0b60143f8c))
* **spanner:** Update grpc-go to v1.59.0 ([81a97b0](https://github.com/googleapis/google-cloud-go/commit/81a97b06cb28b25432e4ece595c55a9857e960b7))


### Documentation

* **spanner:** Updated comment formatting ([24e410e](https://github.com/googleapis/google-cloud-go/commit/24e410efbb6add2d33ecfb6ad98b67dc8894e578))

## [1.51.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.50.0...spanner/v1.51.0) (2023-10-17)


### Features

* **spanner/admin/instance:** Add autoscaling config to the instance proto ([#8701](https://github.com/googleapis/google-cloud-go/issues/8701)) ([56ce871](https://github.com/googleapis/google-cloud-go/commit/56ce87195320634b07ae0b012efcc5f2b3813fb0))


### Bug Fixes

* **spanner:** Update golang.org/x/net to v0.17.0 ([174da47](https://github.com/googleapis/google-cloud-go/commit/174da47254fefb12921bbfc65b7829a453af6f5d))

## [1.50.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.49.0...spanner/v1.50.0) (2023-10-03)


### Features

* **spanner/spansql:** Add support for aggregate functions ([#8498](https://github.com/googleapis/google-cloud-go/issues/8498)) ([d440d75](https://github.com/googleapis/google-cloud-go/commit/d440d75f19286653afe4bc81a5f2efcfc4fa152c))
* **spanner/spansql:** Add support for bit functions, sequence functions and GENERATE_UUID ([#8482](https://github.com/googleapis/google-cloud-go/issues/8482)) ([3789882](https://github.com/googleapis/google-cloud-go/commit/3789882c8b30a6d3100a56c1dcc8844952605637))
* **spanner/spansql:** Add support for SEQUENCE statements ([#8481](https://github.com/googleapis/google-cloud-go/issues/8481)) ([ccd0205](https://github.com/googleapis/google-cloud-go/commit/ccd020598921f1b5550587c95b4ceddf580705bb))
* **spanner:** Add BatchWrite API ([02a899c](https://github.com/googleapis/google-cloud-go/commit/02a899c95eb9660128506cf94525c5a75bedb308))
* **spanner:** Allow non-default service accounts ([#8488](https://github.com/googleapis/google-cloud-go/issues/8488)) ([c90dd00](https://github.com/googleapis/google-cloud-go/commit/c90dd00350fa018dbc5f0af5aabce80e80be0b90))

## [1.49.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.48.0...spanner/v1.49.0) (2023-08-24)


### Features

* **spanner/spannertest:** Support INSERT DML ([#7820](https://github.com/googleapis/google-cloud-go/issues/7820)) ([3dda7b2](https://github.com/googleapis/google-cloud-go/commit/3dda7b27ec536637d8ebaa20937fc8019c930481))


### Bug Fixes

* **spanner:** Transaction was started in a different session ([#8467](https://github.com/googleapis/google-cloud-go/issues/8467)) ([6c21558](https://github.com/googleapis/google-cloud-go/commit/6c21558f75628908a70de79c62aff2851e756e7b))

## [1.48.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.47.0...spanner/v1.48.0) (2023-08-18)


### Features

* **spanner/spansql:** Add complete set of math functions ([#8246](https://github.com/googleapis/google-cloud-go/issues/8246)) ([d7a238e](https://github.com/googleapis/google-cloud-go/commit/d7a238eca2a9b08e968cea57edc3708694673e22))
* **spanner/spansql:** Add support for foreign key actions ([#8296](https://github.com/googleapis/google-cloud-go/issues/8296)) ([d78b851](https://github.com/googleapis/google-cloud-go/commit/d78b8513b13a9a2c04b8097f0d89f85dcfd73797))
* **spanner/spansql:** Add support for IF NOT EXISTS and IF EXISTS clause ([#8245](https://github.com/googleapis/google-cloud-go/issues/8245)) ([96840ab](https://github.com/googleapis/google-cloud-go/commit/96840ab1232bbdb788e37f81cf113ee0f1b4e8e7))
* **spanner:** Add integration tests for Bit Reversed Sequences ([#7924](https://github.com/googleapis/google-cloud-go/issues/7924)) ([9b6e7c6](https://github.com/googleapis/google-cloud-go/commit/9b6e7c6061dc69683d7f558faed7f4249da5b7cb))


### Bug Fixes

* **spanner:** Reset buffer after abort on first SQL statement ([#8440](https://github.com/googleapis/google-cloud-go/issues/8440)) ([d980b42](https://github.com/googleapis/google-cloud-go/commit/d980b42f33968ef25061be50e18038d73b0503b6))
* **spanner:** REST query UpdateMask bug ([df52820](https://github.com/googleapis/google-cloud-go/commit/df52820b0e7721954809a8aa8700b93c5662dc9b))

## [1.47.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.46.0...spanner/v1.47.0) (2023-06-20)


### Features

* **spanner/admin/database:** Add DdlStatementActionInfo and add actions to UpdateDatabaseDdlMetadata ([01eff11](https://github.com/googleapis/google-cloud-go/commit/01eff11eedb3edde69cc33db23e26be6a7e42f10))
* **spanner:** Add databoost property for batch transactions ([#8152](https://github.com/googleapis/google-cloud-go/issues/8152)) ([fc49c78](https://github.com/googleapis/google-cloud-go/commit/fc49c78c9503c6dd4cbcba8c15e887415a744136))
* **spanner:** Add tests for database roles in PG dialect ([#7898](https://github.com/googleapis/google-cloud-go/issues/7898)) ([dc84649](https://github.com/googleapis/google-cloud-go/commit/dc84649c546fe09b0bab09991086c156bd78cb3f))
* **spanner:** Enable client to server compression ([#7899](https://github.com/googleapis/google-cloud-go/issues/7899)) ([3a047d2](https://github.com/googleapis/google-cloud-go/commit/3a047d2a449b0316a9000539ec9797e47cdd5c91))
* **spanner:** Update all direct dependencies ([b340d03](https://github.com/googleapis/google-cloud-go/commit/b340d030f2b52a4ce48846ce63984b28583abde6))


### Bug Fixes

* **spanner:** Fix TestRetryInfoTransactionOutcomeUnknownError flaky behaviour ([#7959](https://github.com/googleapis/google-cloud-go/issues/7959)) ([f037795](https://github.com/googleapis/google-cloud-go/commit/f03779538f949fb4ad93d5247d3c6b3e5b21091a))

## [1.46.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.45.1...spanner/v1.46.0) (2023-05-12)


### Features

* **spanner/admin/database:** Add support for UpdateDatabase in Cloud Spanner ([#7917](https://github.com/googleapis/google-cloud-go/issues/7917)) ([83870f5](https://github.com/googleapis/google-cloud-go/commit/83870f55035d6692e22264b209e39e07fe2823b9))
* **spanner:** Make leader aware routing default enabled for supported RPC requests. ([#7912](https://github.com/googleapis/google-cloud-go/issues/7912)) ([d0d3755](https://github.com/googleapis/google-cloud-go/commit/d0d37550911f37e09ea9204d0648fb64ff3204ff))


### Bug Fixes

* **spanner:** Update grpc to v1.55.0 ([1147ce0](https://github.com/googleapis/google-cloud-go/commit/1147ce02a990276ca4f8ab7a1ab65c14da4450ef))

## [1.45.1](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.45.0...spanner/v1.45.1) (2023-04-21)


### Bug Fixes

* **spanner/spannertest:** Send transaction id in result metadata ([#7809](https://github.com/googleapis/google-cloud-go/issues/7809)) ([e3bbd5f](https://github.com/googleapis/google-cloud-go/commit/e3bbd5f10b3922ab2eb50cb39daccd7bc1891892))
* **spanner:** Context timeout should be wrapped correctly ([#7744](https://github.com/googleapis/google-cloud-go/issues/7744)) ([f8e22f6](https://github.com/googleapis/google-cloud-go/commit/f8e22f6cbba10fc262e87b4d06d5c1289d877503))

## [1.45.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.44.0...spanner/v1.45.0) (2023-04-10)


### Features

* **spanner/spansql:** Add support for missing DDL syntax for ALTER CHANGE STREAM ([#7429](https://github.com/googleapis/google-cloud-go/issues/7429)) ([d34fe02](https://github.com/googleapis/google-cloud-go/commit/d34fe02cfa31520f88dedbd41bbc887e8faa857f))
* **spanner/spansql:** Support fine-grained access control DDL syntax ([#6691](https://github.com/googleapis/google-cloud-go/issues/6691)) ([a7edf6b](https://github.com/googleapis/google-cloud-go/commit/a7edf6b5c62d02b7d5199fc83d435f6a37a8eac5))
* **spanner/spansql:** Support grant/revoke view, change stream, table function ([#7533](https://github.com/googleapis/google-cloud-go/issues/7533)) ([9c61215](https://github.com/googleapis/google-cloud-go/commit/9c612159647d540e694ec9e84cab5cdd1c94d2b8))
* **spanner:** Add x-goog-spanner-route-to-leader header to Spanner RPC contexts for RW/PDML transactions. ([#7500](https://github.com/googleapis/google-cloud-go/issues/7500)) ([fcab05f](https://github.com/googleapis/google-cloud-go/commit/fcab05faa5026896af76b762eed5b7b6b2e7ee07))
* **spanner:** Adding new fields for Serverless analytics ([69067f8](https://github.com/googleapis/google-cloud-go/commit/69067f8c0075099a84dd9d40e438711881710784))
* **spanner:** Enable custom decoding for list value ([#7463](https://github.com/googleapis/google-cloud-go/issues/7463)) ([3aeadcd](https://github.com/googleapis/google-cloud-go/commit/3aeadcd97eaf2707c2f6e288c8b72ef29f49a185))
* **spanner:** Update iam and longrunning deps ([91a1f78](https://github.com/googleapis/google-cloud-go/commit/91a1f784a109da70f63b96414bba8a9b4254cddd))


### Bug Fixes

* **spanner/spansql:** Fix SQL for CREATE CHANGE STREAM TableName; case ([#7514](https://github.com/googleapis/google-cloud-go/issues/7514)) ([fc5fd86](https://github.com/googleapis/google-cloud-go/commit/fc5fd8652771aeca73e7a28ee68134155a5a9499))
* **spanner:** Correcting the proto field Id for field data_boost_enabled ([00fff3a](https://github.com/googleapis/google-cloud-go/commit/00fff3a58bed31274ab39af575876dab91d708c9))

## [1.44.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.43.0...spanner/v1.44.0) (2023-02-01)


### Features

* **spanner/spansql:** Add support for ALTER INDEX statement ([#7287](https://github.com/googleapis/google-cloud-go/issues/7287)) ([fbe1bd4](https://github.com/googleapis/google-cloud-go/commit/fbe1bd4d0806302a48ff4a5822867757893a5f2d))
* **spanner/spansql:** Add support for managing the optimizer statistics package ([#7283](https://github.com/googleapis/google-cloud-go/issues/7283)) ([e528221](https://github.com/googleapis/google-cloud-go/commit/e52822139e2821a11873c2d6af85a5fea07700e8))
* **spanner:** Add support for Optimistic Concurrency Control ([#7332](https://github.com/googleapis/google-cloud-go/issues/7332)) ([48ba16f](https://github.com/googleapis/google-cloud-go/commit/48ba16f3a09893a3527a22838ad1e9ff829da15b))

## [1.43.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.42.0...spanner/v1.43.0) (2023-01-19)


### Features

* **spanner/spansql:** Add support for change stream value_capture_type option ([#7201](https://github.com/googleapis/google-cloud-go/issues/7201)) ([27b3398](https://github.com/googleapis/google-cloud-go/commit/27b33988f078779c2d641f776a11b2095a5ccc51))
* **spanner/spansql:** Support `default_leader` database option ([#7187](https://github.com/googleapis/google-cloud-go/issues/7187)) ([88adaa2](https://github.com/googleapis/google-cloud-go/commit/88adaa216832467560c19e61528b5ce5f1e5ff76))
* **spanner:** Add REST client ([06a54a1](https://github.com/googleapis/google-cloud-go/commit/06a54a16a5866cce966547c51e203b9e09a25bc0))
* **spanner:** Inline begin transaction for ReadWriteTransactions ([#7149](https://github.com/googleapis/google-cloud-go/issues/7149)) ([2ce3606](https://github.com/googleapis/google-cloud-go/commit/2ce360644439a386aeaad7df5f47541667bd621b))


### Bug Fixes

* **spanner:** Fix integration tests data race ([#7229](https://github.com/googleapis/google-cloud-go/issues/7229)) ([a741024](https://github.com/googleapis/google-cloud-go/commit/a741024abd6fb1f073831503c2717b2a44226a59))

## [1.42.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.41.0...spanner/v1.42.0) (2022-12-14)


### Features

* **spanner:** Add database roles ([#5701](https://github.com/googleapis/google-cloud-go/issues/5701)) ([6bb95ef](https://github.com/googleapis/google-cloud-go/commit/6bb95efb7997692a52c321e787e633a5045b21f8))
* **spanner:** Rewrite signatures and type in terms of new location ([620e6d8](https://github.com/googleapis/google-cloud-go/commit/620e6d828ad8641663ae351bfccfe46281e817ad))


### Bug Fixes

* **spanner:** Fallback to check grpc error message if ResourceType is nil for checking sessionNotFound errors ([#7163](https://github.com/googleapis/google-cloud-go/issues/7163)) ([2552e09](https://github.com/googleapis/google-cloud-go/commit/2552e092cff01e0d6b80fefaa7877f77e36db6be))

## [1.41.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.40.0...spanner/v1.41.0) (2022-12-01)


### Features

* **spanner:** Start generating proto stubs ([#7030](https://github.com/googleapis/google-cloud-go/issues/7030)) ([41f446f](https://github.com/googleapis/google-cloud-go/commit/41f446f891a17c97278879f2207fd58996fd038c))

## [1.40.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.39.0...spanner/v1.40.0) (2022-11-03)


### Features

* **spanner/spansql:** Add support for interval arg of some date/timestamp functions ([#6950](https://github.com/googleapis/google-cloud-go/issues/6950)) ([1ce0f7d](https://github.com/googleapis/google-cloud-go/commit/1ce0f7d38778068fd1d9a171377067739f4ea8d6))
* **spanner:** Configurable logger ([#6958](https://github.com/googleapis/google-cloud-go/issues/6958)) ([bd85442](https://github.com/googleapis/google-cloud-go/commit/bd85442bc6fb8c18d1a7c6d73850d220c3973c46)), refs [#6957](https://github.com/googleapis/google-cloud-go/issues/6957)
* **spanner:** PG JSONB support ([#6874](https://github.com/googleapis/google-cloud-go/issues/6874)) ([5b14658](https://github.com/googleapis/google-cloud-go/commit/5b146587939ccc3403945c756cbf68e6f2d41fda))
* **spanner:** Update result_set.proto to return undeclared parameters in ExecuteSql API ([de4e16a](https://github.com/googleapis/google-cloud-go/commit/de4e16a498354ea7271f5b396f7cb2bb430052aa))
* **spanner:** Update transaction.proto to include different lock modes ([caf4afa](https://github.com/googleapis/google-cloud-go/commit/caf4afa139ad7b38b6df3e3b17b8357c81e1fd6c))

## [1.39.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.38.0...spanner/v1.39.0) (2022-09-21)


### Features

* **spanner/admin/database:** Add custom instance config operations ([ec1a190](https://github.com/googleapis/google-cloud-go/commit/ec1a190abbc4436fcaeaa1421c7d9df624042752))
* **spanner/admin/instance:** Add custom instance config operations ([ef2b0b1](https://github.com/googleapis/google-cloud-go/commit/ef2b0b1d4de9beb9005537ae48d7d8e1c0f23b98))
* **spanner/spannersql:** Add backticks when name contains a hypen ([#6621](https://github.com/googleapis/google-cloud-go/issues/6621)) ([e88ca66](https://github.com/googleapis/google-cloud-go/commit/e88ca66ca950e15d9011322dbfca3c88ccceb0ec))
* **spanner/spansql:** Add support for create, alter and drop change … ([#6669](https://github.com/googleapis/google-cloud-go/issues/6669)) ([cc4620a](https://github.com/googleapis/google-cloud-go/commit/cc4620a5ee3a9129a4cdd48d90d4060ba0bbcd58))
* **spanner:** Retry spanner transactions and mutations when RST_STREAM error ([#6699](https://github.com/googleapis/google-cloud-go/issues/6699)) ([1b56cd0](https://github.com/googleapis/google-cloud-go/commit/1b56cd0ec31bc32362259fc722907e092bae081a))


### Bug Fixes

* **spanner/admin/database:** Revert add custom instance config operations (change broke client libraries; reverting before any are released) ([ec1a190](https://github.com/googleapis/google-cloud-go/commit/ec1a190abbc4436fcaeaa1421c7d9df624042752))
* **spanner:** Destroy session when client is closing ([#6700](https://github.com/googleapis/google-cloud-go/issues/6700)) ([a1ce541](https://github.com/googleapis/google-cloud-go/commit/a1ce5410f1e0f4d68dae0ddc790518e9978faf0c))
* **spanner:** Spanner sessions will be cleaned up from the backend ([#6679](https://github.com/googleapis/google-cloud-go/issues/6679)) ([c27097e](https://github.com/googleapis/google-cloud-go/commit/c27097e236abeb8439a67ad9b716d05c001aea2e))

## [1.38.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.37.0...spanner/v1.38.0) (2022-09-03)


### Features

* **spanner/spannertest:** add support for adding and dropping Foreign Keys ([#6608](https://github.com/googleapis/google-cloud-go/issues/6608)) ([ccd3614](https://github.com/googleapis/google-cloud-go/commit/ccd3614f6edbaf3d7d202feb4df220f244550a78))
* **spanner/spansql:** add support for coalesce expressions ([#6461](https://github.com/googleapis/google-cloud-go/issues/6461)) ([bff16a7](https://github.com/googleapis/google-cloud-go/commit/bff16a783c1fd4d7e888d4ee3b5420c1bbf10da1))
* **spanner:** Adds auto-generated CL for googleapis for jsonb ([3bc37e2](https://github.com/googleapis/google-cloud-go/commit/3bc37e28626df5f7ec37b00c0c2f0bfb91c30495))


### Bug Fixes

* **spanner:** pass userAgent to cloud spanner requests ([#6598](https://github.com/googleapis/google-cloud-go/issues/6598)) ([59d162b](https://github.com/googleapis/google-cloud-go/commit/59d162bdfcbe00a060a52930be7185f00e8df2c1))

## [1.37.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.36.0...spanner/v1.37.0) (2022-08-28)


### Features

* **spanner/admin/database:** Add ListDatabaseRoles API to support role based access control ([1ffeb95](https://github.com/googleapis/google-cloud-go/commit/1ffeb9557bf1f18cc131aff40ec7e0e15a9f4ead))
* **spanner/spansql:** add support for nullif expressions ([#6423](https://github.com/googleapis/google-cloud-go/issues/6423)) ([5b7bfeb](https://github.com/googleapis/google-cloud-go/commit/5b7bfebcd4a0fd3cbe355d9d290e6b5101810b7e))
* **spanner:** install grpc rls and xds by default ([#6007](https://github.com/googleapis/google-cloud-go/issues/6007)) ([70d562f](https://github.com/googleapis/google-cloud-go/commit/70d562f25738052e833a46daf6ff7fa1f4a0a746))
* **spanner:** set client wide ReadOptions, ApplyOptions, and TransactionOptions ([#6486](https://github.com/googleapis/google-cloud-go/issues/6486)) ([757f1ca](https://github.com/googleapis/google-cloud-go/commit/757f1cac7a765fe2e7ead872d07eb24baad61c28))


### Bug Fixes

* **spanner/admin/database:** target new spanner db admin service config ([1d6fbcc](https://github.com/googleapis/google-cloud-go/commit/1d6fbcc6406e2063201ef5a98de560bf32f7fb73))

## [1.36.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.35.0...spanner/v1.36.0) (2022-07-23)


### Features

* **spanner/spansql:** add support for IFNULL expressions ([#6389](https://github.com/googleapis/google-cloud-go/issues/6389)) ([09e96ce](https://github.com/googleapis/google-cloud-go/commit/09e96ce1076df4b41d45c3676b7506b318da6b9c))
* **spanner/spansql:** support for parsing a DML file ([#6349](https://github.com/googleapis/google-cloud-go/issues/6349)) ([267a9bb](https://github.com/googleapis/google-cloud-go/commit/267a9bbec55ee8fe885354efc8db8a61a17a8374))

## [1.35.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.34.1...spanner/v1.35.0) (2022-07-19)


### Features

* **spanner/admin/instance:** Adding two new fields for Instance create_time and update_time ([8a1ad06](https://github.com/googleapis/google-cloud-go/commit/8a1ad06572a65afa91a0a77a85b849e766876671))
* **spanner/spansql:** add support for if expressions ([#6341](https://github.com/googleapis/google-cloud-go/issues/6341)) ([56c858c](https://github.com/googleapis/google-cloud-go/commit/56c858cebd683e45d1dd5ab8ae98ef9bfd767edc))


### Bug Fixes

* **spanner:** fix pool.numInUse exceeding MaxOpened ([#6344](https://github.com/googleapis/google-cloud-go/issues/6344)) ([882b325](https://github.com/googleapis/google-cloud-go/commit/882b32593e8c7bff8369b1ff9259c7b408fad661))

## [1.34.1](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.34.0...spanner/v1.34.1) (2022-07-06)


### Bug Fixes

* **spanner/spansql:** Add tests for INSERT parsing ([#6303](https://github.com/googleapis/google-cloud-go/issues/6303)) ([0d19fb5](https://github.com/googleapis/google-cloud-go/commit/0d19fb5d60554b9a90fac52918f784e6c3e13918))

## [1.34.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.33.0...spanner/v1.34.0) (2022-06-17)


### Features

* **spanner/spansql:** add a support for parsing INSERT statement ([#6148](https://github.com/googleapis/google-cloud-go/issues/6148)) ([c6185cf](https://github.com/googleapis/google-cloud-go/commit/c6185cffc7f23741ac4a230aadee74b3def85ced))
* **spanner:** add Session creator role docs: clarify transaction semantics ([4134941](https://github.com/googleapis/google-cloud-go/commit/41349411e601f57dc6d9e246f1748fd86d17bb15))

## [1.33.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.32.0...spanner/v1.33.0) (2022-05-28)


### Bug Fixes

* **spanner/spansql:** fix invalid timestamp literal formats ([#6077](https://github.com/googleapis/google-cloud-go/issues/6077)) ([6ab8bed](https://github.com/googleapis/google-cloud-go/commit/6ab8bed93a978e00a6c195d8cb4d574ca6db27c3))


### Miscellaneous Chores

* **spanner:** release 1.33.0 ([#6104](https://github.com/googleapis/google-cloud-go/issues/6104)) ([54bc54e](https://github.com/googleapis/google-cloud-go/commit/54bc54e9bbdc22e2bbfd9f315885f95987e2c3f2))

## [1.32.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.31.0...spanner/v1.32.0) (2022-05-09)


### Features

* **spanner/spansql:** support DEFAULT keyword ([#5932](https://github.com/googleapis/google-cloud-go/issues/5932)) ([49c19a9](https://github.com/googleapis/google-cloud-go/commit/49c19a956031fa889d024bd57fa34681bc79e743))
* **spanner/spansql:** support JSON literals ([#5968](https://github.com/googleapis/google-cloud-go/issues/5968)) ([b500120](https://github.com/googleapis/google-cloud-go/commit/b500120f3cc5c7b5717f6525a24de72fd317ba66))
* **spanner:** enable row.ToStructLenient to work with STRUCT data type ([#5944](https://github.com/googleapis/google-cloud-go/issues/5944)) ([bca8d50](https://github.com/googleapis/google-cloud-go/commit/bca8d50533115b9995f7b4a63d5d1f9abaf6a753))

## [1.31.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.30.1...spanner/v1.31.0) (2022-04-08)


### Features

* **spanner/spansql:** support case expression ([#5836](https://github.com/googleapis/google-cloud-go/issues/5836)) ([3ffdd62](https://github.com/googleapis/google-cloud-go/commit/3ffdd626e72c6472f337a423b9702baf0c298185))


### Bug Fixes

* **spanner/spannertest:** Improve DDL application delay cancellation. ([#5874](https://github.com/googleapis/google-cloud-go/issues/5874)) ([08f1e72](https://github.com/googleapis/google-cloud-go/commit/08f1e72dbf2ef5a06425f71500d061af246bd490))

### [1.30.1](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.30.0...spanner/v1.30.1) (2022-03-28)


### Bug Fixes

* **spanner:** early unlock of session pool lock during dumping the tracked session handles to avoid deadlock ([#5777](https://github.com/googleapis/google-cloud-go/issues/5777)) ([b007836](https://github.com/googleapis/google-cloud-go/commit/b0078362865159b87bc34c1a7f990a361f1cafcf))

## [1.30.0](https://github.com/googleapis/google-cloud-go/compare/spanner/v1.29.0...spanner/v1.30.0) (2022-03-04)


### Features

* **spanner:** add better version metadata to calls ([#5515](https://github.com/googleapis/google-cloud-go/issues/5515)) ([dcab7c4](https://github.com/googleapis/google-cloud-go/commit/dcab7c4a98ebecfef1f75ec5bddfd7782b28a7c5)), refs [#2749](https://github.com/googleapis/google-cloud-go/issues/2749)
* **spanner:** add file for tracking version ([17b36ea](https://github.com/googleapis/google-cloud-go/commit/17b36ead42a96b1a01105122074e65164357519e))
* **spanner:** add support of PGNumeric with integration tests for PG dialect ([#5700](https://github.com/googleapis/google-cloud-go/issues/5700)) ([f7e02e1](https://github.com/googleapis/google-cloud-go/commit/f7e02e11064d14c04eca18ab808e8fe5194ac355))
* **spanner:** set versionClient to module version ([55f0d92](https://github.com/googleapis/google-cloud-go/commit/55f0d92bf112f14b024b4ab0076c9875a17423c9))

### Bug Fixes

* **spanner/spansql:** support GROUP BY without an aggregation function ([#5717](https://github.com/googleapis/google-cloud-go/issues/5717)) ([c819ee9](https://github.com/googleapis/google-cloud-go/commit/c819ee9ad4695afa31eddcb4bf87764762555cd5))


### Miscellaneous Chores

* **spanner:** release 1.30.0 ([#5715](https://github.com/googleapis/google-cloud-go/issues/5715)) ([a19d182](https://github.com/googleapis/google-cloud-go/commit/a19d182dab5476cf01e719c751e94a73a98c6c4a))

## [1.29.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.28.0...spanner/v1.29.0) (2022-01-06)


### ⚠ BREAKING CHANGES

* **spanner:** fix data race in spanner integration tests (#5276)

### Features

* **spanner/spansql:** support EXTRACT ([#5218](https://www.github.com/googleapis/google-cloud-go/issues/5218)) ([81b7c85](https://www.github.com/googleapis/google-cloud-go/commit/81b7c85a8993a36557ea4eb4ec0c47d1f93c4960))
* **spanner/spansql:** support MOD function ([#5231](https://www.github.com/googleapis/google-cloud-go/issues/5231)) ([0a81fbc](https://www.github.com/googleapis/google-cloud-go/commit/0a81fbc0171af7e828f3e606cbe7b3905ac32213))
* **spanner:** add google-c2p dependence ([5343756](https://www.github.com/googleapis/google-cloud-go/commit/534375668b5b81bae5ef750c96856bef027f9d1e))
* **spanner:** Add ReadRowWithOptions method ([#5240](https://www.github.com/googleapis/google-cloud-go/issues/5240)) ([c276428](https://www.github.com/googleapis/google-cloud-go/commit/c276428bca79702245d422849af6472bb2e74171))
* **spanner:** Adding GFE Latency and Header Missing Count Metrics ([#5199](https://www.github.com/googleapis/google-cloud-go/issues/5199)) ([3d8a9ea](https://www.github.com/googleapis/google-cloud-go/commit/3d8a9ead8d73a4f38524a424a98362c32f56954b))


### Bug Fixes

* **spanner:** result from unmarshal of string and spanner.NullString type from json should be consistent. ([#5263](https://www.github.com/googleapis/google-cloud-go/issues/5263)) ([7eaaa47](https://www.github.com/googleapis/google-cloud-go/commit/7eaaa470fda5dc7cd1ff041d6a898e35fb54920e))


### Tests

* **spanner:** fix data race in spanner integration tests ([#5276](https://www.github.com/googleapis/google-cloud-go/issues/5276)) ([22df34b](https://www.github.com/googleapis/google-cloud-go/commit/22df34b8e7d0d003b3eeaf1c069aee58f30a8dfe))


### Miscellaneous Chores

* **spanner:** release 1.29.0 ([#5292](https://www.github.com/googleapis/google-cloud-go/issues/5292)) ([9f0b900](https://www.github.com/googleapis/google-cloud-go/commit/9f0b9003686d26c66a10c3b54e67b59c2a6327ff))

## [1.28.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.27.0...spanner/v1.28.0) (2021-12-03)


### Features

* **spanner/spannertest:** support JSON_VALUE function ([#5173](https://www.github.com/googleapis/google-cloud-go/issues/5173)) ([ac98735](https://www.github.com/googleapis/google-cloud-go/commit/ac98735cb1adc9384c5b2caeb9aac938db275bf7))
* **spanner/spansql:** support CAST and SAFE_CAST ([#5057](https://www.github.com/googleapis/google-cloud-go/issues/5057)) ([54cbf4c](https://www.github.com/googleapis/google-cloud-go/commit/54cbf4c0a0305e680b213f84487110dfeaf8e7e1))
* **spanner:** add ToStructLenient method to decode to struct fields with no error return with un-matched row's column with struct's exported fields. ([#5153](https://www.github.com/googleapis/google-cloud-go/issues/5153)) ([899ffbf](https://www.github.com/googleapis/google-cloud-go/commit/899ffbf8ce42b1597ca3cd59bfd9f042054b8ae2))

## [1.27.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.26.0...spanner/v1.27.0) (2021-10-19)


### Features

* **spanner:** implement valuer and scanner interfaces ([#4936](https://www.github.com/googleapis/google-cloud-go/issues/4936)) ([4537b45](https://www.github.com/googleapis/google-cloud-go/commit/4537b45d2611ce480abfb5d186b59e7258ec872c))

## [1.26.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.25.0...spanner/v1.26.0) (2021-10-11)


### Features

* **spanner/spannertest:** implement RowDeletionPolicy in spannertest ([#4961](https://www.github.com/googleapis/google-cloud-go/issues/4961)) ([7800a33](https://www.github.com/googleapis/google-cloud-go/commit/7800a3303b97204a0573780786388437bbbf2673)), refs [#4782](https://www.github.com/googleapis/google-cloud-go/issues/4782)
* **spanner/spannertest:** Support generated columns ([#4742](https://www.github.com/googleapis/google-cloud-go/issues/4742)) ([324d11d](https://www.github.com/googleapis/google-cloud-go/commit/324d11d3c19ffbd77848c8e19c972b70ff5e9268))
* **spanner/spansql:** fill in missing hash functions ([#4808](https://www.github.com/googleapis/google-cloud-go/issues/4808)) ([37ee2d9](https://www.github.com/googleapis/google-cloud-go/commit/37ee2d95220efc1aaf0280d0aa2c01ae4b9d4c1b))
* **spanner/spansql:** support JSON data type ([#4959](https://www.github.com/googleapis/google-cloud-go/issues/4959)) ([e84e408](https://www.github.com/googleapis/google-cloud-go/commit/e84e40830752fc8bc0ccdd869fa7b8fd0c80f306))
* **spanner/spansql:** Support multiple joins in query ([#4743](https://www.github.com/googleapis/google-cloud-go/issues/4743)) ([81a308e](https://www.github.com/googleapis/google-cloud-go/commit/81a308e909a3ae97504a49fbc9982f7eeb6be80c))

## [1.25.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.24.1...spanner/v1.25.0) (2021-08-25)


### Features

* **spanner/spansql:** add support for STARTS_WITH function ([#4670](https://www.github.com/googleapis/google-cloud-go/issues/4670)) ([7a56af0](https://www.github.com/googleapis/google-cloud-go/commit/7a56af03d1505d9a29d1185a50e261c0e90fdb1a)), refs [#4661](https://www.github.com/googleapis/google-cloud-go/issues/4661)
* **spanner:** add support for JSON data type ([#4104](https://www.github.com/googleapis/google-cloud-go/issues/4104)) ([ade8ab1](https://www.github.com/googleapis/google-cloud-go/commit/ade8ab111315d84fa140ddde020387a78668dfa4))


### Bug Fixes

* **spanner/spannertest:** Fix the "LIKE" clause handling for prefix and suffix matches ([#4655](https://www.github.com/googleapis/google-cloud-go/issues/4655)) ([a2118f0](https://www.github.com/googleapis/google-cloud-go/commit/a2118f02fb03bfc50952699318f35c23dc234c41))
* **spanner:** invalid numeric should throw an error ([#3926](https://www.github.com/googleapis/google-cloud-go/issues/3926)) ([cde8697](https://www.github.com/googleapis/google-cloud-go/commit/cde8697be01f1ef57806275c0ddf54f87bb9a571))

### [1.24.1](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.24.0...spanner/v1.24.1) (2021-08-11)


### Bug Fixes

* **spanner/spansql:** only add comma after other option ([#4551](https://www.github.com/googleapis/google-cloud-go/issues/4551)) ([3ac1e00](https://www.github.com/googleapis/google-cloud-go/commit/3ac1e007163803d315dcf5db612fe003f6eab978))
* **spanner:** allow decoding null values to spanner.Decoder ([#4558](https://www.github.com/googleapis/google-cloud-go/issues/4558)) ([45ddaca](https://www.github.com/googleapis/google-cloud-go/commit/45ddaca606a372d9293bf2e2b3dc6d4398166c43)), refs [#4552](https://www.github.com/googleapis/google-cloud-go/issues/4552)

## [1.24.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.23.0...spanner/v1.24.0) (2021-07-29)


### Features

* **spanner/spansql:** add ROW DELETION POLICY parsing ([#4496](https://www.github.com/googleapis/google-cloud-go/issues/4496)) ([3d6c6c7](https://www.github.com/googleapis/google-cloud-go/commit/3d6c6c7873e1b75e8b492ede2e561411dc40536a))
* **spanner/spansql:** fix unstable SelectFromTable SQL ([#4473](https://www.github.com/googleapis/google-cloud-go/issues/4473)) ([39bc4ec](https://www.github.com/googleapis/google-cloud-go/commit/39bc4eca655d0180b18378c175d4a9a77fe1602f))
* **spanner/spansql:** support ALTER DATABASE ([#4403](https://www.github.com/googleapis/google-cloud-go/issues/4403)) ([1458dc9](https://www.github.com/googleapis/google-cloud-go/commit/1458dc9c21d98ffffb871943f178678cc3c21306))
* **spanner/spansql:** support table_hint_expr at from_clause on query_statement ([#4457](https://www.github.com/googleapis/google-cloud-go/issues/4457)) ([7047808](https://www.github.com/googleapis/google-cloud-go/commit/7047808794cf463c6a96d7b59ef5af3ed94fd7cf))
* **spanner:** add row.String() and refine error message for decoding a struct array ([#4431](https://www.github.com/googleapis/google-cloud-go/issues/4431)) ([f6258a4](https://www.github.com/googleapis/google-cloud-go/commit/f6258a47a4dfadc02dcdd75b53fd5f88c5dcca30))
* **spanner:** allow untyped nil values in parameterized queries ([#4482](https://www.github.com/googleapis/google-cloud-go/issues/4482)) ([c1ba18b](https://www.github.com/googleapis/google-cloud-go/commit/c1ba18b1b1fc45de6e959cc22a5c222cc80433ee))


### Bug Fixes

* **spanner/spansql:** fix DATE and TIMESTAMP parsing. ([#4480](https://www.github.com/googleapis/google-cloud-go/issues/4480)) ([dec7a67](https://www.github.com/googleapis/google-cloud-go/commit/dec7a67a3e980f6f5e0d170919da87e1bffe923f))

## [1.23.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.22.0...spanner/v1.23.0) (2021-07-08)


### Features

* **spanner/admin/database:** add leader_options to InstanceConfig and default_leader to Database ([7aa0e19](https://www.github.com/googleapis/google-cloud-go/commit/7aa0e195a5536dd060a1fca871bd3c6f946d935e))

## [1.22.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.21.0...spanner/v1.22.0) (2021-06-30)


### Features

* **spanner:** support request and transaction tags ([#4336](https://www.github.com/googleapis/google-cloud-go/issues/4336)) ([f08c73a](https://www.github.com/googleapis/google-cloud-go/commit/f08c73a75e2d2a8b9a0b184179346cb97c82e9e5))
* **spanner:** enable request options for batch read ([#4337](https://www.github.com/googleapis/google-cloud-go/issues/4337)) ([b9081c3](https://www.github.com/googleapis/google-cloud-go/commit/b9081c36ed6495a67f8e458ad884bdb8da5b7fbc))

## [1.21.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.20.0...spanner/v1.21.0) (2021-06-23)


### Miscellaneous Chores

* **spanner:** trigger a release for low cost instance ([#4264](https://www.github.com/googleapis/google-cloud-go/issues/4264)) ([24c4451](https://www.github.com/googleapis/google-cloud-go/commit/24c4451404cdf4a83cc7a35ee1911d654d2ba132))

## [1.20.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.19.0...spanner/v1.20.0) (2021-06-08)


### Features

* **spanner:** add the support of optimizer statistics package ([#2717](https://www.github.com/googleapis/google-cloud-go/issues/2717)) ([29c7247](https://www.github.com/googleapis/google-cloud-go/commit/29c724771f0b19849c76e62d4bc8e9342922bf75))

## [1.19.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.18.0...spanner/v1.19.0) (2021-06-03)


### Features

* **spanner/spannertest:** support multiple aggregations ([#3965](https://www.github.com/googleapis/google-cloud-go/issues/3965)) ([1265dc3](https://www.github.com/googleapis/google-cloud-go/commit/1265dc3289693f79fcb9c5785a424eb510a50007))
* **spanner/spansql:** case insensitive parsing of keywords and functions ([#4034](https://www.github.com/googleapis/google-cloud-go/issues/4034)) ([ddb09d2](https://www.github.com/googleapis/google-cloud-go/commit/ddb09d22a737deea0d0a9ab58cd5d337164bbbfe))
* **spanner:** add a database name getter to client ([#4190](https://www.github.com/googleapis/google-cloud-go/issues/4190)) ([7fce29a](https://www.github.com/googleapis/google-cloud-go/commit/7fce29af404f0623b483ca6d6f2af4c726105fa6))
* **spanner:** add custom instance config to tests ([#4194](https://www.github.com/googleapis/google-cloud-go/issues/4194)) ([e935345](https://www.github.com/googleapis/google-cloud-go/commit/e9353451237e658bde2e41b30e8270fbc5987b39))


### Bug Fixes

* **spanner:** add missing NUMERIC type to the doc for Row ([#4116](https://www.github.com/googleapis/google-cloud-go/issues/4116)) ([9a3b416](https://www.github.com/googleapis/google-cloud-go/commit/9a3b416221f3c8b3793837e2a459b1d7cd9c479f))
* **spanner:** indent code example for Encoder and Decoder ([#4128](https://www.github.com/googleapis/google-cloud-go/issues/4128)) ([7c1f48f](https://www.github.com/googleapis/google-cloud-go/commit/7c1f48f307284c26c10cd5787dbc94136a2a36a6))
* **spanner:** mark SessionPoolConfig.MaxBurst deprecated ([#4115](https://www.github.com/googleapis/google-cloud-go/issues/4115)) ([d60a686](https://www.github.com/googleapis/google-cloud-go/commit/d60a68649f85f1edfbd8f11673bb280813c2b771))

## [1.18.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.17.0...spanner/v1.18.0) (2021-04-29)


### Features

* **spanner/admin/database:** add `progress` field to `UpdateDatabaseDdlMetadata` ([9029071](https://www.github.com/googleapis/google-cloud-go/commit/90290710158cf63de918c2d790df48f55a23adc5))

## [1.17.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.16.0...spanner/v1.17.0) (2021-03-31)


### Features

* **spanner/admin/database:** add tagging request options ([2b02a03](https://www.github.com/googleapis/google-cloud-go/commit/2b02a03ff9f78884da5a8e7b64a336014c61bde7))
* **spanner:** add RPC Priority request options ([b5b4da6](https://www.github.com/googleapis/google-cloud-go/commit/b5b4da6952922440d03051f629f3166f731dfaa3))
* **spanner:** Add support for RPC priority ([#3341](https://www.github.com/googleapis/google-cloud-go/issues/3341)) ([88cf097](https://www.github.com/googleapis/google-cloud-go/commit/88cf097649f1cdf01cab531eabdff7fbf2be3f8f))

## [1.16.0](https://www.github.com/googleapis/google-cloud-go/compare/v1.15.0...v1.16.0) (2021-03-17)


### Features

* **spanner:** add `optimizer_statistics_package` field in `QueryOptions` ([18c88c4](https://www.github.com/googleapis/google-cloud-go/commit/18c88c437bd1741eaf5bf5911b9da6f6ea7cd75d))
* **spanner/admin/database:** add CMEK fields to backup and database ([16597fa](https://github.com/googleapis/google-cloud-go/commit/16597fa1ce549053c7183e8456e23f554a5501de))


### Bug Fixes

* **spanner/spansql:** fix parsing of NOT IN operator ([#3724](https://www.github.com/googleapis/google-cloud-go/issues/3724)) ([7636478](https://www.github.com/googleapis/google-cloud-go/commit/76364784d82073b80929ae60fd42da34c8050820))

## [1.15.0](https://www.github.com/googleapis/google-cloud-go/compare/v1.14.1...v1.15.0) (2021-02-24)


### Features

* **spanner/admin/database:** add CMEK fields to backup and database ([47037ed](https://www.github.com/googleapis/google-cloud-go/commit/47037ed33cd36edfff4ba7c4a4ea332140d5e67b))
* **spanner/admin/database:** add CMEK fields to backup and database ([16597fa](https://www.github.com/googleapis/google-cloud-go/commit/16597fa1ce549053c7183e8456e23f554a5501de))


### Bug Fixes

* **spanner:** parallelize session deletion when closing pool ([#3701](https://www.github.com/googleapis/google-cloud-go/issues/3701)) ([75ac7d2](https://www.github.com/googleapis/google-cloud-go/commit/75ac7d2506e706869ae41cf186b0c873b146e926)), refs [#3685](https://www.github.com/googleapis/google-cloud-go/issues/3685)

### [1.14.1](https://www.github.com/googleapis/google-cloud-go/compare/v1.14.0...v1.14.1) (2021-02-09)


### Bug Fixes

* **spanner:** restore removed scopes ([#3684](https://www.github.com/googleapis/google-cloud-go/issues/3684)) ([232d3a1](https://www.github.com/googleapis/google-cloud-go/commit/232d3a17bdadb92864592351a335ec920a68f9bf))

## [1.14.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.13.0...v1.14.0) (2021-02-09)


### Features

* **spanner/admin/database:** adds PITR fields to backup and database ([0959f27](https://www.github.com/googleapis/google-cloud-go/commit/0959f27e85efe94d39437ceef0ff62ddceb8e7a7))
* **spanner/spannertest:** restructure column alteration implementation ([#3616](https://www.github.com/googleapis/google-cloud-go/issues/3616)) ([176400b](https://www.github.com/googleapis/google-cloud-go/commit/176400be9ab485fb343b8994bc49ac2291d8eea9))
* **spanner/spansql:** add complete set of array functions ([#3633](https://www.github.com/googleapis/google-cloud-go/issues/3633)) ([13d50b9](https://www.github.com/googleapis/google-cloud-go/commit/13d50b93cc8348c54641b594371a96ecdb1bcabc))
* **spanner/spansql:** add complete set of string functions ([#3625](https://www.github.com/googleapis/google-cloud-go/issues/3625)) ([34027ad](https://www.github.com/googleapis/google-cloud-go/commit/34027ada6a718603be2987b4084ce5e0ead6413c))
* **spanner:** add option for returning Spanner commit stats ([c7ecf0f](https://www.github.com/googleapis/google-cloud-go/commit/c7ecf0f3f454606b124e52d20af2545b2c68646f))
* **spanner:** add option for returning Spanner commit stats ([7bdebad](https://www.github.com/googleapis/google-cloud-go/commit/7bdebadbe06774c94ab745dfef4ce58ce40a5582))
* **spanner:** support CommitStats ([#3444](https://www.github.com/googleapis/google-cloud-go/issues/3444)) ([b7c3ca6](https://www.github.com/googleapis/google-cloud-go/commit/b7c3ca6c83cbdca95d734df8aa07c5ddb8ab3db0))


### Bug Fixes

* **spanner/spannertest:** support queries in ExecuteSql ([#3640](https://www.github.com/googleapis/google-cloud-go/issues/3640)) ([8eede84](https://www.github.com/googleapis/google-cloud-go/commit/8eede8411a5521f45a5c3f8091c42b3c5407ea90)), refs [#3639](https://www.github.com/googleapis/google-cloud-go/issues/3639)
* **spanner/spansql:** fix SelectFromJoin behavior ([#3571](https://www.github.com/googleapis/google-cloud-go/issues/3571)) ([e0887c7](https://www.github.com/googleapis/google-cloud-go/commit/e0887c762a4c58f29b3e5b49ee163a36a065463c))

## [1.13.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.12.0...v1.13.0) (2021-01-15)


### Features

* **spanner/spannertest:** implement ANY_VALUE aggregation function ([#3428](https://www.github.com/googleapis/google-cloud-go/issues/3428)) ([e16c3e9](https://www.github.com/googleapis/google-cloud-go/commit/e16c3e9b412762b85483f3831ee586a5e6631313))
* **spanner/spannertest:** implement FULL JOIN ([#3218](https://www.github.com/googleapis/google-cloud-go/issues/3218)) ([99f7212](https://www.github.com/googleapis/google-cloud-go/commit/99f7212bd70bb333c1aa1c7a57348b4dfd80d31b))
* **spanner/spannertest:** implement SELECT ... FROM UNNEST(...) ([#3431](https://www.github.com/googleapis/google-cloud-go/issues/3431)) ([deb466f](https://www.github.com/googleapis/google-cloud-go/commit/deb466f497a1e6df78fcad57c3b90b1a4ccd93b4))
* **spanner/spannertest:** support array literals ([#3438](https://www.github.com/googleapis/google-cloud-go/issues/3438)) ([69e0110](https://www.github.com/googleapis/google-cloud-go/commit/69e0110f4977035cd1a705c3034c3ba96cadf36f))
* **spanner/spannertest:** support AVG aggregation function ([#3286](https://www.github.com/googleapis/google-cloud-go/issues/3286)) ([4788415](https://www.github.com/googleapis/google-cloud-go/commit/4788415c908f58c1cc08c951f1a7f17cdaf35aa2))
* **spanner/spannertest:** support Not Null constraint ([#3491](https://www.github.com/googleapis/google-cloud-go/issues/3491)) ([c36aa07](https://www.github.com/googleapis/google-cloud-go/commit/c36aa0785e798b9339d540e691850ca3c474a288))
* **spanner/spannertest:** support UPDATE DML ([#3201](https://www.github.com/googleapis/google-cloud-go/issues/3201)) ([1dec6f6](https://www.github.com/googleapis/google-cloud-go/commit/1dec6f6a31768a3f70bfec7274828301c22ea10b))
* **spanner/spansql:** define structures and parse UPDATE DML statements ([#3192](https://www.github.com/googleapis/google-cloud-go/issues/3192)) ([23b6904](https://www.github.com/googleapis/google-cloud-go/commit/23b69042c58489df512703259f54d075ba0c0722))
* **spanner/spansql:** support DATE and TIMESTAMP literals ([#3557](https://www.github.com/googleapis/google-cloud-go/issues/3557)) ([1961930](https://www.github.com/googleapis/google-cloud-go/commit/196193034a15f84dc3d3c27901990e8be77fca85))
* **spanner/spansql:** support for parsing generated columns ([#3373](https://www.github.com/googleapis/google-cloud-go/issues/3373)) ([9b1d06f](https://www.github.com/googleapis/google-cloud-go/commit/9b1d06fc90a4c07899c641a893dba0b47a1cead9))
* **spanner/spansql:** support NUMERIC data type ([#3411](https://www.github.com/googleapis/google-cloud-go/issues/3411)) ([1bc65d9](https://www.github.com/googleapis/google-cloud-go/commit/1bc65d9124ba22db5bec4c71b6378c27dfc04724))
* **spanner:** Add a DirectPath fallback integration test ([#3487](https://www.github.com/googleapis/google-cloud-go/issues/3487)) ([de821c5](https://www.github.com/googleapis/google-cloud-go/commit/de821c59fb81e9946216d205162b59de8b5ce71c))
* **spanner:** attempt DirectPath by default ([#3516](https://www.github.com/googleapis/google-cloud-go/issues/3516)) ([bbc61ed](https://www.github.com/googleapis/google-cloud-go/commit/bbc61ed368453b28aaf5bed627ca2499a3591f63))
* **spanner:** include User agent ([#3465](https://www.github.com/googleapis/google-cloud-go/issues/3465)) ([4e1ef1b](https://www.github.com/googleapis/google-cloud-go/commit/4e1ef1b3fb536ef950249cdee02cc0b6c2b56e86))
* **spanner:** run E2E test over DirectPath ([#3466](https://www.github.com/googleapis/google-cloud-go/issues/3466)) ([18e3a4f](https://www.github.com/googleapis/google-cloud-go/commit/18e3a4fe2a0c59c6295db2d85c7893ac51688083))
* **spanner:** support NUMERIC in mutations ([#3328](https://www.github.com/googleapis/google-cloud-go/issues/3328)) ([fa90737](https://www.github.com/googleapis/google-cloud-go/commit/fa90737a2adbe0cefbaba4aa1046a6efbba2a0e9))


### Bug Fixes

* **spanner:** fix session leak ([#3461](https://www.github.com/googleapis/google-cloud-go/issues/3461)) ([11fb917](https://www.github.com/googleapis/google-cloud-go/commit/11fb91711db5b941995737980cef7b48b611fefd)), refs [#3460](https://www.github.com/googleapis/google-cloud-go/issues/3460)

## [1.12.0](https://www.github.com/googleapis/google-cloud-go/compare/spanner/v1.11.0...v1.12.0) (2020-11-10)


### Features

* **spanner:** add metadata to RowIterator ([#3050](https://www.github.com/googleapis/google-cloud-go/issues/3050)) ([9a2289c](https://www.github.com/googleapis/google-cloud-go/commit/9a2289c3a38492bc2e84e0f4000c68a8718f5c11)), closes [#1805](https://www.github.com/googleapis/google-cloud-go/issues/1805)
* **spanner:** export ToSpannerError ([#3133](https://www.github.com/googleapis/google-cloud-go/issues/3133)) ([b951d8b](https://www.github.com/googleapis/google-cloud-go/commit/b951d8bd194b76da0a8bf2ce7cf85b546d2e051c)), closes [#3122](https://www.github.com/googleapis/google-cloud-go/issues/3122)
* **spanner:** support rw-transaction with options ([#3058](https://www.github.com/googleapis/google-cloud-go/issues/3058)) ([5130694](https://www.github.com/googleapis/google-cloud-go/commit/51306948eef9d26cff70453efc3eb500ddef9117))
* **spanner/spannertest:** make SELECT list aliases visible to ORDER BY ([#3054](https://www.github.com/googleapis/google-cloud-go/issues/3054)) ([7d2d83e](https://www.github.com/googleapis/google-cloud-go/commit/7d2d83ee1cce58d4014d5570bc599bcef1ed9c22)), closes [#3043](https://www.github.com/googleapis/google-cloud-go/issues/3043)

## v1.11.0

* Features:
  - feat(spanner): add KeySetFromKeys function (#2837)
* Misc:
  - test(spanner): check for Aborted error (#3039)
  - test(spanner): fix potential race condition in TestRsdBlockingStates (#3017)
  - test(spanner): compare data instead of struct (#3013)
  - test(spanner): fix flaky oc_test.go (#2838)
  - docs(spanner): document NULL value (#2885)
* spansql/spannertest:
  - Support JOINs (all but FULL JOIN) (#2936, #2924, #2896, #3042, #3037, #2995, #2945, #2931)
  - feat(spanner/spansql): parse CHECK constraints (#3046)
  - fix(spanner/spansql): fix parsing of unary minus and plus (#2997)
  - fix(spanner/spansql): fix parsing of adjacent inline and leading comments (#2851)
  - fix(spanner/spannertest): fix ORDER BY combined with SELECT aliases (#3043)
  - fix(spanner/spannertest): generate query output columns in construction order (#2990)
  - fix(spanner/spannertest): correct handling of NULL AND FALSE (#2991)
  - fix(spanner/spannertest): correct handling of tri-state boolean expression evaluation (#2983)
  - fix(spanner/spannertest): fix handling of NULL with LIKE operator (#2982)
  - test(spanner/spannertest): migrate most test code to integration_test.go (#2977)
  - test(spanner/spansql): add fuzz target for ParseQuery (#2909)
  - doc(spanner/spannertest): document the implementation (#2996)
  - perf(spanner/spannertest): speed up no-wait DDL changes (#2994)
  - perf(spanner/spansql): make fewer allocations during SQL (#2969)
* Backward Incompatible Changes
  - chore(spanner/spansql): use ID type for identifiers throughout (#2889)
  - chore(spanner/spansql): restructure FROM, TABLESAMPLE (#2888)

## v1.10.0

* feat(spanner): add support for NUMERIC data type (#2415)
* feat(spanner): add custom type support to spanner.Key (#2748)
* feat(spanner/spannertest): add support for bool parameter types (#2674)
* fix(spanner): update PDML to take sessions from pool (#2736)
* spanner/spansql: update docs on TableAlteration, ColumnAlteration (#2825)
* spanner/spannertest: support dropping columns (#2823)
* spanner/spannertest: implement GetDatabase (#2802)
* spanner/spannertest: fix aggregation in query evaluation for empty inputs (#2803)

## v1.9.0

* Features:
  - feat(spanner): support custom field type (#2614)
* Bugfixes:
  - fix(spanner): call ctx.cancel after stats have been recorded (#2728)
  - fix(spanner): retry session not found for read (#2724)
  - fix(spanner): specify credentials with SPANNER_EMULATOR_HOST (#2701)
  - fix(spanner): update pdml to retry EOS internal error (#2678)
* Misc:
  - test(spanner): unskip tests for emulator (#2675)
* spansql/spannertest:
  - spanner/spansql: restructure types and parsing for column options (#2656)
  - spanner/spannertest: return error for Read with no keys (#2655)

## v1.8.0

* Features:
  - feat(spanner): support of client-level custom retry settings (#2599)
  - feat(spanner): add a statement-based way to run read-write transaction. (#2545)
* Bugfixes:
  - fix(spanner): set 'gccl' to the request header. (#2609)
  - fix(spanner): add the missing resource prefix (#2605)
  - fix(spanner): fix the upgrade of protobuf. (#2583)
  - fix(spanner): do not copy protobuf messages by value. (#2581)
  - fix(spanner): fix the required resource prefix. (#2580)
  - fix(spanner): add extra field to ignore with cmp (#2577)
  - fix(spanner): remove appengine-specific numChannels. (#2513)
* Misc:
  - test(spanner): log warning instead of fail for stress test (#2559)
  - test(spanner): fix failed TestRsdBlockingStates test (#2597)
  - chore(spanner): cleanup mockserver and mockclient (#2414)

## v1.7.0

* Retry:
  - Only retry certain types of internal errors. (#2460)
* Tracing/metrics:
  - Never sample `ping()` trace spans (#2520)
  - Add oc tests for session pool metrics. (#2416)
* Encoding:
  - Allow encoding struct with custom types to mutation (#2529)
* spannertest:
  - Fix evaluation on IN (#2479)
  - Support MIN/MAX aggregation functions (#2411)
* Misc:
  - Fix TestClient_WithGRPCConnectionPoolAndNumChannels_Misconfigured test (#2539)
  - Cleanup backoff files and rename a variable (#2526)
  - Fix TestIntegration_DML test to return err from tx (#2509)
  - Unskip tests for emulator 0.8.0. (#2494)
  - Fix TestIntegration_StartBackupOperation test. (#2418)
  - Fix flakiness in TestIntegration_BatchDML_Error
  - Unskip TestIntegration_BatchDML and TestIntegration_BatchDML_TwoStatements
    for emulator by checking the existence of status.
  - Fix TestStressSessionPool test by taking lock while getting sessions from
    hc.

## v1.6.0

* Sessions:
  - Increase the number of sessions in batches instead of one by one when
    additional sessions are needed. The step size is set to 25, which means
    that whenever the session pool needs at least one more session, it will
    create a batch of 25 sessions.
* Emulator:
  - Run integration tests against the emulator in Kokoro Presubmit.
* RPC retrying:
  - Retry CreateDatabase on retryable codes.
* spannertest:
  - Change internal representation of DATE/TIMESTAMP values.
* spansql:
  - Cleanly parse adjacent comment marker/terminator.
  - Support FROM aliases in SELECT statements.
* Misc:
  - Fix comparing errors in tests.
  - Fix flaky session pool test.
  - Increase timeout in TestIntegration_ReadOnlyTransaction.
  - Fix incorrect instance IDs when deleting instances in tests.
  - Clean up test instances.
  - Clearify docs on Aborted transaction.
  - Fix timeout+staleness bound for test
  - Remove the support for resource-based routing.
  - Fix TestTransaction_SessionNotFound test.

## v1.5.1

* Fix incorrect decreasing metrics, numReads and numWrites.
* Fix an issue that XXX fields/methods are internal to proto and may change
  at any time. XXX_Merge panics in proto v1.4.0. Use proto.Merge instead of
  XXX_Merge.
* spannertest: handle list parameters in RPC interfacea.

## v1.5.0

* Metrics
  - Instrument client library with adding OpenCensus metrics. This allows for
    better monitoring of the session pool.
* Session management
  - Switch the session keepalive method from GetSession to SELECT 1.
* Emulator
  - Use client hooks for admin clients running against an emulator. With
    this change, users can use SPANNER_EMULATOR_HOST for initializing admin
    clients when running against an emulator.
* spansql
  - Add space between constraint name and foreign key def.
* Misc
  - Fix segfault when a non-existent credentials file had been specified.
  - Fix cleaning up instances in integration tests.
  - Fix race condition in batch read-only transaction.
  - Fix the flaky TestLIFOTakeWriteSessionOrder test.
  - Fix ITs to order results in SELECT queries.
  - Fix the documentation of timestamp bounds.
  - Fix the regex issue in managing backups.

## v1.4.0

- Support managed backups. This includes the API methods for CreateBackup,
  GetBackup, UpdateBackup, DeleteBackup and others. Also includes a simple
  wrapper in DatabaseAdminClient to create a backup.
- Update the healthcheck interval. The default interval is updated to 50 mins.
  By default, the first healthcheck is scheduled between 10 and 55 mins and
  the subsequent healthchecks are between 45 and 55 mins. This update avoids
  overloading the backend service with frequent healthchecking.

## v1.3.0

* Query options:
  - Adds the support of providing query options (optimizer version) via
    three ways (precedence follows the order):
    `client-level < environment variables < query-level`. The environment
    variable is set by "SPANNER_OPTIMIZER_VERSION".
* Connection pooling:
  - Use the new connection pooling in gRPC. This change deprecates
    `ClientConfig.numChannels` and users should move to
    `WithGRPCConnectionPool(numChannels)` at their earliest convenience.
    Example:
    ```go
    // numChannels (deprecated):
    err, client := NewClientWithConfig(ctx, database, ClientConfig{NumChannels: 8})

    // gRPC connection pool:
    err, client := NewClientWithConfig(ctx, database, ClientConfig{}, option.WithGRPCConnectionPool(8))
    ```
* Error handling:
  - Do not rollback after failed commit.
  - Return TransactionOutcomeUnknownError if a DEADLINE_EXCEEDED or CANCELED
    error occurs while a COMMIT request is in flight.
* spansql:
  - Added support for IN expressions and OFFSET clauses.
  - Fixed parsing of table constraints.
  - Added support for foreign key constraints in ALTER TABLE and CREATE TABLE.
  - Added support for GROUP BY clauses.
* spannertest:
  - Added support for IN expressions and OFFSET clauses.
  - Added support for GROUP BY clauses.
  - Fixed data race in query execution.
  - No longer rejects reads specifying an index to use.
  - Return last commit timestamp as read timestamp when requested.
  - Evaluate add, subtract, multiply, divide, unary
    negation, unary not, bitwise and/xor/or operations, as well as reporting
    column types for expressions involving any possible arithmetic
    operator.arithmetic expressions.
  - Fixed handling of descending primary keys.
* Misc:
  - Change default healthcheck interval to 30 mins to reduce the GetSession
    calls made to the backend.
  - Add marshal/unmarshal json for nullable types to support NullString,
    NullInt64, NullFloat64, NullBool, NullTime, NullDate.
  - Use ResourceInfo to extract error.
  - Extract retry info from status.

## v1.2.1

- Fix session leakage for ApplyAtLeastOnce. Previously session handles where
  leaked whenever Commit() returned a non-abort, non-session-not-found error,
  due to a missing recycle() call.
- Fix error for WriteStruct with pointers. This fixes a specific check for
  encoding and decoding to pointer types.
- Fix a GRPCStatus issue that returns a Status that has Unknown code if the
  base error is nil. Now, it always returns a Status based on Code field of
  current error.

## v1.2.0

- Support tracking stacktrace of sessionPool.take() that allows the user
  to instruct the session pool to keep track of the stacktrace of each
  goroutine that checks out a session from the pool. This is disabled by
  default, but it can be enabled by setting
  `SessionPoolConfig.TrackSessionHandles: true`.
- Add resource-based routing that includes a step to retrieve the
  instance-specific endpoint before creating the session client when
  creating a new spanner client. This is disabled by default, but it can
  be enabled by setting `GOOGLE_CLOUD_SPANNER_ENABLE_RESOURCE_BASED_ROUTING`.
- Make logger configurable so that the Spanner client can now be configured to
  use a specific logger instead of the standard logger.
- Support encoding custom types that point back to supported basic types.
- Allow decoding Spanner values to custom types that point back to supported
  types.

## v1.1.0

- The String() method of NullString, NullTime and NullDate will now return
  an unquoted string instead of a quoted string. This is a BREAKING CHANGE.
  If you relied on the old behavior, please use fmt.Sprintf("%q", T).
- The Spanner client will now use the new BatchCreateSessions RPC to initialize
  the session pool. This will improve the startup time of clients that are
  initialized with a minimum number of sessions greater than zero
  (i.e. SessionPoolConfig.MinOpened>0).
- Spanner clients that are created with the NewClient method will now default
  to a minimum of 100 opened sessions in the pool
  (i.e. SessionPoolConfig.MinOpened=100). This will improve the performance
  of the first transaction/query that is executed by an application, as a
  session will normally not have to be created as part of the transaction.
  Spanner clients that are created with the NewClientWithConfig method are
  not affected by this change.
- Spanner clients that are created with the NewClient method will now default
  to a write sessions fraction of 0.2 in the pool
  (i.e. SessionPoolConfig.WriteSessions=0.2).
  Spanner clients that are created with the NewClientWithConfig method are
  not affected by this change.
- The session pool maintenance worker has been improved so it keeps better
  track of the actual number of sessions needed. It will now less often delete
  and re-create sessions. This can improve the overall performance of
  applications with a low transaction rate.

## v1.0.0

This is the first tag to carve out spanner as its own module. See:
https://github.com/golang/go/wiki/Modules#is-it-possible-to-add-a-module-to-a-multi-module-repository.
