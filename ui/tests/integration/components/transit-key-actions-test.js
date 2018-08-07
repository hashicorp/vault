import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import Ember from 'ember';
import { encodeString } from 'vault/utils/b64';

const storeStub = Ember.Service.extend({
  callArgs: null,
  keyActionReturnVal: null,
  rootKeyActionReturnVal: null,
  adapterFor() {
    const self = this;
    return {
      keyAction(action, { backend, id, payload }, options) {
        self.set('callArgs', { action, backend, id, payload });
        self.set('callArgsOptions', options);
        const rootResp = Ember.assign({}, self.get('rootKeyActionReturnVal'));
        const resp =
          Object.keys(rootResp).length > 0
            ? rootResp
            : {
                data: Ember.assign({}, self.get('keyActionReturnVal')),
              };
        return Ember.RSVP.resolve(resp);
      },
    };
  },
});

moduleForComponent('transit-key-actions', 'Integration | Component | transit key actions', {
  integration: true,
  beforeEach: function() {
    this.register('service:store', storeStub);
    this.inject.service('store', { as: 'storeService' });
  },
});

test('it requires `key`', function(assert) {
  assert.expectAssertion(
    () => this.render(hbs`{{transit-key-actions}}`),
    /`key` is required for/,
    'asserts without key'
  );
});

test('it renders', function(assert) {
  this.set('key', { backend: 'transit', supportedActions: ['encrypt'] });
  this.render(hbs`{{transit-key-actions selectedAction="encrypt" key=key}}`);
  assert.equal(this.$('[data-test-transit-action="encrypt"]').length, 1, 'renders encrypt');

  this.set('key', { backend: 'transit', supportedActions: ['sign'] });
  this.render(hbs`{{transit-key-actions selectedAction="sign" key=key}}`);
  assert.equal(this.$('[data-test-transit-action="sign"]').length, 1, 'renders sign');
});

test('it renders: signature_algorithm field', function(assert) {
  this.set('key', { backend: 'transit', supportsSigning: true, supportedActions: ['sign', 'verify'] });
  this.set('selectedAction', 'sign');
  this.render(hbs`{{transit-key-actions selectedAction=selectedAction key=key}}`);
  assert.equal(
    this.$('[data-test-signature-algorithm]').length,
    0,
    'does not render signature_algorithm field on sign'
  );
  this.set('selectedAction', 'verify');
  assert.equal(
    this.$('[data-test-signature-algorithm]').length,
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
    this.$('[data-test-signature-algorithm]').length,
    1,
    'renders signature_algorithm field on sign with rsa key'
  );
  this.set('selectedAction', 'verify');
  assert.equal(
    this.$('[data-test-signature-algorithm]').length,
    1,
    'renders signature_algorithm field on verify with rsa key'
  );
});

test('it renders: rotate', function(assert) {
  this.set('key', { backend: 'transit', id: 'akey', supportedActions: ['rotate'] });
  this.render(hbs`{{transit-key-actions selectedAction="rotate" key=key}}`);

  assert.equal(this.$().text().trim(), '', 'renders an empty div');

  this.set('key.canRotate', true);
  assert.equal(
    this.$('button').text().trim(),
    'Rotate encryption key',
    'renders confirm-button when key.canRotate is true'
  );
});

function doEncrypt(assert, actions = [], keyattrs = {}) {
  let keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat(actions) };

  const key = Ember.assign({}, keyDefaults, keyattrs);
  this.set('key', key);
  this.set('selectedAction', 'encrypt');
  this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
  this.render(hbs`{{transit-key-actions selectedAction=selectedAction key=key}}`);

  this.$('#plaintext').val('plaintext').trigger('input');
  this.$('[data-test-transit-b64-toggle="plaintext"]').click();
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

  assert.equal(this.$('#ciphertext').val(), 'secret');
}

test('it encrypts', doEncrypt);

test('it shows key version selection', function(assert) {
  let keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat([]) };
  let keyattrs = { keysForEncryption: [3, 2, 1], latestVersion: 3 };
  const key = Ember.assign({}, keyDefaults, keyattrs);
  this.set('key', key);
  this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
  this.render(hbs`{{transit-key-actions selectedAction="encrypt" key=key}}`);

  this.$('#plaintext').val('plaintext').trigger('input');
  this.$('[data-test-transit-b64-toggle="plaintext"]').click();
  assert.equal(this.$('#key_version').length, 1, 'it renders the key version selector');

  this.$('#key_version').trigger('change');
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

test('it hides key version selection', function(assert) {
  let keyDefaults = { backend: 'transit', id: 'akey', supportedActions: ['encrypt'].concat([]) };
  let keyattrs = { keysForEncryption: [1] };
  const key = Ember.assign({}, keyDefaults, keyattrs);
  this.set('key', key);
  this.set('storeService.keyActionReturnVal', { ciphertext: 'secret' });
  this.render(hbs`{{transit-key-actions selectedAction="encrypt" key=key}}`);

  this.$('#plaintext').val('plaintext').trigger('input');
  this.$('[data-test-transit-b64-toggle="plaintext"]').click();

  assert.equal(
    this.$('#key_version').length,
    0,
    'it does not render the selector when there is only one key'
  );
});

test('it carries ciphertext value over to decrypt', function(assert) {
  const plaintext = 'not so secret';
  doEncrypt.call(this, assert, ['decrypt']);

  this.set('storeService.keyActionReturnVal', { plaintext });
  this.set('selectedAction', 'decrypt');
  assert.equal(this.$('#ciphertext').val(), 'secret', 'keeps ciphertext value');

  this.$('button:submit').click();
  assert.equal(this.$('#plaintext').val(), plaintext, 'renders decrypted value');
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
  assert.equal(this.$('#export').val(), 'wrapped-token', 'wraps by default');
});

test('it can export a key:unwrapped behavior', function(assert) {
  const response = { keys: { a: 'key' } };
  this.set('storeService.keyActionReturnVal', response);
  setupExport.call(this);
  this.$('#wrap-response').click().change();
  this.$('button:submit').click();
  assert.deepEqual(
    JSON.parse(this.$('.CodeMirror').get(0).CodeMirror.getValue()),
    response,
    'prints json response'
  );
});

test('it can export a key: unwrapped, single version', function(assert) {
  const response = { keys: { a: 'key' } };
  this.set('storeService.keyActionReturnVal', response);
  setupExport.call(this);
  this.$('#wrap-response').click().change();
  this.$('#exportVersion').click().change();
  this.$('button:submit').click();
  assert.deepEqual(
    JSON.parse(this.$('.CodeMirror').get(0).CodeMirror.getValue()),
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

test('it includes algorithm param for HMAC', function(assert) {
  this.set('key', {
    backend: 'transit',
    id: 'akey',
    supportedActions: ['hmac'],
    validKeyVersions: [1],
  });
  this.render(hbs`{{transit-key-actions key=key}}`);
  this.$('#algorithm').val('sha2-384').change();
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
