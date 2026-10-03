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

module('Integration | Component | kmip | Page::Scopes', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const { secrets } = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(secrets, 'kmipDeleteScope').resolves();

    const flash = this.owner.lookup('service:flashMessages');
    this.flashSuccessStub = sinon.stub(flash, 'success');
    this.flashDangerStub = sinon.stub(flash, 'danger');

    const router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(router, 'transitionTo');
    this.refreshStub = sinon.stub(router, 'refresh');

    this.scopes = ['scope-1', 'scope-2'];
    this.scopes.meta = {
      currentPage: 1,
      pageSize: 10,
      filteredTotal: this.scopes.length,
      total: this.scopes.length,
    };
    this.filterValue = '';

    const { pathFor } = this.owner.lookup('service:capabilities');
    this.capabilities = this.scopes.reduce((capabilities, name) => {
      const path = pathFor('kmipScope', { backend: this.backend, name });
      capabilities[path] = { canDelete: name === 'scope-1' }; // only scope-1 can be deleted
      return capabilities;
    }, {});

    this.renderComponent = () =>
      render(
        hbs`<Page::Scopes @scopes={{this.scopes}} @capabilities={{this.capabilities}} @filterValue={{this.filterValue}} />`,
        { owner: this.engine }
      );
  });

  test('it should render filter and create action in toolbar', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.filterInput).exists('Renders filter input in toolbar');
    assert.dom('[data-test-scope-create]').exists('Renders create scope action in toolbar');
  });

  test('it should populate filter with arg value', async function (assert) {
    this.filterValue = 'scope-2';
    await this.renderComponent();

    assert.dom(GENERAL.filterInput).hasValue('scope-2', 'Renders filter input with correct value');
  });

  test('it should filter list items', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.filterInput, 'scope-1');
    assert.true(
      this.transitionStub.calledWith({ queryParams: { pageFilter: 'scope-1' } }),
      'Transitions with correct query param on page filter change'
    );
  });

  test('it should render list items', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.listItem()).exists({ count: this.scopes.length }, 'Renders correct number of scopes');
    assert.dom(GENERAL.listItem('scope-1')).containsText('scope-1', 'Renders scope name in list item');

    await click(`${GENERAL.listItem('scope-1')} ${GENERAL.menuTrigger}`);
    assert.dom(GENERAL.menuItem('View scope')).exists('Renders View scope action in more menu');
    assert
      .dom(`${GENERAL.listItem('scope-1')} ${GENERAL.confirmTrigger}`)
      .exists('Renders Delete action in more menu');

    await click(`${GENERAL.listItem('scope-2')} ${GENERAL.menuTrigger}`);
    assert
      .dom(`${GENERAL.listItem('scope-2')} ${GENERAL.confirmTrigger}`)
      .doesNotExist('Does not render Delete action for users without capabilities');
  });

  test('it should delete scope', async function (assert) {
    await this.renderComponent();

    await click(`${GENERAL.listItem('scope-1')} ${GENERAL.menuTrigger}`);
    await click(`${GENERAL.listItem('scope-1')} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);

    assert.true(this.apiStub.calledWith('scope-1', this.backend), 'Calls API to delete scope');
    assert.true(
      this.flashSuccessStub.calledWith('Successfully deleted scope scope-1'),
      'Shows success flash message'
    );
    assert.true(this.refreshStub.called, 'Refreshes the route on delete success');
  });

  test('it should handle delete error', async function (assert) {
    const error = 'An error occurred deleting the scope';
    this.apiStub.rejects(getErrorResponse({ errors: [error] }, 500));

    await this.renderComponent();
    await click(`${GENERAL.listItem('scope-1')} ${GENERAL.menuTrigger}`);
    await click(`${GENERAL.listItem('scope-1')} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);

    assert.true(
      this.flashDangerStub.calledWith(`Error deleting scope scope-1: ${error}`),
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
