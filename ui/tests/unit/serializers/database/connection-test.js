/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

module('Unit | Serializer | database/connection', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    this.store = this.owner.lookup('service:store');
  });
  test('it should serialize only keys that are valid for the database type (elasticsearch)', function (assert) {
    const backend = `db-serializer-test-${this.uid}`;
    const name = `elastic-test-${this.uid}`;
    const record = this.store.createRecord('database/connection', {
      plugin_name: 'elasticsearch-database-plugin',
      backend,
      name,
      allowed_roles: ['readonly'],
      connection_url: 'http://localhost:9200',
      url: 'http://localhost:9200',
      username: 'elastic',
      password: 'changeme',
      tls_ca: 'some-value',
      ca_cert: undefined, // does not send undefined values
    });
    const expectedResult = {
      plugin_name: 'elasticsearch-database-plugin',
      backend,
      name,
      verify_connection: true,
      allowed_roles: ['readonly'],
      url: 'http://localhost:9200',
      username: 'elastic',
      password: 'changeme',
      insecure: false,
    };

    const serializedRecord = record.serialize();
    assert.deepEqual(
      serializedRecord,
      expectedResult,
      'invalid elasticsearch options were not added to the payload'
    );
  });

  test('it should normalize values for the database type (elasticsearch)', function (assert) {
    const serializer = this.owner.lookup('serializer:database/connection');
    const normalized = serializer.normalizeSecrets({
      request_id: 'request-id',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        allowed_roles: ['readonly'],
        connection_details: {
          backend: 'database',
          insecure: false,
          url: 'https://localhost:9200',
          username: 'root',
        },
        password_policy: '',
        plugin_name: 'elasticsearch-database-plugin',
        plugin_version: '',
        root_credentials_rotate_statements: [],
      },
      wrap_info: null,
      warnings: null,
      auth: null,
      mount_type: 'database',
      backend: 'database',
      id: 'elastic-test',
    });
    const expectedResult = {
      allowed_roles: ['readonly'],
      backend: 'database',
      connection_details: {
        backend: 'database',
        insecure: false,
        url: 'https://localhost:9200',
        username: 'root',
      },
      id: 'elastic-test',
      insecure: false,
      name: 'elastic-test',
      password_policy: '',
      plugin_name: 'elasticsearch-database-plugin',
      plugin_version: '',
      root_credentials_rotate_statements: [],
      root_rotation_statements: [],
      url: 'https://localhost:9200',
      username: 'root',
    };
    assert.deepEqual(normalized, expectedResult, `Normalizes and flattens database response`);
  });
});
