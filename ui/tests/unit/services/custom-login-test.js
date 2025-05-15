/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Service | custom-login', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const service = this.owner.lookup('service:custom-login');
    assert.notOk(service);
  });

  // Fetch list of rules unauth

  // Fetch list of rules auth

  // Fetch individual rule

  // Create rule

  // Delete rule

  // In admin namespace, should not be able to post rule for root
});
