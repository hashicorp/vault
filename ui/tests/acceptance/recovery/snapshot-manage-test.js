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
import { overrideResponse } from 'vault/tests/helpers/stubs';

const requestTests = (test) => {
  test('it makes the recovery request in the root namespace', async function (assert) {
    assert.expect(3);

    this.server.post(`/${this.type}/:path`, (schema, req) => {
      assert.strictEqual(
        req.params.path,
        `${this.type}-recovered-data`,
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
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    await click(GENERAL.button('recover'));
  });

  test('it makes the recovery request in the child namespace', async function (assert) {
    assert.expect(3);

    this.server.post(`/${this.type}/:path`, (schema, req) => {
      assert.strictEqual(
        req.params.path,
        `${this.type}-recovered-data`,
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
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    await click(GENERAL.button('recover'));
  });

  test('it makes the recovery request in the grandchild namespace', async function (assert) {
    assert.expect(3);

    this.server.post(`/${this.type}/:path`, (schema, req) => {
      assert.strictEqual(
        req.params.path,
        `${this.type}-recovered-data`,
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
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    await click(GENERAL.button('recover'));
  });

  test('it recovers to a copy in the root namespace', async function (assert) {
    assert.expect(5);

    this.server.post(`/${this.type}/:path`, (schema, req) => {
      assert.strictEqual(req.params.path, `${this.type}-new-path`, 'url param is the new "copy" path');

      const expectedHeaders = {
        'x-vault-recover-snapshot-id': this.snapshot_id,
        'x-vault-recover-source-path': `${this.type}/${this.type}-recovered-data`,
        'x-vault-namespace': '',
      };

      Object.entries(expectedHeaders).forEach(([key, value]) => {
        assert.strictEqual(req.requestHeaders[key], value, `"${key}" header has expected value "${value}"`);
      });
    });

    await click(GENERAL.selectByAttr('mount'));
    await click(`[data-option-index="${this.mountIdx}"]`);
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    // Select "Recover to a new path"
    await click(GENERAL.inputByAttr('copy'));
    assert
      .dom(GENERAL.inputByAttr('copyPath'))
      .hasValue(`${this.type}-recovered-data-copy`, 'it appends "-copy" to original resource path');
    await fillIn(GENERAL.inputByAttr('copyPath'), `${this.type}-new-path`);
    await click(GENERAL.button('recover'));
  });

  test('it recovers to a copy in the child namespace', async function (assert) {
    assert.expect(5);

    this.server.post(`/${this.type}/:path`, (schema, req) => {
      assert.strictEqual(req.params.path, `${this.type}-new-path`, 'url param is the new "copy" path');

      const expectedHeaders = {
        'x-vault-recover-snapshot-id': this.snapshot_id,
        'x-vault-recover-source-path': `${this.type}/${this.type}-recovered-data`,
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
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    // Select "Recover to a new path"
    await click(GENERAL.inputByAttr('copy'));
    assert
      .dom(GENERAL.inputByAttr('copyPath'))
      .hasValue(`${this.type}-recovered-data-copy`, 'it appends "-copy" to original resource path');
    await fillIn(GENERAL.inputByAttr('copyPath'), `${this.type}-new-path`);
    await click(GENERAL.button('recover'));
  });

  test('it recovers to a copy in the grandchild namespace', async function (assert) {
    assert.expect(5);

    this.server.post(`/${this.type}/:path`, (schema, req) => {
      assert.strictEqual(req.params.path, `${this.type}-new-path`, 'url param is the new "copy" path');

      const expectedHeaders = {
        'x-vault-recover-snapshot-id': this.snapshot_id,
        'x-vault-recover-source-path': `${this.type}/${this.type}-recovered-data`,
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
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    // Select "Recover to a new path"
    await click(GENERAL.inputByAttr('copy'));
    assert
      .dom(GENERAL.inputByAttr('copyPath'))
      .hasValue(`${this.type}-recovered-data-copy`, 'it appends "-copy" to original resource path');
    await fillIn(GENERAL.inputByAttr('copyPath'), `${this.type}-new-path`);
    await click(GENERAL.button('recover'));
  });

  test('it makes the read request in the root namespace', async function (assert) {
    assert.expect(3);

    this.server.get(`/${this.type}/:path`, (schema, req) => {
      const { params, queryParams, requestHeaders } = req;
      assert.strictEqual(params.path, `${this.type}-recovered-data`, 'url param is resource path');
      assert.strictEqual(queryParams.read_snapshot_id, this.snapshot_id, 'query param has snapshot_id');
      assert.strictEqual(
        requestHeaders['x-vault-namespace'],
        '',
        `'x-vault-namespace' header is an empty string`
      );
    });

    await click(GENERAL.selectByAttr('mount'));
    await click(`[data-option-index="${this.mountIdx}"]`);
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    await click(GENERAL.button('read'));
  });

  test('it makes the read request in the child namespace', async function (assert) {
    assert.expect(3);

    this.server.get(`/${this.type}/:path`, (schema, req) => {
      const { params, queryParams, requestHeaders } = req;
      assert.strictEqual(params.path, `${this.type}-recovered-data`, 'url param is resource path');
      assert.strictEqual(queryParams.read_snapshot_id, this.snapshot_id, 'query param has snapshot_id');
      assert.strictEqual(
        requestHeaders['x-vault-namespace'],
        'child-ns-1',
        `'x-vault-namespace' header has namespace`
      );
    });

    await click(GENERAL.selectByAttr('namespace'));
    await click('[data-option-index="1"]');
    await click(GENERAL.selectByAttr('mount'));
    await click(`[data-option-index="${this.mountIdx}"]`);
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    await click(GENERAL.button('read'));
  });

  test('it makes the read request in the grandchild namespace', async function (assert) {
    assert.expect(3);

    this.server.get(`/${this.type}/:path`, (schema, req) => {
      const { params, queryParams, requestHeaders } = req;
      assert.strictEqual(params.path, `${this.type}-recovered-data`, 'url param is resource path');
      assert.strictEqual(queryParams.read_snapshot_id, this.snapshot_id, 'query param has snapshot_id');
      assert.strictEqual(
        requestHeaders['x-vault-namespace'],
        'child-ns-1/nested',
        `'x-vault-namespace' header has namespace`
      );
    });

    await click(GENERAL.selectByAttr('namespace'));
    await click('[data-option-index="2"]');
    await click(GENERAL.selectByAttr('mount'));
    await click(`[data-option-index="${this.mountIdx}"]`);
    await fillIn(GENERAL.inputByAttr('resourcePath'), `${this.type}-recovered-data`);
    await click(GENERAL.button('read'));
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
      this.type = 'cubbyhole';
      this.mountIdx = '1.0';
    });

    requestTests(test);
  });

  module('kvv1', function (hooks) {
    hooks.beforeEach(function () {
      this.type = 'kv';
      this.mountIdx = '1.1';
    });

    requestTests(test);

    test('it makes the recovery request to custom mount paths', async function (assert) {
      assert.expect(5);
      // Stub no permissions to render manual input so we can supply a custom path
      this.server.get('/sys/internal/ui/mounts', () => overrideResponse(403));
      // Re-visit route to re-fire request to mounts endpoint
      await visit('/vault/recovery/snapshots');

      const customMountPath = `custom-kvv1`;
      this.server.post(`/${customMountPath}/:path`, (schema, req) => {
        assert.true(true, 'it makes request to custom path URL');
        assert.strictEqual(req.params.path, `kvv1-recovered-data`, ':path param is expected value');

        const expectedHeaders = {
          'x-vault-recover-snapshot-id': this.snapshot_id,
          'x-vault-namespace': '',
        };

        Object.entries(expectedHeaders).forEach(([key, value]) => {
          assert.strictEqual(req.requestHeaders[key], value, `"${key}" header has expected value "${value}"`);
        });
      });

      await click(GENERAL.button('Type'));
      // First select a different type to ensure input resets as expected
      await click(GENERAL.radioByAttr('database'));
      await fillIn(GENERAL.inputByAttr('manual-mount-path'), 'bad-path');
      // Now select type we actually want to use
      await click(GENERAL.button('Type'));
      await click(GENERAL.radioByAttr('kv'));
      assert.dom(GENERAL.inputByAttr('manual-mount-path')).hasValue('', 'input clears when type changes');
      await fillIn(GENERAL.inputByAttr('manual-mount-path'), customMountPath);
      await fillIn(GENERAL.inputByAttr('resourcePath'), `kvv1-recovered-data`);
      await click(GENERAL.button('recover'));
    });
  });

  module('database static-roles', function (hooks) {
    hooks.beforeEach(function () {
      this.type = 'database/static-roles';
      this.mountIdx = '0.0';
    });

    requestTests(test);

    test('it makes the recovery request to custom mount paths', async function (assert) {
      assert.expect(5);
      // Stub no permissions to render manual input so we can supply a custom path
      this.server.get('/sys/internal/ui/mounts', () => overrideResponse(403));
      // Re-visit route to re-fire request to mounts endpoint
      await visit('/vault/recovery/snapshots');

      const customMountPath = `custom-database`;
      this.server.post(`/${customMountPath}/static-roles/:path`, (schema, req) => {
        assert.true(true, 'it makes request to custom path URL');
        assert.strictEqual(req.params.path, `database-recovered-data`, ':path param is expected value');

        const expectedHeaders = {
          'x-vault-recover-snapshot-id': this.snapshot_id,
          'x-vault-namespace': '',
        };

        Object.entries(expectedHeaders).forEach(([key, value]) => {
          assert.strictEqual(req.requestHeaders[key], value, `"${key}" header has expected value "${value}"`);
        });
      });

      await click(GENERAL.button('Type'));
      // First select a different type to ensure input resets as expected
      await click(GENERAL.radioByAttr('kv'));
      await fillIn(GENERAL.inputByAttr('manual-mount-path'), 'bad-path');
      // Now select type we actually want to use
      await click(GENERAL.button('Type'));
      await click(GENERAL.radioByAttr('database'));
      assert.dom(GENERAL.inputByAttr('manual-mount-path')).hasValue('', 'input clears when type changes');
      await fillIn(GENERAL.inputByAttr('manual-mount-path'), customMountPath);
      await fillIn(GENERAL.inputByAttr('resourcePath'), `database-recovered-data`);
      await click(GENERAL.button('recover'));
    });
  });
});
