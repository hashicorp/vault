/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { run } from '@ember/runloop';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, find, fillIn, blur, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { encodeString } from 'vault/utils/b64';
import waitForError from 'vault/tests/helpers/wait-for-error';
import codemirror from 'vault/tests/helpers/codemirror';

const storeStub = Service.extend({
  callArgs: null,
  keyActionReturnVal: null,
  rootKeyActionReturnVal: null,
  adapterFor() {
    const self = this;
    return {
      keyAction(action, { backend, id, payload }, options) {
        self.set('callArgs', { action, backend, id, payload });
        self.set('callArgsOptions', options);
        const rootResp = { ...self.get('rootKeyActionReturnVal') };
        const resp =
          Object.keys(rootResp).length > 0
            ? rootResp
            : {
                data: { ...self.get('keyActionReturnVal') },
              };
        return resolve(resp);
      },
    };
  },
});

module('Integration | Component | transit key actions', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeStub);
      this.storeService = this.owner.lookup('service:store');
    });
  });

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
    const keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat(actions) };

    const key = { ...keyDefaults, ...keyattrs };
    this.set('key', key);
    this.set('selectedAction', 'encrypt');
    this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
    await render(hbs`
    <TransitKeyActions @selectedAction={{this.selectedAction}} @key={{this.key}} />`);

    codemirror('#plaintext-control').setValue('plaintext');
    await click('button[type="submit"]');
    assert.deepEqual(
      this.storeService.callArgs,
      {
        action: 'encrypt',
        backend: 'transit',
        id: 'akey',
        payload: {
          plaintext: encodeString('plaintext'),
        },
      },
      'passes expected args to the adapter'
    );

    assert.strictEqual(find('[data-test-encrypted-value="ciphertext"]').innerText, 'secret');

    // exit modal
    await click('dialog button');
    // Encrypt again, with pre-encoded value and checkbox selected
    const preEncodedValue = encodeString('plaintext');
    codemirror('#plaintext-control').setValue(preEncodedValue);
    await click('input[data-test-transit-input="encodedBase64"]');
    await click('button[type="submit"]');

    assert.deepEqual(
      this.storeService.callArgs,
      {
        action: 'encrypt',
        backend: 'transit',
        id: 'akey',
        payload: {
          plaintext: preEncodedValue,
        },
      },
      'passes expected args to the adapter'
    );
    await click('dialog button');
  }

  test('it encrypts', doEncrypt);

  test('it shows key version selection', async function (assert) {
    const keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat([]) };
    const keyattrs = { keysForEncryption: [3, 2, 1], latestVersion: 3 };
    const key = { ...keyDefaults, ...keyattrs };
    this.set('key', key);
    this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
    await render(hbs`
    <TransitKeyActions @selectedAction="encrypt" @key={{this.key}} />`);

    codemirror().setValue('plaintext');
    assert.dom('#key_version').exists({ count: 1 }, 'it renders the key version selector');

    await triggerEvent('#key_version', 'change');
    await click('button[type="submit"]');
    assert.deepEqual(
      this.storeService.callArgs,
      {
        action: 'encrypt',
        backend: 'transit',
        id: 'akey',
        payload: {
          plaintext: encodeString('plaintext'),
          key_version: '0',
        },
      },
      'includes key_version in the payload'
    );
  });

  test('it hides key version selection', async function (assert) {
    const keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat([]) };
    const keyattrs = { keysForEncryption: [1] };
    const key = { ...keyDefaults, ...keyattrs };
    this.set('key', key);
    this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
    await render(hbs`
    <TransitKeyActions @selectedAction="encrypt" @key={{this.key}} />`);

    codemirror('#plaintext-control').setValue('plaintext');
    assert.dom('#key_version').doesNotExist('it does not render the selector when there is only one key');
  });

  test('it does not carry ciphertext value over to decrypt', async function (assert) {
    assert.expect(4);
    const plaintext = 'not so secret';
    await doEncrypt.call(this, assert, ['decrypt']);

    this.set('storeService.keyActionReturnVal', { plaintext });
    this.set('selectedAction', 'decrypt');
    assert.strictEqual(codemirror('#ciphertext-control').getValue(), '', 'does not prefill ciphertext value');
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
    this.set('storeService.rootKeyActionReturnVal', { wrap_info: { token: 'wrapped-token' } });
    await setupExport.call(this);
    await click('button[type="submit"]');

    assert.deepEqual(
      this.storeService.callArgs,
      {
        action: 'export',
        backend: 'transit',
        id: 'akey',
        payload: {
          param: ['encryption'],
        },
      },
      'passes expected args to the adapter'
    );
    assert.strictEqual(this.storeService.callArgsOptions.wrapTTL, '30m', 'passes value for wrapTTL');
    assert.strictEqual(
      find('[data-test-encrypted-value="export"]').innerText,
      'wrapped-token',
      'wraps by default'
    );
  });

  test('it can export a key:unwrapped behavior', async function (assert) {
    const response = { keys: { a: 'key' } };
    this.set('storeService.keyActionReturnVal', response);
    await setupExport.call(this);
    await click('[data-test-toggle-label="Wrap response"]');
    await click('button[type="submit"]');
    assert.dom('#transit-export-modal').exists('Modal opens after export');
    assert.deepEqual(
      JSON.parse(find('[data-test-encrypted-value="export"]').innerText),
      response,
      'prints json response'
    );
  });

  test('it can export a key: unwrapped, single version', async function (assert) {
    const response = { keys: { a: 'key' } };
    this.set('storeService.keyActionReturnVal', response);
    await setupExport.call(this);
    await click('[data-test-toggle-label="Wrap response"]');
    await click('#exportVersion');
    await triggerEvent('#exportVersion', 'change');
    await click('button[type="submit"]');
    assert.dom('#transit-export-modal').exists('Modal opens after export');
    assert.deepEqual(
      JSON.parse(find('[data-test-encrypted-value="export"]').innerText),
      response,
      'prints json response'
    );
    assert.deepEqual(
      this.storeService.callArgs,
      {
        action: 'export',
        backend: 'transit',
        id: 'akey',
        payload: {
          param: ['encryption', 1],
        },
      },
      'passes expected args to the adapter'
    );
  });

  test('it includes algorithm param for HMAC', async function (assert) {
    // Return mocked data so a11y-testing doesn't get mad about empty copy button contents
    this.set('storeService.rootKeyActionReturnVal', { data: { hmac: 'vault:v1:hmac-token' } });
    this.set('key', {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['hmac'],
      validKeyVersions: [1],
    });
    await render(hbs`
    <TransitKeyActions @key={{this.key}} @selectedAction="hmac" />`);
    await fillIn('#algorithm', 'sha2-384');
    await blur('#algorithm');
    await fillIn('[data-test-component="code-mirror-modifier"] textarea', 'plaintext');
    await click('input[data-test-transit-input="encodedBase64"]');
    await click('button[type="submit"]');
    assert.deepEqual(
      this.storeService.callArgs,
      {
        action: 'hmac',
        backend: 'transit',
        id: 'akey',
        payload: {
          algorithm: 'sha2-384',
          input: 'plaintext',
        },
      },
      'passes expected args to the adapter'
    );
  });
});
