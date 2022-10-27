import { click, currentRouteName, currentURL, find, visit, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

module('Acceptance | policy', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  hooks.afterEach(function () {
    return logout.visit();
  });

  test('it redirects to list if navigating to root', async function (assert) {
    await visit('/vault/policies/acl/root');
    assert.strictEqual(
      currentURL(),
      '/vault/policies/acl',
      'navigation to root show redirects you to policy list'
    );
  });

  test('it does not show delete for default policy', async function (assert) {
    await visit('/vault/policies/acl/default');
    assert.notOk(find('[data-test-policy-delete]'), 'there is no delete button');
  });

  test('it navigates to edit and back to show when toggle is clicked', async function (assert) {
    await visit('/vault/policies/acl/default/show');
    await waitUntil(() => find('[data-test-policy-edit-toggle]'));
    await click('[data-test-policy-edit-toggle]');
    assert.strictEqual(currentURL(), '/vault/policies/acl/default/edit', 'toggle navigates to edit page');
    await click('[data-test-policy-edit-toggle]');
    assert.strictEqual(
      currentURL(),
      '/vault/policies/acl/default/show',
      'toggle navigates from edit to show'
    );
  });

  test('it redirects to acls on index navigation', async function (assert) {
    await visit('/vault/policies/acl');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
  });

  test('it redirects to acls with unknown policy type', async function (assert) {
    await visit('/vault/policies/foo');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
  });

  test('it redirects to acls with unknown policy type and policy name', async function (assert) {
    await visit('/vault/policies/foo/default');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');

    await visit('/vault/policies/foo/default/edit');
    assert.strictEqual(currentRouteName(), 'vault.cluster.policies.index');
    assert.strictEqual(currentURL(), '/vault/policies/acl');
  });
});
