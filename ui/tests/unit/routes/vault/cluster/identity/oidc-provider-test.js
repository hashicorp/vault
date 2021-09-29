import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Route | vault/cluster/identity/oidc-provider', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let route = this.owner.lookup('route:vault/cluster/identity/oidc-provider');
    assert.ok(route);
  });
});
