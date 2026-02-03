/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, click, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | client count config', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(this.router, 'transitionTo');

    const { sys } = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(sys, 'internalClientActivityConfigure').resolves();

    this.renderComponent = (mode, config = {}) => {
      this.mode = mode;
      this.config = {
        enabled: 'enable',
        reporting_enabled: false,
        minimum_retention_months: 48,
        retention_months: 49,
        ...config,
      };
      return render(hbs`<Clients::Config @config={{this.config}} @mode={{this.mode}} />`);
    };
  });

  test('it shows the table with the correct rows by default', async function (assert) {
    await this.renderComponent('show');

    assert.dom('[data-test-clients-config-table]').exists('Clients config table exists');
    const rows = document.querySelectorAll('.info-table-row');
    assert.strictEqual(rows.length, 2, 'renders 2 info table rows');
    assert.ok(
      find('[data-test-row-value="Usage data collection"]').textContent.includes('On'),
      'Enabled value matches model'
    );
    assert.ok(
      find('[data-test-row-value="Retention period"]').textContent.includes('49'),
      'Retention period value matches model'
    );
  });

  test('it should validate retention_months', async function (assert) {
    await this.renderComponent('edit');

    assert.dom('[data-test-input="retention_months"]').hasValue('49', 'Retention months render');
    await fillIn('[data-test-input="retention_months"]', 20);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('retention_months'))
      .hasText(
        'Retention period must be greater than or equal to 48.',
        'Validation error shows for min retention period'
      );

    await fillIn('[data-test-input="retention_months"]', 90);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('retention_months'))
      .hasText(
        'Retention period must be less than or equal to 60.',
        'Validation error shows for max retention period'
      );
  });

  test('it should validate retention_months when minimum_retention_months is 0', async function (assert) {
    await this.renderComponent('edit', { minimum_retention_months: 0 });

    await fillIn('[data-test-input="retention_months"]', '');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('retention_months'))
      .hasText(
        'Retention period must be greater than or equal to 48.',
        'Validation error shows for min retention period'
      );
  });

  test('it should function in edit mode when enabling reporting', async function (assert) {
    const retention_months = 60;

    await this.renderComponent('edit', { enabled: 'disable' });

    assert.dom('[data-test-input="enabled"]').isNotChecked('Data collection checkbox is not checked');
    assert
      .dom('label[for="enabled"]')
      .hasText('Data collection is off', 'Correct label renders when data collection is off');

    await click('[data-test-input="enabled"]');
    await fillIn('[data-test-input="retention_months"]', retention_months);
    await click(GENERAL.submitButton);
    assert
      .dom('[data-test-clients-config-modal="title"]')
      .hasText('Turn usage tracking on?', 'Correct modal title renders');
    assert.dom('[data-test-clients-config-modal="on"]').exists('Correct modal description block renders');

    await click('[data-test-clients-config-modal="continue"]');
    assert.true(
      this.apiStub.calledWith({ enabled: 'enable', retention_months }),
      'API called with correct params'
    );
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.clients.config'),
      'Route transitions correctly on save success'
    );
  });

  test('it should function in edit mode when disabling reporting', async function (assert) {
    await this.renderComponent('edit');

    assert.dom('[data-test-input="enabled"]').isChecked('Data collection checkbox is checked');
    assert
      .dom('label[for="enabled"]')
      .hasText('Data collection is on', 'Correct label renders when data collection is on');

    await click('[data-test-input="enabled"]');
    await click(GENERAL.submitButton);
    assert.dom('[data-test-clients-config-modal]').exists('Modal renders');
    assert
      .dom('[data-test-clients-config-modal="title"]')
      .hasText('Turn usage tracking off?', 'Correct modal title renders');
    assert.dom('[data-test-clients-config-modal="off"]').exists('Correct modal description block renders');

    await click('[data-test-clients-config-modal="cancel"]');
    assert.dom('[data-test-clients-config-modal]').doesNotExist('Modal is hidden on cancel');

    await click(GENERAL.submitButton);
    await click('[data-test-clients-config-modal="continue"]');
    assert.true(
      this.apiStub.calledWith({ enabled: 'disable', retention_months: 49 }),
      'API called with correct params'
    );
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.clients.config'),
      'Route transitions correctly on save success'
    );
  });

  test('it should hide enabled field in edit mode when reporting is enabled', async function (assert) {
    const config = { enabled: 'enable', reporting_enabled: true, minimum_retention_months: 24 };
    await this.renderComponent('edit', config);

    assert.dom('[data-test-input="enabled"]').doesNotExist('Data collection input not shown');
    assert.dom('[data-test-input="retention_months"]').hasValue('49', 'Retention months render');

    await fillIn('[data-test-input="retention_months"]', 5);
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('retention_months'))
      .hasText(
        'Retention period must be greater than or equal to 24.',
        'Validation error shows for incorrect retention period'
      );

    await fillIn('[data-test-input="retention_months"]', 48);
    await click(GENERAL.submitButton);
    assert.true(
      this.apiStub.calledWith({ enabled: 'enable', retention_months: 48 }),
      'API called with correct params'
    );
  });

  test('it should not show modal when data collection has not changed', async function (assert) {
    await this.renderComponent('edit');

    await fillIn('[data-test-input="retention_months"]', 48);
    await click(GENERAL.submitButton);
    assert.true(
      this.apiStub.calledWith({ enabled: 'enable', retention_months: 48 }),
      'API called with correct params'
    );
  });
});
