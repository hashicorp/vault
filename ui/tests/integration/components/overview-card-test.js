/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const CARD_TITLE = 'Card title';
const ACTION_TEXT = 'View card';
const SUBTEXT = 'This is subtext for card';

const SELECTORS = {
  container: '[data-test-overview-card-container]',
  title: '[data-test-overview-card-title]',
  subtitle: '[data-test-overview-card-subtitle]',
  action: '[data-test-action-text]',
  customSubtext: '[data-test-custom-subtext]',
};

module('Integration | Component | overview-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('cardTitle', CARD_TITLE);
    this.set('actionText', ACTION_TEXT);
    this.set('subText', SUBTEXT);
  });

  test('it returns card title', async function (assert) {
    await render(hbs`<OverviewCard @cardTitle={{this.cardTitle}}/>`);
    assert.dom(SELECTORS.title).hasText('Card title');
  });
  test('it returns custom title if both exist', async function (assert) {
    await render(hbs`
      <OverviewCard @cardTitle={{this.cardTitle}}>
        <:customTitle>
          Fancy custom title
        </:customTitle>
      </OverviewCard>
      `);
    assert.dom(SELECTORS.container).hasText('Fancy custom title');
    assert.dom(SELECTORS.container).doesNotIncludeText(this.cardTitle);
  });
  test('it renders card @subText arg, ', async function (assert) {
    await render(hbs`<OverviewCard @cardTitle={{this.cardTitle}}  @subText={{this.subText}} />`);
    assert.dom(SELECTORS.subtitle).hasText('This is subtext for card');
  });
  test('it renders card action text', async function (assert) {
    await render(
      hbs`
      <OverviewCard @cardTitle={{this.cardTitle}}>
        <:action>
        <div data-test-action-text>
        {{this.actionText}}
        </div>
        </:action>
      </OverviewCard>
      `
    );
    assert.dom(SELECTORS.action).hasText('View card');
  });
  test('it renders custom subtext text', async function (assert) {
    await render(
      hbs`
      <OverviewCard @cardTitle={{this.cardTitle}}>
        <:customSubtext>
          <div data-test-custom-subtext>
            Fancy yielded subtext
          </div>
        </:customSubtext>
      </OverviewCard>
      `
    );
    assert.dom(SELECTORS.customSubtext).hasText('Fancy yielded subtext');
  });
});
