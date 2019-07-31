---
layout: "docs"
page_title: "Helm - Kubernetes"
sidebar_current: "docs-platform-k8s-helm"
description: |-
  The Vault Helm chart is the recommended way to install and configure Vault on Kubernetes.
---

# Helm Chart

The [Vault Helm chart](https://github.com/hashicorp/vault-helm)
is the recommended way to install and configure Vault on Kubernetes.
In addition to running Vault itself, the Helm chart is the primary
method for installing and configuring Vault to integrate with other
services such as Consul for High Availability deployments.

This page assumes general knowledge of [Helm](https://helm.sh/) and
how to use it. Using Helm to install Vault will require that Helm is
properly installed and configured with your Kubernetes cluster.

-> **Important:** The Helm chart is new and
may still change significantly over time. Please always run Helm with
`--dry-run` before any install or upgrade to verify changes.

~> **Security Warning:** By default, the chart will install an insecure configuration
of Vault. This provides a less complicated out-of-box experience for new users,
but is not appropriate for a production setup. It is highly recommended to use
a [properly secured Kubernetes cluster](https://kubernetes.io/docs/tasks/administer-cluster/securing-a-cluster/). 
See the [architecture reference](/docs/platform/k8s/run.html#architecture) 
for a Vault Helm production deployment checklist.

## Using the Helm Chart

To use the Helm chart, you must download or clone the
[vault-helm GitHub repository](https://github.com/hashicorp/vault-helm)
and run Helm against the directory. We plan to transition to using a real
Helm repository soon. When running Helm, we highly recommend you always
checkout a specific tagged release of the chart to avoid any
instabilities from master.

Prior to this, you must have Helm installed and configured both in your
Kubernetes cluster and locally on your machine. The steps to do this are
out of the scope of this document, please read the
[Helm documentation](https://helm.sh/) for more information.

Example chart usage:

```sh
# Clone the chart repo
$ git clone https://github.com/hashicorp/vault-helm.git
$ cd vault-helm

# Checkout a tagged version
$ git checkout v0.1.0

# Run Helm
$ helm install --dry-run ./
```

## Configuration (Values)

The chart is highly customizable using
[Helm configuration values](https://docs.helm.sh/using_helm/#customizing-the-chart-before-installing).
Each value has a sane default tuned for an optimal getting started experience
with Vault. Before going into production, please review the parameters below
and consider if they're appropriate for your deployment.

* <a name="v-global" href="#v-global">`global`</a> - These global values affect multiple components of the chart.

  * <a name="v-global-enabled" href="#v-global-enabled">`enabled`</a> (`boolean: true`) - The master enabled/disabled configuration. If this is true, most components will be installed by default. If this is false, no components will be installed by default and manually opt-in is required, such as by setting <a href="#v-">`server.enabled`</a> to true.

  * <a name="v-global-image" href="#v-global-image">`image`</a> (`string: "vault:latest"`) - The name of the Docker image (including any tag) for the containers running Vault. **This should be pinned to a specific version when running in production.** Otherwise, other changes to the chart may inadvertently upgrade your Vault version.

* <a name="v-server" href="#v-server">`server`</a> - Values that configure running a Vault server within Kubernetes.

  * <a name="v-server-resources" href="#v-server-resources">`resources`</a> (`string: null`) - The resource requests and limits (CPU, memory, etc.) for each of the server. This should be a multi-line string mapping directly to a Kubernetes [ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#resourcerequirements-v1-core) object. If this isn't specified, then the pods won't request any specific amount of resources. **Setting this is highly recommended.**

    ```yaml
    # Resources are defined as a formatted multi-line string:
    resources: |
      requests:
        memory: "10Gi"
      limits:
        memory: "10Gi"
    ```

  * <a name="v-server-extraenvironmentvars" href="#v-server-extraenvironmentvars">`extraEnvironmentVars`</a> (`string: null`) - The extra environment variables to be applied to the Vault server.  This should be a multi-line key/value string.

    ```yaml
    # Extra Environment Variables are defined as key/value strings.
     extraEnvironmentVars:
       GOOGLE_REGION: global,
       GOOGLE_PROJECT: myproject,
       GOOGLE_CREDENTIALS: /vault/userconfig/myproject/myproject-creds.json
    ```

  * <a name="v-server-extravolumes" href="#v-server-extravolumes">`extraVolumes`</a> (`array: []`) - A list of extra volumes to mount to Vault servers. This is useful for bringing in extra data that can be referenced by other configurations at a well known path, such as TLS certificates. The value of this should be a list of objects. Each object supports the following keys:

      - <a name="v-server-extravolumes-type" href="#v-server-extravolumes-type">`type`</a> (`string: required`) -
      Type of the volume, must be one of "configMap" or "secret". Case sensitive.

      - <a name="v-server-extravolumes-name" href="#v-server-extravolumes-name">`name`</a> (`string: required`) -
      Name of the configMap or secret to be mounted. This also controls the path
      that it is mounted to. The volume will be mounted to `/vault/userconfig/<name>`.

      - <a name="v-server-extravolumes-load" href="#v-server-extravolumes-load">`load`</a> (`boolean: false`) -
      If true, then the agent will be configured to automatically load HCL/JSON
      configuration files from this volume with `-config-dir`. This defaults
      to false.

        ```yaml
        extraVolumes:
          -  type: "secret"
             name: "consul-certs"
             load: false
        ```

  * <a name="v-server-affinity" href="#v-server-affinity">`affinity`</a> (`string`) - This value defines the [affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity) for server pods. It defaults to allowing only a single pod on each node, which minimizes risk of the cluster becoming unusable if a node is lost. If you need to run more pods per node (for example, testing on Minikube), set this value to `null`.

    ```yaml
    # Recommended default server affinity:
    affinity: |
      podAntiAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
          matchLabels:
            app: {{ template "vault.name" . }}
            release: "{{ .Release.Name }}"
            component: server
          topologyKey: kubernetes.io/hostname
    ```

  * <a name="v-server-service" href="#v-server-service">`extraVolumes`</a> - This configures the `Service` resource create for the Vault server.

      - <a name="v-server-service-enabled" href="#v-server-service-enabled">`enabled`</a> (`boolean: true`) -
      Enables a service to allow other pods running in Kubernetes to communicate with the Vault server.

  * <a name="v-server-datastorage" href="#v-server-datastorage">`dataStorage`</a> - This configures the volume used for storing Vault data when not using external storage such as Consul.

      - <a name="v-server-datastorage-enabled" href="#v-server-datastorage-enabled">`enabled`</a> (`boolean: true`) -
      Enables a persistent volume to be created for storing Vault data when not using an external storage service.

      - <a name="v-server-datastorage-size" href="#v-server-datastorage-size">`size`</a> (`string: 10Gi`) -
      Size of the volume to be created for Vault's data storage when not using an external storage service.

      - <a name="v-server-datastorage-storageclass" href="#v-server-datastorage-storageclass">`storageClass`</a> (`string: null`) -
      Name of the storage class to use when creating the data storage volume.

      - <a name="v-server-datastorage-accessmode" href="#v-server-datastorage-accessmode">`accessMode`</a> (`string: ReadWriteOnce`) -
      Type of access mode of the storage device.  See https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes for more information.

  * <a name="v-server-auditstorage" href="#v-server-auditstorage">`auditStorage`</a> - This configures the volume used for storing Vault's audit logs.  See https://www.vaultproject.io/docs/audit/index.html.

      - <a name="v-server-auditstorage-enabled" href="#v-server-auditstorage-enabled">`enabled`</a> (`boolean: true`) -
      Enables a persistent volume to be created for storing Vault's audit logs.

      - <a name="v-server-auditstorage-size" href="#v-server-auditstorage-size">`size`</a> (`string: 10Gi`) -
      Size of the volume to be created for Vault's audit logs.

      - <a name="v-server-auditstorage-storageclass" href="#v-server-auditstorage-storageclass">`storageClass`</a> (`string: null`) -
      Name of the storage class to use when creating the audit storage volume.

      - <a name="v-server-auditstorage-accessmode" href="#v-server-auditstorage-accessmode">`accessMode`</a> (`string: ReadWriteOnce`) -
      Type of access mode of the storage device.

  * <a name="v-server-dev" href="#v-server-dev">`dev`</a> - This configures `dev` mode for the Vault server.

      - <a name="v-server-dev-enabled" href="#v-server-dev-enabled">`enabled`</a> (`boolean: false`) -
      Enables `dev` mode for the Vault server.  This mode is useful for experimenting with Vault without needing to unseal.

        ~> **Security Warning:** Never, ever, ever run a "dev" mode server in production. It is insecure and will lose data on every restart (since it stores data in-memory). It is only made for development or experimentation.

  * <a name="v-server-standalone" href="#v-server-standalone">`standalone`</a> - This configures `standalone` mode for the Vault server.

      - <a name="v-server-standalone-enabled" href="#v-server-standalone-enabled">`enabled`</a> (`boolean: true`) -
      Enables `standalone` mode for the Vault server.  This mode uses the `file` storage backend and requires a volume for persistence (`dataStorage`).

      - <a name="v-server-standalone-config" href="#v-server-standalone-config">`config`</a> (`string: "{}"`) -
      A raw string of extra HCL or JSON [configuration](https://www.vaultproject.io/docs/configuration/index.html) for Vault servers.
      This will be saved as-is into a ConfigMap that is read by the Vault servers.
      This can be used to add additional configuration that isn't directly exposed by the chart.

        ```yaml
        # ExtraConfig values are formatted as a multi-line string:
        config: |
          api_addr = "http://POD_IP:8200"

          listener "tcp" {
            tls_disable = 1
            address     = "0.0.0.0:8200"
          }

          storage "file" {
            path = "/vault/data"
          }
        ```

        This can also be set using Helm's `--set` flag (vault-helm v0.1.0 and later), using the following syntax:

        ```shell
        --set 'server.standalone.config='{ listener "tcp" { address = "0.0.0.0:8200" }'
        ```

  * <a name="v-server-ha" href="#v-server-ha">`ha`</a> - This configures `ha` mode for the Vault server.

      - <a name="v-server-ha-enabled" href="#v-server-ha-enabled">`enabled`</a> (`boolean: false`) -
      Enables `ha` mode for the Vault server.  This mode uses a highly available backend storage (such as Consul) to store Vault's data.  By default this is configured to use Consul Helm: https://github.com/hashicorp/consul-helm.  For a complete list of storage backends, see the official documentation: https://www.vaultproject.io/docs/configuration/storage/index.html.

      - <a name="v-server-ha-replicas" href="#v-server-ha-replicas">`replicas`</a> (`int: 5`) -
      The number of pods to deploy to create a highly available cluster of Vault servers.

      - <a name="v-server-ha-updatepartition" href="#v-server-ha-updatepartition">`updatePartition`</a> (`int: 0`) -
      If an updatePartition is specified, all Pods with an ordinal that is greater than or equal to the partition will be updated when the StatefulSet’s `.spec.template` is updated.  If set to `0`, this disables parition updates.  For more information see the official documentation: https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#rolling-updates

      - <a name="v-server-ha-config" href="#v-server-ha-config">`config`</a> (`string: "{}"`) -
      A raw string of extra HCL or JSON [configuration](https://www.vaultproject.io/docs/configuration/index.html) for Vault servers.
      This will be saved as-is into a ConfigMap that is read by the Vault servers.
      This can be used to add additional configuration that isn't directly exposed by the chart.

        ```yaml
        # ExtraConfig values are formatted as a multi-line string:
        config: |
          ui = true
          api_addr = "http://POD_IP:8200"
          listener "tcp" {
              tls_disable = 1
              address     = "0.0.0.0:8200"
          }

          storage "consul" {
              path = "vault"
              address = "HOST_IP:8500"
          }
        ```

        This can also be set using Helm's `--set` flag (vault-helm v0.1.0 and later), using the following syntax:

        ```shell
        --set 'server.ha.config='{ listener "tcp" { address = "0.0.0.0:8200" }'
        ```

      - <a name="v-server-ha-disruptionbudget" href="#v-server-ha-disruptionbudget">`disruptionBudget`</a> - Values that configures the disruption budget policy:  https://kubernetes.io/docs/tasks/run-application/configure-pdb/.

           - <a name="v-server-ha-disruptionbudget-enabled" href="#v-server-ha-disruptionbudget-enabled">`enabled`</a> (`boolean: true`) -
           Enables disruption budget policy to limit the number of pods that are down simultaneously from voluntary disruptions.

           - <a name="v-server-ha-disruptionbudget-maxunavailable" href="#v-server-ha-disruptionbudget-maxunavailable">`maxUnavailable`</a> (`int: null`) -
           The maximum number of unavailable pods. By default, this will be automatically
           computed based on the `server.replicas` value to be `(n/2)-1`. If you need to set
           this to `0`, you will need to add a `--set 'server.disruptionBudget.maxUnavailable=0'`
           flag to the helm chart installation command because of a limitation in the Helm
           templating language.

* <a name="v-ui" href="#v-ui">`ui`</a> - Values that configure the Vault UI.

  - <a name="v-ui-enabled" href="#v-ui-enabled">`enabled`</a> (`boolean: false`) - If true, the UI will be enabled. The UI will only be enabled on Vault servers. If `server.enabled` is false, then this setting has no effect. To expose the UI in some way, you must configure `ui.service`.

  - <a name="v-ui-servicetype" href="#v-ui-servicetype">`serviceType`</a> (`string: ClusterIP`) -
  The service type to register. This defaults to `ClusterIP`.
  The available service types are documented on
  [the Kubernetes website](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types).

## Using the Helm Chart to Deploy Vault Enterprise

You can also use this Helm chart to deploy Vault Enterprise by following a few extra steps.

Find the license file that you received in your welcome email. It should have the extension `.hclic`. You will use the contents of this file to install the license in Vault.

In your `values.yaml`, change the value of `global.image` to one of the enterprise [release tags](https://hub.docker.com/r/hashicorp/vault-enterprise/tags).

```yaml
global:
  image: "hashicorp/vault-enterprise:1.2.0-beta2"
```

Next, to install the license, the following requirements must be satisfied:
* Vault is initialized
* Vault is unsealed

-> **Important:** The Helm chart will not auto-initialize and unseal the cluster.  This must be done manually after running the Vault installation.  See the [initialization documentation](https://www.vaultproject.io/docs/commands/operator/init.html) and the [unseal documentation](https://www.vaultproject.io/docs/commands/operator/unseal.html) for more information.

Once initialized and unsealed, run the following to install the license:

First, setup a port-forward tunnel to the Vault cluster:

```bash
$ kubectl port-forward <NAME OF VAULT POD> 8200:8200
```

Next, in a separate terminal, create a `payload.json` file that contains the license key like this example:

```json
{
  "text": "01ABCDEFG..."
}
```

Finally, make an HTTP request to Vault API with the license key:

```bash
$ curl \
    --header "X-Vault-Token: VAULT_LOGIN_TOKEN_HERE" \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8200/v1/sys/license

```

To verify that the license installation worked correctly, make an HTTP request to the Vault API:

```bash
$ curl \
    --header "X-Vault-Token: VAULT_LOGIN_TOKEN_HERE" \
    http://127.0.0.1:8200/v1/sys/license
```

## Helm Chart Examples

The following are different configuration examples to support a variety of 
deployment models.

### Standalone Server with Load Balanced UI

The below `values.yaml` can be used to set up a single server Vault cluster with a LoadBalancer to allow external access to the UI and API.

```
global:
  enabled: true
  image: "vault:1.2.0-beta2"
 
server:
  standalone:
    enabled: true
    config: |
      api_addr = "http://POD_IP:8200"
      listener "tcp" {
        tls_disable = true
        address     = "0.0.0.0:8200"
      }

      storage "file" {
        path = "/vault/data"
      }

  service:
    enabled: true

  dataStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce
    
ui:
  enabled: true
  serviceType: LoadBalancer
```

### Standalone Server with TLS

The below `values.yaml` can be used to set up a single server Vault cluster using TLS.  
This assumes that a Kubernetes `secret` exists with the server certificate, key and 
certificate authority:

```
global:
  enabled: true
  image: "vault:1.2.0-beta2"

server:
  extraVolumes:
  - type: secret
    name: vault-server-tls

  extraEnvironmentVars:
    VAULT_ADDR: "https://localhost:8200"

  standalone:
    enabled: true
    config: |
      api_addr = "https://POD_IP:8200"
      listener "tcp" {
        tls_cert_file = "/vault/userconfig/vault-server-tls/vault.crt"
        tls_key_file  = "/vault/userconfig/vault-server-tls/vault.key"
        tls_client_ca_file = "/vault/userconfig/vault-server-tls/vault.ca"
        address     = "0.0.0.0:8200"
      }

      storage "file" {
        path = "/vault/data"
      }

  service:
    enabled: true

  dataStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce
```

### Standalone Server with Audit Storage

The below `values.yaml` can be used to set up a single server Vault cluster with 
auditing enabled.

```
global:
  enabled: true
  image: "vault:1.2.0-beta2"

server:
  standalone:
    enabled: true
    config: |
      api_addr = "http://POD_IP:8200"
      listener "tcp" {
        tls_disable = true
        address     = "0.0.0.0:8200"
      }

      storage "file" {
        path = "/vault/data"
      }

  service:
    enabled: true

  dataStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce

  auditStorage:
    enabled: true
    size: 10Gi
    storageClass: null
    accessMode: ReadWriteOnce
```

After Vault has been deployed, initialized and unsealed, auditing can be enabled 
by running the following command against the Vault pod:

```bash
$ kubectl exec -ti <POD NAME> --  vault audit enable file file_path=/vault/audit/vault_audit.log
```

### Highly Available Vault Cluster with Consul

The below `values.yaml` can be used to set up a five server Vault cluster using 
Consul as a highly available storage backend, Google Cloud KMS for Auto Unseal.

```
global:
  enabled: true
  image: "vault:1.2.0-beta2"

server:
  extraEnvironmentVars: {}
    GOOGLE_REGION: global,
    GOOGLE_PROJECT: myproject,
    GOOGLE_CREDENTIALS: /vault/userconfig/my-gcp-iam/myproject-creds.json

  extraVolumes: []
    - type: secret
      name: my-gcp-iam
      load: false

  affinity: |
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchLabels:
              app: {{ template "vault.name" . }}
              release: "{{ .Release.Name }}"
              component: server
          topologyKey: kubernetes.io/hostname

  service:
    enabled: true
    
  ha:
    enabled: false
    replicas: 5

    config: |
      ui = true
      api_addr = "http://POD_IP:8200"
      listener "tcp" {
        tls_disable = 1
        address     = "0.0.0.0:8200"
      }
      storage "consul" {
        path = "vault"
        address = "HOST_IP:8500"
      }

      seal "gcpckms" {
         project     = "myproject"
         region      = "global"
         key_ring    = "vault-unseal-kr"
         crypto_key  = "vault-unseal-key"
      }
```
