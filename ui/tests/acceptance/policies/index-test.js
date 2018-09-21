import { currentURL, currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/policies/index';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | policies/acl', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it lists default and root acls', async function(assert) {
    await page.visit({ type: 'acl' });
    assert.equal(currentURL(), '/vault/policies/acl');
    assert.ok(page.findPolicyByName('root'), 'root policy shown in the list');
    assert.ok(page.findPolicyByName('default'), 'default policy shown in the list');
  });

  test('it navigates to show when clicking on the link', async function(assert) {
    await page.visit({ type: 'acl' });
    await page.findPolicyByName('default').click();
    assert.equal(currentRouteName(), 'vault.cluster.policy.show');
    assert.equal(currentURL(), '/vault/policy/acl/default');
  });
});
