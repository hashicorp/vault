/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | secrets backend | unknown item type', function (hooks) {
  setupTest(hooks);

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('reading an unrecognized item type renders a not found error', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.secret-edit');
    sinon.stub(route, 'paramsFor').returns({ backend: 'keymgmt', secret: 'foo' });
    sinon.stub(route, 'modelType').returns('keymgmt/unknown');

    await assert.rejects(
      route.model({}, { to: { queryParams: {} } }),
      (error) => error.httpStatus === 404 && error.path === 'keymgmt/foo',
      'the route throws a 404 instead of querying a model that does not exist'
    );
  });

  test('creating an unrecognized item type renders a not found error', function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.create-root');
    sinon.stub(route, 'paramsFor').returns({ backend: 'keymgmt' });
    sinon.stub(route, 'modelType').returns('keymgmt/unknown');

    assert.throws(
      () => route.createModel({ to: { queryParams: {} } }),
      (error) => error.httpStatus === 404 && error.path === 'keymgmt',
      'the route throws a 404 instead of creating a model that does not exist'
    );
  });
});
