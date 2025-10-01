/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { render, click, fillIn, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | kv-v2 | Page::List', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    const paginationMeta = {
      currentPage: 1,
      lastPage: 2,
      nextPage: 2,
      prevPage: 1,
      total: 5,
      filteredTotal: 5,
      pageSize: 3,
    };
    this.secrets = ['secret-1', 'my-path/', 'secret-2'];
    this.secrets.meta = paginationMeta;
    this.pathToSecret = 'my-kv/';
    this.backend = 'kv-engine';
    this.filterValue = '';
    this.failedDirectoryQuery = false;
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
    ];
    this.capabilities = { canRead: true, canDelete: true };

    this.renderComponent = () =>
      render(
        hbs`
        <Page::List
          @secrets={{this.secrets}}
          @backend={{this.backend}}
          @pathToSecret={{this.pathToSecret}}
          @filterValue={{this.filterValue}}
          @failedDirectoryQuery={{this.failedDirectoryQuery}}
          @breadcrumbs={{this.breadcrumbs}}
          @currentRouteParams={{array this.backend}}
          @capabilities={{this.capabilities}}
        />`,
        { owner: this.engine }
      );

    this.transitionTo = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    setRunOptions({
      rules: {
        // TODO: ConfirmAction renders modal within list when @isInDropdown
        list: { enabled: false },
      },
    });
  });

  test('it should render page title and toolbar elements', async function (assert) {
    await this.renderComponent();

    assert.dom(PAGE.title).includesText(this.backend, 'renders mount path as page title');
    assert.dom(PAGE.secretTab('Secrets')).exists('renders Secrets tab');
    assert.dom(PAGE.secretTab('Configuration')).exists('renders Configuration tab');
    assert.dom(PAGE.list.filter).exists('renders filter input');
    assert.dom(PAGE.list.createSecret).exists('renders create secret action');
  });

  test('it should render 403 state', async function (assert) {
    this.secrets = 403;
    this.failedDirectoryQuery = true;
    await this.renderComponent();

    assert.dom(PAGE.list.filter).doesNotExist('filter input is hidden');
    assert.dom(PAGE.list.overviewCard).exists('renders overview card');
    assert.dom(PAGE.list.overviewInput).hasValue('my-kv/', 'shows correct path in overview card input');

    await typeIn(PAGE.list.overviewInput, 'my-dir/');
    await click(GENERAL.submitButton);
    assert.true(
      this.transitionTo.calledWith('vault.cluster.secrets.backend.kv.list-directory', 'my-kv/my-dir/'),
      'transitions to correct route if path is directory'
    );
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'You do not have the required permissions or the directory does not exist.',
        'alert renders for failed directory query'
      );

    await fillIn(PAGE.list.overviewInput, '');
    await typeIn(PAGE.list.overviewInput, 'secret');
    await click(GENERAL.submitButton);
    assert.true(
      this.transitionTo.calledWith('vault.cluster.secrets.backend.kv.secret.index', 'secret'),
      'transitions to correct route if path is not a directory'
    );
  });

  test('it should render empty states', async function (assert) {
    this.secrets = [];
    await this.renderComponent();
    assert.dom(GENERAL.emptyStateTitle).hasText('No secrets yet', 'empty state renders for no secrets');

    this.filterValue = 'foo';
    await this.renderComponent();
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('There are no secrets matching "foo".', 'empty state renders for no filter results');
  });

  test('it should render paginated secrets', async function (assert) {
    await this.renderComponent();

    assert.dom(PAGE.list.item()).exists({ count: 3 }, 'renders 3 secrets for first page');
    assert.dom(PAGE.list.item('secret-1')).hasText('secret-1', 'secret path renders');
    assert.dom(GENERAL.pagination).exists('renders hds pagination component');
    assert.dom(GENERAL.paginationInfo).hasText('1–3 of 5', 'renders correct page information');
  });

  test('it should render list item menu', async function (assert) {
    await this.renderComponent();

    await click(`${PAGE.list.item('my-path/')} ${PAGE.popup}`);
    assert.dom(PAGE.list.menuItem('Content')).exists('renders content menu item for directory');

    await click(`${PAGE.list.item('secret-1')} ${PAGE.popup}`);
    assert.dom(PAGE.list.menuItem('Overview')).exists('renders overview menu item');
    assert.dom(PAGE.list.menuItem('Secret data')).exists('renders secret data menu item');
    assert.dom(PAGE.list.menuItem('View version history')).exists('renders version history menu item');
    assert.dom(PAGE.list.menuItem('Permanently delete')).exists('renders delete menu item');

    await click(PAGE.list.menuItem('Permanently delete'));
    assert
      .dom(GENERAL.confirmMessage)
      .hasText(
        'This will permanently delete this secret and all its versions.',
        'renders confirm modal on delete click'
      );

    this.deleteStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'kvV2DeleteMetadataAndAllVersions')
      .resolves();
    await click(GENERAL.confirmButton);
    assert.true(
      this.deleteStub.calledWith('my-kv/secret-1', this.backend),
      'makes request to delete secret on confirm'
    );
  });
});
