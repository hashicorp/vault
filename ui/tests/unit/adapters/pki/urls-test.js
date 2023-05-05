/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/urls', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'pki-engine';
  });

  test('it should make request to correct endpoint on update', async function (assert) {
    assert.expect(1);

    this.server.post(`/${this.backend}/config/urls`, () => {
      assert.ok(true, 'request made to correct endpoint on update');
    });

    this.store.pushPayload('pki/urls', {
      modelName: 'pki/urls',
      id: this.backend,
    });

    const model = this.store.peekRecord('pki/urls', this.backend);
    await model.save();
  });

  test('it should make request to correct endpoint on find', async function (assert) {
    assert.expect(1);

    this.server.get(`/${this.backend}/config/urls`, () => {
      assert.ok(true, 'request is made to correct endpoint on find');
      return { data: { id: this.backend } };
    });

    this.store.findRecord('pki/urls', this.backend);
  });
});
