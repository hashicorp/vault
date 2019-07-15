# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)

## 5.14.0 - 2019-07-09
### Added
- Support for network security groups in the Load Balancing service
- Support for network security groups in Core Services
- Support for network security groups on database systems in the Database service
- Support for creating autonomous transaction processing and autonomous data warehouse previews in the Database service
- Support for getting the load balancer attachments of instance pools in the Compute service
- Support for moving resources across compartments in the Resource Manager service
- Support for moving VCN resources across compartments in the Networking service

## 5.13.0 - 2019-07-02
### Added
- Support for moving images, instance configurations, and instance pools across compartments in Core Services
- Support for moving autoscaling configurations across compartments in the Compute Autoscaling service

### Fixed
- Fixed a bug where the Streaming service's endpoints in Tokyo, Seoul, and future regions were not reachable from the SDK

## 5.12.0 - 2019-06-25
### Added
- Support for moving senders across compartments in the Email service
- Support for moving NAT gateway resources across compartments in Core Services

## 5.11.0 - 2019-06-18
### Added
- Support for moving service gateway resources across compartments in Core Services
- Support for moving block storage resources across compartments in Core Services
- Support for key deletion in the Key Management service

## 5.10.0 - 2019-06-11
### Added
- Support for specifying custom boot volume sizes on instance configurations in the Compute Autoscaling service
- Support for 'Autonomous Transaction Processing - Dedicated' features, as well as maintenance run and backup operations on autonomous databases, autonomous container databases, and autonomous Exadata infrastructure in the Database service

## 5.9.0 - 2019-06-04
### Added
- Support for autoscaling autonomous databases and autonomous data warehouses in the Database service
- Support for specifying fault domains as part of instance configurations in the Compute Autoscaling service
- Support for deleting tag definitions and tag namespaces in the Identity service

### Fixed
- Support for regions in realms other than oraclecloud.com in the Load Balancing service

## 5.8.0 - 2019-05-28
### Added
- Support for the Work Requests service, and tracking of a number of Core Services operations through work requests
- Support for emulated volume attachments in Core Services
- Support for changing the compartment of resources in the File Storage service
- Support for tags in list operations in the File Storage service
- Support for returning UI password creation dates in the Identity service

## 5.7.0 - 2019-05-21
### Added
- Support for returning tags when listing instance configurations, instance pools, or autoscaling configurations in the Compute Autoscaling service
- Support for getting the namespace of another tenancy than the caller's tenancy in the Object Storage service
- Support for BGP dynamic routing and providing pre-shared secrets (PSKs) when establishing tunnels in the Networking service

## 5.6.0 - 2019-05-14
### Added
- Support for the Seoul (ICN) region
- Support for logging context fields on data-plane APIs of the Key Management Service
- Support for reverse pagination on list operations of the Email service
- Support for configuring backup retention windows on database backups in the Database service

## 5.5.0 - 2019-05-07
### Added
- Support for the Tokyo (NRT) region

- Support UploadManager for uploading large objects. Sample is available on [Github](https://github.com/oracle/oci-go-sdk/tree/master/example/example_objectstorage_test.go)

## 5.4.0 - 2019-04-16
### Added
- Support for tagging dynamic groups in the Identity service
- Support for updating network ACLs and license types for autonomous databases and autonomous data warehouses in the Database service
- Support for editing static routes and IPSec remote IDs in the Virtual Networking service

## 5.3.0 - 2019-04-09
### Added
- Support for etag and if-match headers (for optimistic concurrency control) in the Email service

## 5.2.0 - 2019-04-02
### Added
- Support for provider service key names on virtual circuits in the FastConnect service
- Support for customer reference names on cross connects and cross connect groups in the FastConnect service

## 5.1.0 - 2019-03-26
### Added
- Support for glob patterns and exclusions for object lifecycle management in the Object Storage service
- Documentation enhancements and corrections for traffic management in the DNS service

### Fixed
- The 'tag' info is always ignored in the returned string of Version() function [Github issue 157](https://github.com/oracle/oci-go-sdk/issues/157)

## 5.0.0 - 2019-03-19
### Added

- Support for specifying metadata on node pools in the Container Engine for Kubernetes service
- Support for provisioning a new autonomous database or autonomous data warehouse as a clone of another in the Database service
### Breaking changes
- The field``CreateAutonomousDatabaseDetails`` is no longer an anonymous field and the type changed from struct to interface in struct ``CreateAutonomousDatabaseRequest``. Here is sample code that shows how to update your code to incorporate this change. 

    - Before

    ```golang
    // create a CreateAutonomousDatabaseRequest
    // There were two ways to initialize the CreateAutonomousDatabaseRequest struct.
    // This breaking change only impact option #2
    request := database.CreateAutonomousDatabaseRequest{}

    // #1. Instantiate CreateAutonomousDatabaseDetails directly: no impact
    details := database.CreateAutonomousDatabaseDetails{}
    details.CompartmentId = common.String(getCompartmentID())
    // ... other properties

    // Set it to the request class
    request.CreateAutonomousDatabaseDetails = details

    // #2. Instantiate CreateAutnomousDatabaseDetails through  anonymous fields: will break
    request.CompartmentId = common.String(getCompartmentID())
    // ... other properties
    ```

    - After

    ```golang
    // #2 no longer supported. Create CreateAutonomousDatabaseDetails directly
    details := database.CreateAutonomousDatabaseDetails{}
    details.CompartmentId = common.String(getCompartmentID())
    // ... other properties

    // and set the details to CreateAutonomousDatabaseBase
    request := database.CreateAutonomousDatabaseRequest{}
    request.CreateAutonomousDatabaseDetails = details
    // ...
    ```


## 4.2.0 - 2019-03-12
### Added
- Support for the Budgets service
- Support for managing multifactor authentication in the Identity service
- Support for managing default tags in the Identity service
- Support for account recovery in the Identity service
- Support for authentication policies in the Identity service
- Support for specifying the workload type when creating autonomous databases in the Database service
- Support for I/O resource management for Exadata database systems in the Database service
- Support for customer-specified timezones on database systems in the Database service

## 4.1.0 - 2019-02-28
### Added
- Support for the Monitoring service
- Support for the Notification service
- Support for the Resource Manager service
- Support for the Compute Autoscaling service
- Support for changing the compartment of a tag namespace in the Identity service
- Support for specifying fault domains in the Database service
- Support for managing instance monitoring in the Compute service
- Support for attaching/detaching load balancers to instance pools in the Compute service

## 4.0.0 - 2019-02-21
### Added
- Support for government-realm regions
- Support for the Streaming service
- Support for tags in the Key Management service
- Support for regional subnets in the Virtual Networking service

### Fixed
- Removed unused Announcements service 'NotificationFollowupDetails' struct and 'GetFollowups' operation
- InstancePrincipals now invalidates a token shortly before its expiration time to avoid making  a service call with an expired token
- Requests with binary bodies that require its body to be included in the signature are now being signed correctly

## 3.7.0 - 2019-02-07
### Added
- Support for the Web Application Acceleration and Security (WAAS) service
- Support for the Health Checks service
- Support for connection strings on Database resources in the Database service
- Support for traffic management in the DNS service
- Support for tagging in the Email service
### Fixed
- Retry context in now cancelable during wait for new retry

## 3.6.0 - 2019-01-31
### Added
- Support for the Announcements service

## 3.5.0 - 2019-01-24
### Added

- Support for renaming databases during restore-from-backup operations in the Database service
- Built-in logging now supports log levels. More information about the changes can be found in the [go-docs page](https://godoc.org/github.com/oracle/oci-go-sdk#hdr-Logging_and_Debugging)
- Support for calling Oracle Cloud Infrastructure services in the ca-toronto-1 region

## 3.4.0 - 2019-01-10
### Added 
- Support for device attributes on volume attachments in the Compute service
- Support for custom header rulesets in the Load Balancing service


## 3.3.0 - 2018-12-13
### Added 
- Support for Data Guard for VM shapes in the Database service
- Support for sparse disk groups for Exadata shapes in the Database service
- Support for a new field, isLatestForMajorVersion, when listing DB versions in the Database service
- Support for in-transit encryption for paravirtualized boot volume and data volume attachments in the Block Storage service
- Support for tagging DNS Zones in the DNS service
- Support for resetting credentials for SCIM clients associated with an Identity provider and updating user capabilities in the Identity service

## 3.2.0 - 2018-11-29
### Added 
- Support for getting bucket statistics in the Object Storage service

### Fixed
- Block Storage service for copying volume backups across regions is now enabled 
- Object storage `PutObject` and `UploadPart` operations now do not override the client's signer

## 3.1.0 - 2018-11-15
### Added
- Support for VCN transit routing in the Networking service 

## 3.0.0 - 2018-11-01
### Added
- Support for modifying the route table, DHCP options and security lists associated with a subnet in the Networking service.
- Support for tagging of File Systems, Mount Targets and Snapshots in the File Storage service.
- Support for nested compartments in the Identity service

### Notes
- The version is bumped due to breaking changes in previous release.

## 2.7.0 - 2018-10-18
### Added
- Support for cost tracking tags in the Identity service
- Support for generating and downloading wallets in the Database service
- Support for creating a standalone backup from an on-premises database in the Database service
- Support for db version and additional connection strings in the Autonomous Transaction Processing and Autonomous Data Warehouse resources of the Database service
- Support for copying volume backups across regions in the Block Storage service
- Support for deleting compartments in the Identity service
- Support for reboot migration for virtual machines in the Compute service
- Support for Instance Pools and Instance Configurations in the Compute service

### Fixed
- The signing algorithm does not lower case the header fields [Github issue 132](https://github.com/oracle/oci-go-sdk/issues/132)
- Raw configuration provider does not check for empty strings [Github issue 134](https://github.com/oracle/oci-go-sdk/issues/134)

### Breaking change
- DbDataSizeInMBs field in Backup and BackupSummary struct was renamed to DatabaseSizeInGBs and type changed from *int to *float64 
    - Before
    ```golang
    // Size of the database in megabytes (MB) at the time the backup was taken.
    DbDataSizeInMBs *int `mandatory:"false" json:"dbDataSizeInMBs"`
    ```

    - After

    ```golang
    // The size of the database in gigabytes at the time the backup was taken.
    DatabaseSizeInGBs *float64 `mandatory:"false" json:"databaseSizeInGBs"`
    ```
- Data type for DatabaseEdition in Backup and BackupSummary struct was changed from *string to BackupDatabaseEditionEnum
    - Before

    ```golang
    // The Oracle Database edition of the DB system from which the database backup was taken.
    DatabaseEdition *string `mandatory:"false" json:"databaseEdition"`
    ```

    - After

    ```golang
     // The Oracle Database edition of the DB system from which the database backup was taken.
     DatabaseEdition BackupDatabaseEditionEnum `mandatory:"false" json:"databaseEdition,omitempty"`
    ```

## 2.6.0 - 2018-10-04
### Added
- Support for trusted partner images through application listings and subscriptions in the Compute service
- Support for object lifecycle policies in the Object Storage service
- Support for copying objects across regions in the Object Storage service
- Support for network address translation (NAT) gateways in the Networking service

## 2.5.0 - 2018-09-27
### Added
- Support for paravirtualized launch mode when importing images in the Compute service
- Support for Key Management service
- Support for encrypting the contents of an Object Storage bucket using a Key Management service key
- Support for specifying a Key Management service key when launching a compute instance in the Compute service
- Support for specifying a Key Management service key when backing up or restoring a block storage volume in the Block Volume service

## 2.4.0 - 2018-09-06
### Added
- Added support for updating metadata fields on an instance in the Compute service

## 2.3.0 - 2018-08-23
### Added
- Support for fault domain in the Identity Service
- Support for Autonomous Data Warehouse and Autonomous Transaction Processing in the Database service
- Support for resizing an offline volume in the Block Storage service
- Nil interface when polymorphic json response object is null

## 2.2.0 - 2018-08-09
### Added
- Support for fault domains in the Compute service
- A sample showing how to use Search service from the SDK is available on [Github](https://github.com/oracle/oci-go-sdk/tree/master/example/example_resourcesearch_test.go)

## 2.1.0 - 2018-07-26
### Added
- Support for the Search service
- Support for specifying a backup policy when creating a boot volume in the Block Storage service

### Fixed
- OCI error is missing opc-request-id value [Github Issue 120](https://github.com/oracle/oci-go-sdk/issues/120)
- Include raw http response when service error occurred

## 2.0.0 - 2018-07-12
### Added
- Support for tagging Load Balancers in the Load Balancing service
- Support for export options in the File Storage service
- Support for retrieving compartment name and user name as part of events in the Audit service

### Fixed
- CreateKubeconfig function should not close http reponse body [Github Issue 116](https://github.com/oracle/oci-go-sdk/issues/116)

### Breaking changes
- Datatype changed from *int to *int64 for several request/response structs. Here is sample code that shows how to update your code to incorporate this change. 

    - Before

    ```golang
    // Update the impacted properties from common.Int to common.Int64.
    // Here is the updates for CreateBootVolumeDetails
    details := core.CreateBootVolumeDetails{
        SizeInGBs: common.Int(10),
    }
    ```

    - After

    ```golang
    details := core.CreateBootVolumeDetails{
        SizeInGBs: common.Int64(10),
    }
    ```

- Impacted packages and structs
    - core
        - BootVolume.(SizeInGBs, SizeInMBs)
        - BootVolumeBackup.(SizeInGBs, UniqueSizeInGBs)
        - CreateBootVolumeDetails.SizeInGBs
        - CreateVolumeDetails.(SizeInGBs, SizeInMBs)
        - Image.SizeInMBs
        - InstanceSourceViaImageDetails.BootVolumeSizeInGBs
        - Volume.(SizeInGBs, SizeInMBs)
        - VolumeBackup.(SizeInGBs, SizeInMBs, UniqueSizeInGBs, UniqueSizeInMbs)
        - VolumeGroup.(SizeInMBs, SizeInGBs)
        - VolumeGroupBackup.(SizeInMBs, SizeInGBs, UniqueSizeInMbs, UniqueSizeInGbs)
    - dns
        - GetDomainRecordsRequest.Limit
        - GetRRSetRequest.Limit
        - GetZoneRecordsRequest.Limit
        - ListZonesRequest.Limit
        - Zone.Serial
        - ZoneSummary.Serial
    - filestorage
        - ExportSet.(MaxFsStatBytes, MaxFsStatFiles)
        - FileSystem.MeteredBytes
        - FileSystemSummary.MeteredBytes
        - UpdateExportSetDetails.(MaxFsStatBytes, MaxFsStatFiles)
    - identity
        - ApiKey.InactiveStatus
        - AuthToken.InactiveStatus
        - Compartment.InactiveStatus
        - CustomerSecretKey.InactiveStatus
        - CustomerSecretKeySummary.InactiveStatus
        - DynamicGroup.InactiveStatus
        - Group.InactiveStatus
        - IdentityProvider.InactiveStatus
        - IdpGroupMapping.InactiveStatus
        - Policy.InactiveStatus
        - Saml2IdentityProvider.InactiveStatus
        - SmtpCredential.InactiveStatus
        - SmtpCredentialSummary.InactiveStatus
        - SwiftPassword.InactiveStatus
        - UiPassword.InactiveStatus
        - User.InactiveStatus
        - UserGroupMembership.InactiveStatus
    - loadbalancer
        - ConnectionConfiguration.IdleTimeout
        - ListLoadBalancerHealthsRequest.Limit
        - ListLoadBalancersRequest.Limit
        - ListPoliciesRequest 
        - ListProtocolsRequest.Limit
        - ListShapesRequest.Limit
        - ListWorkRequestsRequest.Limit
    - objectstorage
        - GetObjectResponse.ContentLength
        - HeadObjectResponse.ContentLength
        - MultipartUploadPartSummary.Size
        - ObjectSummary.Size
        - PutObjectRequest.ContentLength
        - UploadPartRequest.ContentLength

## 1.8.0 - 2018-06-28
### Added
- Support for service gateway management in the Networking service
- Support for backup and clone of boot volumes in the Block Storage service

## 1.7.0 - 2018-06-14
### Added
- Support for the Container Engine service. A sample showing how to use this service from the SDK is available [Github](https://github.com/oracle/oci-go-sdk/tree/master/example/example_containerengine_test.go)

### Fixed
- Empty string was send to backend service for optional enum if it's not set

## 1.6.0 - 2018-05-31
### Added
- Support for the "soft shutdown" instance action in the Compute service
- Support for Auth Token management in the Identity service
- Support for backup or clone of multiple volumes at once using volume groups in the Block Storage service
- Support for launching a database system from a backup in the Database service

### Breaking changes
- ``LaunchDbSystemDetails`` is renamed to ``LaunchDbSystemBase`` and the type changed from struct to interface in ``LaunchDbSystemRequest``. Here is sample code that shows how to update your code to incorporate this change. 

    - Before

    ```golang
    // create a LaunchDbSystemRequest
    // There were two ways to initialize the LaunchDbSystemRequest struct.
    // This breaking change only impact option #2
    request := database.LaunchDbSystemRequest{}

    // #1. explicity create LaunchDbSystemDetails struct (No impact)
    details := database.LaunchDbSystemDetails{}
    details.AvailabilityDomain = common.String(validAD())
    details.CompartmentId = common.String(getCompartmentID())
    // ... other properties
    request.LaunchDbSystemDetails = details

    // #2. use anonymous fields (Will break)
    request.AvailabilityDomain = common.String(validAD())
    request.CompartmentId = common.String(getCompartmentID())
    // ...
    ```

    - After

    ```golang
    // create a LaunchDbSystemRequest
    request := database.LaunchDbSystemRequest{}
    details := database.LaunchDbSystemDetails{}
    details.AvailabilityDomain = common.String(validAD())
    details.CompartmentId = common.String(getCompartmentID())
    // ... other properties

    // set the details to LaunchDbSystemBase
    request.LaunchDbSystemBase = details
    // ...
    ```

## 1.5.0 - 2018-05-17
### Added
- ~~Support for backup or clone of multiple volumes at once using volume groups in the Block Storage service~~
- Support for the ability to optionally specify a compartment filter when listing exports in the File Storage service
- Support for tagging virtual cloud network resources in the Networking service
- Support for specifying the PARAVIRTUALIZED remote volume type when creating a virtual image or launching a new instance in the Compute service
- Support for tilde in private key path in configuration files

## 1.4.0 - 2018-05-03
### Added
- Support for ``event_name`` in Audit Service
- Support for multiple ``hostnames`` for loadbalancer listener in LoadBalance service
- Support for auto-generating opc-request-id for all operations
- Add opc-request-id property for all requests except for Object Storage which use opc-client-request-id

## 1.3.0 - 2018-04-19
### Added
- Support for retry on Oracle Cloud Infrastructure service APIs. Example can be found on [Github](https://github.com/oracle/oci-go-sdk/tree/master/example/example_retry_test.go)
- Support for tagging DbSystem and Database resources in the Database Service
- Support for filtering by DbSystemId in ListDbVersions operation in Database Service

### Fixed
- Fixed a request signing bug for PatchZoneRecords API
- Fixed a bug in DebugLn

## 1.2.0 - 2018-04-05
### Added
- Support for Email Delivery Service. Example can be found on [Github](https://github.com/oracle/oci-go-sdk/tree/master/example/example_email_test.go)
- Support for paravirtualized volume attachments in Core Services
- Support for remote VCN peering across regions
- Support for variable size boot volumes in Core Services
- Support for SMTP credentials in the Identity Service
- Support for tagging Bucket resources in the Object Storage Service

## 1.1.0 - 2018-03-27
### Added
- Support for DNS service
- Support for File Storage service
- Support for PathRouteSets and Listeners in Load Balancing service
- Support for Public IPs in Core Services
- Support for Dynamic Groups in Identity service
- Support for tagging in Core Services and Identity service. Example can be found on [Github](https://github.com/oracle/oci-go-sdk/tree/master/example/example_tagging_test.go)
- Fix ComposingConfigurationProvider to not accept a nil ConfigurationProvider
- Support for passphrase configuration to FileConfiguration provider

## 1.0.0 - 2018-02-28 Initial Release
### Added
- Support for Audit service
- Support for Core Services (Networking, Compute, Block Volume)
- Support for Database service
- Support for IAM service
- Support for Load Balancing service
- Support for Object Storage service
