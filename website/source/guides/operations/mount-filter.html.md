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

Mount filters are new way of controlling which secrets are moved across clusters
and physical regions as a result of replication. With mount filters, users can
select which mounts will be replicated as part of a performance replication
relationship.

By default, all non-local secret engines and associated data are replicated as
part of replication. The mount filter feature allows users to whitelist and/or
blacklist which secret engines are replicated, thereby allowing users to further
control the movement of secrets across their infrastructure.

![Performance Replication](/assets/images/vault-mount-filter-2.png)


## Reference Materials

- Preparing for GDPR Compliance with HashiCorp Vault [webinar](https://www.hashicorp.com/resources/preparing-for-gdpr-compliance-with-hashicorp-vault)
- Preparing for GDPR Compliance with HashiCorp Vault [blog post](https://www.hashicorp.com/blog/preparing-for-gdpr-compliance-with-hashicorp-vault)  
- [Create Mounts Filter (API)](/api/system/replication-performance.html#create-mounts-filter)

## Estimated Time to Complete

10 minutes

## Challenge

[General Data Protection Regulation (GDPR)](https://www.eugdpr.org/) is designed
to strengthen data protections and privacy for all individuals within the
European Union.  It requires that personally identifiable data not be physically
transferred to locations outside the European Union unless the region or country
has an equal rigor of data protection regulation as the EU.

Failure to abide by GDPR will result in fines as high as 20 million EUR or 4% of
the global annual revenue (whichever greater).


## Solution

Vault's mount filter feature allows users to specify which secret engines are
replicated as a part of the [performance
replication](/guides/operations/reference-architecture.html#vault-replication).

Leverage this feature to abide by data movements and sovereignty regulations
while ensuring performance access across geographically distributed regions.

Watch the recording of [***Preparing for GDPR Compliance with HashiCorp
Vault***](https://www.hashicorp.com/resources/preparing-for-gdpr-compliance-with-hashicorp-vault)
webinar.

[![YouTube](/assets/images/vault-mount-filter.png)](https://youtu.be/hmf6sN4W8pE)

## Prerequisites

This intermediate Vault operations guide assumes that you have some working
knowledge of Vault.

You need two Vault Enterprise clusters: one representing the EU cluster, and
another representing the US cluster.


## Steps

The scenario here is that you have a Vault cluster in EU and wish to span across
the United States by setting up a secondary cluster and enable the performance
replication. However, some data must remain in EU and should not be replicated
to the US cluster.

![Guide Scenario](/assets/images/vault-mount-filter-0.png)

Leverage the mount filter feature to whitelist or blacklist secrets engines from
being replicated across the regions.

1. [Segment GDPR and non-GDPR secret engines](#step1)
1. [Enable performance replication with mount filter](#step2)
1. [Verify the mount filter](#step3)
1. [Enable a local secret engine](#step4)

~> Ensure GDPR data is segmented by secret mount and blacklist the movement of
those secret mounts to non-GDPR territories.


### <a name="step1"></a>Step 1: Segment GDPR and non-GDPR secret engines

Enable key/value secret engines:

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

1. Create a mount filter to blacklist `EU_GDPR_data`.

    ```plaintext
    $ vault write sys/replication/performance/primary/mount-filter/secondary  \
           mode="blacklist" paths="EU_GDPR_data/"
    ```

1. Enable performance replication on the **secondary** cluster.

    ```plaintext
    $ vault write sys/replication/performance/secondary/enable token="..."
    ```
    Where the token is the token obtained from the primary cluster.



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
           https://eu-west-1.compute.com:8200/v1/sys/replication/performance/primary/secondary-token
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

1. Create a mount filter to blacklist `EU_GDPR_data`.

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

    Where the token in `payload.json` is the token obtained from the primary
    cluster.

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



### <a name="step3"></a>Step 3: Verify the mount filter

Once the replication completes, verify that you can see the `US_NON_GDPR_data`
but not `EU_GDPR_data`.

#### CLI command

```plaintext
$ vault read US_NON_GDPR_data/US
Key                 Value
---                 -----
refresh_interval    768h
secret              bar

$ vault read EU_GDPR_data/UK
No value found at EU_GDPR_data/UK
```

#### API call using cURL

```plaintext
$ curl --header "X-Vault-Token: ..." \
       https://eu-west-1.compute.com:8200/v1/sys/mounts | jq
{
   "US_NON_GDPR_data/": {
     "accessor": "kv_a5f54d9a",
     "config": {
       "default_lease_ttl": 0,
       "force_no_cache": false,
       "max_lease_ttl": 0,
       "plugin_name": ""
     },
     "description": "",
     "local": false,
     "options": {
       "version": "2"
     },
     "seal_wrap": false,
     "type": "kv"
   },
   ...
}
```

#### Web UI

![Secrets](/assets/images/vault-mount-filter-11.png)


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


## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn best
practices for a production hardening deployment of Vault.
