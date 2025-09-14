/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, triggerKeyEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | disabled-plugin-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.handleDisabledPluginClick = sinon.spy();
    this.handleDisabledPluginKeyDown = sinon.spy();
    this.type = {
      type: 'external-plugin',
      displayName: 'External Plugin',
      glyph: 'folder',
    };
  });

  test('it renders disabled plugin card', async function (assert) {
    await render(hbs`
      <DisabledPluginCard 
        @type={{this.type}} 
        @handleDisabledPluginClick={{this.handleDisabledPluginClick}}
        @handleDisabledPluginKeyDown={{this.handleDisabledPluginKeyDown}}
      />
    `);

    assert.dom(GENERAL.cardContainer('external-plugin')).exists('renders disabled plugin card');
    assert
      .dom(`${GENERAL.cardContainer('external-plugin')} h3`)
      .hasText('External Plugin', 'displays plugin name');
    assert
      .dom(`${GENERAL.cardContainer('external-plugin')} ${GENERAL.icon()}`)
      .exists('displays plugin icon');
    assert.dom(GENERAL.cardContainer('external-plugin')).exists('plugin card is rendered');
  });

  test('it handles click events', async function (assert) {
    await render(hbs`
      <DisabledPluginCard 
        @type={{this.type}} 
        @handleDisabledPluginClick={{this.handleDisabledPluginClick}}
        @handleDisabledPluginKeyDown={{this.handleDisabledPluginKeyDown}}
      />
    `);

    await click(GENERAL.cardContainer('external-plugin'));

    assert.ok(this.handleDisabledPluginClick.calledOnce, 'calls handleDisabledPluginClick on click');
    assert.ok(this.handleDisabledPluginClick.calledWith(this.type), 'passes correct plugin type');
  });

  test('it handles keyboard events', async function (assert) {
    await render(hbs`
      <DisabledPluginCard 
        @type={{this.type}} 
        @handleDisabledPluginClick={{this.handleDisabledPluginClick}}
        @handleDisabledPluginKeyDown={{this.handleDisabledPluginKeyDown}}
      />
    `);

    await triggerKeyEvent(GENERAL.cardContainer('external-plugin'), 'keydown', 'Enter');

    assert.ok(this.handleDisabledPluginKeyDown.calledOnce, 'calls handleDisabledPluginKeyDown on keydown');
    assert.ok(this.handleDisabledPluginKeyDown.calledWith(this.type), 'passes correct arguments');
  });

  test('it renders without icon when no glyph provided', async function (assert) {
    this.type = {
      type: 'no-icon-plugin',
      displayName: 'No Icon Plugin',
    };

    await render(hbs`
      <DisabledPluginCard 
        @type={{this.type}} 
        @handleDisabledPluginClick={{this.handleDisabledPluginClick}}
        @handleDisabledPluginKeyDown={{this.handleDisabledPluginKeyDown}}
      />
    `);

    assert.dom(GENERAL.icon()).doesNotExist('does not render icon when no glyph');
    assert
      .dom(GENERAL.cardContainer('no-icon-plugin'))
      .hasText('No Icon Plugin', 'still displays plugin name');
  });

  test('it renders and displays plugin information correctly', async function (assert) {
    await render(hbs`
      <DisabledPluginCard 
        @type={{this.type}} 
        @handleDisabledPluginClick={{this.handleDisabledPluginClick}}
        @handleDisabledPluginKeyDown={{this.handleDisabledPluginKeyDown}}
      />
    `);

    // Test that the basic elements are rendered
    assert.dom(GENERAL.cardContainer('external-plugin')).exists('renders disabled plugin card');

    assert
      .dom(`${GENERAL.cardContainer('external-plugin')} h3`)
      .hasText('External Plugin', 'displays correct plugin name');

    // Verify it's actually rendered as expected
    assert.dom(GENERAL.cardContainer('external-plugin')).exists();
  });
});
