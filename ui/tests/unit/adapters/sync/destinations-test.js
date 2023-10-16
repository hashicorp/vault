/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { destinationTypes } from 'vault/helpers/sync-destinations';

module('Unit | Adapter | sync/destinations', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it calls the correct endpoint for findRecord', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);

    for (const type of types) {
      this.server.get(`sys/sync/destinations/${type}/my-dest`, () => {
        assert.ok(true, `request is made to GET sys/sync/destinations/${type}/my-dest endpoint on find`);
        return {
          data: {
            connection_details: {},
            name: 'my-dest',
            type,
          },
        };
      });
      this.store.findRecord(`sync/destinations/${type}`, 'my-dest');
    }
  });
});
