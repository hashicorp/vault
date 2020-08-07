import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MODEL = {
  totalEntities: 0,
  httpsRequests: [{ start_time: '2019-04-01T00:00:00Z', total: 5500 }],
  totalTokens: 1,
};

const MODEL_WITH_GRID = {
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
    this.set('modelWithGrid', MODEL_WITH_GRID);
  });

  test('it renders', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{model}}/>`);
    assert.dom('.selectable-card-container').exists();
  });

  test('it renders a card for each of the models and titles are returned', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{model}}/>`);
    assert.dom('.selectable-card').exists({ count: 3 });
    let cardTitles = ['Http Requests', 'Entities', 'Token'];
    let httpRequestsTitle = this.element.querySelectorAll('[data-test-selectable-card-title]');

    httpRequestsTitle.forEach(item => {
      assert.notEqual(cardTitles.indexOf(item.innerText), -1);
    });
  });

  test('it renders with more than one month of data', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{modelWithGrid}}/>`);
    assert.dom('.selectable-card-container.has-grid').exists();
  });

  test('it renders 3 selectable cards when there is more than one month of data', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{modelWithGrid}}/>`);
    assert.dom('.selectable-card').exists({ count: 3 });
  });

  test('it only renders a bar chart with the last 12 months of data', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{modelWithGrid}}/>`);
    assert.dom('rect').exists({ count: 12 });
  });
});
