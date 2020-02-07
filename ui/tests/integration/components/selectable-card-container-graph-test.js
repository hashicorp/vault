import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MODEL = {
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

module('Integration | Component | selectable-card-container-graph', function(hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function() {
    this.set('model', MODEL);
  });

  test('it renders', async function(assert) {
    await render(hbs`<SelectableCardContainer @counters={{model}}/>`);
    // the component due to CSS issues does not use the classNames in the route file
    // thus the test is unable to select by the class name of the component
    assert.dom('.selectable-card').exists({ count: 3 });
  });

  test('only 12 months of data are sent in filteredHttpsRequests', async function(assert) {
    await render(hbs`<SelectableCardContainerGraph @counters={{model}}/>`);
    assert.dom('rect').exists({ count: 12 });
  });
});
