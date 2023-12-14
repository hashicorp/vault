---
layout: docs
page_title: Add a containerized secrets plugin
description: >-
  Add a containerized secrets plugin to your Vault instance.
---

## Add a containerized secrets plugin to Vault

@include 'alerts/beta.mdx'

Run your external secrets plugins in containers to increases the isolation
between the plugin and Vault.

## Before you start

- **You must be running Vault 1.15.0+**.
- **Your Vault server must be running on Linux**.

## Step 1: Install a container engine

If you do not have a container engine available, install one of the supported
container engines:

  - [Docker](https://docs.docker.com/engine/install/) or
    [Rootless Docker](https://docs.docker.com/engine/security/rootless/)
  - [Podman](https://podman.io/docs/installation#installing-on-linux) or
    [Rootless Podman](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md)

## Step 2: Install gVisor

Vault uses the `runsc` runtime from
[gVisor](https://gvisor.dev/docs/user_guide/install/) as the entrypoint to your
container runtime. If you want to use a container runtime other than gVisor, it
must be compatible with `runsc`.


## Step 3: Configure the container runtime 

Update your container engine to use `runsc` for Unix sockets between the host
and plugin binary.

<Tabs>

<Tab heading="Docker">

  1. Install `runsc` as a Docker runtime:
     ```shell-session
     $ sudo runsc install
     ```

  1. Confirm the `runsc` entry in your
     [Docker daemon configuration](https://docs.docker.com/config/daemon) file
     (`/etc/docker/daemon.json`):
      ```json
      {
        "runtimes": {
          "runsc": {
            "path": " /usr/bin/runsc",
            "runtimeArgs": [
              "--host-uds=all"
            ]
          }
        }
      }
      ```

  1. Restart Docker:
      ```shell-session
      sudo systemctl reload docker
      ```

</Tab>

<Tab heading="Rootless Docker">

  1. Install `runsc` as a Docker runtime:
     ```shell-session
     $ sudo runsc install
     ```

  1. Create a configuration directory if it does not exist already:
      ```shell-session
      mkdir -p ~/.config/docker
      ```
  1. Confirm the `runsc` entry in your
     [Docker daemon configuration](https://docs.docker.com/config/daemon) file
     (`~/.config/docker/daemon.json`):
      ```json
      {
        "runtimes": {
          "runsc": {
            "path": /usr/bin/runsc",
            "runtimeArgs": [
              "--host-uds=all"
              "--ignore-cgroups"
            ]
          }
        }
      }
      ```

  1. Restart Docker:
      ```shell-session
      systemctl --user restart docker
      ```

</Tab>

<Tab heading="Podman">

  1. Create an executable script to configure the OCI runtime flags.
     For example:
    ```shell-session
    sudo tee /usr/local/bin/runsc.podman <<EOF
    #!/bin/bash
    /usr/local/bin/runsc --host-uds=all "$\@"
    EOF
    ```

  1. Grant execute permission to `runsc` for Podman:
      ```shell-session
      chmod a+x /usr/local/bin/runsc.podman
      ```

  1. Start the [Docker-compatible Podman API](https://docs.podman.io/en/latest/_static/api.html)
      ```shell-session
      podman --runtime=/usr/local/bin/runsc.podman system service -t 0 &
      ```

</Tab>

<Tab heading="Rootless Podman">

  1. Create a local `bin` directory if it does not exist already:
    ```shell-session
      mkdir -p "$HOME/local/bin"
    ```

  1. Create an executable script to configure the OCI runtime flags:
    ```shell-session
      tee ~/local/bin/runsc.podman <<EOF
      #!/bin/bash
      /usr/local/bin/runsc --host-uds=all --ignore-cgroups "\$@"
      EOF
    ```

  1. Grant execute permission to `runsc` for Podman:
      ```shell-session
        chmod u+x "$HOME/local/bin/runsc.podman"
      ```

  1. Start the [Docker-compatible Podman API](https://docs.podman.io/en/latest/_static/api.html)
      ```shell-session
        podman --runtime="$HOME/local/bin/runsc.podman" system service -t 0 &
      ```

</Tab>

</Tabs>

## Step 3: Build the plugin container

Containerized plugins should run as a binary in the finished container and
behave the same whether run in a container or as a standalone application:

1. Build your plugin locally with v1.5.0+ of the HashiCorp
   [`go-plugin`](https://github.com/hashicorp/go-plugin) library to ensure the
   finished binary is compatible with containerization on Linux.

1. Create a container file for your plugin with the compiled binary as the
   entry-point.

1. Build the image with a unique tag.

<Tip title="The Vault SDK includes go-plugin">

  If you build plugins with the Vault Go SDK, you can update the `go-plugin`
  library by pulling the latest SDK version from the `hashicorp/vault` repo:

  `go install github.com/hashicorp/vault/sdk@latest`

</Tip>

For example, to build a containerized version of the built-in key-value (KV)
secrets plugin for Docker:

1. Install `go` so you can build the Go binary:
   ```shell-session
   $ sudo apt install golang-go
   ```

1. Clone the latest version of the KV secrets plugin from
   `hashicorp/vault-plugin-secrets-kv`:
    ```shell-session
    $ git clone https://github.com/hashicorp/vault-plugin-secrets-kv.git
    ```

1. Build the Go binary for Linux and create an empty Dockerfile under
   `vault-plugin-secrets-kv`:
    ```shell-session
    $ cd vault-plugin-secrets-kv ; CGO_ENABLED=0 GOOS=linux \
    go build -o kv cmd/vault-plugin-secrets-kv/main.go ; touch Dockerfile
    ```

1. Update the empty a `Dockerfile` with the infrastructure build details and the
   compiled binary as the entry-point:
   ```Dockerfile
   FROM <YOUR_LINUX_PLATFORM>
   COPY kv /bin/kv
   ENTRYPOINT [ "/bin/kv" ]
   ```
   For example:

   <CodeBlockConfig hideClipboard>
   
   ```Dockerfile
   FROM ubuntu
   COPY kv /bin/kv
   ENTRYPOINT [ "/bin/kv" ]
   ```
   
   </CodeBlockConfig>

1. Build a Docker image with the tag `kv-container`:
   ```shell-session
   $ docker build -t hashicorp/vault-plugin-secrets-kv:kv-container .
   ```

## Step 5: Register the plugin

Registering a containerized plugin with Vault is similar to registering any
other external plugins as long as the containers are available locally.

1. Store the SHA256 of the plugin image:
   ```shell-session
   export SHA256=$(
     docker images          \
       --no-trunc           \
       --format="{{ .ID }}" \
       <YOUR_DOCKER_IMAGE>:<TAG> | cut -d: -f2
   )
   ```
   For example:
   
   <CodeBlockConfig hideClipboard>

   ```shell-session
   $ export SHA256=$(
     docker images          \
       --no-trunc           \
       --format="{{ .ID }}" \
       hashicorp/vault-plugin-secrets-kv:kv-container | cut -d: -f2
   )
   ```
   
   </CodeBlockConfig>

1. Register the plugin with `vault plugin register` and specify your plugin
   image with the `oci_image` flag:
   ```shell-session
   $ vault plugin register            \
       -sha256="${SHA256}"            \
       -oci_image=<YOUR_PLUGIN_IMAGE> \
       -version=<PLUGIN_VERSION>      \
       <NEW_PLUGIN_TYPE> <NEW_PLUGIN_ID>
   ```
   For example:
   
   <CodeBlockConfig hideClipboard>
   
   ```shell-session
   $ vault plugin register                          \
       -sha256="${SHA256}"                          \
       -oci_image=hashicorp/vault-plugin-secrets-kv \
       -version="kv-container"                      \
       secret my-kv-plugin
   Success! Registered plugin: my-kv-plugin
   ```
   
   </CodeBlockConfig>

1. Enable the new plugin for your Vault instance with `vault secrets enable` and
   the new plugin ID:
   ```shell-session
    $ vault secrets enable <NEW_PLUGIN_ID>
   ```
   For example:
   
   <CodeBlockConfig hideClipboard>
   
   ```shell-session
   $ vault secrets enable my-kv-plugin
   ```
        
   </CodeBlockConfig>

<Tip title="Customize container behavior with registration flags">

  You can provide additional information about the image entrypoint, command,
  and environment with the `-command`, `-args`, and `-env` flags for
  `vault plugin register`.

</Tip>

## Step 6: Test your plugin

Now that the container is registered with Vault, you should be able to interact
with it like any other plugin. Try writing then fetching a new secret with your
new plugin.


1. Use `vault write` to store a secret with your containerized plugin:
   ```shell-session
   $ vault write <NEW_PLUGIN_ID>/<SECRET_PATH> <SECRET_KEY>=<SECRET_VALUE>
   ```
   For example:
   
   <CodeBlockConfig hideClipboard>
   
   ```shell-session
   $ vault write my-kv-plugin/testing subject=containers
   Success! Data written to: my-kv-plugin/testing
   ```
    
   </CodeBlockConfig>

1. Fetch the secret you just wrote:
   ```shell-session
   $ vault read <NEW_PLUGIN_ID>/<SECRET_PATH>
   ```
   For example:
   
   <CodeBlockConfig hideClipboard>
   
   ```shell-session
   $ vault read my-kv-plugin/testing
   ===== Data =====
   Key        Value
   ---        -----
   subject    containers
   ```
   
   </CodeBlockConfig>
