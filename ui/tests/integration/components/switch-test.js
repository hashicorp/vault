import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find, pauseTest } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | switch', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`{{switch}}`);

    assert.equal(this.element.textContent.trim(), '');

    // Template block usage:
    await render(hbs`
      {{#switch}}
        <span id="test-value" class="has-text-grey">template block text</span>
      {{/switch}}
    `);
    await pauseTest();
    assert.dom('[data-test-switch-label]').exists('switch label exists');
    assert.equal(find('#test-value').textContent.trim(), 'template block text', 'yielded text renders');
  });
});
