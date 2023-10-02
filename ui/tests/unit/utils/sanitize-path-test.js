import { module, test } from 'qunit';
import { ensureTrailingSlash, sanitizePath } from 'core/utils/sanitize-path';

module('Unit | Utility | sanitize-path', function () {
  test('it removes spaces and slashes from either side', function (assert) {
    assert.strictEqual(
      sanitizePath(' /foo/bar/baz/ '),
      'foo/bar/baz',
      'removes spaces and slashs on either side'
    );
    assert.strictEqual(sanitizePath('//foo/bar/baz/'), 'foo/bar/baz', 'removes more than one slash');
  });

  test('#ensureTrailingSlash', function (assert) {
    assert.strictEqual(ensureTrailingSlash('foo/bar'), 'foo/bar/', 'adds trailing slash');
    assert.strictEqual(ensureTrailingSlash('baz/'), 'baz/', 'keeps trailing slash if there is one');
  });
});
