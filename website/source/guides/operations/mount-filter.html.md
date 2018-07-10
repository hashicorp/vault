---
layout: "guides"
page_title: "Vault Mount Filter - Guides"
sidebar_current: "guides-operations-mount-filter"
description: |-
  This guide demonstrates how to selectively filter out secret mounts for
  Performance Replication.
---

# Vault Mount Filter

~> **Enterprise Only:** Mount filter feature is a part of _Vault Enterprise Premium_.

Mount filters are a new way of controlling which secrets are moved across
clusters and physical regions as a result of replication. With mount filters,
users can select which secret engines will be replicated as part of a
performance replication relationship.

By default, all non-local secret engines and associated data are replicated as
part of replication. The mount filter feature allows users to whitelist or
blacklist which secret engines are replicated, thereby allowing users to further
control the movement of secrets across their infrastructure.

![Performance Replication](/assets/images/vault-mount-filter-2.png)


## Reference Materials

- Preparing for GDPR Compliance with HashiCorp Vault [webinar](https://www.hashicorp.com/resources/preparing-for-gdpr-compliance-with-hashicorp-vault)
- Preparing for GDPR Compliance with HashiCorp Vault [blog post](https://www.hashicorp.com/blog/preparing-for-gdpr-compliance-with-hashicorp-vault)  
- [Create Mounts Filter (API)](/api/system/replication-performance.html#create-mounts-filter)
- [Performance Replication and Disaster Recovery (DR) Replication](/docs/enterprise/replication/index.html#performance-replication-and-disaster-recovery-dr-replication)

## Estimated Time to Complete

10 minutes

## Challenge

[**General Data Protection Regulation (GDPR)**](https://www.eugdpr.org/) is designed
to strengthen data protection and privacy for all individuals within the
European Union.  It requires that personally identifiable data not be physically
transferred to locations outside the European Union unless the region or country
has an equal rigor of data protection regulation as the EU.

Failure to abide by GDPR will result in fines as high as 20 million EUR or 4% of
the global annual revenue (whichever greater).


## Solution

Leverage Vault's **mount filter** feature to abide by data movements and
sovereignty regulations while ensuring performance access across geographically
distributed regions.

The [***Preparing for GDPR Compliance with HashiCorp
Vault***](https://www.hashicorp.com/resources/preparing-for-gdpr-compliance-with-hashicorp-vault)
webinar discusses the GDPR compliance further in details.

[![YouTube](/assets/images/vault-mount-filter.png)](https://youtu.be/hmf6sN4W8pE)

## Prerequisites

This intermediate Vault operations guide assumes that you have some working
knowledge of Vault.

You need two Vault Enterprise clusters: one representing the EU cluster, and
another representing the US cluster both backed by Consul for storage.


## Steps

**Scenario:**  You have a Vault cluster in EU and wish to span across the United
States by setting up a secondary cluster and enable the performance
replication. However, some data must remain in EU and should ***not*** be
replicated to the US cluster.

![Guide Scenario](/assets/images/vault-mount-filter-0.png)

Leverage the mount filter feature to blacklist the secrets, that are subject to
GDPR, from being replicated across the regions.

1. [Segment GDPR and non-GDPR secret engines](#step1)
1. [Enable performance replication with mount filter](#step2)
1. [Verify the replication mount filter](#step3)
1. [Enable a local secret engine](#step4)

~> **NOTE:** Ensure that GDPR data is segmented by secret mount and blacklist
the movement of those secret mounts to non-GDPR territories.


### <a name="step1"></a>Step 1: Segment GDPR and non-GDPR secret engines

In the EU cluster (primary cluster), enable key/value secret engines:

- At **`EU_GDPR_data`** for GDPR data
- At **`US_NON_GDPR_data`** for non-GDPR data localized for US

#### CLI command

```shell
# For GDPR data
$ vault secrets enable -path=EU_GDPR_data kv-v2

# For non-GDPR data accessible in US
$ vault secrets enable -path=US_NON_GDPR_data kv-v2
```

#### API call using cURL

```shell
# For GDPR data
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"type":"kv-v2"}' \
       https://eu-west-1.compute.com:8200/v1/sys/mounts/EU_GDPR_data

# For non-GDPR data accessible in US
$ curl --header "X-Vault-Token: ..." \
      --request POST \
      --data '{"type":"kv-v2"}' \
      https://eu-west-1.compute.com:8200/v1/sys/mounts/US_NON_GDPR_data
```

#### Web UI

Open a web browser and launch the Vault UI (e.g.
https://eu-west-1.compute.com:8200/ui) and then login.

Select **Enable new engine** and enter corresponding parameter values:

![GDPR KV](/assets/images/vault-mount-filter-3.png)

![Non-GDPR KV](/assets/images/vault-mount-filter-4.png)


Click **Enable Engine** to complete.


### <a name="step2"></a>Step 2: Enable performance replication with mount filter

#### CLI command

1. Enable performance replication on the **primary** cluster.

    ```plaintext
    $ vault write -f sys/replication/performance/primary/enable
    WARNING! The following warnings were returned from Vault:

    * This cluster is being enabled as a primary for replication. Vault will be
    unavailable for a brief period and will resume service shortly.
    ```

1. Generate a secondary token.

    ```plaintext
    $ vault write sys/replication/performance/primary/secondary-token id="secondary"
    Key                              Value
    ---                              -----
    wrapping_token:                  eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyIjoiaHR0cDovLzE3Mi4zMS4yMC4xODA6ODIwMyIsImV4cCI6MTUyNTg0ODAxMywiaWF0IjoxNTI1ODQ2MjEzLCJqdGkiOiJlNTFiMjUxZi01ZTg2LTg4OWEtNGZmMy03NTQzMjRkNTdlMGQiLCJ0eXBlIjoid3JhcHBpbmcifQ.MIGGAkE2dDj3nmaoLHg7oldQ1iZPD0U8doyj3x3mQUVfTl8W99QYG8GM6VGVzhRPGvKctGriuo2oXN_8euWQb01M1y6n7gJBSu-qdXw-v2RieOyopAHls1bWhw4sO9Nlds8IDFA15vqkLXnq2g4_5lvlhxpP7B8dEOHvWXkHG4kJ_mKvrgR0dU0
    wrapping_accessor:               6ded4fb0-5e8c-2a37-1b3e-823673220348
    wrapping_token_ttl:              30m
    wrapping_token_creation_time:    2018-05-09 06:10:13.437421436 +0000 UTC
    wrapping_token_creation_path:    sys/replication/performance/primary/secondary-token
    ```

1. Create a **mount filter** to blacklist `EU_GDPR_data`.

    ```plaintext
    $ vault write sys/replication/performance/primary/mount-filter/secondary  \
           mode="blacklist" paths="EU_GDPR_data/"
    ```

1. Enable performance replication on the **secondary** cluster.

    ```plaintext
    $ vault write sys/replication/performance/secondary/enable token="..."
    ```
    Where the `token` is the `wrapping_token` obtained from the primary cluster.

    !> **NOTE:** This will immediately clear all data in the secondary cluster.

#### API call using cURL

1. Enable performance replication on the **primary** cluster.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"primary_cluster_addr":"https://eu-west-1.compute.com:8200"}' \
           https://eu-west-1.compute.com:8200/v1/sys/replication/performance/primary/enable
    ```

1. Generate a secondary token.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{ "id": "secondary"}' \
           https://eu-west-1.compute.com:8200/v1/sys/replication/performance/primary/secondary-token | jq
      {
       "request_id": "",
       "lease_id": "",
       "renewable": false,
       "lease_duration": 0,
       "data": null,
       "wrap_info": {
         "token": "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyIjoiaHR0cDovLzEyNy4wLjAuMTo4MjAzIiwiZXhwIjoxNTI1ODI5Njc2LCJpYXQiOjE1MjU4Mjc4NzYsImp0aSI6IjAwNmVkMDdjLWQ0MzYtZWViYy01OWYwLTdiMTU0ZGFmMDNiMCIsInR5cGUiOiJ3cmFwcGluZyJ9.MIGHAkF6saWWL-oRQMJIoUnaUOHNkcoHZCBwQs6mSMjPBopMi8DkGCJGBrh4jgV2mSzwFY1r5Ne7O66HmuMpm40MsYqjAQJCANSco_Sx5q6FmQSfoY-HtsVO1_YKWF4O6B7gYCvPKYkODMIwe5orCSgmIDyXHZt-REPm0sfdk4ZNyRCIRK5hDWyQ",
         "accessor": "6ea2a4e2-2926-120f-f288-c2348c78fb3e",
         "ttl": 1800,
         "creation_time": "2018-05-09T01:04:36.514715311Z",
         "creation_path": "sys/replication/performance/primary/secondary-token"
       },
       "warnings": null,
       "auth": null
     }
    ```

1. Create a **mount filter** to blacklist `EU_GDPR_data`.

    ```plaintext
    $ tee payload.json <<EOF
    {
      "mode": "blacklist",
      "paths": [ "EU_GDPR_data/" ]
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload.json \
           https://eu-west-1.compute.com:8200/v1/sys/replication/performance/primary/mount-filter/secondary
    ```

1. Enable performance replication on the **secondary** cluster.

    ```plaintext
    $ tee payload.json <<EOF
    {
      "token": "..."
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload.json \
           https://us-central.compute.com:8201/v1/sys/replication/performance/secondary/enable
    ```

    Where the `token` in `payload.json` is the token obtained from the primary
    cluster.

    !> **NOTE:** This will immediately clear all data in the secondary cluster.


#### Web UI

1. Select **Replication** and check the **Performance** radio button.
  ![Performance Replication - primary](/assets/images/vault-mount-filter-5.png)

1. Click **Enable replication**.

1. Select the **Secondaries** tab, and then click **Add**.
  ![Performance Replication - primary](/assets/images/vault-mount-filter-6.png)

1. Populate the **Secondary ID** field, and then select **Configure performance
mount filtering** to set your mount filter options.  You can filter by
whitelisting or blacklisting. For this example, select **Blacklist**.

1. Check **EU_GDPR_data** to prevent it from being replicated to the secondary
cluster.
  ![Performance Replication - primary](/assets/images/vault-mount-filter-7.png)

1. Click **Generate token**.
  ![Performance Replication - primary](/assets/images/vault-mount-filter-8.png)

1. Click **Copy** to copy the token.

1. Now, launch the Vault UI for the secondary cluster (e.g. https://us-central.compute.com:8201/ui), and then click **Replication**.

1. Check the **Performance** radio button, and then select **secondary** under the **Cluster mode**. Paste the token you copied from the primary.
  ![Performance Replication - secondary](/assets/images/vault-mount-filter-9.png)

1. Click **Enable replication**.

<br>

~> **NOTE:** At this point, the secondary cluster must be unsealed using the
**primary cluster's unseal key**. If the secondary is in an HA cluster, ensure
that each standby is sealed and unsealed with the primary’s unseal keys. The
secondary cluster mirrors the configuration of its primary cluster's backends
such as auth methods, secret engines, audit devices, etc. It uses the primary as
the _source of truth_ and passes token requests to the primary.


Restart the secondary vault server (e.g. `https://us-central.compute.com:8201`)
and unseal it with the primary cluster's unseal key.

```plaintext
$ vault operator unseal
Unseal Key (will be hidden): <primary_cluster_unseal_key>
```

The initial root token on the secondary no longer works. Use the auth methods
configured on the primary cluster to log into the secondary.

**Example:**

Enable and configure the userpass auth method on the **primary** cluster and
create a new username and password.

```shell
# Enable the userpass auth method on the primary
$ vault auth enable userpass

# Create a user with admin policy
$ vault write auth/userpass/users/james password="passw0rd" policy="admin"
```

-> Alternatively, you can [generate a new root token](/guides/operations/generate-root.html)
using the primary cluster's unseal key. However, it is recommended that root
tokens are only used for just enough initial setup or in emergencies.


Log into the **secondary** cluster using the enabled auth method.

```plaintext
$ vault login -method=userpass username=james password="passw0rd"
```


### <a name="step3"></a>Step 3: Verify the replication mount filter

Once the replication completes, verify that the secrets stored in the
`EU_GDPR_data` never get replicated to the US cluster.

#### CLI command

On the **EU** cluster, write some secrets:

```shell
# Write some secret at EU_GDPR_data/secret
$ vault kv put EU_GDPR_data/secret pswd="password"
Key              Value
---              -----
created_time     2018-05-10T18:00:38.912587665Z
deletion_time    n/a
destroyed        false
version          1

# Write some secret at US_NON_GDPR_data/secret
$ vault kv put US_NON_GDPR_data/secret apikey="my-api-key"
Key              Value
---              -----
created_time     2018-05-10T18:04:37.554665851Z
deletion_time    n/a
destroyed        false
version          1
```

From the **US** cluster, read the secrets:

```shell
# Read the secrets at EU_GDPR_data/secret
$ vault kv get EU_GDPR_data/secret
No value found at EU_GDPR_data/secret

# Read the secrets at US_NON_GDPR_data/secret
$ vault kv get US_NON_GDPR_data/secret
====== Metadata ======
Key              Value
---              -----
created_time     2018-05-10T18:09:07.717250408Z
deletion_time    n/a
destroyed        false
version          1

===== Data =====
Key       Value
---       -----
apikey    my-api-key
```


#### API call using cURL

On the **EU** cluster, write some secret:

```shell
# Create the request payload
$ tee payload.json <<EOF
{
  "data": {
    "pswd": "password"
  }
}
EOF

# Write some secret at EU_GDPR_data/secret
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       https://eu-west-1.compute.com:8200/v1/EU_GDPR_data/data/secret

# Create the request payload
$ tee payload-us.json <<EOF
{
 "data": {
   "apikey": "my-api-key"
 }
}
EOF

# Write some secret at US_NON_GDPR_data/secret
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload-us.json \
       https://eu-west-1.compute.com:8200/v1/US_NON_GDPR_data/data/secret
```

From the **US** cluster, read the secrets:

```shell
# Read the secrets at EU_GDPR_data/secret
$ curl --header "X-Vault-Token: ..." \
       https://us-central.compute.com:8201/v1/EU_GDPR_data/data/secret | jq
{
  "errors": []
}

# Read the secrets at US_NON_GDPR_data/secret
$ curl --header "X-Vault-Token: ..." \
       https://us-central.compute.com:8201/v1/US_NON_GDPR_data/data/secret | jq
{
  "request_id": "f5eedb5b-9406-c519-f5a5-070336c10205",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": [
    "Invalid path for a versioned K/V secrets engine. See the API docs for the appropriate API endpoints to use. If using the Vault CLI, use 'vault kv get' for this operation."
  ],
  "auth": null
}
```

#### Web UI

On the **EU** cluster, select **EU_GDPR_data** > **Create secret**:

![Secrets](/assets/images/vault-mount-filter-12.png)

Enter the values and click **Save**.  Repeat the step to write some secrets at
the **US_NON_GDPR_data** path as well.


On the **US** cluster, select **US_NON_GDPR_data**. You should be able to see
the `apikey` under `US_NON_GDPR_data/secret`.

![Secrets](/assets/images/vault-mount-filter-13.png)

The **EU_GDPR_data** data is not replicated, so you won't be able to see the
secrets.


### <a name="step4"></a>Step 4: Enable a local secret engine

When replication is enabled, you can mark the secrets engine local only.  Local
secret engines are not replicated or removed by replication.

Login to the **secondary** cluster and enable key/value secret engine at
`US_ONLY_data` to store secrets only valid for the US region.

#### CLI command

Pass the `-local` flag:

```plaintext
$ vault secrets enable -local -path=US_ONLY_data kv-v2
```

#### API call using cURL

Pass the `local` parameter in the API request:

```plaintext
$ tee payload.json <<EOF
{
  "type": "kv-v2",
  "local": true
}
EOF

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload.json \
       https://us-central.compute.com:8201/v1/sys/mounts/US_ONLY_data
```


#### Web UI

Be sure to select the check box for **Local** to keep it mounted locally within
the cluster.

![Local Secret](/assets/images/vault-mount-filter-10.png)

<br>

-> **NOTE:** `US_ONLY_data` only exists locally in the secondary cluster that
you won't be able to see it from the primary cluster.




## Next steps

Read [Vault Deployment Reference
Architecture](/guides/operations/reference-architecture.html) to learn more
about the recommended deployment practices.
