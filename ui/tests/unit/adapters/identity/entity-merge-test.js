/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { storeMVP } from './_test-cases';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | identity/entity-merge', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  test(`entity-merge#createRecord`, function (assert) {
    assert.expect(2);
    this.server.post('/identity/entity/merge', (_, req) => {
      const { url, method } = req;
      assert.strictEqual(url, `/v1/identity/entity/merge`, ` calls the correct url`);
      assert.strictEqual(method, 'POST', `uses the correct http verb: POST`);
      return {};
    });
    const adapter = this.owner.lookup('adapter:identity/entity-merge');
    adapter.createRecord(storeMVP, { modelName: 'identity/entity-merge' }, { attr: (x) => x });
  });
});
