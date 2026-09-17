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

module('Integration | Component | identity/entities/table', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.api = this.owner.lookup('service:api');
    this.flashMessages = this.owner.lookup('service:flash-messages');

    this.refreshStub = sinon.stub(this.router, 'refresh');

    this.entityUpdateStub = sinon.stub(this.api.identity, 'entityUpdateById').resolves();
    this.entityDeleteStub = sinon.stub(this.api.identity, 'entityDeleteById').resolves();
    this.entityDeleteAliasStub = sinon.stub(this.api.identity, 'entityDeleteAliasById').resolves();
    this.parseErrorStub = sinon.stub(this.api, 'parseError').resolves({ message: 'Test error message' });

    this.flashSuccessStub = sinon.stub(this.flashMessages, 'success');
    this.flashDangerStub = sinon.stub(this.flashMessages, 'danger');

    this.model = [
      {
        id: 'entity-1',
        name: 'test-entity-1',
        disabled: false,
        creation_time: '2026-01-01T10:00:00Z',
        last_update_time: '2026-01-02T11:00:00Z',
        canEdit: true,
        canAddAlias: true,
        aliases: [
          {
            id: 'alias-1',
            name: 'test-alias-1',
            mount_accessor: 'accessor-1',
            mount_type: 'token',
            creation_time: '2026-01-03T12:00:00Z',
            last_update_time: '2026-01-04T13:00:00Z',
          },
          {
            id: 'alias-2',
            name: 'test-alias-2',
            mount_accessor: 'accessor-2',
            mount_type: 'approle',
            creation_time: '2026-01-05T14:00:00Z',
            last_update_time: '2026-01-06T15:00:00Z',
          },
        ],
      },
      {
        id: 'entity-2',
        name: 'test-entity-2',
        disabled: true,
        creation_time: '2026-02-01T10:00:00Z',
        last_update_time: '2026-02-02T11:00:00Z',
        canEdit: true,
        canAddAlias: true,
        aliases: [],
      },
      {
        id: 'entity-3',
        name: 'test-entity-3',
        disabled: false,
        creation_time: '2026-03-01T10:00:00Z',
        last_update_time: '2026-03-02T11:00:00Z',
        canEdit: false,
        canAddAlias: false,
        aliases: [],
      },
    ];

    this.onClickEntity = sinon.stub();
    this.onClickAlias = sinon.stub();

    this.renderComponent = () =>
      render(
        hbs`<Identity::Entities::Table @model={{this.model}} @onClickEntity={{this.onClickEntity}} @onClickAlias={{this.onClickAlias}} />`
      );
  });

  test('it renders table with correct columns', async function (assert) {
    await this.renderComponent();
    const expectedHeaders = [
      'Entity name',
      'Entity/Alias ID',
      'Entity status',
      'Auth method',
      'Last updated',
      'Created at',
      'Actions',
    ];

    expectedHeaders.forEach((header, index) => {
      assert
        .dom(GENERAL.tableColumnHeader(index + 1, { isAdvanced: true }))
        .includesText(header, `Column "${header}" renders`);
    });
  });

  test('it displays entity data correctly', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();
    assert
      .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
      .includesText('test-entity-1', 'displays first entity name');
    assert
      .dom(GENERAL.tableData(parentRowIds[0], 'entityAliasId'))
      .includesText('entity-1', 'displays entity ID');
    assert
      .dom(GENERAL.tableData(parentRowIds[0], 'entityStatus'))
      .hasText('Enabled', 'displays enabled status');

    assert
      .dom(GENERAL.tableData(parentRowIds[1], 'entityName'))
      .includesText('test-entity-2', 'displays second entity name');
    assert
      .dom(GENERAL.tableData(parentRowIds[1], 'entityStatus'))
      .hasText('Disabled', 'displays disabled status');
  });

  test('it displays alias count when entity has aliases', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();
    assert
      .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
      .includesText('2 aliases', 'displays alias count for entity with aliases');
    assert
      .dom(GENERAL.tableData(parentRowIds[1], 'entityName'))
      .doesNotIncludeText('alias', 'does not display alias count when no aliases');
  });

  test('it displays singular "alias" label when entity has exactly one alias', async function (assert) {
    this.model[0].aliases = [this.model[0].aliases[0]];
    await this.renderComponent();
    assert
      .dom(GENERAL.tableData(0, 'entityName'))
      .includesText('1 alias', 'displays singular "alias" for one alias');
  });

  test('it formats dates correctly', async function (assert) {
    await this.renderComponent();
    const expectedCreatedAt = dateFormat(['2026-01-01T10:00:00Z', 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });
    const expectedUpdatedAt = dateFormat(['2026-01-02T11:00:00Z', 'MMM dd, yyyy hh:mm:ss a'], {
      withTimeZone: true,
    });

    assert
      .dom(GENERAL.tableData(0, 'entityCreatedAt'))
      .hasText(expectedCreatedAt, 'formats creation date correctly');
    assert
      .dom(GENERAL.tableData(0, 'entityUpdatedAt'))
      .hasText(expectedUpdatedAt, 'formats update date correctly');
  });

  test('it displays copy button for entity/alias IDs', async function (assert) {
    await this.renderComponent();
    assert
      .dom(`${GENERAL.tableData(0, 'entityAliasId')} ${GENERAL.copyButton}`)
      .exists('copy button exists for entity ID');
  });

  test('it displays entity status badge with correct color', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();
    assert
      .dom(`${GENERAL.tableData(parentRowIds[0], 'entityStatus')} .hds-badge`)
      .hasClass('hds-badge--color-neutral', 'enabled entity has neutral badge color');

    assert
      .dom(`${GENERAL.tableData(parentRowIds[1], 'entityStatus')} .hds-badge`)
      .hasClass('hds-badge--color-warning', 'disabled entity has warning badge color');
  });

  test('it expands to show aliases when clicked', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.tableDataNested(1, 'entityAliasId')).isNotVisible('first alias hidden when collapsed');
    assert
      .dom(GENERAL.tableDataNested(2, 'entityAliasId'))
      .isNotVisible('second alias hidden when collapsed');

    await click(GENERAL.tableExpandableColumn(0, 'entityName'));

    assert.dom(GENERAL.tableDataNested(1, 'entityAliasId')).isVisible('first alias visible after expand');
    assert.dom(GENERAL.tableDataNested(2, 'entityAliasId')).isVisible('second alias visible after expand');

    assert
      .dom(GENERAL.tableDataNested(1, 'entityName'))
      .hasText('Alias: test-alias-1', 'shows "Alias:" label and alias name');
    assert
      .dom(GENERAL.tableDataNested(2, 'entityName'))
      .hasText('Alias: test-alias-2', 'shows second alias name');
  });

  test('it shows popup menu for all entity rows', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();
    assert
      .dom(`${GENERAL.tableRow(parentRowIds[0])} ${GENERAL.menuTrigger}`)
      .exists('popup menu exists for first entity');
    assert
      .dom(`${GENERAL.tableRow(parentRowIds[1])} ${GENERAL.menuTrigger}`)
      .exists('popup menu exists for second entity');
  });

  test('it shows edit and create-alias options when entity canEdit and canAddAlias', async function (assert) {
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);

    assert.dom(GENERAL.menuItem('edit')).containsText('Edit entity', 'shows edit option');
    assert.dom(GENERAL.menuItem('create alias')).containsText('Create alias', 'shows create alias option');
  });

  test('it hides edit and create-alias options when entity cannot edit or add alias', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();
    await click(`${GENERAL.tableRow(parentRowIds[2])} ${GENERAL.menuTrigger}`);

    assert.dom(GENERAL.menuItem('edit')).doesNotExist('hides edit option when canEdit is false');
    assert
      .dom(GENERAL.menuItem('create alias'))
      .doesNotExist('hides create alias option when canAddAlias is false');
  });

  test('it shows disable option for an enabled entity', async function (assert) {
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);

    assert.dom(GENERAL.menuItem('disable')).containsText('Disable entity', 'shows disable option');
    assert.dom(GENERAL.menuItem('enable')).doesNotExist('does not show enable option for enabled entity');
  });

  test('it shows enable option for a disabled entity', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();
    await click(`${GENERAL.tableRow(parentRowIds[1])} ${GENERAL.menuTrigger}`);

    assert.dom(GENERAL.menuItem('enable')).containsText('Enable entity', 'shows enable option');
    assert.dom(GENERAL.menuItem('disable')).doesNotExist('does not show disable option for disabled entity');
  });

  test('it opens disable confirmation modal when clicking disable', async function (assert) {
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('disable'));

    assert.dom(GENERAL.confirmTitle).containsText('Disable this entity?', 'shows correct modal title');
    assert
      .dom(GENERAL.confirmMessage)
      .containsText(
        'Associated tokens will not be revoked, but cannot be used.',
        'shows correct modal message'
      );
  });

  test('it calls toggleEntity when confirming disable', async function (assert) {
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('disable'));
    await click(GENERAL.confirmButton);

    assert.true(this.entityUpdateStub.calledOnce, 'entityUpdateById called once');
    assert.true(
      this.entityUpdateStub.calledWith('entity-1', { disabled: true }),
      'called with correct arguments to disable'
    );
    assert.true(
      this.flashSuccessStub.calledWith('Successfully disabled entity'),
      'shows correct success flash'
    );
    assert.true(
      this.refreshStub.calledWith('vault.cluster.access.identity.entities.index'),
      'refreshes entities index route'
    );
  });

  test('it calls toggleEntity directly when confirming enable (no modal)', async function (assert) {
    await this.renderComponent();
    const parentRowIds = findAllParentRowIds();
    await click(`${GENERAL.tableRow(parentRowIds[1])} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('enable'));

    assert.true(this.entityUpdateStub.calledOnce, 'entityUpdateById called once');
    assert.true(
      this.entityUpdateStub.calledWith('entity-2', { disabled: false }),
      'called with correct arguments to enable'
    );
    assert.true(
      this.flashSuccessStub.calledWith('Successfully enabled entity'),
      'shows correct success flash'
    );
    assert.true(
      this.refreshStub.calledWith('vault.cluster.access.identity.entities.index'),
      'refreshes entities index route'
    );
  });

  test('it shows error message when toggleEntity fails', async function (assert) {
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

  test('it opens delete entity confirmation modal', async function (assert) {
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));

    assert.dom(GENERAL.confirmTitle).containsText('Delete this entity?', 'shows correct modal title');
  });

  test('it calls deleteEntity when confirming delete', async function (assert) {
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    assert.true(this.entityDeleteStub.calledOnce, 'entityDeleteById called once');
    assert.true(this.entityDeleteStub.calledWith('entity-1'), 'called with correct entity ID');
    assert.true(
      this.flashSuccessStub.calledWith('Successfully deleted entity entity-1'),
      'shows correct success flash'
    );
    assert.true(
      this.refreshStub.calledWith('vault.cluster.access.identity.entities.index'),
      'refreshes entities index route'
    );
  });

  test('it shows error message when deleteEntity fails', async function (assert) {
    this.entityDeleteStub.rejects(new Error('API Error'));
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    assert.true(this.flashDangerStub.calledOnce, 'danger flash message shown');
    assert.true(
      this.flashDangerStub.calledWith('Error deleting entity entity-1: Test error message'),
      'shows correct error message'
    );
  });

  test('it closes disable modal on cancel without calling API', async function (assert) {
    await this.renderComponent();
    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('disable'));
    assert.dom(GENERAL.confirmModal).exists('disable modal is open');

    await click(GENERAL.cancelButton);

    assert.dom(GENERAL.confirmModal).doesNotExist('modal is closed after cancel');
    assert.false(this.entityUpdateStub.called, 'entityUpdateById not called');
  });

  test('it closes delete entity modal on cancel without calling API', async function (assert) {
    await this.renderComponent();

    await click(`${GENERAL.tableRow(0)} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));
    assert.dom(GENERAL.confirmModal).exists('delete modal is open');

    await click(GENERAL.cancelButton);

    assert.dom(GENERAL.confirmModal).doesNotExist('modal is closed after cancel');
    assert.false(this.entityDeleteStub.called, 'entityDeleteById not called');
  });

  test('alias rows show edit-alias and delete-alias menu options', async function (assert) {
    // give the alias entity canEdit permission
    this.model[0].aliases[0].canEdit = true;
    await this.renderComponent();

    await click(GENERAL.tableExpandableColumn(0, 'entityName'));
    await click(`${GENERAL.tableDataNested(1, 'popupMenu')} ${GENERAL.menuTrigger}`);

    assert.dom(GENERAL.menuItem('edit alias')).containsText('Edit alias', 'shows edit alias option');
    assert.dom(GENERAL.menuItem('delete')).containsText('Delete alias', 'shows delete alias option');
  });

  test('it opens delete alias confirmation modal', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.tableExpandableColumn(0, 'entityName'));
    await click(`${GENERAL.tableDataNested(1, 'popupMenu')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));

    assert.dom(GENERAL.confirmTitle).containsText('Delete this alias?', 'shows correct alias modal title');
  });

  test('it calls deleteAlias when confirming alias delete', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.tableExpandableColumn(0, 'entityName'));
    await click(`${GENERAL.tableDataNested(1, 'popupMenu')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    assert.true(this.entityDeleteAliasStub.calledOnce, 'entityDeleteAliasById called once');
    assert.true(this.entityDeleteAliasStub.calledWith('alias-1'), 'called with correct alias ID');
    assert.true(
      this.flashSuccessStub.calledWith('Successfully deleted alias alias-1'),
      'shows correct success flash'
    );
    assert.true(
      this.refreshStub.calledWith('vault.cluster.access.identity.entities.index'),
      'refreshes entities index route'
    );
  });

  test('it shows error message when deleteAlias fails', async function (assert) {
    this.entityDeleteAliasStub.rejects(new Error('API Error'));
    await this.renderComponent();

    await click(GENERAL.tableExpandableColumn(0, 'entityName'));
    await click(`${GENERAL.tableDataNested(1, 'popupMenu')} ${GENERAL.menuTrigger}`);
    await click(GENERAL.menuItem('delete'));
    await click(GENERAL.confirmButton);

    assert.true(this.flashDangerStub.calledOnce, 'danger flash message shown');
    assert.true(
      this.flashDangerStub.calledWith('Error deleting alias alias-1: Test error message'),
      'shows correct error message'
    );
  });

  test('it renders the empty state when model is empty', async function (assert) {
    this.model = [];
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No entities found', 'renders empty state title');
  });

  module('filtering', function () {
    test('it filters entities by name', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'test-entity-1');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only matching entity');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
        .includesText('test-entity-1', 'shows filtered entity');
    });

    test('it filters entities by ID', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'entity-2');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only matching entity');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
        .includesText('test-entity-2', 'shows entity by ID');
    });

    test('it filters entities by status', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'Disabled');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows only disabled entity');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityStatus'))
        .hasText('Disabled', 'shows disabled entity');
    });

    test('it filters entities by alias name', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'test-alias-1');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows parent entity when alias matches');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
        .includesText('test-entity-1', 'correct parent entity shown');
    });

    test('it filters entities by alias ID', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'alias-2');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows parent entity when alias ID matches');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
        .includesText('test-entity-1', 'correct parent entity shown');
    });

    test('it filters case-insensitively', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'TEST-ENTITY-1');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'matches case-insensitively');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
        .includesText('test-entity-1', 'shows case-insensitively matched entity');
    });

    test('it shows multiple matching entities', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'test-entity');

      assert.dom(GENERAL.tableParentRow).exists({ count: 3 }, 'shows all matching entities');
    });

    test('it shows empty state when filter matches nothing', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'nonexistent-entity');

      assert.dom(GENERAL.tableParentRow).doesNotExist('table is hidden when nothing matches');
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('No entities found', 'renders "No entities found" empty state');
    });

    test('it shows all entities when filter is cleared', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'test-entity-1');
      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows filtered results');

      await fillIn(GENERAL.filterInput, '');
      assert.dom(GENERAL.tableParentRow).exists({ count: 3 }, 'shows all entities when filter cleared');
    });

    test('it filters by partial name match', async function (assert) {
      await this.renderComponent();
      await fillIn(GENERAL.filterInput, 'entity-1');
      const parentRowIds = findAllParentRowIds();

      assert.dom(GENERAL.tableParentRow).exists({ count: 1 }, 'shows partial match');
      assert
        .dom(GENERAL.tableData(parentRowIds[0], 'entityName'))
        .includesText('test-entity-1', 'matches partial string');
    });
  });
});
