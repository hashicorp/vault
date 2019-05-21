import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import waitForError from 'vault/tests/helpers/wait-for-error';

module('Integration | Component | chevron', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<Chevron />`);
    assert.dom('.hs-icon').exists('renders');

    await render(hbs`<Chevron @isButton={{true}} />`);
    assert.dom('.hs-icon').hasClass('hs-icon-button-right', 'renders');

    await render(hbs`<Chevron @direction='left' @isButton={{true}} />`);
    assert.dom('.hs-icon').doesNotHaveClass('hs-icon-button-right', 'renders');

    let promise = waitForError();
    render(hbs`<Chevron @direction='lol' />`);
    let err = await promise;
    assert.ok(
      err.message.includes('The direction property of <vault@component:chevron:'),
      'asserts about unsupported direction'
    );
  });
});
