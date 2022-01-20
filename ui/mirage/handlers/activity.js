export default function (server) {
  // 1.10 API response
  server.get('/sys/internal/counters/activity', function () {
    return {
      request_id: '26be5ab9-dcac-9237-ec12-269a8ca647d5',
      data: {
        start_time: '2019-11-01T00:00:00Z',
        end_time: '2020-10-31T23:59:59Z',
        total: {
          _comment1: 'total client counts',
          clients: 3517,
          _comment2: 'following 2 fields are deprecated',
          distinct_entities: 1593,
          non_entity_tokens: 1924,
        },
        by_namespace: [
          {
            namespace_id: '5SWT8',
            namespace_path: 'namespacelonglonglong4/',
            counts: {
              distinct_entities: 171,
              non_entity_tokens: 20,
              clients: 191,
            },
            mounts: [
              {
                path: 'auth/method/uMGBU',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/8YJO3',
                counts: {
                  clients: 28,
                  entity_clients: 18,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/Ro774',
                counts: {
                  clients: 25,
                  entity_clients: 15,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/ZIpjT',
                counts: {
                  clients: 20,
                  entity_clients: 10,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/jdRjF',
                counts: {
                  clients: 27,
                  entity_clients: 15,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/yyBoC',
                counts: {
                  clients: 24,
                  entity_clients: 14,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/WLxYp',
                counts: {
                  clients: 18,
                  entity_clients: 11,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/SNM6V',
                counts: {
                  clients: 5,
                  entity_clients: 2,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/vNHtH',
                counts: {
                  clients: 9,
                  entity_clients: 5,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/EqmlO',
                counts: {
                  clients: 0,
                  entity_clients: 0,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'BCl56',
            namespace_path: 'namespace8/',
            counts: {
              distinct_entities: 141,
              non_entity_tokens: 47,
              clients: 188,
            },
            mounts: [
              {
                path: 'auth/method/LpVqc',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/VFHO6',
                counts: {
                  clients: 33,
                  entity_clients: 19,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/utu0r',
                counts: {
                  clients: 25,
                  entity_clients: 16,
                  non_entity_clients: 9,
                },
              },
              {
                path: 'auth/method/xikiW',
                counts: {
                  clients: 25,
                  entity_clients: 13,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/uPSo6',
                counts: {
                  clients: 18,
                  entity_clients: 12,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/Z8fpo',
                counts: {
                  clients: 14,
                  entity_clients: 7,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/5BBm7',
                counts: {
                  clients: 10,
                  entity_clients: 8,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/Eyxkz',
                counts: {
                  clients: 11,
                  entity_clients: 8,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/QBC0w',
                counts: {
                  clients: 13,
                  entity_clients: 9,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/8MdGr',
                counts: {
                  clients: 4,
                  entity_clients: 3,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'yHcL9',
            namespace_path: 'namespace11/',
            counts: {
              distinct_entities: 10,
              non_entity_tokens: 176,
              clients: 186,
            },
            mounts: [
              {
                path: 'auth/method/qT4Wl',
                counts: {
                  clients: 33,
                  entity_clients: 20,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/Vhu56',
                counts: {
                  clients: 22,
                  entity_clients: 8,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/PCc58',
                counts: {
                  clients: 29,
                  entity_clients: 16,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/nPP4c',
                counts: {
                  clients: 20,
                  entity_clients: 13,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/LY3am',
                counts: {
                  clients: 29,
                  entity_clients: 16,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/McQ4X',
                counts: {
                  clients: 20,
                  entity_clients: 15,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/NpjhH',
                counts: {
                  clients: 17,
                  entity_clients: 11,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/ToKO8',
                counts: {
                  clients: 12,
                  entity_clients: 9,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/wfApH',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/L9uWV',
                counts: {
                  clients: 1,
                  entity_clients: 0,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'bw5UO',
            namespace_path: 'namespace6/',
            counts: {
              distinct_entities: 29,
              non_entity_tokens: 155,
              clients: 184,
            },
            mounts: [
              {
                path: 'auth/method/XQUrA',
                counts: {
                  clients: 34,
                  entity_clients: 19,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/1p6HR',
                counts: {
                  clients: 30,
                  entity_clients: 15,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/qRjoJ',
                counts: {
                  clients: 24,
                  entity_clients: 12,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/x9QQB',
                counts: {
                  clients: 27,
                  entity_clients: 16,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/rezK4',
                counts: {
                  clients: 20,
                  entity_clients: 15,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/qWNSS',
                counts: {
                  clients: 23,
                  entity_clients: 15,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/OmQEf',
                counts: {
                  clients: 14,
                  entity_clients: 10,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/PhoAy',
                counts: {
                  clients: 8,
                  entity_clients: 6,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/aUuyM',
                counts: {
                  clients: 3,
                  entity_clients: 3,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/kUj1S',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'F0xGm',
            namespace_path: 'namespace10/',
            counts: {
              distinct_entities: 75,
              non_entity_tokens: 107,
              clients: 182,
            },
            mounts: [
              {
                path: 'auth/method/xYL0l',
                counts: {
                  clients: 34,
                  entity_clients: 19,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/CwWM7',
                counts: {
                  clients: 25,
                  entity_clients: 13,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/swCd0',
                counts: {
                  clients: 29,
                  entity_clients: 17,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/0CZTs',
                counts: {
                  clients: 21,
                  entity_clients: 8,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/9v04G',
                counts: {
                  clients: 21,
                  entity_clients: 13,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/6hAlO',
                counts: {
                  clients: 23,
                  entity_clients: 12,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/ydSdP',
                counts: {
                  clients: 8,
                  entity_clients: 4,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/i0CTY',
                counts: {
                  clients: 8,
                  entity_clients: 6,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/nevwU',
                counts: {
                  clients: 6,
                  entity_clients: 3,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/k2jYC',
                counts: {
                  clients: 7,
                  entity_clients: 4,
                  non_entity_clients: 3,
                },
              },
            ],
          },
          {
            namespace_id: 'MSgZE',
            namespace_path: 'namespace7/',
            counts: {
              distinct_entities: 72,
              non_entity_tokens: 109,
              clients: 181,
            },
            mounts: [
              {
                path: 'auth/method/gD50V',
                counts: {
                  clients: 31,
                  entity_clients: 19,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/iJRmf',
                counts: {
                  clients: 31,
                  entity_clients: 17,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/GrNjy',
                counts: {
                  clients: 18,
                  entity_clients: 12,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/r0Uw3',
                counts: {
                  clients: 23,
                  entity_clients: 11,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/k2lQG',
                counts: {
                  clients: 25,
                  entity_clients: 18,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/hJxto',
                counts: {
                  clients: 15,
                  entity_clients: 8,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/vtDck',
                counts: {
                  clients: 16,
                  entity_clients: 8,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/1CenH',
                counts: {
                  clients: 9,
                  entity_clients: 4,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/M47Ey',
                counts: {
                  clients: 8,
                  entity_clients: 4,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/gVT0t',
                counts: {
                  clients: 5,
                  entity_clients: 5,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'AAidI',
            namespace_path: 'namespace20/',
            counts: {
              distinct_entities: 39,
              non_entity_tokens: 141,
              clients: 180,
            },
            mounts: [
              {
                path: 'auth/method/zolCO',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/6p3g4',
                counts: {
                  clients: 26,
                  entity_clients: 15,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/iKOdR',
                counts: {
                  clients: 22,
                  entity_clients: 12,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/brnKt',
                counts: {
                  clients: 30,
                  entity_clients: 19,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/qK3rr',
                counts: {
                  clients: 17,
                  entity_clients: 12,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/DmAuN',
                counts: {
                  clients: 13,
                  entity_clients: 7,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/krE4t',
                counts: {
                  clients: 21,
                  entity_clients: 13,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/sFrWK',
                counts: {
                  clients: 11,
                  entity_clients: 10,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/bQg4l',
                counts: {
                  clients: 4,
                  entity_clients: 2,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/Jaw0k',
                counts: {
                  clients: 1,
                  entity_clients: 0,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: '5KxXA',
            namespace_path: 'namespace18anotherlong/',
            counts: {
              distinct_entities: 168,
              non_entity_tokens: 11,
              clients: 179,
            },
            mounts: [
              {
                path: 'auth/method/GkDM1',
                counts: {
                  clients: 33,
                  entity_clients: 18,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/7deLa',
                counts: {
                  clients: 30,
                  entity_clients: 15,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/Ash3Y',
                counts: {
                  clients: 30,
                  entity_clients: 17,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/doKJ0',
                counts: {
                  clients: 21,
                  entity_clients: 13,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/9Irmo',
                counts: {
                  clients: 13,
                  entity_clients: 5,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/jdYx5',
                counts: {
                  clients: 18,
                  entity_clients: 12,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/sYe2h',
                counts: {
                  clients: 11,
                  entity_clients: 6,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/Z5F36',
                counts: {
                  clients: 6,
                  entity_clients: 3,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/O0cuK',
                counts: {
                  clients: 11,
                  entity_clients: 6,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/0clSt',
                counts: {
                  clients: 6,
                  entity_clients: 2,
                  non_entity_clients: 4,
                },
              },
            ],
          },
          {
            namespace_id: 'yYNw2',
            namespace_path: 'namespace19/',
            counts: {
              distinct_entities: 50,
              non_entity_tokens: 129,
              clients: 179,
            },
            mounts: [
              {
                path: 'auth/method/zD8lQ',
                counts: {
                  clients: 31,
                  entity_clients: 16,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/Dl96I',
                counts: {
                  clients: 31,
                  entity_clients: 19,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/ElIse',
                counts: {
                  clients: 31,
                  entity_clients: 19,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/AXzhE',
                counts: {
                  clients: 20,
                  entity_clients: 13,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/cNuC6',
                counts: {
                  clients: 22,
                  entity_clients: 12,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/gXtbE',
                counts: {
                  clients: 20,
                  entity_clients: 13,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/PptIE',
                counts: {
                  clients: 12,
                  entity_clients: 7,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/QILdh',
                counts: {
                  clients: 7,
                  entity_clients: 5,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/cClAS',
                counts: {
                  clients: 2,
                  entity_clients: 2,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/YYm3v',
                counts: {
                  clients: 3,
                  entity_clients: 3,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'R6L40',
            namespace_path: 'namespace2/',
            counts: {
              distinct_entities: 121,
              non_entity_tokens: 56,
              clients: 177,
            },
            mounts: [
              {
                path: 'auth/method/824CE',
                counts: {
                  clients: 33,
                  entity_clients: 18,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/r2zb4',
                counts: {
                  clients: 29,
                  entity_clients: 17,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/1zfD6',
                counts: {
                  clients: 25,
                  entity_clients: 15,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/L14lj',
                counts: {
                  clients: 24,
                  entity_clients: 14,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/cTsw9',
                counts: {
                  clients: 19,
                  entity_clients: 12,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/3KTWZ',
                counts: {
                  clients: 19,
                  entity_clients: 11,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/Douf5',
                counts: {
                  clients: 13,
                  entity_clients: 6,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/30eez',
                counts: {
                  clients: 9,
                  entity_clients: 7,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/xSSJz',
                counts: {
                  clients: 5,
                  entity_clients: 3,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/pR3x7',
                counts: {
                  clients: 1,
                  entity_clients: 0,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'Z4Rzh',
            namespace_path: 'namespace1/',
            counts: {
              distinct_entities: 142,
              non_entity_tokens: 33,
              clients: 175,
            },
            mounts: [
              {
                path: 'auth/method/NqMeC',
                counts: {
                  clients: 34,
                  entity_clients: 20,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/S0FaZ',
                counts: {
                  clients: 30,
                  entity_clients: 20,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/vzH3z',
                counts: {
                  clients: 28,
                  entity_clients: 17,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/uP1zV',
                counts: {
                  clients: 26,
                  entity_clients: 14,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/yAga3',
                counts: {
                  clients: 14,
                  entity_clients: 6,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/DTAFz',
                counts: {
                  clients: 17,
                  entity_clients: 9,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/Rk3Pt',
                counts: {
                  clients: 16,
                  entity_clients: 9,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/wnNH5',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/N3BJy',
                counts: {
                  clients: 4,
                  entity_clients: 2,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/C5qsy',
                counts: {
                  clients: 3,
                  entity_clients: 3,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'XGu7R',
            namespace_path: 'namespace12/',
            counts: {
              distinct_entities: 18,
              non_entity_tokens: 157,
              clients: 175,
            },
            mounts: [
              {
                path: 'auth/method/qcuLl',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/KGWiS',
                counts: {
                  clients: 29,
                  entity_clients: 17,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/iM8pi',
                counts: {
                  clients: 21,
                  entity_clients: 13,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/IeyA4',
                counts: {
                  clients: 27,
                  entity_clients: 18,
                  non_entity_clients: 9,
                },
              },
              {
                path: 'auth/method/KGFfV',
                counts: {
                  clients: 21,
                  entity_clients: 13,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/23AQk',
                counts: {
                  clients: 14,
                  entity_clients: 10,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/PqTWe',
                counts: {
                  clients: 11,
                  entity_clients: 9,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/pPSo1',
                counts: {
                  clients: 8,
                  entity_clients: 6,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/HMu5H',
                counts: {
                  clients: 7,
                  entity_clients: 4,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/xpOk3',
                counts: {
                  clients: 2,
                  entity_clients: 0,
                  non_entity_clients: 2,
                },
              },
            ],
          },
          {
            namespace_id: 'IeyJp',
            namespace_path: 'namespace14/',
            counts: {
              distinct_entities: 33,
              non_entity_tokens: 142,
              clients: 175,
            },
            mounts: [
              {
                path: 'auth/method/8NFVo',
                counts: {
                  clients: 30,
                  entity_clients: 16,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/XnNDy',
                counts: {
                  clients: 31,
                  entity_clients: 18,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/RYrzg',
                counts: {
                  clients: 26,
                  entity_clients: 14,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/SOKji',
                counts: {
                  clients: 27,
                  entity_clients: 17,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/CEYXo',
                counts: {
                  clients: 19,
                  entity_clients: 13,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/RPjsj',
                counts: {
                  clients: 15,
                  entity_clients: 10,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/dIqPJ',
                counts: {
                  clients: 9,
                  entity_clients: 6,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/wThqG',
                counts: {
                  clients: 8,
                  entity_clients: 6,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/Sa1dO',
                counts: {
                  clients: 7,
                  entity_clients: 5,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/0JVs1',
                counts: {
                  clients: 3,
                  entity_clients: 3,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'kxU4t',
            namespace_path: 'namespacelonglonglong3/',
            counts: {
              distinct_entities: 151,
              non_entity_tokens: 21,
              clients: 172,
            },
            mounts: [
              {
                path: 'auth/method/lDz9c',
                counts: {
                  clients: 32,
                  entity_clients: 17,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/GtbUu',
                counts: {
                  clients: 23,
                  entity_clients: 13,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/WCyYz',
                counts: {
                  clients: 30,
                  entity_clients: 17,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/j227p',
                counts: {
                  clients: 21,
                  entity_clients: 14,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/9V6aN',
                counts: {
                  clients: 20,
                  entity_clients: 13,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/USYOd',
                counts: {
                  clients: 17,
                  entity_clients: 13,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/8pfWr',
                counts: {
                  clients: 12,
                  entity_clients: 7,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/0L511',
                counts: {
                  clients: 6,
                  entity_clients: 2,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/6d0rw',
                counts: {
                  clients: 6,
                  entity_clients: 4,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/ECHpZ',
                counts: {
                  clients: 5,
                  entity_clients: 5,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: '5xKya',
            namespace_path: 'namespace15/',
            counts: {
              distinct_entities: 73,
              non_entity_tokens: 98,
              clients: 171,
            },
            mounts: [
              {
                path: 'auth/method/u2r0G',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/mKqBV',
                counts: {
                  clients: 29,
                  entity_clients: 18,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/nGOa2',
                counts: {
                  clients: 19,
                  entity_clients: 9,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/46UKX',
                counts: {
                  clients: 21,
                  entity_clients: 11,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/WHW73',
                counts: {
                  clients: 26,
                  entity_clients: 15,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/KcO46',
                counts: {
                  clients: 20,
                  entity_clients: 12,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/y2vSv',
                counts: {
                  clients: 13,
                  entity_clients: 12,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/VNy4X',
                counts: {
                  clients: 3,
                  entity_clients: 1,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/cEDV9',
                counts: {
                  clients: 2,
                  entity_clients: 2,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/CZTaj',
                counts: {
                  clients: 3,
                  entity_clients: 3,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 112,
              non_entity_tokens: 58,
              clients: 170,
            },
            mounts: [
              {
                path: 'auth/method/koO6h',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/iF9oZ',
                counts: {
                  clients: 31,
                  entity_clients: 20,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/N6guZ',
                counts: {
                  clients: 28,
                  entity_clients: 16,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/h2CxN',
                counts: {
                  clients: 16,
                  entity_clients: 9,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/pA5pU',
                counts: {
                  clients: 21,
                  entity_clients: 15,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/xbqJh',
                counts: {
                  clients: 9,
                  entity_clients: 6,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/m7vOo',
                counts: {
                  clients: 10,
                  entity_clients: 6,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/lULhW',
                counts: {
                  clients: 9,
                  entity_clients: 7,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/hB9qn',
                counts: {
                  clients: 10,
                  entity_clients: 8,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/RIEKI',
                counts: {
                  clients: 1,
                  entity_clients: 0,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'DcgzU',
            namespace_path: 'namespace17/',
            counts: {
              distinct_entities: 43,
              non_entity_tokens: 125,
              clients: 168,
            },
            mounts: [
              {
                path: 'auth/method/cdZ64',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/UpXi1',
                counts: {
                  clients: 33,
                  entity_clients: 19,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/6OzPw',
                counts: {
                  clients: 22,
                  entity_clients: 10,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/PkimI',
                counts: {
                  clients: 21,
                  entity_clients: 10,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/7ecN2',
                counts: {
                  clients: 14,
                  entity_clients: 10,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/AYdDo',
                counts: {
                  clients: 16,
                  entity_clients: 8,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/kS9h6',
                counts: {
                  clients: 9,
                  entity_clients: 4,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/dIoMU',
                counts: {
                  clients: 6,
                  entity_clients: 3,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/eXB1u',
                counts: {
                  clients: 7,
                  entity_clients: 3,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/SQ8Ty',
                counts: {
                  clients: 5,
                  entity_clients: 3,
                  non_entity_clients: 2,
                },
              },
            ],
          },
          {
            namespace_id: 'Uc0o8',
            namespace_path: 'namespace16/',
            counts: {
              distinct_entities: 56,
              non_entity_tokens: 112,
              clients: 168,
            },
            mounts: [
              {
                path: 'auth/method/my50c',
                counts: {
                  clients: 33,
                  entity_clients: 18,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/D8zfa',
                counts: {
                  clients: 29,
                  entity_clients: 17,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/w2xnA',
                counts: {
                  clients: 32,
                  entity_clients: 18,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/FwR7Z',
                counts: {
                  clients: 20,
                  entity_clients: 12,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/wwNCu',
                counts: {
                  clients: 17,
                  entity_clients: 10,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/vv2O6',
                counts: {
                  clients: 11,
                  entity_clients: 7,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/zRqUm',
                counts: {
                  clients: 9,
                  entity_clients: 6,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/Yez2v',
                counts: {
                  clients: 8,
                  entity_clients: 3,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/SBBJ2',
                counts: {
                  clients: 5,
                  entity_clients: 3,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/NNSCC',
                counts: {
                  clients: 4,
                  entity_clients: 1,
                  non_entity_clients: 3,
                },
              },
            ],
          },
          {
            namespace_id: 'Rqa3W',
            namespace_path: 'namespace13/',
            counts: {
              distinct_entities: 9,
              non_entity_tokens: 156,
              clients: 165,
            },
            mounts: [
              {
                path: 'auth/method/KPlRb',
                counts: {
                  clients: 33,
                  entity_clients: 20,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/199gy',
                counts: {
                  clients: 29,
                  entity_clients: 19,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/UDpxk',
                counts: {
                  clients: 24,
                  entity_clients: 14,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/bmgSl',
                counts: {
                  clients: 21,
                  entity_clients: 13,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/oyWlP',
                counts: {
                  clients: 20,
                  entity_clients: 12,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/z7Uka',
                counts: {
                  clients: 15,
                  entity_clients: 5,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/ftNn7',
                counts: {
                  clients: 10,
                  entity_clients: 6,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/pvdQ7',
                counts: {
                  clients: 9,
                  entity_clients: 7,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/DsnIn',
                counts: {
                  clients: 4,
                  entity_clients: 2,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/E1YLg',
                counts: {
                  clients: 0,
                  entity_clients: 0,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'aJuQG',
            namespace_path: 'namespace9/',
            counts: {
              distinct_entities: 80,
              non_entity_tokens: 71,
              clients: 151,
            },
            mounts: [
              {
                path: 'auth/method/RCpUn',
                counts: {
                  clients: 31,
                  entity_clients: 18,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/S0O4t',
                counts: {
                  clients: 21,
                  entity_clients: 7,
                  non_entity_clients: 14,
                },
              },
              {
                path: 'auth/method/QqXfg',
                counts: {
                  clients: 25,
                  entity_clients: 12,
                  non_entity_clients: 13,
                },
              },
              {
                path: 'auth/method/CSSoi',
                counts: {
                  clients: 23,
                  entity_clients: 13,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/klonh',
                counts: {
                  clients: 15,
                  entity_clients: 7,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/JyhFQ',
                counts: {
                  clients: 15,
                  entity_clients: 5,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/S66CH',
                counts: {
                  clients: 7,
                  entity_clients: 4,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/6pBz3',
                counts: {
                  clients: 6,
                  entity_clients: 4,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/qHCZa',
                counts: {
                  clients: 6,
                  entity_clients: 4,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/I6OpF',
                counts: {
                  clients: 2,
                  entity_clients: 2,
                  non_entity_clients: 0,
                },
              },
            ],
          },
        ],
        months: [
          {
            timestamp: '2022-01-01T08:00:00.000Z',
            counts: {
              clients: 68,
              entity_clients: 23,
              non_entity_clients: 45,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 14,
                  non_entity_clients: 18,
                  clients: 32,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 4,
                  non_entity_clients: 13,
                  clients: 17,
                },
                mounts: [
                  {
                    path: 'auth/method/KPlRb',
                    counts: {
                      clients: 671,
                      entity_clients: 325,
                      non_entity_clients: 346,
                    },
                  },
                  {
                    path: 'auth/method/199gy',
                    counts: {
                      clients: 270,
                      entity_clients: 185,
                      non_entity_clients: 85,
                    },
                  },
                  {
                    path: 'auth/method/UDpxk',
                    counts: {
                      clients: 6,
                      entity_clients: 3,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/bmgSl',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                ],
                id: 'namespace13/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 3,
                  non_entity_clients: 12,
                  clients: 15,
                },
                mounts: [
                  {
                    path: 'auth/method/uMGBU',
                    counts: {
                      clients: 690,
                      entity_clients: 237,
                      non_entity_clients: 453,
                    },
                  },
                  {
                    path: 'auth/method/8YJO3',
                    counts: {
                      clients: 685,
                      entity_clients: 394,
                      non_entity_clients: 291,
                    },
                  },
                  {
                    path: 'auth/method/Ro774',
                    counts: {
                      clients: 21,
                      entity_clients: 8,
                      non_entity_clients: 13,
                    },
                  },
                  {
                    path: 'auth/method/ZIpjT',
                    counts: {
                      clients: 6,
                      entity_clients: 2,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 2,
                  non_entity_clients: 2,
                  clients: 4,
                },
                mounts: [
                  {
                    path: 'auth/method/qcuLl',
                    counts: {
                      clients: 553,
                      entity_clients: 300,
                      non_entity_clients: 253,
                    },
                  },
                  {
                    path: 'auth/method/KGWiS',
                    counts: {
                      clients: 89,
                      entity_clients: 63,
                      non_entity_clients: 26,
                    },
                  },
                  {
                    path: 'auth/method/iM8pi',
                    counts: {
                      clients: 387,
                      entity_clients: 261,
                      non_entity_clients: 126,
                    },
                  },
                  {
                    path: 'auth/method/IeyA4',
                    counts: {
                      clients: 315,
                      entity_clients: 262,
                      non_entity_clients: 53,
                    },
                  },
                ],
                id: 'namespace12/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 47,
                entity_clients: 11,
                non_entity_clients: 36,
              },
              namespaces: [
                {
                  counts: {
                    clients: 14,
                    entity_clients: 11,
                    non_entity_clients: 3,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 6,
                        entity_clients: 4,
                        non_entity_clients: 813,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 5,
                        entity_clients: 3,
                        non_entity_clients: 113,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 45,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 168,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 11,
                    entity_clients: 9,
                    non_entity_clients: 2,
                  },
                  mounts: [
                    {
                      path: 'auth/method/KPlRb',
                      counts: {
                        clients: 5,
                        entity_clients: 1,
                        non_entity_clients: 346,
                      },
                    },
                    {
                      path: 'auth/method/199gy',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 85,
                      },
                    },
                    {
                      path: 'auth/method/UDpxk',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/bmgSl',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 2,
                      },
                    },
                  ],
                  id: 'namespace13/',
                  path: '',
                },
                {
                  counts: {
                    clients: 10,
                    entity_clients: 8,
                    non_entity_clients: 2,
                  },
                  mounts: [
                    {
                      path: 'auth/method/uMGBU',
                      counts: {
                        clients: 4,
                        entity_clients: 1,
                        non_entity_clients: 453,
                      },
                    },
                    {
                      path: 'auth/method/8YJO3',
                      counts: {
                        clients: 3,
                        entity_clients: 1,
                        non_entity_clients: 291,
                      },
                    },
                    {
                      path: 'auth/method/Ro774',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 13,
                      },
                    },
                    {
                      path: 'auth/method/ZIpjT',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 4,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
                {
                  counts: {
                    clients: 1,
                    entity_clients: 1,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qcuLl',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 253,
                      },
                    },
                    {
                      path: 'auth/method/KGWiS',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 26,
                      },
                    },
                    {
                      path: 'auth/method/iM8pi',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 126,
                      },
                    },
                    {
                      path: 'auth/method/IeyA4',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 53,
                      },
                    },
                  ],
                  id: 'namespace12/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-02-01T08:00:00.000Z',
            counts: {
              clients: 115,
              entity_clients: 95,
              non_entity_clients: 20,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 83,
                  non_entity_clients: 17,
                  clients: 100,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 10,
                  non_entity_clients: 1,
                  clients: 11,
                },
                mounts: [
                  {
                    path: 'auth/method/824CE',
                    counts: {
                      clients: 50,
                      entity_clients: 47,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/r2zb4',
                    counts: {
                      clients: 593,
                      entity_clients: 98,
                      non_entity_clients: 495,
                    },
                  },
                  {
                    path: 'auth/method/1zfD6',
                    counts: {
                      clients: 37,
                      entity_clients: 29,
                      non_entity_clients: 8,
                    },
                  },
                  {
                    path: 'auth/method/L14lj',
                    counts: {
                      clients: 7,
                      entity_clients: 2,
                      non_entity_clients: 5,
                    },
                  },
                ],
                id: 'namespace2/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 1,
                  non_entity_clients: 1,
                  clients: 2,
                },
                mounts: [
                  {
                    path: 'auth/method/uMGBU',
                    counts: {
                      clients: 690,
                      entity_clients: 237,
                      non_entity_clients: 453,
                    },
                  },
                  {
                    path: 'auth/method/8YJO3',
                    counts: {
                      clients: 685,
                      entity_clients: 394,
                      non_entity_clients: 291,
                    },
                  },
                  {
                    path: 'auth/method/Ro774',
                    counts: {
                      clients: 21,
                      entity_clients: 8,
                      non_entity_clients: 13,
                    },
                  },
                  {
                    path: 'auth/method/ZIpjT',
                    counts: {
                      clients: 6,
                      entity_clients: 2,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 1,
                  non_entity_clients: 1,
                  clients: 2,
                },
                mounts: [
                  {
                    path: 'auth/method/qT4Wl',
                    counts: {
                      clients: 1076,
                      entity_clients: 57,
                      non_entity_clients: 1019,
                    },
                  },
                  {
                    path: 'auth/method/Vhu56',
                    counts: {
                      clients: 23,
                      entity_clients: 12,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/PCc58',
                    counts: {
                      clients: 110,
                      entity_clients: 81,
                      non_entity_clients: 29,
                    },
                  },
                  {
                    path: 'auth/method/nPP4c',
                    counts: {
                      clients: 12,
                      entity_clients: 5,
                      non_entity_clients: 7,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 30,
                entity_clients: 26,
                non_entity_clients: 4,
              },
              namespaces: [
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: -1,
                        entity_clients: -1,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/824CE',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/r2zb4',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/1zfD6',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/L14lj',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace2/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/uMGBU',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/8YJO3',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/Ro774',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ZIpjT',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: -1,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qT4Wl',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/Vhu56',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PCc58',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/nPP4c',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace11/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-03-01T08:00:00.000Z',
            counts: {
              clients: 145,
              entity_clients: 121,
              non_entity_clients: 24,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 83,
                  non_entity_clients: 20,
                  clients: 103,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 20,
                  non_entity_clients: 2,
                  clients: 22,
                },
                mounts: [
                  {
                    path: 'auth/method/824CE',
                    counts: {
                      clients: 50,
                      entity_clients: 47,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/r2zb4',
                    counts: {
                      clients: 593,
                      entity_clients: 98,
                      non_entity_clients: 495,
                    },
                  },
                  {
                    path: 'auth/method/1zfD6',
                    counts: {
                      clients: 37,
                      entity_clients: 29,
                      non_entity_clients: 8,
                    },
                  },
                  {
                    path: 'auth/method/L14lj',
                    counts: {
                      clients: 7,
                      entity_clients: 2,
                      non_entity_clients: 5,
                    },
                  },
                ],
                id: 'namespace2/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 15,
                  non_entity_clients: 1,
                  clients: 16,
                },
                mounts: [
                  {
                    path: 'auth/method/uMGBU',
                    counts: {
                      clients: 690,
                      entity_clients: 237,
                      non_entity_clients: 453,
                    },
                  },
                  {
                    path: 'auth/method/8YJO3',
                    counts: {
                      clients: 685,
                      entity_clients: 394,
                      non_entity_clients: 291,
                    },
                  },
                  {
                    path: 'auth/method/Ro774',
                    counts: {
                      clients: 21,
                      entity_clients: 8,
                      non_entity_clients: 13,
                    },
                  },
                  {
                    path: 'auth/method/ZIpjT',
                    counts: {
                      clients: 6,
                      entity_clients: 2,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 3,
                  non_entity_clients: 1,
                  clients: 4,
                },
                mounts: [
                  {
                    path: 'auth/method/qT4Wl',
                    counts: {
                      clients: 1076,
                      entity_clients: 57,
                      non_entity_clients: 1019,
                    },
                  },
                  {
                    path: 'auth/method/Vhu56',
                    counts: {
                      clients: 23,
                      entity_clients: 12,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/PCc58',
                    counts: {
                      clients: 110,
                      entity_clients: 81,
                      non_entity_clients: 29,
                    },
                  },
                  {
                    path: 'auth/method/nPP4c',
                    counts: {
                      clients: 12,
                      entity_clients: 5,
                      non_entity_clients: 7,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 29,
                entity_clients: 21,
                non_entity_clients: 8,
              },
              namespaces: [
                {
                  counts: {
                    clients: 127,
                    entity_clients: 119,
                    non_entity_clients: 8,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 108,
                        entity_clients: 69,
                        non_entity_clients: 813,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 9,
                        entity_clients: 1,
                        non_entity_clients: 113,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 5,
                        entity_clients: 4,
                        non_entity_clients: 45,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 3,
                        entity_clients: 2,
                        non_entity_clients: 168,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 40,
                    entity_clients: 39,
                    non_entity_clients: 1,
                  },
                  mounts: [
                    {
                      path: 'auth/method/824CE',
                      counts: {
                        clients: 35,
                        entity_clients: 10,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/r2zb4',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 495,
                      },
                    },
                    {
                      path: 'auth/method/1zfD6',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 8,
                      },
                    },
                    {
                      path: 'auth/method/L14lj',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 5,
                      },
                    },
                  ],
                  id: 'namespace2/',
                  path: '',
                },
                {
                  counts: {
                    clients: 21,
                    entity_clients: 15,
                    non_entity_clients: 6,
                  },
                  mounts: [
                    {
                      path: 'auth/method/uMGBU',
                      counts: {
                        clients: 10,
                        entity_clients: 2,
                        non_entity_clients: 453,
                      },
                    },
                    {
                      path: 'auth/method/8YJO3',
                      counts: {
                        clients: 4,
                        entity_clients: 3,
                        non_entity_clients: 291,
                      },
                    },
                    {
                      path: 'auth/method/Ro774',
                      counts: {
                        clients: 3,
                        entity_clients: 1,
                        non_entity_clients: 13,
                      },
                    },
                    {
                      path: 'auth/method/ZIpjT',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 4,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
                {
                  counts: {
                    clients: 2,
                    entity_clients: 1,
                    non_entity_clients: 1,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qT4Wl',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 1019,
                      },
                    },
                    {
                      path: 'auth/method/Vhu56',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 11,
                      },
                    },
                    {
                      path: 'auth/method/PCc58',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 29,
                      },
                    },
                    {
                      path: 'auth/method/nPP4c',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 7,
                      },
                    },
                  ],
                  id: 'namespace11/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-04-01T07:00:00.000Z',
            counts: {
              clients: 174,
              entity_clients: 81,
              non_entity_clients: 93,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 53,
                  non_entity_clients: 83,
                  clients: 136,
                },
                mounts: [
                  {
                    path: 'auth/method/uMGBU',
                    counts: {
                      clients: 690,
                      entity_clients: 237,
                      non_entity_clients: 453,
                    },
                  },
                  {
                    path: 'auth/method/8YJO3',
                    counts: {
                      clients: 685,
                      entity_clients: 394,
                      non_entity_clients: 291,
                    },
                  },
                  {
                    path: 'auth/method/Ro774',
                    counts: {
                      clients: 21,
                      entity_clients: 8,
                      non_entity_clients: 13,
                    },
                  },
                  {
                    path: 'auth/method/ZIpjT',
                    counts: {
                      clients: 6,
                      entity_clients: 2,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 14,
                  non_entity_clients: 5,
                  clients: 19,
                },
                mounts: [
                  {
                    path: 'auth/method/qT4Wl',
                    counts: {
                      clients: 1076,
                      entity_clients: 57,
                      non_entity_clients: 1019,
                    },
                  },
                  {
                    path: 'auth/method/Vhu56',
                    counts: {
                      clients: 23,
                      entity_clients: 12,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/PCc58',
                    counts: {
                      clients: 110,
                      entity_clients: 81,
                      non_entity_clients: 29,
                    },
                  },
                  {
                    path: 'auth/method/nPP4c',
                    counts: {
                      clients: 12,
                      entity_clients: 5,
                      non_entity_clients: 7,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 7,
                  non_entity_clients: 4,
                  clients: 11,
                },
                mounts: [
                  {
                    path: 'auth/method/xYL0l',
                    counts: {
                      clients: 854,
                      entity_clients: 659,
                      non_entity_clients: 195,
                    },
                  },
                  {
                    path: 'auth/method/CwWM7',
                    counts: {
                      clients: 197,
                      entity_clients: 109,
                      non_entity_clients: 88,
                    },
                  },
                  {
                    path: 'auth/method/swCd0',
                    counts: {
                      clients: 76,
                      entity_clients: 18,
                      non_entity_clients: 58,
                    },
                  },
                  {
                    path: 'auth/method/0CZTs',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 7,
                  non_entity_clients: 1,
                  clients: 8,
                },
                mounts: [
                  {
                    path: 'auth/method/RCpUn',
                    counts: {
                      clients: 702,
                      entity_clients: 257,
                      non_entity_clients: 445,
                    },
                  },
                  {
                    path: 'auth/method/S0O4t',
                    counts: {
                      clients: 441,
                      entity_clients: 63,
                      non_entity_clients: 378,
                    },
                  },
                  {
                    path: 'auth/method/QqXfg',
                    counts: {
                      clients: 12,
                      entity_clients: 9,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/CSSoi',
                    counts: {
                      clients: 8,
                      entity_clients: 4,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespace9/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 20,
                entity_clients: 5,
                non_entity_clients: 15,
              },
              namespaces: [
                {
                  counts: {
                    clients: 725,
                    entity_clients: 416,
                    non_entity_clients: 309,
                  },
                  mounts: [
                    {
                      path: 'auth/method/uMGBU',
                      counts: {
                        clients: 281,
                        entity_clients: 48,
                        non_entity_clients: 453,
                      },
                    },
                    {
                      path: 'auth/method/8YJO3',
                      counts: {
                        clients: 252,
                        entity_clients: 142,
                        non_entity_clients: 291,
                      },
                    },
                    {
                      path: 'auth/method/Ro774',
                      counts: {
                        clients: 93,
                        entity_clients: 43,
                        non_entity_clients: 13,
                      },
                    },
                    {
                      path: 'auth/method/ZIpjT',
                      counts: {
                        clients: 86,
                        entity_clients: 19,
                        non_entity_clients: 4,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
                {
                  counts: {
                    clients: 178,
                    entity_clients: 99,
                    non_entity_clients: 79,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qT4Wl',
                      counts: {
                        clients: 76,
                        entity_clients: 69,
                        non_entity_clients: 1019,
                      },
                    },
                    {
                      path: 'auth/method/Vhu56',
                      counts: {
                        clients: 64,
                        entity_clients: 38,
                        non_entity_clients: 11,
                      },
                    },
                    {
                      path: 'auth/method/PCc58',
                      counts: {
                        clients: 20,
                        entity_clients: 12,
                        non_entity_clients: 29,
                      },
                    },
                    {
                      path: 'auth/method/nPP4c',
                      counts: {
                        clients: 15,
                        entity_clients: 12,
                        non_entity_clients: 7,
                      },
                    },
                  ],
                  id: 'namespace11/',
                  path: '',
                },
                {
                  counts: {
                    clients: 40,
                    entity_clients: 28,
                    non_entity_clients: 12,
                  },
                  mounts: [
                    {
                      path: 'auth/method/xYL0l',
                      counts: {
                        clients: 24,
                        entity_clients: 9,
                        non_entity_clients: 195,
                      },
                    },
                    {
                      path: 'auth/method/CwWM7',
                      counts: {
                        clients: 10,
                        entity_clients: 8,
                        non_entity_clients: 88,
                      },
                    },
                    {
                      path: 'auth/method/swCd0',
                      counts: {
                        clients: 4,
                        entity_clients: 2,
                        non_entity_clients: 58,
                      },
                    },
                    {
                      path: 'auth/method/0CZTs',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 1,
                      },
                    },
                  ],
                  id: 'namespace10/',
                  path: '',
                },
                {
                  counts: {
                    clients: 34,
                    entity_clients: 21,
                    non_entity_clients: 13,
                  },
                  mounts: [
                    {
                      path: 'auth/method/RCpUn',
                      counts: {
                        clients: 27,
                        entity_clients: 14,
                        non_entity_clients: 445,
                      },
                    },
                    {
                      path: 'auth/method/S0O4t',
                      counts: {
                        clients: 4,
                        entity_clients: 1,
                        non_entity_clients: 378,
                      },
                    },
                    {
                      path: 'auth/method/QqXfg',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/CSSoi',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 4,
                      },
                    },
                  ],
                  id: 'namespace9/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-05-01T07:00:00.000Z',
            counts: {
              clients: 194,
              entity_clients: 187,
              non_entity_clients: 7,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 106,
                  non_entity_clients: 3,
                  clients: 109,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 80,
                  non_entity_clients: 3,
                  clients: 83,
                },
                mounts: [
                  {
                    path: 'auth/method/qcuLl',
                    counts: {
                      clients: 553,
                      entity_clients: 300,
                      non_entity_clients: 253,
                    },
                  },
                  {
                    path: 'auth/method/KGWiS',
                    counts: {
                      clients: 89,
                      entity_clients: 63,
                      non_entity_clients: 26,
                    },
                  },
                  {
                    path: 'auth/method/iM8pi',
                    counts: {
                      clients: 387,
                      entity_clients: 261,
                      non_entity_clients: 126,
                    },
                  },
                  {
                    path: 'auth/method/IeyA4',
                    counts: {
                      clients: 315,
                      entity_clients: 262,
                      non_entity_clients: 53,
                    },
                  },
                ],
                id: 'namespace12/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 1,
                  non_entity_clients: 1,
                  clients: 2,
                },
                mounts: [
                  {
                    path: 'auth/method/XQUrA',
                    counts: {
                      clients: 233,
                      entity_clients: 69,
                      non_entity_clients: 164,
                    },
                  },
                  {
                    path: 'auth/method/1p6HR',
                    counts: {
                      clients: 454,
                      entity_clients: 162,
                      non_entity_clients: 292,
                    },
                  },
                  {
                    path: 'auth/method/qRjoJ',
                    counts: {
                      clients: 49,
                      entity_clients: 2,
                      non_entity_clients: 47,
                    },
                  },
                  {
                    path: 'auth/method/x9QQB',
                    counts: {
                      clients: 122,
                      entity_clients: 89,
                      non_entity_clients: 33,
                    },
                  },
                ],
                id: 'namespace6/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 38,
                entity_clients: 29,
                non_entity_clients: 9,
              },
              namespaces: [
                {
                  counts: {
                    clients: 62,
                    entity_clients: 37,
                    non_entity_clients: 25,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 57,
                        entity_clients: 26,
                        non_entity_clients: 813,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 113,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 45,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 168,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 31,
                    entity_clients: 20,
                    non_entity_clients: 11,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qcuLl',
                      counts: {
                        clients: 27,
                        entity_clients: 24,
                        non_entity_clients: 253,
                      },
                    },
                    {
                      path: 'auth/method/KGWiS',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 26,
                      },
                    },
                    {
                      path: 'auth/method/iM8pi',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 126,
                      },
                    },
                    {
                      path: 'auth/method/IeyA4',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 53,
                      },
                    },
                  ],
                  id: 'namespace12/',
                  path: '',
                },
                {
                  counts: {
                    clients: 30,
                    entity_clients: 21,
                    non_entity_clients: 9,
                  },
                  mounts: [
                    {
                      path: 'auth/method/XQUrA',
                      counts: {
                        clients: 15,
                        entity_clients: 5,
                        non_entity_clients: 164,
                      },
                    },
                    {
                      path: 'auth/method/1p6HR',
                      counts: {
                        clients: 7,
                        entity_clients: 4,
                        non_entity_clients: 292,
                      },
                    },
                    {
                      path: 'auth/method/qRjoJ',
                      counts: {
                        clients: 6,
                        entity_clients: 1,
                        non_entity_clients: 47,
                      },
                    },
                    {
                      path: 'auth/method/x9QQB',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 33,
                      },
                    },
                  ],
                  id: 'namespace6/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-06-01T07:00:00.000Z',
            counts: {
              clients: 232,
              entity_clients: 47,
              non_entity_clients: 185,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 21,
                  non_entity_clients: 102,
                  clients: 123,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 19,
                  non_entity_clients: 59,
                  clients: 78,
                },
                mounts: [
                  {
                    path: 'auth/method/qT4Wl',
                    counts: {
                      clients: 1076,
                      entity_clients: 57,
                      non_entity_clients: 1019,
                    },
                  },
                  {
                    path: 'auth/method/Vhu56',
                    counts: {
                      clients: 23,
                      entity_clients: 12,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/PCc58',
                    counts: {
                      clients: 110,
                      entity_clients: 81,
                      non_entity_clients: 29,
                    },
                  },
                  {
                    path: 'auth/method/nPP4c',
                    counts: {
                      clients: 12,
                      entity_clients: 5,
                      non_entity_clients: 7,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 7,
                  non_entity_clients: 24,
                  clients: 31,
                },
                mounts: [
                  {
                    path: 'auth/method/my50c',
                    counts: {
                      clients: 342,
                      entity_clients: 179,
                      non_entity_clients: 163,
                    },
                  },
                  {
                    path: 'auth/method/D8zfa',
                    counts: {
                      clients: 681,
                      entity_clients: 292,
                      non_entity_clients: 389,
                    },
                  },
                  {
                    path: 'auth/method/w2xnA',
                    counts: {
                      clients: 29,
                      entity_clients: 17,
                      non_entity_clients: 12,
                    },
                  },
                  {
                    path: 'auth/method/FwR7Z',
                    counts: {
                      clients: 40,
                      entity_clients: 4,
                      non_entity_clients: 36,
                    },
                  },
                ],
                id: 'namespace16/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 71,
                entity_clients: 49,
                non_entity_clients: 22,
              },
              namespaces: [
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qT4Wl',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/Vhu56',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PCc58',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/nPP4c',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace11/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/my50c',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/D8zfa',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/w2xnA',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/FwR7Z',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace16/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-07-01T07:00:00.000Z',
            counts: {
              clients: 303,
              entity_clients: 218,
              non_entity_clients: 85,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 120,
                  non_entity_clients: 44,
                  clients: 164,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 65,
                  non_entity_clients: 35,
                  clients: 100,
                },
                mounts: [
                  {
                    path: 'auth/method/qT4Wl',
                    counts: {
                      clients: 1076,
                      entity_clients: 57,
                      non_entity_clients: 1019,
                    },
                  },
                  {
                    path: 'auth/method/Vhu56',
                    counts: {
                      clients: 23,
                      entity_clients: 12,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/PCc58',
                    counts: {
                      clients: 110,
                      entity_clients: 81,
                      non_entity_clients: 29,
                    },
                  },
                  {
                    path: 'auth/method/nPP4c',
                    counts: {
                      clients: 12,
                      entity_clients: 5,
                      non_entity_clients: 7,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 33,
                  non_entity_clients: 6,
                  clients: 39,
                },
                mounts: [
                  {
                    path: 'auth/method/my50c',
                    counts: {
                      clients: 342,
                      entity_clients: 179,
                      non_entity_clients: 163,
                    },
                  },
                  {
                    path: 'auth/method/D8zfa',
                    counts: {
                      clients: 681,
                      entity_clients: 292,
                      non_entity_clients: 389,
                    },
                  },
                  {
                    path: 'auth/method/w2xnA',
                    counts: {
                      clients: 29,
                      entity_clients: 17,
                      non_entity_clients: 12,
                    },
                  },
                  {
                    path: 'auth/method/FwR7Z',
                    counts: {
                      clients: 40,
                      entity_clients: 4,
                      non_entity_clients: 36,
                    },
                  },
                ],
                id: 'namespace16/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 72,
                entity_clients: 46,
                non_entity_clients: 26,
              },
              namespaces: [
                {
                  counts: {
                    clients: 41,
                    entity_clients: 39,
                    non_entity_clients: 2,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 21,
                        entity_clients: 18,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 12,
                        entity_clients: 9,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 3,
                        entity_clients: 2,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 3,
                        entity_clients: 2,
                        non_entity_clients: 1,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 14,
                    entity_clients: 13,
                    non_entity_clients: 1,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qT4Wl',
                      counts: {
                        clients: 8,
                        entity_clients: 6,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/Vhu56',
                      counts: {
                        clients: 3,
                        entity_clients: 1,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/PCc58',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/nPP4c',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace11/',
                  path: '',
                },
                {
                  counts: {
                    clients: 10,
                    entity_clients: 6,
                    non_entity_clients: 4,
                  },
                  mounts: [
                    {
                      path: 'auth/method/my50c',
                      counts: {
                        clients: 6,
                        entity_clients: 5,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/D8zfa',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/w2xnA',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/FwR7Z',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace16/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-08-01T07:00:00.000Z',
            counts: {
              clients: 375,
              entity_clients: 80,
              non_entity_clients: 295,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 63,
                  non_entity_clients: 216,
                  clients: 279,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 14,
                  non_entity_clients: 45,
                  clients: 59,
                },
                mounts: [
                  {
                    path: 'auth/method/xYL0l',
                    counts: {
                      clients: 854,
                      entity_clients: 659,
                      non_entity_clients: 195,
                    },
                  },
                  {
                    path: 'auth/method/CwWM7',
                    counts: {
                      clients: 197,
                      entity_clients: 109,
                      non_entity_clients: 88,
                    },
                  },
                  {
                    path: 'auth/method/swCd0',
                    counts: {
                      clients: 76,
                      entity_clients: 18,
                      non_entity_clients: 58,
                    },
                  },
                  {
                    path: 'auth/method/0CZTs',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 3,
                  non_entity_clients: 34,
                  clients: 37,
                },
                mounts: [
                  {
                    path: 'auth/method/RCpUn',
                    counts: {
                      clients: 702,
                      entity_clients: 257,
                      non_entity_clients: 445,
                    },
                  },
                  {
                    path: 'auth/method/S0O4t',
                    counts: {
                      clients: 441,
                      entity_clients: 63,
                      non_entity_clients: 378,
                    },
                  },
                  {
                    path: 'auth/method/QqXfg',
                    counts: {
                      clients: 12,
                      entity_clients: 9,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/CSSoi',
                    counts: {
                      clients: 8,
                      entity_clients: 4,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespace9/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 0,
                entity_clients: 0,
                non_entity_clients: 0,
              },
              namespaces: [
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/xYL0l',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/CwWM7',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/swCd0',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/0CZTs',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace10/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/RCpUn',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/S0O4t',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/QqXfg',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/CSSoi',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace9/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-09-01T07:00:00.000Z',
            counts: {
              clients: 375,
              entity_clients: 67,
              non_entity_clients: 308,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 34,
                  non_entity_clients: 230,
                  clients: 264,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 24,
                  non_entity_clients: 48,
                  clients: 72,
                },
                mounts: [
                  {
                    path: 'auth/method/xYL0l',
                    counts: {
                      clients: 854,
                      entity_clients: 659,
                      non_entity_clients: 195,
                    },
                  },
                  {
                    path: 'auth/method/CwWM7',
                    counts: {
                      clients: 197,
                      entity_clients: 109,
                      non_entity_clients: 88,
                    },
                  },
                  {
                    path: 'auth/method/swCd0',
                    counts: {
                      clients: 76,
                      entity_clients: 18,
                      non_entity_clients: 58,
                    },
                  },
                  {
                    path: 'auth/method/0CZTs',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 9,
                  non_entity_clients: 30,
                  clients: 39,
                },
                mounts: [
                  {
                    path: 'auth/method/RCpUn',
                    counts: {
                      clients: 702,
                      entity_clients: 257,
                      non_entity_clients: 445,
                    },
                  },
                  {
                    path: 'auth/method/S0O4t',
                    counts: {
                      clients: 441,
                      entity_clients: 63,
                      non_entity_clients: 378,
                    },
                  },
                  {
                    path: 'auth/method/QqXfg',
                    counts: {
                      clients: 12,
                      entity_clients: 9,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/CSSoi',
                    counts: {
                      clients: 8,
                      entity_clients: 4,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespace9/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 114,
                entity_clients: 63,
                non_entity_clients: 51,
              },
              namespaces: [
                {
                  counts: {
                    clients: 7,
                    entity_clients: 6,
                    non_entity_clients: 1,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 2,
                        entity_clients: 2,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 5,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 1,
                    entity_clients: 1,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/xYL0l',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/CwWM7',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/swCd0',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/0CZTs',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace10/',
                  path: '',
                },
                {
                  counts: {
                    clients: 1,
                    entity_clients: 1,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/RCpUn',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/S0O4t',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/QqXfg',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/CSSoi',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace9/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-10-01T07:00:00.000Z',

            counts: {
              clients: 489,
              entity_clients: 134,
              non_entity_clients: 355,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 108,
                  non_entity_clients: 322,
                  clients: 430,
                },
                mounts: [
                  {
                    path: 'auth/method/uMGBU',
                    counts: {
                      clients: 690,
                      entity_clients: 237,
                      non_entity_clients: 453,
                    },
                  },
                  {
                    path: 'auth/method/8YJO3',
                    counts: {
                      clients: 685,
                      entity_clients: 394,
                      non_entity_clients: 291,
                    },
                  },
                  {
                    path: 'auth/method/Ro774',
                    counts: {
                      clients: 21,
                      entity_clients: 8,
                      non_entity_clients: 13,
                    },
                  },
                  {
                    path: 'auth/method/ZIpjT',
                    counts: {
                      clients: 6,
                      entity_clients: 2,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 15,
                  non_entity_clients: 26,
                  clients: 41,
                },
                mounts: [
                  {
                    path: 'auth/method/qT4Wl',
                    counts: {
                      clients: 1076,
                      entity_clients: 57,
                      non_entity_clients: 1019,
                    },
                  },
                  {
                    path: 'auth/method/Vhu56',
                    counts: {
                      clients: 23,
                      entity_clients: 12,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/PCc58',
                    counts: {
                      clients: 110,
                      entity_clients: 81,
                      non_entity_clients: 29,
                    },
                  },
                  {
                    path: 'auth/method/nPP4c',
                    counts: {
                      clients: 12,
                      entity_clients: 5,
                      non_entity_clients: 7,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 11,
                  non_entity_clients: 7,
                  clients: 18,
                },
                mounts: [
                  {
                    path: 'auth/method/xYL0l',
                    counts: {
                      clients: 854,
                      entity_clients: 659,
                      non_entity_clients: 195,
                    },
                  },
                  {
                    path: 'auth/method/CwWM7',
                    counts: {
                      clients: 197,
                      entity_clients: 109,
                      non_entity_clients: 88,
                    },
                  },
                  {
                    path: 'auth/method/swCd0',
                    counts: {
                      clients: 76,
                      entity_clients: 18,
                      non_entity_clients: 58,
                    },
                  },
                  {
                    path: 'auth/method/0CZTs',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 21,
                entity_clients: 11,
                non_entity_clients: 10,
              },
              namespaces: [
                {
                  counts: {
                    clients: 154,
                    entity_clients: 129,
                    non_entity_clients: 25,
                  },
                  mounts: [
                    {
                      path: 'auth/method/uMGBU',
                      counts: {
                        clients: 114,
                        entity_clients: 84,
                        non_entity_clients: 453,
                      },
                    },
                    {
                      path: 'auth/method/8YJO3',
                      counts: {
                        clients: 33,
                        entity_clients: 25,
                        non_entity_clients: 291,
                      },
                    },
                    {
                      path: 'auth/method/Ro774',
                      counts: {
                        clients: 4,
                        entity_clients: 2,
                        non_entity_clients: 13,
                      },
                    },
                    {
                      path: 'auth/method/ZIpjT',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 4,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
                {
                  counts: {
                    clients: 6,
                    entity_clients: 3,
                    non_entity_clients: 3,
                  },
                  mounts: [
                    {
                      path: 'auth/method/qT4Wl',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 1019,
                      },
                    },
                    {
                      path: 'auth/method/Vhu56',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 11,
                      },
                    },
                    {
                      path: 'auth/method/PCc58',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 29,
                      },
                    },
                    {
                      path: 'auth/method/nPP4c',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 7,
                      },
                    },
                  ],
                  id: 'namespace11/',
                  path: '',
                },
                {
                  counts: {
                    clients: 4,
                    entity_clients: 2,
                    non_entity_clients: 2,
                  },
                  mounts: [
                    {
                      path: 'auth/method/xYL0l',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 195,
                      },
                    },
                    {
                      path: 'auth/method/CwWM7',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 88,
                      },
                    },
                    {
                      path: 'auth/method/swCd0',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 58,
                      },
                    },
                    {
                      path: 'auth/method/0CZTs',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 1,
                      },
                    },
                  ],
                  id: 'namespace10/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-11-01T07:00:00.000Z',

            counts: {
              clients: 510,
              entity_clients: 164,
              non_entity_clients: 346,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 66,
                  non_entity_clients: 176,
                  clients: 242,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 51,
                  non_entity_clients: 120,
                  clients: 171,
                },
                mounts: [
                  {
                    path: 'auth/method/lDz9c',
                    counts: {
                      clients: 599,
                      entity_clients: 207,
                      non_entity_clients: 392,
                    },
                  },
                  {
                    path: 'auth/method/GtbUu',
                    counts: {
                      clients: 38,
                      entity_clients: 6,
                      non_entity_clients: 32,
                    },
                  },
                  {
                    path: 'auth/method/WCyYz',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/j227p',
                    counts: {
                      clients: 64,
                      entity_clients: 13,
                      non_entity_clients: 51,
                    },
                  },
                ],
                id: 'namespacelonglonglong3/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 47,
                  non_entity_clients: 50,
                  clients: 97,
                },
                mounts: [
                  {
                    path: 'auth/method/uMGBU',
                    counts: {
                      clients: 690,
                      entity_clients: 237,
                      non_entity_clients: 453,
                    },
                  },
                  {
                    path: 'auth/method/8YJO3',
                    counts: {
                      clients: 685,
                      entity_clients: 394,
                      non_entity_clients: 291,
                    },
                  },
                  {
                    path: 'auth/method/Ro774',
                    counts: {
                      clients: 21,
                      entity_clients: 8,
                      non_entity_clients: 13,
                    },
                  },
                  {
                    path: 'auth/method/ZIpjT',
                    counts: {
                      clients: 6,
                      entity_clients: 2,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 27,
                entity_clients: 8,
                non_entity_clients: 19,
              },
              namespaces: [
                {
                  counts: {
                    clients: 79,
                    entity_clients: 40,
                    non_entity_clients: 39,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 43,
                        entity_clients: 19,
                        non_entity_clients: 813,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 19,
                        entity_clients: 15,
                        non_entity_clients: 113,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 9,
                        entity_clients: 6,
                        non_entity_clients: 45,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 4,
                        entity_clients: 2,
                        non_entity_clients: 168,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 4,
                    entity_clients: 3,
                    non_entity_clients: 1,
                  },
                  mounts: [
                    {
                      path: 'auth/method/lDz9c',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 392,
                      },
                    },
                    {
                      path: 'auth/method/GtbUu',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 32,
                      },
                    },
                    {
                      path: 'auth/method/WCyYz',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/j227p',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 51,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong3/',
                  path: '',
                },
                {
                  counts: {
                    clients: 3,
                    entity_clients: 2,
                    non_entity_clients: 1,
                  },
                  mounts: [
                    {
                      path: 'auth/method/uMGBU',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 453,
                      },
                    },
                    {
                      path: 'auth/method/8YJO3',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 291,
                      },
                    },
                    {
                      path: 'auth/method/Ro774',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 13,
                      },
                    },
                    {
                      path: 'auth/method/ZIpjT',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 4,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-12-01T08:00:00.000Z',
            counts: {
              clients: 537,
              entity_clients: 125,
              non_entity_clients: 412,
            },
            namespaces: [
              {
                counts: {
                  entity_clients: 89,
                  non_entity_clients: 188,
                  clients: 277,
                },
                mounts: [
                  {
                    path: 'auth/method/cdZ64',
                    counts: {
                      clients: 817,
                      entity_clients: 4,
                      non_entity_clients: 813,
                    },
                  },
                  {
                    path: 'auth/method/UpXi1',
                    counts: {
                      clients: 385,
                      entity_clients: 272,
                      non_entity_clients: 113,
                    },
                  },
                  {
                    path: 'auth/method/6OzPw',
                    counts: {
                      clients: 93,
                      entity_clients: 48,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/PkimI',
                    counts: {
                      clients: 172,
                      entity_clients: 4,
                      non_entity_clients: 168,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 24,
                  non_entity_clients: 187,
                  clients: 211,
                },
                mounts: [
                  {
                    path: 'auth/method/u2r0G',
                    counts: {
                      clients: 247,
                      entity_clients: 37,
                      non_entity_clients: 210,
                    },
                  },
                  {
                    path: 'auth/method/mKqBV',
                    counts: {
                      clients: 336,
                      entity_clients: 320,
                      non_entity_clients: 16,
                    },
                  },
                  {
                    path: 'auth/method/nGOa2',
                    counts: {
                      clients: 22,
                      entity_clients: 3,
                      non_entity_clients: 19,
                    },
                  },
                  {
                    path: 'auth/method/46UKX',
                    counts: {
                      clients: 67,
                      entity_clients: 58,
                      non_entity_clients: 9,
                    },
                  },
                ],
                id: 'namespace15/',
                path: '',
              },
              {
                counts: {
                  entity_clients: 12,
                  non_entity_clients: 37,
                  clients: 49,
                },
                mounts: [
                  {
                    path: 'auth/method/uMGBU',
                    counts: {
                      clients: 690,
                      entity_clients: 237,
                      non_entity_clients: 453,
                    },
                  },
                  {
                    path: 'auth/method/8YJO3',
                    counts: {
                      clients: 685,
                      entity_clients: 394,
                      non_entity_clients: 291,
                    },
                  },
                  {
                    path: 'auth/method/Ro774',
                    counts: {
                      clients: 21,
                      entity_clients: 8,
                      non_entity_clients: 13,
                    },
                  },
                  {
                    path: 'auth/method/ZIpjT',
                    counts: {
                      clients: 6,
                      entity_clients: 2,
                      non_entity_clients: 4,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 0,
                entity_clients: 0,
                non_entity_clients: 0,
              },
              namespaces: [
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/cdZ64',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/UpXi1',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6OzPw',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PkimI',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace17/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/u2r0G',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/mKqBV',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/nGOa2',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/46UKX',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace15/',
                  path: '',
                },
                {
                  counts: {
                    clients: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                  },
                  mounts: [
                    {
                      path: 'auth/method/uMGBU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/8YJO3',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/Ro774',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ZIpjT',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
              ],
            },
          },
        ],
      },
    };
  });

  // current (<= 1.9) API response
  server.get('/sys/internal/counters/activity/monthly', function () {
    return {
      data: {
        by_namespace: [
          {
            namespace_id: 'Z4Rzh',
            namespace_path: 'namespace1/',
            counts: {
              distinct_entities: 867,
              non_entity_tokens: 939,
              clients: 1806,
            },
          },
          {
            namespace_id: 'DcgzU',
            namespace_path: 'namespace17/',
            counts: {
              distinct_entities: 966,
              non_entity_tokens: 550,
              clients: 1516,
            },
          },
          {
            namespace_id: '5SWT8',
            namespace_path: 'namespacelonglonglong4/',
            counts: {
              distinct_entities: 996,
              non_entity_tokens: 417,
              clients: 1413,
            },
          },
          {
            namespace_id: 'XGu7R',
            namespace_path: 'namespace12/',
            counts: {
              distinct_entities: 829,
              non_entity_tokens: 540,
              clients: 1369,
            },
          },
          {
            namespace_id: 'yHcL9',
            namespace_path: 'namespace11/',
            counts: {
              distinct_entities: 563,
              non_entity_tokens: 705,
              clients: 1268,
            },
          },
          {
            namespace_id: 'F0xGm',
            namespace_path: 'namespace10/',
            counts: {
              distinct_entities: 925,
              non_entity_tokens: 255,
              clients: 1180,
            },
          },
          {
            namespace_id: 'aJuQG',
            namespace_path: 'namespace9/',
            counts: {
              distinct_entities: 935,
              non_entity_tokens: 239,
              clients: 1174,
            },
          },
          {
            namespace_id: 'bw5UO',
            namespace_path: 'namespace6/',
            counts: {
              distinct_entities: 810,
              non_entity_tokens: 363,
              clients: 1173,
            },
          },
          {
            namespace_id: 'IeyJp',
            namespace_path: 'namespace14/',
            counts: {
              distinct_entities: 774,
              non_entity_tokens: 392,
              clients: 1166,
            },
          },
          {
            namespace_id: 'Uc0o8',
            namespace_path: 'namespace16/',
            counts: {
              distinct_entities: 408,
              non_entity_tokens: 743,
              clients: 1151,
            },
          },
          {
            namespace_id: 'R6L40',
            namespace_path: 'namespace2/',
            counts: {
              distinct_entities: 292,
              non_entity_tokens: 736,
              clients: 1028,
            },
          },
          {
            namespace_id: 'Rqa3W',
            namespace_path: 'namespace13/',
            counts: {
              distinct_entities: 160,
              non_entity_tokens: 803,
              clients: 963,
            },
          },
          {
            namespace_id: 'MSgZE',
            namespace_path: 'namespace7/',
            counts: {
              distinct_entities: 201,
              non_entity_tokens: 657,
              clients: 858,
            },
          },
          {
            namespace_id: 'kxU4t',
            namespace_path: 'namespacelonglonglong3/',
            counts: {
              distinct_entities: 742,
              non_entity_tokens: 26,
              clients: 768,
            },
          },
          {
            namespace_id: '5xKya',
            namespace_path: 'namespace15/',
            counts: {
              distinct_entities: 663,
              non_entity_tokens: 19,
              clients: 682,
            },
          },
          {
            namespace_id: '5KxXA',
            namespace_path: 'namespace18anotherlong/',
            counts: {
              distinct_entities: 470,
              non_entity_tokens: 196,
              clients: 666,
            },
          },
          {
            namespace_id: 'AAidI',
            namespace_path: 'namespace20/',
            counts: {
              distinct_entities: 429,
              non_entity_tokens: 60,
              clients: 489,
            },
          },
          {
            namespace_id: 'BCl56',
            namespace_path: 'namespace8/',
            counts: {
              distinct_entities: 61,
              non_entity_tokens: 201,
              clients: 262,
            },
          },
          {
            namespace_id: 'yYNw2',
            namespace_path: 'namespace19/',
            counts: {
              distinct_entities: 165,
              non_entity_tokens: 85,
              clients: 250,
            },
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 67,
              non_entity_tokens: 9,
              clients: 76,
            },
          },
        ],
        distinct_entities: 11323,
        non_entity_tokens: 7935,
        clients: 19258,
      },
    };
  });
}
