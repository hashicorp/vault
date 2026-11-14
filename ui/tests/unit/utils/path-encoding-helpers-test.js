/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { recoverEngineMountPath } from 'vault/utils/path-encoding-helpers';

module('Unit | Util | path-encoding-helpers', function () {
  test('recoverEngineMountPath restores %2F decoded by Ember', function (assert) {
    assert.strictEqual(recoverEngineMountPath('myeng/with/slash'), 'myeng%2Fwith%2Fslash');
  });

  test('recoverEngineMountPath leaves normal mounts unchanged', function (assert) {
    assert.strictEqual(recoverEngineMountPath('secret'), 'secret');
    assert.strictEqual(recoverEngineMountPath('myeng%2Fwith%2Fslash'), 'myeng%2Fwith%2Fslash');
  });
});
