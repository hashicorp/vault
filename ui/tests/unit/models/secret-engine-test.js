/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Model | secret-engine', function (hooks) {
  setupTest(hooks);
  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  module('modelTypeForKV', function () {
    test('is secret by default', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine');
      assert.strictEqual(model.get('modelTypeForKV'), 'secret');
    });

    test('is secret-v2 for kv v2', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        version: 2,
        type: 'kv',
      });
      assert.strictEqual(model.get('modelTypeForKV'), 'secret-v2');
    });

    test('is secret-v2 for generic v2', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        version: 2,
        type: 'kv',
      });

      assert.strictEqual(model.get('modelTypeForKV'), 'secret-v2');
    });

    test('is secret when v2 if not kv or generic', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        version: 2,
        type: 'ssh',
      });

      assert.strictEqual(model.get('modelTypeForKV'), 'secret');
    });
  });

  module('formFields', function () {
    test('it returns correct fields by default', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: '',
      });

      assert.deepEqual(model.get('formFields'), [
        'type',
        'path',
        'description',
        'accessor',
        'local',
        'sealWrap',
        'config.defaultLeaseTtl',
        'config.maxLeaseTtl',
        'config.allowedManagedKeys',
        'config.auditNonHmacRequestKeys',
        'config.auditNonHmacResponseKeys',
        'config.passthroughRequestHeaders',
        'config.allowedResponseHeaders',
      ]);
    });

    test('it returns correct fields for KV v1', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'kv',
      });

      assert.deepEqual(model.get('formFields'), [
        'type',
        'path',
        'description',
        'accessor',
        'local',
        'sealWrap',
        'config.defaultLeaseTtl',
        'config.maxLeaseTtl',
        'config.allowedManagedKeys',
        'config.auditNonHmacRequestKeys',
        'config.auditNonHmacResponseKeys',
        'config.passthroughRequestHeaders',
        'config.allowedResponseHeaders',
        'version',
      ]);
    });

    test('it returns correct fields for KV v2', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'kv',
        version: '2',
      });

      assert.deepEqual(model.get('formFields'), [
        'type',
        'path',
        'description',
        'accessor',
        'local',
        'sealWrap',
        'config.defaultLeaseTtl',
        'config.maxLeaseTtl',
        'config.allowedManagedKeys',
        'config.auditNonHmacRequestKeys',
        'config.auditNonHmacResponseKeys',
        'config.passthroughRequestHeaders',
        'config.allowedResponseHeaders',
        'version',
        'casRequired',
        'deleteVersionAfter',
        'maxVersions',
      ]);
    });

    test('it returns correct fields for keymgmt', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'keymgmt',
      });

      assert.deepEqual(model.get('formFields'), [
        'type',
        'path',
        'description',
        'accessor',
        'local',
        'sealWrap',
        'config.allowedManagedKeys',
        'config.auditNonHmacRequestKeys',
        'config.auditNonHmacResponseKeys',
        'config.passthroughRequestHeaders',
        'config.allowedResponseHeaders',
      ]);
    });
  });

  module('formFieldGroups', function () {
    test('returns correct values by default', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'aws',
      });

      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path'] },
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.defaultLeaseTtl',
            'config.maxLeaseTtl',
            'config.allowedManagedKeys',
            'config.auditNonHmacRequestKeys',
            'config.auditNonHmacResponseKeys',
            'config.passthroughRequestHeaders',
            'config.allowedResponseHeaders',
          ],
        },
      ]);
    });
    test('returns correct values for KV', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'kv',
      });

      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path', 'maxVersions', 'casRequired', 'deleteVersionAfter'] },
        {
          'Method Options': [
            'version',
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.defaultLeaseTtl',
            'config.maxLeaseTtl',
            'config.allowedManagedKeys',
            'config.auditNonHmacRequestKeys',
            'config.auditNonHmacResponseKeys',
            'config.passthroughRequestHeaders',
            'config.allowedResponseHeaders',
          ],
        },
      ]);
    });

    test('returns correct values for generic', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'generic',
      });

      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path'] },
        {
          'Method Options': [
            'version',
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.defaultLeaseTtl',
            'config.maxLeaseTtl',
            'config.allowedManagedKeys',
            'config.auditNonHmacRequestKeys',
            'config.auditNonHmacResponseKeys',
            'config.passthroughRequestHeaders',
            'config.allowedResponseHeaders',
          ],
        },
      ]);
    });

    test('returns correct values for database', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'database',
      });

      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path', 'config.defaultLeaseTtl', 'config.maxLeaseTtl'] },
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.allowedManagedKeys',
            'config.auditNonHmacRequestKeys',
            'config.auditNonHmacResponseKeys',
            'config.passthroughRequestHeaders',
            'config.allowedResponseHeaders',
          ],
        },
      ]);
    });

    test('returns correct values for pki', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'pki',
      });

      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path', 'config.defaultLeaseTtl', 'config.maxLeaseTtl', 'config.allowedManagedKeys'] },
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.auditNonHmacRequestKeys',
            'config.auditNonHmacResponseKeys',
            'config.passthroughRequestHeaders',
            'config.allowedResponseHeaders',
          ],
        },
      ]);
    });

    test('returns correct values for keymgmt', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'keymgmt',
      });

      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path'] },
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.allowedManagedKeys',
            'config.auditNonHmacRequestKeys',
            'config.auditNonHmacResponseKeys',
            'config.passthroughRequestHeaders',
            'config.allowedResponseHeaders',
          ],
        },
      ]);
    });
  });

  module('engineType', function () {
    test('strips leading ns_ from type', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        // eg. ns_cubbyhole, ns_identity, ns_system
        type: 'ns_identity',
      });
      assert.strictEqual(model.engineType, 'identity');
    });
    test('returns type by default', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'zebras',
      });
      assert.strictEqual(model.engineType, 'zebras');
    });
  });

  module('icon', function () {
    test('returns secrets if no engineType', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: '',
      });
      assert.strictEqual(model.icon, 'secrets');
    });
    test('returns secrets if kmip', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'kmip',
      });
      assert.strictEqual(model.icon, 'secrets');
    });
    test('returns key if keymgmt', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'keymgmt',
      });
      assert.strictEqual(model.icon, 'key');
    });
    test('returns engineType by default', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'ducks',
      });
      assert.strictEqual(model.icon, 'ducks');
    });
  });

  module('shouldIncludeInList', function () {
    test('returns false if excludeList includes type', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'system',
      });
      assert.false(model.shouldIncludeInList);
    });
    test('returns true if excludeList does not include type', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'hippos',
      });
      assert.true(model.shouldIncludeInList);
    });
  });

  module('localDisplay', function () {
    test('returns local if local', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        local: true,
      });
      assert.strictEqual(model.localDisplay, 'local');
    });
    test('returns replicated if !local', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        local: false,
      });
      assert.strictEqual(model.localDisplay, 'replicated');
    });
  });

  module('saveCA', function () {
    test('does not call endpoint if type != ssh', async function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'not-ssh',
      });
      const saveSpy = sinon.spy(model, 'save');
      await model.saveCA({});
      assert.ok(saveSpy.notCalled, 'save not called');
    });
    test('calls save with correct params', async function (assert) {
      assert.expect(4);
      const model = this.store.createRecord('secret-engine', {
        type: 'ssh',
        privateKey: 'private-key',
        publicKey: 'public-key',
        generateSigningKey: true,
      });
      const saveStub = sinon.stub(model, 'save').callsFake((params) => {
        assert.deepEqual(
          params,
          {
            adapterOptions: {
              options: {},
              apiPath: 'config/ca',
              attrsToSend: ['privateKey', 'publicKey', 'generateSigningKey'],
            },
          },
          'send correct params to save'
        );
        return;
      });

      await model.saveCA({});
      assert.strictEqual(model.privateKey, 'private-key', 'value exists before save');
      assert.strictEqual(model.publicKey, 'public-key', 'value exists before save');
      assert.true(model.generateSigningKey, 'value true before save');

      saveStub.restore();
    });
    test('sets properties when isDelete', async function (assert) {
      assert.expect(7);
      const model = this.store.createRecord('secret-engine', {
        type: 'ssh',
        privateKey: 'private-key',
        publicKey: 'public-key',
        generateSigningKey: true,
      });
      const saveStub = sinon.stub(model, 'save').callsFake((params) => {
        assert.deepEqual(
          params,
          {
            adapterOptions: {
              options: { isDelete: true },
              apiPath: 'config/ca',
              attrsToSend: ['privateKey', 'publicKey', 'generateSigningKey'],
            },
          },
          'send correct params to save'
        );
        return;
      });
      assert.strictEqual(model.privateKey, 'private-key', 'value exists before save');
      assert.strictEqual(model.publicKey, 'public-key', 'value exists before save');
      assert.true(model.generateSigningKey, 'value true before save');

      await model.saveCA({ isDelete: true });
      assert.strictEqual(model.privateKey, null, 'value null after save');
      assert.strictEqual(model.publicKey, null, 'value null after save');
      assert.false(model.generateSigningKey, 'value false after save');
      saveStub.restore();
    });
  });

  module('saveZeroAddressConfig', function () {
    test('calls save with correct params', async function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {});
      const saveStub = sinon.stub(model, 'save').callsFake((params) => {
        assert.deepEqual(
          params,
          {
            adapterOptions: {
              adapterMethod: 'saveZeroAddressConfig',
            },
          },
          'send correct params to save'
        );
        return;
      });
      await model.saveZeroAddressConfig();
      saveStub.restore();
    });
  });
});
