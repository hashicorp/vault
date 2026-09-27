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

const GROUP_ID_1 = '66638b30-a05e-560b-18cc-f43af766ce73';
const GROUP_ID_2 = '78938b30-a85e-535b-14ac-f43af766ce73';

const STUB_GROUPS = [
  { id: GROUP_ID_1, name: 'Kochi_design' },
  { id: GROUP_ID_2, name: 'Florida_design' },
];

function groupListResponse(groups) {
  const key_info = {};
  const keys = [];
  for (const g of groups) {
    key_info[g.id] = { name: g.name };
    keys.push(g.id);
  }
  return { data: { key_info, keys }, request_id: 'test' };
}

module('Acceptance | Identity groups list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    this.server.get('/identity/group/id', () => groupListResponse(STUB_GROUPS));
    this.server.get('/identity/group/id/:id', (_, req) => {
      const group = STUB_GROUPS.find((g) => g.id === req.params.id) ?? STUB_GROUPS[0];
      return {
        data: { id: group.id, name: group.name, type: 'internal', alias: null, policies: [] },
        request_id: 'test',
      };
    });
    // Empty capabilities data causes all permission checks to default to allowed.
    this.server.post('/sys/capabilities-self', () => ({ data: {}, request_id: 'test' }));
  });

  // ── Page structure ──────────────────────────────────────────────────────

  test('it renders the page title', async function (assert) {
    await visit('/vault/access/identity/groups');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Groups');
  });

  test('it renders breadcrumbs', async function (assert) {
    await visit('/vault/access/identity/groups');
    assert.dom(GENERAL.breadcrumbLink('Vault')).exists('first breadcrumb links to Vault dashboard');
    assert.dom(GENERAL.currentBreadcrumb('Groups')).exists('last breadcrumb is current page');
  });

  test('it renders the Create group primary action button', async function (assert) {
    await visit('/vault/access/identity/groups');
    assert.dom(GENERAL.button('Create group')).exists();
  });

  // ── Table columns ────────────────────────────────────────────────────────

  test('it renders the Group name and Group ID column headers', async function (assert) {
    await visit('/vault/access/identity/groups');
    assert.dom(GENERAL.tableColumnHeader(1, { isAdvanced: true })).includesText('Group name');
    assert.dom(GENERAL.tableColumnHeader(2, { isAdvanced: true })).includesText('Group ID');
  });

  // ── Filter ───────────────────────────────────────────────────────────────

  test('it shows the filtered empty state when no groups match the filter', async function (assert) {
    await visit('/vault/access/identity/groups');
    await fillIn(GENERAL.filterInput, 'nonexistent-group-xyz');
    assert.dom(GENERAL.emptyStateTitle).includesText('No results for');
  });

  // ── Empty state ──────────────────────────────────────────────────────────

  test('it shows the empty state when there are no groups', async function (assert) {
    this.server.get('/identity/group/id', () => ({
      data: { key_info: {}, keys: [] },
      request_id: 'test',
    }));
    await visit('/vault/access/identity/groups');
    assert.dom(GENERAL.emptyStateTitle).hasText('No groups yet');
  });

  // ── Row actions ──────────────────────────────────────────────────────────

  test('it renders popup menu actions for a group', async function (assert) {
    await visit('/vault/access/identity/groups');
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('view-details')).exists('View details action exists');
    assert.dom(GENERAL.menuItem('edit-details')).exists('Edit details action exists');
    assert.dom(GENERAL.menuItem('delete-group')).exists('Delete group action exists');
  });

  // ── Delete modal ─────────────────────────────────────────────────────────

  test('it shows a confirmation modal when Delete group is clicked', async function (assert) {
    await visit('/vault/access/identity/groups');
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete-group'));
    assert.dom(GENERAL.confirmTitle).hasText('Delete this group?');
  });

  // ── Pagination ───────────────────────────────────────────────────────────

  test('it advances the page query param when the next-page button is clicked', async function (assert) {
    const manyGroups = Array.from({ length: 20 }, (_, i) => ({
      id: `group-id-${i}`,
      name: `group-${i}`,
    }));
    this.server.get('/identity/group/id', () => groupListResponse(manyGroups));
    this.server.get('/identity/group/id/:id', (_, req) => ({
      data: {
        id: req.params.id,
        name: `group-${req.params.id}`,
        type: 'internal',
        alias: null,
        policies: [],
      },
      request_id: 'test',
    }));
    await visit('/vault/access/identity/groups');
    await click(GENERAL.nextPage);
    assert.true(currentURL().includes('page=2'), 'URL contains page=2 after clicking next');
  });

  test('it resets the page query param to 1 on navigation away', async function (assert) {
    await visit('/vault/access/identity/groups?page=3');
    await visit('/vault/dashboard');
    await visit('/vault/access/identity/groups');
    assert.false(currentURL().includes('page=3'), 'page param was reset after navigation');
  });

  test('different pages display different groups', async function (assert) {
    const manyGroups = Array.from({ length: 20 }, (_, i) => ({
      id: `group-id-${i}`,
      name: `group-${i}`,
    }));
    this.server.get('/identity/group/id', () => groupListResponse(manyGroups));
    this.server.get('/identity/group/id/:id', (_, req) => ({
      data: {
        id: req.params.id,
        name: `group-${req.params.id}`,
        type: 'internal',
        alias: null,
        policies: [],
      },
      request_id: 'test',
    }));

    await visit('/vault/access/identity/groups');

    assert.dom('[data-test-identity-link="group-0"]').exists('group-0 is visible on page 1');
    assert.dom('[data-test-identity-link="group-15"]').doesNotExist('group-15 is not visible on page 1');

    await click(GENERAL.nextPage);

    assert.dom('[data-test-identity-link="group-15"]').exists('group-15 is visible on page 2');
    assert.dom('[data-test-identity-link="group-0"]').doesNotExist('group-0 is not visible on page 2');
  });
});
