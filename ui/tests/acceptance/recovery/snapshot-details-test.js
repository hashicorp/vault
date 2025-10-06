/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, currentRouteName, currentURL, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import recoveryHandler from 'vault/mirage/handlers/recovery';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  action: (index) => `.hds-dropdown__list li:nth-of-type(${index}) > :first-child`,
};

module('Acceptance | recovery | snapshot-details', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    recoveryHandler(this.server);

    this.snapshot = this.server.create('snapshot');

    return login();
  });

  test('it redirects to snapshots route when a snapshot is unloaded', async function (assert) {
    await visit('vault/recovery/snapshots');
    await click('[data-test-details-link]');

    await click(GENERAL.button('toggle'));

    await click(SELECTORS.action(2));

    assert.strictEqual(currentURL(), '/vault/recovery/snapshots');
    assert.strictEqual(currentRouteName(), 'vault.cluster.recovery.snapshots.index');
  });

  test('it makes the snapshot status requests in the root namespace', async function (assert) {
    // 2 for initial req + 2 for polling req
    assert.expect(4);

    this.server.get('/sys/storage/raft/snapshot-load/:snapshot_id', (schema, req) => {
      assert.propEqual(req.params, {
        snapshot_id: this.snapshot.snapshot_id,
      });
      assert.strictEqual(req.requestHeaders['x-vault-namespace'], '');

      const record = schema.db['snapshots'].findBy(req.params);
      if (record) {
        delete record.id; // "snapshot_id" is the id
        return { data: record };
      }
      return new Response(404, {}, { errors: [] });
    });

    // go directly to details page to avoid polling requests from manage page
    await visit(`/vault/recovery/snapshots/${this.snapshot.snapshot_id}/details`);
  });

  test('it redirects to snapshot route when manage is selected', async function (assert) {
    await visit('vault/recovery/snapshots');
    await click('[data-test-details-link]');

    await click(GENERAL.button('toggle'));
    await click(SELECTORS.action(1));

    assert.strictEqual(currentURL(), `/vault/recovery/snapshots/${this.snapshot.snapshot_id}/manage`);
    assert.strictEqual(currentRouteName(), 'vault.cluster.recovery.snapshots.snapshot.manage');
  });
});
