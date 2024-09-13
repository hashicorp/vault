/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Adapter | capabilities', function (hooks) {
  setupTest(hooks);

  test('calls the correct url', function (assert) {
    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:capabilities').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.findRecord(null, 'capabilities', 'foo');
    assert.strictEqual(url, '/v1/sys/capabilities-self', 'calls the correct URL');
    assert.deepEqual({ paths: ['foo'] }, options.data, 'data params OK');
    assert.strictEqual(method, 'POST', 'method OK');
  });

  test('enterprise calls the correct url within namespace when userRoot = root', function (assert) {
    const namespaceSvc = this.owner.lookup('service:namespace');
    namespaceSvc.setNamespace('admin');

    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:capabilities').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.findRecord(null, 'capabilities', 'foo');
    assert.strictEqual(url, '/v1/sys/capabilities-self', 'calls the correct URL');
    assert.deepEqual({ paths: ['admin/foo'] }, options.data, 'data params prefix paths with namespace');
    assert.strictEqual(options.namespace, '', 'sent with root namespace');
    assert.strictEqual(method, 'POST', 'method OK');
  });

  test('enterprise calls the correct url within namespace when userRoot is not root', function (assert) {
    const namespaceSvc = this.owner.lookup('service:namespace');
    const auth = this.owner.lookup('service:auth');
    namespaceSvc.setNamespace('admin/bar/baz');
    // Set user root namespace
    auth.setCluster('1');
    auth.set('tokens', ['vault-_root_☃1']);
    auth.setTokenData('vault-_root_☃1', { userRootNamespace: 'admin/bar', backend: { mountPath: 'token' } });

    let url, method, options;
    const adapter = this.owner.factoryFor('adapter:capabilities').create({
      ajax: (...args) => {
        [url, method, options] = args;
        return resolve();
      },
    });

    adapter.findRecord(null, 'capabilities', 'foo');
    assert.strictEqual(url, '/v1/sys/capabilities-self', 'calls the correct URL');
    assert.deepEqual({ paths: ['baz/foo'] }, options.data, 'data params prefix path with relative namespace');
    assert.strictEqual(options.namespace, 'admin/bar', 'sent with root namespace');
    assert.strictEqual(method, 'POST', 'method OK');
  });
});
