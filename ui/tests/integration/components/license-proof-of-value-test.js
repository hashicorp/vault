import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | license-proof-of-value', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<LicenseProofOfValue />`);

    assert.dom().hasText('');

    // Template block usage:
    await render(hbs`
      <LicenseProofOfValue>
        template block text
      </LicenseProofOfValue>
    `);

    assert.dom().hasText('template block text');
  });
});
