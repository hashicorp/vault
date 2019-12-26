import { dotToDash } from 'vault/helpers/dot-to-dash';
import { module, test } from 'qunit';

module('Unit | Helpers | dot-to-dash', function() {
  test('it returns a string unchanged if there are not .s', function(assert) {
    let string = 'foo';
    let result = dotToDash([string]);
    assert.equal(string, result);
  });

  test('it replaces a single . with -', function(assert) {
    let string = 'foo.bar';
    let result = dotToDash([string]);
    assert.equal(result, 'foo-bar');
  });

  test('it replaces multiple . with -', function(assert) {
    let string = 'foo.bar.baz';
    let result = dotToDash([string]);
    assert.equal(result, 'foo-bar-baz');
  });
});
