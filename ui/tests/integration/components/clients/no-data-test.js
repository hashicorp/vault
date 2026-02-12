/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | clients/no-data', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.canUpdate = false;
    this.setConfig = (enabled, reporting_enabled) => ({
      enabled: enabled ? 'default-enabled' : 'default-disabled',
      reporting_enabled,
    });
    this.renderComponent = async () => {
      return render(hbs`<Clients::NoData @config={{this.config}} @canUpdate={{this.canUpdate}} />`);
    };
  });

  test('it renders empty state when enabled', async function (assert) {
    assert.expect(2);
    this.config = this.setConfig(true, false);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No data received');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Tracking is turned on and Vault is gathering data. It should appear here within 30 minutes.');
  });

  test('it renders empty state when reporting is fully enabled', async function (assert) {
    assert.expect(2);
    this.config = this.setConfig(true, true);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No data received');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Tracking is turned on and Vault is gathering data. It should appear here within 30 minutes.');
  });

  test('it renders empty state when reporting is fully disabled', async function (assert) {
    assert.expect(4);
    this.config = this.setConfig(false, false);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('Data tracking is disabled');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Tracking is disabled, and no data is being collected. To turn it on, edit the configuration.'
      );
    assert.dom(GENERAL.linkTo('config')).doesNotExist('Config link does not render without capabilities');

    this.canUpdate = true;
    await this.renderComponent();
    assert.dom(GENERAL.linkTo('config')).exists('Config link renders with update capabilities');
  });

  test('it renders empty state when config data is not available', async function (assert) {
    assert.expect(2);
    this.config = null;
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('Activity configuration data is unavailable');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Reporting status is unknown and could be enabled or disabled. Check the Vault logs for more information.'
      );
  });
});
