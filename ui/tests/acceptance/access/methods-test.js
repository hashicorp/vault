import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/methods';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | /access/', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authPage.login();
  });

  test('it navigates', async function(assert) {
    await page.visit();
    assert.ok(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.ok(page.navLinks.objectAt(0).isActive, 'the first link is active');
    assert.equal(page.navLinks.objectAt(0).text, 'Auth Methods');
  });
});
