/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | kubernetes/config', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.unloadAll('kubernetes/config');
  });

  test('it should make request to correct endpoint when querying record', async function (assert) {
    assert.expect(1);
    this.server.get('/kubernetes-test/config', () => {
      assert.ok('GET request made to correct endpoint when querying record');
    });
    await this.store.queryRecord('kubernetes/config', { backend: 'kubernetes-test' });
  });

  test('it should make request to correct endpoint when creating new record', async function (assert) {
    assert.expect(1);
    this.server.post('/kubernetes-test/config', () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });
    const record = this.store.createRecord('kubernetes/config', { backend: 'kubernetes-test' });
    await record.save();
  });

  test('it should make request to correct endpoint when updating record', async function (assert) {
    assert.expect(1);
    this.server.post('/kubernetes-test/config', () => {
      assert.ok('POST request made to correct endpoint when updating record');
    });
    this.store.pushPayload('kubernetes/config', {
      modelName: 'kubernetes/config',
      backend: 'kubernetes-test',
    });
    const record = this.store.peekRecord('kubernetes/config', 'kubernetes-test');
    await record.save();
  });

  test('it should make request to correct endpoint when deleting record', async function (assert) {
    assert.expect(1);
    this.server.delete('/kubernetes-test/config', () => {
      assert.ok('DELETE request made to correct endpoint when deleting record');
    });
    this.store.pushPayload('kubernetes/config', {
      modelName: 'kubernetes/config',
      backend: 'kubernetes-test',
    });
    const record = this.store.peekRecord('kubernetes/config', 'kubernetes-test');
    await record.destroyRecord();
  });

  test('it should check the config vars endpoint', async function (assert) {
    assert.expect(1);

    this.server.get('/kubernetes-test/check', () => {
      assert.ok('GET request made to config vars check endpoint');
    });

    await this.store.adapterFor('kubernetes/config').checkConfigVars('kubernetes-test');
  });
});
