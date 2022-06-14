import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/line-chart', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('dataset', [
      {
        foo: 1,
        bar: 4,
      },
      {
        foo: 2,
        bar: 8,
      },
      {
        foo: 3,
        bar: 14,
      },
      {
        foo: 4,
        bar: 10,
      },
    ]);
  });

  test('it renders', async function (assert) {
    await render(hbs`
    <div class="chart-container-wide">
      <Clients::LineChart @dataset={{dataset}} @xKey="foo" @yKey="bar" />
      </div>
    `);

    assert.dom('[data-test-line-chart]').exists('Chart is rendered');
    assert.dom('.hover-circle').exists({ count: 4 }, 'Renders dot for each data point');
  });
});
