/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { paginate } from 'core/utils/paginate-list';

module('Integration | Component | generated-item-list', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = paginate([{ id: 'bob' }], { page: 1 });
    this.paths = [];
    this.methodModel = { id: 'userpass', type: 'userpass' };
    this.urls = { delete: (id) => `/auth/userpass/users/${id}` };
    const api = this.owner.lookup('service:api');
    this.delete = sinon.stub(api.request, 'delete').rejects(new Error('denied'));
    this.parseError = sinon.stub(api, 'parseError').resolves({ message: 'permission denied' });
    const flashMessages = this.owner.lookup('service:flash-messages');
    this.danger = sinon.stub(flashMessages, 'danger');
    const controlGroup = this.owner.lookup('service:control-group');
    this.saveTokenFromError = sinon.stub(controlGroup, 'saveTokenFromError');
    this.logFromError = sinon.stub(controlGroup, 'logFromError').returns({
      accessor: 'cg-accessor',
      token: 'cg-token',
      creation_path: 'auth/userpass/users/existing',
      href: '/ui/vault/access/control-groups/cg-accessor',
    });
    this.deleteFromMenu = async () => {
      await render(
        hbs`<GeneratedItemList @model={{this.model}} @itemType="users" @paths={{this.paths}} @methodModel={{this.methodModel}} @urls={{this.urls}} />`
      );
      await click('[data-test-popup-menu-trigger]');
      await click(GENERAL.menuItem('delete'));
    };
  });

  test('it opens one delete confirmation for the selected item', async function (assert) {
    await this.deleteFromMenu();
    assert.dom(GENERAL.confirmButton).exists({ count: 1 }, 'a single confirmation dialog renders');
  });

  test('a failed delete shows an error flash', async function (assert) {
    await this.deleteFromMenu();
    await click(GENERAL.confirmButton);

    assert.true(this.delete.calledOnceWith('/auth/userpass/users/bob'), 'DELETE targets the item');
    assert.true(
      this.danger.calledWith('There was an error deleting this user permission denied'),
      'error flash is shown'
    );
  });

  test('a delete behind a Control Group shows the approval details', async function (assert) {
    const response = {
      isControlGroupError: true,
      accessor: 'cg-accessor',
      token: 'cg-token',
      creation_path: 'auth/userpass/users/bob',
    };
    this.parseError.resolves({ message: 'Control Group encountered', response });
    await this.deleteFromMenu();
    await click(GENERAL.confirmButton);

    assert.true(this.saveTokenFromError.calledWith(response), 'token is saved for a retry');
    assert.dom(GENERAL.messageError).includesText('cg-accessor', 'control group error is shown');
    assert.true(this.danger.notCalled, 'no error flash is shown');
  });
});
