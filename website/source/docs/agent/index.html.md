---
layout: "docs"
page_title: "Vault Agent"
sidebar_title: "Vault Agent"
sidebar_current: "docs-agent"
description: |-
  Vault Agent is a client-side daemon that can be used to perform some Vault
  functionality automatically.
---

# Vault Agent

Vault Agent is a client daemon that can perform useful tasks.

To get help, run:

```text
$ vault agent -h
```
## Auto-Auth

Vault Agent allows for easy authentication to Vault in a wide variety of
environments. Please see the [Auto-Auth docs](/docs/agent/autoauth/index.html)
for information.

Auto-Auth functionality takes place within an `auto_auth` configuration stanza.

## Configuration

These are the currently-available general configuration option:

- `pid_file` `(string: "")` - Path to the file in which the agent's Process ID
  (PID) should be stored

- `exit_after_auth` `(bool: false)` - If set to `true`, the agent will exit
  with code `0` after a single successful auth, where success means that a
  token was retrieved and all sinks successfully wrote it

## Example Configuration

An example configuration, with very contrived values, follows:

```python
pid_file = "./pidfile"

auto_auth {
        method "aws" {
                mount_path = "auth/aws-subaccount"
                config = {
                        type = "iam"
                        role = "foobar"
                }
        }

        sink "file" {
                config = {
                        path = "/tmp/file-foo"
                }
        }

        sink "file" {
                wrap_ttl = "5m" 
                aad_env_var = "TEST_AAD_ENV"
                dh_type = "curve25519"
                dh_path = "/tmp/file-foo-dhpath2"
                config = {
                        path = "/tmp/file-bar"
                }
        }
}
```
