/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { looseEqual } from 'core/helpers/loose-equal';

module('Integration | Helper | loose-equal', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.inputValue = 1234;
    await render(hbs`{{if (loose-equal "1234" 1234) "true" "false"}}`);
    assert.dom(this.element).hasText('true');

    this.inputValue = '4567';
    await render(hbs`{{if (loose-equal "1234" "4567") "true" "false"}}`);
    assert.dom(this.element).hasText('false');
  });

  test('it compares values as expected', async function (assert) {
    assert.true(looseEqual([0, '0']));
    assert.true(looseEqual([0, 0]));
    assert.true(looseEqual(['0', '0']));
    assert.true(looseEqual(['1234', 1234]));
    assert.true(looseEqual(['1234', '1234']));
    assert.true(looseEqual([1234, 1234]));
    assert.true(looseEqual(['abc', 'abc']));
    assert.true(looseEqual(['', '']));

    // == normally returns true for this comparison, we intercept and return false
    assert.false(looseEqual(['', 0]));
    assert.false(looseEqual([0, '']));
  });
});
