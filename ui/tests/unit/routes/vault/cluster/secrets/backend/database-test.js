/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import config from 'vault/config/environment';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

const CONTROL_GROUP_RESPONSE = { isControlGroupError: true, accessor: 'cg-accessor', token: 'cg-token' };

module('Unit | Route | vault/cluster/secrets/backend database', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.api = this.owner.lookup('service:api');
    this.capabilities = sinon.stub(this.owner.lookup('service:capabilities'), 'fetch').resolves({});
  });

  module('list', function (hooks) {
    hooks.beforeEach(function () {
      this.route = this.owner.lookup('route:vault.cluster.secrets.backend.list');
    });

    // Capabilities are checked per item, so only the items on the requested page are checked.
    test('it only fetches capabilities for the connections on the current page', async function (assert) {
      const pageSize = config.APP.DEFAULT_PAGE_SIZE;
      const keys = Array.from({ length: pageSize + 2 }, (_, i) => `conn-${i}`);
      sinon.stub(this.api.secrets, 'databaseListConnections').resolves({ keys });
      this.capabilities.resolves({ [`database/config/conn-${pageSize}`]: { canUpdate: true } });

      const connections = await this.route.fetchDatabaseConnections('database', 2);

      const [paths] = this.capabilities.lastCall.args;
      assert.strictEqual(paths.length, 6, 'three paths for each of the two connections on page 2');
      assert.true(connections[0].canEdit, 'capabilities are added to the page items');
      assert.strictEqual(connections.meta.totalItems, pageSize + 2, 'pagination metadata is kept');
    });

    test('it only fetches capabilities for the roles on the current page', async function (assert) {
      sinon.stub(this.api.secrets, 'databaseListStaticRoles').resolves({ keys: ['static-a', 'other'] });
      sinon.stub(this.api.secrets, 'databaseListRoles').resolves({ keys: ['dynamic-a'] });

      const roles = await this.route.fetchDatabaseRoles('database', 1, '-a');

      const [paths] = this.capabilities.lastCall.args;
      assert.strictEqual(paths.length, 6, 'three paths for each of the two filtered roles');
      assert.deepEqual(
        roles.map(({ id }) => id),
        ['static-a', 'dynamic-a'],
        'only filtered roles are returned'
      );
      assert.strictEqual(roles.meta.filteredTotal, 2, 'pagination metadata is kept');
    });

    test('it throws the control group response when listing roles requires approval', async function (assert) {
      sinon
        .stub(this.api.secrets, 'databaseListStaticRoles')
        .rejects(getErrorResponse(CONTROL_GROUP_RESPONSE, 403));
      sinon.stub(this.api.secrets, 'databaseListRoles').rejects(getErrorResponse(undefined, 404));

      await assert.rejects(
        this.route.fetchDatabaseRoles('database', 1),
        (error) => error.isControlGroupError === true,
        'the application route can start the approval flow'
      );
    });

    test('it keeps the status of a forbidden connection list', async function (assert) {
      sinon
        .stub(this.api.secrets, 'databaseListConnections')
        .rejects(getErrorResponse({ errors: ['permission denied'] }, 403));

      await assert.rejects(
        this.route.fetchDatabaseConnections('database', 1),
        (error) => error.httpStatus === 403 && error.message === 'permission denied',
        'the error substate can render the forbidden view'
      );
    });
  });

  module('secret-edit', function (hooks) {
    hooks.beforeEach(function () {
      this.route = this.owner.lookup('route:vault.cluster.secrets.backend.show');
    });

    test('it throws the control group response when reading a connection requires approval', async function (assert) {
      sinon
        .stub(this.api.secrets, 'databaseReadConnectionConfiguration')
        .rejects(getErrorResponse(CONTROL_GROUP_RESPONSE, 403));

      await assert.rejects(
        this.route.fetchDatabaseConnection('database', 'my-conn'),
        (error) => error.isControlGroupError === true,
        'the application route can start the approval flow'
      );
    });

    test('it throws the control group response when reading a role requires approval', async function (assert) {
      sinon
        .stub(this.api.secrets, 'databaseReadStaticRole')
        .rejects(getErrorResponse(CONTROL_GROUP_RESPONSE, 403));
      sinon.stub(this.api.secrets, 'databaseReadRole').rejects(getErrorResponse(undefined, 404));

      await assert.rejects(
        this.route.fetchDatabaseRole('database', 'my-role'),
        (error) => error.isControlGroupError === true,
        'the application route can start the approval flow'
      );
    });
  });
});
