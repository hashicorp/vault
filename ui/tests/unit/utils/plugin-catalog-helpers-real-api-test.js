/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { addVersionsToEngines } from 'vault/utils/plugin-catalog-helpers';

module('Unit | Utility | plugin-catalog-helpers | real API response', function () {
  test('addVersionsToEngines handles real plugin catalog response correctly', function (assert) {
    const staticEngines = [
      { type: 'aws', displayName: 'AWS', mountCategory: ['secret'], glyph: 'aws-color' },
      { type: 'database', displayName: 'Databases', mountCategory: ['secret'], glyph: 'database' },
      { type: 'kv', displayName: 'KV', mountCategory: ['secret'], glyph: 'key-values' },
    ];

    // Using data from your actual plugin catalog response
    const secretEnginesList = [
      'ad',
      'alicloud',
      'aws',
      'azure',
      'consul',
      'gcp',
      'gcpkms',
      'keymgmt',
      'kmip',
      'kubernetes',
      'kv',
      'ldap',
      'mongodbatlas',
      'nomad',
      'openldap',
      'pki',
      'rabbitmq',
      'ssh',
      'terraform',
      'totp',
      'transform',
      'transit',
    ];

    const secretEnginesDetailed = [
      {
        builtin: true,
        deprecation_status: 'supported',
        name: 'aws',
        type: 'secret',
        version: 'v1.21.0+builtin.vault',
      },
      {
        builtin: true,
        deprecation_status: 'supported',
        name: 'kv',
        type: 'secret',
        version: 'v0.24.1+builtin',
      },
    ];

    const databasePluginsList = [
      'cassandra-database-plugin',
      'couchbase-database-plugin',
      'elasticsearch-database-plugin',
      'hana-database-plugin',
      'influxdb-database-plugin',
      'mongodb-database-plugin',
      'mongodbatlas-database-plugin',
      'mssql-database-plugin',
      'mysql-aurora-database-plugin',
      'mysql-database-plugin',
      'mysql-legacy-database-plugin',
      'mysql-rds-database-plugin',
      'postgresql-database-plugin',
      'redis-database-plugin',
      'redis-elasticache-database-plugin',
      'redshift-database-plugin',
      'snowflake-database-plugin',
    ];

    const databasePluginsDetailed = [
      {
        builtin: true,
        deprecation_status: 'supported',
        name: 'cassandra-database-plugin',
        type: 'database',
        version: 'v1.21.0+builtin.vault',
      },
      {
        builtin: true,
        deprecation_status: 'supported',
        name: 'mysql-database-plugin',
        type: 'database',
        version: 'v1.21.0+builtin.vault',
      },
      {
        builtin: true,
        deprecation_status: 'supported',
        name: 'postgresql-database-plugin',
        type: 'database',
        version: 'v1.21.0+builtin.vault',
      },
    ];

    const result = addVersionsToEngines(
      staticEngines,
      secretEnginesList,
      secretEnginesDetailed,
      databasePluginsList,
      databasePluginsDetailed
    );

    const realEngines = result.filter(
      (e) => !e.type.startsWith('demo-') && !e.type.startsWith('example-') && !e.type.startsWith('test-')
    );

    // Find the database engine
    const databaseEngine = realEngines.find((e) => e.type === 'database');
    assert.ok(databaseEngine, 'Database engine should be present');
    assert.true(
      databaseEngine.isAvailable,
      'Database engine should be available when database plugins are present in catalog'
    );
    assert.strictEqual(databaseEngine.version, 'v1.21.0', 'Database should have cleaned version');
    assert.false(databaseEngine.isExternalPlugin, 'Database should not be marked as external');
    assert.true(databaseEngine.builtin, 'Database should be marked as builtin');

    // Verify AWS and KV work normally too
    const awsEngine = realEngines.find((e) => e.type === 'aws');
    assert.ok(awsEngine, 'AWS engine should be present');
    assert.true(awsEngine.isAvailable, 'AWS should be available');

    const kvEngine = realEngines.find((e) => e.type === 'kv');
    assert.ok(kvEngine, 'KV engine should be present');
    assert.true(kvEngine.isAvailable, 'KV should be available');
  });
});
