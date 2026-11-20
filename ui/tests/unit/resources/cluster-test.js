/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import EmberObject, { computed } from '@ember/object';
import sinon from 'sinon';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import ClusterResource from 'vault/resources/cluster';
import ReplicationAttributesResource from 'vault/resources/replication-attributes';

module('Unit | Resource | cluster', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:version').type = 'enterprise';
    this.sealStatus = { initialized: true, sealed: false, t: 3, n: 5, progress: 0, type: 'shamir' };
    this.replication = { dr: { mode: 'disabled' }, performance: { mode: 'disabled' } };
    this.server.get('/sys/health', () => ({
      standby: false,
      license: { state: 'autoloaded', expiry_time: '2099-01-27T00:00:00Z' },
    }));
    this.server.get('/sys/seal-status', () => this.sealStatus);
    this.server.get('/sys/replication/status', () => ({ data: this.replication }));
    this.cluster = () =>
      new ClusterResource({ id: '1', name: 'vault' }, this.owner.lookup('service:version'));
  });

  test('it reads the cluster state', async function (assert) {
    const cluster = this.cluster();
    assert.true(cluster.needsInit, 'a cluster that has not been read needs init');
    assert.true(cluster.sealed, 'a cluster that has not been read is sealed');

    assert.strictEqual(await cluster.reload(), cluster, 'reload resolves to the cluster');
    assert.false(cluster.needsInit, 'the cluster is initialized');
    assert.true(cluster.unsealed, 'the cluster is unsealed');
    assert.false(cluster.sealed, 'the cluster is not sealed');
    assert.strictEqual(cluster.sealThreshold, 3, 'it reads the seal threshold');
    assert.strictEqual(cluster.sealType, 'shamir', 'it reads the seal type');
    assert.false(cluster.usingRaft, 'it is not using raft');
    assert.false(cluster.standby, 'it reads standby');
    assert.strictEqual(cluster.licenseState, 'autoloaded', 'it reads the license state');
    assert.strictEqual(cluster.licenseExpiry, '2099-01-27T00:00:00Z', 'it reads the license expiry');
    assert.true(cluster.allReplicationDisabled, 'replication is disabled');
    assert.notOk(cluster.anyReplicationEnabled, 'no replication is enabled');
  });

  test('it needs init until the seal status says it is initialized', async function (assert) {
    this.sealStatus = { initialized: false, sealed: true };
    const cluster = this.cluster();
    await cluster.reload();
    assert.true(cluster.needsInit, 'an uninitialized cluster needs init');
    assert.true(cluster.sealed, 'it is sealed');

    this.sealStatus = { initialized: true, sealed: true, storage_type: 'raft', progress: 1 };
    await cluster.reload();
    assert.false(cluster.needsInit, 'an initialized cluster does not');
    assert.true(cluster.sealed, 'it is still sealed');
    assert.true(cluster.usingRaft, 'it is using raft');
    assert.true(cluster.hasProgress, 'it has unseal progress');
  });

  test('it keeps previous values when a request fails', async function (assert) {
    this.replication = { dr: { mode: 'primary', cluster_id: 'dr-id' }, performance: { mode: 'disabled' } };
    const cluster = this.cluster();
    await cluster.reload();
    assert.false(cluster.sealed, 'the cluster is unsealed');
    assert.strictEqual(cluster.dr.mode, 'primary', 'dr is a primary');

    this.server.get('/sys/seal-status', () => new Response(500, {}, { errors: ['error'] }));
    this.server.get('/sys/replication/status', () => new Response(500, {}, { errors: ['error'] }));
    await cluster.reload();
    assert.false(cluster.sealed, 'the seal state is kept');
    assert.strictEqual(cluster.dr.mode, 'primary', 'the replication state is kept');
    assert.notOk(cluster.replicationRedacted, 'the replication is not redacted');
  });

  test('it flags a chroot namespace and redacted replication', async function (assert) {
    this.server.get('/sys/health', () => new Response(500, {}, { errors: ['error'] }));
    this.server.get('/sys/replication/status', () => new Response(403, {}, { errors: ['disabled path'] }));
    const cluster = this.cluster();
    await cluster.reload();

    assert.true(cluster.hasChrootNamespace, 'the cluster is in a chroot namespace');
    assert.true(cluster.replicationRedacted, 'replication is redacted');
    assert.notOk(cluster.dr, 'there is no dr state');
  });

  test('it only reads the replication status for enterprise in the root namespace', async function (assert) {
    let requested = 0;
    this.server.get('/sys/replication/status', () => {
      requested++;
      return { data: this.replication };
    });
    const cluster = this.cluster();
    this.owner.lookup('service:version').type = 'community';
    await cluster.reload();
    assert.strictEqual(requested, 0, 'it is not read for community');

    this.owner.lookup('service:version').type = 'enterprise';
    this.owner.lookup('service:namespace').setNamespace('ns1');
    await cluster.reload();
    assert.strictEqual(requested, 0, 'it is not read outside the root namespace');

    this.owner.lookup('service:namespace').setNamespace('');
    await cluster.reload();
    assert.strictEqual(requested, 1, 'it is read for enterprise in the root namespace');
  });

  test('an enabled replication mode is updated in place and a disabled one is replaced', async function (assert) {
    this.replication = { dr: { mode: 'primary', cluster_id: 'dr-id', last_wal: 1 } };
    const cluster = this.cluster();
    await cluster.reload();
    const { dr } = cluster;

    this.replication = { dr: { mode: 'primary', cluster_id: 'dr-id', last_wal: 2 } };
    await cluster.reload();
    assert.strictEqual(cluster.dr, dr, 'the enabled mode is the same');
    assert.strictEqual(cluster.dr.lastWAL, 2, 'its values are updated');

    this.replication = { dr: { mode: 'disabled' } };
    await cluster.reload();
    assert.notStrictEqual(cluster.dr, dr, 'a different mode is a different object');
    assert.strictEqual(cluster.dr.mode, 'disabled', 'it is disabled');
  });

  test('it ignores a response that is older than one that is already applied', async function (assert) {
    const cluster = this.cluster();
    const older = Promise.withResolvers();
    const newer = Promise.withResolvers();
    const readSealStatus = sinon.stub(cluster, 'readSealStatus');
    readSealStatus.onFirstCall().returns(older.promise);
    readSealStatus.onSecondCall().returns(newer.promise);

    const olderReload = cluster.reload();
    const newerReload = cluster.reload();
    newer.resolve({ initialized: true, sealed: false });
    await newerReload;
    assert.false(cluster.sealed, 'the newer response is applied');

    older.resolve({ initialized: true, sealed: true });
    await olderReload;
    assert.false(cluster.sealed, 'the older response does not replace it');
  });

  test('classic computed properties that depend on its getters are updated by a reload', async function (assert) {
    this.replication = { dr: { mode: 'primary', cluster_id: 'dr-id' } };
    const cluster = this.cluster();
    await cluster.reload();
    const subject = EmberObject.extend({
      state: computed('cluster.sealed', function () {
        return this.cluster.sealed ? 'sealed' : 'open';
      }),
      rootLabel: computed('cluster.dr.merkleRoot', function () {
        return `root ${this.cluster.dr.merkleRoot}`;
      }),
    }).create({ cluster });
    assert.strictEqual(subject.state, 'open', 'it is unsealed');
    assert.strictEqual(subject.rootLabel, 'root undefined', 'there is no merkle root');

    this.sealStatus = { initialized: true, sealed: true };
    this.replication = { dr: { mode: 'primary', cluster_id: 'dr-id', merkle_root: 'abc' } };
    await cluster.reload();
    assert.strictEqual(subject.state, 'sealed', 'a cluster getter is updated');
    assert.strictEqual(subject.rootLabel, 'root abc', 'a replication getter of the same object is updated');
  });

  test('it returns the replication attributes of the mode that is viewed', async function (assert) {
    this.replication = { dr: { mode: 'primary' }, performance: { mode: 'disabled' } };
    const cluster = this.cluster();
    await cluster.reload();

    assert.strictEqual(cluster.replicationAttrs, null, 'there is none when no mode is viewed');
    this.owner.lookup('service:replication-mode').setMode('dr');
    assert.strictEqual(cluster.replicationAttrs, cluster.dr, 'it returns dr');
    assert.strictEqual(cluster.replicationModeForDisplay, 'Disaster recovery', 'it displays dr');
    this.owner.lookup('service:replication-mode').setMode('performance');
    assert.strictEqual(cluster.replicationAttrs, cluster.performance, 'it returns performance');
    assert.strictEqual(cluster.replicationModeForDisplay, 'Performance', 'it displays performance');
    assert.true(cluster.anyReplicationEnabled, 'replication is enabled');
    assert.false(cluster.replicationIsInitializing, 'both modes have a mode');
  });
});

module('Unit | Resource | replication-attributes', function () {
  test('it derives the mode', function (assert) {
    const modes = {
      disabled: { replicationDisabled: true, replicationEnabled: false, modeForUrl: false },
      unsupported: { replicationDisabled: true, replicationEnabled: false, modeForUrl: false },
      primary: { replicationDisabled: false, replicationEnabled: true, modeForUrl: 'primary' },
      secondary: { replicationDisabled: false, replicationEnabled: true, modeForUrl: 'secondary' },
      bootstrapping: { replicationDisabled: false, replicationEnabled: true, modeForUrl: 'bootstrapping' },
    };
    for (const [mode, expected] of Object.entries(modes)) {
      const attrs = new ReplicationAttributesResource({ mode });
      for (const [key, value] of Object.entries(expected)) {
        assert.strictEqual(attrs[key], value, `${mode}: ${key}`);
      }
    }
    assert.true(new ReplicationAttributesResource({ mode: 'primary' }).isPrimary, 'it is a primary');
    assert.true(new ReplicationAttributesResource({ mode: 'secondary' }).isSecondary, 'it is a secondary');
    assert.strictEqual(
      new ReplicationAttributesResource({ mode: 'bootstrapping' }).modeDisplay,
      'Bootstrapping'
    );
    assert.strictEqual(
      new ReplicationAttributesResource({}).modeDisplay,
      'Disabled',
      'it defaults to disabled'
    );
    assert.strictEqual(new ReplicationAttributesResource({}).modeForHeader, 'loading', 'it is loading');
  });

  test('it derives the cluster id and sync progress', function (assert) {
    const attrs = new ReplicationAttributesResource({
      cluster_id: 'abc-123-def',
      mode: 'secondary',
      state: 'merkle-sync',
      merkle_sync_progress: { sync_total_keys: 8, sync_progress: 3 },
    });
    assert.strictEqual(attrs.id, 'abc-123-def', 'the id is the cluster id');
    assert.strictEqual(attrs.clusterIdDisplay, 'abc', 'it shortens the cluster id');
    assert.deepEqual(attrs.syncProgress, { progress: 3, total: 8 }, 'it reads the sync progress');
    assert.strictEqual(attrs.syncProgressPercent, 37, 'it calculates the percentage');

    attrs.update({ state: 'idle' });
    assert.strictEqual(attrs.syncProgress, null, 'there is no progress when not syncing');
    assert.strictEqual(attrs.mode, 'secondary', 'update keeps the other values');
  });

  test('it coerces values like the Ember Data transforms did', function (assert) {
    const attrs = new ReplicationAttributesResource({
      known_secondaries: null,
      last_wal: '455',
      secondaries: undefined,
    });
    assert.deepEqual(attrs.knownSecondaries, [], 'a null array is empty');
    assert.strictEqual(attrs.lastWAL, 455, 'a number is parsed');
    assert.strictEqual(attrs.secondaries, undefined, 'a missing value stays undefined');
    assert.ok(new ReplicationAttributesResource({}).id, 'an id is generated without a cluster id');
  });
});
