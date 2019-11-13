---
layout: "docs"
page_title: "Helm - Kubernetes"
sidebar_current: "docs-platform-k8s-helm"
sidebar_title: "Helm Chart"
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
out of the scope of this document. Please refer to the
[Helm documentation](https://helm.sh/) for more information.

Example chart usage:

```sh
# Clone the chart repo
$ git clone https://github.com/hashicorp/vault-helm.git
$ cd vault-helm

# Checkout a tagged version
$ git checkout v0.2.1

# Run Helm
$ helm install --dry-run ./
```

## Configuration (Values)

The chart is highly customizable using
[Helm configuration values](https://docs.helm.sh/using_helm/#customizing-the-chart-before-installing).
Each value has a default tuned for an optimal getting started experience
with Vault. Before going into production, please review the parameters below
and consider if they're appropriate for your deployment.

* `global` - These global values affect multiple components of the chart.

  * `enabled` (`boolean: true`) - The master enabled/disabled configuration. If this is true, most components will be installed by default. If this is false, no components will be installed by default and manually opting-in is required, such as by setting `server.enabled` to true.

  * `image` (`string: "vault:latest"`) - The name of the Docker image (including any tag) for the containers running Vault. **This should be pinned to a specific version when running in production.** Otherwise, other changes to the chart may inadvertently upgrade your Vault version.
  
  * `imagePullPolicy` (`string: "IfNotPresent"`) - The pull policy for container images.  The default pull policy is `IfNotPresent` which causes the Kubelet to skip pulling an image if it already exists.
  
  * `imagePullSecrets` (`string: ""`) - Defines secrets to be used when pulling images from private registries.

      - `name`: (`string: required`) - 
      Name of the secret containing files required for authentication to private image registries.

  * `tlsDisable` (`boolean: true`) - When set to `true`, changes URLs from `https` to `http` (such as the `VAULT_ADDR=http://127.0.0.1:8200` environment variable set on the Vault pods).

* `server` - Values that configure running a Vault server within Kubernetes.

  * `resources` (`string: null`) - The resource requests and limits (CPU, memory, etc.) for each of the server. This should be a multi-line string mapping directly to a Kubernetes [ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#resourcerequirements-v1-core) object. If this isn't specified, then the pods won't request any specific amount of resources. **Setting this is highly recommended.**

    ```yaml
    # Resources are defined as a formatted multi-line string:
    resources: |
      requests:
        memory: "10Gi"
      limits:
        memory: "10Gi"
    ```
  
  * `ingress` - Values that configure Ingress services for Vault.

    - `enabled` (`boolean: false`) - When set to `true`, an [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) service will be created.
    
    - `annotations` (`string`) - This value defines additional annotations to add to the Ingress service.  This should be formatted as a multi-line string.

        ```yaml
        annotations: |
          kubernetes.io/ingress.class: nginx
          kubernetes.io/tls-acme: "true"
        ```
    * `hosts` - Values that configure the Ingress host rules.

      - `host`: Name of the host to use for Ingress.

      - `paths`: This value defines the types of host rules for the Ingress service.

          ```yaml
          paths:
          - backend:
            serviceName: service2
            servicePort: 80
          ``` 

    * `tls` - Values that configure the Ingress TLS rules.

      - `hosts`: Name of the hosts defined in the Common Name of the TLS Certificate.  This should be formated as a multi-line string.

      - `secretName`: Name of the secret containing the required TLS files such as certificates and keys.

        ```yaml
        hosts:
          - sslexample.foo.com
          - sslexample.bar.com
         secretName: testsecret-tls
        ```

  * `authDelegator` - Values that configure the Cluster Role Binding attached to the Vault service account.

    - `enabled` (`boolean: false`) - When set to `true`, a Cluster Role Binding will be bound to the Vault service account.  This Cluster Role Binding has the necessary privileges for Vault to use the [Kubernetes Auth Method](/docs/auth/kubernetes.html).

  * `extraEnvironmentVars` (`string: null`) - The extra environment variables to be applied to the Vault server.  This should be a multi-line key/value string.

    ```yaml
    # Extra Environment Variables are defined as key/value strings.
     extraEnvironmentVars:
       GOOGLE_REGION: global,
       GOOGLE_PROJECT: myproject,
       GOOGLE_CREDENTIALS: /vault/userconfig/myproject/myproject-creds.json
    ```

  * `extraSecretEnvironmentVars` (`string: null`) - The extra environment variables populated from a secret to be applied to the Vault server.  This should be a multi-line key/value string.

      - `envName` (`string: required`) -
      Name of the environment variable to be populated in the Vault container.
      
      - `secretName` (`string: required`) -
      Name of Kubernetes secret used to populate the environment variable defined by `envName`.

      - `secretKey` (`string: required`) -
      Name of the key where the requested secret value is located in the Kubernetes secret.

    ```yaml
    # Extra Environment Variables populated from a secret.
     extraSecretEnvironmentVars:
      - envName: AWS_SECRET_ACCESS_KEY
        secretName: vault
        secretKey: AWS_SECRET_ACCESS_KEY
    ```

  * `extraVolumes` (`array: []`) - A list of extra volumes to mount to Vault servers. This is useful for bringing in extra data that can be referenced by other configurations at a well known path, such as TLS certificates. The value of this should be a list of objects. Each object supports the following keys:

      - `type` (`string: required`) -
      Type of the volume, must be one of "configMap" or "secret". Case sensitive.

      - `name` (`string: required`) -
      Name of the configMap or secret to be mounted. This also controls the path
      that it is mounted to. The volume will be mounted to `/vault/userconfig/<name>` by default
      unless `path` is configured.
      
      - `path` (`string: /vault/userconfigs`) -
      Name of the path where a configMap or secret is mounted.  If not specified 
      the volume will be mounted to `/vault/userconfig/<name of volume>`.

        ```yaml
        extraVolumes:
          -  type: "secret"
             name: "vault-certs"
             path: "/etc/pki"
        ```

  * `affinity` (`string`) - This value defines the [affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity) for server pods. It defaults to allowing only a single pod on each node, which minimizes risk of the cluster becoming unusable if a node is lost. If you need to run more pods per node (for example, testing on Minikube), set this value to `null`.

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

  * `tolerations` (`array []`) - This value defines the [tolerations](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/) that are acceptable when being scheduled.

        ```yaml
        tolerations:
        - key: "node.kubernetes.io/unreachable"
          operator: "Exists"
          effect: "NoExecute"
          tolerationSeconds: 6000
        ```

  * `nodeSelector` (`string`) - This value defines additional node selection criteria for more control over where the Vault servers are deployed.

        ```yaml
        nodeSelector:
          disktype: ssd
        ```

  * `extraLabels` (`string`) - This value defines additional labels for server pods. This should be formatted as a multi-line string.

        ```yaml
        extraLabels: |
          "sample/label1": "foo"
          "sample/label2": "bar"
        ```

  * `annotations` (`string`) - This value defines additional annotations for server pods. This should be a formatted as a multi-line string.

        ```yaml
        annotations: |
          "sample/annotation1": "foo"
          "sample/annotation2": "bar"
        ```

  * `service` - Values that configure the Kubernetes service created for Vault.

    - `enabled` (`boolean: true`) - When set to `true`, a Kubernetes service will be created for Vault.

    - `clusterIP` (`string`) - ClusterIP controls whether an IP address (cluster IP) is attached to the Vault service within Kubernetes.  By default the Vault service will be given a Cluster IP address, set to `None` to disable.  When disabled Kubernetes will create a "headless" service.  Headless services can be used to communicate with pods directly through DNS instead of a round robin load balancer.
    
    - `port` (`int: 8200`) - Port on which Vault server is listening inside the pod.

    - `targetPort` (`int: 8200`) - Port on which the service is listening.

    - `annotations` (`string`) - This value defines additional annotations for the service. This should be formatted as a multi-line string.

        ```yaml
        annotations: |
          "sample/annotation1": "foo"
          "sample/annotation2": "bar"
        ```

 * `serviceAccount` - Values that configure the Kubernetes service account created for Vault.

    - `annotations` (`string`) - This value defines additional annotations for the service account. This should be formatted as a multi-line string.

        ```yaml
        annotations: |
          "sample/annotation1": "foo"
          "sample/annotation2": "bar"
        ```

  * `extraVolumes` - This configures the `Service` resource created for the Vault server.

      - `enabled` (`boolean: true`) -
      Enables a service to allow other pods running in Kubernetes to communicate with the Vault server.

  * `dataStorage` - This configures the volume used for storing Vault data when not using external storage such as Consul.

      - `enabled` (`boolean: true`) -
      Enables a persistent volume to be created for storing Vault data when not using an external storage service.

      - `size` (`string: 10Gi`) -
      Size of the volume to be created for Vault's data storage when not using an external storage service.

      - `storageClass` (`string: null`) -
      Name of the storage class to use when creating the data storage volume.

      - `accessMode` (`string: ReadWriteOnce`) -
      Type of access mode of the storage device.  See the [official Kubernetes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) for more information.

  * `auditStorage` - This configures the volume used for storing Vault's audit logs.  See the [Vault documentation](/docs/audit) for more information.

      - `enabled` (`boolean: true`) -
      Enables a persistent volume to be created for storing Vault's audit logs.

      - `size` (`string: 10Gi`) -
      Size of the volume to be created for Vault's audit logs.

      - `storageClass` (`string: null`) -
      Name of the storage class to use when creating the audit storage volume.

      - `accessMode` (`string: ReadWriteOnce`) -
      Type of access mode of the storage device.

  * `dev` - This configures `dev` mode for the Vault server.

      - `enabled` (`boolean: false`) -
      Enables `dev` mode for the Vault server.  This mode is useful for experimenting with Vault without needing to unseal.

        ~> **Security Warning:** Never, ever, ever run a "dev" mode server in production. It is insecure and will lose data on every restart (since it stores data in-memory). It is only made for development or experimentation.

  * `standalone` - This configures `standalone` mode for the Vault server.

      - `enabled` (`boolean: true`) -
      Enables `standalone` mode for the Vault server.  This mode uses the `file` storage backend and requires a volume for persistence (`dataStorage`).

      - `config` (`string: "{}"`) -
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
        --set server.standalone.config='{ listener "tcp" { address = "0.0.0.0:8200" }'
        ```

  * `ha` - This configures `ha` mode for the Vault server.

      - `enabled` (`boolean: false`) -
      Enables `ha` mode for the Vault server.  This mode uses a highly available backend storage (such as Consul) to store Vault's data.  By default this is configured to use [Consul Helm](https://github.com/hashicorp/consul-helm).  For a complete list of storage backends, see the [Vault documentation](/docs/configuration).

      - `replicas` (`int: 5`) -
      The number of pods to deploy to create a highly available cluster of Vault servers.

      - `updatePartition` (`int: 0`) -
      If an updatePartition is specified, all Pods with an ordinal that is greater than or equal to the partition will be updated when the StatefulSet’s `.spec.template` is updated.  If set to `0`, this disables parition updates.  For more information see the [official Kubernetes documentation](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#rolling-updates).

      - `config` (`string: "{}"`) -
      A raw string of extra HCL or JSON [configuration](/docs/configuration) for Vault servers.
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
        --set server.ha.config='{ listener "tcp" { address = "0.0.0.0:8200" }'
        ```

      - `disruptionBudget` - Values that configures the disruption budget policy.  See the [official Kubernetes documentation](https://kubernetes.io/docs/tasks/run-application/configure-pdb/) for more information.

           - `enabled` (`boolean: true`) -
           Enables disruption budget policy to limit the number of pods that are down simultaneously from voluntary disruptions.

           - `maxUnavailable` (`int: null`) -
           The maximum number of unavailable pods. By default, this will be automatically
           computed based on the `server.replicas` value to be `(n/2)-1`. If you need to set
           this to `0`, you will need to add a `--set 'server.disruptionBudget.maxUnavailable=0'`
           flag to the helm chart installation command because of a limitation in the Helm
           templating language.

* `ui` - Values that configure the Vault UI.

  - `enabled` (`boolean: false`) - If true, the UI will be enabled. The UI will only be enabled on Vault servers. If `server.enabled` is false, then this setting has no effect. To expose the UI in some way, you must configure `ui.service`.

  - `serviceType` (`string: ClusterIP`) -
  The service type to register. This defaults to `ClusterIP`.
  The available service types are documented on
  [the Kubernetes website](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types).

  - `serviceNodePort` (`int: null`) -
  Sets the Node Port value when using `serviceType: NodePort` on the Vault UI service.

  - `externalPort` (`int: 8200`) -
  Sets the external port value of the service.

  - `loadBalancerSourceRanges` (`string`) - This value defines additional source CIDRs when using `serviceType: LoadBalancer`.  This should be formatted as a multi-line string.

    ```yaml
    loadBalancerSourceRanges:
    - 10.0.0.0/16
    - 120.78.23.3/32
    ```

  - `loadBalancerIP` (`string`) - This value defines the IP address of the load balancer when using `serviceType: LoadBalancer`.
   
  - `annotations` (`string`) - This value defines additional annotations for the UI service. This should be a formatted as a multi-line string.

    ```yaml
    annotations: |
      "sample/annotation1": "foo"
      "sample/annotation2": "bar"
    ```

## Helm Chart Examples

The following are different configuration examples to support a variety of 
deployment models.

### Standalone Server with Load Balanced UI

The below `values.yaml` can be used to set up a single server Vault cluster with a LoadBalancer to allow external access to the UI and API.

```yaml
global:
  enabled: true
  image: "vault:1.2.4"
 
server:
  standalone:
    enabled: true
    config: |
      ui = true

      listener "tcp" {
        tls_disable = 1
        address = "[::]:8200"
        cluster_address = "[::]:8201"
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

```yaml
global:
  enabled: true
  image: "vault:1.2.4"
  tlsDisable: false

server:
  extraVolumes:
  - type: secret
    name: vault-server-tls

  standalone:
    enabled: true
    config: |
      listener "tcp" {
        address = "[::]:8200"
        cluster_address = "[::]:8201"
        tls_cert_file = "/vault/userconfig/vault-server-tls/vault.crt"
        tls_key_file  = "/vault/userconfig/vault-server-tls/vault.key"
        tls_client_ca_file = "/vault/userconfig/vault-server-tls/vault.ca"
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

```yaml
global:
  enabled: true
  image: "vault:1.2.4"

server:
  standalone:
    enabled: true
    config: |
      listener "tcp" {
        tls_disable = true
        address = "[::]:8200"
        cluster_address = "[::]:8201"
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

```yaml
global:
  enabled: true
  image: "vault:1.2.4"

server:
  extraEnvironmentVars:
    GOOGLE_REGION: global
    GOOGLE_PROJECT: myproject
    GOOGLE_APPLICATION_CREDENTIALS: /vault/userconfig/my-gcp-iam/myproject-creds.json

  extraVolumes: []
    - type: secret
      name: my-gcp-iam

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

      listener "tcp" {
        tls_disable = 1
        address = "[::]:8200"
        cluster_address = "[::]:8201"
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
