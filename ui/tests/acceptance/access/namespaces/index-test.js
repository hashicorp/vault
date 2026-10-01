/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, currentURL, fillIn, visit, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { createNS, deleteNS, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';

module('Acceptance | Enterprise | /access/namespaces', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    // dismiss the wizard
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.namespace);
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

  test('the namespace path column sorts namespaces', async function (assert) {
    const namespaces = ['test-sort-namespace-z', 'test-sort-namespace-a'];
    for (const namespace of namespaces) {
      await runCmd(createNS(namespace), false);
    }

    await click(GENERAL.button('refresh-namespace-list'));
    await fillIn(GENERAL.filterInput, 'test-sort-namespace-');
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));

    assert.dom(GENERAL.tableData(0, 'id')).hasText(namespaces[1], 'the first path sorts alphabetically');
    assert.dom(GENERAL.tableData(1, 'id')).hasText(namespaces[0], 'the second path sorts alphabetically');

    for (const namespace of namespaces) {
      await runCmd(deleteNS(namespace), false);
    }
  });

  test('the route should update namespace list after create/delete WITH manual refresh in the CLI', async function (assert) {
    const testNS = 'test-refresh-ns-cli';

    // Setup: Create namespace via the CLI
    await runCmd(createNS(testNS), false);

    // Click the refresh list button on the namespace page
    await click(GENERAL.button('refresh-namespace-list'));
    await fillIn(GENERAL.filterInput, testNS);

    assert.dom('[data-test-list-item]').hasText(testNS, 'Namespace is displayed after refreshing the list');

    // Delete the created namespace via the CLI
    await runCmd(deleteNS(testNS), false);
    await visit('/vault/access/namespaces');
    await fillIn(GENERAL.filterInput, '');

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

    // Create a new namespace in the UI
    await click(GENERAL.button('Create namespace'));
    await fillIn(GENERAL.inputByAttr('path'), testNS);
    await click(GENERAL.submitButton);

    // Verify test-create-ns-ui exists in the Manage Namespace page
    await fillIn(GENERAL.filterInput, testNS);

    assert.dom('[data-test-list-item]').hasText(testNS, 'Namespace is displayed after refreshing the list');

    // Delete the created namespace
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    // Clear filter and refresh
    await fillIn(GENERAL.filterInput, '');
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
    await click(GENERAL.button('refresh-namespace-list'));

    // Search for created namespace// Enter search text
    await fillIn(GENERAL.filterInput, testNS);

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
    await click(GENERAL.button('refresh-namespace-list'));

    // Search for created namespace
    await fillIn(GENERAL.filterInput, testNS);

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

// A read-only token must not be offered create/delete actions it cannot perform,
// while still being able to switch into a namespace it can read.
module('Acceptance | Enterprise | /access/namespaces | read-only token', function (hooks) {
  setupApplicationTest(hooks);

  const READ_ONLY_POLICY = `
    path "sys/namespaces" { capabilities = ["read", "list"] }
    path "sys/namespaces/*" { capabilities = ["read", "list"] }
  `;

  hooks.beforeEach(async function () {
    this.namespace = 'read-only-ns';
    await login();
    await runCmd(createNS(this.namespace), false);
    const token = await runCmd(tokenWithPolicyCmd('namespaces-read-only', READ_ONLY_POLICY));
    await login(token);
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.namespace);
    await visit('/vault/access/namespaces');
  });

  hooks.afterEach(async function () {
    await login();
    await runCmd(deleteNS(this.namespace), false);
  });

  test('it hides create and delete actions but offers switch', async function (assert) {
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Namespaces', 'list view renders');
    assert.dom(GENERAL.button('Create namespace')).doesNotExist('create action is hidden');

    await fillIn(GENERAL.filterInput, this.namespace);
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('switch')).exists('switch action is offered for a readable namespace');
    assert.dom(GENERAL.menuItem('delete')).doesNotExist('delete action is hidden');

    await click(GENERAL.menuItem('switch'));
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'switch lands on the dashboard');
    assert.true(currentURL().includes(`namespace=${this.namespace}`), 'switch targets the namespace');
  });
});

// Create permission is often scoped to a namespace prefix, which capabilities-self cannot evaluate for
// an unnamed namespace, so the header action must still show for these tokens.
module('Acceptance | Enterprise | /access/namespaces | path-scoped token', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const scopedPolicy = `
      path "sys/namespaces" { capabilities = ["read", "list"] }
      path "sys/namespaces/team-*" { capabilities = ["update"] }
    `;
    await login();
    const token = await runCmd(tokenWithPolicyCmd('namespaces-team-scoped', scopedPolicy));
    await login(token);
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.namespace);
  });

  test('it shows the create action when update is scoped to a namespace prefix', async function (assert) {
    await visit('/vault/access/namespaces');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Namespaces', 'list view renders');
    assert.dom(GENERAL.button('Create namespace')).exists('create action is offered');
  });
});

// The api-lock routes share the sys/namespaces prefix but cannot create namespaces, and a token that
// can only list namespaces has no row actions, so neither should leave behind a dead control.
module('Acceptance | Enterprise | /access/namespaces | limited tokens', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.namespace = 'limited-token-ns';
    await login();
    await runCmd(createNS(this.namespace), false);
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.namespace);
  });

  hooks.afterEach(async function () {
    await login();
    await runCmd(deleteNS(this.namespace), false);
  });

  test('it hides the create action for a token that can only lock namespaces', async function (assert) {
    const lockOnlyPolicy = `
      path "sys/namespaces" { capabilities = ["read", "list"] }
      path "sys/namespaces/api-lock/lock/*" { capabilities = ["update"] }
    `;
    await login(await runCmd(tokenWithPolicyCmd('namespaces-lock-only', lockOnlyPolicy)));
    await visit('/vault/access/namespaces');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Namespaces', 'list view renders');
    assert.dom(GENERAL.button('Create namespace')).doesNotExist('create action is hidden');
  });

  test('it omits the row menu when the token has no actions on a namespace', async function (assert) {
    const listOnlyPolicy = 'path "sys/namespaces" { capabilities = ["list"] }';
    await login(await runCmd(tokenWithPolicyCmd('namespaces-list-only', listOnlyPolicy)));
    await visit('/vault/access/namespaces');
    await fillIn(GENERAL.filterInput, this.namespace);
    assert.dom(GENERAL.listItem(this.namespace)).exists('namespace is listed');
    assert.dom(GENERAL.menuTrigger).doesNotExist('no empty row menu is rendered');
  });
});
