import { filterWildcard } from 'vault/helpers/filter-wildcard';
import { module, test } from 'qunit';

module('Unit | Helpers | filter-wildcard', function() {
  test('it returns a count if array contains a wildcard', function(assert) {
    let string = { id: 'foo*' };
    let array = ['foobar', 'foozar', 'boo', 'oof'];
    let result = filterWildcard([string, array]);
    assert.equal(2, result);
  });

  test('it returns zero if no wildcard is string', function(assert) {
    let string = { id: 'foo#' };
    let array = ['foobar', 'foozar', 'boo', 'oof'];
    let result = filterWildcard([string, array]);
    assert.equal(0, result);
  });

  test('it escapes function and does not error if no id is in string', function(assert) {
    let string = '*bar*';
    let array = ['foobar', 'foozar', 'boobarboo', 'oof'];
    let result = filterWildcard([string, array]);
    assert.equal(2, result);
  });
});
