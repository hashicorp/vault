/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | sync/destinations', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it calls the correct endpoint for findRecord', async function (assert) {
    // TODO make destination types a util?
    const destinationTypes = ['aws-sm', 'azure-kv', 'gcp-sm', 'gh', 'vercel-project'];
    assert.expect(destinationTypes.length);

    for (const type of destinationTypes) {
      this.server.get(`sys/sync/destinations/${type}/my-dest`, () => {
        assert.ok(true, `request is made to GET destinations/${type} endpoint on find`);
        return { data: { name: 'my-dest' } };
      });
      this.store.findRecord(`sync/destinations/${type}`, 'my-dest');
    }
  });
});
