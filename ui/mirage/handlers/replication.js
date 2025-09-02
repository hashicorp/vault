/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.get('sys/replication/status', function () {
    return {
      data: {
        dr: {
          cluster_id: 'dr-cluster-id',
          corrupted_merkle_tree: false,
          known_secondaries: ['foobar'],
          last_corruption_check_epoch: '-62135596800',
          last_dr_wal: 98,
          last_reindex_epoch: '0',
          last_wal: 98,
          merkle_root: 'ad721f32ed6789a1e5824841f358a517340a4585',
          mode: 'primary',
          primary_cluster_addr: 'dr-foobar',
          secondaries: [
            {
              api_address: 'http://127.0.0.1:8202',
              cluster_address: 'https://127.0.0.1:8203',
              connection_status: 'disconnected',
              last_heartbeat: '2024-04-09T09:04:22-05:00',
              node_id: 'foobar',
            },
          ],
          ssct_generation_counter: 0,
          state: 'running',
        },
        performance: {
          cluster_id: 'perf-cluster-id',
          corrupted_merkle_tree: false,
          known_secondaries: [],
          last_corruption_check_epoch: '-62135596800',
          last_performance_wal: 98,
          last_reindex_epoch: '0',
          last_wal: 98,
          merkle_root: '618c9136bb443aa584f5d6b90755d42888c9c54a',
          mode: 'primary',
          primary_cluster_addr: 'perf-foobar',
          secondaries: [],
          ssct_generation_counter: 0,
          state: 'running',
        },
      },
    };
  });

  server.get('sys/replication/performance/status', function () {
    return {
      data: {
        cluster_id: 'perf-cluster-id',
        corrupted_merkle_tree: false,
        known_secondaries: ['foobar'],
        last_corruption_check_epoch: '-62135596800',
        last_dr_wal: 98,
        last_reindex_epoch: '0',
        last_wal: 98,
        merkle_root: 'ad721f32ed6789a1e5824841f358a517340a4585',
        mode: 'primary',
        primary_cluster_addr: 'perf-foobar',
        secondaries: [
          {
            api_address: 'http://127.0.0.1:8202',
            cluster_address: 'https://127.0.0.1:8203',
            connection_status: 'disconnected',
            last_heartbeat: '2024-04-09T09:04:22-05:00',
            node_id: 'foobar',
          },
        ],
        ssct_generation_counter: 0,
        state: 'running',
      },
    };
  });

  server.get('sys/replication/dr/status', function () {
    return {
      data: {
        cluster_id: 'dr-cluster-id',
        corrupted_merkle_tree: false,
        known_secondaries: ['foobar'],
        last_corruption_check_epoch: '-62135596800',
        last_dr_wal: 98,
        last_reindex_epoch: '0',
        last_wal: 98,
        merkle_root: 'ad721f32ed6789a1e5824841f358a517340a4585',
        mode: 'primary',
        primary_cluster_addr: 'dr-foobar',
        secondaries: [
          {
            api_address: 'http://127.0.0.1:8202',
            cluster_address: 'https://127.0.0.1:8203',
            connection_status: 'disconnected',
            last_heartbeat: '2024-04-09T09:04:22-05:00',
            node_id: 'foobar',
          },
        ],
        ssct_generation_counter: 0,
        state: 'running',
      },
    };
  });
}
