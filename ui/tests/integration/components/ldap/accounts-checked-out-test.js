/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | ldap | AccountsCheckedOut', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'ldap-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);
    this.authStub = sinon.stub(this.owner.lookup('service:auth'), 'authData');
    this.pathForStub = sinon.stub(this.owner.lookup('service:capabilities'), 'pathFor').returns('test-path');

    this.library = this.server.create('ldap-library', {
      name: 'test-library',
      completeLibraryName: 'test-library',
    });
    this.statuses = [
      {
        account: 'foo.bar',
        available: false,
        library: 'test-library',
        borrower_client_token: '123',
        borrower_entity_id: '456',
      },
      { account: 'bar.baz', available: false, library: 'test-library' },
      { account: 'checked.in', available: true, library: 'test-library' },
    ];
    this.capabilities = { 'test-path': { canUpdate: true } };
    this.onCheckInSuccess = sinon.spy();
    this.isLoadingStatuses = false;
    this.showLibraryColumn = false;

    this.renderComponent = () =>
      render(
        hbs`<AccountsCheckedOut
          @libraries={{array this.library}}
          @capabilities={{this.capabilities}}
          @statuses={{this.statuses}}
          @showLibraryColumn={{this.showLibraryColumn}}
          @onCheckInSuccess={{this.onCheckInSuccess}}
          @isLoadingStatuses={{this.isLoadingStatuses}}
        />
      `,
        {
          owner: this.engine,
        }
      );
  });

  test('it should render empty state when no accounts are checked out', async function (assert) {
    this.statuses = [
      { account: 'foo', available: true, library: 'test-library' },
      { account: 'bar', available: true, library: 'test-library' },
    ];

    await this.renderComponent();

    assert
      .dom('[data-test-empty-state-title]')
      .hasText('No accounts checked out yet', 'Empty state title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText('There is no account that is currently in use.', 'Empty state message renders');
  });

  test('it should filter accounts for root user', async function (assert) {
    this.authStub.value({ entityId: '' });

    await this.renderComponent();

    assert.dom('[data-test-checked-out-account]').exists({ count: 1 }, 'Correct number of accounts render');
    assert
      .dom('[data-test-checked-out-account="bar.baz"]')
      .hasText('bar.baz', 'Account renders that was checked out by root user');
  });

  test('it should filter accounts for non root user', async function (assert) {
    this.authStub.value({ entityId: '456' });

    await this.renderComponent();

    assert.dom('[data-test-checked-out-account]').exists({ count: 1 }, 'Correct number of accounts render');
    assert
      .dom('[data-test-checked-out-account="foo.bar"]')
      .hasText('foo.bar', 'Account renders that was checked out by non root user');
  });

  test('it should display all accounts when check-in enforcement is disabled on library', async function (assert) {
    this.library.disable_check_in_enforcement = true;

    await this.renderComponent();

    assert.dom('[data-test-checked-out-account]').exists({ count: 2 }, 'Correct number of accounts render');
    assert
      .dom('[data-test-checked-out-account="checked.in"]')
      .doesNotExist('checked.in', 'Checked in accounts do not render');
  });

  test('it should display details in table', async function (assert) {
    this.authStub.value({ entityId: '456' });

    await this.renderComponent();

    assert.dom('[data-test-checked-out-account="foo.bar"]').hasText('foo.bar', 'Account renders');
    assert.dom('[data-test-checked-out-library="foo.bar"]').doesNotExist('Library column is hidden');
    assert
      .dom('[data-test-checked-out-account-action="foo.bar"]')
      .includesText('Check-in', 'Check-in action renders');

    this.showLibraryColumn = true;
    await this.renderComponent();

    assert.dom('[data-test-checked-out-library="foo.bar"]').hasText('test-library', 'Library column renders');
  });

  test('it should check in account', async function (assert) {
    assert.expect(2);

    this.library.disable_check_in_enforcement = true;
    const apiStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'ldapLibraryForceCheckIn')
      .resolves();

    await this.renderComponent();

    await click('[data-test-checked-out-account-action="foo.bar"]');
    await click('[data-test-check-in-confirm]');

    assert.true(
      apiStub.calledWith(this.library.name, this.backend, { service_account_names: ['foo.bar'] }),
      'Check in request called with correct parameters'
    );
    assert.true(this.onCheckInSuccess.called, 'onCheckInSuccess callback was called');
  });

  test('it should show loading state when isLoadingStatuses is true', async function (assert) {
    this.isLoadingStatuses = true;

    await this.renderComponent();

    assert
      .dom('.has-padding-l.flex.is-centered .hds-icon')
      .exists('Loading icon is displayed when isLoadingStatuses is true');
    assert.dom('.hds-table').doesNotExist('Table is not rendered while loading');
  });

  test('it should not show loading state when isLoadingStatuses is false', async function (assert) {
    await this.renderComponent();

    assert
      .dom('.has-padding-l.flex.is-centered .hds-icon')
      .doesNotExist('Loading icon is not displayed when isLoadingStatuses is false');
    assert.dom('.hds-table').exists('Table is rendered when not loading');
  });

  test('it should find library by completeLibraryName for hierarchical libraries', async function (assert) {
    // Create a hierarchical library with proper setup
    this.library = {
      name: 'sa-prod',
      completeLibraryName: 'service-account/sa-prod',
      disable_check_in_enforcement: true, // Allow all accounts to show
    };
    // Mock the auth service to simulate a root user (no entity ID)
    this.authStub.value({ entityId: '' });
    // Status should reference the complete library name
    this.statuses = [
      {
        account: 'prod@example.com',
        available: false,
        library: 'service-account/sa-prod', // Complete hierarchical path
        borrower_client_token: '123',
        borrower_entity_id: '', // Root user has no entity ID
      },
    ];
    this.showLibraryColumn = true;

    await this.renderComponent();

    assert
      .dom('[data-test-checked-out-account="prod@example.com"]')
      .hasText('prod@example.com', 'Account renders for hierarchical library');
    assert
      .dom('[data-test-checked-out-library="prod@example.com"]')
      .hasText('service-account/sa-prod', 'Library name displays full hierarchical path correctly');
  });
});
