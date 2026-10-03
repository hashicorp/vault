/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | mount/configure-tabs', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.configRoute = 'pki.configuration.create';
    // This component also accepts @path but that's route related and irrelevant to an integration test
    this.renderComponent = () => {
      return render(
        hbs`<Mount::ConfigureTabs
          @displayName={{this.displayName}}
          @configRoute={{this.configRoute}}
          @path="my-pki-engine"
        />`
      );
    };
  });

  test('it renders when args are undefined', async function (assert) {
    this.configRoute = undefined;
    this.displayName = undefined;
    await this.renderComponent();
    assert.dom(GENERAL.tab('general-settings')).exists().hasText('General settings');
  });

  test('it renders plugin settings tab when @configRoute provided', async function (assert) {
    this.displayName = 'PKI';
    await this.renderComponent();

    assert.dom(GENERAL.tab('general-settings')).exists().hasText('General settings');
    await click(GENERAL.tab('plugin-settings'));
    assert.dom(GENERAL.tab('plugin-settings')).exists().hasText('PKI settings');
  });

  test('it renders fallback when @displayName not provided', async function (assert) {
    this.displayName = '';
    await this.renderComponent();
    assert.dom(GENERAL.tab('plugin-settings')).exists().hasText('Plugin settings');
  });

  test('it hides plugin settings when there is no @configRoute', async function (assert) {
    this.configRoute = '';
    await this.renderComponent();

    assert.dom(GENERAL.tab('general-settings')).exists().hasText('General settings');
    assert.dom(GENERAL.tab('plugin-settings')).doesNotExist();
  });
});
