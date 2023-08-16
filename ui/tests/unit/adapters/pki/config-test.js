/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/config', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'pki-engine';
  });

  const testHelper = (test) => {
    test('it should make request to correct endpoint on update', async function (assert) {
      assert.expect(1);

      this.server.post(`/${this.backend}/config/${this.endpoint}`, () => {
        assert.ok(true, `request made to POST config/${this.endpoint} endpoint on update`);
      });

      this.store.pushPayload(`pki/config/${this.endpoint}`, {
        modelName: `pki/config/${this.endpoint}`,
        id: this.backend,
      });

      const model = this.store.peekRecord(`pki/config/${this.endpoint}`, this.backend);
      await model.save();
    });

    test('it should make request to correct endpoint on find', async function (assert) {
      assert.expect(1);

      this.server.get(`/${this.backend}/config/${this.endpoint}`, () => {
        assert.ok(true, `request is made to GET /config/${this.endpoint} endpoint on find`);
        return { data: { id: this.backend } };
      });

      this.store.findRecord(`pki/config/${this.endpoint}`, this.backend);
    });
  };

  module('cluster', function (hooks) {
    hooks.beforeEach(async function () {
      this.endpoint = 'cluster';
    });
    testHelper(test);
  });

  module('urls', function (hooks) {
    hooks.beforeEach(async function () {
      this.endpoint = 'urls';
    });
    testHelper(test);
  });

  module('crl', function (hooks) {
    hooks.beforeEach(async function () {
      this.endpoint = 'crl';
    });
    testHelper(test);
  });
});
