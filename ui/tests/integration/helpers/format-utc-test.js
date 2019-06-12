import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { formatUtc } from '../../../helpers/format-utc';

module('Integration | Helper | format-utc', function(hooks) {
  setupRenderingTest(hooks);

  test('it formats a UTC date string and maintains the timezone', function(assert) {
    let expected = 'Apr 01 2019, 00:00';
    let dateTime = '2019-04-01T00:00:00Z';
    let result = formatUtc([dateTime, '%b %d %Y, %H:%M']);
    assert.equal(result, expected, 'it displays the date in UTC');
  });
});
