/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | raft storage', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.config = this.server.create('configuration', 'withRaft');
    this.server.get('/sys/internal/ui/resultant-acl', () =>
      this.server.create('configuration', { data: { root: true } })
    );
    this.server.get('/sys/license/features', () => ({ features: [] }));
    await login();
  });

  test('it should render correct number of raft peers', async function (assert) {
    assert.expect(2);

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
    // leave route and return to trigger a fresh fetch of the raft configuration
    await visit('/vault/secrets-engines');
    await visit('/vault/storage/raft');
    assert
      .dom('[data-raft-row]')
      .exists({ count: 1 }, 'Only raft nodes from the latest response are rendered');
  });

  test('it should remove raft peer', async function (assert) {
    assert.expect(3);

    this.server.get('/sys/storage/raft/configuration', () => this.config);
    this.server.post('/sys/storage/raft/remove-peer', (schema, req) => {
      const body = JSON.parse(req.requestBody);
      const removedNodeId = this.config.data.config.servers[1].node_id;
      assert.strictEqual(body.server_id, removedNodeId, 'Remove peer request made with node id');
      // simulate the server actually removing the peer so the follow-up route refresh reflects it
      this.config.data.config.servers = this.config.data.config.servers.filter(
        (server) => server.node_id !== body.server_id
      );
      return {};
    });

    const row = '[data-raft-row]:nth-child(2) [data-test-raft-actions]';
    await visit('/vault/storage/raft');
    assert.dom('[data-raft-row]').exists({ count: 2 }, '2 raft peers render in table');
    await click(`${row} button`);
    await click(`${row} ${GENERAL.confirmTrigger}`);
    await click(GENERAL.confirmButton);
    assert.dom('[data-raft-row]').exists({ count: 1 }, 'Raft peer successfully removed');
  });
});
