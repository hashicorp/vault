/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import HasCapabilityHelper from 'core/helpers/has-capability';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Helpers | has-capability', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.helper = new HasCapabilityHelper(this.owner);
    this.namedArgs = { pathKey: 'customMessages', params: { id: 'foobar' } };
    const capabilities = {
      'sys/config/ui/custom-messages/foobar': {
        canList: false,
        canRead: true,
        canUpdate: true,
        canCreate: false,
        canDelete: false,
      },
    };
    this.hasCapability = (types = [], namedArgs) => {
      return this.helper.compute([capabilities, ...types], namedArgs || this.namedArgs);
    };
  });

  test('it should throw an error if capabilities are not provided', async function (assert) {
    try {
      this.helper.compute([], {})();
    } catch (error) {
      assert.strictEqual(
        error.message,
        'First positional argument must be the capabilities map.',
        'throws error when capabilities map is not provided'
      );
    }
  });

  test('it should throw an error when types are not provided', async function (assert) {
    try {
      this.hasCapability([]);
    } catch (error) {
      assert.strictEqual(
        error.message,
        'At least one capability type is required as a positional argument.',
        'throws error when types are not provided'
      );
    }
  });

  test('it should throw an error for invalid types', async function (assert) {
    try {
      this.hasCapability(['doEverything']);
    } catch (error) {
      assert.strictEqual(
        error.message,
        'Invalid capability types: doEverything. Accepted types are: read, update, delete, list, create, patch, sudo.',
        'throws error when invalid types are provided'
      );
    }
  });

  test('it should throw an error if pathKey is not provided', async function (assert) {
    try {
      this.hasCapability(['read'], {});
    } catch (error) {
      assert.strictEqual(
        error.message,
        'pathKey is a required named arg for path lookup in capabilities map',
        'throws error when pathKey is not found'
      );
    }
  });

  test('it should throw error if path is not found', async function (assert) {
    try {
      this.hasCapability(['read'], { pathKey: 'notAnId' });
    } catch (error) {
      assert.strictEqual(
        error.message,
        'Path not found for key: notAnId',
        'throws error when path is not found'
      );
    }
  });

  test('it should return correct value when evaluating single type', async function (assert) {
    assert.true(this.hasCapability(['read']), 'returns correct true value for single type');
    assert.false(this.hasCapability(['list']), 'returns correct false value for single type');
  });

  test('it should return correct value when evaluating multiple types when all is false', async function (assert) {
    assert.true(
      this.hasCapability(['list', 'read']),
      'returns correct true value for multiple types when at least one is true'
    );
    assert.false(
      this.hasCapability(['create', 'delete']),
      'returns correct false value for multiple types when all are false'
    );
  });

  test('it should return correct value when evaluating multiple types when all is true', async function (assert) {
    this.namedArgs.all = true;
    assert.true(
      this.hasCapability(['update', 'read']),
      'returns correct true value for multiple types when all are true'
    );
    assert.false(
      this.hasCapability(['update', 'read', 'list']),
      'returns correct false value for multiple types when at least one is false'
    );
  });
});
