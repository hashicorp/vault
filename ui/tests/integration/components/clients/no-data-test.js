/**
 * Copyright (c) HashiCorp, Inc.
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
    this.store = this.owner.lookup('service:store');
    this.setConfig = async (data) => {
      // the clients/config model does some funky serializing for the "enabled" param
      // so stubbing the request here instead of just the model for additional coverage
      this.server.get('sys/internal/counters/config', () => {
        return {
          request_id: '25a94b99-b49a-c4ac-cb7b-5ba0eb390a25',
          data,
        };
      });
      return this.store.queryRecord('clients/config', {});
    };
    this.renderComponent = async () => {
      return render(hbs`<Clients::NoData @config={{this.config}} />`);
    };
  });

  test('it renders empty state when enabled is "on"', async function (assert) {
    assert.expect(2);
    const data = {
      enabled: 'default-enabled',
      reporting_enabled: false,
    };
    ``;
    this.config = await this.setConfig(data);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No data received');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Tracking is turned on and Vault is gathering data. It should appear here within 30 minutes.');
  });

  test('it renders empty state when reporting_enabled is true', async function (assert) {
    assert.expect(2);
    const data = {
      enabled: 'default-disabled',
      reporting_enabled: true,
    };
    this.config = await this.setConfig(data);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No data received');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Tracking is turned on and Vault is gathering data. It should appear here within 30 minutes.');
  });

  test('it renders empty state when reporting is fully disabled', async function (assert) {
    assert.expect(2);
    const data = {
      enabled: 'default-disabled',
      reporting_enabled: false,
    };
    this.config = await this.setConfig(data);
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('Data tracking is disabled');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Tracking is disabled, and no data is being collected. To turn it on, edit the configuration.'
      );
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
