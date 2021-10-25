const EXPIRY_DATE = '2021-05-12T23:20:50.52Z';

export default function() {
  this.namespace = 'v1';

  this.get('sys/internal/counters/activity', function(db) {
    let data = {};
    const firstRecord = db['clients/activities'].first();
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
      data: db['clients/configs'].first(),
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

  this.get('/sys/internal/counters/activity/monthly', function() {
    return {
      data: {
        by_namespace: [
          {
            namespace_id: 'rdOZR',
            namespace_path: 'namespace1/',
            counts: {
              distinct_entities: 971,
              non_entity_tokens: 940,
              clients: 1911,
            },
          },
          {
            namespace_id: 'Y5daB',
            namespace_path: 'namespace3/',
            counts: {
              distinct_entities: 975,
              non_entity_tokens: 606,
              clients: 1581,
            },
          },
          {
            namespace_id: '0PY4u',
            namespace_path: 'namespace27/',
            counts: {
              distinct_entities: 839,
              non_entity_tokens: 685,
              clients: 1524,
            },
          },
          {
            namespace_id: 'aqdXs',
            namespace_path: 'namespace226/',
            counts: {
              distinct_entities: 578,
              non_entity_tokens: 838,
              clients: 1416,
            },
          },
          {
            namespace_id: 'k6SUB',
            namespace_path: 'namespace10/',
            counts: {
              distinct_entities: 427,
              non_entity_tokens: 952,
              clients: 1379,
            },
          },
          {
            namespace_id: 'AWujC',
            namespace_path: 'namespace16/',
            counts: {
              distinct_entities: 835,
              non_entity_tokens: 534,
              clients: 1369,
            },
          },
          {
            namespace_id: 'Zt31d',
            namespace_path: 'namespace17/',
            counts: {
              distinct_entities: 696,
              non_entity_tokens: 664,
              clients: 1360,
            },
          },
          {
            namespace_id: 's6tBw',
            namespace_path: 'namespace212/',
            counts: {
              distinct_entities: 557,
              non_entity_tokens: 758,
              clients: 1315,
            },
          },
          {
            namespace_id: 'cj2Bi',
            namespace_path: 'namespacelonglonglong7/',
            counts: {
              distinct_entities: 429,
              non_entity_tokens: 834,
              clients: 1263,
            },
          },
          {
            namespace_id: 'WV0oI',
            namespace_path: 'namespace2/',
            counts: {
              distinct_entities: 417,
              non_entity_tokens: 837,
              clients: 1254,
            },
          },
          {
            namespace_id: '6iD69',
            namespace_path: 'namespace19/',
            counts: {
              distinct_entities: 610,
              non_entity_tokens: 559,
              clients: 1169,
            },
          },
          {
            namespace_id: 'dxh3E',
            namespace_path: 'namespace22/',
            counts: {
              distinct_entities: 537,
              non_entity_tokens: 601,
              clients: 1138,
            },
          },
          {
            namespace_id: 'fDPrT',
            namespace_path: 'namespace26/',
            counts: {
              distinct_entities: 605,
              non_entity_tokens: 527,
              clients: 1132,
            },
          },
          {
            namespace_id: 'QYquU',
            namespace_path: 'namespace2-/',
            counts: {
              distinct_entities: 584,
              non_entity_tokens: 490,
              clients: 1074,
            },
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 252,
              non_entity_tokens: 821,
              clients: 1073,
            },
          },
          {
            namespace_id: 'cNuVQ',
            namespace_path: 'namespace11/',
            counts: {
              distinct_entities: 291,
              non_entity_tokens: 762,
              clients: 1053,
            },
          },
          {
            namespace_id: 'B3POb',
            namespace_path: 'namespace225/',
            counts: {
              distinct_entities: 109,
              non_entity_tokens: 932,
              clients: 1041,
            },
          },
          {
            namespace_id: '1Lam4',
            namespace_path: 'namespace15/',
            counts: {
              distinct_entities: 60,
              non_entity_tokens: 943,
              clients: 1003,
            },
          },
          {
            namespace_id: 'QvkqC',
            namespace_path: 'namespace6/',
            counts: {
              distinct_entities: 896,
              non_entity_tokens: 35,
              clients: 931,
            },
          },
          {
            namespace_id: 'wGvyy',
            namespace_path: 'namespace14/',
            counts: {
              distinct_entities: 108,
              non_entity_tokens: 718,
              clients: 826,
            },
          },
          {
            namespace_id: 'OoX8p',
            namespace_path: 'namespace9/',
            counts: {
              distinct_entities: 290,
              non_entity_tokens: 514,
              clients: 804,
            },
          },
          {
            namespace_id: 'zDzCb',
            namespace_path: 'namespacelonglonglong8/',
            counts: {
              distinct_entities: 774,
              non_entity_tokens: 16,
              clients: 790,
            },
          },
          {
            namespace_id: 'NbPw5',
            namespace_path: 'namespace25/',
            counts: {
              distinct_entities: 433,
              non_entity_tokens: 333,
              clients: 766,
            },
          },
          {
            namespace_id: 'UxWkm',
            namespace_path: 'namespace21/',
            counts: {
              distinct_entities: 614,
              non_entity_tokens: 119,
              clients: 733,
            },
          },
          {
            namespace_id: 't9XKF',
            namespace_path: 'namespace13/',
            counts: {
              distinct_entities: 679,
              non_entity_tokens: 28,
              clients: 707,
            },
          },
          {
            namespace_id: 'AjB3O',
            namespace_path: 'namespace23/',
            counts: {
              distinct_entities: 436,
              non_entity_tokens: 239,
              clients: 675,
            },
          },
          {
            namespace_id: 'zYHxi',
            namespace_path: 'namespace24/',
            counts: {
              distinct_entities: 167,
              non_entity_tokens: 428,
              clients: 595,
            },
          },
          {
            namespace_id: 'XwMN3',
            namespace_path: 'namespace28/',
            counts: {
              distinct_entities: 389,
              non_entity_tokens: 197,
              clients: 586,
            },
          },
          {
            namespace_id: 'fQbwO',
            namespace_path: 'namespace5/',
            counts: {
              distinct_entities: 331,
              non_entity_tokens: 254,
              clients: 585,
            },
          },
          {
            namespace_id: '6e78R',
            namespace_path: 'namespace29/',
            counts: {
              distinct_entities: 92,
              non_entity_tokens: 458,
              clients: 550,
            },
          },
          {
            namespace_id: 'eE0N5',
            namespace_path: 'namespace4/',
            counts: {
              distinct_entities: 32,
              non_entity_tokens: 391,
              clients: 423,
            },
          },
          {
            namespace_id: 'HzfYV',
            namespace_path: 'namespace220/',
            counts: {
              distinct_entities: 277,
              non_entity_tokens: 135,
              clients: 412,
            },
          },
          {
            namespace_id: '9Ykg6',
            namespace_path: 'namespace12/',
            counts: {
              distinct_entities: 201,
              non_entity_tokens: 57,
              clients: 258,
            },
          },
          {
            namespace_id: 'hVPni',
            namespace_path: 'namespace118/',
            counts: {
              distinct_entities: 222,
              non_entity_tokens: 9,
              clients: 231,
            },
          },
          {
            namespace_id: 'DkpBZ',
            namespace_path: 'namespace30/',
            counts: {
              distinct_entities: 42,
              non_entity_tokens: 181,
              clients: 223,
            },
          },
        ],
        end_time: '2021-09-30T23:59:59Z',
        start_time: '2021-09-01T00:00:00Z',
        total: { distinct_entities: 15755, non_entity_tokens: 17395, clients: 33150 },
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

  this.get('/sys/license/status', function() {
    return {
      data: {
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
      },
    };
  });

  this.passthrough();
}
