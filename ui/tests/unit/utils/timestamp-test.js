/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import timestamp from 'core/utils/timestamp';
import sinon from 'sinon';
import { module, test } from 'qunit';

/*
  This test coverage is more for an example than actually covering the utility
*/
module('Unit | Utility | timestamp', function () {
  test('it can be overridden', function (assert) {
    sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2030-03-03T03:30:03')));
    const result = timestamp.now();
    assert.strictEqual(result.toISOString(), new Date('2030-03-03T03:30:03').toISOString());
  });
});
