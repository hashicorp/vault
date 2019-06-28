import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | kmip/role', function(hooks) {
  setupTest(hooks);

  let serializeTests = [
    [
      'operation_all is the only item present after serialization',
      {
        serialize() {
          return { operation_all: true, operation_get: true, operation_create: true };
        },
      },
      {
        operation_all: true,
      },
    ],
    [
      'operation_none is the only item present after serialization',
      {
        serialize() {
          return { operation_none: true, operation_get: true, operation_add_attribute: true };
        },
      },
      {
        operation_none: true,
      },
    ],
    [
      'operation_all and operation_none are removed if not truthy',
      {
        serialize() {
          return {
            operation_all: false,
            operation_none: false,
            operation_get: true,
            operation_add_attribute: true,
            operation_destroy: true,
          };
        },
      },
      {
        operation_get: true,
        operation_add_attribute: true,
        operation_destroy: true,
      },
    ],
  ];
  for (let testCase of serializeTests) {
    let [name, snapshotStub, expected] = testCase;

    test(`adapter serialize: ${name}`, function(assert) {
      let adapter = this.owner.lookup('adapter:kmip/role');
      let result = adapter.serialize(snapshotStub);
      assert.deepEqual(result, expected, 'output matches expected');
    });
  }
});
