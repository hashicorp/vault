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
import GeneratedItemForm from 'vault/forms/generated-item';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

module('Integration | Component | generated-item', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const username = new FormField('username', 'string', { readOnly: true });
    const numUses = new FormField('token_num_uses', 'number');
    this.metadata = {
      formFields: [username, numUses],
      formFieldGroups: [new FormFieldGroup('default', [username, numUses])],
      identifier: 'username',
    };
    this.urls = {
      list: '/auth/userpass/users/',
      read: (id) => `/auth/userpass/users/${id}`,
      create: (id) => `/auth/userpass/users/${id}`,
      update: (id) => `/auth/userpass/users/${id}`,
      delete: (id) => `/auth/userpass/users/${id}`,
    };
    const api = this.owner.lookup('service:api');
    this.post = sinon.stub(api.request, 'post').resolves();
    this.put = sinon.stub(api.request, 'put').resolves();
    this.delete = sinon.stub(api.request, 'delete').resolves();
    this.parseError = sinon.stub(api, 'parseError').resolves({ message: 'permission denied' });
    const router = this.owner.lookup('service:router');
    this.transitionTo = sinon.stub(router, 'transitionTo').returns({ followRedirects: () => {} });
    const flashMessages = this.owner.lookup('service:flash-messages');
    this.success = sinon.stub(flashMessages, 'success');
    this.danger = sinon.stub(flashMessages, 'danger');
    const controlGroup = this.owner.lookup('service:control-group');
    this.saveTokenFromError = sinon.stub(controlGroup, 'saveTokenFromError');
    this.logFromError = sinon.stub(controlGroup, 'logFromError').returns({
      accessor: 'cg-accessor',
      token: 'cg-token',
      creation_path: 'auth/userpass/users/existing',
      href: '/ui/vault/access/control-groups/cg-accessor',
    });
    this.controlGroupResponse = {
      isControlGroupError: true,
      accessor: 'cg-accessor',
      token: 'cg-token',
      creation_path: 'auth/userpass/users/existing',
      creation_time: '2026-11-05T00:00:00Z',
      ttl: 86400,
    };
    this.renderItem = (mode, form) => {
      this.mode = mode;
      this.form = form;
      return render(
        hbs`<GeneratedItem @form={{this.form}} @urls={{this.urls}} @mountPath="userpass" @mode={{this.mode}} @itemType="user" />`
      );
    };
  });

  test('create posts the form payload to the identifier path', async function (assert) {
    const form = new GeneratedItemForm(this.metadata, {}, { isNew: true });
    form.set('username', 'new-user');
    form.set('token_num_uses', '5');
    await this.renderItem('create', form);
    await click(GENERAL.submitButton);

    assert.true(
      this.post.calledOnceWith('/auth/userpass/users/new-user', { token_num_uses: 5 }),
      'POST sends the serialized form'
    );
    assert.true(this.success.calledWith('Successfully saved user new-user.'), 'success flash is shown');
    assert.true(
      this.transitionTo.calledWith('vault.cluster.access.method.item.list'),
      'returns to the item list'
    );
  });

  test('edit puts to the identifier path and renders request errors', async function (assert) {
    this.put.rejects(new Error('denied'));
    const form = new GeneratedItemForm(this.metadata, { username: 'existing', token_num_uses: 1 });
    await this.renderItem('edit', form);
    await click(GENERAL.submitButton);

    assert.true(this.put.calledOnceWith('/auth/userpass/users/existing'), 'PUT targets the existing item');
    assert.dom(GENERAL.messageError).includesText('permission denied', 'the parsed error is shown');
    assert.true(this.transitionTo.notCalled, 'stays on the form after a failed save');
  });

  test('show deletes the item and reports a failed delete', async function (assert) {
    const form = new GeneratedItemForm(this.metadata, { username: 'existing', token_num_uses: 1 });
    await this.renderItem('show', form);

    this.delete.rejects(new Error('denied'));
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.true(
      this.danger.calledWith('There was an error deleting this user permission denied'),
      'a failed delete shows a danger flash'
    );
    assert.true(this.transitionTo.notCalled, 'stays on the item after a failed delete');

    this.delete.resolves();
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.true(this.delete.calledWith('/auth/userpass/users/existing'), 'DELETE targets the item');
    assert.true(this.success.calledWith('Successfully deleted user existing.'), 'success flash is shown');
  });

  test('a save behind a Control Group shows the approval details and stays on the form', async function (assert) {
    this.put.rejects(new Error('locked'));
    this.parseError.resolves({ message: 'Control Group encountered', response: this.controlGroupResponse });
    const form = new GeneratedItemForm(this.metadata, { username: 'existing', token_num_uses: 1 });
    await this.renderItem('edit', form);
    await click(GENERAL.submitButton);

    assert.true(this.saveTokenFromError.calledWith(this.controlGroupResponse), 'token is saved for a retry');
    assert.dom(GENERAL.messageError).includesText('Control Group Error', 'control group error is shown');
    assert.dom(GENERAL.messageError).includesText('cg-accessor', 'the accessor is shown');
    assert.true(this.transitionTo.notCalled, 'stays on the form');
    assert.true(this.success.notCalled, 'no success flash is shown');
  });

  test('a delete behind a Control Group shows the approval details instead of an error flash', async function (assert) {
    this.delete.rejects(new Error('locked'));
    this.parseError.resolves({ message: 'Control Group encountered', response: this.controlGroupResponse });
    const form = new GeneratedItemForm(this.metadata, { username: 'existing', token_num_uses: 1 });
    await this.renderItem('show', form);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);

    assert.true(this.saveTokenFromError.calledWith(this.controlGroupResponse), 'token is saved for a retry');
    assert.dom(GENERAL.messageError).includesText('cg-accessor', 'control group error is shown');
    assert.true(this.danger.notCalled, 'no error flash is shown');
    assert.true(this.transitionTo.notCalled, 'stays on the item');
  });
});
