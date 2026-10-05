/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test, skip } from 'qunit';
import { FEATURE_VERSIONS, getFeatureVersion } from 'vault/utils/feature-versions';

module('Unit | Utility | feature-versions', function () {
  test('FEATURE_VERSIONS is an array', function (assert) {
    assert.ok(Array.isArray(FEATURE_VERSIONS), 'FEATURE_VERSIONS is an array');
  });

  // enable this test once there are features in the map to test against
  skip('every entry in FEATURE_VERSIONS has required fields', function (assert) {
    for (const feature of FEATURE_VERSIONS) {
      assert.ok(feature.key, `entry has a key: ${feature.key}`);
      assert.ok(feature.name, `entry "${feature.key}" has a name`);
      assert.ok(feature.version, `entry "${feature.key}" has a version`);
      assert.strictEqual(typeof feature.key, 'string', `key is a string`);
      assert.strictEqual(typeof feature.name, 'string', `name is a string`);
      assert.strictEqual(typeof feature.version, 'string', `version is a string`);
    }
  });

  test('getFeatureVersion returns the matching entry by key', function (assert) {
    // Temporarily add a known entry for testing purposes
    FEATURE_VERSIONS.push({ key: 'test-feature', name: 'Test Feature', version: '2.0.0' });

    const result = getFeatureVersion('test-feature');
    assert.ok(result, 'returns an entry for a known key');
    assert.strictEqual(result?.key, 'test-feature', 'returned entry has correct key');
    assert.strictEqual(result?.name, 'Test Feature', 'returned entry has correct name');
    assert.strictEqual(result?.version, '2.0.0', 'returned entry has correct version');

    // Clean up
    const idx = FEATURE_VERSIONS.findIndex((f) => f.key === 'test-feature');
    FEATURE_VERSIONS.splice(idx, 1);
  });

  test('getFeatureVersion returns undefined for an unknown key', function (assert) {
    const result = getFeatureVersion('this-key-does-not-exist');
    assert.strictEqual(result, undefined, 'returns undefined for unknown key');
  });

  test('all keys in FEATURE_VERSIONS are unique', function (assert) {
    const keys = FEATURE_VERSIONS.map((f) => f.key);
    const uniqueKeys = new Set(keys);
    assert.strictEqual(keys.length, uniqueKeys.size, 'no duplicate keys in FEATURE_VERSIONS');
  });
});
