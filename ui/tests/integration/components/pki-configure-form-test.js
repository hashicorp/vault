import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | pki-configure-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  test('it renders', async function (assert) {
    await render(hbs`<PkiConfigureForm />`, { owner: this.engine });

    assert.dom('[data-test-pki-config-option]').exists({ count: 3 }, 'Three configuration options are shown');

    await click('[data-test-pki-config-option="import"]');
    assert.dom('[data-test-pki-config-option="import"]').isChecked('Selected item is checked');

    await click('[data-test-pki-config-option="generate-csr"]');
    assert.dom('[data-test-pki-config-option="generate-csr"]').isChecked('Selected item is checked');
  });
});
