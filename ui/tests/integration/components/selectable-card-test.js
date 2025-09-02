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
    this.onClick = sinon.spy();
  });

  test('it renders', async function (assert) {
    await render(hbs`<SelectableCard @onClick={{this.onClick}}/>`);
    await click('.selectable-card');
    assert.ok(this.onClick.calledOnce, 'calls on click');
  });

  test('it renders block content', async function (assert) {
    await render(hbs`<SelectableCard  @onClick={{this.onClick}}>hello</SelectableCard>`);
    assert.dom('.selectable-card').hasText('hello');
  });

  test('it does not process click event on disabled card', async function (assert) {
    await render(hbs`<SelectableCard @onClick={{this.onClick}} @disabled={{true}}>disabled</SelectableCard>`);
    await click('.selectable-card');
    assert.notOk(this.onClick.calledOnce, 'does not call the click event');
  });
});
