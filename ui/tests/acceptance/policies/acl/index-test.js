/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// Two non-root, non-default policies that will appear in the list.
const STUB_POLICIES = ['my-policy', 'another-policy'];

module('Acceptance | ACL policies list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    // Dismiss the wizard so it does not obscure the list view.
    this.owner.lookup('service:wizard').dismiss('acl-policy');
    // Stub the list endpoint used by the route model.
    // Mirage namespace is 'v1' so paths without a leading slash resolve to /v1/<path>.
    // The API client calls GET /v1/sys/policies/acl/ (trailing slash) + ?list=true.
    // ACL policies are never empty (root and default always exist), so the
    // stub intentionally returns a non-empty set of keys.
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: ['default', 'root', ...STUB_POLICIES] },
      request_id: 'test',
    }));
    // Stub the capabilities endpoint so capability checks don't fail.
    // VoidApiResponse.value() returns the raw JSON body; the capabilities
    // service destructures { data } from it.  An empty data object causes
    // all capability checks to default to true.
    this.server.post('sys/capabilities-self', () => ({
      data: {},
      request_id: 'test',
    }));
  });

  // ── Page structure ───────────────────────────────────────────────────────

  test('it renders the page title', async function (assert) {
    // Verifies Page::ListView renders the title from listViewConfig.
    await visit('/vault/policies/acl');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('ACL policies');
  });

  test('it renders breadcrumbs', async function (assert) {
    // Verifies both breadcrumb entries are present on the policies index page.
    await visit('/vault/policies/acl');
    assert.dom(GENERAL.breadcrumbLink('Vault')).exists('first breadcrumb links to Vault dashboard');
    assert.dom(GENERAL.currentBreadcrumb('ACL policies')).exists('last breadcrumb is current page');
  });

  test('it renders the Create ACL policy primary action button', async function (assert) {
    // Verifies the primary action button is rendered by Page::ListView.
    await visit('/vault/policies/acl');
    assert.dom(GENERAL.button('Create ACL policy')).exists();
  });

  // ── Table columns ────────────────────────────────────────────────────────

  test('it renders the Policy name column header', async function (assert) {
    // Verifies the Policy name column is rendered in the table header.
    // ListTable uses Hds::AdvancedTable so the isAdvanced flag is required.
    await visit('/vault/policies/acl');
    assert.dom(GENERAL.tableColumnHeader(1, { isAdvanced: true })).includesText('Policy name');
  });

  // ── Name column rendering ────────────────────────────────────────────────

  test('it renders policy names as links for non-root policies', async function (assert) {
    // Verifies that non-root policies show a navigable link in the Name column.
    await visit('/vault/policies/acl');
    assert.dom('[data-test-policy-link="my-policy"]').exists('non-root policy renders as a link');
  });

  test('it renders the root policy name as plain text without a link', async function (assert) {
    // The root policy is a special built-in and must not be navigable or deletable.
    await visit('/vault/policies/acl');
    assert.dom('[data-test-policy-name]').exists('root policy name element exists');
    assert.dom('[data-test-policy-link="root"]').doesNotExist('root policy has no link');
  });

  // ── Filter ───────────────────────────────────────────────────────────────

  test('it shows the filtered empty state when no policies match the filter', async function (assert) {
    // Verifies the filteredEmptyTitle is displayed when filter text matches nothing.
    await visit('/vault/policies/acl');
    await fillIn(GENERAL.filterInput, 'nonexistent-policy-xyz');
    assert.dom(GENERAL.emptyStateTitle).includesText('No results for');
  });

  // ── Row actions — root policy ─────────────────────────────────────────────

  test('it does not render a popup menu for the root policy', async function (assert) {
    // The root policy row must suppress all dropdown actions.
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: ['root'] },
      request_id: 'test',
    }));
    await visit('/vault/policies/acl');
    assert.dom(GENERAL.menuTrigger).doesNotExist('root policy has no popup menu');
  });

  // ── Row actions — regular policy ─────────────────────────────────────────

  test('it renders popup menu actions for a regular policy', async function (assert) {
    // Verifies View, Edit, Download, and Delete actions are present for a regular policy.
    // Stub a single non-special policy so the first (and only) menu trigger belongs to it.
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: ['my-policy'] },
      request_id: 'test',
    }));
    this.server.post('sys/capabilities-self', () => ({
      data: { 'sys/policies/acl/my-policy': ['read', 'update', 'delete'] },
      request_id: 'test',
    }));
    await visit('/vault/policies/acl');
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('view-policy')).exists('View policy action exists');
    assert.dom(GENERAL.menuItem('edit-policy')).exists('Edit policy action exists');
    assert.dom(GENERAL.menuItem('download-policy')).exists('Download policy action exists');
    assert.dom(GENERAL.menuItem('delete-policy')).exists('Delete policy action exists');
  });

  test('it suppresses Delete for the default policy', async function (assert) {
    // default and default-ceiling policies must not show the Delete action.
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: ['default'] },
      request_id: 'test',
    }));
    this.server.post('sys/capabilities-self', () => ({
      data: { 'sys/policies/acl/default': ['read', 'update', 'delete'] },
      request_id: 'test',
    }));
    await visit('/vault/policies/acl');
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('delete-policy')).doesNotExist('default policy has no Delete action');
  });

  // ── Delete modal ─────────────────────────────────────────────────────────

  test('it shows a confirmation modal when Delete is clicked', async function (assert) {
    // Verifies the ConfirmModal appears when the critical Delete action is triggered.
    // Stub a single non-special policy so the first trigger belongs to it.
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: ['my-policy'] },
      request_id: 'test',
    }));
    this.server.post('sys/capabilities-self', () => ({
      data: { 'sys/policies/acl/my-policy': ['read', 'update', 'delete'] },
      request_id: 'test',
    }));
    await visit('/vault/policies/acl');
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete-policy'));
    assert.dom(GENERAL.confirmTitle).hasText('Delete policy?');
  });

  // ── Pagination ───────────────────────────────────────────────────────────

  test('it advances the page query param when the next-page button is clicked', async function (assert) {
    // Verifies that clicking the pagination next-page control updates the URL
    // query param. Page::ListView wires onPageChange → router.transitionTo so
    // the page QP stays in sync with what the user sees.
    // Stub > DEFAULT_PAGE_SIZE (15) policies so the pagination control renders.
    const manyPolicies = Array.from({ length: 20 }, (_, i) => `policy-${i}`);
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: manyPolicies },
      request_id: 'test',
    }));
    await visit('/vault/policies/acl');
    await click(GENERAL.nextPage);
    assert.true(currentURL().includes('page=2'), 'URL contains page=2 after clicking next');
  });

  test('it resets the page query param to 1 on navigation away', async function (assert) {
    // Verifies resetController resets the page QP when navigating away.
    await visit('/vault/policies/acl?page=3');
    await visit('/vault/dashboard');
    await visit('/vault/policies/acl');
    assert.false(currentURL().includes('page=3'), 'page param was reset after navigation');
  });

  test('different pages display different policies', async function (assert) {
    // Regression guard: page 1 and page 2 must show distinct, non-overlapping
    // rows. With DEFAULT_PAGE_SIZE=15 in test env, policies 0-14 are on page 1
    // and policies 15-19 are on page 2.
    const manyPolicies = Array.from({ length: 20 }, (_, i) => `policy-${i}`);
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: manyPolicies },
      request_id: 'test',
    }));

    await visit('/vault/policies/acl');

    // Capture the first policy name visible on page 1 (policy-0).
    assert.dom(GENERAL.listItem('policy-0')).exists('policy-0 is visible on page 1');
    assert.dom(GENERAL.listItem('policy-15')).doesNotExist('policy-15 is not visible on page 1');

    // Navigate to page 2.
    await click(GENERAL.nextPage);

    // Page 2 must show policy-15 and must NOT show policy-0.
    assert.dom(GENERAL.listItem('policy-15')).exists('policy-15 is visible on page 2');
    assert.dom(GENERAL.listItem('policy-0')).doesNotExist('policy-0 is not visible on page 2');
  });

  // ── Sorting ─────────────────────────────────────────────────────────────────

  test('sorting descending sorts the global dataset before pagination', async function (assert) {
    // Regression for VAULT-50816: sorting must apply to the full dataset, not
    // just the items on the current page. With 20 policies (policy-0..policy-19)
    // sorted Z→A with natural sorting, the globally last item numerically ('policy-19')
    // must appear on page 1 — not stranded on page 2 as it was before the fix.
    const manyPolicies = Array.from({ length: 20 }, (_, i) => `policy-${i}`);
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: manyPolicies },
      request_id: 'test',
    }));

    await visit('/vault/policies/acl');

    // Click the "Policy name" column sort button once (asc), then again (desc).
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));

    // Under natural descending sort, "policy-19" is globally first.
    // It must be visible on page 1.
    assert.dom(GENERAL.listItem('policy-19')).exists('policy-19 is on page 1 when globally sorted Z→A');

    // "policy-0" sorts last and must NOT be visible on page 1.
    assert
      .dom(GENERAL.listItem('policy-0'))
      .doesNotExist('policy-0 is not on page 1 when sorted Z→A (belongs on a later page)');
  });

  test('sortBy and sortOrder query params persist across page transitions', async function (assert) {
    // Regression for VAULT-50816: the active sort column/direction must survive
    // the model refresh triggered when the page query param changes.  Before the
    // fix, navigating to page 2 reset the sort indicator back to unsorted (⇅).
    const manyPolicies = Array.from({ length: 20 }, (_, i) => `policy-${i}`);
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: manyPolicies },
      request_id: 'test',
    }));

    await visit('/vault/policies/acl');

    // Sort ascending — first click on the sortable "Policy name" header.
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));

    // The URL must now include sortBy and sortOrder.
    assert.true(currentURL().includes('sortBy=name'), 'sortBy=name is in the URL');
    assert.true(currentURL().includes('sortOrder=asc'), 'sortOrder=asc is in the URL');

    // Navigate to page 2.
    await click(GENERAL.nextPage);

    // Sort params must still be present after the page transition.
    assert.true(currentURL().includes('sortBy=name'), 'sortBy=name remains after navigating to page 2');
    assert.true(currentURL().includes('sortOrder=asc'), 'sortOrder=asc remains after navigating to page 2');

    // The column header sort indicator must still be active (aria-sort="ascending").
    assert
      .dom('.hds-advanced-table__th:nth-child(1)')
      .hasAttribute('aria-sort', 'ascending', 'sort indicator remains active on page 2');
  });

  test('sortBy and sortOrder query params reset on navigation away', async function (assert) {
    // Verifies resetController clears sort QPs when leaving the route so a
    // subsequent visit starts unsorted.
    await visit('/vault/policies/acl?sortBy=name&sortOrder=asc');
    await visit('/vault/dashboard');
    await visit('/vault/policies/acl');
    assert.false(currentURL().includes('sortBy='), 'sortBy was reset after navigation');
    assert.false(currentURL().includes('sortOrder='), 'sortOrder was reset after navigation');
  });

  // ── Regression: VAULT-50826 ──────────────────────────────────────────────

  test('it navigates to the correct ACL policy edit page from the list', async function (assert) {
    // Regression guard for VAULT-50826: clicking "Edit policy" on a list row must
    // navigate to THAT policy's edit page, not the one last loaded by the show route.
    //
    // Root cause: policy/edit calls this.modelFor('vault.cluster.policy.show'),
    // which returns route.currentModel — a value Ember retains on the route singleton
    // after navigating away. When the show route was previously activated for a
    // different policy, edit.model() blindly reuses that stale object.
    //
    // Minimum reproduction: show(policy-alpha) → list → edit(policy-beta).
    this.server.get('sys/policies/acl/', () => ({
      data: { keys: ['policy-alpha', 'policy-beta'] },
    }));
    this.server.get('sys/policies/acl/:name', (_schema, request) => ({
      data: {
        name: request.params.name,
        policy: `path "secret/*" { capabilities = ["read"] }`,
      },
    }));

    // Step 1: visit the show page for policy-alpha so show.currentModel is set on
    // the route singleton to { name: 'policy-alpha', ... }.
    await visit('/vault/policy/acl/policy-alpha');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('policy-alpha', 'show page loaded for policy-alpha');

    // Step 2: navigate to the list (show deactivates, but currentModel persists).
    await click(GENERAL.breadcrumbLink('ACL policies'));

    // Step 3: click "Edit policy" on policy-beta.
    // show.currentModel is still policy-alpha — the edit route must not reuse it.
    const triggers = document.querySelectorAll(GENERAL.menuTrigger);
    await click(triggers[1]);
    await click(GENERAL.menuItem('edit-policy'));

    assert.true(currentURL().includes('/vault/policy/acl/policy-beta/edit'), 'URL is for policy-beta');
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText('policy-beta', 'edit page shows policy-beta, not the stale policy-alpha');
  });
});
