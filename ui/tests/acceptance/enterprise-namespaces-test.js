/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  click,
  settled,
  visit,
  fillIn,
  currentURL,
  findAll,
  triggerKeyEvent,
  find,
} from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { runCmd, createNSFromPaths, deleteNSFromPaths } from 'vault/tests/helpers/commands';
import { login, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from '../helpers/general-selectors';
import { NAMESPACE_PICKER_SELECTORS } from '../helpers/namespace-picker';

module('Acceptance | Enterprise | namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async () => {
    await login();
  });

  test('it focuses the search input field when user toggles namespace picker', async function (assert) {
    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the search input field is focused
    const searchInput = find(NAMESPACE_PICKER_SELECTORS.searchInput);
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

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Simulate typing into the search input
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'beep/boop');

    assert
      .dom(NAMESPACE_PICKER_SELECTORS.searchInput)
      .hasValue('beep/boop', 'The search input field has the correct value');

    // Simulate pressing Enter
    await triggerKeyEvent(NAMESPACE_PICKER_SELECTORS.searchInput, 'keydown', 'Enter');

    // Verify navigation to the matching namespace
    assert.strictEqual(
      this.owner.lookup('service:router').currentURL,
      '/vault/dashboard?namespace=beep%2Fboop',
      'Navigates to the correct namespace when Enter is pressed'
    );

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths(namespaces);
  });

  test('it filters namespaces based on search input', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespaces = ['beep/boop/bop'];
    await createNSFromPaths(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

    // Verify all namespaces are displayed initially
    assert.dom(NAMESPACE_PICKER_SELECTORS.link()).exists('Namespace link(s) exist');
    const allNamespaces = findAll(NAMESPACE_PICKER_SELECTORS.link());

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 3 namespaces are displayed after searching for "beep"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'beep');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      3,
      'Display 3 namespaces matching "beep" after searching'
    );

    // Verify 1 namespace is displayed after searching for "bop"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'bop');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      1,
      'Display 1 namespace matching "bop" after searching'
    );

    // Verify no namespaces are displayed after searching for "other"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'other');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      0,
      'No namespaces are displayed after searching for "other"'
    );

    // Clear the search input & verify all namespaces are displayed again
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, '');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      allNamespaces.length,
      'All namespaces are displayed after clearing search input'
    );

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths(namespaces);
  });

  test('it updates the namespace list after clicking "Refresh list"', async function (assert) {
    // Open the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 0 namespaces are displayed after searching for "beep"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'beep');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      0,
      'No namespaces are displayed after searching for "beep"'
    );

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Create 'beep' namespace via the CLI
    const namespaces = ['beep'];
    await createNSFromPaths(namespaces);

    // Open the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Refresh the list of namespaces
    assert.dom(NAMESPACE_PICKER_SELECTORS.refreshList).exists('Refresh list button exists');
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 1 namespace is displayed after searching for "beep"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'beep');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link('beep')).length,
      1,
      '1 namespace is displayed after searching for "beep"'
    );

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Delete the 'beep' namespace via the CLI
    await deleteNSFromPaths(namespaces);

    // Open the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Refresh the list of namespaces
    assert.dom(NAMESPACE_PICKER_SELECTORS.refreshList).exists('Refresh list button exists');
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 0 namespaces are displayed after searching for "beep"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'beep');
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      0,
      'No namespaces are displayed after searching for "beep"'
    );

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));
  });

  test('it displays the "Manage" button with the correct URL', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespaces = ['beep'];
    await createNSFromPaths(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

    // Verify the "Manage" button is rendered and has the correct URL
    assert
      .dom('[href="/ui/vault/access/namespaces"]')
      .exists('The "Manage" button is displayed with the correct URL');

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
    await click(GENERAL.toggleInput('namespace-id'));

    // Verify that the root namespace is selected by default
    assert.dom(NAMESPACE_PICKER_SELECTORS.link()).hasText('root', 'root renders as current namespace');
    assert
      .dom(`${NAMESPACE_PICKER_SELECTORS.link()} svg${GENERAL.icon('check')}`)
      .exists('The root namespace is selected');

    // Verify that the foo namespace does not exist in the namespace picker
    assert
      .dom(NAMESPACE_PICKER_SELECTORS.link(namespace))
      .exists({ count: 0 }, 'foo should not exist in the namespace picker');

    // Logout and log back into root
    await logout();
    await login();

    // Open the namespace picker & verify that the foo namespace does exist
    await click(GENERAL.toggleInput('namespace-id'));
    assert
      .dom(NAMESPACE_PICKER_SELECTORS.link(namespace))
      .exists({ count: 1 }, 'foo should exist in the namespace picker');

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths([namespace]);
  });

  // This test originated from this PR: https://github.com/hashicorp/vault/pull/7186
  test('it displays namespaces whether you log in with a namespace prefixed with / or not', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespaces = ['beep/boop/bop'];
    await createNSFromPaths(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

    // Login with a namespace prefixed with /
    await loginNs('/beep/boop');
    await settled();

    assert
      .dom(GENERAL.toggleInput('namespace-id'))
      .hasText('boop', `shows the namespace 'boop' in the toggle component`);

    // Open the namespace picker & wait for it to render
    await click(GENERAL.toggleInput('namespace-id'));
    assert.dom(`svg${GENERAL.icon('check')}`).exists('The check icon is rendered');

    // Find the selected element with the check icon & ensure it exists
    const checkIcon = find(`${NAMESPACE_PICKER_SELECTORS.link()} ${GENERAL.icon('check')}`);
    assert.dom(checkIcon).exists('A selected namespace link with the check icon exists');

    // Get the selected namespace with the data-test-namespace-link attribute & ensure it exists
    const selectedNamespace = checkIcon?.closest(NAMESPACE_PICKER_SELECTORS.link());
    assert.dom(selectedNamespace).exists('The selected namespace link exists');

    // Verify that the selected namespace has the correct data-test-namespace-link attribute and path value
    assert.strictEqual(
      selectedNamespace.getAttribute('data-test-namespace-link'),
      'beep/boop',
      'The current namespace does not begin or end with /'
    );

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths(namespaces);
  });

  test('it shows the regular namespace toolbar when not managed', async function (assert) {
    // This test is the opposite of the test in managed-namespace-test
    await logout();
    assert.strictEqual(currentURL(), '/vault/auth', 'Does not redirect');
    assert.dom(AUTH_FORM.managedNsRoot).doesNotExist('Managed namespace indicator does not exist');
    assert.dom('input[name="namespace"]').hasAttribute('placeholder', '/ (root)');
    await fillIn('input[name="namespace"]', '/foo/bar ');
    const encodedNamespace = encodeURIComponent('foo/bar');
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
    await click(GENERAL.toggleInput('namespace-id'));
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, namespace);

    assert
      .dom(NAMESPACE_PICKER_SELECTORS.link(namespace))
      .exists({ count: 1 }, 'Namespace exists in the namespace picker');

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Verify that the namespace exists in the manage namespaces page
    await fillIn(GENERAL.filterInputExplicit, namespace);
    await click(GENERAL.filterInputExplicitSearch);

    assert.dom(GENERAL.menuTrigger).exists();
    await click(GENERAL.menuTrigger);

    // Delete the namespace
    const deleteNamespaceButton = '.hds-dropdown-list-item:nth-of-type(2)';
    assert.dom(deleteNamespaceButton).hasText('Delete', 'Delete namespace button exists');
    await click(`${deleteNamespaceButton} button`);

    assert.dom(GENERAL.confirmButton).hasText('Confirm', 'Confirm namespace deletion button is shown');
    await click(GENERAL.confirmButton);

    // Verify that the namespace does not exist in the nmanage namespace page
    assert.strictEqual(
      currentURL(),
      `/vault/access/namespaces?page=1&pageFilter=${namespace}`,
      'Should remain on the manage namespaces page after deletion'
    );

    assert
      .dom('.list-item-row')
      .exists({ count: 0 }, 'Namespace should be deleted and not displayed in the list');

    // Verify that the namespace does not exist in the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, namespace);
    assert
      .dom(NAMESPACE_PICKER_SELECTORS.link())
      .exists({ count: 0 }, 'Deleted namespace does not exist in the namespace picker');
  });

  test('it should show root in namespace picker when the user explicitly logs into root namespace', async function (assert) {
    // Explicitly set root as the namespace to login to
    await loginNs('root');

    assert
      .dom(NAMESPACE_PICKER_SELECTORS.toggle)
      .hasText('root', `shows the namespace 'root' in the toggle component`);

    // Verify user is in root namespace
    assert.true(
      this.owner.lookup('service:namespace').inRootNamespace,
      'Verifies that the user is in the root namespace'
    );
  });
});
