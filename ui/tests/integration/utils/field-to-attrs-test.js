import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

module('Integration | Util | field to attrs', function(hooks) {
  setupTest(hooks);

  const PATH_ATTR = { type: 'string', name: 'path', options: {} };
  const DESCRIPTION_ATTR = { type: 'string', name: 'description', options: { editType: 'textarea' } };
  const DEFAULT_LEASE_ATTR = {
    type: undefined,
    name: 'config.defaultLeaseTtl',
    options: { label: 'Default Lease TTL', editType: 'ttl' },
  };

  const OTHER_DEFAULT_LEASE_ATTR = {
    type: undefined,
    name: 'otherConfig.defaultLeaseTtl',
    options: { label: 'Default Lease TTL', editType: 'ttl' },
  };
  const MAX_LEASE_ATTR = {
    type: undefined,
    name: 'config.maxLeaseTtl',
    options: { label: 'Max Lease TTL', editType: 'ttl' },
  };
  const OTHER_MAX_LEASE_ATTR = {
    type: undefined,
    name: 'otherConfig.maxLeaseTtl',
    options: { label: 'Max Lease TTL', editType: 'ttl' },
  };

  test('it extracts attrs', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    run(() => {
      const [attr] = expandAttributeMeta(model, ['path']);
      assert.deepEqual(attr, PATH_ATTR, 'returns attribute meta');
    });
  });

  test('it extracts more than one attr', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    run(() => {
      const [path, desc] = expandAttributeMeta(model, ['path', 'description']);
      assert.deepEqual(path, PATH_ATTR, 'returns attribute meta');
      assert.deepEqual(desc, DESCRIPTION_ATTR, 'returns attribute meta');
    });
  });

  test('it extracts fieldGroups', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    run(() => {
      const groups = fieldToAttrs(model, [{ default: ['path'] }, { Options: ['description'] }]);
      const expected = [{ default: [PATH_ATTR] }, { Options: [DESCRIPTION_ATTR] }];
      assert.deepEqual(groups, expected, 'expands all given groups');
    });
  });

  test('it extracts arrays as fieldGroups', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    run(() => {
      const groups = fieldToAttrs(model, [
        { default: ['path', 'description'] },
        { Options: ['description'] },
      ]);
      const expected = [{ default: [PATH_ATTR, DESCRIPTION_ATTR] }, { Options: [DESCRIPTION_ATTR] }];
      assert.deepEqual(groups, expected, 'expands all given groups');
    });
  });

  test('it extracts model-fragment attributes with brace expansion', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    run(() => {
      const [attr] = expandAttributeMeta(model, ['config.{defaultLeaseTtl}']);
      assert.deepEqual(attr, DEFAULT_LEASE_ATTR, 'properly extracts model fragment attr');
    });

    run(() => {
      const [defaultLease, maxLease] = expandAttributeMeta(model, ['config.{defaultLeaseTtl,maxLeaseTtl}']);
      assert.deepEqual(defaultLease, DEFAULT_LEASE_ATTR, 'properly extracts default lease');
      assert.deepEqual(maxLease, MAX_LEASE_ATTR, 'properly extracts max lease');
    });
  });

  test('it extracts model-fragment attributes with double brace expansion', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    run(() => {
      const [configDefault, configMax, otherConfigDefault, otherConfigMax] = expandAttributeMeta(model, [
        '{config,otherConfig}.{defaultLeaseTtl,maxLeaseTtl}',
      ]);
      assert.deepEqual(configDefault, DEFAULT_LEASE_ATTR, 'properly extracts config.defaultLeaseTTL');
      assert.deepEqual(
        otherConfigDefault,
        OTHER_DEFAULT_LEASE_ATTR,
        'properly extracts otherConfig.defaultLeaseTTL'
      );

      assert.deepEqual(configMax, MAX_LEASE_ATTR, 'properly extracts config.maxLeaseTTL');
      assert.deepEqual(otherConfigMax, OTHER_MAX_LEASE_ATTR, 'properly extracts otherConfig.maxLeaseTTL');
    });
  });

  test('it extracts model-fragment attributes with dot notation', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    run(() => {
      const [attr] = expandAttributeMeta(model, ['config.defaultLeaseTtl']);
      assert.deepEqual(attr, DEFAULT_LEASE_ATTR, 'properly extracts model fragment attr');
    });

    run(() => {
      const [defaultLease, maxLease] = expandAttributeMeta(model, [
        'config.defaultLeaseTtl',
        'config.maxLeaseTtl',
      ]);
      assert.deepEqual(defaultLease, DEFAULT_LEASE_ATTR, 'properly extracts model fragment attr');
      assert.deepEqual(maxLease, MAX_LEASE_ATTR, 'properly extracts model fragment attr');
    });
  });

  test('it extracts fieldGroups from model-fragment attributes with brace expansion', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    const expected = [
      { default: [PATH_ATTR, DEFAULT_LEASE_ATTR, MAX_LEASE_ATTR] },
      { Options: [DESCRIPTION_ATTR] },
    ];
    run(() => {
      const groups = fieldToAttrs(model, [
        { default: ['path', 'config.{defaultLeaseTtl,maxLeaseTtl}'] },
        { Options: ['description'] },
      ]);
      assert.deepEqual(groups, expected, 'properly extracts fieldGroups with brace expansion');
    });
  });

  test('it extracts fieldGroups from model-fragment attributes with dot notation', function(assert) {
    const model = run(() => this.owner.lookup('service:store').createRecord('test-form-model'));
    const expected = [
      { default: [DEFAULT_LEASE_ATTR, PATH_ATTR, MAX_LEASE_ATTR] },
      { Options: [DESCRIPTION_ATTR] },
    ];
    run(() => {
      const groups = fieldToAttrs(model, [
        { default: ['config.defaultLeaseTtl', 'path', 'config.maxLeaseTtl'] },
        { Options: ['description'] },
      ]);
      assert.deepEqual(groups, expected, 'properly extracts fieldGroups with dot notation');
    });
  });
});
