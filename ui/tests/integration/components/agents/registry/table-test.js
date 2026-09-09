/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { dateFormat } from 'core/helpers/date-format';

const findAllParentRowIds = () =>
  findAll(GENERAL.tableParentRow).map((row) => row.getAttribute('data-test-table-row'));

module('Integration | Component | agents/registry-table', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.api = this.owner.lookup('service:api');
    this.flashMessages = this.owner.lookup('service:flash-messages');

    // Stub router methods
    this.transitionToStub = sinon.stub(this.router, 'transitionTo');
    this.refreshStub = sinon.stub(this.router, 'refresh');

    // Stub API methods
    this.entityUpdateStub = sinon.stub(this.api.identity, 'entityUpdateById').resolves();
    this.registrationDeleteStub = sinon.stub(this.api.secrets, 'registrationDeleteByName').resolves();
    this.parseErrorStub = sinon.stub(this.api, 'parseError').resolves({ message: 'Test error message' });

    // Stub flash messages
    this.flashSuccessStub = sinon.stub(this.flashMessages, 'success');
    this.flashDangerStub = sinon.stub(this.flashMessages, 'danger');

    // Stub on click handler
    this.onClick = sinon.stub();

    // Mock data
    const agents = [
      {
        id: 'agent-1',
        display_name: 'test-agent-1',
        entity_id: 'entity-1',
        entity: {
          id: 'entity-1',
          name: 'test-entity-1',
          disabled: false,
          creation_time: '2026-01-01T10:00:00Z',
          last_update_time: '2026-01-02T11:00:00Z',
          aliases: [
            {
              id: 'alias-1',
              name: 'test-alias-1',
              mount_accessor: 'accessor-1',
              mount_path: 'token/',
              mount_type: 'token',
              creation_time: '2026-01-03T12:00:00Z',
              last_update_time: '2026-01-04T13:00:00Z',
            },
            {
              id: 'alias-2',
              name: 'test-alias-2',
              mount_accessor: 'accessor-2',
              mount_path: 'approle/',
              mount_type: 'approle',
              creation_time: '2026-01-05T14:00:00Z',
              last_update_time: '2026-01-06T15:00:00Z',
            },
          ],
        },
      },
      {
        id: 'agent-2',
        display_name: 'test-agent-2',
        entity_id: 'entity-2',
        entity: {
          id: 'entity-2',
          name: 'test-entity-2',
          disabled: true,
          creation_time: '2026-02-01T10:00:00Z',
          last_update_time: '2026-02-02T11:00:00Z',
          aliases: [],
        },
      },
      {
        id: 'agent-3',
        display_name: 'test-agent-3',
        entity_id: 'entity-3',
        // No entity - orphaned agent
      },
    ];
    this.model = {
      agents,
      page: 1,
      pageFilter: '',
      pageSize: 5,
    };

    this.renderComponent = () =>
      render(hbs`<Agents::Registry::Table @model={{this.model}} @onClick={{this.onClick}} />`);
  });

  test('it renders table with correct columns', async function (assert) {
    await this.renderComponent();

    const expectedHeaders = [
      'Agent name',
      'Agentic entity in Vault',
      'Entity/Alias ID',
      'Entity status',
      'Entity created at',
      'Entity updated at',
      'Actions',
    ];

    expectedHeaders.forEach((header, index) => {
      assert
        .dom(GENERAL.tableColumnHeader(index + 1, { isAdvanced: true }))
        .includesText(header, `Column ${header} renders`);
    });
  });

  test('it displays agent data correctly', async function (assert) {
    await this.renderComponent();
    const parenRowIds = findAllParentRowIds();

    // Check first agent row
    assert
      .dom(GENERAL.tableData(parenRowIds[0], 'agentName'))
      .includesText('test-agent-1 2 aliases', 'displays agent name');
    assert
      .dom(GENERAL.tableData(parenRowIds[0], 'entityAliasName'))
      .hasText('test-entity-1', 'displays entity name');
    assert
      .dom(GENERAL.tableData(parenRowIds[0], 'entityAliasId'))
      .includesText('entity-1', 'displays entity ID');
    assert
      .dom(GENERAL.tableData(parenRowIds[0], 'entityStatus'))
      .hasText('Enabled', 'displays entity status');

    // Check second agent row (disabled entity)
    assert
      .dom(GENERAL.tableData(parenRowIds[1], 'agentName'))
      .hasText('test-agent-2', 'displays second agent name');
    assert
      .dom(GENERAL.tableData(parenRowIds[1], 'entityStatus'))
      .hasText('Disabled', 'displays disabled status');

    // Check third agent row (no entity)
    assert
      .dom(GENERAL.tableData(parenRowIds[2], 'agentName'))
      .hasText('test-agent-3', 'displays third agent name');
    assert
      .dom(GENERAL.tableData(parenRowIds[2], 'entityAliasName'))
      .hasText('/', 'displays slash for missing entity');
  });

  test('it displays alias count when agent has aliases', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    assert
      .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
      .includesText('2 aliases', 'displays alias count for agent with aliases');
    assert
      .dom(GENERAL.tableData(parentRowIds[1], 'agentName'))
      .doesNotIncludeText('aliases', 'does not display alias count when no aliases');
  });

  test('it formats dates correctly', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    const expectedCreatedAt = dateFormat(['2026-01-01T10:00:00Z', 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });
    const expectedUpdatedAt = dateFormat(['2026-01-02T11:00:00Z', 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });

    assert
      .dom(GENERAL.tableData(parentRowIds[0], 'entityCreatedAt'))
      .hasText(expectedCreatedAt, 'formats creation date correctly');
    assert
      .dom(GENERAL.tableData(parentRowIds[0], 'entityUpdatedAt'))
      .hasText(expectedUpdatedAt, 'formats update date correctly');
  });

  test('it displays copy button for entity/alias IDs', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    assert
      .dom(`${GENERAL.tableData(parentRowIds[0], 'entityAliasId')} ${GENERAL.copyButton}`)
      .exists('copy button exists for entity ID');
  });

  test('it expands to show aliases when clicked', async function (assert) {
    await this.renderComponent();

    // Initially, aliases should not be visible
    assert
      .dom(GENERAL.tableDataNested(1, 'entityAliasId'))
      .isNotVisible('first alias is hidden when nested rows are collapsed');
    assert
      .dom(GENERAL.tableDataNested(2, 'entityAliasId'))
      .isNotVisible('second alias is hidden when nested rows are collapsed');

    // Click to expand nested aliases of first agent
    await click(GENERAL.tableExpandableColumn(0, 'agentName'));

    // Now aliases should be visible
    assert
      .dom(GENERAL.tableDataNested(1, 'entityAliasId'))
      .isVisible('first alias renders when nested rows are expanded');
    assert
      .dom(GENERAL.tableDataNested(2, 'entityAliasId'))
      .isVisible('second alias is renders when nested rows are expanded');

    // Check alias data
    assert
      .dom(GENERAL.tableDataNested(1, 'entityAliasName'))
      .hasText('Alias: test-alias-1', 'shows alias label');
    assert
      .dom(GENERAL.tableDataNested(2, 'entityAliasName'))
      .hasText('Alias: test-alias-2', 'shows second alias name');
  });

  test('it shows popup menu for agents with entities', async function (assert) {
    await this.renderComponent();

    const parentRowIds = findAllParentRowIds();

    // First agent (enabled entity) should have popup menu
    assert
      .dom(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`)
      .exists('popup menu exists for first agent');

    // Third agent (no entity) should have popup menu
    assert
      .dom(`${GENERAL.tableRow(parentRowIds[2])} ${GENERAL.menuTrigger}`)
      .exists('popup menu exists for agent without entity');
  });

  test('it shows correct menu options for enabled entity', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);

    assert.dom(GENERAL.menuItem()).exists({ count: 2 }, 'shows both menu options');
    assert.dom(GENERAL.menuItem('disable')).containsText('Disable entity', 'shows disable option');
    assert.dom(GENERAL.menuItem('delete')).containsText('Delete from registry', 'shows delete option');
  });

  test('it shows correct menu options for disabled entity', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[1])} ${GENERAL.menuTrigger}`);

    assert.dom(GENERAL.menuItem()).exists({ count: 2 }, 'shows both menu options');
    assert.dom(GENERAL.menuItem('enable')).containsText('Enable entity', 'shows enable option');
    assert.dom(GENERAL.menuItem('delete')).containsText('Delete from registry', 'shows delete option');
  });

  test('it opens disable confirmation modal when clicking disable', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('disable'));

    assert
      .dom(GENERAL.confirmTitle)
      .containsText('Disable associated entity for this agent registry record?', 'shows correct modal title');
    assert
      .dom(GENERAL.confirmMessage)
      .containsText(
        'Disabling this entity will prevent the agent from accessing credentials',
        'shows correct modal message'
      );
  });

  test('it calls toggleEntity when confirming disable', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('disable'));
    await click(GENERAL.confirmButton);

    assert.true(this.entityUpdateStub.calledOnce, 'entityUpdateById called once');
    assert.true(
      this.entityUpdateStub.calledWith('entity-1', { disabled: true }),
      'called with correct arguments to disable'
    );
    assert.true(
      this.flashSuccessStub.calledWith('Successfully disabled entity'),
      'shows correct success message'
    );
    assert.true(this.refreshStub.calledWith('vault.cluster.agents.registry.index'), 'refreshes the route');
  });

  test('it calls toggleEntity when confirming enable', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[1])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('enable'));

    // For enabled entity, it should call directly without modal
    assert.true(this.entityUpdateStub.calledOnce, 'entityUpdateById called once');
    assert.true(
      this.entityUpdateStub.calledWith('entity-2', { disabled: false }),
      'called with correct arguments to enable'
    );
    assert.true(
      this.flashSuccessStub.calledWith('Successfully enabled entity'),
      'shows correct success message'
    );
    assert.true(this.refreshStub.calledWith('vault.cluster.agents.registry.index'), 'refreshes the route');
  });

  test('it shows error message when toggle entity fails', async function (assert) {
    this.entityUpdateStub.rejects(new Error('API Error'));

    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[1])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('enable'));

    assert.true(this.flashDangerStub.calledOnce, 'danger flash message shown');
    assert.true(
      this.flashDangerStub.calledWith('Error disabling entity: Test error message'),
      'shows correct error message'
    );
  });

  test('it opens delete confirmation modal when clicking delete', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));

    assert
      .dom(GENERAL.confirmTitle)
      .containsText('Delete this agent registry record?', 'shows correct modal title');
    assert
      .dom(GENERAL.confirmMessage)
      .containsText(
        'Deleting test-agent-1 will permanently remove',
        'shows correct modal message with agent name'
      );
  });

  test('it calls deleteAgent when confirming delete', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    assert.true(this.registrationDeleteStub.calledOnce, 'registrationDeleteByName called once');
    assert.true(
      this.registrationDeleteStub.calledWith('test-agent-1', 'agent-registry'),
      'called with correct arguments'
    );
    assert.true(this.flashSuccessStub.calledOnce, 'success flash message shown');
    assert.true(
      this.flashSuccessStub.calledWith('Successfully deleted agent test-agent-1'),
      'shows correct success message'
    );
    assert.true(this.refreshStub.calledWith('vault.cluster.agents.registry.index'), 'refreshes the route');
  });

  test('it shows error message when delete fails', async function (assert) {
    this.registrationDeleteStub.rejects(new Error('API Error'));

    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    assert.true(this.flashDangerStub.calledOnce, 'danger flash message shown');
    assert.true(
      this.flashDangerStub.calledWith('Error deleting agent test-agent-1: Test error message'),
      'shows correct error message'
    );
  });

  test('it closes modal when clicking cancel on disable confirmation', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('disable'));

    assert.dom(GENERAL.confirmModal).exists('modal is open');

    await click(GENERAL.cancelButton);

    assert.dom(GENERAL.confirmModal).doesNotExist('modal is closed');
    assert.false(this.entityUpdateStub.called, 'entityUpdateById not called');
  });

  test('it closes modal when clicking cancel on delete confirmation', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    await click(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));

    assert.dom(GENERAL.confirmModal).exists('modal is open');

    await click(GENERAL.cancelButton);

    assert.dom(GENERAL.confirmModal).doesNotExist('modal is closed');
    assert.false(this.registrationDeleteStub.called, 'registrationDeleteByName not called');
  });

  test('it handles pagination page change', async function (assert) {
    // add more agents so a second page is available
    for (let i = 4; i < 10; i++) {
      this.model.agents.push({
        name: `agent-${i}`,
        display_name: `test-agent-${i}`,
        entity_id: `entity-${i}`,
      });
    }
    await this.renderComponent();

    await click(GENERAL.nextPage);

    assert.true(this.transitionToStub.calledOnce, 'transitionTo called');
    assert.true(
      this.transitionToStub.calledWith('vault.cluster.agents.registry.index', {
        queryParams: { page: 2 },
      }),
      'transitions to next page'
    );
  });

  test('it handles page size change', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.paginationSizeSelector, '25');

    assert.true(this.transitionToStub.calledOnce, 'transitionTo called');
    assert.true(
      this.transitionToStub.calledWith('vault.cluster.agents.registry.index', {
        queryParams: { page: 1, pageSize: 25 },
      }),
      'transitions with new page size and resets to page 1'
    );
  });

  test('it displays badge with correct color for entity status', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    // Enabled entity should have no warning color
    assert
      .dom(`${GENERAL.tableData(parentRowIds[0], 'entityStatus')} .hds-badge`)
      .exists('badge exists for enabled entity');
    assert
      .dom(`${GENERAL.tableData(parentRowIds[0], 'entityStatus')} .hds-badge`)
      .hasClass('hds-badge--color-neutral', 'enabled entity has neutral color');

    // Disabled entity should have warning color
    assert
      .dom(`${GENERAL.tableData(parentRowIds[1], 'entityStatus')} .hds-badge`)
      .exists('badge exists for disabled entity');
    assert
      .dom(`${GENERAL.tableData(parentRowIds[1], 'entityStatus')} .hds-badge`)
      .hasClass('hds-badge--color-warning', 'disabled entity has warning color');
  });

  test('it does not show popup menu for alias rows', async function (assert) {
    await this.renderComponent();

    // Expand first agent to show aliases
    await click(GENERAL.tableExpandableColumn(0, 'agentName'));

    // Check that alias rows don't have popup menus
    const allMenuTriggers = this.element.querySelectorAll(GENERAL.menuTrigger);
    assert.strictEqual(allMenuTriggers.length, 3, 'only 3 popup menus exist (one per agent, not aliases)');
  });

  test('it handles agents with no entity gracefully', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();

    // Third agent has no entity
    assert
      .dom(GENERAL.tableData(parentRowIds[2], 'agentName'))
      .hasText('test-agent-3', 'displays agent name');
    assert
      .dom(GENERAL.tableData(parentRowIds[2], 'entityAliasName'))
      .hasText('/', 'shows slash for missing entity name');
    assert
      .dom(GENERAL.tableData(parentRowIds[2], 'entityAliasId'))
      .hasText('/', 'shows slash for missing entity ID');
    assert
      .dom(GENERAL.tableData(parentRowIds[2], 'entityStatus'))
      .hasText('/', 'shows slash for missing entity status');
    assert
      .dom(GENERAL.tableData(parentRowIds[2], 'entityCreatedAt'))
      .hasText('/', 'shows slash for missing creation date');
    assert
      .dom(GENERAL.tableData(parentRowIds[2], 'entityUpdatedAt'))
      .hasText('/', 'shows slash for missing update date');
  });

  test('it handles click events for entity and alias', async function (assert) {
    await this.renderComponent();

    const agent = this.model.agents[0];
    await click(GENERAL.button(`agent ${agent.display_name}`));
    assert.strictEqual(this.onClick.lastCall.args[0].agentName, agent.display_name);

    await click(GENERAL.tableExpandableColumn(0, 'agentName'));
    await click(GENERAL.button(`alias ${agent.entity.aliases[0].name}`));
    assert.true(this.onClick.lastCall.args[0].isAlias);
  });

  module('filtering', function () {
    test('it filters agents by agent name', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'test-agent-1');
      const parentRowIds = findAllParentRowIds();

      // Only first agent should be visible
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only one agent');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
        .includesText('test-agent-1', 'shows filtered agent');
    });

    test('it filters agents by entity name', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'test-entity-2');
      const parentRowIds = findAllParentRowIds();

      // Only second agent should be visible
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only one agent');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityAliasName'))
        .hasText('test-entity-2', 'shows filtered entity');
    });

    test('it filters agents by entity ID', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'entity-2');
      const parentRowIds = findAllParentRowIds();

      // Only second agent should be visible
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only one agent');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
        .hasText('test-agent-2', 'shows agent with matching entity ID');
    });

    test('it filters agents by entity status', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'Disabled');
      const parentRowIds = findAllParentRowIds();

      // Only second agent (disabled) should be visible
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only one agent');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityStatus'))
        .hasText('Disabled', 'shows disabled agent');
    });

    test('it filters agents by alias name', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'test-alias-1');
      const parentRowIds = findAllParentRowIds();

      // Only first agent should be visible (has matching alias)
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only one agent');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
        .includesText('test-agent-1', 'shows agent with matching alias');
    });

    test('it filters agents by alias ID', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'alias-2');
      const parentRowIds = findAllParentRowIds();

      // Only first agent should be visible (has matching alias ID)
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only one agent');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
        .includesText('test-agent-1', 'shows agent with matching alias ID');
    });

    test('it filters case-insensitively', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'TEST-AGENT-1');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows one agent with case-insensitive match');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
        .includesText('test-agent-1', 'matches case-insensitively');
    });

    test('it shows multiple matching agents', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'test-agent');

      // All three agents match "test-agent"
      assert.dom(GENERAL.tableParentRow).exists({ count: 3 }, 'shows all matching agents');
    });

    test('it shows no results when filter matches nothing', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'nonexistent-agent');

      assert.dom(GENERAL.tableRow()).doesNotExist('table is hidden when filter matches nothing');
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('No agents found', 'renders "No agents found" empty state message');
    });

    test('it shows all agents when filter is cleared', async function (assert) {
      await this.renderComponent();

      // Apply filter
      await fillIn(GENERAL.filterInput, 'test-agent-1');
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows filtered results');

      // Clear filter
      await fillIn(GENERAL.filterInput, '');
      assert.dom(GENERAL.tableParentRow).exists({ count: 3 }, 'shows all agents when filter is cleared');
    });

    test('it filters by partial matches', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'agent-1');
      const parentRowIds = findAllParentRowIds();

      // Should match test-agent-1
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows partial match');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
        .includesText('test-agent-1', 'matches partial string');
    });

    test('it filters by formatted dates', async function (assert) {
      await this.renderComponent();

      // Filter by part of the formatted date (e.g., "Jan 01" from "2026-01-01T10:00:00Z")
      await fillIn(GENERAL.filterInput, 'Jan 01');
      const parentRowIds = findAllParentRowIds();

      // First agent has creation date of Jan 01, 2026
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'filters by formatted date');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'agentName'))
        .includesText('test-agent-1', 'shows agent with matching date');
    });

    test('it filters before pagination is applied', async function (assert) {
      // Add more agents to test pagination with filtering
      for (let i = 4; i <= 15; i++) {
        this.model.agents.push({
          id: `agent-${i}`,
          display_name: `test-agent-${i}`,
          entity_id: `entity-${i}`,
          entity: {
            id: `entity-${i}`,
            name: `test-entity-${i}`,
            disabled: false,
            creation_time: '2026-01-01T10:00:00Z',
            last_update_time: '2026-01-02T11:00:00Z',
            aliases: [],
          },
        });
      }

      await this.renderComponent();

      // Filter to agents containing "-1" (should match agent-1, agent-10, agent-11, agent-12, agent-13, agent-14, agent-15)
      await fillIn(GENERAL.filterInput, '-1');

      // Should show filtered results across pages
      // With pageSize of 5, first page should show 5 results
      assert.dom(GENERAL.tableParentRow).exists({ count: 5 }, 'shows first page of filtered results');

      // Navigate to next page
      await click(GENERAL.nextPage);

      // Should show remaining filtered results
      assert.dom(GENERAL.tableParentRow).exists({ count: 2 }, 'shows second page of filtered results');
    });

    test('it shows alias rows when parent matches filter', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.filterInput, 'test-agent-1');

      // Expand the filtered agent
      await click(GENERAL.tableExpandableColumn(0, 'agentName'));

      // Aliases should be visible
      assert
        .dom(GENERAL.tableDataNested(1, 'entityAliasName'))
        .isVisible('shows aliases when parent agent matches filter');
      assert
        .dom(GENERAL.tableDataNested(2, 'entityAliasName'))
        .isVisible('shows second alias when parent agent matches filter');
    });

    test('it includes parent when alias matches filter', async function (assert) {
      await this.renderComponent();

      // Filter by alias name that only exists in first agent
      await fillIn(GENERAL.filterInput, 'test-alias-1');

      // Parent agent should be visible
      assert.dom(GENERAL.tableRow(0)).exists({ count: 1 }, 'shows parent agent when alias matches');
      assert
        .dom(GENERAL.tableData(0, 'agentName'))
        .includesText('test-agent-1', 'shows correct parent agent');

      // Expand to verify the alias is there
      await click(GENERAL.tableExpandableColumn(0, 'agentName'));
      assert
        .dom(GENERAL.tableDataNested(1, 'entityAliasName'))
        .hasText('Alias: test-alias-1', 'matching alias is present');
    });
  });
});
