/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, visit, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { disableReplication, enableReplication } from 'vault/tests/helpers/replication';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';

// To allow a user to login and create a secondary dr cluster we demote a primary dr cluster.
// We stub this demotion so we do not break the dev process for all future tests.
// All DR secondary assertions are done in one test to avoid the lengthy setup and teardown process.
module('Acceptance | Enterprise | replication-secondaries', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    await settled();
    await disableReplication('dr');
    await settled();
    await disableReplication('performance');
    await settled();
  });

  hooks.afterEach(async function () {
    // For the tests following this, return to a good state.
    // We've reset mirage with this.server.shutdown() but re-poll the cluster to get the latest state.
    this.server.shutdown();
    await pollCluster(this.owner);
    await disableReplication('dr');
    await settled();
    await pollCluster(this.owner);
    await logout();
  });

  test('DR secondary: manage tab, details tab, and analytics are not run', async function (assert) {
    // Log in and set up a DR primary
    await login();
    await settled();
    await enableReplication('dr', 'primary');
    await pollCluster(this.owner);
    await click('[data-test-replication-link="manage"]');

    // Stub the demote action so it does not actually demote the cluster
    this.server.post('/sys/replication/dr/demote', () => {
      return { request_id: 'fake-demote', data: { success: true } };
    });
    // Stub endpoints for DR secondary state
    this.server.post('/sys/capabilities-self', () => ({ capabilities: [] }));
    this.server.get('/sys/replication/status', () => ({
      request_id: '2f50313f-be70-493d-5883-c84c2d6f05ce',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        dr: {
          cluster_id: '7222cbbf-3fb3-949b-8e03-cd5a15babde6',
          corrupted_merkle_tree: false,
          known_primary_cluster_addrs: null,
          last_corruption_check_epoch: '-62135596800',
          last_reindex_epoch: '0',
          merkle_root: 'd3ae75bde029e05d435f92b9ecc5641c1b027cc4',
          mode: 'secondary',
          primaries: [],
          primary_cluster_addr: '',
          secondary_id: '',
          ssct_generation_counter: 0,
          state: 'idle',
        },
        performance: {
          mode: 'disabled',
        },
      },
      wrap_info: null,
      warnings: null,
      auth: null,
      mount_type: '',
    }));
    this.server.get('/sys/replication/dr/status', () => ({
      data: { mode: 'secondary', cluster_id: 'dr-cluster-id' },
    }));
    this.server.get('/sys/health', () => ({
      initialized: true,
      sealed: false,
      standby: false,
      performance_standby: false,
      replication_performance_mode: 'disabled',
      replication_dr_mode: 'secondary',
      server_time_utc: 1754948244,
      version: '1.21.0-beta1+ent',
      enterprise: true,
      cluster_name: 'vault-cluster-64853bcd',
      cluster_id: '113a6c47-077f-bea7-0e8e-70a91821e85a',
      last_wal: 82,
      license: {
        state: 'autoloaded',
        expiry_time: '2029-01-27T00:00:00Z',
        terminated: false,
      },
      echo_duration_ms: 0,
      clock_skew_ms: 0,
      replication_primary_canary_age_ms: 0,
      removed_from_cluster: false,
    }));
    this.server.get('/sys/seal-status', () => ({
      type: 'shamir',
      initialized: true,
      sealed: false,
      t: 1,
      n: 1,
      progress: 0,
      nonce: '',
      version: '1.21.0-beta1+ent',
      build_date: '2025-08-11T14:11:00Z',
      migration: false,
      cluster_name: 'vault-cluster-64853bcd',
      cluster_id: '113a6c47-077f-bea7-0e8e-70a91821e85a',
      recovery_seal: false,
      storage_type: 'raft',
      removed_from_cluster: false,
    }));

    await click(GENERAL.button('demote'));
    await pollCluster(this.owner); // We must poll the cluster to stimulate a cluster reload. This is skipped in ember testing so must be forced.

    // Spy on the route's addAnalyticsService method
    const clusterRoute = this.owner.lookup('route:vault.cluster');
    const addAnalyticsSpy = sinon.spy(clusterRoute, 'addAnalyticsService');

    // Visit the DR secondary view. This is the route used only by DR secondaries
    await visit('/vault/replication-dr-promote');

    assert
      .dom('[data-test-promote-description]')
      .hasText(
        'Promote this cluster to a Disaster Recovery primary',
        'shows the correct description for a DR secondary'
      );
    assert.dom(GENERAL.badge('secondary')).includesText('secondary', 'shows the DR secondary mode badge');

    await click(GENERAL.linkTo('Details'));
    assert
      .dom('[data-test-replication-secondary-card]')
      .hasClass(
        'has-error-border',
        'shows error border on status because the DR secondary is not connected to a primary.'
      );

    // Assert addAnalyticsService was NOT called
    assert.false(addAnalyticsSpy.called, 'addAnalyticsService should not be called on DR secondary');

    // Restore spy
    addAnalyticsSpy.restore();
  });
});
