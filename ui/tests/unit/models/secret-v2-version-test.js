/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Model | secret-v2-version', function (hooks) {
  setupTest(hooks);

  test('deleted is true for a past deletionTime', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-v2-version', {
          deletionTime: '2000-10-14T00:00:00.000000Z',
        })
      );
      assert.true(model.get('deleted'));
    });
  });

  test('deleted is false for a future deletionTime', function (assert) {
    assert.expect(1);
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-v2-version', {
          deletionTime: '2999-10-14T00:00:00.000000Z',
        })
      );
      assert.false(model.get('deleted'));
    });
  });
});
