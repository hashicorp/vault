/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | ldap/library', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.adapter = this.store.adapterFor('ldap/library');
  });

  test('it should make request to correct endpoint when listing records', async function (assert) {
    assert.expect(1);

    this.server.get('/ldap-test/library', (schema, req) => {
      assert.ok(req.queryParams.list, 'GET request made to correct endpoint when listing records');
      return { data: { keys: ['test-library'] } };
    });

    await this.store.query('ldap/library', { backend: 'ldap-test' });
  });

  test('it should make request to correct endpoint when querying record', async function (assert) {
    assert.expect(1);

    this.server.get('/ldap-test/library/test-library', () => {
      assert.ok('GET request made to correct endpoint when querying record');
    });

    await this.store.queryRecord('ldap/library', { backend: 'ldap-test', name: 'test-library' });
  });

  test('it should make request to correct endpoint when creating new record', async function (assert) {
    assert.expect(1);

    this.server.post('/ldap-test/library/test-library', () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });

    await this.store.createRecord('ldap/library', { backend: 'ldap-test', name: 'test-library' }).save();
  });

  test('it should make request to correct endpoint when updating record', async function (assert) {
    assert.expect(1);

    this.server.post('/ldap-test/library/test-library', () => {
      assert.ok('POST request made to correct endpoint when updating record');
    });

    this.store.pushPayload('ldap/library', {
      modelName: 'ldap/library',
      backend: 'ldap-test',
      name: 'test-library',
    });

    await this.store.peekRecord('ldap/library', 'test-library').save();
  });

  test('it should make request to correct endpoint when deleting record', async function (assert) {
    assert.expect(1);

    this.server.delete('/ldap-test/library/test-library', () => {
      assert.ok('DELETE request made to correct endpoint when deleting record');
    });

    this.store.pushPayload('ldap/library', {
      modelName: 'ldap/library',
      backend: 'ldap-test',
      name: 'test-library',
    });

    await this.store.peekRecord('ldap/library', 'test-library').destroyRecord();
  });

  test('it should make request to correct endpoint when fetching check-out status', async function (assert) {
    assert.expect(1);

    this.server.get('/ldap-test/library/test-library/status', () => {
      assert.ok('GET request made to correct endpoint when fetching check-out status');
    });

    await this.adapter.fetchStatus('ldap-test', 'test-library');
  });

  test('it should make request to correct endpoint when checking out library', async function (assert) {
    assert.expect(1);

    this.server.post('/ldap-test/library/test-library/check-out', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.strictEqual(json.ttl, '1h', 'POST request made to correct endpoint when checking out library');
      return {
        data: { password: 'test', service_account_name: 'foo@bar.com' },
      };
    });

    await this.adapter.checkOutAccount('ldap-test', 'test-library', '1h');
  });

  test('it should make request to correct endpoint when checking in service accounts', async function (assert) {
    assert.expect(1);

    this.server.post('/ldap-test/library/test-library/check-in', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      assert.deepEqual(
        json.service_account_names,
        ['foo@bar.com'],
        'POST request made to correct endpoint when checking in service accounts'
      );
      return {
        data: {
          'check-ins': ['foo@bar.com'],
        },
      };
    });

    await this.adapter.checkInAccount('ldap-test', 'test-library', ['foo@bar.com']);
  });
});
