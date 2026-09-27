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

// Two RGP policies that will appear in the list.
const STUB_POLICIES = ['my-rgp-policy', 'another-rgp-policy'];

module('Acceptance | RGP policies list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    // Stub the version service so the route treats this as an Enterprise instance
    // with Sentinel enabled — otherwise the route returns null and the page is empty.
    const version = this.owner.lookup('service:version');
    version.features = ['Sentinel'];
    // Stub the list endpoint used by the route model.
    // The API client calls GET /v1/sys/policies/rgp/ (trailing slash) + ?list=true.
    this.server.get('sys/policies/rgp/', () => ({
      data: { keys: STUB_POLICIES },
      request_id: 'test',
    }));
    // Stub the capabilities endpoint so capability checks don't fail.
    // An empty data object causes all capability checks to default to true.
    this.server.post('sys/capabilities-self', () => ({
      data: {},
      request_id: 'test',
    }));
  });

  // ── Page structure ───────────────────────────────────────────────────────

  test('it renders the page title', async function (assert) {
    // Verifies Page::ListView renders the title from listViewConfig.
    await visit('/vault/policies/rgp');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('RGP policies');
  });

  test('it renders breadcrumbs', async function (assert) {
    // Verifies both breadcrumb entries are present on the policies index page.
    await visit('/vault/policies/rgp');
    assert.dom(GENERAL.breadcrumbLink('Vault')).exists('first breadcrumb links to Vault dashboard');
    assert.dom(GENERAL.currentBreadcrumb('RGP policies')).exists('last breadcrumb is current page');
  });

  test('it renders the Create RGP policy primary action button', async function (assert) {
    // Verifies the primary action button is rendered by Page::ListView.
    await visit('/vault/policies/rgp');
    assert.dom(GENERAL.button('Create RGP policy')).exists();
  });

  // ── Table columns ────────────────────────────────────────────────────────

  test('it renders the Policy name column header', async function (assert) {
    // Verifies the Policy name column is rendered in the table header.
    // ListTable uses Hds::AdvancedTable so the isAdvanced flag is required.
    await visit('/vault/policies/rgp');
    assert.dom(GENERAL.tableColumnHeader(1, { isAdvanced: true })).includesText('Policy name');
  });

  // ── Name column rendering ────────────────────────────────────────────────

  test('it renders policy names as links', async function (assert) {
    // Verifies that RGP policies show a navigable link in the Name column.
    await visit('/vault/policies/rgp');
    assert.dom('[data-test-policy-link="my-rgp-policy"]').exists('policy renders as a link');
  });

  // ── Filter ───────────────────────────────────────────────────────────────

  test('it shows the filtered empty state when no policies match the filter', async function (assert) {
    // Verifies the filteredEmptyTitle is displayed when filter text matches nothing.
    await visit('/vault/policies/rgp');
    await fillIn(GENERAL.filterInput, 'nonexistent-policy-xyz');
    assert.dom(GENERAL.emptyStateTitle).includesText('No results for');
  });

  // ── Empty state — no policies ────────────────────────────────────────────

  test('it shows the empty state when there are no RGP policies', async function (assert) {
    // Verifies the noDataTitle is displayed when the list is empty.
    this.server.get('sys/policies/rgp/', () => ({
      data: { keys: [] },
      request_id: 'test',
    }));
    await visit('/vault/policies/rgp');
    assert.dom(GENERAL.emptyStateTitle).hasText('No RGP policies yet');
  });

  // ── Sentinel gate ────────────────────────────────────────────────────────

  test('it shows the upgrade page when Sentinel is not available', async function (assert) {
    // When the "Sentinel" feature flag is absent the route template renders
    // <UpgradePage> instead of the list. The has-feature helper reads from
    // version.features which is populated by GET /sys/license/features on every
    // cluster route transition. Stubbing that endpoint to return an empty feature
    // list and clearing the cached value forces the service to re-fetch and pick
    // up the empty response, so has-feature "Sentinel" evaluates to false.
    const version = this.owner.lookup('service:version');
    version.features = [];
    this.server.get('/sys/license/features', () => ({ features: [] }));
    await visit('/vault/policies/rgp');
    assert.dom(GENERAL.emptyStateTitle).includesText('Upgrade to use');
  });

  // ── Row actions ──────────────────────────────────────────────────────────

  test('it renders popup menu actions for a policy', async function (assert) {
    // Verifies View, Edit, Download, and Delete actions are present.
    this.server.get('sys/policies/rgp/', () => ({
      data: { keys: ['my-rgp-policy'] },
      request_id: 'test',
    }));
    this.server.post('sys/capabilities-self', () => ({
      data: { 'sys/policies/rgp/my-rgp-policy': ['read', 'update', 'delete'] },
      request_id: 'test',
    }));
    await visit('/vault/policies/rgp');
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('view-policy')).exists('View policy action exists');
    assert.dom(GENERAL.menuItem('edit-policy')).exists('Edit policy action exists');
    assert.dom(GENERAL.menuItem('download-policy')).exists('Download policy action exists');
    assert.dom(GENERAL.menuItem('delete-policy')).exists('Delete policy action exists');
  });

  // ── Delete modal ─────────────────────────────────────────────────────────

  test('it shows a confirmation modal when Delete is clicked', async function (assert) {
    // Verifies the ConfirmModal appears when the critical Delete action is triggered.
    this.server.get('sys/policies/rgp/', () => ({
      data: { keys: ['my-rgp-policy'] },
      request_id: 'test',
    }));
    this.server.post('sys/capabilities-self', () => ({
      data: { 'sys/policies/rgp/my-rgp-policy': ['read', 'update', 'delete'] },
      request_id: 'test',
    }));
    await visit('/vault/policies/rgp');
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
    const manyPolicies = Array.from({ length: 20 }, (_, i) => `rgp-policy-${i}`);
    this.server.get('sys/policies/rgp/', () => ({
      data: { keys: manyPolicies },
      request_id: 'test',
    }));
    await visit('/vault/policies/rgp');
    await click(GENERAL.nextPage);
    assert.true(currentURL().includes('page=2'), 'URL contains page=2 after clicking next');
  });

  test('it resets the page query param to 1 on navigation away', async function (assert) {
    // Verifies resetController resets the page QP when navigating away.
    await visit('/vault/policies/rgp?page=3');
    await visit('/vault/dashboard');
    await visit('/vault/policies/rgp');
    assert.false(currentURL().includes('page=3'), 'page param was reset after navigation');
  });

  test('different pages display different policies', async function (assert) {
    // Regression guard: page 1 and page 2 must show distinct, non-overlapping rows.
    const manyPolicies = Array.from({ length: 20 }, (_, i) => `rgp-policy-${i}`);
    this.server.get('sys/policies/rgp/', () => ({
      data: { keys: manyPolicies },
      request_id: 'test',
    }));

    await visit('/vault/policies/rgp');

    assert.dom(GENERAL.listItem('rgp-policy-0')).exists('rgp-policy-0 is visible on page 1');
    assert.dom(GENERAL.listItem('rgp-policy-15')).doesNotExist('rgp-policy-15 is not visible on page 1');

    await click(GENERAL.nextPage);

    assert.dom(GENERAL.listItem('rgp-policy-15')).exists('rgp-policy-15 is visible on page 2');
    assert.dom(GENERAL.listItem('rgp-policy-0')).doesNotExist('rgp-policy-0 is not visible on page 2');
  });
});
