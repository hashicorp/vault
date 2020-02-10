import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, pauseTest } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MODEL = {
  totalEntities: 0,
  httpsRequests: [{ start_time: '2019-04-01T00:00:00Z', total: 5500 }],
  totalTokens: 1,
};

module('Integration | Component | selectable-card-container', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('model', MODEL);
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
});
