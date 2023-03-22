/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const CARD_TITLE = 'Card title';
const ACTION_TEXT = 'View card';
const SUBTEXT = 'This is subtext for card';

module('Integration | Component overview-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('cardTitle', CARD_TITLE);
    this.set('actionText', ACTION_TEXT);
    this.set('subText', SUBTEXT);
  });

  test('it returns card title, ', async function (assert) {
    await render(hbs`<OverviewCard @cardTitle={{this.cardTitle}}/>`);
    const titleText = this.element.querySelector('.title').innerText;
    assert.strictEqual(titleText, 'Card title');
  });
  test('it returns card subtext, ', async function (assert) {
    await render(hbs`<OverviewCard @cardTitle={{this.cardTitle}}  @subText={{this.subText}} />`);
    const titleText = this.element.querySelector('p').innerText;
    assert.strictEqual(titleText, 'This is subtext for card');
  });
  test('it returns card action text', async function (assert) {
    await render(hbs`<OverviewCard @cardTitle={{this.cardTitle}} @actionText={{this.actionText}}/>`);
    const titleText = this.element.querySelector('a').innerText;
    assert.strictEqual(titleText, 'View card ');
  });
});
