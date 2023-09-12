import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Route | vault/cluster/dashboard', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const route = this.owner.lookup('route:vault/cluster/dashboard');
    assert.ok(route);
  });
});
