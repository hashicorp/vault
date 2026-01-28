/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, visit, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { createNS, deleteNS, runCmd } from 'vault/tests/helpers/commands';
import localStorage from 'vault/lib/local-storage';

module('Acceptance | Enterprise | /access/namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async () => {
    await login();
    // dismiss the wizard
    localStorage.setItem('dismissed-wizards', ['namespace']);
    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');
  });

  test('the route url navigates to namespace index page', async function (assert) {
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.namespaces.index',
      'navigates to the correct route'
    );

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Namespaces', 'Page title is displayed correctly');
  });

  test('the route displays the breadcrumb trail', async function (assert) {
    assert.dom(GENERAL.breadcrumb).exists({ count: 2 }, 'Only two breadcrumb is displayed');
    assert.dom(GENERAL.breadcrumbAtIdx(0)).hasText('Vault', 'Breadcrumb trail is displayed correctly');
    assert
      .dom(GENERAL.currentBreadcrumb('Namespaces'))
      .hasText('Namespaces', 'Namespace breadcrumb trail is displayed correctly');
  });

  test('the route should update namespace list after create/delete WITH manual refresh in the CLI', async function (assert) {
    const testNS = 'test-refresh-ns-cli';

    // Setup: Create namespace via the CLI
    await runCmd(createNS(testNS), false);

    // Click the refresh list button on the namespace page
    await click(GENERAL.button('refresh-namespace-list'));
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
    assert.dom('[data-test-list-item]').hasText(testNS, 'Namespace is displayed after refreshing the list');

    // Delete the created namespace via the CLI
    await runCmd(deleteNS(testNS), false);
    await visit('/vault/access/namespaces');

    // Search for the deleted namespace
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));

    // Click the refresh list button from the namespace page
    await click(GENERAL.button('refresh-namespace-list'));
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText(
        'No namespaces yet',
        'Empty state is displayed when searching for the namespace we have created in the CLI but have not refreshed the list yet'
      );
  });

  test('the route should update namespace list after create/delete WITHOUT manual refresh in the UI', async function (assert) {
    const testNS = 'test-create-ns-ui';

    // Verify test-create-ns does not exist in the Manage Namespace page
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
    await waitFor(GENERAL.emptyStateTitle, {
      timeout: 2000,
      timeoutMessage: 'timed out waiting for empty state title to render',
    });
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText(
        'No namespaces yet',
        'Empty state is displayed when searching for the namespace we have created in the UI but have not refreshed the list yet'
      );

    // Create a new namespace in the UI
    await click(GENERAL.linkTo('create-namespace'));
    await fillIn(GENERAL.inputByAttr('path'), testNS);
    await click(GENERAL.submitButton);

    // Verify test-create-ns-ui exists in the Manage Namespace page
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
    assert.dom('[data-test-list-item]').hasText(testNS, 'Namespace is displayed after refreshing the list');

    // Delete the created namespace
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);
    await click(GENERAL.button('refresh-namespace-list'));

    // Verify test-create-ns does not exist in the Manage Namespace page
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No namespaces yet', 'Empty state is displayed indicating the namespace was deleted');
  });

  test('the route should show "delete" option menu for each namespace', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const testNS = 'asdf';
    await runCmd(createNS(testNS), false);

    // Search for created namespace// Enter search text
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
    await click(GENERAL.button('refresh-namespace-list'));

    // Verify the menu options
    await waitFor(GENERAL.menuTrigger, {
      timeout: 2000,
      timeoutMessage: 'timed out waiting for menu trigger to render',
    });
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('delete')).exists('Delete namespace option is displayed');

    // Cleanup: Delete namespace(s) via the CLI
    await runCmd(deleteNS(testNS), false);
  });

  test('the route should switch to the selected namespace on click "Switch to namespace"', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const testNS = 'test-create-ns-switch';
    await runCmd(createNS(testNS), false);

    // Search for created namespace
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
    await click(GENERAL.button('refresh-namespace-list'));

    // Switch namespace
    await waitFor(GENERAL.menuTrigger);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('switch'));

    // Verify that we switched namespaces
    await click(GENERAL.button('namespace-picker'));
    assert.dom('[data-test-badge-namespace]').hasText(testNS, 'Namespace badge shows the correct namespace');
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'navigates to the correct route');

    // Cleanup: Delete namespace(s) via the CLI
    await visit('vault/dashboard'); // navigate to "root" before deleting
    await runCmd(deleteNS(testNS), false);
  });
});
