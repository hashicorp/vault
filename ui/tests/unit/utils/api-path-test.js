import apiPath from 'vault/utils/api-path';
import { module, test } from 'qunit';

module('Unit | Util | api path', function () {
  test('it returns a function', function (assert) {
    const ret = apiPath`foo`;
    assert.strictEqual(typeof ret, 'function');
  });

  test('it iterpolates strings from passed context object', function (assert) {
    const ret = apiPath`foo/${'one'}/${'two'}`;
    const result = ret({ one: 1, two: 2 });

    assert.strictEqual(result, 'foo/1/2', 'returns the expected string');
  });

  test('it throws when the key is not found in the context', function (assert) {
    const ret = apiPath`foo/${'one'}/${'two'}`;
    assert.throws(() => {
      ret({ one: 1 });
    }, /Error: Assertion Failed: Expected 2 keys in apiPath context, only recieved one/);
  });
});
