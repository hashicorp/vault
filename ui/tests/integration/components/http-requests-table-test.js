import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const COUNTERS = [
  { start_time: '2019-04-01T00:00:00Z', total: 5500 },
  { start_time: '2019-05-01T00:00:00Z', total: 4500 },
  { start_time: '2019-06-01T00:00:00Z', total: 5000 },
];

module('Integration | Component | http-requests-table', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('counters', COUNTERS);
  });

  test('it renders', async function(assert) {
    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);

    assert.dom('.http-requests-table').exists();
  });

  test('it does not show Change column with less than one month of data', async function(assert) {
    const one_month_counter = [
      {
        start_time: '2019-05-01T00:00:00Z',
        total: 50,
      },
    ];
    this.set('one_month_counter', one_month_counter);

    await render(hbs`<HttpRequestsTable @counters={{one_month_counter}}/>`);

    assert.dom('.http-requests-table').exists();
    assert.dom('[data-test-change]').doesNotExist();
  });

  test('it shows Change column for more than one month of data', async function(assert) {
    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);

    assert.dom('[data-test-change]').exists();
  });

  test('it shows the percent change between each time window', async function(assert) {
    const simple_counters = [
      { start_time: '2019-04-01T00:00:00Z', total: 1 },
      { start_time: '2019-05-01T00:00:00Z', total: 2 },
      { start_time: '2019-06-01T00:00:00Z', total: 1 },
      { start_time: '2019-07-01T00:00:00Z', total: 1 },
    ];
    this.set('counters', simple_counters);

    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);
    // the expectedValues are in reverse chronological order because that is the order
    // that the table shows its data.
    let expectedValues = ['', '-50%', '100%', ''];

    this.element.querySelectorAll('[data-test-change]').forEach((td, i) => {
      return assert.equal(td.textContent.trim(), expectedValues[i]);
    });
  });
});
