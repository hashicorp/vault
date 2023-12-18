/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | ldap/config', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.unloadAll('ldap/config');
  });

  test('it should make request to correct endpoint when querying record', async function (assert) {
    assert.expect(1);
    this.server.get('/ldap-test/config', () => {
      assert.ok('GET request made to correct endpoint when querying record');
    });
    await this.store.queryRecord('ldap/config', { backend: 'ldap-test' });
  });

  test('it should make request to correct endpoint when creating new record', async function (assert) {
    assert.expect(1);
    this.server.post('/ldap-test/config', () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });
    const record = this.store.createRecord('ldap/config', { backend: 'ldap-test' });
    await record.save();
  });

  test('it should make request to correct endpoint when updating record', async function (assert) {
    assert.expect(1);
    this.server.post('/ldap-test/config', () => {
      assert.ok('POST request made to correct endpoint when updating record');
    });
    this.store.pushPayload('ldap/config', {
      modelName: 'ldap/config',
      backend: 'ldap-test',
    });
    const record = this.store.peekRecord('ldap/config', 'ldap-test');
    await record.save();
  });

  test('it should make request to correct endpoint when deleting record', async function (assert) {
    assert.expect(1);
    this.server.delete('/ldap-test/config', () => {
      assert.ok('DELETE request made to correct endpoint when deleting record');
    });
    this.store.pushPayload('ldap/config', {
      modelName: 'ldap/config',
      backend: 'ldap-test',
    });
    const record = this.store.peekRecord('ldap/config', 'ldap-test');
    await record.destroyRecord();
  });
});
