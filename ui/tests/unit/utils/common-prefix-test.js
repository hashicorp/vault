/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import commonPrefix from 'core/utils/common-prefix';
import { module, test } from 'qunit';

module('Unit | Util | common prefix', function () {
  test('it returns empty string if called with no args or an empty array', function (assert) {
    let returned = commonPrefix();
    assert.strictEqual(returned, '', 'returns an empty string');
    returned = commonPrefix([]);
    assert.strictEqual(returned, '', 'returns an empty string for an empty array');
  });

  test('it returns empty string if there are no common prefixes', function (assert) {
    const secrets = ['asecret', 'secret2', 'secret3'].map((s) => ({ id: s }));
    const returned = commonPrefix(secrets);
    assert.strictEqual(returned, '', 'returns an empty string');
  });

  test('it returns the longest prefix', function (assert) {
    const secrets = ['secret1', 'secret2', 'secret3'].map((s) => ({ id: s }));
    let returned = commonPrefix(secrets);
    assert.strictEqual(returned, 'secret', 'finds secret prefix');
    const greetings = ['hello-there', 'hello-hi', 'hello-howdy'].map((s) => ({ id: s }));
    returned = commonPrefix(greetings);
    assert.strictEqual(returned, 'hello-', 'finds hello- prefix');
  });

  test('it can compare an attribute that is not "id" to calculate the longest prefix', function (assert) {
    const secrets = ['secret1', 'secret2', 'secret3'].map((s) => ({ name: s }));
    const returned = commonPrefix(secrets, 'name');
    assert.strictEqual(returned, 'secret', 'finds secret prefix from name attribute');
  });
});
