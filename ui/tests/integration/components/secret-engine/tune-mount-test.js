import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | secret-engine/tune-mount', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<SecretEngine::TuneMount />`);

    assert.dom().hasText('');

    // Template block usage:
    await render(hbs`
      <SecretEngine::TuneMount>
        template block text
      </SecretEngine::TuneMount>
    `);

    assert.dom().hasText('template block text');
  });
});
