import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';

const SELECTORS = {
  option: '[data-test-pki-config-option]',
  optionByKey: (key) => `[data-test-pki-config-option="${key}"]`,
  cancelButton: '[data-test-pki-config-cancel]',
  saveButton: '[data-test-pki-config-save]',
};
module('Integration | Component | pki-configure-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  test('it renders', async function (assert) {
    await render(hbs`<PkiConfigureForm @config={{this.config}} />`, { owner: this.engine });

    assert.dom(SELECTORS.option).exists({ count: 3 }, 'Three configuration options are shown');
    assert.dom(SELECTORS.cancelButton).exists('Cancel link is shown');
    assert.dom(SELECTORS.saveButton).isDisabled('Done button is disabled');

    await click(SELECTORS.optionByKey('import'));
    assert.dom(SELECTORS.optionByKey('import')).isChecked('Selected item is checked');

    await click(SELECTORS.optionByKey('generate-csr'));
    assert.dom(SELECTORS.optionByKey('generate-csr')).isChecked('Selected item is checked');
  });
});
