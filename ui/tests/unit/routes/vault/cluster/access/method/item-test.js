/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import config from 'vault/config/environment';
import { fetchGeneratedItemForm } from 'vault/routes/vault/cluster/access/method/item/show';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

const CONTROL_GROUP_RESPONSE = {
  isControlGroupError: true,
  accessor: 'cg-accessor',
  token: 'cg-token',
  creation_path: 'auth/userpass/users/bob',
};

module('Unit | Route | vault/cluster/access/method/item', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.api = this.owner.lookup('service:api');
    this.requestGet = sinon.stub(this.api.request, 'get');
    this.parseError = sinon.stub(this.api, 'parseError');
    const username = new FormField('username', 'string', { readOnly: true });
    const numUses = new FormField('token_num_uses', 'number');
    this.itemModel = {
      itemType: 'user',
      formMetadata: {
        formFields: [username, numUses],
        formFieldGroups: [new FormFieldGroup('default', [username, numUses])],
        identifier: 'username',
      },
      urls: {
        list: '/auth/userpass/users/',
        read: (id) => `/auth/userpass/users/${id}`,
      },
    };
    this.json = (body) => ({ json: () => Promise.resolve(body) });
  });

  module('list', function (hooks) {
    hooks.beforeEach(function () {
      this.route = this.owner.lookup('route:vault.cluster.access.method.item.list');
      sinon.stub(this.route, 'modelFor').returns(this.itemModel);
      this.params = sinon.stub(this.route, 'paramsFor').returns({});
    });

    test('it lists and paginates item keys', async function (assert) {
      const keys = Array.from({ length: 101 }, (_, i) => `user-${i}`);
      this.requestGet.resolves(this.json({ data: { keys } }));
      this.params.returns({ page: '2' });

      const items = await this.route.model();
      const pageSize = config.APP.DEFAULT_PAGE_SIZE;
      assert.true(this.requestGet.calledOnceWith('/auth/userpass/users/?list=true'), 'lists the item path');
      assert.deepEqual(
        items[0],
        { id: `user-${pageSize}` },
        'the requested page starts after the first page'
      );
      assert.strictEqual(items.meta.currentPage, 2, 'pagination metadata is attached');
      assert.strictEqual(items.meta.totalItems, 101, 'total includes every key');
    });

    test('it filters item keys', async function (assert) {
      this.requestGet.resolves(this.json({ data: { keys: ['alice', 'bob', 'bobby'] } }));
      this.params.returns({ pageFilter: 'bob' });

      const items = await this.route.model();
      assert.deepEqual(
        items.map(({ id }) => id),
        ['bob', 'bobby'],
        'only matching keys are returned'
      );
    });

    test('it returns an empty list when the mount has no items', async function (assert) {
      this.requestGet.rejects(new Error('not found'));
      this.parseError.resolves({ status: 404 });
      assert.deepEqual(await this.route.model(), [], 'a 404 is an empty list');
    });

    test('it rethrows a Control Group error for the application route to handle', async function (assert) {
      this.requestGet.rejects(new Error('locked'));
      this.parseError.resolves({ status: 403, response: CONTROL_GROUP_RESPONSE });
      await assert.rejects(
        this.route.model(),
        (error) => error === CONTROL_GROUP_RESPONSE,
        'the control group response is thrown'
      );
    });

    test('it rethrows other errors', async function (assert) {
      const error = new Error('permission denied');
      this.requestGet.rejects(error);
      this.parseError.resolves({ status: 403, response: { errors: ['permission denied'] } });
      await assert.rejects(this.route.model(), (thrown) => thrown === error, 'the original error is thrown');
    });
  });

  module('show', function (hooks) {
    hooks.beforeEach(function () {
      this.route = this.owner.lookup('route:vault.cluster.access.method.item.show');
      sinon.stub(this.route, 'modelFor').returns(this.itemModel);
    });

    test('it reads the item into a form with the identifier from the URL', async function (assert) {
      this.requestGet.resolves(this.json({ data: { token_num_uses: 5 } }));

      const form = await fetchGeneratedItemForm(this.route, 'bob');
      assert.true(this.requestGet.calledOnceWith('/auth/userpass/users/bob'), 'reads the item path');
      assert.strictEqual(form.id, 'bob', 'identifier is set from the URL');
      assert.strictEqual(form.token_num_uses, 5, 'response data is loaded');
      assert.false(form.isNew, 'form is for an existing item');
    });

    test('it rethrows a Control Group error for the application route to handle', async function (assert) {
      this.requestGet.rejects(new Error('locked'));
      this.parseError.resolves({ status: 403, response: CONTROL_GROUP_RESPONSE });
      await assert.rejects(
        fetchGeneratedItemForm(this.route, 'bob'),
        (error) => error === CONTROL_GROUP_RESPONSE,
        'the control group response is thrown'
      );
    });

    test('it rethrows other errors', async function (assert) {
      const error = new Error('permission denied');
      this.requestGet.rejects(error);
      this.parseError.resolves({ status: 403, response: { errors: ['permission denied'] } });
      await assert.rejects(
        fetchGeneratedItemForm(this.route, 'bob'),
        (thrown) => thrown === error,
        'the original error is thrown'
      );
    });
  });
});
