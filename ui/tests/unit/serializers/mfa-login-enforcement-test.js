import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Serializer | mfa-login-enforcement', function (hooks) {
  setupTest(hooks);

  test('it should transform property names for hasMany relationships', function (assert) {
    const serverData = {
      name: 'foo',
      mfa_method_ids: ['1'],
      auth_method_types: ['userpass'],
      auth_method_accessors: ['auth_approle_17a552c6'],
      identity_entity_ids: ['2', '3'],
      identity_group_ids: ['4', '5', '6'],
    };
    const tranformedData = {
      name: 'foo',
      mfa_methods: ['1'],
      auth_method_types: ['userpass'],
      auth_method_accessors: ['auth_approle_17a552c6'],
      identity_entities: ['2', '3'],
      identity_groups: ['4', '5', '6'],
    };
    const mutableData = { ...serverData };
    const serializer = this.owner.lookup('serializer:mfa-login-enforcement');

    serializer.transformHasManyKeys(mutableData, 'model');
    assert.deepEqual(mutableData, tranformedData, 'hasMany property names are transformed for model');

    serializer.transformHasManyKeys(mutableData, 'server');
    assert.deepEqual(mutableData, serverData, 'hasMany property names are transformed for server');
  });
});
