import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Model | transit key', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let model = run(() => this.owner.lookup('service:store').createRecord('transit-key'));
    assert.ok(!!model);
  });

  test('supported actions', function(assert) {
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('transit-key', {
        supportsEncryption: true,
        supportsDecryption: true,
        supportsSigning: false,
      })
    );

    let supportedActions = model.get('supportedActions');
    assert.deepEqual(['encrypt', 'decrypt', 'datakey', 'rewrap', 'hmac', 'verify'], supportedActions);
  });

  test('encryption key versions', function(assert) {
    let done = assert.async();
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('transit-key', {
        keyVersions: [1, 2, 3, 4, 5],
        minDecryptionVersion: 1,
        latestVersion: 5,
      })
    );
    assert.deepEqual([5, 4, 3, 2, 1], model.get('encryptionKeyVersions'), 'lists all available versions');
    run(() => {
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
    let model = run(() =>
      this.owner.lookup('service:store').createRecord('transit-key', {
        keyVersions: [1, 2, 3, 4, 5],
        minDecryptionVersion: 1,
        latestVersion: 5,
      })
    );

    assert.deepEqual(
      [5, 4, 3, 2, 1],
      model.get('keysForEncryption'),
      'lists all available versions when no min is set'
    );

    run(() => {
      model.set('minEncryptionVersion', 4);
      assert.deepEqual([5, 4], model.get('keysForEncryption'), 'calculates using minEncryptionVersion');
      done();
    });
  });
});
