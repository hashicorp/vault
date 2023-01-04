import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';

module('Integration | Component | pki-config/import', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  test('it renders', async function (assert) {
    await render(hbs`<PkiConfig::Import />`, { owner: this.engine });

    assert.dom('[data-test-pki-config-import-form]').exists({ count: 1 }, 'Import form exists');
    assert.dom('[data-test-pki-config-save]').isNotDisabled('Save button not disabled');
    assert.dom('[data-test-pki-config-cancel]').exists({ count: 1 }, 'cancel button exists');
  });
});
