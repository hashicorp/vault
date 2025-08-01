import { Factory, trait } from 'miragejs';
import timestamp from 'core/utils/timestamp';
import { addYears, subMonths } from 'date-fns';

export default Factory.extend({
  // Core health status fields
  initialized: true,
  sealed: false,
  standby: false,
  performance_standby: false,
  replication_performance_mode: 'disabled',
  replication_dr_mode: 'disabled',

  // Server information
  server_time_utc: () => Math.floor(timestamp.now() / 1000),
  version: '1.15.0+ent',
  cluster_name: 'vault-cluster-e779cd7c',
  cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',

  // Storage and WAL
  last_wal: 121,

  // Enterprise features
  enterprise: true,
  license: () => ({
    expiry: addYears(timestamp.now(), 1).toISOString(),
    state: 'autoloaded',
  }),

  // Migration status (when applicable)
  migration: false,

  // Echo request ID for debugging
  echo_duration_ms: null,
  clock_skew_ms: null,

  // Traits for different health states
  isSealed: trait({
    sealed: true,
    initialized: true,
    standby: false,
  }),

  isUninitialized: trait({
    initialized: false,
    sealed: true,
    standby: false,
  }),

  isStandby: trait({
    standby: true,
    performance_standby: false,
    sealed: false,
    initialized: true,
  }),

  isPerformanceStandby: trait({
    standby: false,
    performance_standby: true,
    sealed: false,
    initialized: true,
  }),

  isDrPrimary: trait({
    replication_dr_mode: 'primary',
    replication_performance_mode: 'disabled',
  }),

  isDrSecondary: trait({
    replication_dr_mode: 'secondary',
    replication_performance_mode: 'disabled',
  }),

  isPerformancePrimary: trait({
    replication_performance_mode: 'primary',
    replication_dr_mode: 'disabled',
  }),

  isPerformanceSecondary: trait({
    replication_performance_mode: 'secondary',
    replication_dr_mode: 'disabled',
  }),

  isCommunity: trait({
    enterprise: false,
    version: '1.15.0',
    license: null,
  }),

  isExpired: trait({
    license: () => ({
      expiry: subMonths(timestamp.now(), 1).toISOString(),
      state: 'expired',
    }),
  }),
});
