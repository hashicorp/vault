/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Model | secret-engine', function (hooks) {
  setupTest(hooks);
  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
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

    test('it returns correct fields for aws', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'aws',
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
        'config.identityTokenKey',
      ]);
    });
  });

  module('formFieldGroups', function () {
    test('returns correct values by default', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'cubbyhole',
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

    test('returns correct values for aws', function (assert) {
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
            'config.identityTokenKey',
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
    test('returns default icon if no engineType', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: '',
      });
      assert.strictEqual(model.icon, 'lock', 'uses default icon');
    });
    test('returns default icon if kmip', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'kmip',
      });
      assert.strictEqual(model.icon, 'lock');
    });
    test('returns key if keymgmt', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'keymgmt',
      });
      assert.strictEqual(model.icon, 'key');
    });
    test('returns default when engine type is not in list of mountable engines', function (assert) {
      assert.expect(1);
      const model = this.store.createRecord('secret-engine', {
        type: 'ducks',
      });
      assert.strictEqual(model.icon, 'lock');
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
