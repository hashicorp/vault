/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | tool-tip', function (hooks) {
  setupRenderingTest(hooks);

  test.skip('it should open and close on mouse enter and leave', async function (assert) {
    await render(hbs`
      <ToolTip as |T|>
        <T.Trigger>
          <span data-test-trigger>Trigger</span>
        </T.Trigger>
        <T.Content>
          <span data-test-content>Tooltip</span>
        </T.Content>
      </ToolTip>
    `);

    assert.dom('[data-test-content]').doesNotExist('Tooltip is hidden');
    await triggerEvent('[data-test-trigger]', 'mouseenter');
    assert.dom('[data-test-content]').exists('Tooltip renders on mouse enter');
    await triggerEvent('[data-test-content]', 'mouseleave');
    assert.dom('[data-test-content]').doesNotExist('Tooltip is hidden on mouse leave');
  });
});
