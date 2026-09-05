/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, find, fillIn, blur, triggerEvent, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { encodeString } from 'vault/utils/b64';
import waitForError from 'vault/tests/helpers/wait-for-error';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | transit key actions', function (hooks) {
  setupRenderingTest(hooks);

  test('it requires `key`', async function (assert) {
    const promise = waitForError();
    render(hbs`
      <TransitKeyActions />`);
    const err = await promise;
    assert.ok(err.message.includes('@key is required for'), 'asserts without key');
  });

  test('it renders', async function (assert) {
    this.set('key', { backend: 'transit', supportedActions: ['encrypt'] });
    await render(hbs`<TransitKeyActions @selectedAction="encrypt" @key={{this.key}} />`);
    assert.dom('[data-test-transit-action="encrypt"]').exists({ count: 1 }, 'renders encrypt');

    this.set('key', { backend: 'transit', supportedActions: ['sign'] });
    await render(hbs`<TransitKeyActions @selectedAction="sign" @key={{this.key}} />`);
    assert.dom('[data-test-transit-action="sign"]').exists({ count: 1 }, 'renders sign');
  });

  test('it renders: signature_algorithm field', async function (assert) {
    this.set('key', { backend: 'transit', supportsSigning: true, supportedActions: ['sign', 'verify'] });
    this.set('selectedAction', 'sign');
    await render(hbs`
    <TransitKeyActions @selectedAction={{this.selectedAction}} @key={{this.key}} />`);
    assert
      .dom('[data-test-signature-algorithm]')
      .doesNotExist('does not render signature_algorithm field on sign');
    this.set('selectedAction', 'verify');
    assert
      .dom('[data-test-signature-algorithm]')
      .doesNotExist('does not render signature_algorithm field on verify');

    this.set('selectedAction', 'sign');
    this.set('key', {
      type: 'rsa-2048',
      supportsSigning: true,
      backend: 'transit',
      supportedActions: ['sign', 'verify'],
    });
    assert
      .dom('[data-test-signature-algorithm]')
      .exists({ count: 1 }, 'renders signature_algorithm field on sign with rsa key');
    this.set('selectedAction', 'verify');
    assert
      .dom('[data-test-signature-algorithm]')
      .exists({ count: 1 }, 'renders signature_algorithm field on verify with rsa key');
  });

  test('it renders: padding_scheme field for rsa key types', async function (assert) {
    const supportedActions = ['datakey', 'decrypt', 'encrypt'];
    const supportedKeyTypes = ['rsa-2048', 'rsa-3072', 'rsa-4096'];

    for (const key of supportedKeyTypes) {
      this.set('key', {
        type: key,
        backend: 'transit',
        supportedActions,
      });
      for (const action of this.key.supportedActions) {
        this.selectedAction = action;
        await render(hbs`
    <TransitKeyActions @selectedAction={{this.selectedAction}} @key={{this.key}} />`);
        assert
          .dom('[data-test-padding-scheme]')
          .hasValue(
            'oaep',
            `key type: ${key} renders padding_scheme field with default value for action: ${action}`
          );
      }
    }
  });
  test('it renders: decrypt_padding_scheme and encrypt_padding_scheme fields for rsa key types', async function (assert) {
    this.selectedAction = 'rewrap';
    const supportedKeyTypes = ['rsa-2048', 'rsa-3072', 'rsa-4096'];
    const SELECTOR = (type) => `[data-test-padding-scheme="${type}"]`;
    for (const key of supportedKeyTypes) {
      this.set('key', {
        type: key,
        backend: 'transit',
        supportedActions: [this.selectedAction],
      });
      await render(hbs`
    <TransitKeyActions @selectedAction={{this.selectedAction}} @key={{this.key}} />`);
      assert
        .dom(SELECTOR('encrypt'))
        .hasValue('oaep', `key type: ${key} renders ${SELECTOR('encrypt')} field with default value`);
      assert
        .dom(SELECTOR('decrypt'))
        .hasValue('oaep', `key type: ${key} renders ${SELECTOR('decrypt')} field with default value`);
    }
  });

  async function doEncrypt(assert, actions = [], keyattrs = {}) {
    const keyDefaults = {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['encrypt'].concat(actions),
    };

    const key = { ...keyDefaults, ...keyattrs };
    this.set('key', key);
    this.set('selectedAction', 'encrypt');

    this.apiStub = sinon.stub(this.owner.lookup('service:api').secrets, 'transitEncrypt').resolves({
      data: { ciphertext: 'secret' },
    });

    await render(hbs`
    <TransitKeyActions
      @selectedAction={{this.selectedAction}}
      @key={{this.key}}
    />
  `);

    let editor;
    await waitFor('.cm-editor');
    editor = codemirror('#plaintext-control');
    setCodeEditorValue(editor, 'plaintext');

    await click('button[type="submit"]');

    assert.true(this.apiStub.calledOnce, 'calls the API to encrypt');

    assert.true(
      this.apiStub.calledWith('akey', 'transit', {
        plaintext: encodeString('plaintext'),
      }),
      'passes expected args to transitEncrypt'
    );

    assert.strictEqual(find('[data-test-encrypted-value="ciphertext"]').innerText, 'secret');

    // exit modal
    await click('dialog button');

    // Encrypt again, with pre-encoded value and checkbox selected
    const preEncodedValue = encodeString('plaintext');

    await waitFor('.cm-editor');
    editor = codemirror('#plaintext-control');
    setCodeEditorValue(editor, preEncodedValue);

    await click('input[data-test-transit-input="encodedBase64"]');
    await click('button[type="submit"]');

    assert.strictEqual(this.apiStub.callCount, 2, 'calls the API to encrypt again');

    assert.true(
      this.apiStub.secondCall.calledWith('akey', 'transit', {
        plaintext: preEncodedValue,
      }),
      'passes pre-encoded value without re-encoding'
    );

    await click('dialog button');
  }

  test('it encrypts', doEncrypt);

  test('it shows key version selection', async function (assert) {
    const keyDefaults = {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['encrypt'],
    };
    const keyattrs = {
      keysForEncryption: [3, 2, 1],
      latestVersion: 3,
    };

    const key = { ...keyDefaults, ...keyattrs };
    this.set('key', key);
    const encryptStub = sinon.stub(this.owner.lookup('service:api').secrets, 'transitEncrypt').resolves({
      data: { ciphertext: 'secret' },
    });

    await render(hbs`
    <TransitKeyActions
      @selectedAction="encrypt"
      @key={{this.key}}
    />
  `);

    await waitFor('.cm-editor');

    const editor = codemirror();
    setCodeEditorValue(editor, 'plaintext');

    assert.dom('#key_version').exists({ count: 1 }, 'it renders the key version selector');

    await triggerEvent('#key_version', 'change');
    await click('button[type="submit"]');

    assert.true(encryptStub.calledOnce, 'calls transitEncrypt');

    assert.true(
      encryptStub.calledWith('akey', 'transit', {
        plaintext: encodeString('plaintext'),
        key_version: '0',
      }),
      'includes key_version in the payload'
    );
  });

  test('it hides key version selection', async function (assert) {
    const keyDefaults = {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['encrypt'],
    };

    const keyattrs = { keysForEncryption: [1] };
    const key = { ...keyDefaults, ...keyattrs };

    this.set('key', key);

    await render(hbs`
    <TransitKeyActions
      @selectedAction="encrypt"
      @key={{this.key}}
    />
  `);

    await waitFor('.cm-editor');

    const editor = codemirror('#plaintext-control');
    setCodeEditorValue(editor, 'plaintext');

    assert.dom('#key_version').doesNotExist('it does not render the selector when there is only one key');
  });

  const setupExport = async function () {
    this.set('key', {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['export'],
      exportKeyTypes: ['encryption'],
      validKeyVersions: [1],
    });
    await render(hbs`
    <TransitKeyActions @key={{this.key}} />`);
  };

  test('it can export a key:default behavior', async function (assert) {
    const exportStub = sinon.stub(this.owner.lookup('service:api').secrets, 'transitExportKey').resolves({
      wrap_info: { token: 'wrapped-token' },
    });

    await setupExport.call(this);
    await click('button[type="submit"]');

    assert.true(exportStub.calledOnce);
    assert.deepEqual(
      exportStub.firstCall.args,
      [
        'akey',
        'encryption-key',
        'transit',
        {
          headers: {
            'X-Vault-Wrap-TTL': '30m',
          },
          wrapTTL: '30m',
        },
      ],
      'passes expected args to api service'
    );

    assert.strictEqual(
      find('[data-test-encrypted-value="export"]').innerText,
      'wrapped-token',
      'wraps by default'
    );
  });

  test('it can export a key:unwrapped behavior', async function (assert) {
    const response = {
      data: {
        keys: { a: 'key' },
        type: 'encryption',
        name: 'akey',
      },
    };
    sinon.stub(this.owner.lookup('service:api').secrets, 'transitExportKey').resolves(response);

    await setupExport.call(this);
    await click('[data-test-toggle-label="Wrap response"]');
    await click(GENERAL.submitButton);

    assert.dom('#transit-export-modal').exists('Modal opens after export');
    assert.deepEqual(
      JSON.parse(find('[data-test-encrypted-value="export"]').innerText),
      {
        keys: { a: 'key' },
        type: 'encryption',
        name: 'akey',
      },
      'prints json response'
    );
  });

  test('it can export a key: unwrapped, single version', async function (assert) {
    const response = {
      data: {
        keys: { a: 'key' },
        type: 'encryption',
        name: 'akey',
      },
    };
    const exportVersionStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'transitExportKeyVersion')
      .resolves(response);

    await setupExport.call(this);
    await click('[data-test-toggle-label="Wrap response"]');
    await click('#exportVersion');
    await triggerEvent('#exportVersion', 'change');
    await click(GENERAL.submitButton);

    assert.dom('#transit-export-modal').exists('Modal opens after export');
    assert.deepEqual(
      JSON.parse(find('[data-test-encrypted-value="export"]').innerText),
      {
        keys: { a: 'key' },
        type: 'encryption',
        name: 'akey',
      },
      'prints json response'
    );

    assert.true(exportVersionStub.calledOnce);
    assert.deepEqual(
      exportVersionStub.firstCall.args,
      ['akey', 'encryption-key', 1, 'transit', {}],
      'passes expected args to api service'
    );
  });

  test('it includes algorithm param for HMAC', async function (assert) {
    const hmacStub = sinon.stub(this.owner.lookup('service:api').secrets, 'transitGenerateHmac').resolves({
      data: {
        hmac: 'vault:v1:hmac-token',
      },
    });

    this.set('key', {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['hmac'],
      validKeyVersions: [1],
    });

    await render(hbs`
    <TransitKeyActions
      @key={{this.key}}
      @selectedAction="hmac"
    />
  `);

    await fillIn('#algorithm', 'sha2-384');
    await blur('#algorithm');

    await waitFor('.cm-editor');
    const editor = codemirror();
    setCodeEditorValue(editor, 'plaintext');

    await click('input[data-test-transit-input="encodedBase64"]');
    await click(GENERAL.submitButton);

    assert.true(hmacStub.calledOnce, 'calls transitGenerateHmac');

    assert.deepEqual(
      hmacStub.firstCall.args,
      [
        'akey',
        'transit',
        {
          algorithm: 'sha2-384',
          input: 'plaintext',
        },
      ],
      'passes expected args to the API'
    );
  });
});
