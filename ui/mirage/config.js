const EXPIRY_DATE = '2021-05-12T23:20:50.52Z';

export default function() {
  this.namespace = 'v1';

  // this.get('sys/internal/counters/activity', function(db) {
  //   let data = {};
  //   const firstRecord = db['clients/activities'].first();
  //   if (firstRecord) {
  //     data = firstRecord;
  //   }
  //   return {
  //     data,
  //     request_id: '0001',
  //   };
  // });

  // this.get('sys/internal/counters/config', function(db) {
  //   return {
  //     request_id: '00001',
  //     data: db['clients/configs'].first(),
  //   };
  // });

  // this.get('/sys/internal/ui/feature-flags', db => {
  //   const featuresResponse = db.features.first();
  //   return {
  //     data: {
  //       feature_flags: featuresResponse ? featuresResponse.feature_flags : null,
  //     },
  //   };
  // });

  this.get('/sys/internal/counters/activity/monthly', function() {
    return {
      data: {
        by_namespace: [
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 832,
              non_entity_tokens: 990,
              clients: 1822,
            },
          },
          {
            namespace_id: 'eE0N5',
            namespace_path: 'namespace4/',
            counts: {
              distinct_entities: 948,
              non_entity_tokens: 750,
              clients: 1698,
            },
          },
          {
            namespace_id: 'hVPni',
            namespace_path: 'namespace118/',
            counts: {
              distinct_entities: 619,
              non_entity_tokens: 808,
              clients: 1427,
            },
          },
          {
            namespace_id: 'AWujC',
            namespace_path: 'namespace16/',
            counts: {
              distinct_entities: 492,
              non_entity_tokens: 903,
              clients: 1395,
            },
          },
          {
            namespace_id: 'NbPw5',
            namespace_path: 'namespace25/',
            counts: {
              distinct_entities: 957,
              non_entity_tokens: 423,
              clients: 1380,
            },
          },
          {
            namespace_id: '0PY4u',
            namespace_path: 'namespace27/',
            counts: {
              distinct_entities: 559,
              non_entity_tokens: 805,
              clients: 1364,
            },
          },
          {
            namespace_id: 'rdOZR',
            namespace_path: 'namespace1/',
            counts: {
              distinct_entities: 365,
              non_entity_tokens: 931,
              clients: 1296,
            },
          },
          {
            namespace_id: '6iD69',
            namespace_path: 'namespace19/',
            counts: {
              distinct_entities: 340,
              non_entity_tokens: 936,
              clients: 1276,
            },
          },
          {
            namespace_id: 'wGvyy',
            namespace_path: 'namespace14/',
            counts: {
              distinct_entities: 383,
              non_entity_tokens: 891,
              clients: 1274,
            },
          },
          {
            namespace_id: 'cNuVQ',
            namespace_path: 'namespace11/',
            counts: {
              distinct_entities: 460,
              non_entity_tokens: 789,
              clients: 1249,
            },
          },
          {
            namespace_id: 't9XKF',
            namespace_path: 'namespace13/',
            counts: {
              distinct_entities: 922,
              non_entity_tokens: 318,
              clients: 1240,
            },
          },
          {
            namespace_id: 'QvkqC',
            namespace_path: 'namespace6/',
            counts: {
              distinct_entities: 558,
              non_entity_tokens: 645,
              clients: 1203,
            },
          },
          {
            namespace_id: 'HzfYV',
            namespace_path: 'namespace220/',
            counts: {
              distinct_entities: 735,
              non_entity_tokens: 376,
              clients: 1111,
            },
          },
          {
            namespace_id: 'XwMN3',
            namespace_path: 'namespace28/',
            counts: {
              distinct_entities: 452,
              non_entity_tokens: 658,
              clients: 1110,
            },
          },
          {
            namespace_id: 's6tBw',
            namespace_path: 'namespace212/',
            counts: {
              distinct_entities: 602,
              non_entity_tokens: 483,
              clients: 1085,
            },
          },
          {
            namespace_id: 'QYquU',
            namespace_path: 'namespace2-/',
            counts: {
              distinct_entities: 304,
              non_entity_tokens: 716,
              clients: 1020,
            },
          },
          {
            namespace_id: 'fQbwO',
            namespace_path: 'namespace5/',
            counts: {
              distinct_entities: 44,
              non_entity_tokens: 942,
              clients: 986,
            },
          },
          {
            namespace_id: 'cj2Bi',
            namespace_path: 'namespacelonglonglong7/',
            counts: {
              distinct_entities: 538,
              non_entity_tokens: 389,
              clients: 927,
            },
          },
          {
            namespace_id: 'Y5daB',
            namespace_path: 'namespace3/',
            counts: {
              distinct_entities: 692,
              non_entity_tokens: 211,
              clients: 903,
            },
          },
          {
            namespace_id: 'aqdXs',
            namespace_path: 'namespace226/',
            counts: {
              distinct_entities: 788,
              non_entity_tokens: 109,
              clients: 897,
            },
          },
          {
            namespace_id: 'DkpBZ',
            namespace_path: 'namespace30/',
            counts: {
              distinct_entities: 379,
              non_entity_tokens: 499,
              clients: 878,
            },
          },
          {
            namespace_id: '9Ykg6',
            namespace_path: 'namespace12/',
            counts: {
              distinct_entities: 167,
              non_entity_tokens: 613,
              clients: 780,
            },
          },
          {
            namespace_id: 'zDzCb',
            namespace_path: 'namespacelonglonglong8/',
            counts: {
              distinct_entities: 453,
              non_entity_tokens: 299,
              clients: 752,
            },
          },
          {
            namespace_id: 'dxh3E',
            namespace_path: 'namespace22/',
            counts: {
              distinct_entities: 158,
              non_entity_tokens: 512,
              clients: 670,
            },
          },
          {
            namespace_id: 'fDPrT',
            namespace_path: 'namespace26/',
            counts: {
              distinct_entities: 9,
              non_entity_tokens: 656,
              clients: 665,
            },
          },
          {
            namespace_id: 'k6SUB',
            namespace_path: 'namespace10/',
            counts: {
              distinct_entities: 278,
              non_entity_tokens: 369,
              clients: 647,
            },
          },
          {
            namespace_id: 'UxWkm',
            namespace_path: 'namespace21/',
            counts: {
              distinct_entities: 588,
              non_entity_tokens: 7,
              clients: 595,
            },
          },
          {
            namespace_id: 'OoX8p',
            namespace_path: 'namespace9/',
            counts: {
              distinct_entities: 519,
              non_entity_tokens: 67,
              clients: 586,
            },
          },
          {
            namespace_id: 'WV0oI',
            namespace_path: 'namespace2/',
            counts: {
              distinct_entities: 302,
              non_entity_tokens: 196,
              clients: 498,
            },
          },
          {
            namespace_id: '6e78R',
            namespace_path: 'namespace29/',
            counts: {
              distinct_entities: 43,
              non_entity_tokens: 455,
              clients: 498,
            },
          },
          {
            namespace_id: 'AjB3O',
            namespace_path: 'namespace23/',
            counts: {
              distinct_entities: 234,
              non_entity_tokens: 206,
              clients: 440,
            },
          },
          {
            namespace_id: 'B3POb',
            namespace_path: 'namespace225/',
            counts: {
              distinct_entities: 300,
              non_entity_tokens: 114,
              clients: 414,
            },
          },
          {
            namespace_id: 'Zt31d',
            namespace_path: 'namespace17/',
            counts: {
              distinct_entities: 218,
              non_entity_tokens: 51,
              clients: 269,
            },
          },
          {
            namespace_id: 'zYHxi',
            namespace_path: 'namespace24/',
            counts: {
              distinct_entities: 240,
              non_entity_tokens: 4,
              clients: 244,
            },
          },
          {
            namespace_id: '1Lam4',
            namespace_path: 'namespace15/',
            counts: {
              distinct_entities: 21,
              non_entity_tokens: 161,
              clients: 182,
            },
          },
        ],
        distinct_entities: 15499,
        non_entity_tokens: 17282,
        clients: 32781,
      },
    };
  });

  // this.get('/sys/health', function() {
  //   return {
  //     initialized: true,
  //     sealed: false,
  //     standby: false,
  //     license: {
  //       expiry: '2021-05-12T23:20:50.52Z',
  //       state: 'stored',
  //     },
  //     performance_standby: false,
  //     replication_performance_mode: 'disabled',
  //     replication_dr_mode: 'disabled',
  //     server_time_utc: 1622562585,
  //     version: '1.9.0+ent',
  //     cluster_name: 'vault-cluster-e779cd7c',
  //     cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
  //     last_wal: 121,
  //   };
  // });

  this.passthrough();
}
