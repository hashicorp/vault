/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import { SecretsApiTransitListKeysListEnum } from '@hashicorp/vault-client-typescript';

module('Unit | Route | transit list', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.route = this.owner.lookup('route:vault.cluster.secrets.backend.list-root');
    this.api = this.owner.lookup('service:api');
    this.list = sinon
      .stub(this.api.secrets, 'transitListKeys')
      .resolves({ keys: ['aes-key', 'signing-key'] });
    sinon.stub(this.route, 'paramsFor').returns({ backend: 'transit-manual' });
    sinon.stub(this.route, 'modelFor').returns({ engineType: 'transit' });
    sinon.stub(this.route, 'secretParam').returns('');
    this.legacyQuery = sinon.stub(this.owner.lookup('service:pagination'), 'lazyPaginatedQuery');
    sinon.stub(this.owner.lookup('service:capabilities'), 'fetch').resolves({
      'transit-manual/keys/aes-key': { canRead: true, canUpdate: true },
      'transit-manual/keys/signing-key': { canRead: true, canUpdate: false },
    });
  });

  test('existing keys are listed without the removed Ember Data adapter', async function (assert) {
    const { secrets } = await this.route.model({ page: 1 });
    assert.true(
      this.list.calledWithExactly('transit-manual', SecretsApiTransitListKeysListEnum.TRUE),
      'lists keys at the Transit endpoint'
    );
    assert.false(this.legacyQuery.called, 'does not query the removed model');
    assert.deepEqual(
      secrets.map(({ id }) => id),
      ['aes-key', 'signing-key'],
      'existing keys appear'
    );
    assert.strictEqual(secrets.meta.total, 2, 'the list renders its non-empty state');
    assert.true(secrets[0].canEdit, 'update capability enables editing');
    assert.false(secrets[1].canEdit, 'read-only capability does not enable editing');
  });

  test('filtering and subsequent visits use fresh API results', async function (assert) {
    const filtered = await this.route.model({ page: 1, pageFilter: 'AES' });
    assert.deepEqual(
      filtered.secrets.map(({ id }) => id),
      ['aes-key'],
      'filtering is case insensitive'
    );
    this.list.resolves({ keys: ['aes-key', 'signing-key', 'new-key'] });
    const refreshed = await this.route.model({ page: 1 });
    assert.strictEqual(refreshed.secrets.meta.total, 3, 'returning to the list fetches current keys');
  });

  test('only a missing key collection is treated as empty', async function (assert) {
    const error = new Error('request failed');
    this.list.rejects(error);
    const parse = sinon.stub(this.api, 'parseError').resolves({ status: 404 });
    const empty = await this.route.model({ page: 1 });
    assert.strictEqual(empty.secrets.length, 0, 'a 404 is an empty collection');
    parse.resolves({ status: 403 });
    await assert.rejects(this.route.model({ page: 1 }), /request failed/, 'permission errors propagate');
    parse.resolves({ status: 500 });
    await assert.rejects(this.route.model({ page: 1 }), /request failed/, 'server errors propagate');
  });
});
