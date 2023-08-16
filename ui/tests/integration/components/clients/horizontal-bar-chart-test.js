import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { findAll, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/horizontal-bar-chart', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    this.set('chartLegend', [
      { label: 'entity clients', key: 'entity_clients' },
      { label: 'non-entity clients', key: 'non_entity_clients' },
    ]);
  });

  test('it renders chart and tooltip', async function (assert) {
    const totalObject = { clients: 5, entity_clients: 2, non_entity_clients: 3 };
    const dataArray = [
      { label: 'second', clients: 3, entity_clients: 1, non_entity_clients: 2 },
      { label: 'first', clients: 2, entity_clients: 1, non_entity_clients: 1 },
    ];
    this.set('totalUsageCounts', totalObject);
    this.set('totalClientsData', dataArray);

    await render(hbs`
    <Clients::HorizontalBarChart
      @dataset={{this.totalClientsData}}
      @chartLegend={{chartLegend}}
      @totalUsageCounts={{totalUsageCounts}}
    />`);

    assert.dom('[data-test-horizontal-bar-chart]').exists();
    const dataBars = findAll('[data-test-horizontal-bar-chart] rect.data-bar');
    const actionBars = findAll('[data-test-horizontal-bar-chart] rect.action-bar');

    assert.equal(actionBars.length, dataArray.length, 'renders correct number of hover bars');
    assert.equal(dataBars.length, dataArray.length * 2, 'renders correct number of data bars');

    const textLabels = this.element.querySelectorAll('[data-test-horizontal-bar-chart] .tick text');
    const textTotals = this.element.querySelectorAll('[data-test-horizontal-bar-chart] text.total-value');
    textLabels.forEach((label, index) => {
      assert.dom(label).hasText(dataArray[index].label, 'label renders correct text');
    });
    textTotals.forEach((label, index) => {
      assert.dom(label).hasText(`${dataArray[index].clients}`, 'total value renders correct number');
    });

    for (let [i, bar] of actionBars.entries()) {
      let percent = Math.round((dataArray[i].clients / totalObject.clients) * 100);
      await triggerEvent(bar, 'mouseover');
      let tooltip = document.querySelector('.ember-modal-dialog');
      assert.dom(tooltip).includesText(`${percent}%`, 'tooltip renders correct percentage');
    }
  });

  test('it renders data with a large range', async function (assert) {
    const totalObject = { clients: 5929393, entity_clients: 1391997, non_entity_clients: 4537396 };
    const dataArray = [
      { label: 'second', clients: 5929093, entity_clients: 1391896, non_entity_clients: 4537100 },
      { label: 'first', clients: 300, entity_clients: 101, non_entity_clients: 296 },
    ];
    this.set('totalUsageCounts', totalObject);
    this.set('totalClientsData', dataArray);

    await render(hbs`
    <Clients::HorizontalBarChart
      @dataset={{this.totalClientsData}}
      @chartLegend={{chartLegend}}
      @totalUsageCounts={{totalUsageCounts}}
    />`);

    assert.dom('[data-test-horizontal-bar-chart]').exists();
    const dataBars = findAll('[data-test-horizontal-bar-chart] rect.data-bar');
    const actionBars = findAll('[data-test-horizontal-bar-chart] rect.action-bar');

    assert.equal(actionBars.length, dataArray.length, 'renders correct number of hover bars');
    assert.equal(dataBars.length, dataArray.length * 2, 'renders correct number of data bars');

    for (let [i, bar] of actionBars.entries()) {
      let percent = Math.round((dataArray[i].clients / totalObject.clients) * 100);
      await triggerEvent(bar, 'mouseover');
      let tooltip = document.querySelector('.ember-modal-dialog');
      assert.dom(tooltip).includesText(`${percent}%`, 'tooltip renders correct percentage');
    }
  });
});
