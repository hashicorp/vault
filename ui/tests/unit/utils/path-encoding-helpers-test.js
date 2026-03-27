/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { encodePath, normalizePath } from 'vault/utils/path-encoding-helpers';
import { module, test } from 'qunit';

module('Unit | Utility | path-encoding-helpers', function () {
  module('encodePath', function () {
    test('it encodes special characters but preserves slashes', function (assert) {
      // This is the key behavior for the URL encoding fix:
      // Slashes should be preserved, other special characters encoded
      assert.strictEqual(encodePath('simple'), 'simple', 'simple path unchanged');
      assert.strictEqual(encodePath('nested/path'), 'nested/path', 'slashes preserved');
      assert.strictEqual(encodePath('deep/nested/path'), 'deep/nested/path', 'multiple slashes preserved');
    });

    test('it encodes spaces and special characters', function (assert) {
      assert.strictEqual(encodePath('path with spaces'), 'path%20with%20spaces', 'spaces encoded');
      assert.strictEqual(encodePath('path?with=query'), 'path%3Fwith%3Dquery', 'query chars encoded');
      assert.strictEqual(encodePath('path#with#hash'), 'path%23with%23hash', 'hash encoded');
    });

    test('it handles nested paths with special characters', function (assert) {
      // Slashes preserved, other special chars in each segment encoded
      assert.strictEqual(
        encodePath('dir with space/secret'),
        'dir%20with%20space/secret',
        'space in directory encoded, slash preserved'
      );
      assert.strictEqual(
        encodePath('foo/bar?baz'),
        'foo/bar%3Fbaz',
        'question mark encoded, slash preserved'
      );
    });

    test('it encodes percent signs correctly', function (assert) {
      // This tests the edge case where a user wants a literal %2f in their path
      assert.strictEqual(encodePath('foo%2fbar'), 'foo%252fbar', 'percent in path is encoded');
      assert.strictEqual(
        encodePath('hello/foo%2fbar/world'),
        'hello/foo%252fbar/world',
        'percent encoded, slashes preserved'
      );
    });

    test('it handles empty and null paths', function (assert) {
      assert.strictEqual(encodePath(''), '', 'empty string returns empty');
      assert.strictEqual(encodePath(null), null, 'null returns null');
      assert.strictEqual(encodePath(undefined), undefined, 'undefined returns undefined');
    });

    test('it handles paths with only slashes', function (assert) {
      assert.strictEqual(encodePath('/'), '/', 'single slash preserved');
      assert.strictEqual(encodePath('//'), '//', 'double slash preserved');
      assert.strictEqual(encodePath('/foo/'), '/foo/', 'leading/trailing slashes preserved');
    });
  });

  module('normalizePath', function () {
    test('it decodes encoded characters', function (assert) {
      assert.strictEqual(normalizePath('simple'), 'simple', 'simple path unchanged');
      assert.strictEqual(normalizePath('path%20with%20spaces'), 'path with spaces', 'spaces decoded');
      assert.strictEqual(normalizePath('path%3Fwith%3Dquery'), 'path?with=query', 'query chars decoded');
    });

    test('it preserves slashes during decoding', function (assert) {
      assert.strictEqual(normalizePath('nested/path'), 'nested/path', 'slashes preserved');
      assert.strictEqual(
        normalizePath('dir%20with%20space/secret'),
        'dir with space/secret',
        'slashes preserved, encoded parts decoded'
      );
    });

    test('it handles empty paths', function (assert) {
      assert.strictEqual(normalizePath(''), '', 'empty string returns empty');
      assert.strictEqual(normalizePath(null), '', 'null returns empty string');
      assert.strictEqual(normalizePath(undefined), '', 'undefined returns empty string');
    });
  });
});
