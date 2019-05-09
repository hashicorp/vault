import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | icon', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<Icon class="i-con" />`);
    assert.dom('.i-con').exists('renders');

    await render(hbs`<Icon @glyph="vault-logo" />`);
    assert.dom('.vault-logo').exists('inlines the SVG');

    await render(hbs`<Icon class="ah" aria-hidden="true" />`);
    assert.dom('.ah').hasAttribute('aria-hidden', 'true', 'renders aria-hidden');

    await render(hbs`<Icon class="al" aria-label="Testing" />`);
    assert.dom('.al').hasAttribute('aria-label', 'Testing', 'renders aria-label');
  });
});
