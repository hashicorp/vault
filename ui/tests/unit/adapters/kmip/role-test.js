/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | kmip/role', function (hooks) {
  setupTest(hooks);

  // these are only some of the actual editable fields
  const editableFields = ['tlsTtl', 'operationAll', 'operationNone', 'operationGet', 'operationCreate'];
  const serializeTests = [
    [
      'operation_all is the only operation item present after serialization',
      {
        serialize() {
          return { operation_all: true, operation_get: true, operation_create: true, tls_ttl: '10s' };
        },
        record: {
          editableFields,
        },
      },
      {
        operation_all: true,
        tls_ttl: '10s',
      },
    ],
    [
      'serialize does not include nonOperationFields values if they are not set',
      {
        serialize() {
          return { operation_all: true, operation_get: true, operation_create: true };
        },
        record: {
          editableFields,
        },
      },
      {
        operation_all: true,
      },
    ],
    [
      'operation_none is the only operation item present after serialization',
      {
        serialize() {
          return { operation_none: true, operation_get: true, operation_add_attribute: true, tls_ttl: '10s' };
        },
        record: {
          editableFields,
        },
      },
      {
        operation_none: true,
        tls_ttl: '10s',
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
        record: {
          editableFields,
        },
      },
      {
        operation_get: true,
        operation_add_attribute: true,
        operation_destroy: true,
      },
    ],
  ];
  for (const testCase of serializeTests) {
    const [name, snapshotStub, expected] = testCase;
    test(`adapter serialize: ${name}`, function (assert) {
      const adapter = this.owner.lookup('adapter:kmip/role');
      const result = adapter.serialize(snapshotStub);
      assert.deepEqual(result, expected, 'output matches expected');
    });
  }
});
