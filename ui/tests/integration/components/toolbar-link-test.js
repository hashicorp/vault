/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | toolbar-link', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<ToolbarLink @route="/secrets">Link</ToolbarLink>`);
    assert.dom('.toolbar-link').hasText('Link');
    assert.dom(GENERAL.icon('chevron-right')).exists('Default chevron right icon renders');
    assert.dom(GENERAL.tooltip('toolbar-link')).doesNotExist('Tooltip is hidden when not provided');
  });

  test('it should render icons', async function (assert) {
    await render(hbs`
      <ToolbarLink
        @route="/secrets"
        @type='add'
      >
        Test Link
      </ToolbarLink>
    `);

    assert.dom(GENERAL.icon('plus')).exists('Icon can be overridden to show plus sign');
  });
});
