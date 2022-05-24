// base handlers used in mirage config when a specific handler is not specified
const EXPIRY_DATE = '2021-05-12T23:20:50.52Z';

export default function (server) {
  server.get('/sys/internal/ui/feature-flags', (db) => {
    const featuresResponse = db.features.first();
    return {
      data: {
        feature_flags: featuresResponse ? featuresResponse.feature_flags : null,
      },
    };
  });

  server.get('/sys/health', function () {
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

  server.get('/sys/license/status', function () {
    return {
      data: {
        autoloading_used: false,
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
      },
    };
  });

  server.get('sys/namespaces', function () {
    return {
      data: {
        keys: [
          'ns1/',
          'ns2/',
          'ns3/',
          'ns4/',
          'ns5/',
          'ns6/',
          'ns7/',
          'ns8/',
          'ns9/',
          'ns10/',
          'ns11/',
          'ns12/',
          'ns13/',
          'ns14/',
          'ns15/',
          'ns16/',
          'ns17/',
          'ns18/',
        ],
      },
    };
  });
}
