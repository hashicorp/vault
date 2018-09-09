import { currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/policy/edit';

module('Acceptance | policy/acl/:name/edit', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it redirects to list if navigating to root', function(assert) {
    page.visit({ type: 'acl', name: 'root' });
    assert.equal(currentURL(), '/vault/policies/acl', 'navigation to root show redirects you to policy list');
  });

  test('it does not show delete for default policy', function(assert) {
    page.visit({ type: 'acl', name: 'default' });
    assert.notOk(page.deleteIsPresent, 'there is no delete button');
  });

  test('it navigates to show when the toggle is clicked', function(assert) {
    page.visit({ type: 'acl', name: 'default' }).toggleEdit();
    assert.equal(currentURL(), '/vault/policy/acl/default', 'toggle navigates from edit to show');
  });
});
