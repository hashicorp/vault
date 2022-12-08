import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Adapter | pki/certificate/generate', function (hooks) {
  setupTest(hooks);

  // Replace this with your real tests.
  test('it exists', function (assert) {
    const adapter = this.owner.lookup('adapter:pki/certificate/generate');
    assert.ok(adapter);
  });
});
