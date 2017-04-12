---
layout: "docs"
page_title: "DynamoDB - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-dynamodb"
description: |-
  The DynamoDB storage backend is used to persist Vault's data in DynamoDB
  table.
---

# DynamoDB Storage Backend

The DynamoDB storage backend is used to persist Vault's data in
[DynamoDB][dynamodb] table.

- **High Availability** – the DynamoDB storage backend supports high
  availability. Because DynamoDB uses the time on the Vault node to implement
  the session lifetimes on its locks, significant clock skew across Vault nodes
  could cause contention issues on the lock.

- **Community Supported** – the DynamoDB storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "dynamodb" {
  ha_enabled = "true"
  region     = "us-west-2"
  table      = "vault-data"
}
```

For more information about the read/write capacity of DynamoDB tables, please
see the [official AWS DynamoDB documentation][dynamodb-rw-capacity].

## `dynamodb` Parameters

- `endpoint` `(string: "")` – Specifies an alternative, AWS compatible, DynamoDB
  endpoint. This can also be provided via the environment variable
  `AWS_DYNAMODB_ENDPOINT`.

- `ha_enabled` `(bool: false)` – Specifies whether this backend should be used
  to run Vault in high availability mode. This can also be provided via the
  environment variable `DYNAMODB_HA_ENABLED`.

- `max_parallel` `(string: "128")` – Specifies the maximum number of concurrent
  requests.

- `region` `(string "us-east-1")` – Specifies the AWS region. This can also be
  provided via the environment variable `AWS_DEFAULT_REGION`.

- `read_capacity` `(int: 5)` – Specifies the maximum number of reads consumed
  per second on the table. This can also be provided via the environment
  variable `AWS_DYNAMODB_READ_CAPACITY`.

- `table` `(string: "vault-dynamodb-backend")` – Specifies the name of the
  DynamoDB table in which to store Vault data. If the specified table does not
  yet exist, it will be created during initialization. This can also be
  provided via the environment variable `AWS_DYNAMODB_TABLE`.

- `write_capacity` `(int: 5)` – Specifies the maximum number of writes performed
  per second on the table. This can also be provided via the environment
  variable `AWS_DYNAMODB_WRITE_CAPACITY`.

The following settings are used for authenticating to AWS. If you are
running your Vault server on an EC2 instance, you can also make use of the EC2
instance profile service to provide the credentials Vault will use to make
DynamoDB API calls. Leaving the `access_key` and `secret_key` fields empty will
cause Vault to attempt to retrieve credentials from the AWS metadata service.

- `access_key` `(string: <required>)` – Specifies the AWS access key. This can
  also be provided via the environment variable `AWS_ACCESS_KEY_ID`.

- `secret_key` `(string: <required>)` – Specifies the AWS secret key. This can
  also be provided via the environment variable `AWS_SECRET_ACCESS_KEY`.

- `session_token` `(string: "")` – Specifies the AWS session token. This can
  also be provided via the environment variable `AWS_SESSION_TOKEN`.

This backend also supports the following high availability parameters. These are
discussed in more detail in the [HA concepts page](/docs/concepts/ha.html).

- `cluster_addr` `(string: "")` – Specifies the address to advertise to other
  Vault servers in the cluster for request forwarding. This can also be provided
  via the environment variable `VAULT_CLUSTER_ADDR`. This is a full URL, like
  `redirect_addr`, but Vault will ignore the scheme (all cluster members always
  use TLS with a private key/certificate).

- `disable_clustering` `(bool: false)` – Specifies whether clustering features
  such as request forwarding are enabled. Setting this to true on one Vault node
  will disable these features _only when that node is the active node_.

- `redirect_addr` `(string: <required>)` – Specifies the address (full URL) to
  advertise to other Vault servers in the cluster for client redirection. This
  can also be provided via the environment variable `VAULT_REDIRECT_ADDR`.

## `dynamodb` Examples

### Custom Table and Read-Write Capacity

This example shows using a custom table name and read/write capacity.

```hcl
storage "dynamodb" {
  table = "my-vault-data"

  read_capacity  = 10
  write_capacity = 15
}
```

### Enabling High Availability

This example show enabling high availability for the DynamoDB storage backend.

```hcl
storage "dynamodb" {
  ha_enabled    = "true"
  redirect_addr = "vault-leader.my-company.internal"
}
```

[dynamodb]: https://aws.amazon.com/dynamodb/
[dynamodb-rw-capacity]: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html#ProvisionedThroughput
