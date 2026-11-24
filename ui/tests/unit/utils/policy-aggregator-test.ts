/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { aggregatePolicies } from 'vault/utils/policy-aggregator';

module('Unit | Utility | policy-aggregator', function () {
  const stubs = {
    a: '# this one has a comment\npath "/foo" {\n    capabilities = ["create"]\n}\npath "/foo/bar" {\n    capabilities = ["create", "read", "update"]\n}',
    b: 'path "/foo" {\n    capabilities = ["delete"]\n}\npath "/bar" {\n    capabilities = ["patch", "delete"]\n}',
  };

  test('it aggregates policies without duplicating paths', function (assert) {
    const result = aggregatePolicies([stubs.a, stubs.b]);

    assert.deepEqual(
      result.policy,
      {
        '/foo': ['create', 'delete'],
        '/foo/bar': ['create', 'read', 'update'],
        '/bar': ['delete', 'patch'],
      },
      'policy object has correct structure with combined capabilities'
    );
  });

  test('it combines capabilities for duplicate paths', function (assert) {
    const result = aggregatePolicies([stubs.a, stubs.b]);

    assert.ok(result.policy['/foo']?.includes('create'), 'includes create from policy a');
    assert.ok(result.policy['/foo']?.includes('delete'), 'includes delete from policy b');
    assert.strictEqual(result.policy['/foo']?.length, 2, 'has exactly 2 capabilities');
  });

  test('it does not duplicate capabilities', function (assert) {
    const duplicatePolicy = 'path "/foo" {\n    capabilities = ["create"]\n}';
    const result = aggregatePolicies([stubs.a, duplicatePolicy]);

    const createCount = result.policy['/foo']?.filter((cap) => cap === 'create').length;
    assert.strictEqual(createCount, 1, 'create capability appears only once');
  });

  test('it generates correct policy string output', function (assert) {
    const result = aggregatePolicies([stubs.a, stubs.b]);

    assert.ok(result.policyString.includes('path "/foo"'), 'includes /foo path');
    assert.ok(result.policyString.includes('path "/foo/bar"'), 'includes /foo/bar path');
    assert.ok(result.policyString.includes('path "/bar"'), 'includes /bar path');
    assert.ok(result.policyString.includes('capabilities = '), 'includes capabilities keyword');
  });

  test('it sorts capabilities alphabetically', function (assert) {
    const result = aggregatePolicies([stubs.a, stubs.b]);

    assert.deepEqual(result.policy['/foo'], ['create', 'delete'], 'capabilities are sorted alphabetically');
    assert.deepEqual(result.policy['/bar'], ['delete', 'patch'], 'capabilities are sorted alphabetically');
  });

  test('it handles empty policy strings', function (assert) {
    const result = aggregatePolicies([]);

    assert.deepEqual(result.policy, {}, 'returns empty object for no policies');
    assert.strictEqual(result.policyString, '', 'returns empty string for no policies');
  });

  test('it handles policies with comments', function (assert) {
    const result = aggregatePolicies([stubs.a]);

    assert.ok(result.policy['/foo'], 'parses policy with comments correctly');
    assert.deepEqual(result.policy['/foo'], ['create'], 'extracts capabilities from policy with comments');
  });

  test('it handles single policy', function (assert) {
    const singlePolicy = 'path "/test" {\n    capabilities = ["read", "list"]\n}';
    const result = aggregatePolicies([singlePolicy]);

    assert.deepEqual(
      result.policy,
      {
        '/test': ['list', 'read'],
      },
      'correctly processes single policy'
    );
  });

  test('policy string format matches input format', function (assert) {
    const result = aggregatePolicies([stubs.a, stubs.b]);

    // Check that output format matches expected HCL format
    const lines = result.policyString.split('\n');
    const pathLines = lines.filter((line) => line.startsWith('path '));

    assert.strictEqual(pathLines.length, 3, 'has correct number of path declarations');

    // Verify format of each path block
    pathLines.forEach((line) => {
      assert.ok(line.match(/^path "\/[^"]*" \{$/), 'path line has correct format');
    });
  });
});
