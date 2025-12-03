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
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | ldap | Page::Library::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', () => ({
      data: {
        capabilities: ['root'],
      },
    }));

    this.store = this.owner.lookup('service:store');

    this.store.pushPayload('ldap/library', {
      modelName: 'ldap/library',
      backend: 'ldap-test',
      ...this.server.create('ldap-library', { name: 'test-library' }),
    });
    this.model = this.store.peekRecord('ldap/library', 'test-library');

    this.breadcrumbs = [
      { label: 'ldap-test', route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      { label: 'test-library' },
    ];
  });

  test('it should render page header, tabs and toolbar actions', async function (assert) {
    assert.expect(10);

    this.server.delete(`/${this.model.backend}/library/${this.model.name}`, () => {
      assert.ok(true, 'Request made to delete library');
      return;
    });

    await render(hbs`<Page::Library::Details @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText(this.model.name, 'Library name renders in header');
    assert.dom(GENERAL.breadcrumbAtIdx(0)).containsText(this.model.backend, 'Overview breadcrumb renders');
    assert.dom(GENERAL.breadcrumbAtIdx(1)).containsText('Libraries', 'Libraries breadcrumb renders');
    assert
      .dom('[data-test-breadcrumbs] li:nth-child(3)')
      .containsText(this.model.name, 'Library breadcrumb renders');

    assert.dom(GENERAL.tab('accounts')).hasText('Accounts', 'Accounts tab renders');
    assert.dom('[data-test-tab="config"]').hasText('Configuration', 'Configuration tab renders');

    await click(GENERAL.dropdownToggle('Manage'));
    assert.dom(GENERAL.menuItem('Delete library')).hasText('Delete library', 'Delete action renders');
    assert.dom(GENERAL.menuItem('Edit library')).hasText('Edit library', 'Edit action renders');

    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    await click(GENERAL.menuItem('Delete library'));
    await click(GENERAL.confirmButton);
    assert.ok(
      transitionStub.calledWith('vault.cluster.secrets.backend.ldap.libraries'),
      'Transitions to libraries route on delete success'
    );
  });
});
