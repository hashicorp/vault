/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { createNSFromPaths, deleteNSFromPaths } from 'vault/tests/helpers/commands';

module('Acceptance | Enterprise | /access/namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async () => {
    await login();
  });

  test('the route url navigates to namespace index page', async function (assert) {
    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.namespaces.index',
      'navigates to the correct route'
    );
    assert.dom(GENERAL.pageTitle).hasText('Namespaces', 'Page title is displayed correctly');
  });

  test('the route displays the breadcrumb trail', async function (assert) {
    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    assert.dom(GENERAL.breadcrumb).exists({ count: 1 }, 'Only one breadcrumb is displayed');
    assert.dom(GENERAL.breadcrumb).hasText('Namespaces', 'Breadcrumb trail is displayed correctly');
  });

  test('the route should update namespace list after create/delete WITH manual refresh in the CLI', async function (assert) {
    const testNS = 'test-refresh-ns-cli';

    // Setup: Create namespace via the CLI
    const namespaces = [testNS];
    await createNSFromPaths(namespaces);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText(
        'No namespaces yet',
        'Empty state is displayed when searching for the namespace we have created in the CLI but have not refreshed the list yet'
      );
    // Click the refresh list button on the namespace page
    await click(GENERAL.button('refresh-namespace-list'));
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
    assert.dom('[data-test-list-item]').hasText(testNS, 'Namespace is displayed after refreshing the list');

    // Delete the created namespace via the CLI
    await deleteNSFromPaths(namespaces);
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

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    // Verify test-create-ns does not exist in the Manage Namespace page
    await fillIn(GENERAL.filterInputExplicit, testNS);
    await click(GENERAL.button('Search'));
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

    // Verify test-create-ns does not exist in the Manage Namespace page
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No namespaces yet', 'Empty state is displayed indicating the namespace was deleted');
  });

  test('the route should show "delete" and "switch to namespace" option menus for each namespace', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const namespace = 'asdf';
    await createNSFromPaths([namespace]);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    // Hack: Trigger refresh sys/internal/namespace namespaces endpoint that is only hit on the namespace picker (not the namespace route)
    await click(GENERAL.toggleInput('namespace-picker'));
    await click(GENERAL.button('Refresh list'));

    // Enter search text
    await fillIn(GENERAL.filterInputExplicit, namespace);
    await click(GENERAL.button('Search'));

    // Verify the menu options
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('switch')).exists('Switch to namespace option is displayed');
    assert.dom(GENERAL.menuItem('delete')).exists('Delete namespace option is displayed');

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths([namespace]);
  });

  test('the route should switch to the selected namespace on click "Switch to namespace"', async function (assert) {
    // Setup: Create namespace(s) via the CLI
    const testNS = 'test-create-ns-switch';
    await createNSFromPaths([testNS]);

    // Go to the manage namespaces page
    await visit('/vault/access/namespaces');

    // Hack: Trigger refresh internal namespaces endpoint
    await click(GENERAL.toggleInput('namespace-picker'));
    await click(GENERAL.button('Refresh list'));

    // Switch namespace
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('switch'));

    // Verify that we switched namespaces
    await click(GENERAL.toggleInput('namespace-picker'));
    assert.dom('[data-test-badge-namespace]').hasText(testNS);
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'navigates to the correct route');

    // Cleanup: Delete namespace(s) via the CLI
    await deleteNSFromPaths([testNS]);
  });
});
