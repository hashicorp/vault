/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Pretender from 'pretender';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { storeMVP } from './_test-cases';

module('Unit | Adapter | identity/entity-merge', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.server = new Pretender(function () {
      this.post('/v1/**', (response) => {
        return [response, { 'Content-Type': 'application/json' }, JSON.stringify({})];
      });
    });
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test(`entity-merge#createRecord`, function (assert) {
    assert.expect(2);
    const adapter = this.owner.lookup('adapter:identity/entity-merge');
    adapter.createRecord(storeMVP, { modelName: 'identity/entity-merge' }, { attr: (x) => x });
    const { url, method } = this.server.handledRequests[0];
    assert.strictEqual(url, `/v1/identity/entity/merge`, ` calls the correct url`);
    assert.strictEqual(method, 'POST', `uses the correct http verb: POST`);
  });
});
