import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import sinon from 'sinon';

const selectors = {
  form: '[data-test-pki-config-import-form]',
  saveButton: '[data-test-pki-config-save]',
  cancelButton: '[data-test-pki-config-cancel]',
};

module('Integration | Component | pki-config/import', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.config = this.store.createRecord('pki/config', {});
  });

  test('it renders', async function (assert) {
    const saveSpy = sinon.spy();
    const cancelSpy = sinon.spy();
    this.set('onSave', saveSpy);
    this.set('onCancel', cancelSpy);
    await render(
      hbs`<PkiConfig::Import @config={{this.config}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );

    assert.dom(selectors.form).exists({ count: 1 }, 'Import form exists');
    assert.dom(selectors.saveButton).isNotDisabled('Save button not disabled');
    assert.dom(selectors.cancelButton).exists({ count: 1 }, 'cancel button exists');
    await click(selectors.cancelButton);
    assert.ok(cancelSpy.calledOnce, 'cancel called on button click');
    assert.ok(saveSpy.notCalled, 'save not called when cancel clicked');
    await click(selectors.saveButton);
    assert.ok(saveSpy.calledOnce);
  });
});
