import Ember from 'ember';
import { moduleForModel, test } from 'ember-qunit';

moduleForModel('transit-key', 'Unit | Model | transit key');

test('it exists', function(assert) {
  let model = this.subject();
  assert.ok(!!model);
});

test('supported actions', function(assert) {
  let model = this.subject({
    supportsEncryption: true,
    supportsDecryption: true,
    supportsSigning: false,
  });

  let supportedActions = model.get('supportedActions');
  assert.deepEqual(['encrypt', 'decrypt', 'datakey', 'rewrap', 'hmac', 'verify'], supportedActions);
});

test('encryption key versions', function(assert) {
  let done = assert.async();
  let model = this.subject({
    keyVersions: [1, 2, 3, 4, 5],
    minDecryptionVersion: 1,
    latestVersion: 5,
  });
  assert.deepEqual([5, 4, 3, 2, 1], model.get('encryptionKeyVersions'), 'lists all available versions');
  Ember.run(() => {
    model.set('minDecryptionVersion', 3);
    assert.deepEqual(
      [5, 4, 3],
      model.get('encryptionKeyVersions'),
      'adjusts to a change in minDecryptionVersion'
    );
    done();
  });
});

test('keys for encryption', function(assert) {
  let done = assert.async();
  let model = this.subject({
    keyVersions: [1, 2, 3, 4, 5],
    minDecryptionVersion: 1,
    latestVersion: 5,
  });

  assert.deepEqual(
    [5, 4, 3, 2, 1],
    model.get('keysForEncryption'),
    'lists all available versions when no min is set'
  );

  Ember.run(() => {
    model.set('minEncryptionVersion', 4);
    assert.deepEqual([5, 4], model.get('keysForEncryption'), 'calculates using minEncryptionVersion');
    done();
  });
});
