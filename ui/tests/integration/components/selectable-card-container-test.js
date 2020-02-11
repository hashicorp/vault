import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MODEL = {
  totalEntities: 0,
  httpsRequests: [{ start_time: '2019-04-01T00:00:00Z', total: 5500 }],
  totalTokens: 1,
};

const MODEL_GRAPH = {
  httpsRequests: [
    { start_time: '2018-11-01T00:00:00Z', total: 5500 },
    { start_time: '2018-12-01T00:00:00Z', total: 4500 },
    { start_time: '2019-01-01T00:00:00Z', total: 5000 },
    { start_time: '2019-02-01T00:00:00Z', total: 5000 },
    { start_time: '2019-03-01T00:00:00Z', total: 5000 },
    { start_time: '2019-04-01T00:00:00Z', total: 5500 },
    { start_time: '2019-05-01T00:00:00Z', total: 4500 },
    { start_time: '2019-06-01T00:00:00Z', total: 5000 },
    { start_time: '2019-07-01T00:00:00Z', total: 5000 },
    { start_time: '2019-08-01T00:00:00Z', total: 5000 },
    { start_time: '2019-09-01T00:00:00Z', total: 5000 },
    { start_time: '2019-10-01T00:00:00Z', total: 5000 },
    { start_time: '2019-11-01T00:00:00Z', total: 5000 },
    { start_time: '2019-12-01T00:00:00Z', total: 5000 },
    { start_time: '2020-01-01T00:00:00Z', total: 5000 },
    { start_time: '2020-02-01T00:00:00Z', total: 5000 },
  ],
  totalEntities: 0,
  totalTokens: 1,
};

module('Integration | Component | selectable-card-container', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('model', MODEL);
    this.set('modelGraph', MODEL_GRAPH);
  });

  test('it renders', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{model}}/>`);
    assert.dom('.selectable-card-container').exists();
  });

  test('it renders a card for each of the models and titles are returned', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{model}}/>`);
    assert.dom('.selectable-card').exists({ count: 3 });

    assert.dom(`[data-test-selectable-card-title=Requests]`).exists();
    assert.dom(`[data-test-selectable-card-title=Entities]`).exists();
    assert.dom(`[data-test-selectable-card-title=Tokens]`).exists();
  });

  test('it renders with more than one month of data', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{modelGraph}}/>`);
    assert.dom('.selectable-card-container-graph').exists();
  });

  test('it renders 3 selectable cards when there is more than one month of data', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{modelGraph}}/>`);
    assert.dom('.selectable-card').exists({ count: 3 });
  });

  test('it only renders a bar chart with the last 12 months of data', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{modelGraph}}/>`);
    assert.dom('rect').exists({ count: 12 });
  });
});
