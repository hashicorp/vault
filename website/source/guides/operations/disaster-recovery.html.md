---
layout: "guides"
page_title: "Vault Disaster Recovery Replication Setup - Guides"
sidebar_title: "Disaster Recovery Setup"
sidebar_current: "guides-operations-dr"
description: |-
  This guide demonstrates step-by-step instruction of setting up a disaster
  recovery (DR) replications.
---

# Vault Disaster Recovery Replication

~> **Enterprise Only:** Disaster Recovery Replication is a part of _Vault Enterprise Pro_.

It is inevitable for organizations to have a disaster recovery (DR) strategy to
protect their Vault deployment against catastrophic failure of an entire
cluster. Vault Enterprise supports multi-datacenter deployment where you can
replicate data across datacenters for performance as well as disaster recovery.

A cluster is the basic unit of Vault Enterprise replication which follows the
leader-follower model. A leader cluster is referred to as the **primary**
cluster and is considered the _system of record_. Data is streamed from the
primary cluster to all **secondary** (follower) clusters.

![Replication Pattern](/img/vault-ref-arch-8.png)

~> **Important:** In DR replication, secondary clusters ***do not forward***
service read or write requests until they are promoted and become a new primary
- they essentially act as a warm standby cluster.


The [Mount Filter](/guides/operations/mount-filter.html) guide provides step-by-step
instructions on setting up performance replication.  This guide focuses on DR
replication setup.


## Reference Materials

- [Performance Replication and Disaster Recovery (DR) Replication](/docs/enterprise/replication/index.html#performance-replication-and-disaster-recovery-dr-replication)
- [DR Replication API](/api/system/replication-dr.html)
- [Replication Setup & Guidance](/guides/operations/replication.html)
- [Vault HA guide](/guides/operations/vault-ha-consul.html)



## Estimated Time to Complete

10 minutes


## Prerequisites

This intermediate Vault operations guide assumes that you have some working
knowledge of Vault.

You need two Vault Enterprise clusters: one behaves as the **primary cluster**,
and another becomes the **secondary**.

![DR Prerequisites](/img/vault-dr-0.png)

## Steps

This guide walk through the following operations:

1. [Enable DR Primary Replication](#step1)
1. [Enable DR Secondary Replication](#step2)
1. [Promote DR Secondary to Primary](#step3)
1. [Demote DR Primary to Secondary](#step4)
1. [Disable DR Primary](#step5)


### <a name="step1"></a>Step 1: Enable DR Primary Replication

#### CLI command

1. Enable DR replication on the **primary** cluster.

    ```plaintext
    $ vault write -f sys/replication/dr/primary/enable
    WARNING! The following warnings were returned from Vault:

    * This cluster is being enabled as a primary for replication. Vault will be
    unavailable for a brief period and will resume service shortly.
    ```

1. Generate a secondary token.

    ```plaintext
    $ vault write sys/replication/dr/primary/secondary-token id="secondary"
    ```

    The output should look similar to:

    ```plaintext
    Key                              Value
    ---                              -----
    wrapping_token:                  eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyIjoiaHR0cDovLzEzLjU3LjIwLjQxOjgyMDAiLCJleHAiOjE1MjkzMzkzMzEsImlhdCI6MTUyOTMzNzUzMSwianRpIjoiZDZmMmMzZTItMTZjNS1mNTU0LWYxMzAtNzMzZDE0OWNiNTIzIiwidHlwZSI6IndyYXBwaW5nIn0.MIGIAkIArsC3s1x7GYnEbaYwAbYUj-Wgp4B3Q3kVXL0BbaKvsECySV4Pwtm--i24OSQfI9zAlsG8ZypOWJdngRa59wlhWdQCQgG22-I-aNWPehjsqmwwEADU-u37LUrR6O0MsUCqtfWYwIM9o7PFP1wMZ4JwDGftQXUH6hIrkXZDxnnGsSCJ1Vl75w
    wrapping_accessor:               bab0ea36-23f6-d21d-4ca6-a9c3673766a3
    wrapping_token_ttl:              30m
    wrapping_token_creation_time:    2018-06-18 15:58:51.645117216 +0000 UTC
    wrapping_token_creation_path:    sys/replication/dr/primary/secondary-token
    ```

    -> Copy the generated **`wrapping_token`** which you will need to enable the DR
    secondary cluster.

#### API call using cURL

1. Enable DR replication on the **primary** cluster by invoking **`/sys/replication/dr/primary/enable`** endpoint.

    **Example:**

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{}' \
           https://cluster-A.example.com:8200/v1/sys/replication/dr/primary/enable
     {
       "request_id": "ef38af20-9c1f-138a-2d03-bbb6410fb0fc",
       "lease_id": "",
       "renewable": false,
       "lease_duration": 0,
       "data": null,
       "wrap_info": null,
       "warnings": [
         "This cluster is being enabled as a primary for replication. Vault will be
         unavailable for a brief period and will resume service shortly."
       ],
       "auth": null
     }
    ```

1. Generate a secondary token by invoking **`/sys/replication/dr/primary/secondary-token`** endpoint.

    **Example:**

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{ "id": "secondary"}' \
           https://cluster-A.example.com:8200/v1/sys/replication/dr/primary/secondary-token | jq
     {
       "request_id": "",
       "lease_id": "",
       "renewable": false,
       "lease_duration": 0,
       "data": null,
       "wrap_info": {
         "token": "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyIjoiaHR0cDovLzEzLjU3LjIwLjQxOjgyMDAiLCJleHAiOjE1MjkzNDQzMjcsImlhdCI6MTUyOTM0MjUyNywianRpIjoiYmRiZTJiNzEtODgwMS05YjZjLTNjMTQtMzVkNDI3NDQ3MjEzIiwidHlwZSI6IndyYXBwaW5nIn0.MIGIAkIBmESVVq_83l9hixTN7Ot0v5XQMsQfi1zV9APooZWkLvbS2olBWSQnskykQQH6GskMOi-ypOlAabqxWmfoCLA8-TICQgHRdkbJGgAQtWmjc8Z-ZEgymMv8YZq6qQxbUtPXloyM-cf_1Y1qmdGDYWtjPqoF5m1Bt_WkAJl9MguVb04QMWSotw",
         "accessor": "7e56e9da-178c-119d-1d01-807a203fa0b3",
         "ttl": 1800,
         "creation_time": "2018-06-18T17:22:07.129747708Z",
         "creation_path": "sys/replication/dr/primary/secondary-token"
       },
       "warnings": null,
       "auth": null
     }
    ```

    -> Copy the generated **`token`** which you will need to enable the DR
    secondary cluster.

#### Web UI

Open a web browser and launch the Vault UI (e.g.
https://cluster-A.example.com:8200/ui) and then login.

1. Select **Replication** and check the **Disaster Recovery (DR)** radio button.
  ![DR Replication - primary](/img/vault-dr-1.png)

1. Click **Enable replication**.

1. Select the **Secondaries** tab, and then click **Add**.
  ![DR Replication - primary](/img/vault-dr-2.png)

1. Populate the **Secondary ID** field, and click **Generate token**.
  ![DR Replication - primary](/img/vault-dr-3.png)

1. Click **Copy** to copy the token which you will need to enable the DR secondary cluster.
  ![DR Replication - primary](/img/vault-dr-4.png)


<br>

### <a name="step2"></a>Step 2: Enable DR Secondary Replication

The following operations must be performed on the DR secondary cluster.

#### CLI command

1. Enable DR replication on the **secondary** cluster.

    ```plaintext
    $ vault write sys/replication/dr/secondary/enable token="..."
    ```
    Where the `token` is the `wrapping_token` obtained from the primary cluster.

    Expected output:

    ```plaintext
    WARNING! The following warnings were returned from Vault:

    * Vault has successfully found secondary information; it may take a while to
    perform setup tasks. Vault will be unavailable until these tasks and initial
    sync complete.
    ```

    !> **NOTE:** This will immediately clear all data in the secondary cluster.

#### API call using cURL

1. Enable DR replication on the **secondary** cluster.

    ```plaintext
    $ tee payload.json <<EOF
    {
      "token": "..."
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload.json \
           https://cluster-B.example.com:8200/v1/sys/replication/dr/secondary/enable | jq
    {
       "request_id": "7a9730c1-b6fc-6557-5c0a-081e1f89ed2d",
       "lease_id": "",
       "renewable": false,
       "lease_duration": 0,
       "data": null,
       "wrap_info": null,
       "warnings": [
         "Vault has successfully found secondary information; it may take a while
         to perform setup tasks. Vault will be unavailable until these tasks and
         initial sync complete."
       ],
       "auth": null
     }
    ```

    Where the `token` in `payload.json` is the token obtained from the primary
    cluster.

    !> **NOTE:** This will immediately clear all data in the secondary cluster.


#### Web UI

1. Now, launch the Vault UI for the **secondary** cluster (e.g. https://cluster-B.example.com:8200/ui) and click **Replication**.

1. Check the **Disaster Recovery (DR)** radio button and select **secondary** under the **Cluster mode**. Paste the token you copied from the primary in the **Secondary activation token** field.
  ![DR Replication - secondary](/img/vault-dr-5.png)

1. Click **Enable replication**.
  ![DR Replication - secondary](/img/vault-dr-5.2.png)

  !> **NOTE:** This will immediately clear all data in the secondary cluster.



<br>

### <a name="step3"></a>Step 3: Promote DR Secondary to Primary

This step walks you through the promotion of the secondary cluster to become the
new primary when a catastrophic failure causes the primary cluster to be
inoperable. Refer to the [_Important Note about Automated DR
Failover_](#important) section for more background information.

First, you must generate a **DR operation token** which you need to promote the
secondary cluster. The process, outlined below using API calls, is the similar to [_Generating a Root Token (via CLI)_](/guides/operations/generate-root.html).

#### From Terminal

1. Generate an one time password (OTP) to use:

    ```plaintext
    $ vault operator generate-root -dr-token -generate-otp
    HenFLWmt0AgrjWJp/RECzQ==
    ```

1. Start the DR operation token generation process by invoking **`/sys/replication/dr/secondary/generate-operation-token/attempt`** endpoint.

    **Example:**

    ```plaintext
    $ tee payload.json <<EOF
    {
      "otp": "HenFLWmt0AgrjWJp/RECzQ=="
    }
    EOF

    $ curl --request PUT \
         --data @payload.json \
         https://cluster-B.example.com:8200/v1/sys/replication/dr/secondary/generate-operation-token/attempt | jq
     {
       "nonce": "455bf989-6575-1262-c0d0-a94eaf60bdd0",
       "started": true,
       "progress": 0,
       "required": 3,
       "complete": false,
       "encoded_token": "",
       "encoded_root_token": "",
       "pgp_fingerprint": ""
     }
    ```

    -> Distribute the generated **`nonce`** to each unseal key holder.

1. In order to generate a DR operation token, a quorum of unseal keys must be
entered by each key holder via **`/sys/replication/dr/secondary/generate-operation-token/update`** endpoint.

    **Example:**

    ```plaintext
    $ tee payload_key1.json <<EOF
    {
      "key": "<primary_unseal_key_1>",
      "nonce": "455bf989-6575-1262-c0d0-a94eaf60bdd0"
    }
    EOF

    $ curl --request PUT \
           --data @payload_key1.json \
           https://cluster-B.example.com:8200/v1/sys/replication/dr/secondary/generate-operation-token/update | jq
     {
       "nonce": "455bf989-6575-1262-c0d0-a94eaf60bdd0",
       "started": true,
       "progress": 1,
       "required": 3,
       "complete": false,
       "encoded_token": "",
       "encoded_root_token": "",
       "pgp_fingerprint": ""
     }
    ```

    This operation must be executed by each unseal key holder. Once the quorum
    has been reached, the output contains the encoded DR operation token
    (`encoded_token`).

    **Example:**

    ```plaintext
    $ curl --request PUT \
         --data @payload_key3.json \
         https://cluster-B.example.com:8200/v1/sys/replication/dr/secondary/generate-operation-token/update | jq
    {
      "nonce": "455bf989-6575-1262-c0d0-a94eaf60bdd0",
      "started": true,
      "progress": 3,
      "required": 3,
      "complete": true,
      "encoded_token": "dKNQqNmh3JfJcSZdGlkttQ==",
      "encoded_root_token": "",
      "pgp_fingerprint": ""
    }
    ```

1. Decode the generated DR operation token (`encoded_token`).

    **Example:**

    ```plaintext
    $ vault operator generate-root -dr-token \
            -decode="dKNQqNmh3JfJcSZdGlkttQ==" \
            -otp="HenFLWmt0AgrjWJp/RECzQ=="

    23e02f22-2ae6-94cc-d93f-5ee295e03e9d
    ```

1. Finally, promote the DR secondary to become the primary by invoking the
**`sys/replication/dr/secondary/promote`** endpoint. The request payload must
contains the DR operation token.  

    **Example:**

    ```plaintext
    $ tee payload.json <<EOF
    {
	     "dr_operation_token": "23e02f22-2ae6-94cc-d93f-5ee295e03e9d"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
             --request POST \
             --data @payload.json \
             https://cluster-B.example.com:8200/v1/sys/replication/dr/secondary/promote | jq
    {
      "request_id": "3879546b-1dc7-8490-521b-80104ad761b5",
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": null,
      "wrap_info": null,
      "warnings": [
        "This cluster is being promoted to a replication primary. Vault will be unavailable
        for a brief period and will resume service shortly."
      ],
      "auth": null
    }
    ```


#### Web UI

1. Click on **Generate OTP** to generate an OTP.  Then click **Copy OTP**.
    ![DR Replication - secondary](/img/vault-dr-6.png)

1. Click **Generate Operation Token**.

1. A quorum of unseal keys must be entered to create a new operation token for
the DR secondary.

    ![DR Replication - secondary](/img/vault-dr-7.png)

    -> This operation must be performed by each unseal-key holder.


1. Once the quorum has been reached, the output displays the encoded DR operation token.  Click **Copy CLI command**.

    ![DR Replication - secondary](/img/vault-dr-8.png)

1. Execute the CLI command from a terminal to generate a DR operation token
using the OTP generated earlier. (Be sure to enter your OTP in the command.)

    **Example:**

    ```
    $ vault operator generate-root -dr-token \
            -otp="vZpZZf5UI1nvB3A5/7Xq9A==" \          
            -decode="cuplaFGYduDEY6ZVC5IfaA=="

    cf703c0d-afcc-55b9-2b64-d66cf427f59c
    ```

1. Now, click **Promote** tab, and then enter the generated DR operation token.

    ![DR Replication - secondary](/img/vault-dr-9-1.png)

1. Click **Promote cluster**.

    When you prompted, "_Are you sure you want to promote this cluster?_", click **Promote cluster** again to complete.

    ![DR Replication - secondary](/img/vault-dr-9.png)

<br>

> Once the secondary cluster was successfully promoted, you should be able to
log in using the original primary cluster's root token or via configured
authentication method.  If desired, generate a [new root
token](/guides/operations/generate-root.html).



### <a name="step4"></a>Step 4: Demote DR Primary to Secondary

If the _original_ DR primary cluster becomes operational again, you may want to
utilize the cluster by making it a DR secondary cluster. This step explains how
to demote the original DR primary cluster to a secondary.

~> Remember that there is only **one** primary cluster available to the clients
in DR replication.

#### CLI command

Execute the following command to demote the original DR primary cluster to a
secondary.

```plaintext
$ vault write -f sys/replication/dr/primary/demote

WARNING! The following warnings were returned from Vault:

  * This cluster is being demoted to a replication secondary. Vault will be
  unavailable for a brief period and will resume service shortly.
```

This secondary cluster will not attempt to connect to a primary (see the
update-primary call), but will maintain knowledge of its cluster ID and can be
reconnected to the same DR replication set without wiping local storage.


#### API call using cURL

Invoke the **`sys/replication/dr/secondary/enable`** endpoint to demote the
original DR primary cluster to a secondary.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       https://cluster-A.example.com:8200/v1/sys/replication/dr/primary/demote | jq
{
   "request_id": "8a40adac-6eb7-c798-48d0-f7cdd25fdd6f",
   "lease_id": "",
   "renewable": false,
   "lease_duration": 0,
   "data": null,
   "wrap_info": null,
   "warnings": [
     "This cluster is being demoted to a replication secondary. Vault will be unavailable for a brief period and will resume service shortly."
   ],
   "auth": null
}
```

This secondary cluster will not attempt to connect to a primary (see the
update-primary call), but will maintain knowledge of its cluster ID and can be
reconnected to the same DR replication set without wiping local storage.

#### Web UI

Select **Replication** and click **Demote cluster**.

![DR Replication - demotion](/img/vault-dr-10.png)

When you prompted, "_Are you sure you want to demote this cluster?_", click
**Demote cluster** again to complete.

![DR Replication - demotion](/img/vault-dr-12.png)


### <a name="step5"></a>Step 5: Disable DR Primary

Once the DR secondary cluster was promoted to be the **new primary**, you may
want to disable the DR replication on the _original_ primary when it becomes
operational again.

~> Remember that there is only **one** primary cluster available to the clients
in DR replication.


#### CLI command

Execute the following command to disable DR replication.

```plaintext
$ vault write -f sys/replication/dr/primary/disable

WARNING! The following warnings were returned from Vault:

  * This cluster is having replication disabled. Vault will be unavailable for
    a brief period and will resume service shortly.
```

Any secondaries will no longer be able to connect.


#### API call using cURL

Invoke the **`sys/replication/dr/primary/disable`** endpoint to disable DR
replication.

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       https://cluster-A.example.com:8200/v1/sys/replication/dr/primary/disable | jq
{
   "request_id": "92a5f57a-2f7b-11be-b9dd-0f028396fba8",
   "lease_id": "",
   "renewable": false,
   "lease_duration": 0,
   "data": null,
   "wrap_info": null,
   "warnings": [
     "This cluster is having replication disabled. Vault will be unavailable for a brief period and will resume service shortly."
   ],
   "auth": null
}
```

Any secondaries will no longer be able to connect.


#### Web UI

Select **Replication** and click **Disable replication**.

![DR Replication - demotion](/img/vault-dr-11.png)

When you prompted, "_Are you sure you want to disable replication on this
cluster?_", click **Disable** again to complete.

![DR Replication - demotion](/img/vault-dr-13.png)

Any secondaries will no longer be able to connect.

!> **Caution:** Once this is done, re-enabling the DR replication as a primary
will change the cluster's ID.  Its connecting secondaries will require a wipe of
the underlying storage even if they have connected before. If re-enabling DR
replication as a secondary, its underlying storage will be wiped when connected
to a primary.


## <a name="important"></a>Important Note about Automated DR Failover

Vault does not support an automatic failover/promotion of a DR secondary
cluster, and this is a deliberate choice due to the difficulty in accurately
evaluating why a failover should or shouldn't happen. For example, imagine a
DR secondary loses its connection to the primary. Is it because the primary is
down, or is it because networking between the two has failed?

If the DR secondary promotes itself and clients start connecting to it, you now
have two active clusters whose data sets will immediately start diverging.
There's no way to understand simply from one perspective or the other which one
of them is right.

Vault's API supports programmatically performing various replication operations
which allows the customer to write their own logic about automating some of
these operations based on experience within their own environments. You can
review the available replication APIs at the following links:

- [Vault Replication API](/api/system/replication.html)
- [DR Replication API](/api/system/replication-dr.html)
- [Performance Replication API](/api/system/replication-performance.html)



## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn more
about the guidance on hardening the production deployments of Vault.
