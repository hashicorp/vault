import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Route | vault/cluster/oidc-provider', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let route = this.owner.lookup('route:vault/cluster/oidc-provider');
    assert.ok(route);
  });
});
