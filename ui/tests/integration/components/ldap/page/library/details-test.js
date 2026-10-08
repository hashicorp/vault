/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupDetailsTest } from 'vault/tests/helpers/ldap/ldap-helpers';

module('Integration | Component | ldap | Page::Library::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);
  setupDetailsTest(hooks);

  hooks.beforeEach(function () {
    this.breadcrumbs = [
      { label: 'ldap-test', route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      { label: 'test-library' },
    ];
  });

  test('it should render page header, tabs and toolbar actions', async function (assert) {
    assert.expect(10);

    const apiStub = sinon.stub(this.owner.lookup('service:api').secrets, 'ldapLibraryDelete').resolves();

    await render(hbs`<Page::Library::Details @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(this.name, 'Library name renders in header');
    assert.dom(GENERAL.breadcrumbAtIdx(0)).containsText(this.backend, 'Overview breadcrumb renders');
    assert.dom(GENERAL.breadcrumbAtIdx(1)).containsText('Libraries', 'Libraries breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(3)')
      .containsText(this.name, 'Library breadcrumb renders');

    assert.dom(GENERAL.tab('accounts')).hasText('Accounts', 'Accounts tab renders');
    assert.dom('[data-test-tab="config"]').hasText('Configuration', 'Configuration tab renders');

    await click(GENERAL.dropdownToggle('Manage'));
    assert.dom(GENERAL.menuItem('Delete library')).hasText('Delete library', 'Delete action renders');
    assert.dom(GENERAL.menuItem('Edit library')).hasText('Edit library', 'Edit action renders');

    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    await click(GENERAL.menuItem('Delete library'));
    await click(GENERAL.confirmButton);

    assert.true(apiStub.calledWith(this.name, this.backend), 'Delete API called with correct parameters');
    assert.ok(
      transitionStub.calledWith('vault.cluster.secrets.backend.ldap.libraries'),
      'Transitions to libraries route on delete success'
    );
  });

  // An empty dropdown would open a blank popover, so the toggle is omitted entirely.
  test('it should hide the Manage dropdown when no actions are permitted', async function (assert) {
    const { pathFor } = this.owner.lookup('service:capabilities');
    const libraryPath = pathFor('ldapLibrary', { backend: this.backend, name: this.name });
    this.model.capabilities[libraryPath] = { canRead: true, canUpdate: false, canDelete: false };

    await render(hbs`<Page::Library::Details @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    assert.dom(GENERAL.dropdownToggle('Manage')).doesNotExist('Manage is hidden');
  });

  test('it should show the Manage dropdown with only the permitted action', async function (assert) {
    const { pathFor } = this.owner.lookup('service:capabilities');
    const libraryPath = pathFor('ldapLibrary', { backend: this.backend, name: this.name });
    this.model.capabilities[libraryPath] = { canRead: true, canUpdate: false, canDelete: true };

    await render(hbs`<Page::Library::Details @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    await click(GENERAL.dropdownToggle('Manage'));

    assert.dom(GENERAL.menuItem()).exists({ count: 1 }, 'only one action renders');
    assert.dom(GENERAL.menuItem('Delete library')).exists('Delete library renders');
  });
});
