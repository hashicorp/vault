/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | plugin-documentation-flyout', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders documentation flyout for secret engine', async function (assert) {
    this.set('isOpen', true);
    this.set('pluginName', 'aws');
    this.set('pluginType', 'secret');
    this.set('displayName', 'AWS Secrets Engine');
    this.set('onClose', () => {});

    await render(hbs`
      <PluginDocumentationFlyout
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @pluginType={{this.pluginType}}
        @displayName={{this.displayName}}
        @onClose={{this.onClose}}
      />
    `);

    assert.dom('[data-test-modal]').exists('Modal is rendered');
    assert
      .dom('[data-test-modal-header]')
      .containsText('AWS Secrets Engine Plugin Information', 'Header shows plugin name');
    assert
      .dom('[data-test-modal-body]')
      .containsText('secrets engine is not currently enabled', 'Body explains plugin status');
    assert.dom('[data-test-code-block]').containsText('vault secrets enable aws', 'CLI command is shown');
  });

  test('it renders documentation flyout for auth method', async function (assert) {
    this.set('isOpen', true);
    this.set('pluginName', 'github');
    this.set('pluginType', 'auth');
    this.set('displayName', 'GitHub Auth Method');
    this.set('onClose', () => {});

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
      .dom('[data-test-modal-body]')
      .containsText('auth method is not currently enabled', 'Body explains auth method status');
    assert
      .dom('[data-test-code-block]')
      .containsText('vault auth enable github', 'Auth CLI command is shown');
  });

  test('it generates correct documentation URL', async function (assert) {
    this.set('isOpen', true);
    this.set('pluginName', 'aws');
    this.set('pluginType', 'secret');
    this.set('onClose', () => {});

    await render(hbs`
      <PluginDocumentationFlyout
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @pluginType={{this.pluginType}}
        @onClose={{this.onClose}}
      />
    `);

    assert
      .dom('[data-test-documentation-link]')
      .hasAttribute(
        'href',
        'https://developer.hashicorp.com/vault/docs/secrets/aws',
        'Documentation link points to correct URL'
      );
    assert
      .dom('[data-test-register-plugins-link]')
      .hasAttribute(
        'href',
        'https://developer.hashicorp.com/vault/docs/plugins/plugin-development',
        'Plugin development link points to correct URL'
      );
    assert
      .dom('[data-test-register-plugins-link]')
      .containsText('Plugin Development Guide', 'Plugin development link has correct text');
  });

  test('it calls onClose when close button is clicked', async function (assert) {
    assert.expect(1);

    this.set('isOpen', true);
    this.set('pluginName', 'aws');
    this.set('pluginType', 'secret');
    this.set('onClose', () => {
      assert.ok(true, 'onClose callback is called');
    });

    await render(hbs`
      <PluginDocumentationFlyout
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @pluginType={{this.pluginType}}
        @onClose={{this.onClose}}
      />
    `);

    await click('[data-test-modal-close]');
  });

  test('it uses pluginName as displayName when displayName is not provided', async function (assert) {
    this.set('isOpen', true);
    this.set('pluginName', 'aws');
    this.set('pluginType', 'secret');
    this.set('onClose', () => {});

    await render(hbs`
      <PluginDocumentationFlyout
        @isOpen={{this.isOpen}}
        @pluginName={{this.pluginName}}
        @pluginType={{this.pluginType}}
        @onClose={{this.onClose}}
      />
    `);

    assert
      .dom('[data-test-modal-header]')
      .containsText('aws Plugin Information', 'Uses plugin name when display name not provided');
  });
});
