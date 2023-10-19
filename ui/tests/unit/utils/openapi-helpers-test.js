import { module, test } from 'qunit';
import { _getPathParam, pathToHelpUrlSegment } from 'vault/utils/openapi-helpers';

module('Unit | Utility | OpenAPI helper utils', function () {
  test(`pathToHelpUrlSegment`, function (assert) {
    assert.expect(5);
    [
      { path: '/auth/{username}', result: '/auth/example' },
      { path: '{username}/foo', result: 'example/foo' },
      { path: 'foo/{username}/bar', result: 'foo/example/bar' },
      { path: '', result: '' },
      { path: undefined, result: '' },
    ].forEach((test) => {
      assert.strictEqual(pathToHelpUrlSegment(test.path), test.result, `translates ${test.path}`);
    });
  });

  test(`_getPathParam`, function (assert) {
    assert.expect(7);
    [
      { path: '/auth/{username}', result: 'username' },
      { path: '{unicorn}/foo', result: 'unicorn' },
      { path: 'foo/{bigfoot}/bar', result: 'bigfoot' },
      { path: '{alphabet}/bowl/{soup}', result: 'alphabet' },
      { path: 'no/params', result: false },
      { path: '', result: false },
      { path: undefined, result: false },
    ].forEach((test) => {
      assert.strictEqual(_getPathParam(test.path), test.result, `returns first param for ${test.path}`);
    });
  });
});
