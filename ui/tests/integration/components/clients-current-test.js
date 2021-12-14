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
      activity: {},
    });
    this.model = model;
    this.tab = 'current';
  });

  test('it shows empty state when disabled and no data available', async function (assert) {
    Object.assign(this.model.config, { enabled: 'Off', queriesAvailable: false });
    await render(hbs`<Clients::History @tab={{this.tab}} @model={{this.model}} />`);

    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert.dom('[data-test-empty-state-title]').hasText('Tracking is disabled');
  });

  test('it shows zeroes when enabled and no data', async function (assert) {
    Object.assign(this.model.config, { enabled: 'On', queriesAvailable: false });
    Object.assign(this.model.activity, {
      clients: 0,
      distinct_entities: 0,
      non_entity_tokens: 0,
    });
    await render(hbs`<Clients::History @tab={{this.tab}} @model={{this.model}} />`);
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Empty state does not exist');
    assert.dom('[data-test-client-count-stats]').exists('Client count data exists');
  });

  test('it shows zeroed data when enabled but no counts', async function (assert) {
    Object.assign(this.model.config, { queriesAvailable: true, enabled: 'On' });
    Object.assign(this.model.activity, {
      clients: 1234,
      total: 1234,
    });
    await render(hbs`<Clients::History @tab={{this.tab}} @model={{this.model}} />`);
    assert.dom('[data-test-pricing-metrics-form]').doesNotExist('Date range component should not exists');
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Empty state does not exist');
    assert.dom('[data-test-client-count-stats]').exists('Client count data exists');
    assert.dom('[data-test-stat-text-container]').includesText('0');
  });

  test('it shows data when available from query', async function (assert) {
    Object.assign(this.model.config, { queriesAvailable: true, configPath: { canRead: true } });
    Object.assign(this.model.activity, {
      clients: 1234,
      distinct_entities: 234,
      non_entity_tokens: 232,
    });

    await render(hbs`<Clients::History @tab={{this.tab}} @model={{this.model}} />`);
    assert.dom('[data-test-pricing-metrics-form]').doesNotExist('Date range component should not exists');
    assert.dom('[data-test-tracking-disabled]').doesNotExist('Flash message does not exists');
    assert.dom('[data-test-client-count-stats]').exists('Client count data exists');
  });
});
