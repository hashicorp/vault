/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, currentRouteName, currentURL, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | recovery | snapshots-load', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.get('/sys/storage/raft/configuration', () => this.server.create('configuration', 'withRaft'));

    return login();
  });

  test('it redirects to snapshot route when a snapshot is loaded', async function (assert) {
    this.server.get('/sys/storage/raft/snapshot-load', () => {
      return { data: { keys: [] } };
    });

    await visit('vault/recovery/snapshots');

    await click(`${GENERAL.emptyStateActions} a`);

    assert.strictEqual(currentURL(), '/vault/recovery/snapshots/load');
    assert.strictEqual(currentRouteName(), 'vault.cluster.recovery.snapshots.load');
  });
});
