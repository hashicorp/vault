---
layout: "docs"
page_title: "Agent Sidecar Injector Overview"
sidebar_current: "docs-platform-k8s-injector"
sidebar_title: "Agent Injector"
description: |-
  The Vault Agent Sidecar Injector is a Kubernetes admission webhook that adds Vault Agent containers to pods for consuming Vault secrets.
---

# Agent Sidecar Injector

The Vault Agent Injector alters pod specifications to include Vault Agent
containers that render Vault secrets to a shared memory volume using 
[Vault Agent Templates](/docs/agent/template/index.html).
By rendering secrets to a shared volume, containers within the pod can consume 
Vault secrets without being Vault aware.

The injector is a [Kubernetes Mutation Webhook Controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/).
The controller intercepts pod events and applies mutations to the pod if annotations exist within
the request.  This functionality is provided by the [vault-k8s](https://github.com/hashicorp/vault-k8s)
project and can be automatically installed and configured using the
[Vault Helm](https://github.com/hashicorp/vault-helm) chart.

## Overview

The Vault Agent Injector works by intercepting pod `CREATE` and `UPDATE`
events in Kubernetes.  The controller parses the event and looks for the metadata
annotation `vault.hashicorp.com/agent-inject: true`.  If found, the controller will
alter the pod specification based on other annotations present.

### Mutations

At a minimum, every container in the pod will be configured to mount a shared
memory volume.  This volume is mounted to `/vault/secrets` and will be used by the Vault
Agent containers for sharing secrets with the other containers in the pod.

Next, two types of Vault Agent containers can be injected: init and sidecar.  The
init container will prepopulate the shared memory volume with the requested
secrets prior to the other containers starting.  The sidecar container will
continue to authenticate and render secrets to the same location as the pod runs.
Using annotations, the initialization and sidecar containers may be disabled.

Last, two additional types of volumes can be optionally mounted to the Vault Agent
containers.  The first is secret volume containing TLS requirements such as client
and CA (certificate authority) certificates and keys. This volume is useful when
communicating and verifying the Vault server's authenticity using TLS.  The second
is a configuration map containing Vault Agent configuration files.  This volume is
useful to customize Vault Agent beyond what the provided annotations offer.

### Authenticating with Vault

The primary method of authentication with Vault when using the Vault Agent Injector
is the service account attached to the pod.  At this time, no other authentication
method is supported by the controller.

The service account must be bound to a Vault role and a policy granting access to
the secrets desired.

A service account must be present to use the Vault Agent Injector.  It is *not*
recommended to bind Vault roles to the default service account provided to pods
if no service account is defined.

### Requesting Secrets

There are two methods of configuring the Vault Agent containers to render secrets:

* the `vault.hashicorp.com/agent-inject-secret` annotation, or
* a configuration map containing Vault Agent configuration files.

Only one of these methods may be used at any time.

#### Secrets via Annotations

To configure secret injection using annotations, the user must supply:

* one or more _secret_ annotations, and
* the Vault role used to access those secrets. 

The annotation must have the format:

```yaml
vault.hashicorp.com/agent-inject-secret-<unique-name>: /path/to/secret
```

The unique name will be the filename of the rendered secret and must be unique if
multiple secrets are defined by the user.  For example, consider the following
secret annotations:

```yaml
vault.hashicorp.com/agent-inject-secret-foo: database/roles/app
vault.hashicorp.com/agent-inject-secret-bar: consul/creds/app
vault.hashicorp.com/role: "app"
```

The first annotation will be rendered to `/vault/secrets/foo` and the second
annotation will be rendered to `/vault/secrets/bar`.

It's possible to set the file format of the using the annotation.  For example the
following secret will be rendered to `/vault/secrets/foo.txt`:

```yaml
vault.hashicorp.com/agent-inject-secret-foo.txt: database/roles/app
vault.hashicorp.com/role: "app"
```

The secret unique name must consist of alphanumeric characters, `.`, `_` or `-`.

##### Secret Templates

~> Vault Agent uses the Consul Template project to render secrets.  For more information
  on writing templates, see the [Consul Template documentation](https://github.com/hashicorp/consul-template).

How the secret is rendered to the file is also configurable.  To configure the template
used, the user must supply a _template_ annotation using the same unique name of
the secret.  The annotation must have the following format:

```yaml
vault.hashicorp.com/agent-inject-template-<unique-name>: |
  < 
    TEMPLATE
    HERE
  >
```

For example, consider the following:

```yaml
vault.hashicorp.com/agent-inject-secret-foo: "database/roles/app"
vault.hashicorp.com/agent-inject-template-foo: |
  {{- with secret "database/creds/db-app" -}}
  postgres://{{ .Data.username }}:{{ .Data.password }}@postgres:5432/mydb?sslmode=disable
  {{- end }}
vault.hashicorp.com/role: "app"
```

The rendered secret would look like this within the container:

```bash
$ cat /vault/secrets/foo
postgres://v-kubernet-pg-app-q0Z7WPfVN:A1a-BUEuQR52oAqPrP1J@postgres:5432/mydb?sslmode=disable
```

~> The default left and right template delimiters are `{{` and `}}`.

If no template is provided the following generic template is used:

```
{{ with secret "/path/to/secret" }}
    {{ range $k, $v := .Data }}
        {{ $k }}: {{ $v }}
    {{ end }}
{{ end }}
```

For example, the following annotation will use the default template to render
PostgreSQL secrets found at the configured path:

```yaml
vault.hashicorp.com/agent-inject-secret-foo: "database/roles/pg-app"
vault.hashicorp.com/role: "app"
```

The rendered secret would look like this within the container:

```bash
$ cat /vault/secrets/foo
password: A1a-BUEuQR52oAqPrP1J
username: v-kubernet-pg-app-q0Z7WPfVNqqTJuoDqCTY-1576529094
```

#### Vault Agent Configuration Map

For advanced use cases, it may be required to define Vault Agent configuration
files to mount instead of using secret and template annotations.  The Vault Agent 
Injector supports mounting ConfigMaps by specifying the name using the `vault.hashicorp.com/agent-configmap`
annotation.  The configuration files will be mounted to `/vault/configs`. 

The configuration map must contain either one or both of the following files:

* `config-init.hcl` used by the init container.  This must have `exit_after_auth` set to `true`.
* `config.hcl` used by the sidecar container.  This must have `exit_after_auth` set to `false`.

An example of mounting a Vault Agent configmap [can be found here](/docs/platforms/k8s/injector/examples.html#configmap-example).

## Annotations

* `vault.hashicorp.com/agent-inject` - configures whether injection is explicitly
  enabled or disabled for a pod. This should be set to a `true` or `false` value.
  Defaults to `false`.

* `vault.hashicorp.com/agent-inject-status` - blocks further mutations
  by adding the value `injected` to the pod after a successful mutation.

* `vault.hashicorp.com/agent-inject-secret` - configures Vault Agent
  to retrieve the secrets from Vault required by the container. The name of the
  secret is any unique string after `vault.hashicorp.com/agent-inject-secret-`,
  such as `vault.hashicorp.com/agent-inject-secret-foobar`. The value is the path
  in Vault where the secret is located.

* `vault.hashicorp.com/agent-inject-template` - configures Vault Agent
  what template to use for rendering the secrets.  The name of the template is any
  unique string after `vault.hashicorp.com/agent-inject-template-`, such as 
  `vault.hashicorp.com/agent-inject-template-foobar`. This should map to the same
  unique value provided in `vault.hashicorp.com/agent-inject-secret-`. If not provided,
  a default generic template is used.

* `vault.hashicorp.com/role` - configures the Vault role used by the Vault Agent
  auto-auth method.  Required when `vault.hashicorp.com/agent-configmap` is not set.

* `vault.hashicorp.com/agent-configmap` - name of the configuration map where Vault
  Agent configuration file and templates can be found.

* `vault.hashicorp.com/agent-pre-populate` - configures whether an init container
  is included to pre-populate the shared memory volume with secrets prior to the
  containers starting.

* `vault.hashicorp.com/agent-pre-populate-only` - configures whether an init container
  is the only injected container. If true, no sidecar container will be injected
  at runtime of the pod.

* `vault.hashicorp.com/agent-image` - name of the Vault docker image to use. This
  value overrides the default image configured in the controller and is usually
  not needed.  Defaults to `vault:1.3.1`.

* `vault.hashicorp.com/service` - name of the Vault service to use.  This value
  overrides the default service configured in the controller and is usually not
  needed.

* `vault.hashicorp.com/agent-limits-cpu` - configures the CPU limits on the Vault
  Agent containers.  Defaults to `500m`.

* `vault.hashicorp.com/agent-limits-mem` - configures the memory limits on the Vault
  Agent containers.  Defaults to `128Mi`.

* `vault.hashicorp.com/agent-requests-cpu` - configures the CPU requests on the
  Vault Agent containers.  Defaults to `250m`.

* `vault.hashicorp.com/agent-requests-mem` - configures the memory requests on the
  Vault Agent containers.  Defaults to `64Mi`.

* `vault.hashicorp.com/tls-secret` - name of the Kubernetes secret containing TLS
  Client and CA certificates and keys.  This is mounted to `/vault/tls`.

* `vault.hashicorp.com/tls-server-name` - name of the Vault server to verify the
  authenticity of the server when communicating with Vault over TLS.

* `vault.hashicorp.com/tls-skip-verify` - if true, configures the Vault Agent to
  skip verification of Vault's TLS certificate.  It's not recommended to set this
  value to true in a production environment.

* `vault.hashicorp.com/ca-cert` - path of the CA certificate used to verify Vault's
  TLS.

* `vault.hashicorp.com/ca-key` - path of the CA public key used to verify Vault's
  TLS.

* `vault.hashicorp.com/client-cert` - path of the client certificate used when
  communicating with Vault via mTLS.

* `vault.hashicorp.com/client-key` - path of the client public key used when communicating
  with Vault via mTLS.

* `vault.hashicorp.com/client-max-retries` - configures number of Vault Agent retry
  attempts when 5xx errors are encountered.  Defaults to 2.

* `vault.hashicorp.com/client-timeout` - configures the request timeout threshold,
  in seconds, of the Vault Agent when communicating with Vault.  Defaults to `60s` 
  and accepts value types of `60`, `60s` or `1m`.
