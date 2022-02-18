import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/vertical-bar-chart', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    this.set('chartLegend', [
      { label: 'entity clients', key: 'entity_clients' },
      { label: 'non-entity clients', key: 'non_entity_clients' },
    ]);
  });

  test('it renders', async function (assert) {
    const barChartData = [
      { month: 'january', clients: 200, entity_clients: 91, non_entity_clients: 50, new_clients: 5 },
      { month: 'february', clients: 300, entity_clients: 101, non_entity_clients: 150, new_clients: 5 },
    ];
    this.set('barChartData', barChartData);

    await render(hbs`   
      <Clients::VerticalBarChart 
        @dataset={{barChartData}} 
        @chartLegend={{chartLegend}} 
      />
    `);
    assert.dom('[data-test-vertical-bar-chart]').exists();
  });
});
