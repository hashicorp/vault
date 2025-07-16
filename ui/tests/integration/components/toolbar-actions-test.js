/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-actions', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<ToolbarActions>These are the toolbar actions</ToolbarActions>`);

    assert.dom(this.element).hasText('These are the toolbar actions');
    assert.dom('.toolbar-actions').exists();
  });
});
