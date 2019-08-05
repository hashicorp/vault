---
layout: "docs"
page_title: "Kubernetes"
sidebar_current: "docs-platform-k8s-index"
sidebar_title: "Kubernetes"
description: |-
  This section documents the official integration between Vault and Kubernetes.
---

# Kubernetes

Vault can be deployed into Kubernetes using the official HashiCorp Vault Helm chart.  
The helm chart allows users to deploy Vault in various configurations:

* Dev mode: a single in-memory Vault server for testing Vault
* Standalone mode (default): a single Vault server persisting to a volume using the file storage backend
* HA mode: a cluster of Vault servers that use an HA storage backend such as Consul (default)

## Use Cases

**Running a Vault Service:** The Vault server cluster can run directly on Kubernetes.  
This can be used by applications running within Kubernetes as well as external to 
Kubernetes, as long as they can communicate to the server via the network.

**Accessing and Storing Secrets:** Applications using the Vault service running in 
Kubernetes can access and store secrets from Vault using a number of different 
[secret engines](https://www.vaultproject.io/docs/secrets/index.html) and 
[authentication methods](https://www.vaultproject.io/docs/auth/index.html).

**Running a Highly Available Vault Service:**  By using pod affinities, highly available 
backend storage (such as Consul) and [auto-unseal](https://www.vaultproject.io/docs/concepts/seal.html#auto-unseal), 
Vault can become a highly available service in Kubernetes.

**Encryption as a Service:** Applications using the Vault service running in Kubernetes 
can leverage the [Transit secret engine](https://www.vaultproject.io/docs/secrets/transit/index.html) 
as "encryption as a service".  This allows applications to offload encryption needs 
to Vault before storage data at rest.

**Audit Logs for Vault:** Operators can choose to attach a persistent volume 
to the Vault cluster which can be used to [store audit logs](https://www.vaultproject.io/docs/audit/index.html).

**And more!** Vault can run directly on Kubernetes, so in addition to the
native integrations provided by Vault itself, any other tool built for
Kubernetes can choose to leverage Vault.
