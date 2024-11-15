# triton-go

`triton-go` is a client SDK for Go applications using Joyent's Triton Compute
and Object Storage (Manta) APIs.

The Triton Go SDK is used in the following open source projects.

- [Consul](https://www.consul.io/docs/agent/cloud-auto-join.html#joyent-triton)
- [Packer](http://github.com/hashicorp/packer)
- [Vault](http://github.com/hashicorp/vault)
- [Terraform](http://github.com/hashicorp/terraform)
- [Terraform Triton Provider](https://github.com/terraform-providers/terraform-provider-triton)
- [Docker Machine](https://github.com/joyent/docker-machine-driver-triton)
- [Triton Kubernetes](https://github.com/joyent/triton-kubernetes)
- [HashiCorp go-discover](https://github.com/hashicorp/go-discover)

## Usage

Triton uses [HTTP Signature][4] to sign the Date header in each HTTP request
made to the Triton API. Currently, requests can be signed using either a private
key file loaded from disk (using an [`authentication.PrivateKeySigner`][5]), or
using a key stored with the local SSH Agent (using an [`SSHAgentSigner`][6].

To construct a Signer, use the `New*` range of methods in the `authentication`
package. In the case of `authentication.NewSSHAgentSigner`, the parameters are
the fingerprint of the key with which to sign, and the account name (normally
stored in the `TRITON_ACCOUNT` environment variable). There is also support for
passing in a username, this will allow you to use an account other than the main
Triton account. For example:

```go
input := authentication.SSHAgentSignerInput{
    KeyID:       "a4:c6:f3:75:80:27:e0:03:a9:98:79:ef:c5:0a:06:11",
    AccountName: "AccountName",
    Username:    "Username",
}
sshKeySigner, err := authentication.NewSSHAgentSigner(input)
if err != nil {
    log.Fatalf("NewSSHAgentSigner: %s", err)
}
```

An appropriate key fingerprint can be generated using `ssh-keygen`.

```
ssh-keygen -Emd5 -lf ~/.ssh/id_rsa.pub | cut -d " " -f 2 | sed 's/MD5://'
```

Each top level package, `account`, `compute`, `identity`, `network`, all have
their own separate client. In order to initialize a package client, simply pass
the global `triton.ClientConfig` struct into the client's constructor function.

```go
config := &triton.ClientConfig{
    TritonURL:   os.Getenv("TRITON_URL"),
    MantaURL:    os.Getenv("MANTA_URL"),
    AccountName: accountName,
    Username:    os.Getenv("TRITON_USER"),
    Signers:     []authentication.Signer{sshKeySigner},
}

c, err := compute.NewClient(config)
if err != nil {
    log.Fatalf("compute.NewClient: %s", err)
}
```

Constructing `compute.Client` returns an interface which exposes `compute` API
resources. The same goes for all other packages. Reference their unique
documentation for more information.

The same `triton.ClientConfig` will initialize the Manta `storage` client as
well...

```go
c, err := storage.NewClient(config)
if err != nil {
    log.Fatalf("storage.NewClient: %s", err)
}
```

## Error Handling

If an error is returned by the HTTP API, the `error` returned from the function
will contain an instance of `errors.APIError` in the chain. Error wrapping
is performed using the [pkg/errors][7] library.

## Acceptance Tests

Acceptance Tests run directly against the Triton API (e.g. CloudAPI), so you
will need an installation of Triton in order to run them, or
[COAL](https://github.com/joyent/triton/blob/master/docs/developer-guide/coal-setup.md)
(Cloud On A Laptop). The tests create real resources and thus could cost real
money if you are using a paid Triton account! It is also possible that the
acceptance tests will leave behind resources, so extra attention will be needed
to clean up these test resources.

The acceptance tests depend upon a few things in your Triton setup:

- the [Ubuntu 16.04](https://docs.joyent.com/public-cloud/instances/virtual-machines/images/linux/ubuntu-certified#1604-xenial-images) LX image to be installed
- a generic package with Memory in the range of 128MB to 1024MB
- a public (external) network to provision with
- if your Triton setup is used for testing - then running
  `sdcadm post-setup dev-sample-data` ensure these conditions are met

In order to run acceptance tests, the following environment variables must be
set:

- `TRITON_TEST` - must be set to any value in order to indicate desire to create
  resources
- `TRITON_URL` - the base endpoint for the Triton API
- `TRITON_ACCOUNT` - the account name for the Triton API
- `TRITON_KEY_ID` - the fingerprint of the SSH key identifying the key
- `TRITON_SKIP_TLS_VERIFY` - skips the checking of the server's TLS certificate.
  This is insecure and should only be used for internal/testing.

Additionally, you may set these optional environment variables:

- `TRITON_KEY_MATERIAL` to the contents of an unencrypted private key. If this
    is set, the PrivateKeySigner (see above) will be used - if not the
    SSHAgentSigner will be used. You can also set `TRITON_USER` to run the tests
    against an account other than the main Triton account.
- `TRITON_TRACE_HTTP` - when set this will print the HTTP requests and responses
  to stderr.
- `TRITON_VERBOSE_TESTS` - turns on extra logging for test runs

### Example Run

The verbose output has been removed for brevity here.

```
$ HTTP_PROXY=http://localhost:8888 \
    TRITON_TEST=1 \
    TRITON_URL=https://us-sw-1.api.joyent.com \
    TRITON_ACCOUNT=AccountName \
    TRITON_KEY_ID=a4:c6:f3:75:80:27:e0:03:a9:98:79:ef:c5:0a:06:11 \
    go test -v -run "TestAccKey"
=== RUN   TestAccKey_Create
--- PASS: TestAccKey_Create (12.46s)
=== RUN   TestAccKey_Get
--- PASS: TestAccKey_Get (4.30s)
=== RUN   TestAccKey_Delete
--- PASS: TestAccKey_Delete (15.08s)
PASS
ok  	github.com/joyent/triton-go	31.861s
```

## Example API

There's an `examples/` directory available with sample code setup for many of
the APIs within this library. Most of these can be run using `go run` and
referencing your SSH key file use by your active `triton` CLI profile.

```sh
$ eval "$(triton env us-sw-1)"
$ TRITON_KEY_FILE=~/.ssh/triton-id_rsa go run examples/compute/instances.go
```

The following is a complete example of how to initialize the `compute` package
client and list all instances under an account. More detailed usage of this
library follows.

```go


package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "time"

    triton "github.com/joyent/triton-go"
    "github.com/joyent/triton-go/authentication"
    "github.com/joyent/triton-go/compute"
)

func main() {
    keyID := os.Getenv("TRITON_KEY_ID")
    accountName := os.Getenv("TRITON_ACCOUNT")
    keyMaterial := os.Getenv("TRITON_KEY_MATERIAL")
    userName := os.Getenv("TRITON_USER")

    var signer authentication.Signer
    var err error

    if keyMaterial == "" {
        input := authentication.SSHAgentSignerInput{
            KeyID:       keyID,
            AccountName: accountName,
            Username:    userName,
        }
        signer, err = authentication.NewSSHAgentSigner(input)
        if err != nil {
            log.Fatalf("Error Creating SSH Agent Signer: {{err}}", err)
        }
    } else {
        var keyBytes []byte
        if _, err = os.Stat(keyMaterial); err == nil {
            keyBytes, err = ioutil.ReadFile(keyMaterial)
            if err != nil {
                log.Fatalf("Error reading key material from %s: %s",
                    keyMaterial, err)
            }
            block, _ := pem.Decode(keyBytes)
            if block == nil {
                log.Fatalf(
                    "Failed to read key material '%s': no key found", keyMaterial)
            }

            if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
                log.Fatalf(
                    "Failed to read key '%s': password protected keys are\n"+
                        "not currently supported. Please decrypt the key prior to use.", keyMaterial)
            }

        } else {
            keyBytes = []byte(keyMaterial)
        }

        input := authentication.PrivateKeySignerInput{
            KeyID:              keyID,
            PrivateKeyMaterial: keyBytes,
            AccountName:        accountName,
            Username:           userName,
        }
        signer, err = authentication.NewPrivateKeySigner(input)
        if err != nil {
            log.Fatalf("Error Creating SSH Private Key Signer: {{err}}", err)
        }
    }

    config := &triton.ClientConfig{
        TritonURL:   os.Getenv("TRITON_URL"),
        AccountName: accountName,
        Username:    userName,
        Signers:     []authentication.Signer{signer},
    }

    c, err := compute.NewClient(config)
    if err != nil {
        log.Fatalf("compute.NewClient: %s", err)
    }

    listInput := &compute.ListInstancesInput{}
    instances, err := c.Instances().List(context.Background(), listInput)
    if err != nil {
        log.Fatalf("compute.Instances.List: %v", err)
    }
    numInstances := 0
    for _, instance := range instances {
        numInstances++
        fmt.Println(fmt.Sprintf("-- Instance: %v", instance.Name))
    }
}

```

[4]: https://github.com/joyent/node-http-signature/blob/master/http_signing.md
[5]: https://godoc.org/github.com/joyent/triton-go/authentication
[6]: https://godoc.org/github.com/joyent/triton-go/authentication
[7]: https://github.com/pkg/errors
