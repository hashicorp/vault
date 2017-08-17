---
layout: "guides"
page_title: "Disaster Recovery and Failure Scenarios"
sidebar_current: "disaster-recovery-and-failure-scenarios"
description: |-
  Learn how to create fault-tolerant infrastructures that can respond to catastrophic failure using Vault Replication.
---

# Replication Setup &amp; Guidance

If you're unfamiliar with Vault Replication concepts, please first look at the
[general information page](/docs/vault-enterprise/replication/index.html). More
details can be found in the
[replication internals](/docs/internals/replication.html) document.

Vault replication also includes a complete API. For more information, please see
the [Vault Replication API documentation](/api/system/replication.html)


## Reference Architecture

For the purposes of discussing the following failure scenarios, we will use the sample
architecture below of a 4 cluster Vault Enterprise environment spread across the continental
United States.


[![Vault Replication Reference Architecture](/assets/images/vault-rekey-vs-rotate.svg)](/assets/images/dr.png)

Cluster A serves as the primary cluster for this relationship. It is replicating data to Cluster B,
a performance secondary which is providing scalable performance to applications and users located 
geographically closer on the East Coast. Cluster C, which is located in Oregon, is also replicating
data from Cluster A and serves as a disaster recovery secondary. 

Cluster B is also further protected by serving as a DR primary to Cluster D, which is serving as a DR
secondary of Cluster B and is located regionally closer to its primary.

## Failure Scenarios and Responses

| Failure Scenario                                        | Reading Secrets                                                                                                                                                                                                | Writing Secrets                                                                                                                                                                | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
|---------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Scenario 1: Temporary network interruption between A->B | For applications reading from B: uninterrupted as this cluster serves as a replica of A.   Applications attempting to Read from A should be redirected to B. This can be done using a load balancer or Consul. | B will schedule writes for A, which will attempt to be confirmed once connectivity is restored.                                                                                | In the event of a longer term network outage, the Vault Enterprise user may want to consider promoting D as the new primary so that the data written on B in the meantime will be persistent and protected longer term.  B would then need to be demoted and serve as a performance secondary to D.  Once operations resume with A, this process would need to be repeated with A serving first as a DR secondary of D to first sync state, then conduct the promotion and demotion of B and D accordingly to reorganize the cluster to its original environment. |
| Scenario 2: Network interruption between A->C           | No impact: C is not read/writeable as of Vault 0.8.                                                                                                                                                            | No impact: C is not read/writeable as of Vault 0.8.                                                                                                                            | Once connectivity between A->C is restored, C will update with any changes made to A automatically.                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| Scenario 3: Catastrophic failure of Cluster A           | Secrets may only be read off of B until the new primary is promoted                                                                                                                                            | Secrets may be committed to Cluster B, but these will not be reflected outside of B and D unless connectivity is restored with A or C is immediately promoted and synced to B. | In this situation, Cluster C should immediately be promoted to becoming the new primary.  To do this, promote C as the new primary and establish B as a new performance secondary to C such that C->B. D will then reflect the changes in B automatically as part of its previous B->D relationship.   Once A is restored, you can reorganize the environment to its original state as per the last steps of Scenario 1.                                                                                                                                          |
| Scenario 4: Catastrophic failure of Cluster B           | Secrets may be read off of cluster A.                                                                                                                                                                          | Secrets may be written to cluster A.                                                                                                                                           | In this situation, cluster D should be promoted to allow for a regional performance secondary of A.  To do this, establish A->D such that D is serving as a performance secondary to A. When B resumes operation, you can resume the original architecture of your environment by demoting D, terminating the replication relationship of A->D, and returning D to its role as a DR secondary of B.                                                                                                                                                               |