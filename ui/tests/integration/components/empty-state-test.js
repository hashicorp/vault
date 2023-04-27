/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | empty-state', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`{{empty-state}}`);

    assert.dom(this.element).hasText('');

    // Template block usage:
    await render(hbs`
      {{#empty-state
        title="Empty State Title"
        message="This is the empty state message"
      }}
        Actions Link
      {{/empty-state}}
    `);

    assert.dom('.empty-state-title').hasText('Empty State Title', 'renders empty state title');
    assert
      .dom('.empty-state-message')
      .hasText('This is the empty state message', 'renders empty state message');
    assert.dom('.empty-state-actions').hasText('Actions Link', 'renders empty state actions');
  });
});
