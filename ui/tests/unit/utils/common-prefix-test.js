import commonPrefix from 'core/utils/common-prefix';
import { module, test } from 'qunit';

module('Unit | Util | common prefix', function() {
  test('it returns empty string if called with no args or an empty array', function(assert) {
    let returned = commonPrefix();
    assert.equal(returned, '', 'returns an empty string');
    returned = commonPrefix([]);
    assert.equal(returned, '', 'returns an empty string for an empty array');
  });

  test('it returns empty string if there are no common prefixes', function(assert) {
    let secrets = ['asecret', 'secret2', 'secret3'].map(s => ({ id: s }));
    let returned = commonPrefix(secrets);
    assert.equal(returned, '', 'returns an empty string');
  });

  test('it returns the longest prefix', function(assert) {
    let secrets = ['secret1', 'secret2', 'secret3'].map(s => ({ id: s }));
    let returned = commonPrefix(secrets);
    assert.equal(returned, 'secret', 'finds secret prefix');
    let greetings = ['hello-there', 'hello-hi', 'hello-howdy'].map(s => ({ id: s }));
    returned = commonPrefix(greetings);
    assert.equal(returned, 'hello-', 'finds hello- prefix');
  });

  test('it can compare an attribute that is not "id" to calculate the longest prefix', function(assert) {
    let secrets = ['secret1', 'secret2', 'secret3'].map(s => ({ name: s }));
    let returned = commonPrefix(secrets, 'name');
    assert.equal(returned, 'secret', 'finds secret prefix from name attribute');
  });
});
