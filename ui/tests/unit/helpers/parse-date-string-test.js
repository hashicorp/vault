import { parseDateString } from 'vault/helpers/parse-date-string';
import { module, test } from 'qunit';
import { compareAsc } from 'date-fns';

module('Unit | Helpers | parse-date-string', function() {
  test('it returns the first of the month when date like MM-yyyy passed in', function(assert) {
    let expected = new Date(2020, 3, 1);
    let result = parseDateString('04-2020');
    assert.equal(compareAsc(expected, result), 0);
  });

  test('it can handle a date format like MM/yyyy', function(assert) {
    let expected = new Date(2020, 11, 1);
    let result = parseDateString('12/2020', '/');
    assert.equal(compareAsc(expected, result), 0);
  });

  test('it throws an error with passed separator if bad format', function(assert) {
    let result;
    try {
      result = parseDateString('01-12-2020');
    } catch (e) {
      result = e.message;
    }
    assert.equal('Please use format MM-yyyy', result);
  });

  test('it throws an error with wrong separator', function(assert) {
    let result;
    try {
      result = parseDateString('12/2020', '.');
    } catch (e) {
      result = e.message;
    }
    assert.equal('Please use format MM.yyyy', result);
  });

  test('it throws an error if month is invalid', function(assert) {
    let result;
    try {
      result = parseDateString('13-2020');
    } catch (e) {
      result = e.message;
    }
    assert.equal('Not a valid month value', result);
  });
});
