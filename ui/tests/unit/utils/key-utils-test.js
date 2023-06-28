import {
  ancestorKeysForKey,
  keyIsFolder,
  keyPartsForKey,
  keyWithoutParentKey,
  parentKeyForKey,
} from 'vault/utils/key-utils';
import { module, test } from 'qunit';

module('Unit | Utility | key-utils', function () {
  test('keyIsFolder', function (assert) {
    let result = keyIsFolder('foo');
    assert.false(result, 'not folder');

    result = keyIsFolder('foo/');
    assert.true(result, 'is folder');

    result = keyIsFolder('foo/bar');
    assert.false(result, 'not folder');
  });
  test('keyPartsForKey', function (assert) {
    let result = keyPartsForKey('');
    assert.strictEqual(result, null, 'falsy value returns null');

    result = keyPartsForKey('foo');
    assert.strictEqual(result, null, 'returns null if not a folder');

    result = keyPartsForKey('foo/bar');
    assert.deepEqual(result, ['foo', 'bar'], 'returns parts of key');

    result = keyPartsForKey('foo/bar/');
    assert.deepEqual(result, ['foo', 'bar'], 'returns parts of key when ends in slash');
  });
  test('parentKeyForKey', function (assert) {
    let result = parentKeyForKey('my/very/nested/secret/path');
    assert.strictEqual(result, 'my/very/nested/secret/', 'returns parent path for key');

    result = parentKeyForKey('my/nested/secret/');
    assert.strictEqual(result, 'my/nested/', 'returns correct parents');

    result = parentKeyForKey('my-secret');
    assert.strictEqual(result, '', 'returns empty string when no parents');
  });
  test('keyWithoutParentKey', function (assert) {
    let result = keyWithoutParentKey('my/very/nested/secret/path');
    assert.strictEqual(result, 'path', 'returns key without parent key');

    result = keyWithoutParentKey('my-secret');
    assert.strictEqual(result, 'my-secret', 'returns path when no parent');

    result = keyWithoutParentKey('folder/');
    assert.strictEqual(result, 'folder/', 'returns path as-is when folder without parent');
  });
  test('ancestorKeysForKey', function (assert) {
    const expected = ['my/', 'my/very/', 'my/very/nested/', 'my/very/nested/secret/'];
    let result = ancestorKeysForKey('my/very/nested/secret/path');
    assert.deepEqual(result, expected, 'returns array of ancestor paths');

    result = ancestorKeysForKey('foobar');
    assert.deepEqual(result, [], 'returns empty array for root path');
  });
});
