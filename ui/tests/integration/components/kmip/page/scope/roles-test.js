/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kmip | Page::Scope::Roles', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const { secrets } = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(secrets, 'kmipDeleteRole').resolves();

    const flash = this.owner.lookup('service:flashMessages');
    this.flashSuccessStub = sinon.stub(flash, 'success');
    this.flashDangerStub = sinon.stub(flash, 'danger');

    const router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(router, 'transitionTo');
    this.refreshStub = sinon.stub(router, 'refresh');

    this.scope = 'scope-1';
    this.roles = ['role-1', 'role-2'];
    this.roles.meta = {
      currentPage: 1,
      pageSize: 10,
      filteredTotal: this.roles.length,
      total: this.roles.length,
    };
    this.filterValue = '';

    const { pathFor } = this.owner.lookup('service:capabilities');
    this.capabilities = this.roles.reduce((capabilities, name) => {
      const path = pathFor('kmipRole', { backend: this.backend, scope: this.scope, name });
      const hasPermission = name === 'role-1'; // only role-1 has permissions to test conditional rendering
      capabilities[path] = { canDelete: hasPermission, canUpdate: hasPermission };
      return capabilities;
    }, {});

    this.renderComponent = () =>
      render(
        hbs`<Page::Scope::Roles @roles={{this.roles}} @scope={{this.scope}} @capabilities={{this.capabilities}} @filterValue={{this.filterValue}} />`,
        { owner: this.engine }
      );
  });

  test('it should render filter and create action in toolbar', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.filterInput).exists('Renders filter input in toolbar');
    assert.dom('[data-test-role-create]').exists('Renders create role action in toolbar');
  });

  test('it should populate filter with arg value', async function (assert) {
    this.filterValue = 'role-1';
    await this.renderComponent();

    assert.dom(GENERAL.filterInput).hasValue('role-1', 'Renders filter input with correct value');
  });

  test('it should filter list items', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.filterInput, 'role-1');
    assert.true(
      this.transitionStub.calledWith({ queryParams: { pageFilter: 'role-1' } }),
      'Transitions with correct query param on page filter change'
    );
  });

  test('it should render list items', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.listItem()).exists({ count: this.roles.length }, 'Renders correct number of roles');
    assert.dom(GENERAL.listItem('role-1')).containsText('role-1', 'Renders role name in list item');

    await click(`${GENERAL.listItem('role-1')} ${GENERAL.menuTrigger}`);
    assert.dom(GENERAL.menuItem('View credentials')).exists('Renders View credentials action in more menu');
    assert.dom(GENERAL.menuItem('View role')).exists('Renders View role action in more menu');
    assert.dom(GENERAL.menuItem('Edit role')).exists('Renders Edit role action in more menu');
    assert
      .dom(`${GENERAL.listItem('role-1')} ${GENERAL.confirmTrigger}`)
      .exists('Renders Delete action in more menu');

    await click(`${GENERAL.listItem('role-2')} ${GENERAL.menuTrigger}`);
    assert
      .dom(GENERAL.menuItem('Edit role'))
      .doesNotExist('Edit role action does not render with no update capability');
    assert
      .dom(`${GENERAL.listItem('role-2')} ${GENERAL.confirmTrigger}`)
      .doesNotExist('Delete action does not render with no delete capability');
  });

  test('it should delete role', async function (assert) {
    await this.renderComponent();

    await click(`${GENERAL.listItem('role-1')} ${GENERAL.menuTrigger}`);
    await click(`${GENERAL.listItem('role-1')} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);

    assert.true(this.apiStub.calledWith('role-1', this.scope, this.backend), 'Calls API to delete role');
    assert.true(
      this.flashSuccessStub.calledWith('Successfully deleted role role-1'),
      'Shows success flash message'
    );
    assert.true(this.refreshStub.called, 'Refreshes the route on delete success');
  });

  test('it should handle delete error', async function (assert) {
    const error = 'An error occurred deleting the role';
    this.apiStub.rejects(getErrorResponse({ errors: [error] }, 500));

    await this.renderComponent();
    await click(`${GENERAL.listItem('role-1')} ${GENERAL.menuTrigger}`);
    await click(`${GENERAL.listItem('role-1')} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);

    assert.true(
      this.flashDangerStub.calledWith(`Error deleting role role-1: ${error}`),
      'Shows flash message on delete error'
    );
  });

  test('it should render pagination controls', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.pagination).exists('Renders pagination controls');
    assert.dom(GENERAL.paginationInfo).hasText('1–2 of 2', 'Renders correct pagination info');
    assert.dom(GENERAL.paginationSizeSelector).doesNotExist('Pagination size selector does not render');
  });
});
