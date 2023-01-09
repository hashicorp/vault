import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Model | pki/config', function (hooks) {
  setupTest(hooks);

  // Replace this with your real tests.
  test('it exists', function (assert) {
    const store = this.owner.lookup('service:store');
    const model = store.createRecord('pki/config', {});
    assert.ok(model);
  });
});
