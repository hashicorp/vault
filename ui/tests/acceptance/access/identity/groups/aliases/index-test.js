/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const ALIAS_ID_1 = 'aaaabbbb-1d5-5a70-0f3c-4ec61f7b1111';
const ALIAS_ID_2 = 'ccccdddd-1d5-9b5c-0f3c-4ec61f7b2222';

const STUB_ALIASES = [
  {
    id: ALIAS_ID_1,
    name: 'Sample_1256',
    mount_type: 'approle',
    mount_accessor: 'auth_approle_427f93376',
  },
  {
    id: ALIAS_ID_2,
    name: 'test_155',
    mount_type: 'userrole',
    mount_accessor: 'auth_userrole_427f93376',
  },
];

function aliasListResponse(aliases) {
  const key_info = {};
  const keys = [];
  for (const a of aliases) {
    key_info[a.id] = {
      name: a.name,
      mount_type: a.mount_type,
      mount_accessor: a.mount_accessor,
    };
    keys.push(a.id);
  }
  return { data: { key_info, keys }, request_id: 'test' };
}

module('Acceptance | Identity group aliases list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    this.server.get('/identity/group-alias/id', () => aliasListResponse(STUB_ALIASES));
    // Empty capabilities data causes all permission checks to default to allowed.
    this.server.post('/sys/capabilities-self', () => ({ data: {}, request_id: 'test' }));
  });

  // ── Page structure ──────────────────────────────────────────────────────

  test('it renders the page title', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Groups');
  });

  test('it renders breadcrumbs', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    assert.dom(GENERAL.breadcrumbLink('Vault')).exists('first breadcrumb links to Vault dashboard');
    assert.dom(GENERAL.currentBreadcrumb('Groups')).exists('last breadcrumb is current page');
  });

  test('it renders the Create group primary action button', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    assert.dom(GENERAL.button('Create group')).exists();
  });

  // ── Tab navigation ───────────────────────────────────────────────────────

  test('it renders Groups and Aliases tabs', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    assert.dom('[data-test-tab="Groups"]').exists('Groups tab exists');
    assert.dom('[data-test-tab="Aliases"]').exists('Aliases tab exists');
  });

  // ── Table columns ────────────────────────────────────────────────────────

  test('it renders the Aliases name, Aliases ID, and Mount type column headers', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    assert.dom(GENERAL.tableColumnHeader(1, { isAdvanced: true })).includesText('Aliases name');
    assert.dom(GENERAL.tableColumnHeader(2, { isAdvanced: true })).includesText('Aliases ID');
    assert.dom(GENERAL.tableColumnHeader(3, { isAdvanced: true })).includesText('Mount type');
  });

  // ── Filter ───────────────────────────────────────────────────────────────

  test('it shows the filtered empty state when no aliases match the filter', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    await fillIn(GENERAL.filterInput, 'nonexistent-alias-xyz');
    assert.dom(GENERAL.emptyStateTitle).includesText('No results for');
  });

  // ── Empty state ──────────────────────────────────────────────────────────

  test('it shows the empty state when there are no group aliases', async function (assert) {
    this.server.get('/identity/group-alias/id', () => ({
      data: { key_info: {}, keys: [] },
      request_id: 'test',
    }));
    await visit('/vault/access/identity/groups/aliases');
    assert.dom(GENERAL.emptyStateTitle).hasText('No group aliases yet');
  });

  // ── Row actions ──────────────────────────────────────────────────────────

  test('it renders popup menu actions for an alias', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    await click(GENERAL.menuItem(STUB_ALIASES[0].name));
    assert.dom(GENERAL.menuItem('view-details')).exists('View details action exists');
    assert.dom(GENERAL.menuItem('edit-details')).exists('Edit details action exists');
    assert.dom(GENERAL.menuItem('delete')).exists('Remove alias action exists');
  });

  // ── Remove modal ─────────────────────────────────────────────────────────

  test('it shows a confirmation modal when Remove alias is clicked', async function (assert) {
    await visit('/vault/access/identity/groups/aliases');
    await click(GENERAL.menuItem(STUB_ALIASES[0].name));
    await click(GENERAL.menuItem('delete'));
    assert.dom(GENERAL.confirmTitle).hasText('Remove this alias');
  });

  // ── Pagination ───────────────────────────────────────────────────────────

  test('it advances the page query param when the next-page button is clicked', async function (assert) {
    const manyAliases = Array.from({ length: 20 }, (_, i) => ({
      id: `alias-id-${i}`,
      name: `alias-${i}`,
      mount_type: 'approle',
      mount_accessor: `auth_approle_${i}`,
    }));
    this.server.get('/identity/group-alias/id', () => aliasListResponse(manyAliases));
    await visit('/vault/access/identity/groups/aliases');
    await click(GENERAL.nextPage);
    assert.true(currentURL().includes('page=2'), 'URL contains page=2 after clicking next');
  });

  test('it resets the page query param to 1 on navigation away', async function (assert) {
    await visit('/vault/access/identity/groups/aliases?page=3');
    await visit('/vault/dashboard');
    await visit('/vault/access/identity/groups/aliases');
    assert.false(currentURL().includes('page=3'), 'page param was reset after navigation');
  });
});
