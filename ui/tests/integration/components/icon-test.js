import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import waitForError from 'vault/tests/helpers/wait-for-error';

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

    await render(hbs`<Icon @glyph="vault-logo" @size="s"/>`);
    assert.dom('.hs-icon').hasClass('hs-icon-s', 'adds the size class');

    let promise = waitForError();
    render(hbs`<Icon @glyph="vault-logo" @size="no"/>`);
    let err = await promise;
    assert.ok(err.message.includes('The size property of'), "errors when passed a size that's not allowed");
  });
});
