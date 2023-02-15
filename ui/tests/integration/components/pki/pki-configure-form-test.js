import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-configure-form';
import sinon from 'sinon';

module('Integration | Component | pki-configure-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  hooks.beforeEach(function () {
    this.cancelSpy = sinon.spy();
  });

  test('it renders', async function (assert) {
    await render(hbs`<PkiConfigureForm @onCancel={{this.cancelSpy}} @config={{this.config}} />`, {
      owner: this.engine,
    });

    assert.dom(SELECTORS.option).exists({ count: 3 }, 'Three configuration options are shown');
    assert.dom(SELECTORS.cancelButton).exists('Cancel link is shown');
    assert.dom(SELECTORS.saveButton).isDisabled('Done button is disabled');

    await click(SELECTORS.optionByKey('import'));
    assert.dom(SELECTORS.optionByKey('import')).isChecked('Selected item is checked');

    await click(SELECTORS.optionByKey('generate-csr'));
    assert.dom(SELECTORS.optionByKey('generate-csr')).isChecked('Selected item is checked');
  });
});
