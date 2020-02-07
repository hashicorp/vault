import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
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
    // the component due to CSS issues does not use the classNames in the route file
    // thus the test is unable to select by the class name of the component
    assert.dom('.selectable-card').exists({ count: 3 });
  });
});
