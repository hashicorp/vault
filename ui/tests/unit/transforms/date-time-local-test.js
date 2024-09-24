/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';

import { setupTest } from 'vault/tests/helpers';

module('Unit | Transform | date time local', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.transform = this.owner.lookup('transform:date-time-local');
  });

  test('it serializes correctly for the API', function (assert) {
    assert.ok(this.transform);
    let serialized = this.transform.serialize('2024-01-31T00:00');
    assert.strictEqual(
      serialized,
      new Date('2024-01-31T00:00').toISOString(),
      'should serialize a string that is not in ISO format'
    );
    serialized = this.transform.serialize(new Date('2024-03-30T17:11:00Z'));
    assert.strictEqual(serialized, '2024-03-30T17:11:00.000Z', 'should serialize a date object');
    serialized = this.transform.serialize('2024-03-30T17:11:00.000Z');
    assert.strictEqual(serialized, '2024-03-30T17:11:00.000Z', 'should always show an ISO string');
  });
});
