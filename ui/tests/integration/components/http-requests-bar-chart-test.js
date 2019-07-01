import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const COUNTERS = [
  { start_time: '2019-04-01T00:00:00Z', total: 5500 },
  { start_time: '2019-05-01T00:00:00Z', total: 4500 },
  { start_time: '2019-06-01T00:00:00Z', total: 5000 },
];

module('Integration | Component | http-requests-bar-chart', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('counters', COUNTERS);
  });

  test('it renders', async function(assert) {
    await render(hbs`<HttpRequestsBarChart @counters={{counters}}/>`);

    assert.dom('.http-requests-bar-chart').exists();
  });

  test('it renders the correct number of bars, ticks, and gridlines', async function(assert) {
    await render(hbs`<HttpRequestsBarChart @counters={{counters}}/>`);

    assert.equal(this.element.querySelectorAll('.bar').length, 6, 'it renders the bars and shadow bars');
    assert.equal(this.element.querySelectorAll('.tick').length, 9), 'it renders the ticks and gridlines';
  });

  test('it formats the ticks', async function(assert) {
    await render(hbs`<HttpRequestsBarChart @counters={{counters}}/>`);

    assert.equal(
      this.element.querySelector('.x-axis>.tick').textContent,
      'Apr 2019',
      'x axis ticks should should show the month and year'
    );
    assert.equal(
      this.element.querySelectorAll('.y-axis>.tick')[1].textContent,
      '2k',
      'y axis ticks should round to the nearest thousand'
    );
  });

  test('it renders a tooltip', async function(assert) {
    await render(hbs`<HttpRequestsBarChart @counters={{counters}}/>`);
    await triggerEvent('.shadow-bars>.bar', 'mouseenter');
    const tooltipLabel = document.querySelector('.d3-tooltip .date');

    assert.equal(tooltipLabel.textContent, 'April 2019', 'it shows the tooltip with the formatted date');
  });
});
