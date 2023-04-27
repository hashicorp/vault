/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, skip } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | transform-role-edit', function (hooks) {
  setupRenderingTest(hooks);

  skip('it renders', async function (assert) {
    // TODO: Fill out these tests, merging without to unblock other work

    await render(hbs`{{transform-role-edit}}`);

    assert.dom(this.element).hasText('');

    // Template block usage:
    await render(hbs`
      {{#transform-role-edit}}
        template block text
      {{/transform-role-edit}}
    `);

    assert.dom(this.element).hasText('template block text');
  });
});
