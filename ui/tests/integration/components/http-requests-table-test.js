import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const COUNTERS = [
  {
    start_time: '2019-05-01T00:00:00Z',
    total: 50,
  },
  {
    start_time: '2019-04-01T00:00:00Z',
    total: 45,
  },
  {
    start_time: '2019-03-01T00:00:00Z',
    total: 55,
  },
];

module('Integration | Component | http-requests-table', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('counters', COUNTERS);
  });

  test('it renders', async function(assert) {
    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);

    assert.ok(this.element.textContent.trim());
  });

  test('it does not show Change column with less than one month of data', async function(assert) {
    const one_month_counter = [
      {
        start_time: '2019-05-01T00:00:00Z',
        total: 50,
      },
    ];
    await render(hbs`<HttpRequestsTable @counters={{one_month_counter}}/>`);

    assert.notOk(this.element.textContent.includes('Change'));
  });

  test('it shows Change column for more than one month of data', async function(assert) {
    await render(hbs`<HttpRequestsTable @counters={{counters}}/>`);

    assert.ok(this.element.textContent.includes('Change'));
  });
});
