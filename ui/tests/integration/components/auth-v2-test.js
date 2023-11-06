import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | auth-v2', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<AuthV2 />`);

    assert.dom(this.element).hasText('');

    // Template block usage:
    await render(hbs`
      <AuthV2>
        template block text
      </AuthV2>
    `);

    assert.dom(this.element).hasText('template block text');
  });
});
