/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { destinationTypes } from 'vault/helpers/sync-destinations';

module('Unit | Serializer | sync/destination', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.serializer = this.owner.lookup('serializer:sync/destination');
  });

  test('it normalizes the payload from the server', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);

    for (const destinationType of types) {
      const { name, type, id, ...connection_details } = this.server.create(
        'sync-destination',
        destinationType
      );
      const serverData = { request_id: id, data: { name, type, connection_details } };
      const normalized = this.serializer.normalizeItems(serverData);
      const expected = { type, name, ...connection_details };
      assert.propEqual(
        normalized,
        expected,
        `connection details key is removed and params are added to ${destinationType} data object`
      );
    }
  });
});
