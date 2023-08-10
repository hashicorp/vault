/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { stringifyObjectValues } from 'vault/components/console/log-object';

module('Integration | Component | console/log object', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    const objectContent = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
    const data = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
    stringifyObjectValues(data);
    const expectedText = 'Key Value one two three four seven {"five":"six"} eight [5,6]';
    this.set('content', objectContent);

    await render(hbs`{{console/log-object content=this.content}}`);
    assert.dom('pre').includesText(expectedText);
  });
});
