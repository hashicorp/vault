/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | tools', function (hooks) {
  setupTest(hooks);

  test('wrapping api urls', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:tools').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    let clientToken;
    let data = { foo: 'bar' };
    adapter.toolAction('wrap', data, { wrapTTL: '30m' });
    assert.strictEqual(url, '/v1/sys/wrapping/wrap', 'wrapping:wrap url OK');
    assert.strictEqual(method, 'POST', 'wrapping:wrap method OK');
    assert.deepEqual({ data: data, wrapTTL: '30m', clientToken }, options, 'wrapping:wrap options OK');

    data = { token: 'token' };
    adapter.toolAction('lookup', data);
    assert.strictEqual(url, '/v1/sys/wrapping/lookup', 'wrapping:lookup url OK');
    assert.strictEqual(method, 'POST', 'wrapping:lookup method OK');
    assert.deepEqual({ data, clientToken }, options, 'wrapping:lookup options OK');

    adapter.toolAction('unwrap', data);
    assert.strictEqual(url, '/v1/sys/wrapping/unwrap', 'wrapping:unwrap url OK');
    assert.strictEqual(method, 'POST', 'wrapping:unwrap method OK');
    assert.deepEqual({ data, clientToken }, options, 'wrapping:unwrap options OK');

    adapter.toolAction('rewrap', data);
    assert.strictEqual(url, '/v1/sys/wrapping/rewrap', 'wrapping:rewrap url OK');
    assert.strictEqual(method, 'POST', 'wrapping:rewrap method OK');
    assert.deepEqual({ data, clientToken }, options, 'wrapping:rewrap options OK');
  });

  test('tools api urls', function (assert) {
    let url, method;
    const adapter = this.owner.factoryFor('adapter:tools').create({
      ajax: (...args) => {
        [url, method] = args;
        return resolve();
      },
    });

    adapter.toolAction('hash', { input: 'someBase64' });
    assert.strictEqual(url, '/v1/sys/tools/hash', 'sys tools hash: url OK');
    assert.strictEqual(method, 'POST', 'sys tools hash: method OK');

    adapter.toolAction('random', { bytes: '32' });
    assert.strictEqual(url, '/v1/sys/tools/random', 'sys tools random: url OK');
    assert.strictEqual(method, 'POST', 'sys tools random: method OK');
  });
});
