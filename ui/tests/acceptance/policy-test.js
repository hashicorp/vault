import { currentURL, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';

module('Acceptance | policies', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  hooks.afterEach(function() {
    return authLogout();
  });

  test('it redirects to acls with unknown policy type', async function(assert) {
    await visit('/vault/policies/foo');
    assert.equal(currentRouteName(), 'vault.cluster.policies.index');
    assert.equal(currentURL(), '/vault/policies/acl');
  });
});
