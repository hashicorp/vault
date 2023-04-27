/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Unit | Adapter | oidc/key', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'oidc/key';
    this.data = {
      name: 'foo-key',
      rotation_period: '12h',
      verification_ttl: 43200,
    };
    this.path = '/identity/oidc/key/foo-key';
  });

  testHelper(test);

  test('it should make request to correct endpoint on rotate', async function (assert) {
    assert.expect(1);

    this.server.post(`${this.path}/rotate`, (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.strictEqual(json.verification_ttl, '30m', 'request made to correct endpoint on rotate');
    });

    await this.store.adapterFor('oidc/key').rotate(this.data.name, '30m');
  });
});
