/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint qunit/no-conditional-assertions: "warn" */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';

const SECONDARIES = [
  { node_id: 'secondary-1', api_address: 'https://127.0.0.1:52304', connection_status: 'connected' },
  { node_id: '2nd', connection_status: 'disconnected' },
  { node_id: '_three_', api_address: 'http://127.0.0.1:8202', connection_status: 'connected' },
];

module('Integration | Component | replication known-secondaries-table', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'replication');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    this.set('secondaries', SECONDARIES);
  });

  test('it renders a table of known secondaries', async function (assert) {
    await render(hbs`<KnownSecondariesTable @secondaries={{this.secondaries}} />`, this.context);

    assert.dom('[data-test-known-secondaries-table]').exists();
  });

  test('it shows the secondary URL and connection_status', async function (assert) {
    assert.expect(13);
    await render(hbs`<KnownSecondariesTable @secondaries={{this.secondaries}} />`, this.context);

    SECONDARIES.forEach((secondary) => {
      assert
        .dom(`[data-test-secondaries-node="${secondary.node_id}"]`)
        .hasText(secondary.node_id, 'shows a table row and ID for each known secondary');
      const expectedAPIAddr = secondary.api_address || 'URL unavailable';
      const expectedTag = secondary.api_address ? 'a' : 'p';
      assert.dom(`[data-test-secondaries-api-address="${secondary.node_id}"]`).hasText(expectedAPIAddr);
      assert
        .dom(`[data-test-secondaries-api-address="${secondary.node_id}"] ${expectedTag}`)
        .exists('has correct tag');

      assert
        .dom(`[data-test-secondaries-connection-status="${secondary.node_id}"]`)
        .hasText(secondary.connection_status, 'shows the connection status');
    });

    assert
      .dom(`[data-test-secondaries-api-address="secondary-1"] a`)
      .hasAttribute(
        'href',
        'https://127.0.0.1:52304/ui/',
        'secondary with API address has correct href attribute'
      );
  });
});
