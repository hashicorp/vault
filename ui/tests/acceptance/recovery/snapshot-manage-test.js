/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { click, fillIn, visit } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import recoveryHandler from 'vault/mirage/handlers/recovery';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const requestTests = (test) => {
  test('it makes the recovery requests in the root namespace', async function (assert) {
    assert.expect(3);

    this.server.post(`/${this.resourceType}/:path`, (schema, req) => {
      assert.strictEqual(
        req.params.path,
        `${this.resourceType}-recovered-data`,
        'url param has expected resource path'
      );

      const expectedHeaders = {
        'x-vault-recover-snapshot-id': this.snapshot_id,
        'x-vault-namespace': '',
      };

      Object.entries(expectedHeaders).forEach(([key, value]) => {
        assert.strictEqual(req.requestHeaders[key], value, `"${key}" header has expected value "${value}"`);
      });
    });

    await click(GENERAL.selectByAttr('mount'));
    await click(`[data-option-index="${this.mountIdx}"]`);
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.resourceType}-recovered-data`);
    await click(GENERAL.button('recover'));
  });

  test('it makes the recovery requests in the child namespace', async function (assert) {
    assert.expect(3);

    this.server.post(`/${this.resourceType}/:path`, (schema, req) => {
      assert.strictEqual(
        req.params.path,
        `${this.resourceType}-recovered-data`,
        'url param has expected resource path'
      );

      const expectedHeaders = {
        'x-vault-recover-snapshot-id': this.snapshot_id,
        'x-vault-namespace': 'child-ns-1',
      };

      Object.entries(expectedHeaders).forEach(([key, value]) => {
        assert.strictEqual(req.requestHeaders[key], value, `"${key}" header has expected value "${value}"`);
      });
    });

    await click(GENERAL.selectByAttr('namespace'));
    await click('[data-option-index="1"]');
    await click(GENERAL.selectByAttr('mount'));
    await click(`[data-option-index="${this.mountIdx}"]`);
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.resourceType}-recovered-data`);
    await click(GENERAL.button('recover'));
  });

  test('it makes the recovery requests in the grandchild namespace', async function (assert) {
    assert.expect(3);

    this.server.post(`/${this.resourceType}/:path`, (schema, req) => {
      assert.strictEqual(
        req.params.path,
        `${this.resourceType}-recovered-data`,
        'url param has expected resource path'
      );

      const expectedHeaders = {
        'x-vault-recover-snapshot-id': this.snapshot_id,
        'x-vault-namespace': 'child-ns-1/nested',
      };

      Object.entries(expectedHeaders).forEach(([key, value]) => {
        assert.strictEqual(req.requestHeaders[key], value, `"${key}" header has expected value "${value}"`);
      });
    });

    await click(GENERAL.selectByAttr('namespace'));
    await click('[data-option-index="2"]');
    await click(GENERAL.selectByAttr('mount'));
    await click(`[data-option-index="${this.mountIdx}"]`);
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.resourceType}-recovered-data`);
    await click(GENERAL.button('recover'));
  });
};

module('Acceptance | recovery | snapshot-manage', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    recoveryHandler(this.server);

    const snapshot = this.server.create('snapshot');
    this.snapshot_id = snapshot.snapshot_id;
    await login();
    await visit('/vault/recovery/snapshots');
  });

  module('cubbyhole', function (hooks) {
    hooks.beforeEach(function () {
      this.resourceType = 'cubbyhole';
      this.mountIdx = '1.0';
    });

    requestTests(test);
  });

  module('kvv1', function (hooks) {
    hooks.beforeEach(function () {
      this.resourceType = 'kv';
      this.mountIdx = '1.1';
    });

    requestTests(test);
  });

  module('database static-roles', function (hooks) {
    hooks.beforeEach(function () {
      this.resourceType = 'database/static-roles';
      this.mountIdx = '0.0';
    });

    requestTests(test);
  });
});
