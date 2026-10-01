/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, visit, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';

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

  test('it filters groups by full or partial group ID', async function (assert) {
    await visit('/vault/access/identity/groups');
    await fillIn(GENERAL.filterInput, GROUP_ID_2);
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'only the matching group renders');
    assert.dom('[data-test-identity-link="Florida_design"]').exists('group matching the full ID renders');

    await fillIn(GENERAL.filterInput, GROUP_ID_1.slice(0, 8));
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'only the matching group renders');
    assert.dom('[data-test-identity-link="Kochi_design"]').exists('group matching the partial ID renders');
  });

  test('it filters groups by name', async function (assert) {
    await visit('/vault/access/identity/groups');
    await fillIn(GENERAL.filterInput, 'kochi');
    assert.dom(GENERAL.tableRow()).exists({ count: 1 }, 'only the matching group renders');
    assert.dom('[data-test-identity-link="Kochi_design"]').exists();
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

  // ── Sorting ──────────────────────────────────────────────────────────────

  test('sorting by name applies to the full dataset before pagination', async function (assert) {
    // 20 groups named group-00..group-19; ascending sort should put group-00 on page 1.
    // Descending sort should push group-19 to page 1 and group-00 off it.
    const manyGroups = Array.from({ length: 20 }, (_, i) => ({
      id: `group-id-${i}`,
      name: `group-${String(i).padStart(2, '0')}`,
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

    // Click the Group name sort button (column 1) once → ascending
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));
    await waitFor('[data-test-identity-link]');

    assert.dom('[data-test-identity-link="group-00"]').exists('group-00 is on page 1 after ascending sort');
    assert
      .dom('[data-test-identity-link="group-19"]')
      .doesNotExist('group-19 is not on page 1 after ascending sort');

    // Click again → descending; group-19 should now be on page 1
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));
    await waitFor('[data-test-identity-link]');

    assert.dom('[data-test-identity-link="group-19"]').exists('group-19 is on page 1 after descending sort');
    assert
      .dom('[data-test-identity-link="group-00"]')
      .doesNotExist('group-00 is not on page 1 after descending sort');
  });

  test('sort query params persist when navigating to a different page', async function (assert) {
    // After sorting then changing page, sortBy and sortOrder must remain in the URL.
    const manyGroups = Array.from({ length: 20 }, (_, i) => ({
      id: `group-id-${i}`,
      name: `group-${String(i).padStart(2, '0')}`,
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
    await click(GENERAL.tableColumnHeaderSortButton(1, { isAdvanced: true }));
    await click(GENERAL.nextPage);

    const url = currentURL();
    assert.true(url.includes('sortBy=name'), 'sortBy=name persists after page change');
    assert.true(url.includes('sortOrder=asc'), 'sortOrder=asc persists after page change');
    assert.true(url.includes('page=2'), 'page=2 is set');
  });

  test('sort query params reset when navigating away from the page', async function (assert) {
    await visit('/vault/access/identity/groups?sortBy=name&sortOrder=asc');
    await visit('/vault/dashboard');
    await visit('/vault/access/identity/groups');

    const url = currentURL();
    assert.false(url.includes('sortBy='), 'sortBy param is cleared after navigation away');
    assert.false(url.includes('sortOrder='), 'sortOrder param is cleared after navigation away');
  });
});

// Read-only tokens must not see create or edit actions for groups.
module('Acceptance | Identity groups list view | read-only token', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    // External, so the details view would otherwise offer "Add alias".
    this.groupId = await runCmd('write -field=id identity/group name=read-only-group type=external');
    const token = await runCmd(
      tokenWithPolicyCmd('groups-read-only', 'path "identity/*" { capabilities = ["read", "list"] }')
    );
    await login(token);
  });

  hooks.afterEach(async function () {
    await login();
    await runCmd('delete identity/group/name/read-only-group');
  });

  test('it hides the create action on the list view', async function (assert) {
    await visit('/vault/access/identity/groups');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Groups', 'list view renders');
    assert.dom('[data-test-identity-link="read-only-group"]').exists('group is listed');
    assert.dom(GENERAL.button('Create group')).doesNotExist('create action is hidden');
  });

  test('it hides the edit and add alias actions on the group details view', async function (assert) {
    await visit(`/vault/access/identity/groups/${this.groupId}/details`);
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('read-only-group', 'details view renders');
    assert.dom('[data-test-entity-edit-link]').doesNotExist('edit action is hidden');
    assert.dom('[data-test-entity-create-link]').doesNotExist('add alias action is hidden');
  });
});

// Capability gating must not hide group details actions from a user who is allowed to use them.
module('Acceptance | Identity group details | root token', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    this.groupId = await runCmd('write -field=id identity/group name=root-external-group type=external');
  });

  hooks.afterEach(async function () {
    await runCmd('delete identity/group/name/root-external-group');
  });

  test('it shows the edit and add alias actions on the group details view', async function (assert) {
    await visit(`/vault/access/identity/groups/${this.groupId}/details`);
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('root-external-group', 'details view renders');
    assert.dom('[data-test-entity-edit-link]').exists('edit action is shown');
    assert.dom('[data-test-entity-create-link]').exists('add alias action is shown');
  });
});
