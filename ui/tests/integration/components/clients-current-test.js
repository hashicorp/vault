import { module, test } from 'qunit';
import EmberObject from '@ember/object';
import { render } from '@ember/test-helpers';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | client count current', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    let model = EmberObject.create({
      config: {},
      monthly: {},
      versionHistory: [],
    });
    this.model = model;
  });

  test('it shows empty state when disabled and no data available', async function (assert) {
    Object.assign(this.model.config, { enabled: 'Off', queriesAvailable: false });
    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::Dashboard @model={{this.model}} />
    <Clients::Current @model={{this.model}} />
    `);
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('Tracking is disabled');
  });

  test('it shows empty state when enabled and no data', async function (assert) {
    Object.assign(this.model.config, { enabled: 'On', queriesAvailable: false });
    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::Current @model={{this.model}} />`);
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('No partial history');
  });

  test('it shows zeroed data when enabled but no counts', async function (assert) {
    Object.assign(this.model.config, { queriesAvailable: true, enabled: 'On' });
    Object.assign(this.model.monthly, {
      total: { clients: 0, entity_clients: 0, non_entity_clients: 0 },
    });
    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::Current @model={{this.model}} />
    `);
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Empty state does not exist');
    assert.dom('[data-test-usage-stats]').exists('Client count data exists');
    assert.dom('[data-test-stat-text-container]').includesText('0');
  });

  test('it shows data when available from query', async function (assert) {
    Object.assign(this.model.config, { queriesAvailable: true, configPath: { canRead: true } });
    Object.assign(this.model.monthly, {
      total: {
        clients: 1234,
        entity_clients: 234,
        non_entity_clients: 232,
      },
    });

    await render(hbs`
    <div id="modal-wormhole"></div>
    <Clients::Current @model={{this.model}} />`);
    assert.dom('[data-test-pricing-metrics-form]').doesNotExist('Date range component should not exists');
    assert.dom('[data-test-tracking-disabled]').doesNotExist('Flash message does not exists');
    assert.dom('[data-test-usage-stats]').exists('Client count data exists');
  });
});
