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
          clients: 300,
          _comment2: 'following 2 fields are deprecated',
          distinct_entities: 200,
          non_entity_tokens: 100,
        },
        by_namespace: [
          {
            namespace_id: 'Z4Rzh',
            namespace_path: 'namespace1/',
            counts: {
              distinct_entities: 867,
              non_entity_tokens: 939,
              clients: 1806,
            },
            mounts: [
              {
                path: 'auth/method/NqMeC',
                counts: {
                  clients: 1728,
                  entity_clients: 1378,
                  non_entity_clients: 350,
                },
              },
              {
                path: 'auth/method/S0FaZ',
                counts: {
                  clients: 11,
                  entity_clients: 6,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/vzH3z',
                counts: {
                  clients: 8,
                  entity_clients: 6,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/uP1zV',
                counts: {
                  clients: 43,
                  entity_clients: 36,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/yAga3',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/DTAFz',
                counts: {
                  clients: 5,
                  entity_clients: 2,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/Rk3Pt',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/wnNH5',
                counts: {
                  clients: 5,
                  entity_clients: 2,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/N3BJy',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/C5qsy',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'DcgzU',
            namespace_path: 'namespace17/',
            counts: {
              distinct_entities: 966,
              non_entity_tokens: 550,
              clients: 1516,
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
              {
                path: 'auth/method/7ecN2',
                counts: {
                  clients: 25,
                  entity_clients: 24,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/AYdDo',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/kS9h6',
                counts: {
                  clients: 4,
                  entity_clients: 1,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/dIoMU',
                counts: {
                  clients: 8,
                  entity_clients: 2,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/eXB1u',
                counts: {
                  clients: 7,
                  entity_clients: 1,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/SQ8Ty',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: '5SWT8',
            namespace_path: 'namespacelonglonglong4/',
            counts: {
              distinct_entities: 996,
              non_entity_tokens: 417,
              clients: 1413,
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
              {
                path: 'auth/method/jdRjF',
                counts: {
                  clients: 5,
                  entity_clients: 3,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/yyBoC',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/WLxYp',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/SNM6V',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/vNHtH',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/EqmlO',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'XGu7R',
            namespace_path: 'namespace12/',
            counts: {
              distinct_entities: 829,
              non_entity_tokens: 540,
              clients: 1369,
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
              {
                path: 'auth/method/KGFfV',
                counts: {
                  clients: 20,
                  entity_clients: 9,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/23AQk',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/PqTWe',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/pPSo1',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/HMu5H',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/xpOk3',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'yHcL9',
            namespace_path: 'namespace11/',
            counts: {
              distinct_entities: 563,
              non_entity_tokens: 705,
              clients: 1268,
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
              {
                path: 'auth/method/LY3am',
                counts: {
                  clients: 8,
                  entity_clients: 5,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/McQ4X',
                counts: {
                  clients: 6,
                  entity_clients: 4,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/NpjhH',
                counts: {
                  clients: 28,
                  entity_clients: 4,
                  non_entity_clients: 24,
                },
              },
              {
                path: 'auth/method/ToKO8',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/wfApH',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/L9uWV',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'F0xGm',
            namespace_path: 'namespace10/',
            counts: {
              distinct_entities: 925,
              non_entity_tokens: 255,
              clients: 1180,
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
              {
                path: 'auth/method/9v04G',
                counts: {
                  clients: 33,
                  entity_clients: 21,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/6hAlO',
                counts: {
                  clients: 11,
                  entity_clients: 10,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/ydSdP',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/i0CTY',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/nevwU',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/k2jYC',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'aJuQG',
            namespace_path: 'namespace9/',
            counts: {
              distinct_entities: 935,
              non_entity_tokens: 239,
              clients: 1174,
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
              {
                path: 'auth/method/klonh',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/JyhFQ',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/S66CH',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/6pBz3',
                counts: {
                  clients: 4,
                  entity_clients: 2,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/qHCZa',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/I6OpF',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'bw5UO',
            namespace_path: 'namespace6/',
            counts: {
              distinct_entities: 810,
              non_entity_tokens: 363,
              clients: 1173,
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
              {
                path: 'auth/method/rezK4',
                counts: {
                  clients: 119,
                  entity_clients: 74,
                  non_entity_clients: 45,
                },
              },
              {
                path: 'auth/method/qWNSS',
                counts: {
                  clients: 106,
                  entity_clients: 42,
                  non_entity_clients: 64,
                },
              },
              {
                path: 'auth/method/OmQEf',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/PhoAy',
                counts: {
                  clients: 9,
                  entity_clients: 8,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/aUuyM',
                counts: {
                  clients: 47,
                  entity_clients: 32,
                  non_entity_clients: 15,
                },
              },
              {
                path: 'auth/method/kUj1S',
                counts: {
                  clients: 32,
                  entity_clients: 1,
                  non_entity_clients: 31,
                },
              },
            ],
          },
          {
            namespace_id: 'IeyJp',
            namespace_path: 'namespace14/',
            counts: {
              distinct_entities: 774,
              non_entity_tokens: 392,
              clients: 1166,
            },
            mounts: [
              {
                path: 'auth/method/8NFVo',
                counts: {
                  clients: 143,
                  entity_clients: 28,
                  non_entity_clients: 115,
                },
              },
              {
                path: 'auth/method/XnNDy',
                counts: {
                  clients: 777,
                  entity_clients: 637,
                  non_entity_clients: 140,
                },
              },
              {
                path: 'auth/method/RYrzg',
                counts: {
                  clients: 113,
                  entity_clients: 68,
                  non_entity_clients: 45,
                },
              },
              {
                path: 'auth/method/SOKji',
                counts: {
                  clients: 114,
                  entity_clients: 34,
                  non_entity_clients: 80,
                },
              },
              {
                path: 'auth/method/CEYXo',
                counts: {
                  clients: 10,
                  entity_clients: 6,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/RPjsj',
                counts: {
                  clients: 5,
                  entity_clients: 4,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/dIqPJ',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/wThqG',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/Sa1dO',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/0JVs1',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'Uc0o8',
            namespace_path: 'namespace16/',
            counts: {
              distinct_entities: 408,
              non_entity_tokens: 743,
              clients: 1151,
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
              {
                path: 'auth/method/wwNCu',
                counts: {
                  clients: 3,
                  entity_clients: 1,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/vv2O6',
                counts: {
                  clients: 48,
                  entity_clients: 32,
                  non_entity_clients: 16,
                },
              },
              {
                path: 'auth/method/zRqUm',
                counts: {
                  clients: 4,
                  entity_clients: 1,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/Yez2v',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/SBBJ2',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/NNSCC',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'R6L40',
            namespace_path: 'namespace2/',
            counts: {
              distinct_entities: 292,
              non_entity_tokens: 736,
              clients: 1028,
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
              {
                path: 'auth/method/cTsw9',
                counts: {
                  clients: 291,
                  entity_clients: 88,
                  non_entity_clients: 203,
                },
              },
              {
                path: 'auth/method/3KTWZ',
                counts: {
                  clients: 14,
                  entity_clients: 7,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/Douf5',
                counts: {
                  clients: 8,
                  entity_clients: 1,
                  non_entity_clients: 7,
                },
              },
              {
                path: 'auth/method/30eez',
                counts: {
                  clients: 19,
                  entity_clients: 11,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/xSSJz',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/pR3x7',
                counts: {
                  clients: 6,
                  entity_clients: 5,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'Rqa3W',
            namespace_path: 'namespace13/',
            counts: {
              distinct_entities: 160,
              non_entity_tokens: 803,
              clients: 963,
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
              {
                path: 'auth/method/oyWlP',
                counts: {
                  clients: 6,
                  entity_clients: 4,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/z7Uka',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/ftNn7',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/pvdQ7',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/DsnIn',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/E1YLg',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'MSgZE',
            namespace_path: 'namespace7/',
            counts: {
              distinct_entities: 201,
              non_entity_tokens: 657,
              clients: 858,
            },
            mounts: [
              {
                path: 'auth/method/gD50V',
                counts: {
                  clients: 246,
                  entity_clients: 73,
                  non_entity_clients: 173,
                },
              },
              {
                path: 'auth/method/iJRmf',
                counts: {
                  clients: 525,
                  entity_clients: 19,
                  non_entity_clients: 506,
                },
              },
              {
                path: 'auth/method/GrNjy',
                counts: {
                  clients: 45,
                  entity_clients: 33,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/r0Uw3',
                counts: {
                  clients: 23,
                  entity_clients: 15,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/k2lQG',
                counts: {
                  clients: 12,
                  entity_clients: 9,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/hJxto',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/vtDck',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/1CenH',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/M47Ey',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/gVT0t',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'kxU4t',
            namespace_path: 'namespacelonglonglong3/',
            counts: {
              distinct_entities: 742,
              non_entity_tokens: 26,
              clients: 768,
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
              {
                path: 'auth/method/9V6aN',
                counts: {
                  clients: 51,
                  entity_clients: 21,
                  non_entity_clients: 30,
                },
              },
              {
                path: 'auth/method/USYOd',
                counts: {
                  clients: 7,
                  entity_clients: 2,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/8pfWr',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/0L511',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/6d0rw',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/ECHpZ',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: '5xKya',
            namespace_path: 'namespace15/',
            counts: {
              distinct_entities: 663,
              non_entity_tokens: 19,
              clients: 682,
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
              {
                path: 'auth/method/WHW73',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/KcO46',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/y2vSv',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/VNy4X',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/cEDV9',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/CZTaj',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: '5KxXA',
            namespace_path: 'namespace18anotherlong/',
            counts: {
              distinct_entities: 470,
              non_entity_tokens: 196,
              clients: 666,
            },
            mounts: [
              {
                path: 'auth/method/GkDM1',
                counts: {
                  clients: 101,
                  entity_clients: 77,
                  non_entity_clients: 24,
                },
              },
              {
                path: 'auth/method/7deLa',
                counts: {
                  clients: 329,
                  entity_clients: 177,
                  non_entity_clients: 152,
                },
              },
              {
                path: 'auth/method/Ash3Y',
                counts: {
                  clients: 126,
                  entity_clients: 66,
                  non_entity_clients: 60,
                },
              },
              {
                path: 'auth/method/doKJ0',
                counts: {
                  clients: 89,
                  entity_clients: 25,
                  non_entity_clients: 64,
                },
              },
              {
                path: 'auth/method/9Irmo',
                counts: {
                  clients: 7,
                  entity_clients: 6,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/jdYx5',
                counts: {
                  clients: 10,
                  entity_clients: 6,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/sYe2h',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/Z5F36',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/O0cuK',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/0clSt',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'AAidI',
            namespace_path: 'namespace20/',
            counts: {
              distinct_entities: 429,
              non_entity_tokens: 60,
              clients: 489,
            },
            mounts: [
              {
                path: 'auth/method/zolCO',
                counts: {
                  clients: 351,
                  entity_clients: 170,
                  non_entity_clients: 181,
                },
              },
              {
                path: 'auth/method/6p3g4',
                counts: {
                  clients: 81,
                  entity_clients: 10,
                  non_entity_clients: 71,
                },
              },
              {
                path: 'auth/method/iKOdR',
                counts: {
                  clients: 13,
                  entity_clients: 4,
                  non_entity_clients: 9,
                },
              },
              {
                path: 'auth/method/brnKt',
                counts: {
                  clients: 21,
                  entity_clients: 15,
                  non_entity_clients: 6,
                },
              },
              {
                path: 'auth/method/qK3rr',
                counts: {
                  clients: 17,
                  entity_clients: 13,
                  non_entity_clients: 4,
                },
              },
              {
                path: 'auth/method/DmAuN',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/krE4t',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/sFrWK',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/bQg4l',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/Jaw0k',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'BCl56',
            namespace_path: 'namespace8/',
            counts: {
              distinct_entities: 61,
              non_entity_tokens: 201,
              clients: 262,
            },
            mounts: [
              {
                path: 'auth/method/LpVqc',
                counts: {
                  clients: 104,
                  entity_clients: 81,
                  non_entity_clients: 23,
                },
              },
              {
                path: 'auth/method/VFHO6',
                counts: {
                  clients: 31,
                  entity_clients: 23,
                  non_entity_clients: 8,
                },
              },
              {
                path: 'auth/method/utu0r',
                counts: {
                  clients: 50,
                  entity_clients: 20,
                  non_entity_clients: 30,
                },
              },
              {
                path: 'auth/method/xikiW',
                counts: {
                  clients: 27,
                  entity_clients: 15,
                  non_entity_clients: 12,
                },
              },
              {
                path: 'auth/method/uPSo6',
                counts: {
                  clients: 16,
                  entity_clients: 15,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/Z8fpo',
                counts: {
                  clients: 24,
                  entity_clients: 23,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/5BBm7',
                counts: {
                  clients: 4,
                  entity_clients: 3,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/Eyxkz',
                counts: {
                  clients: 3,
                  entity_clients: 2,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/QBC0w',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/8MdGr',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
            ],
          },
          {
            namespace_id: 'yYNw2',
            namespace_path: 'namespace19/',
            counts: {
              distinct_entities: 165,
              non_entity_tokens: 85,
              clients: 250,
            },
            mounts: [
              {
                path: 'auth/method/zD8lQ',
                counts: {
                  clients: 63,
                  entity_clients: 28,
                  non_entity_clients: 35,
                },
              },
              {
                path: 'auth/method/Dl96I',
                counts: {
                  clients: 131,
                  entity_clients: 91,
                  non_entity_clients: 40,
                },
              },
              {
                path: 'auth/method/ElIse',
                counts: {
                  clients: 13,
                  entity_clients: 3,
                  non_entity_clients: 10,
                },
              },
              {
                path: 'auth/method/AXzhE',
                counts: {
                  clients: 7,
                  entity_clients: 4,
                  non_entity_clients: 3,
                },
              },
              {
                path: 'auth/method/cNuC6',
                counts: {
                  clients: 4,
                  entity_clients: 2,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/gXtbE',
                counts: {
                  clients: 26,
                  entity_clients: 15,
                  non_entity_clients: 11,
                },
              },
              {
                path: 'auth/method/PptIE',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/QILdh',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/cClAS',
                counts: {
                  clients: 2,
                  entity_clients: 1,
                  non_entity_clients: 1,
                },
              },
              {
                path: 'auth/method/YYm3v',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
            ],
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 67,
              non_entity_tokens: 9,
              clients: 76,
            },
            mounts: [
              {
                path: 'auth/method/koO6h',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/iF9oZ',
                counts: {
                  clients: 65,
                  entity_clients: 60,
                  non_entity_clients: 5,
                },
              },
              {
                path: 'auth/method/N6guZ',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/h2CxN',
                counts: {
                  clients: 3,
                  entity_clients: 1,
                  non_entity_clients: 2,
                },
              },
              {
                path: 'auth/method/pA5pU',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/xbqJh',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/m7vOo',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/lULhW',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/hB9qn',
                counts: {
                  clients: 1,
                  entity_clients: 1,
                  non_entity_clients: 0,
                },
              },
              {
                path: 'auth/method/RIEKI',
                counts: {
                  clients: 1,
                  entity_clients: 1,
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
              clients: 3611,
              entity_clients: 671,
              non_entity_clients: 2940,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 682,
                  entity_clients: 663,
                  non_entity_clients: 19,
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
                  {
                    path: 'auth/method/WHW73',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/KcO46',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/y2vSv',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/VNy4X',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/cEDV9',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/CZTaj',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace15/',
                path: '',
              },
              {
                counts: {
                  clients: 1413,
                  entity_clients: 996,
                  non_entity_clients: 417,
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
                  {
                    path: 'auth/method/jdRjF',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/yyBoC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/WLxYp',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SNM6V',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/vNHtH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/EqmlO',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
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
                    {
                      path: 'auth/method/WHW73',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/KcO46',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/y2vSv',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/VNy4X',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/cEDV9',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/CZTaj',
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
                    {
                      path: 'auth/method/jdRjF',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/yyBoC',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/WLxYp',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SNM6V',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vNHtH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-02-01T08:00:00.000Z',
            counts: {
              clients: 3697,
              entity_clients: 1664,
              non_entity_clients: 2033,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 768,
                  entity_clients: 742,
                  non_entity_clients: 26,
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
                  {
                    path: 'auth/method/9V6aN',
                    counts: {
                      clients: 51,
                      entity_clients: 21,
                      non_entity_clients: 30,
                    },
                  },
                  {
                    path: 'auth/method/USYOd',
                    counts: {
                      clients: 7,
                      entity_clients: 2,
                      non_entity_clients: 5,
                    },
                  },
                  {
                    path: 'auth/method/8pfWr',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/0L511',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/6d0rw',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/ECHpZ',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespacelonglonglong3/',
                path: '',
              },
              {
                counts: {
                  clients: 1413,
                  entity_clients: 996,
                  non_entity_clients: 417,
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
                  {
                    path: 'auth/method/jdRjF',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/yyBoC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/WLxYp',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SNM6V',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/vNHtH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/EqmlO',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 86,
                entity_clients: 42,
                non_entity_clients: 44,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 4,
                        entity_clients: 2,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/9V6aN',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 30,
                      },
                    },
                    {
                      path: 'auth/method/USYOd',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 5,
                      },
                    },
                    {
                      path: 'auth/method/8pfWr',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/0L511',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/6d0rw',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ECHpZ',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/jdRjF',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/yyBoC',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/WLxYp',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SNM6V',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vNHtH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                  id: 'namespacelonglonglong4/',
                  path: '',
                },
              ],
            },
          },
          {
            timestamp: '2022-03-01T08:00:00.000Z',
            counts: {
              clients: 3861,
              entity_clients: 1251,
              non_entity_clients: 2610,
            },
            namespaces: [
              {
                counts: {
                  clients: 1413,
                  entity_clients: 996,
                  non_entity_clients: 417,
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
                  {
                    path: 'auth/method/jdRjF',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/yyBoC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/WLxYp',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SNM6V',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/vNHtH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/EqmlO',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  clients: 1268,
                  entity_clients: 563,
                  non_entity_clients: 705,
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
                  {
                    path: 'auth/method/LY3am',
                    counts: {
                      clients: 8,
                      entity_clients: 5,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/McQ4X',
                    counts: {
                      clients: 6,
                      entity_clients: 4,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/NpjhH',
                    counts: {
                      clients: 28,
                      entity_clients: 4,
                      non_entity_clients: 24,
                    },
                  },
                  {
                    path: 'auth/method/ToKO8',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/wfApH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/L9uWV',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  clients: 1180,
                  entity_clients: 925,
                  non_entity_clients: 255,
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
                  {
                    path: 'auth/method/9v04G',
                    counts: {
                      clients: 33,
                      entity_clients: 21,
                      non_entity_clients: 12,
                    },
                  },
                  {
                    path: 'auth/method/6hAlO',
                    counts: {
                      clients: 11,
                      entity_clients: 10,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/ydSdP',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/i0CTY',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/nevwU',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/k2jYC',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 164,
                entity_clients: 42,
                non_entity_clients: 122,
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
                    {
                      path: 'auth/method/jdRjF',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/yyBoC',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/WLxYp',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SNM6V',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vNHtH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/LY3am',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/McQ4X',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/NpjhH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 24,
                      },
                    },
                    {
                      path: 'auth/method/ToKO8',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/wfApH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/L9uWV',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/9v04G',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 12,
                      },
                    },
                    {
                      path: 'auth/method/6hAlO',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/ydSdP',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/i0CTY',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/nevwU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/k2jYC',
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
              ],
            },
          },
          {
            timestamp: '2022-04-01T07:00:00.000Z',
            counts: {
              clients: 3870,
              entity_clients: 1399,
              non_entity_clients: 2471,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 1180,
                  entity_clients: 925,
                  non_entity_clients: 255,
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
                  {
                    path: 'auth/method/9v04G',
                    counts: {
                      clients: 33,
                      entity_clients: 21,
                      non_entity_clients: 12,
                    },
                  },
                  {
                    path: 'auth/method/6hAlO',
                    counts: {
                      clients: 11,
                      entity_clients: 10,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/ydSdP',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/i0CTY',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/nevwU',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/k2jYC',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
              {
                counts: {
                  clients: 1174,
                  entity_clients: 935,
                  non_entity_clients: 239,
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
                  {
                    path: 'auth/method/klonh',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/JyhFQ',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/S66CH',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/6pBz3',
                    counts: {
                      clients: 4,
                      entity_clients: 2,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/qHCZa',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/I6OpF',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace9/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 9,
                entity_clients: 5,
                non_entity_clients: 4,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/9v04G',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6hAlO',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ydSdP',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/i0CTY',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/nevwU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/k2jYC',
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
                    {
                      path: 'auth/method/klonh',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/JyhFQ',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/S66CH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6pBz3',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/qHCZa',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/I6OpF',
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
            timestamp: '2022-05-01T07:00:00.000Z',
            counts: {
              clients: 3870,
              entity_clients: 416,
              non_entity_clients: 3454,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 1180,
                  entity_clients: 925,
                  non_entity_clients: 255,
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
                  {
                    path: 'auth/method/9v04G',
                    counts: {
                      clients: 33,
                      entity_clients: 21,
                      non_entity_clients: 12,
                    },
                  },
                  {
                    path: 'auth/method/6hAlO',
                    counts: {
                      clients: 11,
                      entity_clients: 10,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/ydSdP',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/i0CTY',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/nevwU',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/k2jYC',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
              {
                counts: {
                  clients: 1174,
                  entity_clients: 935,
                  non_entity_clients: 239,
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
                  {
                    path: 'auth/method/klonh',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/JyhFQ',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/S66CH',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/6pBz3',
                    counts: {
                      clients: 4,
                      entity_clients: 2,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/qHCZa',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/I6OpF',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
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
                    {
                      path: 'auth/method/9v04G',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6hAlO',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ydSdP',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/i0CTY',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/nevwU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/k2jYC',
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
                    {
                      path: 'auth/method/klonh',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/JyhFQ',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/S66CH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/6pBz3',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/qHCZa',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/I6OpF',
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
            timestamp: '2022-06-01T07:00:00.000Z',
            counts: {
              clients: 3935,
              entity_clients: 1852,
              non_entity_clients: 2083,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 1268,
                  entity_clients: 563,
                  non_entity_clients: 705,
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
                  {
                    path: 'auth/method/LY3am',
                    counts: {
                      clients: 8,
                      entity_clients: 5,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/McQ4X',
                    counts: {
                      clients: 6,
                      entity_clients: 4,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/NpjhH',
                    counts: {
                      clients: 28,
                      entity_clients: 4,
                      non_entity_clients: 24,
                    },
                  },
                  {
                    path: 'auth/method/ToKO8',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/wfApH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/L9uWV',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  clients: 1151,
                  entity_clients: 408,
                  non_entity_clients: 743,
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
                  {
                    path: 'auth/method/wwNCu',
                    counts: {
                      clients: 3,
                      entity_clients: 1,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/vv2O6',
                    counts: {
                      clients: 48,
                      entity_clients: 32,
                      non_entity_clients: 16,
                    },
                  },
                  {
                    path: 'auth/method/zRqUm',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/Yez2v',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SBBJ2',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/NNSCC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace16/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 65,
                entity_clients: 15,
                non_entity_clients: 50,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
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
                    {
                      path: 'auth/method/LY3am',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/McQ4X',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/NpjhH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ToKO8',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/wfApH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/L9uWV',
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
                    {
                      path: 'auth/method/wwNCu',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vv2O6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/zRqUm',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/Yez2v',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SBBJ2',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/NNSCC',
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
              clients: 3935,
              entity_clients: 969,
              non_entity_clients: 2966,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 1268,
                  entity_clients: 563,
                  non_entity_clients: 705,
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
                  {
                    path: 'auth/method/LY3am',
                    counts: {
                      clients: 8,
                      entity_clients: 5,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/McQ4X',
                    counts: {
                      clients: 6,
                      entity_clients: 4,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/NpjhH',
                    counts: {
                      clients: 28,
                      entity_clients: 4,
                      non_entity_clients: 24,
                    },
                  },
                  {
                    path: 'auth/method/ToKO8',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/wfApH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/L9uWV',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  clients: 1151,
                  entity_clients: 408,
                  non_entity_clients: 743,
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
                  {
                    path: 'auth/method/wwNCu',
                    counts: {
                      clients: 3,
                      entity_clients: 1,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/vv2O6',
                    counts: {
                      clients: 48,
                      entity_clients: 32,
                      non_entity_clients: 16,
                    },
                  },
                  {
                    path: 'auth/method/zRqUm',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/Yez2v',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SBBJ2',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/NNSCC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace16/',
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
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
                    {
                      path: 'auth/method/LY3am',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/McQ4X',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/NpjhH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ToKO8',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/wfApH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/L9uWV',
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
                    {
                      path: 'auth/method/wwNCu',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vv2O6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/zRqUm',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/Yez2v',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SBBJ2',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/NNSCC',
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
            timestamp: '2022-08-01T07:00:00.000Z',
            counts: {
              clients: 4058,
              entity_clients: 1170,
              non_entity_clients: 2888,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 1369,
                  entity_clients: 829,
                  non_entity_clients: 540,
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
                  {
                    path: 'auth/method/KGFfV',
                    counts: {
                      clients: 20,
                      entity_clients: 9,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/23AQk',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/PqTWe',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/pPSo1',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/HMu5H',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/xpOk3',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace12/',
                path: '',
              },
              {
                counts: {
                  clients: 1173,
                  entity_clients: 810,
                  non_entity_clients: 363,
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
                  {
                    path: 'auth/method/rezK4',
                    counts: {
                      clients: 119,
                      entity_clients: 74,
                      non_entity_clients: 45,
                    },
                  },
                  {
                    path: 'auth/method/qWNSS',
                    counts: {
                      clients: 106,
                      entity_clients: 42,
                      non_entity_clients: 64,
                    },
                  },
                  {
                    path: 'auth/method/OmQEf',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/PhoAy',
                    counts: {
                      clients: 9,
                      entity_clients: 8,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/aUuyM',
                    counts: {
                      clients: 47,
                      entity_clients: 32,
                      non_entity_clients: 15,
                    },
                  },
                  {
                    path: 'auth/method/kUj1S',
                    counts: {
                      clients: 32,
                      entity_clients: 1,
                      non_entity_clients: 31,
                    },
                  },
                ],
                id: 'namespace6/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 123,
                entity_clients: 6,
                non_entity_clients: 117,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/KGFfV',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 11,
                      },
                    },
                    {
                      path: 'auth/method/23AQk',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PqTWe',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/pPSo1',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/HMu5H',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/xpOk3',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/rezK4',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 45,
                      },
                    },
                    {
                      path: 'auth/method/qWNSS',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 64,
                      },
                    },
                    {
                      path: 'auth/method/OmQEf',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/PhoAy',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/aUuyM',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 15,
                      },
                    },
                    {
                      path: 'auth/method/kUj1S',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 31,
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
            timestamp: '2022-09-01T07:00:00.000Z',
            counts: {
              clients: 5035,
              entity_clients: 2168,
              non_entity_clients: 2867,
            },
            namespaces: [
              {
                counts: {
                  clients: 1413,
                  entity_clients: 996,
                  non_entity_clients: 417,
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
                  {
                    path: 'auth/method/jdRjF',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/yyBoC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/WLxYp',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SNM6V',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/vNHtH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/EqmlO',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  clients: 1268,
                  entity_clients: 563,
                  non_entity_clients: 705,
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
                  {
                    path: 'auth/method/LY3am',
                    counts: {
                      clients: 8,
                      entity_clients: 5,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/McQ4X',
                    counts: {
                      clients: 6,
                      entity_clients: 4,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/NpjhH',
                    counts: {
                      clients: 28,
                      entity_clients: 4,
                      non_entity_clients: 24,
                    },
                  },
                  {
                    path: 'auth/method/ToKO8',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/wfApH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/L9uWV',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
              {
                counts: {
                  clients: 1180,
                  entity_clients: 925,
                  non_entity_clients: 255,
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
                  {
                    path: 'auth/method/9v04G',
                    counts: {
                      clients: 33,
                      entity_clients: 21,
                      non_entity_clients: 12,
                    },
                  },
                  {
                    path: 'auth/method/6hAlO',
                    counts: {
                      clients: 11,
                      entity_clients: 10,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/ydSdP',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/i0CTY',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/nevwU',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/k2jYC',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace10/',
                path: '',
              },
              {
                counts: {
                  clients: 1174,
                  entity_clients: 935,
                  non_entity_clients: 239,
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
                  {
                    path: 'auth/method/klonh',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/JyhFQ',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/S66CH',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/6pBz3',
                    counts: {
                      clients: 4,
                      entity_clients: 2,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/qHCZa',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/I6OpF',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace9/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 977,
                entity_clients: 280,
                non_entity_clients: 697,
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
                    {
                      path: 'auth/method/jdRjF',
                      counts: {
                        clients: 13,
                        entity_clients: 5,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/yyBoC',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/WLxYp',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SNM6V',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vNHtH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/LY3am',
                      counts: {
                        clients: 3,
                        entity_clients: 1,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/McQ4X',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/NpjhH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 24,
                      },
                    },
                    {
                      path: 'auth/method/ToKO8',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/wfApH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/L9uWV',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/9v04G',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 12,
                      },
                    },
                    {
                      path: 'auth/method/6hAlO',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/ydSdP',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/i0CTY',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/nevwU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/k2jYC',
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
                    {
                      path: 'auth/method/klonh',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/JyhFQ',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/S66CH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/6pBz3',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/qHCZa',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/I6OpF',
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
              clients: 5225,
              entity_clients: 859,
              non_entity_clients: 4366,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 1028,
                  entity_clients: 292,
                  non_entity_clients: 736,
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
                  {
                    path: 'auth/method/cTsw9',
                    counts: {
                      clients: 291,
                      entity_clients: 88,
                      non_entity_clients: 203,
                    },
                  },
                  {
                    path: 'auth/method/3KTWZ',
                    counts: {
                      clients: 14,
                      entity_clients: 7,
                      non_entity_clients: 7,
                    },
                  },
                  {
                    path: 'auth/method/Douf5',
                    counts: {
                      clients: 8,
                      entity_clients: 1,
                      non_entity_clients: 7,
                    },
                  },
                  {
                    path: 'auth/method/30eez',
                    counts: {
                      clients: 19,
                      entity_clients: 11,
                      non_entity_clients: 8,
                    },
                  },
                  {
                    path: 'auth/method/xSSJz',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/pR3x7',
                    counts: {
                      clients: 6,
                      entity_clients: 5,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace2/',
                path: '',
              },
              {
                counts: {
                  clients: 1413,
                  entity_clients: 996,
                  non_entity_clients: 417,
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
                  {
                    path: 'auth/method/jdRjF',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/yyBoC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/WLxYp',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SNM6V',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/vNHtH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/EqmlO',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  clients: 1268,
                  entity_clients: 563,
                  non_entity_clients: 705,
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
                  {
                    path: 'auth/method/LY3am',
                    counts: {
                      clients: 8,
                      entity_clients: 5,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/McQ4X',
                    counts: {
                      clients: 6,
                      entity_clients: 4,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/NpjhH',
                    counts: {
                      clients: 28,
                      entity_clients: 4,
                      non_entity_clients: 24,
                    },
                  },
                  {
                    path: 'auth/method/ToKO8',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/wfApH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/L9uWV',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace11/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 190,
                entity_clients: 87,
                non_entity_clients: 103,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/cTsw9',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 203,
                      },
                    },
                    {
                      path: 'auth/method/3KTWZ',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 7,
                      },
                    },
                    {
                      path: 'auth/method/Douf5',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 7,
                      },
                    },
                    {
                      path: 'auth/method/30eez',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 8,
                      },
                    },
                    {
                      path: 'auth/method/xSSJz',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/pR3x7',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/jdRjF',
                      counts: {
                        clients: 2,
                        entity_clients: 1,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/yyBoC',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/WLxYp',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SNM6V',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vNHtH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/LY3am',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/McQ4X',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/NpjhH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 24,
                      },
                    },
                    {
                      path: 'auth/method/ToKO8',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/wfApH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/L9uWV',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
            timestamp: '2022-11-01T07:00:00.000Z',
            counts: {
              clients: 5225,
              entity_clients: 2432,
              non_entity_clients: 2793,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 1028,
                  entity_clients: 292,
                  non_entity_clients: 736,
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
                  {
                    path: 'auth/method/cTsw9',
                    counts: {
                      clients: 291,
                      entity_clients: 88,
                      non_entity_clients: 203,
                    },
                  },
                  {
                    path: 'auth/method/3KTWZ',
                    counts: {
                      clients: 14,
                      entity_clients: 7,
                      non_entity_clients: 7,
                    },
                  },
                  {
                    path: 'auth/method/Douf5',
                    counts: {
                      clients: 8,
                      entity_clients: 1,
                      non_entity_clients: 7,
                    },
                  },
                  {
                    path: 'auth/method/30eez',
                    counts: {
                      clients: 19,
                      entity_clients: 11,
                      non_entity_clients: 8,
                    },
                  },
                  {
                    path: 'auth/method/xSSJz',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/pR3x7',
                    counts: {
                      clients: 6,
                      entity_clients: 5,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace2/',
                path: '',
              },
              {
                counts: {
                  clients: 1413,
                  entity_clients: 996,
                  non_entity_clients: 417,
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
                  {
                    path: 'auth/method/jdRjF',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/yyBoC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/WLxYp',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SNM6V',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/vNHtH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/EqmlO',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  clients: 1268,
                  entity_clients: 563,
                  non_entity_clients: 705,
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
                  {
                    path: 'auth/method/LY3am',
                    counts: {
                      clients: 8,
                      entity_clients: 5,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/McQ4X',
                    counts: {
                      clients: 6,
                      entity_clients: 4,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/NpjhH',
                    counts: {
                      clients: 28,
                      entity_clients: 4,
                      non_entity_clients: 24,
                    },
                  },
                  {
                    path: 'auth/method/ToKO8',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/wfApH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/L9uWV',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace11/',
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: -1,
                        entity_clients: -1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
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
                    {
                      path: 'auth/method/cTsw9',
                      counts: {
                        clients: -3,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/3KTWZ',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/Douf5',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/30eez',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/xSSJz',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/pR3x7',
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
                    {
                      path: 'auth/method/jdRjF',
                      counts: {
                        clients: -2,
                        entity_clients: -1,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/yyBoC',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/WLxYp',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SNM6V',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vNHtH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/LY3am',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/McQ4X',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/NpjhH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ToKO8',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/wfApH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/L9uWV',
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
            timestamp: '2022-12-01T08:00:00.000Z',
            counts: {
              clients: 5261,
              entity_clients: 2152,
              non_entity_clients: 3109,
            },
            namespaces: [
              {
                counts: {
                  clients: 1516,
                  entity_clients: 966,
                  non_entity_clients: 550,
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
                  {
                    path: 'auth/method/7ecN2',
                    counts: {
                      clients: 25,
                      entity_clients: 24,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/AYdDo',
                    counts: {
                      clients: 3,
                      entity_clients: 2,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/kS9h6',
                    counts: {
                      clients: 4,
                      entity_clients: 1,
                      non_entity_clients: 3,
                    },
                  },
                  {
                    path: 'auth/method/dIoMU',
                    counts: {
                      clients: 8,
                      entity_clients: 2,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/eXB1u',
                    counts: {
                      clients: 7,
                      entity_clients: 1,
                      non_entity_clients: 6,
                    },
                  },
                  {
                    path: 'auth/method/SQ8Ty',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                ],
                id: 'namespace17/',
                path: '',
              },
              {
                counts: {
                  clients: 963,
                  entity_clients: 160,
                  non_entity_clients: 803,
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
                  {
                    path: 'auth/method/oyWlP',
                    counts: {
                      clients: 6,
                      entity_clients: 4,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/z7Uka',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/ftNn7',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/pvdQ7',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/DsnIn',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/E1YLg',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace13/',
                path: '',
              },
              {
                counts: {
                  clients: 1413,
                  entity_clients: 996,
                  non_entity_clients: 417,
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
                  {
                    path: 'auth/method/jdRjF',
                    counts: {
                      clients: 5,
                      entity_clients: 3,
                      non_entity_clients: 2,
                    },
                  },
                  {
                    path: 'auth/method/yyBoC',
                    counts: {
                      clients: 2,
                      entity_clients: 1,
                      non_entity_clients: 1,
                    },
                  },
                  {
                    path: 'auth/method/WLxYp',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/SNM6V',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/vNHtH',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/EqmlO',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespacelonglonglong4/',
                path: '',
              },
              {
                counts: {
                  clients: 1369,
                  entity_clients: 829,
                  non_entity_clients: 540,
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
                  {
                    path: 'auth/method/KGFfV',
                    counts: {
                      clients: 20,
                      entity_clients: 9,
                      non_entity_clients: 11,
                    },
                  },
                  {
                    path: 'auth/method/23AQk',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/PqTWe',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/pPSo1',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/HMu5H',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                  {
                    path: 'auth/method/xpOk3',
                    counts: {
                      clients: 1,
                      entity_clients: 1,
                      non_entity_clients: 0,
                    },
                  },
                ],
                id: 'namespace12/',
                path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 36,
                entity_clients: 35,
                non_entity_clients: 1,
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
                    {
                      path: 'auth/method/7ecN2',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/AYdDo',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/kS9h6',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 3,
                      },
                    },
                    {
                      path: 'auth/method/dIoMU',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/eXB1u',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 6,
                      },
                    },
                    {
                      path: 'auth/method/SQ8Ty',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
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
                    {
                      path: 'auth/method/oyWlP',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/z7Uka',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/ftNn7',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/pvdQ7',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/DsnIn',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/jdRjF',
                      counts: {
                        clients: 1,
                        entity_clients: 1,
                        non_entity_clients: 2,
                      },
                    },
                    {
                      path: 'auth/method/yyBoC',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 1,
                      },
                    },
                    {
                      path: 'auth/method/WLxYp',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/SNM6V',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/vNHtH',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
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
                    {
                      path: 'auth/method/KGFfV',
                      counts: {
                        clients: -1,
                        entity_clients: -1,
                        non_entity_clients: 11,
                      },
                    },
                    {
                      path: 'auth/method/23AQk',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/PqTWe',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/pPSo1',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/HMu5H',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                    {
                      path: 'auth/method/xpOk3',
                      counts: {
                        clients: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                      },
                    },
                  ],
                  id: 'namespace12/',
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
