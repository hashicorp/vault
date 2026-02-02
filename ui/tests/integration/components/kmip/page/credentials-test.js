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

module('Integration | Component | kmip | Page::Credentials', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const { secrets } = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(secrets, 'kmipRevokeClientCertificate').resolves();

    const flash = this.owner.lookup('service:flashMessages');
    this.flashSuccessStub = sinon.stub(flash, 'success');
    this.flashDangerStub = sinon.stub(flash, 'danger');

    const router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(router, 'transitionTo');
    this.refreshStub = sinon.stub(router, 'refresh');

    this.credentials = ['12345', '54321'];
    this.credentials.meta = {
      currentPage: 1,
      pageSize: 10,
      filteredTotal: this.credentials.length,
      total: this.credentials.length,
    };
    this.filterValue = '';
    this.capabilities = { canDelete: true };
    this.roleName = 'role-1';
    this.scopeName = 'scope-1';

    this.renderComponent = () =>
      render(
        hbs`<Page::Credentials @credentials={{this.credentials}} @capabilities={{this.capabilities}} @roleName={{this.roleName}} @scopeName={{this.scopeName}} @filterValue={{this.filterValue}} />`,
        { owner: this.engine }
      );
  });

  test('it should render filter and generate action in toolbar', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.filterInput).exists('Renders filter input in toolbar');
    assert.dom('[data-test-generate-credentials]').exists('Renders generate action in toolbar');
  });

  test('it should populate filter with arg value', async function (assert) {
    this.filterValue = '222333444';
    await this.renderComponent();

    assert.dom(GENERAL.filterInput).hasValue('222333444', 'Renders filter input with correct value');
  });

  test('it should filter list items', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.filterInput, '23');
    assert.true(
      this.transitionStub.calledWith({ queryParams: { pageFilter: '23' } }),
      'Transitions with correct query param on page filter change'
    );
  });

  test('it should render list items', async function (assert) {
    await this.renderComponent();

    assert
      .dom(GENERAL.listItem())
      .exists({ count: this.credentials.length }, 'Renders correct number of credentials');
    assert.dom(GENERAL.listItem('12345')).containsText('12345', 'Renders serial_number in list item');

    await click(`${GENERAL.listItem('12345')} ${GENERAL.menuTrigger}`);
    assert.dom(GENERAL.menuItem('View credentials')).exists('Renders View credentials action in more menu');
    assert
      .dom(`${GENERAL.listItem('12345')} ${GENERAL.confirmTrigger}`)
      .exists('Renders Revoke action in more menu');

    this.capabilities.canDelete = false;
    await click(`${GENERAL.listItem('12345')} ${GENERAL.menuTrigger}`);
    assert
      .dom(`${GENERAL.listItem('12345')} ${GENERAL.confirmTrigger}`)
      .doesNotExist('Does not render Delete action for users without capabilities');
  });

  test('it should revoke credentials', async function (assert) {
    await this.renderComponent();

    await click(`${GENERAL.listItem('12345')} ${GENERAL.menuTrigger}`);
    await click(`${GENERAL.listItem('12345')} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);

    assert.true(
      this.apiStub.calledWith(this.roleName, this.scopeName, this.backend, { serial_number: '12345' }),
      'Calls API to revoke credentials'
    );
    assert.true(
      this.flashSuccessStub.calledWith('Successfully revoked credentials 12345'),
      'Shows success flash message'
    );
    assert.true(this.refreshStub.called, 'Refreshes the route on revoke success');
  });

  test('it should handle delete error', async function (assert) {
    const error = 'An error occurred revoking credentials';
    this.apiStub.rejects(getErrorResponse({ errors: [error] }, 500));

    await this.renderComponent();
    await click(`${GENERAL.listItem('12345')} ${GENERAL.menuTrigger}`);
    await click(`${GENERAL.listItem('12345')} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);

    assert.true(
      this.flashDangerStub.calledWith(`Error revoking credentials 12345: ${error}`),
      'Shows flash message on revoke error'
    );
  });

  test('it should render pagination controls', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.pagination).exists('Renders pagination controls');
    assert.dom(GENERAL.paginationInfo).hasText('1–2 of 2', 'Renders correct pagination info');
    assert.dom(GENERAL.paginationSizeSelector).doesNotExist('Pagination size selector does not render');
  });
});
