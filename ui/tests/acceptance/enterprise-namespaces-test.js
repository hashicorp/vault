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
import { runCmd, createNS, deleteNS } from 'vault/tests/helpers/commands';
import { login, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { GENERAL } from '../helpers/general-selectors';
import { NAMESPACE_PICKER_SELECTORS } from '../helpers/namespace-picker';

import sinon from 'sinon';

async function createNamespaces(namespaces) {
  for (const ns of namespaces) {
    // Note: iterate through the namespace parts to create the full namespace path
    const parts = ns.split('/');
    let currentPath = '';

    for (const part of parts) {
      // Visit the parent namespace
      const url = `/vault/dashboard${currentPath && `?namespace=${currentPath.replaceAll('/', '%2F')}`}`;
      await visit(url);

      currentPath = currentPath ? `${currentPath}/${part}` : part;

      // Create the current namespace
      await runCmd(createNS(part), false);
      await settled();
    }

    // Reset to the root namespace
    const url = '/vault/dashboard';
    await visit(url);
  }
}

async function deleteNamespaces(namespaces) {
  // Reset to the root namespace
  const url = '/vault/dashboard';
  await visit(url);

  for (const ns of namespaces) {
    // Note: delete the parent namespace to delete all child namespaces
    const part = ns.split('/')[0];
    await runCmd(deleteNS(part), false);
    await settled();
  }
}

module('Acceptance | Enterprise | namespaces', function (hooks) {
  setupApplicationTest(hooks);

  let fetchSpy;

  hooks.beforeEach(() => {
    fetchSpy = sinon.spy(window, 'fetch');
    return login();
  });

  hooks.afterEach(() => {
    fetchSpy.restore();
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
    // Test Setup
    const namespaces = ['beep/boop'];
    await createNamespaces(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(GENERAL.buttonByAttr('refresh-namespaces'));
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

    // Test Cleanup
    await deleteNamespaces(namespaces);
  });

  test('it filters namespaces based on search input', async function (assert) {
    // Test Setup
    const namespaces = ['beep/boop/bop'];
    await createNamespaces(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(GENERAL.buttonByAttr('refresh-namespaces'));

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

    // Test Cleanup
    await deleteNamespaces(namespaces);
  });

  test('it updates the namespace list after clicking "Refresh list"', async function (assert) {
    // Test Setup
    const namespaces = ['beep'];
    await createNamespaces(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the namespace list was fetched on load
    let listNamespaceRequests = fetchSpy
      .getCalls()
      .filter((call) => call.args[0].includes('/v1/sys/internal/ui/namespaces'));
    assert.strictEqual(
      listNamespaceRequests.length,
      1,
      'The network call to the specific endpoint was made twice (once on load, once on refresh)'
    );

    // Refresh the list of namespaces
    assert.dom(GENERAL.buttonByAttr('refresh-namespaces')).exists('Refresh list button exists');
    await click(GENERAL.buttonByAttr('refresh-namespaces'));

    // Verify that the namespace list was fetched on refresh
    listNamespaceRequests = fetchSpy
      .getCalls()
      .filter((call) => call.args[0].includes('/v1/sys/internal/ui/namespaces'));
    assert.strictEqual(
      listNamespaceRequests.length,
      2,
      'The network call to the specific endpoint was made twice (once on load, once on refresh)'
    );

    // Test Cleanup
    await deleteNamespaces(namespaces);
  });

  test('it displays the "Manage" button with the correct URL', async function (assert) {
    // Test Setup
    const namespaces = ['beep'];
    await createNamespaces(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(GENERAL.buttonByAttr('refresh-namespaces'));

    // Verify the "Manage" button is rendered and has the correct URL
    assert
      .dom('[href="/ui/vault/access/namespaces"]')
      .exists('The "Manage" button is displayed with the correct URL');

    // Test Cleanup
    await deleteNamespaces(namespaces);
  });

  // This test originated from this PR: https://github.com/hashicorp/vault/pull/7186
  test('it clears namespaces when you log out', async function (assert) {
    // Test Setup
    const namespaces = ['foo'];
    await createNamespaces(namespaces);

    const ns = 'foo';
    await runCmd(createNS(ns), false);
    const token = await runCmd(`write -field=client_token auth/token/create policies=default`);
    await login(token);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    assert.dom(NAMESPACE_PICKER_SELECTORS.link()).hasText('root', 'root renders as current namespace');
    assert
      .dom(`${NAMESPACE_PICKER_SELECTORS.link()} svg${GENERAL.icon('check')}`)
      .exists('The root namespace is selected');

    // Test Cleanup
    await deleteNamespaces(namespaces);
  });

  // This test originated from this PR: https://github.com/hashicorp/vault/pull/7186
  test('it displays namespaces whether you log in with a namespace prefixed with / or not', async function (assert) {
    // Test Setup
    const namespaces = ['beep/boop/bop'];
    await createNamespaces(namespaces);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);
    await click(GENERAL.buttonByAttr('refresh-namespaces'));

    // Login with a namespace prefixed with /
    await loginNs('/beep/boop');
    await settled();

    assert
      .dom(NAMESPACE_PICKER_SELECTORS.toggle)
      .hasText('boop', `shows the namespace 'boop' in the toggle component`);

    // Open the namespace picker & wait for it to render
    await click(NAMESPACE_PICKER_SELECTORS.toggle);
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

    // Test Cleanup
    await deleteNamespaces(namespaces);
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
    // Test Setup
    const namespaces = ['test-delete-me'];
    await createNamespaces(namespaces);

    await visit('/vault/access/namespaces');

    const searchInput = GENERAL.filterInputExplicit;
    const searchButton = GENERAL.filterInputExplicitSearch;

    await fillIn(searchInput, 'test-delete-me');
    await click(searchButton);

    assert.dom(GENERAL.menuTrigger).exists();
    await click(GENERAL.menuTrigger);

    // Verify that the user can delete the namespace
    const deleteNamespaceButton = '.hds-dropdown-list-item:nth-of-type(1)';
    assert.dom(deleteNamespaceButton).hasText('Delete', 'Allow users to delete the namespace');
    await click(`${deleteNamespaceButton} button`);

    assert.dom(GENERAL.confirmButton).hasText('Confirm', 'Allow users to delete the namespace');
    await click(GENERAL.confirmButton);

    assert.strictEqual(
      currentURL(),
      '/vault/access/namespaces?page=1&pageFilter=test-delete-me',
      'Should remain on the manage namespaces page after deletion'
    );

    assert
      .dom('.list-item-row')
      .exists({ count: 0 }, 'Namespace should be deleted and not displayed in the list');
  });
});
