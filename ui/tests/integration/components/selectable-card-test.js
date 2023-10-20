/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component selectable-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.cardTitle = 'Connections';
    this.onClick = sinon.spy();
  });

  test('it renders', async function (assert) {
    await render(hbs`<SelectableCard @cardTitle={{this.cardTitle}} @onClick={{this.onClick}}/>`);
    const titleText = this.element.querySelector('h3').innerText;
    assert.strictEqual(titleText, 'Connections');
    await click(this.element);
    assert.ok(this.onClick.calledOnce(), 'calls on click');
  });

  test('it renders block content', async function (assert) {
    await render(hbs`<SelectableCard  @onClick={{this.onClick}} >hello</SelectableCard>`);
    assert.strictEqual(this.element.innerText, 'hello');
  });
});
