/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | plugin-documentation-flyout', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClose = sinon.spy();
  });

  test('it does not render when closed', async function (assert) {
    this.isOpen = false;

    await render(hbs`
      <PluginDocumentationFlyout 
        @isOpen={{this.isOpen}}
        @onClose={{this.onClose}}
      />
    `);

    // Check that no flyout-related elements are present
    assert.dom('.hds-flyout').doesNotExist('no flyout container when closed');
    assert.dom(GENERAL.flyout).doesNotExist('no flyout element when closed');
    assert.dom('.hds-flyout__header').doesNotExist('no flyout header when closed');
    assert.dom('.hds-flyout__body').doesNotExist('no flyout body when closed');
  });

  test('it renders plugin-specific content', async function (assert) {
    this.isOpen = true;
    this.pluginName = 'aws';
    this.pluginType = 'secret';
    this.displayName = 'AWS';

    await render(hbs`
      <PluginDocumentationFlyout 
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @pluginType={{this.pluginType}}
        @displayName={{this.displayName}}
        @onClose={{this.onClose}}
      />
    `);

    assert.dom(GENERAL.flyout).exists('renders flyout when open');
    assert
      .dom(`${GENERAL.flyout} .hds-flyout__header`)
      .containsText('AWS plugin information', 'shows plugin-specific header');
    assert
      .dom(`${GENERAL.flyout} .hds-flyout__body`)
      .containsText('AWS secrets engine', 'shows plugin type in body');
    assert
      .dom(`${GENERAL.flyout} .hds-flyout__body`)
      .containsText('not currently enabled', 'explains plugin is not enabled');
  });

  test('it renders generic external plugins content', async function (assert) {
    this.isOpen = true;
    // No pluginName means this is for external plugins help

    await render(hbs`
      <PluginDocumentationFlyout 
        @isOpen={{this.isOpen}}
        @onClose={{this.onClose}}
      />
    `);

    assert.dom(GENERAL.flyout).exists('renders flyout when open');
    assert
      .dom(`${GENERAL.flyout} .hds-flyout__header`)
      .containsText('External plugins information', 'shows external plugins header');
    assert
      .dom(`${GENERAL.flyout} .hds-flyout__body`)
      .containsText('External plugins are plugins found in the plugin catalog', 'explains external plugins');
    assert.dom(GENERAL.linkTo('Register and enable external plugins')).exists('shows register plugins link');
    assert.dom(GENERAL.linkTo('Plugin development guide')).exists('shows plugin development link');
  });

  test('it renders auth method content', async function (assert) {
    this.isOpen = true;
    this.pluginName = 'ldap';
    this.pluginType = 'auth';
    this.displayName = 'LDAP';

    await render(hbs`
      <PluginDocumentationFlyout 
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @pluginType={{this.pluginType}}
        @displayName={{this.displayName}}
        @onClose={{this.onClose}}
      />
    `);

    assert
      .dom(`${GENERAL.flyout} .hds-flyout__body`)
      .containsText('LDAP auth method', 'shows auth method type in body');
  });

  test('it handles close action', async function (assert) {
    this.isOpen = true;

    await render(hbs`
      <PluginDocumentationFlyout 
        @isOpen={{this.isOpen}}
        @onClose={{this.onClose}}
      />
    `);

    await click(GENERAL.button('Close'));
    assert.ok(this.onClose.calledOnce, 'calls onClose when close button clicked');
  });

  test('it shows documentation links', async function (assert) {
    this.isOpen = true;
    this.pluginName = 'aws';

    await render(hbs`
      <PluginDocumentationFlyout 
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @onClose={{this.onClose}}
      />
    `);

    assert
      .dom(GENERAL.linkTo('Register and enable external plugins'))
      .exists('shows register plugins documentation link');
    assert
      .dom(GENERAL.linkTo('Register and enable external plugins'))
      .hasAttribute(
        'href',
        'https://developer.hashicorp.com/vault/docs/plugins/register',
        'has correct documentation URL'
      );
    assert
      .dom(GENERAL.linkTo('Register and enable external plugins'))
      .hasAttribute('target', '_blank', 'opens in new tab');
  });

  test('it uses fallback display name when not provided', async function (assert) {
    this.isOpen = true;
    this.pluginName = 'my-custom-plugin';
    // No displayName provided

    await render(hbs`
      <PluginDocumentationFlyout 
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @onClose={{this.onClose}}
      />
    `);

    assert
      .dom(`${GENERAL.flyout} .hds-flyout__header`)
      .containsText('my-custom-plugin plugin information', 'uses plugin name as fallback');
  });
});
