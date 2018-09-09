import { resolve } from 'rsvp';
import { assign } from '@ember/polyfills';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, find, findAll, fillIn, blur, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { encodeString } from 'vault/utils/b64';

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
        const rootResp = assign({}, self.get('rootKeyActionReturnVal'));
        const resp =
          Object.keys(rootResp).length > 0
            ? rootResp
            : {
                data: assign({}, self.get('keyActionReturnVal')),
              };
        return resolve(resp);
      },
    };
  },
});

module('Integration | Component | transit key actions', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.register('service:store', storeStub);
    this.storeService = this.owner.lookup('service:store');
  });

  test('it requires `key`', function(assert) {
    assert.expectAssertion(
      () => this.render(hbs`{{transit-key-actions}}`),
      /`key` is required for/,
      'asserts without key'
    );
  });

  test('it renders', async function(assert) {
    this.set('key', { backend: 'transit', supportedActions: ['encrypt'] });
    await render(hbs`{{transit-key-actions selectedAction="encrypt" key=key}}`);
    assert.equal(findAll('[data-test-transit-action="encrypt"]').length, 1, 'renders encrypt');

    this.set('key', { backend: 'transit', supportedActions: ['sign'] });
    await render(hbs`{{transit-key-actions selectedAction="sign" key=key}}`);
    assert.equal(findAll('[data-test-transit-action="sign"]').length, 1, 'renders sign');
  });

  test('it renders: signature_algorithm field', async function(assert) {
    this.set('key', { backend: 'transit', supportsSigning: true, supportedActions: ['sign', 'verify'] });
    this.set('selectedAction', 'sign');
    await render(hbs`{{transit-key-actions selectedAction=selectedAction key=key}}`);
    assert.equal(
      findAll('[data-test-signature-algorithm]').length,
      0,
      'does not render signature_algorithm field on sign'
    );
    this.set('selectedAction', 'verify');
    assert.equal(
      findAll('[data-test-signature-algorithm]').length,
      0,
      'does not render signature_algorithm field on verify'
    );

    this.set('selectedAction', 'sign');
    this.set('key', {
      type: 'rsa-2048',
      supportsSigning: true,
      backend: 'transit',
      supportedActions: ['sign', 'verify'],
    });
    assert.equal(
      findAll('[data-test-signature-algorithm]').length,
      1,
      'renders signature_algorithm field on sign with rsa key'
    );
    this.set('selectedAction', 'verify');
    assert.equal(
      findAll('[data-test-signature-algorithm]').length,
      1,
      'renders signature_algorithm field on verify with rsa key'
    );
  });

  test('it renders: rotate', async function(assert) {
    this.set('key', { backend: 'transit', id: 'akey', supportedActions: ['rotate'] });
    await render(hbs`{{transit-key-actions selectedAction="rotate" key=key}}`);

    assert.equal(find('*').textContent.trim(), '', 'renders an empty div');

    this.set('key.canRotate', true);
    assert.equal(
      find('button').textContent.trim(),
      'Rotate encryption key',
      'renders confirm-button when key.canRotate is true'
    );
  });

  async function doEncrypt(assert, actions = [], keyattrs = {}) {
    let keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat(actions) };

    const key = assign({}, keyDefaults, keyattrs);
    this.set('key', key);
    this.set('selectedAction', 'encrypt');
    this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
    this.render(hbs`{{transit-key-actions selectedAction=selectedAction key=key}}`);

    await fillIn('#plaintext', 'plaintext');
    await click('[data-test-transit-b64-toggle="plaintext"]');
    this.$('button:submit').click();
    assert.deepEqual(
      this.get('storeService.callArgs'),
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

    assert.equal(find('#ciphertext').value, 'secret');
  }

  test('it encrypts', doEncrypt);

  test('it shows key version selection', async function(assert) {
    let keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat([]) };
    let keyattrs = { keysForEncryption: [3, 2, 1], latestVersion: 3 };
    const key = assign({}, keyDefaults, keyattrs);
    this.set('key', key);
    this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
    await render(hbs`{{transit-key-actions selectedAction="encrypt" key=key}}`);

    await fillIn('#plaintext', 'plaintext');
    await click('[data-test-transit-b64-toggle="plaintext"]');
    assert.equal(findAll('#key_version').length, 1, 'it renders the key version selector');

    await triggerEvent('#key_version', 'change');
    this.$('button:submit').click();
    assert.deepEqual(
      this.get('storeService.callArgs'),
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

  test('it hides key version selection', async function(assert) {
    let keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat([]) };
    let keyattrs = { keysForEncryption: [1] };
    const key = assign({}, keyDefaults, keyattrs);
    this.set('key', key);
    this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
    await render(hbs`{{transit-key-actions selectedAction="encrypt" key=key}}`);

    await fillIn('#plaintext', 'plaintext');
    await click('[data-test-transit-b64-toggle="plaintext"]');

    assert.equal(
      findAll('#key_version').length,
      0,
      'it does not render the selector when there is only one key'
    );
  });

  test('it carries ciphertext value over to decrypt', function(assert) {
    const plaintext = 'not so secret';
    doEncrypt.call(this, assert, ['decrypt']);

    this.set('storeService.keyActionReturnVal', { plaintext });
    this.set('selectedAction', 'decrypt');
    assert.equal(find('#ciphertext').value, 'secret', 'keeps ciphertext value');

    this.$('button:submit').click();
    assert.equal(find('#plaintext').value, plaintext, 'renders decrypted value');
  });

  const setupExport = function() {
    this.set('key', {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['export'],
      exportKeyTypes: ['encryption'],
      validKeyVersions: [1],
    });
    this.render(hbs`{{transit-key-actions key=key}}`);
  };

  test('it can export a key:default behavior', function(assert) {
    this.set('storeService.rootKeyActionReturnVal', { wrap_info: { token: 'wrapped-token' } });
    setupExport.call(this);
    this.$('button:submit').click();

    assert.deepEqual(
      this.get('storeService.callArgs'),
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
    assert.equal(this.get('storeService.callArgsOptions.wrapTTL'), '30m', 'passes value for wrapTTL');
    assert.equal(find('#export').value, 'wrapped-token', 'wraps by default');
  });

  test('it can export a key:unwrapped behavior', async function(assert) {
    const response = { keys: { a: 'key' } };
    this.set('storeService.keyActionReturnVal', response);
    setupExport.call(this);
    await click('#wrap-response').change();
    this.$('button:submit').click();
    assert.deepEqual(
      JSON.parse(findAll('.CodeMirror')[0].CodeMirror.getValue()),
      response,
      'prints json response'
    );
  });

  test('it can export a key: unwrapped, single version', async function(assert) {
    const response = { keys: { a: 'key' } };
    this.set('storeService.keyActionReturnVal', response);
    setupExport.call(this);
    await click('#wrap-response').change();
    await click('#exportVersion').change();
    this.$('button:submit').click();
    assert.deepEqual(
      JSON.parse(findAll('.CodeMirror')[0].CodeMirror.getValue()),
      response,
      'prints json response'
    );
    assert.deepEqual(
      this.get('storeService.callArgs'),
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

  test('it includes algorithm param for HMAC', async function(assert) {
    this.set('key', {
      backend: 'transit',
      id: 'akey',
      supportedActions: ['hmac'],
      validKeyVersions: [1],
    });
    await render(hbs`{{transit-key-actions key=key}}`);
    await fillIn('#algorithm', 'sha2-384');
    await blur('#algorithm');
    this.$('button:submit').click();
    assert.deepEqual(
      this.get('storeService.callArgs'),
      {
        action: 'hmac',
        backend: 'transit',
        id: 'akey',
        payload: {
          algorithm: 'sha2-384',
        },
      },
      'passes expected args to the adapter'
    );
  });
});
