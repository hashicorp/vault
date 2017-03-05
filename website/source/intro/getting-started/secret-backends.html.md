---
layout: "intro"
page_title: "Secret Backends"
sidebar_current: "gettingstarted-secretbackends"
description: |-
  Secret backends are what create, read, update, and delete secrets.
---

# Secret Backends

Previously, we saw how to read and write arbitrary secrets to Vault.  To
do this, we used the `secret/` prefix. This prefix specifies which
backend to use. By default, Vault mounts a backend called _generic_ to
`secret/`. The generic backend reads and writes raw data to the backend
storage.

Vault supports other backends in addition to the _generic_ backend, and this feature
in particular is what makes Vault unique. For example, the _aws_ backend
generates AWS access keys dynamically, on demand. Another example --
this type of backend does not yet exist -- is a backend that
reads and writes data directly to an
[HSM](https://en.wikipedia.org/wiki/Hardware_security_module).
As Vault matures, more and more backends will be added.

To represent backends, Vault behaves much like a filesystem: backends
are mounted at specific paths. For example, the _generic_ backend is
mounted at the `secret/` prefix.

On this page, we'll learn about the mount system and the operations
that can be performed with it. We use this as prerequisite knowledge
for the next page, where we'll create dynamic secrets.

## Mount a Backend

To start, let's mount another _generic_ backend. Just like a normal
filesystem, Vault can mount a backend multiple times at different
mount points. This is useful if you want different access control policies
(covered later) or configurations for different paths.

To mount the backend:

```
$ vault mount generic
Successfully mounted 'generic' at 'generic'!
```

By default, the mount point will be the same name as the backend. This
is because 99% of the time, you don't want to customize this mount point.
In this example, we mounted the _generic_ backend at `generic/`.

You can inspect mounts using `vault mounts`:

```
$ vault mounts
Path      Type     Description
generic/  generic
secret/   generic  generic secret storage
sys/      system   system endpoints used for control, policy and debugging
```

You can see the `generic/` path we just mounted, as well as the built-in
secret path. You can also see the `sys/` path. We won't cover this in
this guide, but this mount point is used to interact with the Vault core
system.

Spend some time reading and writing secrets to the new mount point to
convince yourself it works. As a bonus, write to the `secret/` endpoint
and observe that those values are unavailable via the `generic/` path: they share the
same backend, but do not share any data. In addition to this, backends
(of the same type or otherwise) _cannot_ access the data of other backends;
they can only access data within their mount point.

## Unmount a Backend

Once you're sufficiently convinced mounts behave as you expect, you can
unmount it. When a backend is unmounted, all of its secrets are revoked
and its data is deleted. If either of these operations fail, the backend
remains mounted.

```
$ vault unmount generic/
Successfully unmounted 'generic/' if it was mounted
```

In addition to unmounting, you can remount a backend. Remounting a
backend changes its mount point. This is still a disruptive command: the
stored data is retained, but all secrets are revoked since secrets are
closely tied to their mount paths.

## What is a Secret Backend?

Now that you've mounted and unmounted a backend, you might wonder:
"what is a secret backend? what is the point of this mounting system?"

Vault behaves a lot like a [virtual filesystem](https://en.wikipedia.org/wiki/Virtual_file_system).
The read/write/delete operations are forwarded to the backend, and the
backend can choose to react to these operations however it wishes.
For example, the _generic_ backend simply passes this through to the
storage backend (after encrypting data first).

However, the _aws_ backend (which you'll see soon), will read/write IAM
policies and access tokens. So, while you might do a `vault read aws/deploy`,
this isn't reading from the physical path `aws/deploy`. Instead, the AWS
backend is dynamically generating an access key based on the `deploy` policy.

This abstraction is incredibly powerful. It lets Vault interface directly
with physical systems such as the backend as well as things such as SQL
databases, HSMs, etc. But in addition to these physical systems, Vault
can interact with more unique environments: AWS IAM, dynamic SQL user creation,
etc. all using the same read/write interface.

## Next

You now know about secret backends and how to operate on the mount table.
This is important knowledge to move forward and learn about other secret
backends.

Next, we'll use the
[AWS backend to generate dynamic secrets](/intro/getting-started/dynamic-secrets.html).
