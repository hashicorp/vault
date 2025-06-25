/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, findAll, currentURL, visit, settled, waitUntil } from '@ember/test-helpers';
import ss from 'vault/tests/pages/components/search-select';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const searchSelect = create(ss);
/**
 * Enables replication mode.
 * @param {string} type - The replication type ('performance' or 'dr').
 * @param {string} mode - The cluster mode ('primary' or 'secondary').
 */
export async function enableReplication(type, mode) {
  await visit('/vault/replication');
  await click(`[data-test-replication-type-select="${type}"]`);
  await fillIn('[data-test-replication-cluster-mode-select]', mode);
  await click(GENERAL.submitButton);
}

/**
 * Adds a secondary cluster.
 * @param {string} secondaryName - The name of the secondary cluster.
 * @param {string} mountFilterMode - The mount filter mode ('deny' or 'allow').
 */
export async function addSecondary(secondaryName, mountFilterMode = null) {
  await click('[data-test-replication-link="secondaries"]');
  await click('[data-test-secondary-add]');
  await fillIn('[data-test-input="Secondary ID"]', secondaryName);

  if (mountFilterMode) {
    await click(`#${mountFilterMode}`);
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
  }

  await click('[data-test-secondary-add]');
}

export const disableReplication = async (type, assert) => {
  // disable performance replication
  await visit(`/vault/replication/${type}`);

  if (findAll('[data-test-replication-link="manage"]').length) {
    await click('[data-test-replication-link="manage"]');

    await click('[data-test-disable-replication] button');

    const typeDisplay = type === 'dr' ? 'Disaster Recovery' : 'Performance';
    await fillIn('[data-test-confirmation-modal-input="Disable Replication?"]', typeDisplay);
    await click(GENERAL.confirmButton);
    await settled(); // eslint-disable-line

    if (assert) {
      assert
        .dom(GENERAL.latestFlashContent)
        .hasText(
          'This cluster is having replication disabled. Vault will be unavailable for a brief period and will resume service shortly.'
        );
      assert.ok(
        await waitUntil(() => currentURL() === '/vault/replication'),
        'redirects to the replication page'
      );
    }
    await settled();
  }
};

export const STATUS_DISABLED_RESPONSE = {
  dr: mockReplicationBlock(),
  performance: mockReplicationBlock(),
};

/**
 * Mock replication block returns the expected payload for a given replication type
 * @param {string} mode disabled | primary | secondary
 * @param {string} status connected | disconnected
 * @returns expected object for a single replication type, eg dr or performance values
 */
export function mockReplicationBlock(mode = 'disabled', status = 'connected') {
  switch (mode) {
    case 'primary':
      return {
        cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
        known_secondaries: ['4'],
        last_wal: 455,
        merkle_root: 'aaaaaabbbbbbbccccccccddddddd',
        mode: 'primary',
        primary_cluster_addr: '',
        secondaries: [
          {
            api_address: 'https://127.0.0.1:49277',
            cluster_address: 'https://127.0.0.1:49281',
            connection_status: status,
            last_heartbeat: '2020-06-10T15:40:46-07:00',
            node_id: '4',
          },
        ],
        state: 'stream-wals',
      };
    case 'secondary':
      return {
        cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
        known_primary_cluster_addrs: ['https://127.0.0.1:8201'],
        last_remote_wal: 291,
        merkle_root: 'aaaaaabbbbbbbccccccccddddddd',
        corrupted_merkle_tree: false,
        last_corruption_check_epoch: '1694456090',
        mode: 'secondary',
        primaries: [
          {
            api_address: 'https://127.0.0.1:49244',
            cluster_address: 'https://127.0.0.1:8201',
            connection_status: status,
            last_heartbeat: '2020-06-10T15:40:46-07:00',
          },
        ],
        primary_cluster_addr: 'https://127.0.0.1:8201',
        secondary_id: '2',
        state: 'stream-wals',
      };
    default:
      return {
        mode: 'disabled',
      };
  }
}
