/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import RoleAwsForm from 'vault/forms/aws-role';
import config from 'vault/config/environment';

module('Unit | Route | aws roles', function (hooks) {
  setupTest(hooks);

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('new AWS roles use a form instead of the generic secret model', function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.create-root');
    sinon.stub(route, 'paramsFor').returns({ backend: 'aws' });
    sinon.stub(route, 'modelType').returns('role-aws');
    const form = route.createModel({ to: { queryParams: {} } });
    assert.true(form instanceof RoleAwsForm);
    assert.true(form.isNew);
    assert.strictEqual(form.data.backend, 'aws');
    assert.strictEqual(form.data.credential_type, 'iam_user');
  });

  test('AWS role routes do not load legacy models before API-backed reads', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.secret-edit');
    sinon.stub(route, 'paramsFor').returns({ backend: 'aws' });
    route.routeName = 'vault.cluster.secrets.backend.show';
    sinon.stub(route, 'modelFor').returns({ type: 'aws' });
    const findRecord = sinon.spy(route.store, 'findRecord');
    await route.beforeModel();
    assert.false(findRecord.called, 'Does not load a legacy Ember Data model');
  });

  test('role lists batch capabilities, filter names and keep role identities', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');
    sinon.stub(route.api.secrets, 'awsListRoles').resolves({ keys: ['alpha', 'beta'] });
    const fetch = sinon.stub(route.capabilitiesService, 'fetch').resolves({
      'aws/roles/alpha': { canRead: true, canUpdate: false, canDelete: false },
      'aws/creds/alpha': { canUpdate: true },
      'aws/roles/beta': { canRead: false, canUpdate: true, canDelete: true },
      'aws/creds/beta': { canUpdate: false },
    });
    const roles = await route.fetchAwsRoles('aws', 1, 'alpha');
    assert.deepEqual(fetch.firstCall.args[0], [
      'aws/roles/alpha',
      'aws/creds/alpha',
      'aws/roles/beta',
      'aws/creds/beta',
    ]);
    assert.strictEqual(fetch.callCount, 1, 'Permissions are fetched in one batch');
    assert.deepEqual(
      roles.map((role) => role.id),
      ['alpha']
    );
    assert.true(roles[0].canRead);
    assert.true(roles[0].canGenerate);
    assert.false(roles[0].canEdit);
    assert.false(roles[0].canDelete);
  });

  test('reads populate the role form from the generated API response', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.secret-edit');
    const read = sinon.stub(route.api.secrets, 'awsReadRole').resolves({
      credential_type: 'assumed_role',
      role_arns: ['arn:example'],
    });
    const form = await route.fetchAwsRole('aws', 'example');
    assert.true(read.calledOnceWithExactly('example', 'aws'), 'Reads the selected role and mount');
    assert.strictEqual(form.data.credential_type, 'assumed_role', 'Uses the returned credential type');
    assert.deepEqual(form.data.role_arns, ['arn:example'], 'Preserves role configuration');
    assert.strictEqual(form.data.backend, 'aws', 'Retains the mount for saving');
    assert.strictEqual(form.id, 'example', 'Retains the role identity');
    assert.false(form.isNew, 'Existing roles are edited rather than created');
  });

  test('empty role lists do not request capabilities', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');
    sinon.stub(route.api.secrets, 'awsListRoles').resolves({ keys: [] });
    const fetch = sinon.spy(route.capabilitiesService, 'fetch');
    const roles = await route.fetchAwsRoles('aws', 1);
    assert.deepEqual(roles, [], 'Empty lists render an empty state');
    assert.strictEqual(roles.meta.totalItems, 0, 'Pagination reports no roles');
    assert.false(fetch.called, 'No empty capabilities request is sent');
  });

  test('role list pagination retains names and permissions on subsequent pages', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');
    const pageSize = config.APP.DEFAULT_PAGE_SIZE;
    const names = Array.from({ length: pageSize + 2 }, (_, index) => `role-${index}`);
    sinon.stub(route.api.secrets, 'awsListRoles').resolves({ keys: names });
    sinon.stub(route.capabilitiesService, 'fetch').resolves({
      [`aws/roles/role-${pageSize}`]: { canRead: true, canUpdate: true, canDelete: true },
      [`aws/creds/role-${pageSize}`]: { canUpdate: true },
    });
    const roles = await route.fetchAwsRoles('aws', 2);
    assert.deepEqual(
      roles.map((role) => role.name),
      names.slice(pageSize),
      'Second page has remaining roles'
    );
    assert.strictEqual(roles.meta.currentPage, 2, 'Reports the requested page');
    assert.strictEqual(roles.meta.totalItems, names.length, 'Reports the complete role count');
    assert.true(roles[0].canEdit, 'Preserves permissions after pagination');
    assert.false(roles[1].canGenerate, 'Missing permissions do not enable generation');
  });

  test('only a missing role list is treated as empty', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');
    const missing = new Error('No roles');
    const denied = new Error('List denied');
    const list = sinon.stub(route.api.secrets, 'awsListRoles');
    list.onFirstCall().rejects(missing);
    list.onSecondCall().rejects(denied);
    const parse = sinon.stub(route.api, 'parseError');
    parse.withArgs(missing).resolves({ status: 404 });
    parse.withArgs(denied).resolves({ status: 403 });
    assert.deepEqual(await route.fetchAwsRoles('aws', 1), [], '404 produces an empty list');
    await assert.rejects(route.fetchAwsRoles('aws', 1), denied, 'Permission failures propagate');
  });

  test('capability failures are not disguised as an empty role list', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.list');
    const error = new Error('Capabilities unavailable');
    sinon.stub(route.api.secrets, 'awsListRoles').resolves({ keys: ['example'] });
    sinon.stub(route.capabilitiesService, 'fetch').rejects(error);
    sinon.stub(route.api, 'parseError').resolves({ status: 500 });
    await assert.rejects(route.fetchAwsRoles('aws', 1), error, 'Failed permissions surface as errors');
  });

  test('credential role reads use the generated API response and only tolerate read denial', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.credentials');
    sinon.stub(route, 'modelFor').returns({ id: 'aws', type: 'aws' });
    const denied = new Error('Read denied');
    const unavailable = new Error('Role unavailable');
    const read = sinon.stub(route.api.secrets, 'awsReadRole');
    read.onFirstCall().resolves({
      credential_type: 'session_token',
    });
    read.onSecondCall().rejects(denied);
    read.onThirdCall().rejects(unavailable);
    const parse = sinon.stub(route.api, 'parseError');
    parse.withArgs(denied).resolves({ status: 403 });
    parse.withArgs(unavailable).resolves({ status: 500 });
    const model = await route.model({ secret: 'example' });
    assert.strictEqual(model.awsRoleType, 'session_token', 'Uses the returned type for generation');
    assert.deepEqual(read.firstCall.args, ['example', 'aws'], 'Reads the correct role and mount');
    assert.strictEqual(
      await route.getAwsRole('aws', 'example'),
      undefined,
      '403 permits manual type selection'
    );
    await assert.rejects(route.getAwsRole('aws', 'example'), unavailable, 'Other errors propagate');
  });

  test('list deletion failures retain the role and display a danger notification', async function (assert) {
    const controller = this.owner.lookup('controller:vault.cluster.secrets.backend.list');
    controller.set('backendType', 'aws');
    const error = new Error('Deletion denied');
    const remove = sinon.stub(controller.api.secrets, 'awsDeleteRole').rejects(error);
    sinon.stub(controller.api, 'parseError').resolves({ message: error.message });
    const danger = sinon.spy(controller.flashMessages, 'danger');
    const success = sinon.spy(controller.flashMessages, 'success');
    const reload = sinon.stub(controller, 'send');
    await controller.actions.delete.call(controller, { id: 'example', backend: 'aws' });
    assert.true(remove.calledOnceWithExactly('example', 'aws'), 'Deletes the selected role');
    assert.true(danger.calledOnceWithExactly('Deletion denied'), 'Reports the API failure');
    assert.false(success.called, 'Does not report success');
    assert.false(reload.called, 'Keeps the list unchanged on failure');
  });
});
