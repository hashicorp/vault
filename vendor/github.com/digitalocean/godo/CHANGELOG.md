# Change Log

## [v1.130.0] - 2024-11-14

- #755 - @vsharma6855  - Add Missing Database Configs for Postgresql and MYSQL
- #754 - @blesswinsamuel - APPS-9858 Add method to obtain websocket URL to get console access into components

## [v1.129.0] - 2024-11-06

- #752 - @andrewsomething - Support maps in Stringify
- #749 - @loosla - [droplets]: add droplet backup policies
- #730 - @rak16 - DOCR-1201: Add new RegistriesService to support methods for multiple-registry open beta
- #748 - @andrewsomething - Support Droplet GPU information

## [v1.128.0] - 2024-10-24

- #746 - @blesswinsamuel - Add archive field to AppSpec to archive/restore apps
- #745 - @asaha2 - Add load balancer monitoring endpoints
- #744 - @asaha2 - Adjust delete dangerous
- #743 - @asaha2 - Introduce droplet autoscale godo methods
- #740 - @blesswinsamuel - Add maintenance field to AppSpec to enable/disable maintenance mode
- #739 - @markusthoemmes - Add protocol to AppSpec and pending to detect responses

## [v1.127.0] - 2024-10-18

- #737 - @loosla - [databases]: change Opensearch ism_history_max_docs type to int64 to …
- #735 - @loosla - [databases]: add a missing field to Opensearch advanced configuration
- #729 - @loosla - [databases]: add support for Opensearch advanced configuration

## [v1.126.0] - 2024-09-25

- #732 - @gottwald - DOKS: add custom CIDR fields
- #727 - @loosla - [databases]: add support for Kafka advanced configuration

## [v1.125.0] - 2024-09-17

- #726 - @loosla - [databases]: add support for MongoDB advanced configuration
- #724 - @andrewsomething - Bump go version to 1.22
- #723 - @jauderho - Update Go dependencies and remove replace statements

## [v1.124.0] - 2024-09-10

- #721 - @vsharma6855 - [DBAAS] | Add API endpoint for applying cluster patches

## [v1.123.0] - 2024-09-06

- #719 - @andrewsomething - apps: mark ListTiers and GetTier as deprecated

## [v1.122.0] - 2024-09-04

- #717 - @danaelhe - DB: Fix Logsink Attribute Types
- #716 - @bhardwajRahul - Databases: Add support for OpenSearch ACL

## [v1.121.0] - 2024-08-20

- #715 - @danaelhe - Databases: Bring back Logsink Support
- #710 - @bhardwajRahul - Update GODO to include new Openseach index crud changes
- #712 - @danaelhe - Database: Namespace logsink
- #711 - @danaelhe - Databases: Add Logsinks CRUD support

## [v1.120.0] - 2024-08-08

- #708 - @markusthoemmes - APPS-9201 Add `UpdateAllSourceVersions` parameter to update app calls
- #706 - @andrewsomething - database: Add Size to DatabaseReplica struct

## [v1.119.0] - 2024-07-24

- #704 - @ElanHasson - APPS-9133 - Add support for OPENSEARCH as a database engine option
- #703 - @dependabot[bot] - Bump github.com/hashicorp/go-retryablehttp from 0.7.4 to 0.7.7
- #699 - @ElanHasson - APPS-8790 Add support to App Platform Log Forwarding for an OpenSearch DBaaS cluster destination.

## [v1.118.0] - 2024-06-04

**Note**: This release contains features in closed beta (#700).

- #701 - @llDrLove - Rename control plane permission to control plane firewall
- #700 - @bbassingthwaite - Add ProxyProtocol to LoadBalancer HealthCheck

## [v1.117.0] - 2024-06-04

- #696 - @llDrLove - Support specifying control plane firewall rules when creating or updating DOKS clusters
- #697 - @asaha2 - Add support for lb internal network type
- #695 - @ElanHasson - APPS-8732 - Update documentation on App Platform OpenSearch endpoint structure.
- #692 - @ElanHasson - APPS-8732 - Add OpenSearch as a Log Destination for App Platform.

## [v1.116.0] - 2024-05-16

- #693 - @guptado - Introduce VPC peering methods

## [v1.115.0] - 2024-05-08

- #688 - @asaha2 - load balancers: support glb active-passive fail-over settings, currently in closed beta

## [v1.114.0] - 2024-04-12

- #686 - @greeshmapill - APPS-8386: Add comments to mark deprecation of unused instance size fields
- #685 - @jcodybaker - APPS-8711: container termination controls
- #682 - @dependabot[bot] - Bump golang.org/x/net from 0.17.0 to 0.23.0

## [v1.113.0] - 2024-04-12

- #679 - @bhardwajRahul - Enable ui_connection parameter for Opensearch
- #678 - @bhardwajRahul - Enable Opensearch option in Godo

## [v1.112.0] - 2024-04-08

- #672 - @dependabot[bot] - Bump google.golang.org/protobuf from 1.28.0 to 1.33.0
- #675 - @bhardwajRahul - Add ListDatabaseEvents to Godo

## [v1.111.0] - 2024-04-02

- #674 - @asaha2 - load balancers: introduce glb settings in godo, currently in closed beta

## [v1.110.0] - 2024-03-14

- #667 - @dwilsondo - Include DBaaS metrics credential endpoint operations
- #670 - @guptado - [NETPROD-3583] Added name param in ListOption to get resource by name
- #671 - @greeshmapill - APPS-8383: Add deprecation intent and bandwidth allowance to app instance size spec

## [v1.109.0] - 2024-02-09

- #668 - @greeshmapill - APPS-8315: Update app instance size spec
- #665 - @jcodybaker - APPS-8263: methods for managing App Platform dev DBs
- #663 - @dwilsondo - Include replica connection info on DBaaS clusters & DBaaS PG pools
- #662 - @ddatta-do - load balancer : add regional network as new LB type

## [v1.108.0] - 2024-01-17

- #660 - @dweinshenker - Enable CRUD operations for replicas with storage_size_mib

## [v1.107.0] - 2023-12-07

- #658 - @markusthoemmes - APPS-8033 Add the RUN_RESTARTED log type
- #656 - @dweinshenker - Enhancement: add database user update
- #657 - @markusthoemmes - apps: Add registry_credentials field, GHCR registry type and the egress spec

## [v1.106.0] - 2023-11-14

- #654 - @dweinshenker - Remove unclean_leader_election_enable for topic configuration

## [v1.105.1] - 2023-11-07

- #652 - @andrewsomething - Retry on HTTP/2 internal errors.
- #648 - @alexandear - test: use fmt.Fprintf instead of fmt.Fprintf(fmt.Sprintf(...))
- #651 - @alexandear - test: Replace deprecated io/ioutil with io
- #647 - @alexandear - test: add missing error check

## [v1.105.0] - 2023-10-16

- #643 - @dweinshenker - Add support for scalable storage on database clusters
- #641 - @dweinshenker - Fix Kafka Partition Count
- #645 - @gregmankes - APPS-7325 - update app godo spec
- #642 - @dependabot[bot] - Bump golang.org/x/net from 0.7.0 to 0.17.0

## [v1.104.1] - 2023-10-10

* #640 - @andrewsomething - Drop required Go version to 1.20 and document policy.
* #640 - @andrewsomething - Fix library version.

## [v1.104.0] - 2023-10-10

- #637 - @mikesmithgh - chore: change uptime alert comparison type
- #638 - @markusthoemmes - APPS-7700 Add ability to specify digest for an image

## [v1.103.0] - 2023-10-03

- #635 - @andrewsomething - Bump github.com/stretchr/testify to v1.8.4
- #634 - @andrewsomething - Bump Go version to v1.21.0
- #632 - @danaelhe - Make Retrys by Default for NewFromToken()
- #633 - @dwilsondo - Add DBaaS engine Kafka
- #621 - @testwill - chore:  use fmt.Fprintf instead of fmt.Fprint(fmt.Sprintf(...))

## [v1.102.1] - 2023-08-17

- #629 - @andrewsomething - Provide a custom retryablehttp.ErrorHandler for more consistent returns using retries.

## [v1.102.0] - 2023-08-14

- #624 - @danaelhe - Update README.md with Retryable Info
- #626 - @andrewsomething - Allow configuring go-retryablehttp.Logger
- #625 - @andrewsomething - Export the HTTP client

## [v1.101.0] - 2023-08-09

- #619 - @danaelhe - Add retryablehttp Client Option

## [v1.100.0] - 2023-07-20

- #618 - @asaha - load balancers: introduce new type field
- #620 - @andrewsomething - account: add name field.

## [v1.99.0] - 2023-04-24

- #616 - @bentranter - Bump CI version for Go 1.20
- #615 - @bentranter - Remove beta support for tokens API
- #604 - @dvigueras - Add support for "Validate a Container Registry Name"
- #613 - @ibilalkayy - updated the README file by showing up the build status icon

## [v1.98.0] - 2023-03-09

- #608 - @anitgandhi - client: don't process body upon 204 response
- #607 - @gregmankes - add apps rewrites/redirects to app spec

## [v1.97.0] - 2023-02-10

- #601 - @jcodybaker - APPS-6813: update app platform - pending_deployment + timing
- #602 - @jcodybaker - Use App Platform active deployment for GetLogs if not specified

## [v1.96.0] - 2023-01-23

- #599 - @markpaulson - Adding PromoteReplicaToPrimary to client interface.

## [v1.95.0] - 2023-01-23

- #595 - @dweinshenker - Add UpgradeMajorVersion to godo

## [v1.94.0] - 2022-01-23

- #596 - @DMW2151 - DBAAS-3906: Include updatePool for DB Clusters
- #593 - @danaelhe - Add Uptime Checks and Alerts Support

## [v1.93.0] - 2022-12-15

- #591 - @andrewsomething - tokens: Add initial support for new API.

## [v1.92.0] - 2022-12-14

- #589 - @wez470 - load-balancers: Minor doc fixup
- #585 - @StephenVarela - Add firewall support for load balancers
- #587 - @StephenVarela - Support new http alerts for load balancers
- #586 - @andrewsomething - godo.go: Sort service lists.
- #583 - @ddebarros - Adds support for functions trigger API

## [v1.91.1] - 2022-11-23

- #582 - @StephenVarela - Load Balancers: Support new endpoints for http alerts

## [v1.90.0] - 2022-11-16

- #571 - @kraai - Add WaitForAvailable
- #579 - @bentranter - Deprecate old pointer helpers, use generic one
- #580 - @StephenVarela - LBAAS Fixup default http idle timeout behaviour
- #578 - @StephenVarela - LBAAS-2430 Add support for HTTP idle timeout seconds
- #577 - @ddebarros - Functions api support

## [v1.89.0] - 2022-11-02

- #575 - @ghostlandr - apps: add option to get projects data from Apps List endpoint

## [v1.88.0] - 2022-10-31

- #573 - @kamaln7 - apps: add ListBuildpacks, UpgradeBuildpack
- #572 - @ghostlandr - Apps: add project id as a parameter to CreateApp and to the App struct
- #570 - @kraai - Fix copy-and-paste error in comment
- #568 - @StephenVarela - LBAAS-2321 Add project_id to load balancers structs

## [v1.87.0] - 2022-10-12

- #564 - @DWizGuy58 - Add public monitoring alert policies for dbaas
- #565 - @dylanrhysscott - CON-5657 (Re-)expose public HA enablement flags in godo
- #563 - @andrewsomething - Add option to configure a rate.Limiter for the client.

## [v1.86.0] - 2022-09-23

- #561 - @jonfriesen - apps: add docr image deploy on push

## [v1.85.0] - 2022-09-21

- #560 - @andrewsomething - Bump golang.org/x/net (fixes: #557).
- #559 - @kamaln7 - apps: update component spec interfaces
- #555 - @kamaln7 - apps: add accessor methods and spec helpers
- #556 - @kamaln7 - update CI for go 1.18 & 1.19

## [v1.84.1] - 2022-09-16

- #554 - @andrewsomething - reserved IPs: project_id should have omitempty in create req.

## [v1.84.0] - 2022-09-16

- #552 - @andrewsomething - reserved IPs: Expose project_id and locked attributes.
- #549 - @rpmoore - adding the replica id to the database replica model

## [v1.83.0] - 2022-08-10

- #546 - @DWizGuy58 - Add support for database options

## [v1.82.0] - 2022-08-04

- #544 - @andrewsomething - apps: Add URN() method.
- #542 - @andrewsomething - databases: Support advanced config endpoints.
- #543 - @nicktate - Ntate/detection models
- #541 - @andrewsomething - droplets: Support listing Droplets filtered by name.
- #540 - @bentranter - Update links to API documentation

## [v1.81.0] - 2022-06-15

- #532 - @senorprogrammer - Add support for Reserved IP addresses
- #538 - @bentranter - util: update droplet create example
- #537 - @rpmoore - Adding project_id to databases
- #536 - @andrewsomething - account: Now may include info on current team.
- #535 - @ElanHasson - APPS-5636 Update App Platform for functions and Starter Tier App Proposals.

## [v1.80.0] - 2022-05-23

- #533 - @ElanHasson - APPS-5636 - App Platform updates

## [v1.79.0] - 2022-04-29

- #530 - @anitgandhi - monitoring: alerts for Load Balancers TLS conns/s utilization
- #529 - @ChiefMateStarbuck - Test against Go 1.18
- #528 - @senorprogrammer - Remove DisablePublicNetworking option from the Create path
- #527 - @senorprogrammer - Remove the WithFloatingIPAddress create option

## [v1.78.0] - 2022-03-31

- #522 - @jcodybaker - app platform: add support for features field

## [v1.77.0] - 2022-03-16

- #518 - @rcj4747 - apps: Update apps protos

## [v1.76.0] - 2022-03-09

- #516 - @CollinShoop - Add registry region support

## [v1.75.0] - 2022-01-27

- #508 - @ElanHasson - Synchronize public protos and add multiple specs

## [v1.74.0] - 2022-01-20

- #506 - @ZachEddy - Add new component type to apps-related structs

## [v1.73.0] - 2021-12-03

- #501 - @CollinShoop - Add support for Registry ListManifests and ListRepositoriesV2

## [v1.72.0] - 2021-11-29

- #500 - @ElanHasson - APPS-4420: Add PreservePathPrefix to AppRouteSpec

## [v1.71.0] - 2021-11-09

- #498 - @bojand - apps: update spec to include log destinations

## [v1.70.0] - 2021-11-01

- #491 - @andrewsomething - Add support for retrieving Droplet monitoring metrics.
- #494 - @alexandear - Refactor tests: replace t.Errorf with assert/require
- #495 - @alexandear - Fix typos and grammar issues in comments
- #492 - @andrewsomething - Update golang.org/x/net
- #486 - @abeltay - Fix typo on "DigitalOcean"

## [v1.69.1] - 2021-10-06

- #484 - @sunny-b - k8s/godo: remove ha field from update request

## [v1.69.0] - 2021-10-04

- #482 - @dikshant - godo/load-balancers: add DisableLetsEncryptDNSRecords field for LBaaS

## [v1.68.0] - 2021-09-29

- #480 - @sunny-b - kubernetes: add support for HA control plane

## [v1.67.0] - 2021-09-22

- #478 - @sunny-b - kubernetes: add supported_features field to the kubernetes/options response
- #477 - @wez470 - Add size unit to LB API.

## [v1.66.0] - 2021-09-21

- #473 - @andrewsomething - Add Go 1.17.x to test matrix and drop unsupported versions.
- #472 - @bsnyder788 - insights: add private (in/out)bound and public inbound bandwidth aler…
- #470 - @gottwald - domains: remove invalid json struct tag option

## [v1.65.0] - 2021-08-05

- #468 - @notxarb - New alerts feature for App Platform
- #467 - @andrewsomething - docs: Update links to API documentation.
- #466 - @andrewsomething - Mark Response.Monitor as deprecated.

## [v1.64.2] - 2021-07-23

- #464 - @bsnyder788 - insights: update HTTP method for alert policy update

## [v1.64.1] - 2021-07-19

- #462 - @bsnyder788 - insights: fix alert policy update endpoint

## [v1.64.0] - 2021-07-19

- #460 - @bsnyder788 - insights: add CRUD APIs for alert policies

## [v1.63.0] - 2021-07-06

- #458 - @ZachEddy - apps: Add tail_lines query parameter to GetLogs function

## [v1.62.0] - 2021-06-07

- #454 - @house-lee - add with_droplet_agent option to create requests

## [v1.61.0] - 2021-05-12

- #452 - @caiofilipini - Add support for DOKS clusters as peers in Firewall rules
- #448 - @andrewsomething - flip: Set omitempty for Region in FloatingIPCreateRequest.
- #451 - @andrewsomething - CheckResponse: Add RequestID from header to ErrorResponse when missing from body.
- #450 - @nanzhong - dbaas: handle ca certificates as base64 encoded
- #449 - @nanzhong - dbaas: add support for getting cluster CA
- #446 - @kamaln7 - app spec: update cors policy

## [v1.60.0] - 2021-04-04

- #443 - @andrewsomething - apps: Support pagination.
- #442 - @andrewsomething - dbaas: Support restoring from a backup.
- #441 - @andrewsomething - k8s: Add URN method to KubernetesCluster.

## [v1.59.0] - 2021-03-29

- #439 - @andrewsomething - vpcs: Support listing members of a VPC.
- #438 - @andrewsomething - Add Go 1.16.x to the testing matrix.

## [v1.58.0] - 2021-02-17

- #436 - @MorrisLaw - kubernetes: add name field to associated resources
- #434 - @andrewsomething - sizes: Add description field.
- #433 - @andrewsomething - Deprecate Name field in godo.DropletCreateVolume

## [v1.57.0] - 2021-01-15

- #429 - @varshavaradarajan - kubernetes: support optional cascading deletes for clusters
- #430 - @jonfriesen - apps: updates apps.gen.go for gitlab addition
- #431 - @nicktate - apps: update proto to support dockerhub registry type

## [v1.56.0] - 2021-01-08

- #422 - @kamaln7 - apps: add ProposeApp method

## [v1.55.0] - 2021-01-07

- #425 - @adamwg - registry: Support the storage usage indicator
- #423 - @ChiefMateStarbuck - Updated README example
- #421 - @andrewsomething - Add some basic input cleaning to NewFromToken
- #420 - @bentranter - Don't set "Content-Type" header on GET requests

## [v1.54.0] - 2020-11-24

- #417 - @waynr - registry: add support for garbage collection types

## [v1.53.0] - 2020-11-20

- #414 - @varshavaradarajan - kubernetes: add clusterlint support
- #413 - @andrewsomething - images: Support updating distribution and description.

## [v1.52.0] - 2020-11-05

- #411 - @nicktate - apps: add unspecified type to image source registry types
- #409 - @andrewsomething - registry: Add support for updating a subscription.
- #408 - @nicktate - apps: update spec to include image source
- #407 - @kamaln7 - apps: add the option to force build a new deployment

## [v1.51.0] - 2020-11-02

- #405 - @adamwg - registry: Support subscription options
- #398 - @reeseconor - Add support for caching dependencies between GitHub Action runs
- #404 - @andrewsomething - CONTRIBUTING.md: Suggest using github-changelog-generator.

## [v1.50.0] - 2020-10-26

- #400 - @waynr - registry: add garbage collection support
- #402 - @snormore - apps: add catchall_document static site spec field and failed-deploy job type
- #401 - @andrewlouis93 - VPC: adds option to set a VPC as the regional default

## [v1.49.0] - 2020-10-21

- #383 - @kamaln7 - apps: add ListRegions, Get/ListTiers, Get/ListInstanceSizes
- #390 - @snormore - apps: add service spec internal_ports

## [v1.48.0] - 2020-10-16

- #388 - @varshavaradarajan - kubernetes - change docr integration api routes
- #386 - @snormore - apps: pull in recent updates to jobs and domains

## [v1.47.0] - 2020-10-14

- #384 kubernetes - add registry related doks apis - @varshavaradarajan
- #385 Fixed some typo in apps.gen.go and databases.go file - @devil-cyber
- #382 Add GetKubeConfigWithExpiry (#334) - @ivanlemeshev
- #381 Fix golint issues #377 - @sidsbrmnn
- #380 refactor: Cyclomatic complexity issue - @DonRenando
- #379 Run gofmt to fix some issues in codebase - @mycodeself

## [v1.46.0] - 2020-10-05

- #373 load balancers: add LB size field, currently in closed beta - @anitgandhi

## [v1.45.0] - 2020-09-25

**Note**: This release contains breaking changes to App Platform features currently in closed beta.

- #369 update apps types to latest - @kamaln7
- #368 Kubernetes: add taints field to node pool create and update requests - @timoreimann
- #367 update apps types, address marshaling bug - @kamaln7

## [v1.44.0] - 2020-09-08

- #364 apps: support aggregate deployment logs - @kamaln7

## [v1.43.0] - 2020-09-08

- #362 update apps types - @kamaln7

## [v1.42.1] - 2020-08-06

- #360 domains: Allow for SRV records with port 0. - @andrewsomething

## [v1.42.0] - 2020-07-22

- #357 invoices: add category to InvoiceItem - @rbutler
- #358 apps: add support for following logs - @nanzhong

## [v1.41.0] - 2020-07-17

- #355 kubernetes: Add support for surge upgrades - @varshavaradarajan

## [v1.40.0] - 2020-07-16

- #347 Make Rate limits thread safe - @roidelapluie
- #353 Reuse TCP connection - @itsksaurabh

## [v1.39.0] - 2020-07-14

- #345, #346 Add app platform support [beta] - @nanzhong

## [v1.38.0] - 2020-06-18

- #341 Install 1-click applications on a Kubernetes cluster - @keladhruv
- #340 Add RecordsByType, RecordsByName and RecordsByTypeAndName to the DomainsService - @viola

## [v1.37.0] - 2020-06-01

- #336 registry: URL encode repository names when building URLs. @adamwg
- #335 Add 1-click service and request. @scottcrawford03

## [v1.36.0] - 2020-05-12

- #331 Expose expiry_seconds for Registry.DockerCredentials. @andrewsomething

## [v1.35.1] - 2020-04-21

- #328 Update vulnerable x/crypto dependency - @bentranter

## [v1.35.0] - 2020-04-20

- #326 Add TagCount field to registry/Repository - @nicktate
- #325 Add DOCR EA routes - @nicktate
- #324 Upgrade godo to Go 1.14 - @bentranter

## [v1.34.0] - 2020-03-30

- #320 Add VPC v3 attributes - @viola

## [v1.33.1] - 2020-03-23

- #318 upgrade github.com/stretchr/objx past 0.1.1 - @hilary

## [v1.33.0] - 2020-03-20

- #310 Add BillingHistory service and List endpoint - @rbutler
- #316 load balancers: add new enable_backend_keepalive field - @anitgandhi

## [v1.32.0] - 2020-03-04

- #311 Add reset database user auth method - @zbarahal-do

## [v1.31.0] - 2020-02-28

- #305 invoices: GetPDF and GetCSV methods - @rbutler
- #304 Add NewFromToken convenience method to init client - @bentranter
- #301 invoices: Get, Summary, and List methods - @rbutler
- #299 Fix param expiry_seconds for kubernetes.GetCredentials request - @velp

## [v1.30.0] - 2020-02-03

- #295 registry: support the created_at field - @adamwg
- #293 doks: node pool labels - @snormore

## [v1.29.0] - 2019-12-13

- #288 Add Balance Get method - @rbutler
- #286,#289 Deserialize meta field - @timoreimann

## [v1.28.0] - 2019-12-04

- #282 Add valid Redis eviction policy constants - @bentranter
- #281 Remove databases info from top-level godoc string - @bentranter
- #280 Fix VolumeSnapshotResourceType value volumesnapshot -> volume_snapshot - @aqche

## [v1.27.0] - 2019-11-18

- #278 add mysql user auth settings for database users - @gregmankes

## [v1.26.0] - 2019-11-13

- #272 dbaas: get and set mysql sql mode - @mikejholly

## [v1.25.0] - 2019-11-13

- #275 registry/docker-credentials: add support for the read/write parameter - @kamaln7
- #273 implement the registry/docker-credentials endpoint - @kamaln7
- #271 Add registry resource - @snormore

## [v1.24.1] - 2019-11-04

- #264 Update isLast to check p.Next - @aqche

## [v1.24.0] - 2019-10-30

- #267 Return []DatabaseFirewallRule in addition to raw response. - @andrewsomething

## [v1.23.1] - 2019-10-30

- #265 add support for getting/setting firewall rules - @gregmankes
- #262 remove ResolveReference call - @mdanzinger
- #261 Update CONTRIBUTING.md - @mdanzinger

## [v1.22.0] - 2019-09-24

- #259 Add Kubernetes GetCredentials method - @snormore

## [v1.21.1] - 2019-09-19

- #257 Upgrade to Go 1.13 - @bentranter

## [v1.21.0] - 2019-09-16

- #255 Add DropletID to Kubernetes Node instance - @snormore
- #254 Add tags to Database, DatabaseReplica - @Zyqsempai

## [v1.20.0] - 2019-09-06

- #252 Add Kubernetes autoscale config fields - @snormore
- #251 Support unset fields on Kubernetes cluster and node pool updates - @snormore
- #250 Add Kubernetes GetUser method - @snormore

## [v1.19.0] - 2019-07-19

- #244 dbaas: add private-network-uuid field to create request

## [v1.18.0] - 2019-07-17

- #241 Databases: support for custom VPC UUID on migrate @mikejholly
- #240 Add the ability to get URN for a Database @stack72
- #236 Fix omitempty typos in JSON struct tags @amccarthy1

## [v1.17.0] - 2019-06-21

- #238 Add support for Redis eviction policy in Databases @mikejholly

## [v1.16.0] - 2019-06-04

- #233 Add Kubernetes DeleteNode method, deprecate RecycleNodePoolNodes @bouk

## [v1.15.0] - 2019-05-13

- #231 Add private connection fields to Databases - @mikejholly
- #223 Introduce Go modules - @andreiavrammsd

## [v1.14.0] - 2019-05-13

- #229 Add support for upgrading Kubernetes clusters - @adamwg

## [v1.13.0] - 2019-04-19

- #213 Add tagging support for volume snapshots - @jcodybaker

## [v1.12.0] - 2019-04-18

- #224 Add maintenance window support for Kubernetes- @fatih

## [v1.11.1] - 2019-04-04

- #222 Fix Create Database Pools json fields - @sunny-b

## [v1.11.0] - 2019-04-03

- #220 roll out vpc functionality - @jheimann

## [v1.10.1] - 2019-03-27

- #219 Fix Database Pools json field - @sunny-b

## [v1.10.0] - 2019-03-20

- #215 Add support for Databases - @mikejholly

## [v1.9.0] - 2019-03-18

- #214 add support for enable_proxy_protocol. - @mregmi

## [v1.8.0] - 2019-03-13

- #210 Expose tags on storage volume create/list/get. - @jcodybaker

## [v1.7.5] - 2019-03-04

- #207 Add support for custom subdomains for Spaces CDN [beta] - @xornivore

## [v1.7.4] - 2019-02-08

- #202 Allow tagging volumes - @mchitten

## [v1.7.3] - 2018-12-18

- #196 Expose tag support for creating Load Balancers.

## [v1.7.2] - 2018-12-04

- #192 Exposes more options for Kubernetes clusters.

## [v1.7.1] - 2018-11-27

- #190 Expose constants for the state of Kubernetes clusters.

## [v1.7.0] - 2018-11-13

- #188 Kubernetes support [beta] - @aybabtme

## [v1.6.0] - 2018-10-16

- #185 Projects support [beta] - @mchitten

## [v1.5.0] - 2018-10-01

- #181 Adding tagging images support - @hugocorbucci

## [v1.4.2] - 2018-08-30

- #178 Allowing creating domain records with weight of 0 - @TFaga
- #177 Adding `VolumeLimit` to account - @lxfontes

## [v1.4.1] - 2018-08-23

- #176 Fix cdn flush cache API endpoint - @sunny-b

## [v1.4.0] - 2018-08-22

- #175 Add support for Spaces CDN - @sunny-b

## [v1.3.0] - 2018-05-24

- #170 Add support for volume formatting - @adamwg

## [v1.2.0] - 2018-05-08

- #166 Remove support for Go 1.6 - @iheanyi
- #165 Add support for Let's Encrypt Certificates - @viola

## [v1.1.3] - 2018-03-07

- #156 Handle non-json errors from the API - @aknuds1
- #158 Update droplet example to use latest instance type - @dan-v

## [v1.1.2] - 2018-03-06

- #157 storage: list volumes should handle only name or only region params - @andrewsykim
- #154 docs: replace first example with fully-runnable example - @xmudrii
- #152 Handle flags & tag properties of domain record - @jaymecd

## [v1.1.1] - 2017-09-29

- #151 Following user agent field recommendations - @joonas
- #148 AsRequest method to create load balancers requests - @lukegb

## [v1.1.0] - 2017-06-06

### Added

- #145 Add FirewallsService for managing Firewalls with the DigitalOcean API. - @viola
- #139 Add TTL field to the Domains. - @xmudrii

### Fixed

- #143 Fix oauth2.NoContext depreciation. - @jbowens
- #141 Fix DropletActions on tagged resources. - @xmudrii

## [v1.0.0] - 2017-03-10

### Added

- #130 Add Convert to ImageActionsService. - @xmudrii
- #126 Add CertificatesService for managing certificates with the DigitalOcean API. - @viola
- #125 Add LoadBalancersService for managing load balancers with the DigitalOcean API. - @viola
- #122 Add GetVolumeByName to StorageService. - @protochron
- #113 Add context.Context to all calls. - @aybabtme
