/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | pki-generate-csr', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.owner.lookup('service:secretMountPath').update('pki-test');
    this.store = this.owner.lookup('service:store');
    this.onComplete = () => {};
    this.model = this.owner
      .lookup('service:store')
      .createRecord('pki/action', { actionType: 'generate-csr' });

    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
        'pki-test/issuers/generate/intermediate/exported': ['root'],
      },
    }));
    setRunOptions({
      rules: {
        // something strange happening here
        'link-name': { enabled: false },
      },
    });
  });

  test('it should render fields and save', async function (assert) {
    assert.expect(9);

    this.server.post('/pki-test/issuers/generate/intermediate/exported', (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.strictEqual(payload.common_name, 'foo', 'Request made to correct endpoint on save');
      return {
        request_id: '123',
        data: {},
      };
    });

    await render(hbs`<PkiGenerateCsr @model={{this.model}} @onComplete={{this.onComplete}} />`, {
      owner: this.engine,
    });

    const fields = [
      'type',
      'commonName',
      'excludeCnFromSans',
      'format',
      'subjectSerialNumber',
      'addBasicConstraints',
    ];
    fields.forEach((key) => {
      assert.dom(`[data-test-input="${key}"]`).exists(`${key} form field renders`);
    });

    assert.dom('[data-test-toggle-group]').exists({ count: 3 }, 'Toggle groups render');

    await fillIn('[data-test-input="type"]', 'exported');
    await fillIn('[data-test-input="commonName"]', 'foo');
    await click('[data-test-save]');

    const savedRecord = this.store.peekAll('pki/action').firstObject;
    assert.false(savedRecord.isNew, 'record is saved');
  });

  test('it should display validation errors', async function (assert) {
    assert.expect(4);

    this.onCancel = () => assert.ok(true, 'onCancel action fires');

    await render(
      hbs`<PkiGenerateCsr @model={{this.model}} @onCancel={{this.onCancel}} @onComplete={{this.onComplete}} />`,
      {
        owner: this.engine,
      }
    );

    await click('[data-test-save]');

    assert
      .dom('[data-test-field-validation="type"]')
      .hasText('Type is required.', 'Type validation error renders');
    assert
      .dom('[data-test-field="commonName"] [data-test-inline-alert]')
      .hasText('Common name is required.', 'Common name validation error renders');
    assert.dom('[data-test-alert]').hasText('There are 2 errors with this form.', 'Alert renders');

    await click('[data-test-cancel]');
  });

  test('it should show generated CSR for type=exported', async function (assert) {
    assert.expect(6);
    this.model.id = '1235-someId';
    this.model.csr = '-----BEGIN CERTIFICATE REQUEST-----...-----END CERTIFICATE REQUEST-----';
    this.model.keyId = '9179de78-1275-a1cf-ebb0-a4eb2e376636';
    this.model.privateKey = '-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----';
    this.model.privateKeyType = 'rsa';
    this.onComplete = () => assert.ok(true, 'onComplete action fires');

    await render(hbs`<PkiGenerateCsr @model={{this.model}} @onComplete={{this.onComplete}} />`, {
      owner: this.engine,
    });
    assert
      .dom('[data-test-next-steps-csr]')
      .hasText(
        'Next steps Copy the CSR below for a parent issuer to sign and then import the signed certificate back into this mount. The private_key is only available once. Make sure you copy and save it now.',
        'renders Next steps alert banner'
      );

    assert
      .dom('[data-test-value-div="CSR"] [data-test-certificate-card] button')
      .hasAttribute('data-test-copy-button', this.model.csr, 'it renders copyable csr');
    assert
      .dom('[data-test-value-div="Key ID"] button')
      .hasAttribute('data-test-copy-button', this.model.keyId, 'it renders copyable key_id');
    assert
      .dom('[data-test-value-div="Private key"] [data-test-certificate-card] button')
      .hasAttribute('data-test-copy-button', this.model.privateKey, 'it renders copyable private_key');
    assert
      .dom('[data-test-value-div="Private key type"]')
      .hasText(this.model.privateKeyType, 'renders private_key_type');
    await click('[data-test-done]');
  });

  test('it should show generated CSR for type=internal', async function (assert) {
    assert.expect(5);
    this.model.id = '1235-someId';
    this.model.csr = '-----BEGIN CERTIFICATE REQUEST-----...-----END CERTIFICATE REQUEST-----';
    this.model.keyId = '9179de78-1275-a1cf-ebb0-a4eb2e376636';
    this.onComplete = () => {};

    await render(hbs`<PkiGenerateCsr @model={{this.model}} @onComplete={{this.onComplete}} />`, {
      owner: this.engine,
    });
    assert
      .dom('[data-test-next-steps-csr]')
      .hasText(
        'Next steps Copy the CSR below for a parent issuer to sign and then import the signed certificate back into this mount.',
        'renders Next steps alert banner'
      );
    assert
      .dom('[data-test-value-div="CSR"] [data-test-certificate-card] button')
      .hasAttribute('data-test-copy-button', this.model.csr, 'it renders copyable csr');
    assert
      .dom('[data-test-value-div="Key ID"] button')
      .hasAttribute('data-test-copy-button', this.model.keyId, 'it renders copyable key_id');
    assert.dom('[data-test-value-div="Private key"]').hasText('internal', 'does not render private key');
    assert
      .dom('[data-test-value-div="Private key type"]')
      .hasText('internal', 'does not render private key type');
  });
});
