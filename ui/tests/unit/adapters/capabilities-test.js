/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
});
