import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/usage-stats', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders defaults', async function (assert) {
    await render(hbs`<Clients::UsageStats />`);

    assert.dom('[data-test-stat-text]').exists({ count: 3 }, 'Renders 3 Stat texts even with no data passed');
    assert.dom('[data-test-stat-text="total-clients"]').exists('Total clients exists');
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('0', 'Value defaults to zero');
    assert.dom('[data-test-stat-text="entity-clients"]').exists('Entity clients exists');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('0', 'Value defaults to zero');
    assert.dom('[data-test-stat-text="non-entity-clients"]').exists('Non entity clients exists');
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText('0', 'Value defaults to zero');
    assert
      .dom('a')
      .hasAttribute('href', 'https://developer.hashicorp.com/vault/tutorials/monitoring/usage-metrics');
  });

  test('it renders with data', async function (assert) {
    this.set('counts', {
      clients: 17,
      entity_clients: 7,
      non_entity_clients: 10,
    });
    await render(hbs`<Clients::UsageStats @totalUsageCounts={{this.counts}} />`);

    assert.dom('[data-test-stat-text]').exists({ count: 3 }, 'Renders 3 Stat texts even with no data passed');
    assert.dom('[data-test-stat-text="total-clients"]').exists('Total clients exists');
    assert
      .dom('[data-test-stat-text="total-clients"] .stat-value')
      .hasText('17', 'Total clients shows passed value');
    assert.dom('[data-test-stat-text="entity-clients"]').exists('Entity clients exists');
    assert
      .dom('[data-test-stat-text="entity-clients"] .stat-value')
      .hasText('7', 'entity clients shows passed value');
    assert.dom('[data-test-stat-text="non-entity-clients"]').exists('Non entity clients exists');
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText('10', 'non entity clients shows passed value');
  });
});
