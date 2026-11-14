/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/secrets-redirect-with-path', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.route = this.owner.lookup('route:vault/cluster/secrets-redirect-with-path');
    this.originalReplaceWith = this.router.replaceWith;
    this.router.replaceWith = sinon.stub();
  });

  hooks.afterEach(function () {
    this.router.replaceWith = this.originalReplaceWith;
  });

  test('it preserves the query string, including the KV version, when redirecting a path', function (assert) {
    this.route.beforeModel({
      to: { params: { path: 'engine/kv/secret/details' }, queryParams: { version: '1' } },
    });

    assert.true(
      this.router.replaceWith.calledWithExactly('/vault/secrets-engines/engine/kv/secret/details?version=1'),
      'redirects to the new route with the version query param intact'
    );
  });

  test('it redirects without a query string when there are no query params', function (assert) {
    this.route.beforeModel({
      to: { params: { path: 'engine/kv/secret/details' }, queryParams: {} },
    });

    assert.true(
      this.router.replaceWith.calledWithExactly('/vault/secrets-engines/engine/kv/secret/details'),
      'redirects to the new route with no trailing question mark'
    );
  });

  test('it redirects to the base secrets page when there is no path', function (assert) {
    this.route.beforeModel({ to: { params: {}, queryParams: {} } });

    assert.true(
      this.router.replaceWith.calledWithExactly('vault.cluster.secrets'),
      'redirects to the base secrets route'
    );
  });
});
