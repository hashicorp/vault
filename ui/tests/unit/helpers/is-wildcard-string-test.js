import { isWildcardString } from 'vault/helpers/is-wildcard-string';
import { module, test } from 'qunit';

module('Unit | Helpers | is-wildcard-string', function() {
  test('it returns true if regular string with wildcard', function(assert) {
    let string = 'foom#*eep';
    let result = isWildcardString(string);
    assert.equal(result, true);
  });

  test('it returns false if no wildcard', function(assert) {
    let string = 'foo.bar';
    let result = isWildcardString(string);
    assert.equal(result, false);
  });

  test('it returns true if string with id as in searchSelect selected has wildcard', function(assert) {
    let string = { id: 'foo.bar*baz' };
    let result = isWildcardString(string);
    assert.equal(result, true);
  });

  test('it returns true if string object has name and no id', function(assert) {
    let string = { name: 'foo.bar*baz' };
    let result = isWildcardString(string);
    assert.equal(result, true);
  });

  test('it returns true if string object has name and id with at least one wildcard', function(assert) {
    let string = { id: '7*', name: 'seven' };
    let result = isWildcardString(string);
    assert.equal(result, true);
  });

  test('it returns true if string object has name and id with wildcard in name not id', function(assert) {
    let string = { id: '7', name: 'sev*n' };
    let result = isWildcardString(string);
    assert.equal(result, true);
  });
});
