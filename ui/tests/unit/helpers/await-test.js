/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import AwaitHelper from 'vault/helpers/await';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { waitUntil } from '@ember/test-helpers';
import { Promise } from 'rsvp';
import { later } from '@ember/runloop';
import sinon from 'sinon';

// recompute triggers a rerender which isn't going to work for unit tests
// override method to trigger compute instead
class AwaitHelperForTesting extends AwaitHelper {
  compute([promise]) {
    // cache original promise arg to simulate rerender when calling recompute
    this._promiseArg = promise;
    return super.compute(...arguments);
  }
  recompute() {
    this.compute([this._promiseArg]);
  }
}

module('Unit | Helpers | await', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.helper = new AwaitHelperForTesting();
    this.spy = sinon.spy(this.helper, 'compute');
  });

  test('it returns value when input is not a promise', async function (assert) {
    this.helper.compute(['foo']);
    assert.strictEqual(this.spy.returnValues[0], 'foo', 'Input value returned when not promise');
  });

  test('it returns null default value and then resolved value', async function (assert) {
    const promise = new Promise((resolve) => resolve('foo'));
    this.helper.compute([promise]);
    await waitUntil(() => this.spy.returnValues[1]);
    assert.strictEqual(this.spy.returnValues[0], null, 'Default value returned while promise resolves');
    assert.strictEqual(this.spy.returnValues[1], 'foo', 'Resolved value is returned');
  });

  test('it returns rejected value', async function (assert) {
    const promise = new Promise((resolve, reject) => reject('bar'));
    this.helper.compute([promise]);
    await waitUntil(() => this.spy.returnValues[1]);
    assert.strictEqual(this.spy.returnValues[1], 'bar', 'Rejected value is returned');
  });

  test('it returns then value', async function (assert) {
    const promise = new Promise((resolve) => resolve('foo')).then(() => 'new resolve value');
    this.helper.compute([promise]);
    await waitUntil(() => this.spy.returnValues[1]);
    assert.strictEqual(this.spy.returnValues[1], 'new resolve value', 'Value from then is returned');
  });

  test('it returns catch value', async function (assert) {
    const promise = new Promise((resolve, reject) => reject('bar')).catch(() => 'new reject value');
    this.helper.compute([promise]);
    await waitUntil(() => this.spy.returnValues[1]);
    assert.strictEqual(this.spy.returnValues[1], 'new reject value', 'Value from catch is returned');
  });

  test('it always returns value from latest promise', async function (assert) {
    const promise1 = new Promise((resolve) => later(() => resolve('foo'), 500));
    const promise2 = new Promise((resolve) => resolve('bar'));
    this.helper.compute([promise1]);
    this.helper.compute([promise2]);
    // allow first promise time to resolve
    await waitUntil(() => later(() => true, 500));
    assert.strictEqual(this.spy.returnValues[2], 'bar', 'Latest promise value is returned');
  });
});
