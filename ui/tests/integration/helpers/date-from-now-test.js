import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { dateFromNow } from '../../../helpers/date-from-now';

module('Integration | Helper | date-from-now', function(hooks) {
  setupRenderingTest(hooks);

  test('it works', function(assert) {
    let result = dateFromNow([1481022124443]);
    assert.ok(typeof result === 'string', 'it is a string');
  });

  test('you can include a suffix', function(assert) {
    let result = dateFromNow([1481022124443], { addSuffix: true });
    assert.ok(result.includes(' ago'));
  });
});
