import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Route | configuration/index', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const route = this.owner.lookup('route:configuration/index');
    assert.ok(route);
  });
});
