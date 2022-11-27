import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, click } from '@ember/test-helpers';
import { resolve } from 'rsvp';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | client count config', function (hooks) {
  setupRenderingTest(hooks);

  const createAttr = (name, type, options) => {
    return {
      name,
      type,
      options,
    };
  };

  const generateModel = (overrides) => {
    return {
      enabled: 'On',
      retentionMonths: 24,
      defaultReportMonths: 12,
      configAttrs: [
        createAttr('enabled', 'string', { editType: 'boolean' }),
        createAttr('retentionMonths', 'number'),
      ],
      changedAttributes: () => ({}),
      save: () => {},
      ...overrides,
    };
  };

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.router.reopen({
      transitionTo() {
        return {
          followRedirects() {
            return resolve();
          },
        };
      },
    });
    const model = generateModel();
    this.model = model;
  });

  test('it shows the table with the correct rows by default', async function (assert) {
    await render(hbs`<Clients::Config @model={{this.model}} />`);

    assert.dom('[data-test-pricing-metrics-config-table]').exists('Pricing metrics config table exists');
    const rows = document.querySelectorAll('.info-table-row');
    assert.strictEqual(rows.length, 2, 'renders 2 infotable rows');
    assert.ok(
      find('[data-test-row-value="Usage data collection"]').textContent.includes('On'),
      'Enabled value matches model'
    );
    assert.ok(
      find('[data-test-row-value="Retention period"]').textContent.includes('24'),
      'Retention period value matches model'
    );
  });

  test('TODO: it shows the config edit form when mode = edit', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Config @model={{this.model}} @mode="edit" />
    `);

    assert.dom('[data-test-pricing-metrics-config-form]').exists('Pricing metrics config form exists');
    const fields = document.querySelectorAll('[data-test-field]');
    assert.strictEqual(fields.length, 2, 'renders 2 fields');
  });

  test('it shows a modal with correct messaging when disabling', async function (assert) {
    // Simulates the model when enabled value has been changed from On to Off
    const simModel = generateModel({
      enabled: 'Off',
      changedAttributes: () => ({ enabled: ['On', 'Off'] }),
    });
    this.set('model', simModel);
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Config @model={{this.model}} @mode="edit" />
    `);

    await click('[data-test-edit-metrics-config-save]');
    assert.dom('.modal.is-active').exists('Modal appears');
    assert.ok(
      find('[data-test-modal-title]').textContent.includes('Turn usage tracking off?'),
      'Modal confirming turn tracking off'
    );
    await click('[data-test-metrics-config-cancel]');
    assert.dom('.modal.is-active').doesNotExist('Modal goes away');
  });

  test('it shows a modal with correct messaging when enabling', async function (assert) {
    // Simulates the model when enabled value has been changed from On to Off
    const simModel = generateModel({
      changedAttributes: () => ({ enabled: ['Off', 'On'] }),
    });
    this.set('model', simModel);
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Config @model={{this.model}} @mode="edit" />
    `);

    await click('[data-test-edit-metrics-config-save]');
    assert.dom('.modal.is-active').exists('Modal appears');
    assert.ok(
      find('[data-test-modal-title]').textContent.includes('Turn usage tracking on?'),
      'Modal confirming turn tracking on'
    );
    await click('[data-test-metrics-config-cancel]');
    assert.dom('.modal.is-active').doesNotExist('Modal goes away');
  });

  test('it does not show a modal on save if enable left unchanged', async function (assert) {
    // Simulates the model when something other than enabled changed
    const simModel = generateModel({
      changedAttributes: () => ({ retentionMonths: [24, '48'] }),
    });
    this.set('model', simModel);
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Config @model={{this.model}} @mode="edit" />
    `);

    await click('[data-test-edit-metrics-config-save]');
    assert.dom('.modal.is-active').doesNotExist('No modal appears');
  });
});
