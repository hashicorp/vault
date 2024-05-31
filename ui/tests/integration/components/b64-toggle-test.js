/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | b64 toggle', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`{{b64-toggle}}`);
    assert.dom('button').exists({ count: 1 });
  });

  test('it toggles encoding on the passed string', async function (assert) {
    this.set('value', 'value');
    await render(hbs`{{b64-toggle value=this.value}}`);
    await click('button');
    assert.strictEqual(this.value, btoa('value'), 'encodes to base64');
    await click('button');
    assert.strictEqual(this.value, 'value', 'decodes from base64');
  });

  test('it toggles encoding starting with base64', async function (assert) {
    this.set('value', btoa('value'));
    await render(hbs`{{b64-toggle value=this.value initialEncoding='base64'}}`);
    assert.ok(find('button').textContent.includes('Decode'), 'renders as on when in b64 mode');
    await click('button');
    assert.strictEqual(this.value, 'value', 'decodes from base64');
  });

  test('it detects changes to value after encoding', async function (assert) {
    this.set('value', btoa('value'));
    await render(hbs`{{b64-toggle value=this.value initialEncoding='base64'}}`);
    assert.ok(find('button').textContent.includes('Decode'), 'renders as on when in b64 mode');
    this.set('value', btoa('value') + '=');
    assert.ok(find('button').textContent.includes('Encode'), 'toggles off since value has changed');
    this.set('value', btoa('value'));
    assert.ok(
      find('button').textContent.includes('Decode'),
      'toggles on since value is equal to the original'
    );
  });

  test('it does not toggle when the value is empty', async function (assert) {
    this.set('value', '');
    await render(hbs`{{b64-toggle value=this.value}}`);
    await click('button');
    assert.ok(find('button').textContent.includes('Encode'));
  });
});
