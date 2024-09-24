/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, click, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | client count config', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(this.router, 'transitionTo');
    const store = this.owner.lookup('service:store');
    this.createModel = (enabled = 'enable', reporting_enabled = false, minimum_retention_months = 48) => {
      store.pushPayload('clients/config', {
        modelName: 'clients/config',
        id: 'foo',
        data: {
          enabled,
          reporting_enabled,
          minimum_retention_months,
          retention_months: 49,
        },
      });
      this.model = store.peekRecord('clients/config', 'foo');
    };
  });

  test('it shows the table with the correct rows by default', async function (assert) {
    this.createModel();

    await render(hbs`<Clients::Config @model={{this.model}} />`);

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

  test('it should function in edit mode when reporting is disabled', async function (assert) {
    assert.expect(13);
    const retentionMonths = 60;
    this.server.put('/sys/internal/counters/config', (schema, req) => {
      const { enabled, retention_months } = JSON.parse(req.requestBody);
      const expected = { enabled: 'enable', retention_months: retentionMonths };
      assert.deepEqual({ enabled, retention_months }, expected, 'Correct data sent in PUT request (1)');
      return {};
    });

    this.createModel('disable');

    await render(hbs`
      <Clients::Config @model={{this.model}} @mode="edit" />
    `);

    assert.dom('[data-test-input="enabled"]').isNotChecked('Data collection checkbox is not checked');
    assert
      .dom('label[for="enabled"]')
      .hasText('Data collection is off', 'Correct label renders when data collection is off');
    assert.dom('[data-test-input="retentionMonths"]').hasValue('49', 'Retention months render');

    await click('[data-test-input="enabled"]');
    await fillIn('[data-test-input="retentionMonths"]', 20);
    await click('[data-test-clients-config-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'Retention period must be greater than or equal to 48.',
        'Validation error shows for min retention period'
      );
    await fillIn('[data-test-input="retentionMonths"]', 90);
    await click('[data-test-clients-config-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'Retention period must be less than or equal to 60.',
        'Validation error shows for max retention period'
      );

    await fillIn('[data-test-input="retentionMonths"]', retentionMonths);
    await click('[data-test-clients-config-save]');
    assert
      .dom('[data-test-clients-config-modal="title"]')
      .hasText('Turn usage tracking on?', 'Correct modal title renders');
    assert.dom('[data-test-clients-config-modal="on"]').exists('Correct modal description block renders');

    await click('[data-test-clients-config-modal="continue"]');
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.clients.config'),
      'Route transitions correctly on save success'
    );

    await click('[data-test-input="enabled"]');
    await click('[data-test-clients-config-save]');
    assert.dom('[data-test-clients-config-modal]').exists('Modal renders');
    assert
      .dom('[data-test-clients-config-modal="title"]')
      .hasText('Turn usage tracking off?', 'Correct modal title renders');
    assert.dom('[data-test-clients-config-modal="off"]').exists('Correct modal description block renders');

    await click('[data-test-clients-config-modal="cancel"]');
    assert.dom('[data-test-clients-config-modal]').doesNotExist('Modal is hidden on cancel');
  });

  test('it should be hidden in edit mode when reporting is enabled', async function (assert) {
    assert.expect(4);

    this.server.put('/sys/internal/counters/config', (schema, req) => {
      const { enabled, retention_months } = JSON.parse(req.requestBody);
      const expected = { enabled: 'enable', retention_months: 48 };
      assert.deepEqual({ enabled, retention_months }, expected, 'Correct data sent in PUT request (2)');
      return {};
    });

    this.createModel('enable', true, 24);

    await render(hbs`
      <Clients::Config @model={{this.model}} @mode="edit" />
    `);

    assert.dom('[data-test-input="enabled"]').doesNotExist('Data collection input not shown ');
    assert.dom('[data-test-input="retentionMonths"]').hasValue('49', 'Retention months render');

    await fillIn('[data-test-input="retentionMonths"]', 5);
    await click('[data-test-clients-config-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'Retention period must be greater than or equal to 24.',
        'Validation error shows for incorrect retention period'
      );

    await fillIn('[data-test-input="retentionMonths"]', 48);
    await click('[data-test-clients-config-save]');
  });

  test('it should not show modal when data collection is not changed', async function (assert) {
    assert.expect(1);

    this.server.put('/sys/internal/counters/config', (schema, req) => {
      const { enabled, retention_months } = JSON.parse(req.requestBody);
      const expected = { enabled: 'enable', retention_months: 48 };
      assert.deepEqual({ enabled, retention_months }, expected, 'Correct data sent in PUT request (3)');
      return {};
    });

    this.createModel();

    await render(hbs`
      <Clients::Config @model={{this.model}} @mode="edit" />
    `);
    await fillIn('[data-test-input="retentionMonths"]', 48);
    await click('[data-test-clients-config-save]');
  });
});
