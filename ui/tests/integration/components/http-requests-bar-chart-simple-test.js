import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const FILTERED_HTTPS_REQUESTS = [
  { start_time: '2018-11-01T00:00:00Z', total: 5500 },
  { start_time: '2018-12-01T00:00:00Z', total: 4500 },
  { start_time: '2019-01-01T00:00:00Z', total: 5000 },
  { start_time: '2019-02-01T00:00:00Z', total: 5000 },
];

module('Integration | Component | http-requests-bar-chart-simple', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('filteredHttpsRequests', FILTERED_HTTPS_REQUESTS);
  });

  test('it renders and the correct number of bars are showing', async function(assert) {
    await render(hbs`<HttpRequestsBarChartSimple @counters={{filteredHttpsRequests}}/>`);

    assert.dom('rect').exists({ count: FILTERED_HTTPS_REQUESTS.length });
    assert.dom('.http-requests-bar-chart-simple').exists();
  });
});
