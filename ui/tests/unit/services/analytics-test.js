import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Service | analytics', function (hooks) {
  setupTest(hooks);

  // TODO: Replace this with your real tests.
  test('it exists', function (assert) {
    const service = this.owner.lookup('service:analytics');
    assert.ok(service);
  });
});
