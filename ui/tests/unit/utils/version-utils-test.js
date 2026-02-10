/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  areVersionsEqual,
  cleanVersion,
  compareVersions,
  getHighestVersion,
  isValidVersion,
  isVersionGreater,
  parseVersion,
  sortVersions,
} from 'vault/utils/version-utils';

module('Unit | Utility | version-utils', function () {
  test('cleanVersion removes prefixes and suffixes correctly', function (assert) {
    assert.strictEqual(cleanVersion('v1.2.3'), '1.2.3', 'removes v prefix');
    assert.strictEqual(cleanVersion('1.2.3+ent'), '1.2.3', 'removes +ent suffix');
    assert.strictEqual(cleanVersion('v1.2.3+builtin'), '1.2.3', 'removes v prefix and +builtin suffix');
    assert.strictEqual(cleanVersion('v1.2.3-beta1+ent'), '1.2.3', 'removes v prefix and -beta1+ent suffix');
    assert.strictEqual(cleanVersion('1.2.3'), '1.2.3', 'leaves clean version unchanged');
  });

  test('parseVersion converts version strings to numeric arrays', function (assert) {
    assert.deepEqual(parseVersion('1.2.3'), [1, 2, 3], 'parses basic version');
    assert.deepEqual(parseVersion('v1.0.0+ent'), [1, 0, 0], 'parses version with prefix and suffix');
    assert.deepEqual(parseVersion('1.2'), [1, 2], 'parses two-part version');
    assert.deepEqual(parseVersion('1.2.3.4'), [1, 2, 3, 4], 'parses four-part version');
    assert.deepEqual(parseVersion('1.0.x'), [1, 0, 0], 'handles non-numeric parts as 0');
  });

  test('compareVersions works correctly', function (assert) {
    // Equal versions
    assert.strictEqual(compareVersions('1.2.3', '1.2.3'), 0, '1.2.3 equals 1.2.3');
    assert.strictEqual(compareVersions('v1.2.3+ent', '1.2.3'), 0, 'v1.2.3+ent equals 1.2.3');

    // First version greater
    assert.ok(compareVersions('1.2.4', '1.2.3') > 0, '1.2.4 > 1.2.3');
    assert.ok(compareVersions('1.3.0', '1.2.9') > 0, '1.3.0 > 1.2.9');
    assert.ok(compareVersions('2.0.0', '1.9.9') > 0, '2.0.0 > 1.9.9');

    // Second version greater
    assert.ok(compareVersions('1.2.3', '1.2.4') < 0, '1.2.3 < 1.2.4');
    assert.ok(compareVersions('1.2.9', '1.3.0') < 0, '1.2.9 < 1.3.0');
    assert.ok(compareVersions('1.9.9', '2.0.0') < 0, '1.9.9 < 2.0.0');

    // Different lengths
    assert.ok(compareVersions('1.2.3', '1.2') > 0, '1.2.3 > 1.2');
    assert.ok(compareVersions('1.2', '1.2.1') < 0, '1.2 < 1.2.1');
  });

  test('sortVersions sorts correctly', function (assert) {
    const versions = ['v1.0.0+ent', 'v0.18.0+ent', 'v0.19.0+ent', 'v1.1.0+ent'];

    // Ascending order (default)
    const ascending = sortVersions(versions);
    assert.deepEqual(
      ascending,
      ['v0.18.0+ent', 'v0.19.0+ent', 'v1.0.0+ent', 'v1.1.0+ent'],
      'sorts ascending'
    );

    // Descending order
    const descending = sortVersions(versions, true);
    assert.deepEqual(
      descending,
      ['v1.1.0+ent', 'v1.0.0+ent', 'v0.19.0+ent', 'v0.18.0+ent'],
      'sorts descending'
    );

    // Original array unchanged
    assert.deepEqual(
      versions,
      ['v1.0.0+ent', 'v0.18.0+ent', 'v0.19.0+ent', 'v1.1.0+ent'],
      'original array unchanged'
    );
  });

  test('getHighestVersion returns the latest version', function (assert) {
    const versions = ['v1.0.0+ent', 'v0.18.0+ent', 'v0.19.0+ent', 'v1.1.0+ent'];
    assert.strictEqual(getHighestVersion(versions), 'v1.1.0+ent', 'returns highest version');
    assert.strictEqual(getHighestVersion([]), null, 'returns null for empty array');
    assert.strictEqual(getHighestVersion(['v1.0.0']), 'v1.0.0', 'returns single version');
  });

  test('isVersionGreater compares versions correctly', function (assert) {
    assert.true(isVersionGreater('1.2.4', '1.2.3'), '1.2.4 > 1.2.3');
    assert.true(isVersionGreater('v1.0.0+ent', '0.9.0'), 'v1.0.0+ent > 0.9.0');
    assert.false(isVersionGreater('1.2.3', '1.2.4'), '1.2.3 not > 1.2.4');
    assert.false(isVersionGreater('1.2.3', '1.2.3'), '1.2.3 not > 1.2.3');
  });

  test('areVersionsEqual compares versions correctly', function (assert) {
    assert.true(areVersionsEqual('1.2.3', '1.2.3'), '1.2.3 equals 1.2.3');
    assert.true(areVersionsEqual('v1.2.3+ent', '1.2.3'), 'v1.2.3+ent equals 1.2.3');
    assert.false(areVersionsEqual('1.2.3', '1.2.4'), '1.2.3 not equal 1.2.4');
  });

  test('edge cases are handled correctly', function (assert) {
    // Empty strings
    assert.strictEqual(compareVersions('', ''), 0, 'empty strings are equal');
    assert.strictEqual(cleanVersion(''), '', 'empty string returns empty');

    // Only prefixes/suffixes
    assert.strictEqual(cleanVersion('v'), '', 'only prefix returns empty');
    assert.strictEqual(cleanVersion('+ent'), '', 'only suffix returns empty');
  });

  test('isValidVersion validates version strings correctly', function (assert) {
    // Valid versions
    assert.true(isValidVersion('0.17'), 'Basic semver is valid');
    assert.true(isValidVersion('0.17.0'), 'Full semver is valid');
    assert.true(isValidVersion('v0.17.1'), 'Version with v prefix is valid');
    assert.true(isValidVersion('1.2.3+ent'), 'Version with build metadata is valid');
    assert.true(isValidVersion('2.0.0-beta'), 'Version with pre-release is valid');

    // Invalid versions
    assert.false(isValidVersion(''), 'Empty string is invalid');
    assert.false(isValidVersion('   '), 'Whitespace only is invalid');
    assert.false(isValidVersion('null'), 'String "null" is invalid');
    assert.false(isValidVersion('invalid'), 'Non-numeric string is invalid');
    assert.false(isValidVersion(null), 'null is invalid');
    assert.false(isValidVersion(undefined), 'undefined is invalid');
  });
});
