/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/certificate/base', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.owner.lookup('service:secretMountPath').update('pki-test');
  });

  test('it should make request to correct endpoint on queryRecord', async function (assert) {
    assert.expect(1);

    this.server.get('/pki-test/cert/1234', () => {
      assert.ok(true, 'Request made to correct endpoint on queryRecord');
      return { data: {} };
    });

    await this.store.queryRecord('pki/certificate/base', { backend: 'pki-test', id: '1234' });
  });

  test('it should make request to correct endpoint on query', async function (assert) {
    assert.expect(1);

    this.server.get('/pki-test/certs', (schema, req) => {
      assert.strictEqual(req.queryParams.list, 'true', 'Request made to correct endpoint on query');
      return { data: { keys: [] } };
    });

    await this.store.query('pki/certificate/base', { backend: 'pki-test' });
  });

  test('it should make request to correct endpoint on update', async function (assert) {
    assert.expect(1);

    this.store.pushPayload('pki/certificate/base', {
      modelName: 'pki/certificate/base',
      data: {
        serial_number: '1234',
      },
    });

    this.server.post('pki-test/revoke', (schema, req) => {
      assert.deepEqual(
        JSON.parse(req.requestBody),
        { serial_number: '1234' },
        'Request made to correct endpoint on update'
      );
      return { data: {} };
    });

    await this.store.peekRecord('pki/certificate/base', '1234').save();
  });
});
