/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/role', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.unloadAll('pki/role');
  });

  test('it should make request to correct endpoint when updating a record', async function (assert) {
    assert.expect(1);
    this.server.post('/pki-not-hardcoded/roles/pki-test', () => {
      assert.ok(true, 'POST request made to correct endpoint when updating a record');
    });

    this.store.pushPayload('pki/role', {
      modelName: 'pki/role',
      backend: 'pki-not-hardcoded',
      id: 'pki-test',
      name: 'pki-test',
    });
    const record = this.store.peekRecord('pki/role', 'pki-test');
    await record.save();
  });
});
