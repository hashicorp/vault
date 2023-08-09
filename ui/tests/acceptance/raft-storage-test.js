/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | raft storage', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.config = this.server.create('configuration', 'withRaft');
    this.server.get('/sys/internal/ui/resultant-acl', () =>
      this.server.create('configuration', { data: { root: true } })
    );
    this.server.get('/sys/license/features', () => ({}));
    await authPage.login();
  });

  test('it should render correct number of raft peers', async function (assert) {
    assert.expect(3);

    let didRemovePeer = false;
    this.server.get('/sys/storage/raft/configuration', () => {
      if (didRemovePeer) {
        this.config.data.config.servers.pop();
      } else {
        // consider peer removed by external means (cli) after initial request
        didRemovePeer = true;
      }
      return this.config;
    });

    await visit('/vault/storage/raft');
    assert.dom('[data-raft-row]').exists({ count: 2 }, '2 raft peers render in table');
    // leave route and return to trigger config fetch
    await visit('/vault/secrets');
    await visit('/vault/storage/raft');
    const store = this.owner.lookup('service:store');
    assert.strictEqual(
      store.peekAll('server').length,
      2,
      'Store contains 2 server records since remove peer was triggered externally'
    );
    assert.dom('[data-raft-row]').exists({ count: 1 }, 'Only raft nodes from response are rendered');
  });

  test('it should remove raft peer', async function (assert) {
    assert.expect(3);

    this.server.get('/sys/storage/raft/configuration', () => this.config);
    this.server.post('/sys/storage/raft/remove-peer', (schema, req) => {
      const body = JSON.parse(req.requestBody);
      assert.strictEqual(
        body.server_id,
        this.config.data.config.servers[1].node_id,
        'Remove peer request made with node id'
      );
      return {};
    });

    await visit('/vault/storage/raft');
    assert.dom('[data-raft-row]').exists({ count: 2 }, '2 raft peers render in table');
    await click('[data-raft-row]:nth-child(2) [data-test-popup-menu-trigger]');
    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button]');
    assert.dom('[data-raft-row]').exists({ count: 1 }, 'Raft peer successfully removed');
  });
});
