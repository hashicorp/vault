/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { currentRouteName, visit } from '@ember/test-helpers';
import { Response } from 'miragejs';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { mockReplication } from 'vault/tests/helpers/replication-mocks';
import { mockReplicationBlock } from 'vault/tests/helpers/replication';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';

// These tests describe how the cluster state is read and refreshed. They only use requests, the
// current cluster and the pages, so they hold for any implementation of the cluster record.

const HEALTH_QUERY = ['standbycode', 'sealedcode', 'uninitcode', 'drsecondarycode', 'performancestandbycode'];
const RAFT_LINK = '[data-test-sidebar-nav-link="Raft storage"]';

module('Acceptance | cluster state', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.sealStatus = {
      type: 'shamir',
      initialized: true,
      sealed: false,
      t: 3,
      n: 5,
      progress: 0,
      version: '1.21.0+ent',
      storage_type: 'file',
    };
    this.recorded = [];
    // the cluster state endpoints are recorded so the tests can check how they are requested
    const record = (path, request) =>
      this.recorded.push({
        path,
        query: request.queryParams,
        headers: Object.fromEntries(
          Object.entries(request.requestHeaders || {}).map(([key, value]) => [key.toLowerCase(), value])
        ),
      });
    this.setup = async (options = {}) => {
      this.state = mockReplication(this.server, options);
      this.server.get('/sys/health', (schema, request) => {
        record('/sys/health', request);
        return {
          initialized: true,
          sealed: false,
          standby: false,
          version: '1.21.0+ent',
          enterprise: true,
          cluster_name: 'vault-cluster',
          cluster_id: '113a6c47-077f-bea7-0e8e-70a91821e85a',
          license: { state: 'autoloaded', expiry_time: '2099-01-27T00:00:00Z', terminated: false },
        };
      });
      this.server.get('/sys/seal-status', (schema, request) => {
        record('/sys/seal-status', request);
        return this.sealStatus;
      });
      this.server.get('/sys/replication/status', (schema, request) => {
        record('/sys/replication/status', request);
        if (options.redacted) return new Response(403, {}, { errors: ['disabled path'] });
        return { data: { dr: this.state.dr, performance: this.state.performance } };
      });
      if (options.login !== false) await login();
    };
    this.cluster = () => this.owner.lookup('service:current-cluster').cluster;
  });

  test('it reads the cluster state without a token from the root namespace', async function (assert) {
    await this.setup({ dr: 'primary' });
    const start = this.recorded.length;
    await pollCluster(this.owner);
    const requests = this.recorded.slice(start);

    for (const path of ['/sys/health', '/sys/seal-status', '/sys/replication/status']) {
      const matches = requests.filter((request) => request.path === path);
      assert.strictEqual(matches.length, 1, `${path} is requested once`);
      assert.notOk(matches[0].headers['x-vault-token'], `${path} is sent without a token`);
      assert.notOk(matches[0].headers['x-vault-namespace'], `${path} is sent without a namespace`);
    }
    const health = requests.find((request) => request.path === '/sys/health');
    for (const param of HEALTH_QUERY) {
      assert.strictEqual(String(health.query[param]), '200', `health asks for ${param}=200`);
    }
  });

  test('it does not read the replication status outside the root namespace', async function (assert) {
    await this.setup({ dr: 'primary' });
    this.owner.lookup('service:namespace').setNamespace('ns1');
    const start = this.recorded.length;
    await pollCluster(this.owner);
    const requests = this.recorded.slice(start);

    assert.strictEqual(
      requests.filter((request) => request.path === '/sys/replication/status').length,
      0,
      'replication status is not requested'
    );
    for (const path of ['/sys/health', '/sys/seal-status']) {
      const matches = requests.filter((request) => request.path === path);
      assert.strictEqual(matches.length, 1, `${path} is still requested`);
      assert.notOk(matches[0].headers['x-vault-namespace'], `${path} is sent to the root namespace`);
    }
  });

  test('a failed health request marks the cluster as a chroot namespace', async function (assert) {
    await this.setup({ dr: 'primary', login: false });
    this.sealStatus = { ...this.sealStatus, storage_type: 'raft' };
    this.server.get('/sys/health', () => new Response(500, {}, { errors: ['namespace not available'] }));
    await login();
    await visit('/vault/dashboard');

    assert.true(this.cluster().hasChrootNamespace, 'the cluster is marked as a chroot namespace');
    assert.true(this.cluster().usingRaft, 'the cluster uses raft');
    assert.dom(RAFT_LINK).doesNotExist('root namespace links are hidden');
  });

  test('redacted replication endpoints hide replication', async function (assert) {
    await this.setup({ redacted: true });
    await visit('/vault/dashboard');

    assert.true(this.cluster().replicationRedacted, 'the cluster is marked as redacted');

    await visit('/vault/replication');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.dashboard',
      'the replication pages redirect to the dashboard'
    );
  });

  test('a replication status error that is not a redaction does not mark the cluster as redacted', async function (assert) {
    await this.setup({ login: false });
    this.server.get('/sys/replication/status', () => new Response(500, {}, { errors: ['internal error'] }));
    await login();

    assert.notOk(this.cluster().replicationRedacted, 'the cluster is not marked as redacted');
  });

  test('a cluster whose seal status cannot be read is treated as uninitialized', async function (assert) {
    await this.setup({ login: false });
    this.server.get('/sys/seal-status', () => new Response(500, {}, { errors: ['internal error'] }));
    await visit('/vault');

    assert.strictEqual(currentRouteName(), 'vault.cluster.init', 'it redirects to init');
  });

  test('a seal status that cannot be read on a refresh keeps the previous seal state', async function (assert) {
    await this.setup();
    assert.false(this.cluster().needsInit, 'the cluster is initialized');
    assert.false(this.cluster().sealed, 'the cluster is unsealed');

    this.server.get('/sys/seal-status', () => new Response(500, {}, { errors: ['internal error'] }));
    await pollCluster(this.owner);

    assert.false(this.cluster().needsInit, 'the cluster is still initialized');
    assert.false(this.cluster().sealed, 'the cluster is still unsealed');
  });

  test('a refresh updates the current cluster in place', async function (assert) {
    await this.setup();
    const cluster = this.cluster();
    assert.false(cluster.sealed, 'the cluster is unsealed');
    assert.strictEqual(cluster.sealThreshold, 3, 'it reads the seal threshold');

    this.sealStatus = { ...this.sealStatus, sealed: true, t: 2, progress: 1, type: 'awskms' };
    await pollCluster(this.owner);

    assert.strictEqual(this.cluster(), cluster, 'it is the same cluster');
    assert.true(cluster.sealed, 'it is sealed');
    assert.strictEqual(cluster.sealThreshold, 2, 'the seal threshold is updated');
    assert.strictEqual(cluster.sealProgress, 1, 'the seal progress is updated');
    assert.strictEqual(cluster.sealType, 'awskms', 'the seal type is updated');
  });

  test('a refresh updates the replication state', async function (assert) {
    await this.setup({ dr: 'primary' });
    const cluster = this.cluster();
    const dr = cluster.dr;
    assert.strictEqual(dr.mode, 'primary', 'dr is a primary');
    assert.true(cluster.anyReplicationEnabled, 'replication is enabled');
    assert.false(cluster.allReplicationDisabled, 'replication is not disabled');
    assert.strictEqual(cluster.performance.mode, 'disabled', 'performance is disabled');

    this.state.dr = { ...mockReplicationBlock('primary'), last_wal: 999 };
    await pollCluster(this.owner);
    assert.strictEqual(cluster.dr, dr, 'an enabled mode is updated in place');
    assert.strictEqual(cluster.dr.lastWAL, 999, 'its values are updated');

    this.state.dr = mockReplicationBlock('disabled');
    await pollCluster(this.owner);
    assert.strictEqual(cluster.dr.mode, 'disabled', 'dr is disabled');
    assert.true(cluster.allReplicationDisabled, 'all replication is disabled');
    assert.notStrictEqual(cluster.dr, dr, 'a disabled mode is a different record');
  });

  test('the raft storage link follows the storage type', async function (assert) {
    await this.setup();
    await visit('/vault/dashboard');
    assert.dom(RAFT_LINK).doesNotExist('there is no raft link for other storage');

    this.sealStatus = { ...this.sealStatus, storage_type: 'raft' };
    await pollCluster(this.owner);
    assert.true(this.cluster().usingRaft, 'the cluster uses raft');
    assert.dom(RAFT_LINK).exists('it shows the raft link');
  });
});
