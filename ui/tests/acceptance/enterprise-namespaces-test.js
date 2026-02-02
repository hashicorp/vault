/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  click,
  settled,
  visit,
  fillIn,
  currentURL,
  triggerKeyEvent,
  find,
  waitFor,
  waitUntil,
} from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { runCmd, createNSFromPaths, deleteNSFromPaths } from 'vault/tests/helpers/commands';
import { login, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import localStorage from 'vault/lib/local-storage';

module('Acceptance | Enterprise | namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async () => {
    await login();
    // dismiss wizard
    localStorage.setItem('dismissed-wizards', ['namespace']);
  });

  test('it focuses the search input field when user toggles namespace picker', async function (assert) {
    await click(GENERAL.button('namespace-picker'));

    // Verify that the search input field is focused
    const searchInput = find(GENERAL.inputByAttr('Search namespaces'));
    assert.strictEqual(
      document.activeElement,
      searchInput,
      'The search input field is focused on component load'
    );
  });

  test('it navigates to the matching namespace when Enter is pressed', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespaces = ['beep/boop'];
    await createNSFromPaths(namespaces);

    await click(GENERAL.button('namespace-picker'));
    await click(GENERAL.button('Refresh list'));

    // Simulate typing into the search input
    await fillIn(GENERAL.inputByAttr('Search namespaces'), 'beep/boop');

    assert
      .dom(GENERAL.inputByAttr('Search namespaces'))
      .hasValue('beep/boop', 'The search input field has the correct value');

    // Simulate pressing Enter
    await triggerKeyEvent(GENERAL.inputByAttr('Search namespaces'), 'keydown', 'Enter');

    // Verify navigation to the matching namespace
    assert.strictEqual(
      this.owner.lookup('service:router').currentURL,
      '/vault/dashboard?namespace=beep%2Fboop',
      'Navigates to the correct namespace when Enter is pressed'
    );

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths(namespaces);
  });

  // This test originated from this PR: https://github.com/hashicorp/vault/pull/7186
  test('it clears namespaces when you log out', async function (assert) {
    // Test Setup
    const namespace = 'foo';
    await createNSFromPaths([namespace]);

    const token = await runCmd(`write -field=client_token auth/token/create policies=default`);
    await login(token);

    // Open the namespace picker
    await click(GENERAL.button('namespace-picker'));

    // Verify that the root namespace is selected by default
    assert.dom(GENERAL.button('root')).hasAttribute('aria-selected', 'true', 'root is selected by default');
    assert
      .dom(`${GENERAL.button('root')} svg${GENERAL.icon('check')}`)
      .exists('The root namespace has a check icon indicating it is selected');

    // Verify that the foo namespace does not exist in the namespace picker
    assert.dom(GENERAL.button(namespace)).doesNotExist('foo should not exist in the namespace picker');

    // Logout and log back into root
    await logout();
    await login();

    // Open the namespace picker & verify that the foo namespace does exist
    await click(GENERAL.button('namespace-picker'));
    assert.dom(GENERAL.button(namespace)).exists('foo should exist in the namespace picker');

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths([namespace]);
  });

  // This test originated from this PR: https://github.com/hashicorp/vault/pull/7186
  test('it displays namespaces whether you log in with a namespace prefixed with / or not', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespaces = ['beep/boop/bop'];
    await createNSFromPaths(namespaces);

    await click(GENERAL.button('namespace-picker'));
    await click(GENERAL.button('Refresh list'));

    // Login with a namespace prefixed with /
    await loginNs('/beep/boop');
    await settled();

    assert
      .dom(GENERAL.button('namespace-picker'))
      .hasText('boop', `shows the namespace 'boop' in the toggle component`);

    // Open the namespace picker
    await click(GENERAL.button('namespace-picker'));

    // Find the selected element with the check icon & ensure it exists
    assert
      .dom(`${GENERAL.button('beep/boop')} ${GENERAL.icon('check')}`)
      .exists('The selected namespace link exists with the check icon');

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths(namespaces);
  });

  test('it shows the regular namespace toolbar when not managed', async function (assert) {
    // This test is the opposite of the test in managed-namespace-test
    await logout();
    assert.strictEqual(currentURL(), '/vault/auth', 'Does not redirect');
    assert.dom(AUTH_FORM.managedNsRoot).doesNotExist('Managed namespace indicator does not exist');
    assert.dom(GENERAL.inputByAttr('namespace')).hasAttribute('placeholder', '/ (root)');
    await fillIn(GENERAL.inputByAttr('namespace'), '/foo/bar ');
    const encodedNamespace = encodeURIComponent('foo/bar');
    await waitUntil(() => currentURL() === `/vault/auth?namespace=${encodedNamespace}`);
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${encodedNamespace}`,
      'correctly sanitizes namespace'
    );
  });

  test('it should allow the user to delete a namespace', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespace = 'test-delete-me';
    await createNSFromPaths([namespace]);

    await visit('/vault/access/namespaces');

    // Verify that the namespace exists in the namespace picker
    await click(GENERAL.button('namespace-picker'));
    await waitFor(GENERAL.button('Refresh list'));
    await click(GENERAL.button('Refresh list'));
    await fillIn(GENERAL.inputByAttr('Search namespaces'), namespace);

    assert.dom(GENERAL.button(namespace)).exists('Namespace exists in the namespace picker');

    // Close the namespace picker
    await click(GENERAL.button('namespace-picker'));

    // Verify that the namespace exists in the manage namespaces page
    await fillIn(GENERAL.filterInputExplicit, namespace);
    await click(GENERAL.button('Search'));

    // Delete the namespace
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    assert.strictEqual(
      currentURL(),
      `/vault/access/namespaces?page=1&pageFilter=${namespace}`,
      'Should remain on the manage namespaces page after deletion'
    );
    // Verify that the namespace no longer exists on the namespace page
    assert.dom(GENERAL.emptyStateTitle).hasText('No namespaces yet', 'Namespace deletion successful');

    // Verify that the namespace does not exist in the namespace picker
    await click(GENERAL.button('namespace-picker'));
    await waitFor(GENERAL.button('Refresh list'));
    await click(GENERAL.button('Refresh list'));
    await fillIn(GENERAL.inputByAttr('Search namespaces'), namespace);
    assert
      .dom(GENERAL.button(namespace))
      .doesNotExist('Deleted namespace does not exist in the namespace picker');
  });

  test('it should show root in namespace picker when the user explicitly logs into root namespace', async function (assert) {
    // Explicitly set root as the namespace to login to
    await loginNs('root');

    assert
      .dom(GENERAL.button('namespace-picker'))
      .hasText('root', `shows the namespace 'root' in the toggle component`);

    // Verify user is in root namespace
    assert.true(
      this.owner.lookup('service:namespace').inRootNamespace,
      'Verifies that the user is in the root namespace'
    );
  });
});
