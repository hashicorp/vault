import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const COUNTERS = [
  { start_time: '2018-04-01T00:00:00Z', total: 5500 },
  { start_time: '2019-05-01T00:00:00Z', total: 4500 },
  { start_time: '2019-06-01T00:00:00Z', total: 5000 },
];

module('Integration | Component | http-requests-dropdown', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('counters', COUNTERS);
  });

  test('it renders with options', async function(assert) {
    await render(hbs`<HttpRequestsDropdown @counters={{counters}} />`);

    assert.dom('[data-test-date-range]').hasValue('All', 'shows all data by default');

    assert.equal(
      this.element.querySelector('[data-test-date-range]').options.length,
      4,
      'it adds an option for each year in the data set'
    );
  });
});
