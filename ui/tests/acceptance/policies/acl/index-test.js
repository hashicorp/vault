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
});
