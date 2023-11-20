import { module, test } from 'qunit';
import { ensureTrailingSlash, getRelativePath, sanitizePath } from 'core/utils/sanitize-path';

module('Unit | Utility | sanitize-path', function () {
  test('it removes spaces and slashes from either side', function (assert) {
    assert.strictEqual(
      sanitizePath(' /foo/bar/baz/ '),
      'foo/bar/baz',
      'removes spaces and slashes on either side'
    );
    assert.strictEqual(sanitizePath('//foo/bar/baz/'), 'foo/bar/baz', 'removes more than one slash');
    assert.strictEqual(sanitizePath(undefined), '', 'handles falsey values');
  });

  test('#ensureTrailingSlash', function (assert) {
    assert.strictEqual(ensureTrailingSlash('foo/bar'), 'foo/bar/', 'adds trailing slash');
    assert.strictEqual(ensureTrailingSlash('baz/'), 'baz/', 'keeps trailing slash if there is one');
  });

  test('#getRelativePath', function (assert) {
    assert.strictEqual(getRelativePath('/', undefined), '', 'works with minimal inputs');
    assert.strictEqual(getRelativePath('/baz/bar/', undefined), 'baz/bar', 'sanitizes the output');
    assert.strictEqual(getRelativePath('recipes/cookies/choc-chip/', 'recipes/'), 'cookies/choc-chip');
    assert.strictEqual(getRelativePath('/recipes/cookies/choc-chip/', 'recipes/cookies'), 'choc-chip');
    assert.strictEqual(getRelativePath('/admin/bop/boop/admin_foo/baz/', 'admin'), 'bop/boop/admin_foo/baz');
  });
});
