/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

import hbs from 'htmlbars-inline-precompile';
import engineDisplayData from 'vault/helpers/engines-display-data';

// The `configurable` array is hardcoded to validate that ALL_ENGINES metadata is correctly
// defined to render the tabs correctly.
const configurable = ['aws', 'azure', 'gcp', 'ldap', 'ssh'];

module('Integration | Component | secret-engines/configure-tabs', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.isConfigured = undefined;
    // This component accepts more args, but they are route related and instead asserted by acceptance tests
    this.renderComponent = (type) => {
      this.engineMetadata = type ? engineDisplayData(type) : undefined;
      return render(
        hbs`<SecretEngine::ConfigureTabs
          @engineMetadata={{this.engineMetadata}}
          @isConfigured={{this.isConfigured}}
        />`
      );
    };
  });

  test('it renders when args are undefined', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.tab('general-settings')).exists().hasText('General settings');
  });

  for (const { type } of filterEnginesByMountCategory({ mountCategory: 'secret' })) {
    if (configurable.includes(type)) {
      test(`${type} (configurable): it renders expected tabs when not configured`, async function (assert) {
        await this.renderComponent(type);
        assert.dom(GENERAL.tab('general-settings')).exists().hasText('General settings');
        assert
          .dom(GENERAL.tab('plugin-settings'))
          .exists()
          .hasText(`${this.engineMetadata.displayName} settings`);
      });

      test(`${type} (configurable): it renders expected tabs when configured`, async function (assert) {
        this.isConfigured = true;
        await this.renderComponent(type);

        assert.dom(GENERAL.tab('general-settings')).exists().hasText('General settings');
        assert
          .dom(GENERAL.tab('plugin-settings'))
          .exists()
          .hasText(`${this.engineMetadata.displayName} settings`);
      });
    } else {
      // NON-CONFIGURABLE ENGINES
      test(`${type} it hides plugin settings when not configurable`, async function (assert) {
        await this.renderComponent(type);

        assert.dom(GENERAL.tab('general-settings')).exists().hasText('General settings');
        assert.dom(GENERAL.tab('plugin-settings')).doesNotExist();
      });
    }
  }
});
