---
layout: "docs"
page_title: "Upgrading Vault - Guides"
sidebar_title: "Upgrade Guides"
sidebar_current: "docs-upgrading"
description: |-
  These are general upgrade instructions for Vault for both non-HA and HA
  setups. Please ensure that you also read the version-specific upgrade notes.
---

# Upgrading Vault

These are general upgrade instructions for Vault for both non-HA and HA setups.
_Please ensure that you also read any version-specific upgrade notes which can be
found in the sidebar._

!> **IMPORTANT NOTE:** Always back up your data before upgrading! Vault does not
make backward-compatibility guarantees for its data store. Simply replacing the
newly-installed Vault binary with the previous version will not cleanly
downgrade Vault, as upgrades may perform changes to the underlying data
structure that make the data incompatible with a downgrade. If you need to roll
back to a previous version of Vault, you should roll back your data store as
well.

## Testing the Upgrade

It's always a good idea to try to ensure that the upgrade will be successful in
your environment. The ideal way to do this is to take a snapshot of your data
and load it into a test cluster. However, if you are issuing secrets to third
party resources (cloud credentials, database credentials, etc.) ensure that you
do not allow external network connectivity during testing, in case credentials
expire. This prevents the test cluster from trying to revoke these resources
along with the non-test cluster.

## Non-HA Installations

Upgrading non-HA installations of Vault is as simple as replacing the Vault
binary with the new version and restarting Vault. Any upgrade tasks that can be
performed for you will be taken care of when Vault is unsealed.

Always use `SIGINT` or `SIGTERM` to properly shut down Vault.

Be sure to also read and follow any instructions in the version-specific
upgrade notes.

## HA Installations

This is our recommended upgrade procedure, and the procedure we use internally
at HashiCorp. However, you should consider how to apply these steps to your
particular setup since HA setups can differ on whether a load balancer is in
use, what addresses clients are being given to connect to Vault (standby +
leader, leader-only, or discovered via service discovery), etc.

Whatever method you use, you should ensure that you never fail over from a
newer version of Vault to an older version. Our suggested procedure is designed
to prevent this.

Please note that Vault does not support true zero-downtime upgrades, but with
proper upgrade procedure the downtime should be very short (a few hundred
milliseconds to a second depending on how the speed of access to the storage
backend).

Perform these steps on each standby:

1. Properly shut down Vault on the standby node via `SIGINT` or `SIGTERM`
2. Replace the Vault binary with the new version
3. Start the standby node
4. Unseal the standby node

At this point all standby nodes will be upgraded and ready to take over. The
upgrade will not be complete until one of the upgraded standby nodes takes over
active duty. To do this:

1. Properly shut down the remaining (active) node. Note: it is _**very
   important**_ that you shut the node down properly. This causes the HA lock to
   be released, allowing a standby node to take over with a very short delay.
   If you kill Vault without letting it release the lock, a standby node will
   not be able to take over until the lock's timeout period has expired. This
   is backend-specific but could be ten seconds or more.
2. Replace the Vault binary with the new version; ensure that `mlock()` capability is added to the new binary with [setcap](https://www.vaultproject.io/docs/configuration/index.html#disable_mlock)
3. Start the node
4. Unseal the node (it will now be a standby)

Internal upgrade tasks will happen after one of the upgraded standby nodes
takes over active duty.

Be sure to also read and follow any instructions in the version-specific
upgrade notes.

## Replication Installations

-> **Note:** Prior to any upgrade, be sure to also read and follow any instructions in the version-specific
upgrade notes which are found in the navigation menu for this documentation.

Upgrading installations of Vault which participate in [Enterprise Replication](/docs/enterprise/replication/index.html) requires the following basic order of operations:

- **Upgrade the replication secondary instances first** using appropriate
  guidance from the previous sections depending on whether each secondary
  instance is Non-HA or HA
- Verify functionality of each secondary instance after upgrading
- When satisfied with functionality of upgraded secondary instances, upgrade
  the primary instance
