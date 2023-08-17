import { InvalidError } from '@ember-data/adapter/error';
import errorMessage from 'core/utils/error-message';
import { module, test } from 'qunit';

module('Unit | Utility | error-message', function () {
  // TODO: Replace this with your real tests.
  test('it works with regular Error', function (assert) {
    const mockError = new Error('permission denied');
    const result = errorMessage(mockError);
    assert.strictEqual(result, 'permission denied');
  });
  test('it works with AdapterError', function (assert) {
    const mockError = new InvalidError(['method not allowed', 'permission denied']);
    const result = errorMessage(mockError);
    assert.strictEqual(result, 'method not allowed, permission denied');
  });
  test('it works with default fallback message', function (assert) {
    const mockError = { foo: 'bar' };
    const result = errorMessage(mockError);
    assert.strictEqual(result, 'An error occurred, please try again');
  });
  test('it works with custom fallback message', function (assert) {
    const mockError = { foo: 'bar' };
    const result = errorMessage(mockError, 'Did you try turning it off and on again?');
    assert.strictEqual(result, 'Did you try turning it off and on again?');
  });
});
