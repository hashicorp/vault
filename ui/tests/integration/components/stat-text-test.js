import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | StatText', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<StatText />`);

    assert.equal(this.element.textContent.trim(), '');
  });

  test('it renders passed in attributes', async function(assert) {
    this.set('label', 'A Label');
    this.set('value', '9,999');
    this.set('size', 'l');
    this.set('subText', 'This is my description');

    await render(
      hbs`<StatText @label={{this.label}} @size={{this.size}} @value={{this.value}} @subText={{this.subText}}/>`
    );

    assert.dom('.stat-label').hasText(this.label, 'renders label');
    assert.dom('.stat-text').hasText(this.subText, 'renders subtext');
    assert.dom('.stat-value').hasText(this.value, 'renders value');
  });
});
