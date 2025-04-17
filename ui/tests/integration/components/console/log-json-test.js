/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { assertCodeBlockValue } from 'vault/tests/helpers/codemirror';

module('Integration | Component | console/log json', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.on('myAction', function(val) { ... });
    const objectContent = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
    const expectedText = JSON.stringify(objectContent, null, 2);

    this.set('content', objectContent);

    await render(hbs`<Console::LogJson @content={{this.content}} />`);
    assertCodeBlockValue(assert, '.hds-code-block__code', expectedText);
  });
});
