---
layout: "guides"
page_title: "Consul Storage Migration"
sidebar_current: "guides-consul-storage-migration"
description: |-
  Migrate Vault data stored in Consul.
---

# Consul Storage Migration

[Consul](https://www.consul.io/) is the HashiCorp recommended storage backend for Vault; it provides a range of advantages over other storage backends.

> If you are not currently using Consul for your Vault storage backend, be sure to check out [Consul Enterprise](https://www.consul.io/docs/enterprise/index.html).

When using Consul as a Vault storage backend a number of of options for migration or backup/restore of Vault data are available to you. Here are some basic Vault data migration techniques which you can use to inform your own Vault data retention and migration strategies.

## Consul KV Export

For Consul server clusters which are used for more use cases than only storing Vault, you generally want to be more be precise about what is restored to avoid disturbing the other Consul data.

You can migrate only the Vault data and skip all other Consul data that would be included with Consul snapshot by using the [consul kv](https://www.consul.io/docs/commands/kv.html) command's [export](https://www.consul.io/docs/commands/kv/export.html) and [import](https://www.consul.io/docs/commands/kv/import.html) sub-commands. (requires Consul version 0.7.3+)

Here's a brief example of the process.

### Source Consul Server Cluster

Export the Vault key/value path (named `vault` by default) from the source Consul server cluster; this command can be run on any of the Consul servers or Consul client agents associated with the servers that store the Vault data:

```
$ consul kv export vault/ > vault.json
```

**NOTE**: Check the value of `path` in your Vault configuration file's `storage` stanza in the case that your Vault uses a different path in Consul's K/V, and use that path instead.

### Destination Consul Server Cluster

From a Consul client agent or server associated with the **destination** Vault instance that you wish to restore to, import the JSON data into Consul.

The correct path will be automatically created in the key/value store:

```
$ consul kv import @vault.json
Imported: vault/audit/0e9483e5-609d-efcd-3e14-c337e48e1f14/salt
...
Imported: vault/wal/logs/00000002/7912
```

-> **IMPORTANT**: Please see the **Post Restoration** section for details on what to do after restoring Vault data.

## Consul Snapshots

If your Consul cluster is used exclusively for Vault data, then you can simply save and restore Consul snapshots as a backup/restoration solution.

This solution also has a nice [Automated Agent in Consul Enterprise](https://www.consul.io/docs/commands/snapshot/agent.html) that helps to ensure snapshots are taken on your desired schedule, and retained in your specified destinations.

See `consul snapshot --help` or the [Consul Snapshot documentation](https://www.consul.io/docs/commands/snapshot.html) for more information about the command

Here is a brief example of manually saving and restoring a Consul snapshot.

### Source Consul Server Cluster

On the **source** Consul server cluster that contains the Vault data to be saved in a snapshot:

```
 $ consul snapshot save backup.snap
 Saved and verified snapshot to index 362428
```

### Destinatation Consul Server Cluster

On the **destination** Consul server cluster that will contain the Vault data to be restored from snapshot:

```
 $ consul snapshot restore backup.snap
Restored snapshot
```

-> **IMPORTANT**: Please see the **Post Restoration** section for details on what to do after restoring Vault data.

## Disaster Recovery Mode Replication

[Vault Enterprise][vault-enterprise] offers a third option which can be realized with the Consul storage backend, and that is disaster recovery mode (DR) replication.

This process essentially boils down to:

1. Enable replication on the source Vault cluster as a *Disaster Recovery mode* **Primary**
2. Configure a DR mode secondary cluster
3. Replication of Vault data will occur between the primary and secondary

Before using this method, you should be familiar with our replication documentation and in particular, the following resources:

- [Vault Replication][vault-replication]
- [Replication Setup & Guidance][vault-replication-guide]
- [Vault Replication API][vault-replication-api]

-> **IMPORTANT**: Please see the **Post Restoration** section for details on what to do after restoring Vault data.

## Post Restoration

This section contains important notes which you should familiarize yourself with prior to performing any migration or backup/restore of Vault data in Consul.

Please be aware of the following caveats and conditions around restoring Vault data before you proceed with a backup and restoration.

### High Availability Mode Lock

In an High Availability Vault cluster, the active node will have held the cluster leadership lock at the time of the data export or snapshot. After restoring Vault data to Consul, you must manually remove this lock so that the Vault cluster can elect a new leader.

Execute this `consul kv` command immediately after restoration of Vault data to Consul:

```
$ consul kv delete vault/core/lock
```

See `consul kv delete --help` or the [Consul KV Delete documentation](https://www.consul.io/docs/commands/kv/delete.html) for more details on the command.

### Dynamic Secret Backends

With use of dynamic secret backends there could be user credentials in for example a database secret backend like PostgreSQL that Vault doesn't have knowledge of.

For example, if users were created after you took the snapshot, then Vault would not be aware of them after restoring the snapshot. There is not much that can be done to mitigate this; frequent snapshotting of Vault data can help.

If you're using dynamic secret backends, after restoring Vault data, you should go through active Vault users, and revoke them all to force your clients to get new credentials and generate a new lease.

### Deleted Users After Snapshot

If users for a given backend were deleted after you took the snapshot you are restoring from, you could experience issues with Vault automatically revoking their leases, which appear in the logs as revocation errors along with `User not found` or `no such user`.

In these cases, you'll need to manually force revocation of the user by their lease ID. Here's an example:

```
$ vault revoke -force -prefix ce9e899b-49d0-9646-9769-381909fea995
Success! Revoked the secret with ID 'ce9e899b-49d0-9646-9769-381909fea995', if it existed.
```

If you want to use the `vault` command to revoke, see `vault revoke --help` for more details on the `-force` flag syntax.

To learn more about doing this programmatically, see the [Revoke Force API documentation][vault-revoke-force-api].

### Restoration After Vault is Rekeyed

Some tips which can help with the scenario where restoration of an older Vault export or snapshot occurs after Vault is rekeyed:

- Use key manager to store unseal keys so you have a versioned history of them
- When transmitted PGP encrypted keys, just use email so you have a history of the unseal keys there
- Archive PGP encrypted unseal keys into a backup and store it somewhere in the event you have to do an older restore
- You can even maintain a history of PGP-keys stored in Vault

## Resources

1. [Consul KV Export](https://www.consul.io/docs/commands/kv/export.html)
2. [Consul KV Import](https://www.consul.io/docs/commands/kv/import.html)
3. [Consul Enterprise Snapshot Agent](https://www.consul.io/docs/commands/snapshot/agent.html)
4. [consul snapshot command](https://www.consul.io/docs/commands/snapshot.html)
5. [Vault Enterprise][vault-enterprise]
5. [Vault Replication][vault-replication]
6. [Replication Setup & Guidance][vault-replication-guide]
7. [Vault Replication API][vault-replication-api]
8. [Revoke Force API documentation][vault-revoke-force-api]


[vault-enterprise]: /docs/enterprise/index.html
[vault-replication]: /docs/enterprise/replication/index.html
[vault-replication-guide]: /guides/replication.html
[vault-replication-api]: /api/system/replication.html
[vault-revoke-force-api]: /api/system/leases.html#revoke-force