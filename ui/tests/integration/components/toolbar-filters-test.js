/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { isPresent } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | toolbar-filters', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<ToolbarFilters>These are the toolbar filters</ToolbarFilters>`);

    assert.dom(this.element).hasText('These are the toolbar filters');
    assert.ok(isPresent('.toolbar-filters'));
  });
});
