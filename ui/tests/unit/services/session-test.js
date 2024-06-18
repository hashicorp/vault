import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Service | session', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const service = this.owner.lookup('service:session');
    assert.ok(service);
  });
});
