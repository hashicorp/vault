/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, currentRouteName, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { Response } from 'miragejs';

module('Acceptance | recovery snapshots', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.get('/sys/storage/raft/configuration', () => this.server.create('configuration', 'withRaft'));

    this.store = this.owner.lookup('service:store');
    return login();
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
    this.server.get('/sys/storage/raft/snapshot-load', () =>
      overrideResponse(null, {
        data: {
          keys: ['1234'],
        },
      })
    );

    this.server.get('/sys/storage/raft/snapshot-load/1234', () =>
      overrideResponse(null, {
        data: {
          status: 'ready',
          expires_at: new Date(),
          snapshot_id: '1234',
        },
      })
    );

    await visit('vault/recovery/snapshots');

    assert.strictEqual(currentURL(), '/vault/recovery/snapshots/1234/manage');
    assert.strictEqual(currentRouteName(), 'vault.cluster.recovery.snapshots.snapshot.manage');
  });
});
