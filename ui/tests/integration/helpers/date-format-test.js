import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { dateFormat } from '../../../helpers/date-format';

module('Integration | Helper | date-format', function(hooks) {
  setupRenderingTest(hooks);

  test('it is able to format a date object', function(assert) {
    let today = new Date();
    let result = dateFormat([today, 'YYYY']);
    assert.ok(typeof result === 'string');
    assert.ok(result !== 'Invalid Date', 'it is not an invalid date');
    assert.ok(Number(result) >= 2017);
  });

  test('it supports date timestamps', function(assert) {
    let today = new Date().getTime();
    let result = dateFormat([today, 'YYYY']);
    assert.ok(Number(result) >= 2017);
  });

  test('it supports date strings', function(assert) {
    let today = new Date().toString();
    let result = dateFormat([today, 'YYYY']);
    assert.ok(Number(result) >= 2017);
  });
});
