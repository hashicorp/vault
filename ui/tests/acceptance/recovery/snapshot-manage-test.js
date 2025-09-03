/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, fillIn, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import recoveryHandler from 'vault/mirage/handlers/recovery';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | recovery | snapshot-manage', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    recoveryHandler(this.server);

    this.snapshot = this.server.create('snapshot');

    return login();
  });

  test('it makes the recovery requests in the root namespace', async function (assert) {
    assert.expect(2);

    this.server.post('/cubbyhole/:path', (schema, req) => {
      assert.propEqual(req.queryParams, {
        recover_snapshot_id: this.snapshot.snapshot_id,
      });
      assert.strictEqual(req.requestHeaders['x-vault-namespace'], '');
    });

    await visit('/vault/recovery/snapshots');

    await click(GENERAL.selectByAttr('mount'));
    await click('[data-option-index]');
    await fillIn(GENERAL.inputByAttr('resourcePath'), 'recovered-secret');

    await click(GENERAL.button('recover'));
  });

  test('it makes the recovery requests in the child namespace', async function (assert) {
    assert.expect(2);

    this.server.post('/cubbyhole/:path', (schema, req) => {
      assert.propEqual(req.queryParams, {
        recover_snapshot_id: this.snapshot.snapshot_id,
      });
      assert.strictEqual(req.requestHeaders['x-vault-namespace'], 'child-ns-1');
    });

    await visit('/vault/recovery/snapshots');

    await click(GENERAL.selectByAttr('namespace'));
    await click('[data-option-index="1"]');
    await click(GENERAL.selectByAttr('mount'));
    await click('[data-option-index]');
    await fillIn(GENERAL.inputByAttr('resourcePath'), 'recovered-secret');

    await click(GENERAL.button('recover'));
  });

  test('it makes the recovery requests in the grandchild namespace', async function (assert) {
    assert.expect(2);

    this.server.post('/cubbyhole/:path', (schema, req) => {
      assert.propEqual(req.queryParams, {
        recover_snapshot_id: this.snapshot.snapshot_id,
      });
      assert.strictEqual(req.requestHeaders['x-vault-namespace'], 'child-ns-1/nested');
    });

    await visit('/vault/recovery/snapshots');

    await click(GENERAL.selectByAttr('namespace'));
    await click('[data-option-index="2"]');
    await click(GENERAL.selectByAttr('mount'));
    await click('[data-option-index]');
    await fillIn(GENERAL.inputByAttr('resourcePath'), 'recovered-secret');

    await click(GENERAL.button('recover'));
  });
});
