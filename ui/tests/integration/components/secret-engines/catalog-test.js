/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | secret-engines/catalog', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.setMountType = sinon.spy();
    this.pluginCatalogData = null;
    this.pluginCatalogError = false;
  });

  test('it renders secret engines catalog', async function (assert) {
    const expectedEngines = filterEnginesByMountCategory({
      mountCategory: 'secret',
      isEnterprise: false,
    }).filter((engine) => engine.type !== 'cubbyhole');

    // Dynamic assertion count: 1 for title + number of engines
    assert.expect(1 + expectedEngines.length);

    await render(
      hbs`<SecretEngines::Catalog 
        @setMountType={{this.setMountType}} 
        @pluginCatalogData={{this.pluginCatalogData}} 
        @pluginCatalogError={{this.pluginCatalogError}} 
      />`
    );

    assert.dom(GENERAL.breadcrumbs).exists('renders breadcrumbs');

    for (const engine of expectedEngines) {
      assert.dom(GENERAL.cardContainer(engine.type)).exists(`renders ${engine.displayName} engine card`);
    }
  });

  test('it calls setMountType when engine is selected', async function (assert) {
    await render(
      hbs`<SecretEngines::Catalog 
        @setMountType={{this.setMountType}} 
        @pluginCatalogData={{this.pluginCatalogData}} 
        @pluginCatalogError={{this.pluginCatalogError}} 
      />`
    );

    await click(GENERAL.cardContainer('kv'));

    assert.true(this.setMountType.calledOnce, 'setMountType was called');
    assert.true(this.setMountType.calledWith('kv'), 'setMountType was called with kv');
  });

  test('it shows plugin catalog error when provided', async function (assert) {
    this.pluginCatalogError = true;

    await render(
      hbs`<SecretEngines::Catalog 
        @setMountType={{this.setMountType}} 
        @pluginCatalogData={{this.pluginCatalogData}} 
        @pluginCatalogError={{this.pluginCatalogError}} 
      />`
    );

    assert.dom(GENERAL.inlineAlert).exists('shows plugin catalog error alert');
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'Plugin information unavailable Unable to fetch current plugin information. Using static plugin data instead. Some plugins may not show current details.',
        'shows correct error title'
      );
  });

  test('it shows flyout when clicking disabled plugin', async function (assert) {
    // Set up plugin catalog data that creates both enabled and disabled engines
    // An engine is disabled when it's not found in the plugin catalog detailed array
    this.pluginCatalogData = {
      detailed: [
        // Include only some engines, leaving others as "disabled"
        {
          name: 'kv',
          type: 'secret',
          builtin: true,
          deprecation_status: 'supported',
          version: 'v1.0.0',
        },
        // AWS engine is NOT included, so it will be marked as isAvailable: false
      ],
    };

    await render(
      hbs`<SecretEngines::Catalog 
        @setMountType={{this.setMountType}} 
        @pluginCatalogData={{this.pluginCatalogData}} 
        @pluginCatalogError={{this.pluginCatalogError}} 
      />`
    );

    // Initially, flyout should not be visible
    assert.dom(GENERAL.flyout).doesNotExist('flyout is not shown initially');

    // Find a disabled plugin card - since AWS is not in our catalog data,
    // it should be rendered as disabled
    const awsCard = document.querySelector(GENERAL.cardContainer('aws'));

    // Look for any disabled cards regardless of AWS card presence
    const disabledCards = document.querySelectorAll(
      '.selectable-engine-card.disabled, .selectable-engine-card[style*="opacity"]'
    );

    let clickedCard = false;

    if (awsCard) {
      await click(awsCard);
      clickedCard = true;

      // After clicking disabled plugin, flyout should appear
      assert.dom(GENERAL.flyout).exists('flyout appears after clicking disabled plugin');
    } else if (disabledCards.length > 0) {
      await click(disabledCards[0]);
      clickedCard = true;
      assert.dom(GENERAL.flyout).exists('flyout appears after clicking any disabled plugin');
    }

    // Always verify we completed the test successfully
    assert.ok(clickedCard, 'successfully clicked a disabled plugin card');
  });
});
