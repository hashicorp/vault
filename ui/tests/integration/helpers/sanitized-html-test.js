/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Helper | sanitized-html', function (hooks) {
  setupRenderingTest(hooks);

  test('it does not alter if string is safe', async function (assert) {
    this.set('inputValue', 'height: 15.33px');

    await render(hbs`{{sanitized-html this.inputValue}}`);
    assert.dom(this.element).hasText('height: 15.33px');
  });

  test('it strips unsafe HTML before rendering safe HTML', async function (assert) {
    this.set(
      'inputValue',
      '<main data-test-thing>This is something<script data-test-script>window.alert(`h4cK3d`)</script></main>'
    );

    await render(hbs`{{sanitized-html this.inputValue}}`);
    assert.dom('[data-test-thing]').hasTagName('main');
    assert.dom('[data-test-thing]').hasText('This is something', 'preserves non-problematic content');
    assert.dom('[data-test-script]').doesNotExist('Script is stripped from render');
  });

  test('it does not invoke functions passed as value', async function (assert) {
    this.set('inputValue', () => {
      window.alert('h4cK3d');
    });
    await render(hbs`{{sanitized-html this.inputValue}}`);
    assert.dom(this.element).hasText("() => { window.alert('h4cK3d'); }");
  });
});
