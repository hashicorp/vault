/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | secret', function (hooks) {
  setupTest(hooks);

  test('secret api urls', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:secret').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve({});
      },
    });

    adapter.query({}, 'secret', { id: '', backend: 'secret' });
    assert.strictEqual(url, '/v1/secret/', 'query generic url OK');
    assert.strictEqual(method, 'GET', 'query generic method OK');
    assert.deepEqual(options, { data: { list: true } }, 'query generic url OK');

    adapter.queryRecord({}, 'secret', { id: 'foo', backend: 'secret' });
    assert.strictEqual(url, '/v1/secret/foo', 'queryRecord generic url OK');
    assert.strictEqual(method, 'GET', 'queryRecord generic method OK');
  });
});
