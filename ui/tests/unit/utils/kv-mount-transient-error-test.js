/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { isTransientKvMountError } from 'vault/utils/kv-mount-transient-error';

module('Unit | Utility | kv-mount-transient-error', function () {
  test('it matches the perf standby/secondary KV upgrade 400 message', function (assert) {
    assert.true(
      isTransientKvMountError(
        400,
        'Waiting for the primary to upgrade from non-versioned to versioned data. This backend will be unavailable for a brief period and will resume service when the primary is finished.'
      )
    );
  });

  test('it matches the non-replicated KV upgrade 400 message', function (assert) {
    assert.true(
      isTransientKvMountError(
        400,
        'Upgrading from non-versioned to versioned data. This backend will be unavailable for a brief period and will resume service shortly.'
      )
    );
  });

  test('it matches the mount-not-yet-routable 404 message', function (assert) {
    assert.true(
      isTransientKvMountError(404, 'no handler for route "my-kv/metadata/". route entry not found.')
    );
  });

  test('it does not match an unrelated 400', function (assert) {
    assert.false(isTransientKvMountError(400, 'permission denied'));
  });

  test('it does not match an unrelated 404, such as the plugin pin lookup', function (assert) {
    assert.false(isTransientKvMountError(404, 'no pinned version for this plugin'));
  });

  test('it does not match a 403', function (assert) {
    assert.false(
      isTransientKvMountError(403, 'Waiting for the primary to upgrade from non-versioned to versioned data.')
    );
  });

  test('it returns false when message is missing', function (assert) {
    assert.false(isTransientKvMountError(400, undefined));
    assert.false(isTransientKvMountError(404, ''));
  });
});
