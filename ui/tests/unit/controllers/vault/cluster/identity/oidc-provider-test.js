import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Controller | vault/cluster/identity/oidc-provider', function(hooks) {
  setupTest(hooks);

  // TODO: Replace this with your real tests.
  test('it exists', function(assert) {
    let controller = this.owner.lookup('controller:vault/cluster/identity/oidc-provider');
    assert.ok(controller);
  });
});
