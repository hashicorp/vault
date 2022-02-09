import { module, test } from 'qunit';
import EmberObject from '@ember/object';
import { render } from '@ember/test-helpers';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | client count history', function (hooks) {
  // TODO CMB add tests for calendar widget showing
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    let model = EmberObject.create({
      config: {},
      activity: {},
      versionHistory: [],
    });
    this.model = model;
  });

  test('it shows empty state when disabled and no data available', async function (assert) {
    Object.assign(this.model.config, { enabled: 'Off', queriesAvailable: false });
    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::History @model={{this.model}} />`);

    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('Data tracking is disabled');
  });

  test('it shows empty state when enabled and no data available', async function (assert) {
    Object.assign(this.model.config, { enabled: 'On', queriesAvailable: false });
    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::History @model={{this.model}} />`);

    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No monthly history');
  });

  test('it shows empty state when no data for queried date range', async function (assert) {
    Object.assign(this.model.config, { queriesAvailable: true });
    Object.assign(this.model, { startTimeFromLicense: ['2021', 5] });
    Object.assign(this.model.activity, {
      byNamespace: [
        {
          label: 'namespace24/',
          clients: 8301,
          entity_clients: 4387,
          non_entity_clients: 3914,
          mounts: [],
        },
        {
          label: 'namespace88/',
          clients: 7752,
          entity_clients: 3632,
          non_entity_clients: 4120,
          mounts: [],
        },
      ],
    });
    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::History @model={{this.model}} />`);
    assert.dom('[data-test-start-date-editor]').exists('Billing start date editor exists');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No data received');
  });

  test('it shows warning when disabled and data available', async function (assert) {
    Object.assign(this.model.config, { queriesAvailable: true, enabled: 'Off' });
    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::History @model={{this.model}} />`);

    assert.dom('[data-test-start-date-editor]').exists('Billing start date editor exists');
    assert.dom('[data-test-tracking-disabled]').exists('Flash message exists');
    assert.dom('[data-test-tracking-disabled] .message-title').hasText('Tracking is disabled');
  });

  test('it shows data when available from query', async function (assert) {
    Object.assign(this.model.config, { queriesAvailable: true, configPath: { canRead: true } });
    Object.assign(this.model, { startTimeFromLicense: ['2021', 5] });
    Object.assign(this.model.activity, {
      byNamespace: [
        { label: 'nsTest5/', clients: 2725, entity_clients: 1137, non_entity_clients: 1588 },
        { label: 'nsTest1/', clients: 200, entity_clients: 100, non_entity_clients: 100 },
      ],
      total: {
        clients: 1234,
        entity_clients: 234,
        non_entity_clients: 232,
      },
    });

    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::History @model={{this.model}} />`);

    assert.dom('[data-test-start-date-editor]').exists('Billing start date editor exists');
    assert.dom('[data-test-tracking-disabled]').doesNotExist('Flash message does not exists');
    assert.dom('[data-test-usage-stats]').exists('Client count data exists');
    assert.dom('[data-test-horizontal-bar-chart]').exists('Horizontal bar chart exists');
  });
});
