import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | /metrics/selectable-card-container-graph', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`{{/metrics/selectable-card-container-graph}}`);

    assert.equal(this.element.textContent.trim(), '');

    // Template block usage:
    await render(hbs`
      {{#/metrics/selectable-card-container-graph}}
        template block text
      {{//metrics/selectable-card-container-graph}}
    `);

    assert.equal(this.element.textContent.trim(), 'template block text');
  });
});
