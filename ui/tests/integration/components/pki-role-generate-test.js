import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import Sinon from 'sinon';

const SELECTORS = {
  form: '[data-test-pki-generate-cert-form]',
  commonNameField: '[data-test-input="commonName"]',
  optionsToggle: '[data-test-toggle-group="Options"]',
  generateButton: '[data-test-pki-generate-button]',
  cancelButton: '[data-test-pki-generate-cancel]',
  downloadButton: '[data-test-pki-cert-download-button]',
  revokeButton: '[data-test-pki-cert-revoke-button]',
  serialNumber: '[data-test-value-div="Serial number"]',
  certificate: '[data-test-value-div="Certificate"]',
};

module('Integration | Component | pki-role-generate', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/certificate/generate', { role: 'my-role' });
  });

  test('it should render the component with the form by default', async function (assert) {
    assert.expect(4);
    this.set('onSuccess', Sinon.spy());
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiRoleGenerate
          @model={{this.model}}
          @onSuccess={{this.onSuccess}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.form).exists('shows the cert generate form');
    assert.dom(SELECTORS.commonNameField).exists('shows the common name field');
    assert.dom(SELECTORS.optionsToggle).exists('toggle exists');
    await fillIn(SELECTORS.commonNameField, 'example.com');
    assert.strictEqual(this.model.commonName, 'example.com', 'Filling in the form updates the model');
  });

  test('it should render the component displaying the cert', async function (assert) {
    assert.expect(5);
    const record = this.store.createRecord('pki/certificate/generate', {
      role: 'my-role',
      serialNumber: 'abcd-efgh-ijkl',
    });
    this.set('model', record);
    this.set('onSuccess', Sinon.spy());
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiRoleGenerate
          @model={{this.model}}
          @onSuccess={{this.onSuccess}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.form).doesNotExist('Does not show the form');
    assert.dom(SELECTORS.downloadButton).exists('shows the download button');
    assert.dom(SELECTORS.revokeButton).exists('shows the revoke button');
    assert.dom(SELECTORS.certificate).exists({ count: 1 }, 'shows certificate info row');
    assert.dom(SELECTORS.serialNumber).hasText('abcd-efgh-ijkl', 'shows serial number info row');
  });
});
