/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | transit key', function (hooks) {
  setupTest(hooks);

  test('transit api urls', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:transit-key').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve({});
      },
    });

    adapter.query({}, 'transit-key', { id: '', backend: 'transit' });
    assert.strictEqual(url, '/v1/transit/keys/', 'query list url OK');
    assert.strictEqual(method, 'GET', 'query list method OK');
    assert.deepEqual(options, { data: { list: true } }, 'query generic url OK');

    adapter.queryRecord({}, 'transit-key', { id: 'foo', backend: 'transit' });
    assert.strictEqual(url, '/v1/transit/keys/foo', 'queryRecord generic url OK');
    assert.strictEqual(method, 'GET', 'queryRecord generic method OK');

    adapter.keyAction('rotate', { backend: 'transit', id: 'foo', payload: {} });
    assert.strictEqual(url, '/v1/transit/keys/foo/rotate', 'keyAction:rotate url OK');

    adapter.keyAction('encrypt', { backend: 'transit', id: 'foo', payload: {} });
    assert.strictEqual(url, '/v1/transit/encrypt/foo', 'keyAction:encrypt url OK');

    adapter.keyAction('datakey', { backend: 'transit', id: 'foo', payload: { param: 'plaintext' } });
    assert.strictEqual(url, '/v1/transit/datakey/plaintext/foo', 'keyAction:datakey url OK');

    adapter.keyAction('export', { backend: 'transit', id: 'foo', payload: { param: ['hmac'] } });
    assert.strictEqual(url, '/v1/transit/export/hmac-key/foo', 'transitAction:export, no version url OK');

    adapter.keyAction('export', { backend: 'transit', id: 'foo', payload: { param: ['hmac', 10] } });
    assert.strictEqual(
      url,
      '/v1/transit/export/hmac-key/foo/10',
      'transitAction:export, with version url OK'
    );
  });
});
