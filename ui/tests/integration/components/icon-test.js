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

    await render(hbs`<Icon @name="vault-logo" />`);
    assert.dom('.vault-logo').exists('inlines the SVG');

    await render(hbs`<Icon class="ah" aria-hidden="true" />`);
    assert.dom('.ah').hasAttribute('aria-hidden', 'true', 'renders aria-hidden');

    await render(hbs`<Icon class="al" aria-label="Testing" />`);
    assert.dom('.al').hasAttribute('aria-label', 'Testing', 'renders aria-label');

    await render(hbs`<Icon @name="vault-logo" @sizeClass="s"/>`);
    assert.dom('.hs-icon').hasClass('hs-icon-s', 'adds the size class');

    let promise = waitForError();
    render(hbs`<Icon @name="vault-logo" @sizeClass="no"/>`);
    let err = await promise;
    assert.ok(
      err.message.includes('The sizeClass property of'),
      "errors when passed a sizeClass that's not allowed"
    );
  });

  test('it should render FlightIcon', async function(assert) {
    assert.expect(5);

    await render(hbs`<Icon @name="x" @sizeClass="xl" />`);
    assert.dom('.flight-icon').exists('FlightIcon renders when provided name of icon in set');
    assert.dom('.flight-icon').hasClass('hs-icon-xl', 'hs icon class applied to component');
    assert.dom('.flight-icon').hasAttribute('width', '24', 'Correct size applied based on sizeClass');

    await render(hbs`<Icon @name="x" @size="24" />`);
    assert.dom('.flight-icon').hasAttribute('width', '24', 'Size applied to svg');

    const promise = waitForError();
    render(hbs`<Icon @name="x" @size="12"/>`);
    const err = await promise;
    assert.ok(
      err.message.includes(`must be either '16' or '24'`),
      "errors when passed a size that's not allowed"
    );
  });
});
