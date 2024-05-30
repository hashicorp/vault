/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | console', function (hooks) {
  setupTest(hooks);

  test('it builds the correct URL', function (assert) {
    const adapter = this.owner.lookup('adapter:console');
    const sysPath = 'sys/health';
    const awsPath = 'aws/roles/my-other-role';
    assert.strictEqual(adapter.buildURL(sysPath), '/v1/sys/health');
    assert.strictEqual(adapter.buildURL(awsPath), '/v1/aws/roles/my-other-role');
  });
});
