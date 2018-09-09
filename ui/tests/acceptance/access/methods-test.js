import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/access/methods';

module('Acceptance | /access/', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it navigates', function(assert) {
    page.visit();
    assert.ok(currentRouteName(), 'vault.cluster.access.methods', 'navigates to the correct route');
    assert.ok(page.navLinks(0).isActive, 'the first link is active');
    assert.equal(page.navLinks(0).text, 'Auth Methods');
  });
});
