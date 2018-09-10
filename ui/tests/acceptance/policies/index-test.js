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
    let policies = page.policies();
    assert.equal(currentURL(), '/vault/policies/acl');
    assert.ok(policies.findByName('root'), 'root policy shown in the list');
    assert.ok(policies.findByName('default'), 'default policy shown in the list');
  });

  test('it navigates to show when clicking on the link', async function(assert) {
    await page.visit({ type: 'acl' });
    await page
      .policies()
      .findByName('default')
      .click();
    assert.equal(currentRouteName(), 'vault.cluster.policy.show');
    assert.equal(currentURL(), '/vault/policy/acl/default');
  });
});
