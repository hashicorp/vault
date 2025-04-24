/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hasCapability } from 'core/helpers/has-capability';
import { module, test } from 'qunit';

module('Unit | Helpers | has-capability', function (hooks) {
  hooks.beforeEach(function () {
    this.id = 'foobar';
    this.capabilities = {
      [`test/path/${this.id}`]: {
        canList: false,
        canRead: true,
        canUpdate: true,
        canCreate: false,
        canDelete: false,
      },
    };
  });

  test('it should return true if capabilities are not found in path', async function (assert) {
    assert.true(hasCapability(), 'returns true when capabilities are not provided');
    assert.true(hasCapability(this.capabilities), 'returns true when types are not provided');
    assert.true(hasCapability(this.capabilities, ['read']), 'returns true when id is not provided');
    assert.true(
      hasCapability(this.capabilities, ['read'], 'notAnId'),
      'returns true when id is not found in map'
    );
    assert.true(
      hasCapability(this.capabilities, ['doEverything'], this.path),
      'returns true when type is not found in map'
    );
  });

  test('it should return correct value when evaluating single type', async function (assert) {
    assert.true(
      hasCapability(this.capabilities, ['read'], this.id),
      'returns correct true value for single type'
    );
    assert.false(
      hasCapability(this.capabilities, ['list'], this.id),
      'returns correct false value for single type'
    );
  });

  test('it should return correct value when evaluating multiple types when all is false', async function (assert) {
    assert.true(
      hasCapability(this.capabilities, ['list', 'read'], this.id),
      'returns correct true value for multiple types when at least one is true'
    );
    assert.false(
      hasCapability(this.capabilities, ['create', 'delete'], this.id),
      'returns correct false value for multiple types when all are false'
    );
  });

  test('it should return correct value when evaluating multiple types when all is true', async function (assert) {
    assert.true(
      hasCapability(this.capabilities, ['update', 'read'], this.id, true),
      'returns correct true value for multiple types when all are true'
    );
    assert.false(
      hasCapability(this.capabilities, ['update', 'read', 'list'], this.id, true),
      'returns correct false value for multiple types when at least one is false'
    );
  });
});
