/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import sinon from 'sinon';

module('Integration | Component | ldap | AccountsCheckedOut', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    this.store = this.owner.lookup('service:store');
    this.authStub = sinon.stub(this.owner.lookup('service:auth'), 'authData');

    this.store.pushPayload('ldap/library', {
      modelName: 'ldap/library',
      backend: 'ldap-test',
      ...this.server.create('ldap-library', { name: 'test-library' }),
    });
    this.library = this.store.peekRecord('ldap/library', 'test-library');
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
    this.renderComponent = () => {
      return render(
        hbs`
          <div id="modal-wormhole"></div>
          <AccountsCheckedOut @libraries={{array this.library}} @statuses={{this.statuses}} @showLibraryColumn={{this.showLibraryColumn}} />
        `,
        {
          owner: this.engine,
        }
      );
    };
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
    this.authStub.value({});

    await this.renderComponent();

    assert.dom('[data-test-checked-out-account]').exists({ count: 1 }, 'Correct number of accounts render');
    assert
      .dom('[data-test-checked-out-account="bar.baz"]')
      .hasText('bar.baz', 'Account renders that was checked out by root user');
  });

  test('it should filter accounts for non root user', async function (assert) {
    this.authStub.value({ entity_id: '456' });

    await this.renderComponent();

    assert.dom('[data-test-checked-out-account]').exists({ count: 1 }, 'Correct number of accounts render');
    assert
      .dom('[data-test-checked-out-account="foo.bar"]')
      .hasText('foo.bar', 'Account renders that was checked out by non root user');
  });

  test('it should display all accounts when check-in enforcement is disabled on library', async function (assert) {
    this.library.disable_check_in_enforcement = 'Disabled';

    await this.renderComponent();

    assert.dom('[data-test-checked-out-account]').exists({ count: 2 }, 'Correct number of accounts render');
    assert
      .dom('[data-test-checked-out-account="checked.in"]')
      .doesNotExist('checked.in', 'Checked in accounts do not render');
  });

  test('it should display details in table', async function (assert) {
    this.authStub.value({ entity_id: '456' });

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

    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.library.disable_check_in_enforcement = 'Disabled';

    this.server.post('/ldap-test/library/test-library/check-in', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.deepEqual(
        json.service_account_names,
        ['foo.bar'],
        'Check-in request made with correct account names'
      );
    });

    await this.renderComponent();

    await click('[data-test-checked-out-account-action="foo.bar"]');
    await click('[data-test-check-in-confirm]');

    const didTransition = transitionStub.calledWith(
      'vault.cluster.secrets.backend.ldap.libraries.library.details.accounts'
    );
    assert.true(didTransition, 'Transitions to accounts route on check-in success');
  });
});
