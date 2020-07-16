import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { dateFromNow } from '../../../helpers/date-from-now';

module('Integration | Helper | date-from-now', function(hooks) {
  setupRenderingTest(hooks);

  test('it accepts a number', function(assert) {
    let result = dateFromNow([1481022124443]);
    assert.ok(result.includes('years'));
  });

  test('it accepts a Date', function(assert) {
    let result = dateFromNow([new Date('2006-06-06')]);
    assert.ok(result.includes('years'));
  });

  test('fails gracefully with strings', function(assert) {
    let result = dateFromNow(['foo']);
    assert.equal(result, '');
  });

  test('you can include a suffix', function(assert) {
    let result = dateFromNow([1481022124443], { addSuffix: true });
    assert.ok(result.includes(' ago'));
  });
});
