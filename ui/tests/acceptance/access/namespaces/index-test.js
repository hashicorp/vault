/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, visit, click, fillIn, currentURL, findAll } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { createNSFromPaths, deleteNSFromPaths } from 'vault/tests/helpers/commands';
import { NAMESPACE_PICKER_SELECTORS } from 'vault/tests/helpers/namespace-picker';

module('Acceptance | Enterprise | /access/namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async () => {
    await login();
  });

  test('it navigates to namespaces page', async function (assert) {
    assert.expect(1);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.namespaces.index',
      'navigates to the correct route'
    );
  });

  test('it displays the breadcrumb trail', async function (assert) {
    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    assert.dom(GENERAL.breadcrumb).exists({ count: 1 }, 'Only one breadcrumb is displayed');
    assert.dom(GENERAL.breadcrumb).hasText('Namespaces', 'Breadcrumb trail is displayed correctly');
  });

  test('it should render correct number of namespaces', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespaces = [
      'ns1',
      'ns2',
      'ns3',
      'ns4',
      'ns5',
      'ns6',
      'ns7',
      'ns8',
      'ns9',
      'ns10',
      'ns11',
      'ns12',
      'ns13',
      'ns14',
      'ns15',
      'ns16',
      'ns17',
      'ns18',
    ];
    await createNSFromPaths(namespaces);

    assert.expect(3);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    const store = this.owner.lookup('service:store');

    // Default page size is 15
    assert.strictEqual(store.peekAll('namespace').length, 15, 'Store has 15 namespaces records');
    assert.dom('.list-item-row').exists({ count: 15 }, 'Should display 15 namespaces');
    assert.dom('.hds-pagination').exists();

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths(namespaces);
  });

  test('it should show button to refresh namespace list', async function (assert) {
    const testNS = 'test-refresh-ns';

    // Setup: Create namespace via the CLI
    const namespaces = [testNS];
    await createNSFromPaths(namespaces);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    // Open the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 0 namespaces are displayed after searching for "test-refresh-ns"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, testNS);
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      0,
      `No namespaces are displayed after searching for "${testNS}"`
    );

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Click the refresh list button
    assert
      .dom(GENERAL.button('refresh-namespace-list'))
      .hasText('Refresh list', 'Refresh button is rendered correctly');
    await click(GENERAL.button('refresh-namespace-list'));

    // Open the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 1 namespace is displayed after searching for "test-refresh-ns"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, testNS);
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      1,
      `1 namespace is displayed after searching for "${testNS}"`
    );

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Cleanup: Delete namespace via the CLI
    await deleteNSFromPaths(namespaces);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    // Open the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 1 namespace is displayed after searching for "test-refresh-ns"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, testNS);
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      1,
      `1 namespace is displayed after searching for "${testNS}"`
    );

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Click the refresh list button
    assert
      .dom(GENERAL.button('refresh-namespace-list'))
      .hasText('Refresh list', 'Refresh button is rendered correctly');
    await click(GENERAL.button('refresh-namespace-list'));

    // Open the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));

    // Verify the search input field exists
    assert.dom(NAMESPACE_PICKER_SELECTORS.searchInput).exists('The namespace search field exists');

    // Verify 0 namespaces are displayed after searching for "test-refresh-ns"
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, testNS);
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      0,
      `No namespaces are displayed after searching for "${testNS}"`
    );

    // Close the namespace picker
    await click(GENERAL.toggleInput('namespace-id'));
  });

  test('it should show button to create new namespace', async function (assert) {
    const createNamespaceLink = GENERAL.linkTo('create-namespace');

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    assert
      .dom(createNamespaceLink)
      .hasText('Create namespace', 'Create namespace button is rendered correctly');
    assert
      .dom(createNamespaceLink)
      .hasAttribute(
        'href',
        '/ui/vault/access/namespaces/create',
        'Create namespace button has the correct href attribute'
      );
  });

  test('it should filter namespaces based on search input', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespaces = ['parent', 'other-parent'];
    await createNSFromPaths(namespaces);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    // Enter search text
    await fillIn(GENERAL.filterInputExplicit, 'other');
    assert.dom(GENERAL.filterInputExplicit).hasValue('other', 'Search input contains the entered text');

    // Click the search button
    await click(GENERAL.filterInputExplicitSearch);

    // Verify the filtered results
    assert.dom('.list-item-row').exists({ count: 1 }, 'Filtered results are displayed correctly');
    assert
      .dom('.list-item-row')
      .hasText('other-parent', 'Correct namespace is displayed in the filtered results');

    // Verify the URL query param is updated
    assert.strictEqual(
      currentURL(),
      '/vault/access/namespaces?page=1&pageFilter=other',
      'URL query param is updated to reflect the search field as pageFilter'
    );

    // Clear the search input
    await fillIn(GENERAL.filterInputExplicit, '');
    await click(GENERAL.filterInputExplicitSearch);

    assert.dom(GENERAL.filterInputExplicit).hasValue('', 'Search input is cleared');
    assert
      .dom('.list-item-row')
      .exists({ count: 2 }, 'All namespaces are displayed after clearing the search input');
    assert.strictEqual(
      currentURL(),
      '/vault/access/namespaces?page=1',
      'URL query param is updated to remove pageFilter'
    );

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths(namespaces);
  });

  test('it should show options menu for each namespace', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespace = 'meep';
    await createNSFromPaths([namespace]);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    // Enter search text
    await fillIn(GENERAL.filterInputExplicit, namespace);
    await click(GENERAL.filterInputExplicitSearch);

    await click(GENERAL.button('refresh-namespace-list'));

    assert.dom(GENERAL.menuTrigger).exists('Namespace options menu is displayed');
    await click(GENERAL.menuTrigger);
    assert.dom('.hds-dropdown-list-item').exists({ count: 2 }, 'Should display 2 options in the menu.');

    // Verify that the user can switch to the namespace
    const switchNamespaceButton = '.hds-dropdown-list-item:nth-of-type(1)';
    assert
      .dom(switchNamespaceButton)
      .hasText('Switch to namespace', 'Allow users to switch to different namespace');
    assert
      .dom(`${switchNamespaceButton} a`)
      .hasAttribute(
        'href',
        `http://localhost:7357/ui/vault/dashboard?namespace=${namespace}`,
        'Switch namespace button has the correct href attribute'
      );

    // Verify that the user can delete the namespace
    const deleteNamespaceButton = '.hds-dropdown-list-item:nth-of-type(2)';
    assert.dom(deleteNamespaceButton).hasText('Delete', 'Allow users to delete the namespace');

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths([namespace]);
  });
});
