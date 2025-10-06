/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

module('Unit | Helper | engineDisplayData', function () {
  test('it returns correct display data for a known engine type', function (assert) {
    const awsData = engineDisplayData('aws');
    const expected = ALL_ENGINES.find((e) => e.type === 'aws');
    assert.propEqual(awsData, expected, 'Returns correct display data for aws');
  });

  test('it returns correct display data for an ent only engine', function (assert) {
    const kmipData = engineDisplayData('kmip');
    assert.true(kmipData.requiresEnterprise, 'KMIP requires enterprise');
    assert.strictEqual(kmipData.displayName, 'KMIP', 'KMIP displayName is correct');
  });

  test('it returns fallback display data for unknown engine type', function (assert) {
    const { displayName, type, mountCategory, glyph, isOldEngine } = engineDisplayData('not-an-engine');
    assert.strictEqual(displayName, 'not-an-engine', 'it returns passed type as fallback displayName');
    assert.strictEqual(type, 'unknown', 'it returns "unknown"" as fallback type');
    assert.propEqual(mountCategory, ['secret', 'auth'], 'mountCategory is correct');
    assert.strictEqual(glyph, 'lock', 'default glyph is a lock');
    assert.true(isOldEngine, 'isOldEngine is true');
  });

  test('it returns fallback display data for empty string', function (assert) {
    const { displayName, type, mountCategory, glyph, isOldEngine } = engineDisplayData('');
    assert.strictEqual(displayName, 'Unknown plugin', 'it returns fallback displayName for empty string');
    assert.strictEqual(type, 'unknown', 'it returns fallback type for empty string');
    assert.propEqual(mountCategory, ['secret', 'auth'], 'mountCategory is correct');
    assert.strictEqual(glyph, 'lock', 'default glyph is a lock');
    assert.true(isOldEngine, 'isOldEngine is true');
  });

  test('it returns fallback display data for undefined', function (assert) {
    const { displayName, type, mountCategory, glyph, isOldEngine } = engineDisplayData(undefined);
    assert.strictEqual(displayName, 'Unknown plugin', 'it returns fallback displayName for undefined');
    assert.strictEqual(type, 'unknown', 'it returns fallback type for undefined');
    assert.propEqual(mountCategory, ['secret', 'auth'], 'mountCategory is correct');
    assert.strictEqual(glyph, 'lock', 'default glyph is a lock');
    assert.true(isOldEngine, 'isOldEngine is true');
  });

  test('it returns fallback display data for null', function (assert) {
    const { displayName, type, mountCategory, glyph, isOldEngine } = engineDisplayData(null);
    assert.strictEqual(displayName, 'Unknown plugin', 'it returns fallback displayName for null');
    assert.strictEqual(type, 'unknown', 'it returns fallback type for null');
    assert.propEqual(mountCategory, ['secret', 'auth'], 'mountCategory is correct');
    assert.strictEqual(glyph, 'lock', 'default glyph is a lock');
    assert.true(isOldEngine, 'isOldEngine is true');
  });
});
