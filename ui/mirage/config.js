const EXPIRY_DATE = '2021-05-12T23:20:50.52Z';

export default function() {
  this.namespace = 'v1';

  this.get('sys/internal/counters/activity', function(db) {
    let data = {};
    const firstRecord = db['metrics/activities'].first();
    if (firstRecord) {
      data = firstRecord;
    }
    return {
      data,
      request_id: '0001',
    };
  });

  this.get('sys/internal/counters/config', function(db) {
    return {
      request_id: '00001',
      data: db['metrics/configs'].first(),
    };
  });

  this.get('/sys/internal/ui/feature-flags', db => {
    const featuresResponse = db.features.first();
    return {
      data: {
        feature_flags: featuresResponse ? featuresResponse.feature_flags : null,
      },
    };
  });

  this.get('/sys/license/status', function() {
    return {
      autoloading_used: false,
      stored: {
        expiration_time: EXPIRY_DATE,
        features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
        license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
        performance_standby_count: 0,
        start_time: '2020-04-28T00:00:00Z',
      },
      persisted_autoload: {
        expiration_time: EXPIRY_DATE,
        features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
        license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
        performance_standby_count: 0,
        start_time: '2020-04-28T00:00:00Z',
      },
      autoloaded: {
        expiration_time: EXPIRY_DATE,
        features: ['DR Replication', 'Namespaces', 'Lease Count Quotas', 'Automated Snapshots'],
        license_id: '0eca7ef8-ebc0-f875-315e-3cc94a7870cf',
        performance_standby_count: 0,
        start_time: '2020-04-28T00:00:00Z',
      },
    };
  });

  this.get('/sys/health', function() {
    return {
      initialized: true,
      sealed: false,
      standby: false,
      license: {
        expiry: '2021-05-12T23:20:50.52Z',
        state: 'stored',
      },
      performance_standby: false,
      replication_performance_mode: 'disabled',
      replication_dr_mode: 'disabled',
      server_time_utc: 1622562585,
      version: '1.9.0+ent',
      cluster_name: 'vault-cluster-e779cd7c',
      cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
      last_wal: 121,
    };
  });

  this.passthrough();
}
