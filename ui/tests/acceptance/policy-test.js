import { currentURL, currentRouteName, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

module('Acceptance | policies', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  hooks.afterEach(function() {
    return logout.visit();
  });

  test('it redirects to acls with unknown policy type', async function(assert) {
    await visit('/vault/policies/foo');
    assert.equal(currentRouteName(), 'vault.cluster.policies.index');
    assert.equal(currentURL(), '/vault/policies/acl');
  });
});
