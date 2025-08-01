import { Factory, trait } from 'miragejs';
import timestamp from 'core/utils/timestamp';

export default Factory.extend({
  // Core replication status
  mode: 'disabled',
  dr: () => ({
    mode: 'disabled',
  }),
  performance: () => ({
    mode: 'disabled',
  }),

  // Traits for different replication states
  drPrimary: trait({
    mode: 'primary',
    dr: () => ({
      mode: 'primary',
      cluster_id: 'dr-primary-cluster-id',
      known_secondaries: ['dr-secondary-1', 'dr-secondary-2'],
      last_wal: 1234,
      last_remote_wal: 1230,
      merkle_root: 'merkle-root-hash',
      state: 'running',
    }),
    performance: () => ({ mode: 'disabled' }),
  }),

  drSecondary: trait({
    mode: 'secondary',
    dr: () => ({
      mode: 'secondary',
      cluster_id: 'dr-secondary-cluster-id',
      known_primary_cluster_addr: 'https://primary.vault.com:8201',
      primary_cluster_addr: 'https://primary.vault.com:8201',
      last_wal: 1234,
      last_remote_wal: 1235,
      merkle_root: 'merkle-root-hash',
      state: 'stream-wals',
      connection_state: 'ready',
      last_heartbeat: () => timestamp.now(),
    }),
    performance: () => ({ mode: 'disabled' }),
  }),

  performancePrimary: trait({
    mode: 'primary',
    dr: () => ({ mode: 'disabled' }),
    performance: () => ({
      mode: 'primary',
      cluster_id: 'perf-primary-cluster-id',
      known_secondaries: ['perf-secondary-1', 'perf-secondary-2'],
      last_wal: 1234,
      last_remote_wal: 1230,
      merkle_root: 'merkle-root-hash',
      state: 'running',
    }),
  }),

  performanceSecondary: trait({
    mode: 'secondary',
    dr: () => ({ mode: 'disabled' }),
    performance: () => ({
      mode: 'secondary',
      cluster_id: 'perf-secondary-cluster-id',
      known_primary_cluster_addr: 'https://primary.vault.com:8201',
      primary_cluster_addr: 'https://primary.vault.com:8201',
      last_wal: 1234,
      last_remote_wal: 1235,
      merkle_root: 'merkle-root-hash',
      state: 'stream-wals',
      connection_state: 'ready',
      last_heartbeat: () => timestamp.now(),
    }),
  }),

  bothPrimary: trait({
    mode: 'primary',
    dr: () => ({
      mode: 'primary',
      cluster_id: 'dr-primary-cluster-id',
      known_secondaries: ['dr-secondary-1'],
      last_wal: 1234,
      last_remote_wal: 1230,
      merkle_root: 'merkle-root-hash',
      state: 'running',
    }),
    performance: () => ({
      mode: 'primary',
      cluster_id: 'perf-primary-cluster-id',
      known_secondaries: ['perf-secondary-1'],
      last_wal: 1234,
      last_remote_wal: 1230,
      merkle_root: 'merkle-root-hash',
      state: 'running',
    }),
  }),

  bothSecondary: trait({
    mode: 'secondary',
    dr: () => ({
      mode: 'secondary',
      cluster_id: 'dr-secondary-cluster-id',
      known_primary_cluster_addr: 'https://primary.vault.com:8201',
      primary_cluster_addr: 'https://primary.vault.com:8201',
      last_wal: 1234,
      last_remote_wal: 1235,
      merkle_root: 'merkle-root-hash',
      state: 'stream-wals',
      connection_state: 'ready',
      last_heartbeat: () => timestamp.now(),
    }),
    performance: () => ({
      mode: 'secondary',
      cluster_id: 'perf-secondary-cluster-id',
      known_primary_cluster_addr: 'https://primary.vault.com:8201',
      primary_cluster_addr: 'https://primary.vault.com:8201',
      last_wal: 1234,
      last_remote_wal: 1235,
      merkle_root: 'merkle-root-hash',
      state: 'stream-wals',
      connection_state: 'ready',
      last_heartbeat: () => timestamp.now(),
    }),
  }),

  withIssues: trait({
    mode: 'secondary',
    dr: () => ({
      mode: 'secondary',
      cluster_id: 'dr-secondary-cluster-id',
      known_primary_cluster_addr: 'https://primary.vault.com:8201',
      primary_cluster_addr: 'https://primary.vault.com:8201',
      last_wal: 1200,
      last_remote_wal: 1235,
      merkle_root: 'merkle-root-hash',
      state: 'idle',
      connection_state: 'transient_failure',
      last_heartbeat: () => timestamp.now() - 300000, // 5 minutes ago
    }),
    performance: () => ({ mode: 'disabled' }),
  }),
});
