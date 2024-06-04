/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, triggerEvent } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-link', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<ToolbarLink @route="/secrets">Link</ToolbarLink>`);

    assert.dom(this.element).hasText('Link');
    assert.ok(isPresent('.toolbar-link'));
    assert.ok(isPresent('.icon'));
  });

  test('it should render icons', async function (assert) {
    assert.expect(2);

    await render(hbs`
      <ToolbarLink
        @route="/secrets"
        @type={{this.type}}
      >
        Test Link
      </ToolbarLink>
    `);

    assert.dom('[data-test-icon="chevron-right"]').exists('Default chevron right icon renders');
    this.set('type', 'add');
    assert.dom('[data-test-icon="plus"]').exists('Icon can be overriden to show plus sign');
  });

  test('it should disable and show tooltip if provided', async function (assert) {
    assert.expect(3);

    await render(hbs`
      <ToolbarLink
        @route="/secrets"
        @disabled={{true}}
        @disabledTooltip={{this.tooltip}}
      >
        Test Link
      </ToolbarLink>
    `);

    assert.dom('a').hasClass('disabled', 'Link can be disabled');
    assert.dom('[data-test-popup-menu-trigger]').doesNotExist('Tooltip is hidden when not provided');
    this.set('tooltip', 'Test tooltip');
    await triggerEvent('.ember-basic-dropdown-trigger', 'mouseenter');
    assert.dom('[data-test-disabled-tooltip]').hasText(this.tooltip, 'Tooltip renders');
  });
});
