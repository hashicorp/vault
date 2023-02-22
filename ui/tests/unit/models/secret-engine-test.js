import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Model | secret-engine', function (hooks) {
  setupTest(hooks);

  test('modelTypeForKV is secret by default', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() => this.owner.lookup('service:store').createRecord('secret-engine'));
      assert.strictEqual(model.get('modelTypeForKV'), 'secret');
    });
  });

  test('modelTypeForKV is secret-v2 for kv v2', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          version: 2,
          type: 'kv',
        })
      );
      assert.strictEqual(model.get('modelTypeForKV'), 'secret-v2');
    });
  });

  test('modelTypeForKV is secret-v2 for generic v2', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          version: 2,
          type: 'kv',
        })
      );
      assert.strictEqual(model.get('modelTypeForKV'), 'secret-v2');
    });
  });

  test('formFieldGroups returns correct values by default', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          type: 'aws',
        })
      );
      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path'] },
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}',
          ],
        },
      ]);
    });
  });

  test('formFieldGroups returns correct values for KV', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          type: 'kv',
        })
      );
      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path', 'maxVersions', 'casRequired', 'deleteVersionAfter'] },
        {
          'Method Options': [
            'version',
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}',
          ],
        },
      ]);
    });
  });

  test('formFieldGroups returns correct values for generic', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          type: 'generic',
        })
      );
      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path'] },
        {
          'Method Options': [
            'version',
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}',
          ],
        },
      ]);
    });
  });

  test('formFieldGroups returns correct values for database', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          type: 'database',
        })
      );
      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path', 'config.{defaultLeaseTtl}', 'config.{maxLeaseTtl}'] },
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.{auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}',
          ],
        },
      ]);
    });
  });

  test('formFieldGroups returns correct values for keymgmt', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          type: 'keymgmt',
        })
      );
      assert.deepEqual(model.get('formFieldGroups'), [
        { default: ['path'] },
        {
          'Method Options': [
            'description',
            'config.listingVisibility',
            'local',
            'sealWrap',
            'config.{auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}',
          ],
        },
      ]);
    });
  });
});
