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

15 minutes

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

You need two Vault Enterprise clusters to perform the steps in this guide.


## Steps

Assuming that you have a Vault cluster in EU and wish to span across United
States by setting up a secondary cluster. Enable the performance replication
from the primary cluster in EU to the secondary cluster in US to provide
localized access.


Leverage the mount filter feature to whitelist or blacklist secrets engines from
being replicated across the regions.

1. Segment GDPR and non-GDPR secret engines
2. Configure the mount filters


~> Ensure GDPR data is segmented by secret mount and blacklist the movement of
those secret mounts to non-GDPR territories.


### <a name="step1"></a>Step 1: Segment GDPR and non-GDPR secret engine

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


![GDPR KV](/assets/images/vault-mount-filter-4.png)


## Next steps

Read [Production Hardening](/guides/operations/production.html) to learn best
practices for a production hardening deployment of Vault.
