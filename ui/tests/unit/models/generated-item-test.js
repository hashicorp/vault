/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';

import { setupTest } from 'vault/tests/helpers';

module('Unit | Model | generated item', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('generated-item', {});
    assert.ok(model, 'generated-item model exists');
  });
});
