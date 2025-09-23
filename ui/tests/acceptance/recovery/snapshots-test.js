/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, currentRouteName, currentURL, click } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { Response } from 'miragejs';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { addDays } from 'date-fns';

module('Acceptance | recovery | snapshots', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('enterprise: it renders empty state when raft storage is not in use', async function (assert) {
    this.server.get('/sys/storage/raft/snapshot-load', () => {
      return overrideResponse(400, JSON.stringify({ errors: ['raft storage is not in use'] }));
    });
    await visit('/vault/recovery/snapshots');
    assert.strictEqual(currentURL(), '/vault/recovery/snapshots');
    assert.dom('header').exists('it renders header despite route throwing an error');
    assert.dom(GENERAL.emptyStateTitle).hasText('Raft storage required');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Raft storage must be used in order to recover data from a snapshot.');
    assert.dom(GENERAL.emptyStateActions).hasText('Snapshot management');
  });

  test('it renders promo for community versions', async function (assert) {
    const version = this.owner.lookup('service:version');
    version.type = 'community';
    this.server.get('/sys/storage/raft/snapshot-load', () => {
      // This assertion is intentionally setup to fail if a request is made to this endpoint
      // because community versions should NOT request the snapshot-load endpoint
      assert.true(false, 'it does not make a request to snapshot-load on CE versions');
    });

    await visit('/vault/recovery/snapshots');
    assert
      .dom(`${GENERAL.navLink('Secrets Recovery')} .hds-badge`)
      .hasText('Enterprise', 'side nav link renders "Enterprise" badge');
    assert.strictEqual(currentURL(), '/vault/recovery/snapshots');
    assert.dom(GENERAL.emptyStateTitle).hasText('Secrets Recovery is an enterprise feature');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Secrets Recovery allows you to restore accidentally deleted or lost secrets from a snapshot. The snapshots can be provided via upload or loaded from external storage.'
      );
    assert.dom(GENERAL.emptyStateActions).hasText('Learn more about upgrading');
    assert.dom(GENERAL.badge('enterprise')).exists();
  });

  module('enterprise: with raft configured', function (hooks) {
    hooks.beforeEach(function () {
      this.server.get('/sys/storage/raft/configuration', () =>
        this.server.create('configuration', 'withRaft')
      );
    });

    test('it renders empty state when no snapshots are loaded', async function (assert) {
      this.server.get('/sys/storage/raft/snapshot-load', () => {
        return new Response(404, { 'Content-Type': 'application/json' }, JSON.stringify({ errors: [] }));
      });

      await visit('/vault/recovery/snapshots');

      assert.strictEqual(currentURL(), '/vault/recovery/snapshots');

      assert.dom(GENERAL.emptyStateTitle).hasText('Upload a snapshot to get started');
      assert.dom(GENERAL.emptyStateActions).hasText('Upload snapshot');
    });

    test('it redirects to snapshot route when a snapshot is loaded', async function (assert) {
      this.server.get('/sys/storage/raft/snapshot-load', () => {
        return { data: { keys: ['1234'] } };
      });

      this.server.get('/sys/storage/raft/snapshot-load/1234', () => {
        return {
          data: {
            status: 'ready',
            expires_at: addDays(new Date(), 3).toISOString(),
            snapshot_id: '1234',
          },
        };
      });

      await click(GENERAL.navLink('Secrets Recovery'));
      assert.strictEqual(currentURL(), '/vault/recovery/snapshots/1234/manage');
      assert.strictEqual(currentRouteName(), 'vault.cluster.recovery.snapshots.snapshot.manage');
    });
  });
});
