/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, visit, click, fillIn, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | Enterprise | /access/namespaces', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  const searchInput = GENERAL.filterInputExplicit;
  const searchButton = GENERAL.filterInputExplicitSearch;

  hooks.beforeEach(function () {
    return login();
  });

  test('it navigates to namespaces page', async function (assert) {
    assert.expect(1);
    await visit('/vault/access/namespaces');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.namespaces.index',
      'navigates to the correct route'
    );
  });

  test('it displays the breadcrumb trail', async function (assert) {
    await visit('/vault/access/namespaces');
    assert.dom(GENERAL.breadcrumb).exists({ count: 1 }, 'Only one breadcrumb is displayed');
    assert.dom(GENERAL.breadcrumb).hasText('Namespaces', 'Breadcrumb trail is displayed correctly');
  });

  test('it should render correct number of namespaces', async function (assert) {
    assert.expect(3);
    await visit('/vault/access/namespaces');
    const store = this.owner.lookup('service:store');
    // Default page size is 15
    assert.strictEqual(store.peekAll('namespace').length, 15, 'Store has 15 namespaces records');
    assert.dom('.list-item-row').exists({ count: 15 }, 'Should display 15 namespaces');
    assert.dom('.hds-pagination').exists();
  });

  test('it should show button to refresh namespace list', async function (assert) {
    let refreshNetworkRequestTriggered;
    const refreshNamespaceButton = GENERAL.buttonByAttr('refresh-namespace-list');

    this.server.get('/sys/internal/ui/namespaces', () => {
      refreshNetworkRequestTriggered = true;
      return;
    });

    await visit('/vault/access/namespaces');

    assert.dom(refreshNamespaceButton).hasText('Refresh list', 'Refresh button is rendered correctly');

    refreshNetworkRequestTriggered = false;
    await click(refreshNamespaceButton);
    assert.true(
      refreshNetworkRequestTriggered,
      'Get namespaces network request was made when refresh button was clicked'
    );
  });

  test('it should show button to create new namespace', async function (assert) {
    const createNamespaceLink = GENERAL.linkTo('create-namespace');

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
    await visit('/vault/access/namespaces');

    // Enter search text
    await fillIn(searchInput, 'ns4');
    assert.dom(searchInput).hasValue('ns4', 'Search input contains the entered text');

    // Click the search button
    await click(searchButton);

    // Verify the filtered results
    assert.dom('.list-item-row').exists({ count: 1 }, 'Filtered results are displayed correctly');
    assert.dom('.list-item-row').hasText('ns4', 'Correct namespace is displayed in the filtered results');

    // Verify the URL query param is updated
    assert.strictEqual(
      currentURL(),
      '/vault/access/namespaces?page=1&pageFilter=ns4',
      'URL query param is updated to reflect the search field as pageFilter'
    );

    // Clear the search input
    await fillIn(searchInput, '');
    await click(searchButton);
    assert.dom(searchInput).hasValue('', 'Search input is cleared');
    assert
      .dom('.list-item-row')
      .exists({ count: 15 }, 'All namespaces are displayed after clearing the search input');
    assert.strictEqual(
      currentURL(),
      '/vault/access/namespaces?page=1',
      'URL query param is updated to remove pageFilter'
    );
  });

  test('it should show options menu for each namespace', async function (assert) {
    await visit('/vault/access/namespaces');
    assert.dom(GENERAL.menuTrigger).exists();
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
        'http://localhost:7357/ui/vault/dashboard?namespace=ns1',
        'Switch namespace button has the correct href attribute'
      );

    // Verify that the user can delete the namespace
    const deleteNamespaceButton = '.hds-dropdown-list-item:nth-of-type(2)';
    assert.dom(deleteNamespaceButton).hasText('Delete', 'Allow users to delete the namespace');
  });

  test('it should hide the switch to namespace option for unaccessible namespaces', async function (assert) {
    await visit('/vault/access/namespaces');

    // Search for a namespace that is not accessible
    await fillIn(searchInput, 'ns12');
    await click(searchButton);

    assert.dom(GENERAL.menuTrigger).exists();
    await click(GENERAL.menuTrigger);

    // Verify that only the delete option is available for the unaccessible namespace
    assert.dom('.hds-dropdown-list-item').exists({ count: 1 }, 'Should display 1 option in the menu.');

    // Verify that the user can delete the namespace
    const deleteNamespaceButton = '.hds-dropdown-list-item:nth-of-type(1)';
    assert.dom(deleteNamespaceButton).hasText('Delete', 'Allow users to delete the namespace');
  });
});
