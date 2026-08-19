/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import agentRegistryHandler from 'vault/mirage/handlers/agent-registry';
import { visit, click, fillIn } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { Response } from 'miragejs';
import { overrideResponse } from 'vault/tests/helpers/stubs';

const ROUTE = '/vault/agents/registry';

module('Acceptance | agents | registry', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    agentRegistryHandler(this.server);
    return login();
  });

  test('it shows the wizard, dismisses it, then shows the empty state', async function (assert) {
    // No agents seeded — wizard should appear automatically
    await visit(ROUTE);

    assert.dom(GENERAL.wizardIntro).exists('wizard intro card is shown when no agents are registered');
    assert
      .dom(GENERAL.wizardIntro)
      .includesText(
        'Register and govern AI agent identities in Vault',
        'wizard shows the correct heading text'
      );
    assert.dom(GENERAL.button('Skip')).exists('Skip button is present');
    assert.dom(GENERAL.button('enable')).exists('"Register via CLI" button is present');

    // Dismiss the wizard via the Skip button
    await click(GENERAL.button('Skip'));

    assert.dom(GENERAL.wizardIntro).doesNotExist('wizard is hidden after dismissal');
    assert.dom(GENERAL.emptyStateTitle).hasText('No agents yet', 'empty state title renders');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Your agents will be listed here. Register an agent in CLI or API to get started.',
        'empty state message renders'
      );
  });

  test('it filters agents by name and paginates through the list', async function (assert) {
    // Create enough agents to fill two pages (default page size is 10)
    this.server.createList('agent-registry-registration', 12);
    await visit(ROUTE);

    // All 12 agents exist — first page shows the first 10 (pagination control is rendered)
    assert.dom(GENERAL.pagination).exists('pagination control is visible');
    assert.dom(GENERAL.paginationInfo).includesText('1–10', 'first page shows rows 1–10 out of 12 total');
    assert.dom(GENERAL.button('agent agent-1')).exists('agent-1 is on the first page');
    assert.dom(GENERAL.button('agent agent-11')).doesNotExist('agent-11 is not on the first page');

    // Advance to page 2 using the next-page arrow
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.paginationInfo).includesText('11–12', 'second page shows rows 11–12');
    assert.dom(GENERAL.button('agent agent-11')).exists('agent-11 is on the second page');
    assert.dom(GENERAL.button('agent agent-12')).exists('agent-12 is on the second page');

    // Return to page 1 and filter by name to a single result
    await click(GENERAL.prevPage);
    await fillIn('[aria-label="Filter agent registry table"]', 'agent-3');

    assert.dom(GENERAL.button('agent agent-3')).exists('matching agent is visible after filtering');
    assert.dom(GENERAL.button('agent agent-1')).doesNotExist('non-matching agent is hidden after filtering');

    // Clearing the filter restores the full list
    await fillIn('[aria-label="Filter agent registry table"]', '');
    assert.dom(GENERAL.button('agent agent-1')).exists('agent-1 is visible again after clearing filter');
  });

  test('it goes through the delete workflow', async function (assert) {
    this.server.createList('agent-registry-registration', 2);

    // Make the first delete attempt fail
    this.server.delete('agent-registry/registration/display-name/:name', () => overrideResponse(403));

    await visit(ROUTE);

    assert.dom(GENERAL.button('agent agent-1')).exists('agent-1 is present before deletion');

    // First attempt — server rejects the request
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete'));
    assert.dom(GENERAL.confirmModal).exists('confirm modal appears');
    await click(GENERAL.confirmButton);

    assert.dom(GENERAL.button('agent agent-1')).exists('agent-1 remains in the list after a failed delete');
    assert
      .dom(GENERAL.latestFlashContent)
      .includesText('Error deleting agent', 'error flash message is shown when the API returns 403');

    // Remove the override so the next attempt succeeds
    this.server.delete('agent-registry/registration/display-name/:name', (schema, request) => {
      const { name } = request.params;
      const registration = schema.db.agentRegistryRegistrations.findBy({ display_name: name });
      schema.db.agentRegistryRegistrations.remove(registration.id);
      return new Response(204);
    });

    // Second attempt — server accepts the request
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete'));
    assert.dom(GENERAL.confirmModal).exists('confirm modal appears again');
    await click(GENERAL.confirmButton);

    assert
      .dom(GENERAL.button('agent agent-1'))
      .doesNotExist('agent-1 is removed from the list after a successful delete');
    assert
      .dom(GENERAL.latestFlashContent)
      .includesText('Successfully deleted agent', 'success flash message is shown');
  });

  test('it goes through the details flyout tabs', async function (assert) {
    // Create one agent with a fixed entity_id so mirage overrides are predictable
    this.server.create('agent-registry-registration', {
      id: 'reg-1',
      display_name: 'agent-1',
      entity_id: 'entity-abc',
      owner: 'owner@example.com',
      ceiling_policies: ['default'],
    });

    // Override entity read to return a deterministic entity with one alias
    this.server.get('identity/entity/id/entity-abc', () => ({
      data: {
        id: 'entity-abc',
        name: 'entity-0',
        disabled: false,
        policies: ['entity-policy'],
        group_ids: [],
        merged_entity_ids: null,
        metadata: { organization: 'test-org' },
        creation_time: new Date().toISOString(),
        last_update_time: new Date().toISOString(),
        aliases: [
          {
            id: 'alias-abc',
            name: 'alias-0',
            mount_type: 'userpass',
            mount_path: 'auth/userpass/',
            mount_accessor: 'auth_userpass_abc',
            creation_time: new Date().toISOString(),
            last_update_time: new Date().toISOString(),
            custom_metadata: { contact_email: 'alias@example.com' },
          },
        ],
      },
    }));

    await visit(ROUTE);

    // Open the flyout by clicking the agent name link
    await click(GENERAL.button('agent agent-1'));
    assert.dom(GENERAL.flyout).exists('flyout opens when an agent name is clicked');
    assert.dom(GENERAL.flyout).includesText('agent-1', 'flyout header shows the agent name');

    // --- Agent details tab (selected by default) ---
    assert.dom(GENERAL.hdsTab('agent')).exists('Agent details tab is present');
    assert.dom(GENERAL.hdsTabPanel('agent')).exists('Agent details panel is visible by default');
    assert.dom(GENERAL.cardContainer('agent-status')).exists('status card renders on agent tab');
    assert.dom(GENERAL.cardContainer('agent-details')).exists('agent id card renders on agent tab');
    assert
      .dom(GENERAL.cardContainer('agent-status'))
      .includesText('Enabled', 'status card shows entity state');
    assert
      .dom(GENERAL.cardContainer('agent-status'))
      .includesText('owner@example.com', 'status card shows owner');

    // --- Entity details tab ---
    await click(GENERAL.hdsTab('entity'));
    assert.dom(GENERAL.hdsTabPanel('entity')).exists('Entity details panel is visible');
    assert.dom(GENERAL.cardContainer('merged-ids')).exists('merged IDs card renders on entity tab');
    assert.dom(GENERAL.widget('entity-metadata')).exists('entity metadata section renders on entity tab');

    // --- Policies tab ---
    await click(GENERAL.hdsTab('policies'));
    assert.dom(GENERAL.hdsTabPanel('policies')).exists('Policies panel is visible');
    assert.dom(GENERAL.table('policy-list')).exists('policy list table renders on policies tab');
    assert.dom(GENERAL.table('policy-list')).includesText('default', 'ceiling policy row renders');
    assert.dom(GENERAL.table('policy-list')).includesText('entity-policy', 'entity policy row renders');

    // --- Alias details tab (only shown when entity has aliases) ---
    await click(GENERAL.hdsTab('alias'));
    assert.dom(GENERAL.hdsTabPanel('alias')).exists('Alias details panel is visible');
    assert.dom(GENERAL.cardContainer('alias-navigation')).exists('alias navigation card renders');
    assert.dom(GENERAL.cardContainer('alias-auth-method')).exists('alias auth method card renders');
    assert
      .dom(GENERAL.cardContainer('alias-auth-method'))
      .includesText('userpass', 'alias auth method card shows mount type');

    // Close the flyout and verify it is dismissed
    await click(GENERAL.button('Close details'));
    assert.dom(GENERAL.flyout).doesNotExist('flyout is dismissed after clicking Close details');

    // Override the entity to have no aliases, then re-open the flyout
    this.server.get('identity/entity/id/entity-abc', () => ({
      data: {
        id: 'entity-abc',
        name: 'entity-0',
        disabled: false,
        policies: ['entity-policy'],
        group_ids: [],
        merged_entity_ids: null,
        metadata: {},
        creation_time: new Date().toISOString(),
        last_update_time: new Date().toISOString(),
        aliases: [],
      },
    }));
    await click(GENERAL.button('agent agent-1'));
    assert.dom(GENERAL.hdsTab('alias')).doesNotExist('alias tab is hidden when entity has no aliases');
  });
});
