/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import sinon from 'sinon';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | box-radio', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('type', 'aws');
    this.set('displayName', 'An Option');
    this.set('mountType', '');
    this.set('disabled', false);
  });

  test('it renders', async function (assert) {
    const spy = sinon.spy();
    this.set('onRadioChange', spy);
    await render(hbs`<BoxRadio
      @type={{this.type}}
      @glyph={{this.type}}
      @displayName={{this.displayName}}
      @onRadioChange={{this.onRadioChange}}
      @disabled={{this.disabled}}
    />`);

    assert.dom(this.element).hasText('An Option', 'shows the display name of the option');
    assert.dom('.tooltip').doesNotExist('tooltip does not exist when disabled is false');
    await click('[data-test-mount-type="aws"]');
    assert.ok(spy.calledOnce, 'calls the radio change function when option clicked');
  });

  test('it renders correctly when disabled', async function (assert) {
    const spy = sinon.spy();
    this.set('onRadioChange', spy);
    await render(hbs`<BoxRadio
      @type={{this.type}}
      @glyph={{this.type}}
      @displayName={{this.displayName}}
      @onRadioChange={{this.onRadioChange}}
      @disabled={{true}}
    />`);

    assert.dom(this.element).hasText('An Option', 'shows the display name of the option');
    assert.dom('.ember-basic-dropdown-trigger').exists('tooltip exists');
    await click('[data-test-mount-type="aws"]');
    assert.ok(spy.notCalled, 'does not call the radio change function when option is clicked');
  });
});
