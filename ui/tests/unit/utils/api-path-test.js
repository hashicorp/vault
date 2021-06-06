import apiPath from 'vault/utils/api-path';
import { module, test } from 'qunit';

module('Unit | Util | api path', function() {
  test('it returns a function', function(assert) {
    let ret = apiPath`foo`;
    assert.ok(typeof ret === 'function');
  });

  test('it iterpolates strings from passed context object', function(assert) {
    let ret = apiPath`foo/${'one'}/${'two'}`;
    let result = ret({ one: 1, two: 2 });

    assert.equal(result, 'foo/1/2', 'returns the expected string');
  });

  test('it throws when the key is not found in the context', function(assert) {
    let ret = apiPath`foo/${'one'}/${'two'}`;
    assert.throws(() => {
      ret({ one: 1 });
    }, /Error: Assertion Failed: Expected 2 keys in apiPath context, only recieved one/);
  });
});
