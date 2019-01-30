---
layout: "docs"
page_title: "Vault Agent Auto-Auth File Sink"
sidebar_title: "File"
sidebar_current: "docs-agent-autoauth-sinks-file"
description: |-
  File sink for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth File Sink 

The `file` sink writes tokens, optionally response-wrapped and/or encrypted, to
a file. This may be a local file or a file mapped via some other process (NFS,
Gluster, CIFS, etc.).

Once the sink writes the file, it is up to the client to control lifecycle;
generally it is best for the client to remove the file as soon as it is seen.

It is also best practice to write the file to a ramdisk, ideally an encrypted
ramdisk, and use appropriate filesystem permissions. The file is currently
always written with `0640` permissions.

## Configuration

- `path` `(string: required)` - The path to use to write the token file
