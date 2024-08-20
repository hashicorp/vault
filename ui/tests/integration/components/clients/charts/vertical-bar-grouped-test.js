import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/charts/vertical-bar-grouped', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<Clients::Charts::VerticalBarGrouped />`);

    assert.dom().hasText('');

    // Template block usage:
    await render(hbs`<Clients::Charts::VerticalBarGrouped>
  template block text
</Clients::Charts::VerticalBarGrouped>`);

    assert.dom().hasText('template block text');
  });
});
