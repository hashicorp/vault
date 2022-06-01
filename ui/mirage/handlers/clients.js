import {
  formatISO,
  isAfter,
  isBefore,
  sub,
  isSameMonth,
  startOfMonth,
  endOfMonth,
  addMonths,
  subMonths,
  differenceInCalendarMonths,
} from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import formatRFC3339 from 'date-fns/formatRFC3339';

const NEW_DATE = new Date();
const COUNTS_START = subMonths(NEW_DATE, 12); // pretend vault user started cluster 1 year ago

// for testing, we're in the middle of a license/billing period
const LICENSE_START = startOfMonth(subMonths(NEW_DATE, 6));
const LICENSE_END = endOfMonth(addMonths(NEW_DATE, 6));

// upgrade happened 1 month after license start
const UPGRADE_DATE = addMonths(LICENSE_START, 1);
// Oldest to newest
const MOCK_MONTHLY_DATA = [
  {
    timestamp: formatISO(UPGRADE_DATE),
    counts: {
      distinct_entities: 0,
      entity_clients: 10433,
      non_entity_tokens: 0,
      non_entity_clients: 7555,
      clients: 17988,
    },
    namespaces: [
      {
        namespace_id: 'PU6JB',
        namespace_path: 'test-ns-2/',
        counts: {
          distinct_entities: 0,
          entity_clients: 3458,
          non_entity_tokens: 0,
          non_entity_clients: 1631,
          clients: 5089,
        },
        mounts: [
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 948,
              non_entity_tokens: 0,
              non_entity_clients: 714,
              clients: 1662,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 899,
              non_entity_tokens: 0,
              non_entity_clients: 301,
              clients: 1200,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 692,
              non_entity_tokens: 0,
              non_entity_clients: 474,
              clients: 1166,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 919,
              non_entity_tokens: 0,
              non_entity_clients: 142,
              clients: 1061,
            },
          },
        ],
      },
      {
        namespace_id: '3lq5r',
        namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2428,
          non_entity_tokens: 0,
          non_entity_clients: 1841,
          clients: 4269,
        },
        mounts: [
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 969,
              non_entity_tokens: 0,
              non_entity_clients: 396,
              clients: 1365,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 794,
              non_entity_tokens: 0,
              non_entity_clients: 501,
              clients: 1295,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 289,
              non_entity_tokens: 0,
              non_entity_clients: 666,
              clients: 955,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 376,
              non_entity_tokens: 0,
              non_entity_clients: 278,
              clients: 654,
            },
          },
        ],
      },
      {
        namespace_id: 'sJRLj',
        namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2384,
          non_entity_tokens: 0,
          non_entity_clients: 1278,
          clients: 3662,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 853,
              non_entity_tokens: 0,
              non_entity_clients: 553,
              clients: 1406,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 677,
              non_entity_tokens: 0,
              non_entity_clients: 182,
              clients: 859,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 582,
              non_entity_tokens: 0,
              non_entity_clients: 175,
              clients: 757,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 272,
              non_entity_tokens: 0,
              non_entity_clients: 368,
              clients: 640,
            },
          },
        ],
      },
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 0,
          entity_clients: 943,
          non_entity_tokens: 0,
          non_entity_clients: 1595,
          clients: 2538,
        },
        mounts: [
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 318,
              non_entity_tokens: 0,
              non_entity_clients: 735,
              clients: 1053,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 362,
              non_entity_tokens: 0,
              non_entity_clients: 415,
              clients: 777,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 158,
              non_entity_tokens: 0,
              non_entity_clients: 325,
              clients: 483,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 105,
              non_entity_tokens: 0,
              non_entity_clients: 120,
              clients: 225,
            },
          },
        ],
      },
      {
        namespace_id: 'opmJ1',
        namespace_path: 'test-ns-1/',
        counts: {
          distinct_entities: 0,
          entity_clients: 1220,
          non_entity_tokens: 0,
          non_entity_clients: 1210,
          clients: 2430,
        },
        mounts: [
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 697,
              non_entity_tokens: 0,
              non_entity_clients: 516,
              clients: 1213,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 154,
              non_entity_tokens: 0,
              non_entity_clients: 480,
              clients: 634,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 223,
              non_entity_tokens: 0,
              non_entity_clients: 97,
              clients: 320,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 146,
              non_entity_tokens: 0,
              non_entity_clients: 117,
              clients: 263,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 0,
        entity_clients: 5032,
        non_entity_tokens: 0,
        non_entity_clients: 2888,
        clients: 7920,
      },
      namespaces: [
        {
          namespace_id: 'sJRLj',
          namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1907,
            non_entity_tokens: 0,
            non_entity_clients: 354,
            clients: 2261,
          },
          mounts: [
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 753,
                non_entity_tokens: 0,
                non_entity_clients: 138,
                clients: 891,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 516,
                non_entity_tokens: 0,
                non_entity_clients: 91,
                clients: 607,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 474,
                non_entity_tokens: 0,
                non_entity_clients: 1,
                clients: 475,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 164,
                non_entity_tokens: 0,
                non_entity_clients: 124,
                clients: 288,
              },
            },
          ],
        },
        {
          namespace_id: '3lq5r',
          namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 843,
            non_entity_tokens: 0,
            non_entity_clients: 748,
            clients: 1591,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 597,
                non_entity_tokens: 0,
                non_entity_clients: 369,
                clients: 966,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 185,
                non_entity_tokens: 0,
                non_entity_clients: 156,
                clients: 341,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 59,
                non_entity_tokens: 0,
                non_entity_clients: 223,
                clients: 282,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 2,
                non_entity_tokens: 0,
                non_entity_clients: 0,
                clients: 2,
              },
            },
          ],
        },
        {
          namespace_id: 'PU6JB',
          namespace_path: 'test-ns-2/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1291,
            non_entity_tokens: 0,
            non_entity_clients: 268,
            clients: 1559,
          },
          mounts: [
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 518,
                non_entity_tokens: 0,
                non_entity_clients: 78,
                clients: 596,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 291,
                non_entity_tokens: 0,
                non_entity_clients: 85,
                clients: 376,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 276,
                non_entity_tokens: 0,
                non_entity_clients: 59,
                clients: 335,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 206,
                non_entity_tokens: 0,
                non_entity_clients: 46,
                clients: 252,
              },
            },
          ],
        },
        {
          namespace_id: 'opmJ1',
          namespace_path: 'test-ns-1/',
          counts: {
            distinct_entities: 0,
            entity_clients: 663,
            non_entity_tokens: 0,
            non_entity_clients: 778,
            clients: 1441,
          },
          mounts: [
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 478,
                non_entity_tokens: 0,
                non_entity_clients: 309,
                clients: 787,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 129,
                non_entity_tokens: 0,
                non_entity_clients: 313,
                clients: 442,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 32,
                non_entity_tokens: 0,
                non_entity_clients: 79,
                clients: 111,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 24,
                non_entity_tokens: 0,
                non_entity_clients: 77,
                clients: 101,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 0,
            entity_clients: 328,
            non_entity_tokens: 0,
            non_entity_clients: 740,
            clients: 1068,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 143,
                non_entity_tokens: 0,
                non_entity_clients: 273,
                clients: 416,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 136,
                non_entity_tokens: 0,
                non_entity_clients: 142,
                clients: 278,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 19,
                non_entity_tokens: 0,
                non_entity_clients: 216,
                clients: 235,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 30,
                non_entity_tokens: 0,
                non_entity_clients: 109,
                clients: 139,
              },
            },
          ],
        },
      ],
    },
  },
  {
    timestamp: formatISO(addMonths(UPGRADE_DATE, 1)),
    counts: {
      distinct_entities: 0,
      entity_clients: 10285,
      non_entity_tokens: 0,
      non_entity_clients: 10425,
      clients: 20710,
    },
    namespaces: [
      {
        namespace_id: '3lq5r',
        namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2335,
          non_entity_tokens: 0,
          non_entity_clients: 2644,
          clients: 4979,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 939,
              non_entity_tokens: 0,
              non_entity_clients: 649,
              clients: 1588,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 461,
              non_entity_tokens: 0,
              non_entity_clients: 870,
              clients: 1331,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 172,
              non_entity_tokens: 0,
              non_entity_clients: 990,
              clients: 1162,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 763,
              non_entity_tokens: 0,
              non_entity_clients: 135,
              clients: 898,
            },
          },
        ],
      },
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 0,
          entity_clients: 2054,
          non_entity_tokens: 0,
          non_entity_clients: 2747,
          clients: 4801,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 540,
              non_entity_tokens: 0,
              non_entity_clients: 941,
              clients: 1481,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 722,
              non_entity_tokens: 0,
              non_entity_clients: 507,
              clients: 1229,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 611,
              non_entity_tokens: 0,
              non_entity_clients: 520,
              clients: 1131,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 181,
              non_entity_tokens: 0,
              non_entity_clients: 779,
              clients: 960,
            },
          },
        ],
      },
      {
        namespace_id: 'PU6JB',
        namespace_path: 'test-ns-2/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2788,
          non_entity_tokens: 0,
          non_entity_clients: 1720,
          clients: 4508,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 643,
              non_entity_tokens: 0,
              non_entity_clients: 814,
              clients: 1457,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 811,
              non_entity_tokens: 0,
              non_entity_clients: 385,
              clients: 1196,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 932,
              non_entity_tokens: 0,
              non_entity_clients: 72,
              clients: 1004,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 402,
              non_entity_tokens: 0,
              non_entity_clients: 449,
              clients: 851,
            },
          },
        ],
      },
      {
        namespace_id: 'sJRLj',
        namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 1162,
          non_entity_tokens: 0,
          non_entity_clients: 2187,
          clients: 3349,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 483,
              non_entity_tokens: 0,
              non_entity_clients: 839,
              clients: 1322,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 373,
              non_entity_tokens: 0,
              non_entity_clients: 858,
              clients: 1231,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 271,
              non_entity_tokens: 0,
              non_entity_clients: 154,
              clients: 425,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 35,
              non_entity_tokens: 0,
              non_entity_clients: 336,
              clients: 371,
            },
          },
        ],
      },
      {
        namespace_id: 'opmJ1',
        namespace_path: 'test-ns-1/',
        counts: {
          distinct_entities: 0,
          entity_clients: 1946,
          non_entity_tokens: 0,
          non_entity_clients: 1127,
          clients: 3073,
        },
        mounts: [
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 322,
              non_entity_tokens: 0,
              non_entity_clients: 537,
              clients: 859,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 685,
              non_entity_tokens: 0,
              non_entity_clients: 132,
              clients: 817,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 321,
              non_entity_tokens: 0,
              non_entity_clients: 385,
              clients: 706,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 618,
              non_entity_tokens: 0,
              non_entity_clients: 73,
              clients: 691,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 0,
        entity_clients: 5315,
        non_entity_tokens: 0,
        non_entity_clients: 5724,
        clients: 11039,
      },
      namespaces: [
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 0,
            entity_clients: 1253,
            non_entity_tokens: 0,
            non_entity_clients: 1529,
            clients: 2782,
          },
          mounts: [
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 452,
                non_entity_tokens: 0,
                non_entity_clients: 433,
                clients: 885,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 134,
                non_entity_tokens: 0,
                non_entity_clients: 732,
                clients: 866,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 472,
                non_entity_tokens: 0,
                non_entity_clients: 361,
                clients: 833,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 195,
                non_entity_tokens: 0,
                non_entity_clients: 3,
                clients: 198,
              },
            },
          ],
        },
        {
          namespace_id: '3lq5r',
          namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1032,
            non_entity_tokens: 0,
            non_entity_clients: 1652,
            clients: 2684,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 93,
                non_entity_tokens: 0,
                non_entity_clients: 849,
                clients: 942,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 692,
                non_entity_tokens: 0,
                non_entity_clients: 117,
                clients: 809,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 159,
                non_entity_tokens: 0,
                non_entity_clients: 596,
                clients: 755,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 88,
                non_entity_tokens: 0,
                non_entity_clients: 90,
                clients: 178,
              },
            },
          ],
        },
        {
          namespace_id: 'opmJ1',
          namespace_path: 'test-ns-1/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1482,
            non_entity_tokens: 0,
            non_entity_clients: 742,
            clients: 2224,
          },
          mounts: [
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 278,
                non_entity_tokens: 0,
                non_entity_clients: 372,
                clients: 650,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 546,
                non_entity_tokens: 0,
                non_entity_clients: 52,
                clients: 598,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 313,
                non_entity_tokens: 0,
                non_entity_clients: 264,
                clients: 577,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 345,
                non_entity_tokens: 0,
                non_entity_clients: 54,
                clients: 399,
              },
            },
          ],
        },
        {
          namespace_id: 'PU6JB',
          namespace_path: 'test-ns-2/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1090,
            non_entity_tokens: 0,
            non_entity_clients: 600,
            clients: 1690,
          },
          mounts: [
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 632,
                non_entity_tokens: 0,
                non_entity_clients: 7,
                clients: 639,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 180,
                non_entity_tokens: 0,
                non_entity_clients: 317,
                clients: 497,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 214,
                non_entity_tokens: 0,
                non_entity_clients: 239,
                clients: 453,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 64,
                non_entity_tokens: 0,
                non_entity_clients: 37,
                clients: 101,
              },
            },
          ],
        },
        {
          namespace_id: 'sJRLj',
          namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 458,
            non_entity_tokens: 0,
            non_entity_clients: 1201,
            clients: 1659,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 134,
                non_entity_tokens: 0,
                non_entity_clients: 827,
                clients: 961,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 232,
                non_entity_tokens: 0,
                non_entity_clients: 112,
                clients: 344,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 4,
                non_entity_tokens: 0,
                non_entity_clients: 193,
                clients: 197,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 88,
                non_entity_tokens: 0,
                non_entity_clients: 69,
                clients: 157,
              },
            },
          ],
        },
      ],
    },
  },
  {
    timestamp: formatISO(addMonths(UPGRADE_DATE, 2)),
    counts: {
      distinct_entities: 0,
      entity_clients: 9721,
      non_entity_tokens: 0,
      non_entity_clients: 11472,
      clients: 21193,
    },
    namespaces: [
      {
        namespace_id: 'sJRLj',
        namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2321,
          non_entity_tokens: 0,
          non_entity_clients: 2864,
          clients: 5185,
        },
        mounts: [
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 616,
              non_entity_tokens: 0,
              non_entity_clients: 940,
              clients: 1556,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 965,
              non_entity_tokens: 0,
              non_entity_clients: 393,
              clients: 1358,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 587,
              non_entity_tokens: 0,
              non_entity_clients: 724,
              clients: 1311,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 153,
              non_entity_tokens: 0,
              non_entity_clients: 807,
              clients: 960,
            },
          },
        ],
      },
      {
        namespace_id: 'PU6JB',
        namespace_path: 'test-ns-2/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2711,
          non_entity_tokens: 0,
          non_entity_clients: 1883,
          clients: 4594,
        },
        mounts: [
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 619,
              non_entity_tokens: 0,
              non_entity_clients: 925,
              clients: 1544,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 661,
              non_entity_tokens: 0,
              non_entity_clients: 814,
              clients: 1475,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 954,
              non_entity_tokens: 0,
              non_entity_clients: 7,
              clients: 961,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 477,
              non_entity_tokens: 0,
              non_entity_clients: 137,
              clients: 614,
            },
          },
        ],
      },
      {
        namespace_id: '3lq5r',
        namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 1426,
          non_entity_tokens: 0,
          non_entity_clients: 2978,
          clients: 4404,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 780,
              non_entity_tokens: 0,
              non_entity_clients: 696,
              clients: 1476,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 369,
              non_entity_tokens: 0,
              non_entity_clients: 977,
              clients: 1346,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 200,
              non_entity_tokens: 0,
              non_entity_clients: 753,
              clients: 953,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 77,
              non_entity_tokens: 0,
              non_entity_clients: 552,
              clients: 629,
            },
          },
        ],
      },
      {
        namespace_id: 'opmJ1',
        namespace_path: 'test-ns-1/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2213,
          non_entity_tokens: 0,
          non_entity_clients: 1851,
          clients: 4064,
        },
        mounts: [
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 610,
              non_entity_tokens: 0,
              non_entity_clients: 893,
              clients: 1503,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 957,
              non_entity_tokens: 0,
              non_entity_clients: 136,
              clients: 1093,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 262,
              non_entity_tokens: 0,
              non_entity_clients: 605,
              clients: 867,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 384,
              non_entity_tokens: 0,
              non_entity_clients: 217,
              clients: 601,
            },
          },
        ],
      },
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 0,
          entity_clients: 1050,
          non_entity_tokens: 0,
          non_entity_clients: 1896,
          clients: 2946,
        },
        mounts: [
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 74,
              non_entity_tokens: 0,
              non_entity_clients: 978,
              clients: 1052,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 238,
              non_entity_tokens: 0,
              non_entity_clients: 530,
              clients: 768,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 367,
              non_entity_tokens: 0,
              non_entity_clients: 199,
              clients: 566,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 371,
              non_entity_tokens: 0,
              non_entity_clients: 189,
              clients: 560,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 0,
        entity_clients: 4637,
        non_entity_tokens: 0,
        non_entity_clients: 5789,
        clients: 10426,
      },
      namespaces: [
        {
          namespace_id: 'sJRLj',
          namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1170,
            non_entity_tokens: 0,
            non_entity_clients: 1525,
            clients: 2695,
          },
          mounts: [
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 511,
                non_entity_tokens: 0,
                non_entity_clients: 720,
                clients: 1231,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 76,
                non_entity_tokens: 0,
                non_entity_clients: 479,
                clients: 555,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 431,
                non_entity_tokens: 0,
                non_entity_clients: 72,
                clients: 503,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 152,
                non_entity_tokens: 0,
                non_entity_clients: 254,
                clients: 406,
              },
            },
          ],
        },
        {
          namespace_id: 'opmJ1',
          namespace_path: 'test-ns-1/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1482,
            non_entity_tokens: 0,
            non_entity_clients: 1129,
            clients: 2611,
          },
          mounts: [
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 553,
                non_entity_tokens: 0,
                non_entity_clients: 884,
                clients: 1437,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 664,
                non_entity_tokens: 0,
                non_entity_clients: 11,
                clients: 675,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 153,
                non_entity_tokens: 0,
                non_entity_clients: 230,
                clients: 383,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 112,
                non_entity_tokens: 0,
                non_entity_clients: 4,
                clients: 116,
              },
            },
          ],
        },
        {
          namespace_id: 'PU6JB',
          namespace_path: 'test-ns-2/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1067,
            non_entity_tokens: 0,
            non_entity_clients: 968,
            clients: 2035,
          },
          mounts: [
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_tokens: 0,
                non_entity_clients: 865,
                clients: 865,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 659,
                non_entity_tokens: 0,
                non_entity_clients: 4,
                clients: 663,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 292,
                non_entity_tokens: 0,
                non_entity_clients: 40,
                clients: 332,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 116,
                non_entity_tokens: 0,
                non_entity_clients: 59,
                clients: 175,
              },
            },
          ],
        },
        {
          namespace_id: '3lq5r',
          namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 547,
            non_entity_tokens: 0,
            non_entity_clients: 1270,
            clients: 1817,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 36,
                non_entity_tokens: 0,
                non_entity_clients: 546,
                clients: 582,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 61,
                non_entity_tokens: 0,
                non_entity_clients: 467,
                clients: 528,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 286,
                non_entity_tokens: 0,
                non_entity_clients: 216,
                clients: 502,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 164,
                non_entity_tokens: 0,
                non_entity_clients: 41,
                clients: 205,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 0,
            entity_clients: 371,
            non_entity_tokens: 0,
            non_entity_clients: 897,
            clients: 1268,
          },
          mounts: [
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 152,
                non_entity_tokens: 0,
                non_entity_clients: 490,
                clients: 642,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 70,
                non_entity_tokens: 0,
                non_entity_clients: 160,
                clients: 230,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 84,
                non_entity_tokens: 0,
                non_entity_clients: 141,
                clients: 225,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 65,
                non_entity_tokens: 0,
                non_entity_clients: 106,
                clients: 171,
              },
            },
          ],
        },
      ],
    },
  },
  {
    counts: {
      distinct_entities: 0,
      entity_clients: 10873,
      non_entity_tokens: 0,
      non_entity_clients: 9343,
      clients: 20216,
    },
    namespaces: [
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 0,
          entity_clients: 1303,
          non_entity_tokens: 0,
          non_entity_clients: 3388,
          clients: 4691,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 721,
              non_entity_tokens: 0,
              non_entity_clients: 980,
              clients: 1701,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 377,
              non_entity_tokens: 0,
              non_entity_clients: 838,
              clients: 1215,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 127,
              non_entity_tokens: 0,
              non_entity_clients: 877,
              clients: 1004,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 78,
              non_entity_tokens: 0,
              non_entity_clients: 693,
              clients: 771,
            },
          },
        ],
      },
      {
        namespace_id: 'opmJ1',
        namespace_path: 'test-ns-1/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2404,
          non_entity_tokens: 0,
          non_entity_clients: 2085,
          clients: 4489,
        },
        mounts: [
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 830,
              non_entity_tokens: 0,
              non_entity_clients: 779,
              clients: 1609,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 926,
              non_entity_tokens: 0,
              non_entity_clients: 311,
              clients: 1237,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 82,
              non_entity_tokens: 0,
              non_entity_clients: 896,
              clients: 978,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 566,
              non_entity_tokens: 0,
              non_entity_clients: 99,
              clients: 665,
            },
          },
        ],
      },
      {
        namespace_id: '3lq5r',
        namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 3076,
          non_entity_tokens: 0,
          non_entity_clients: 1396,
          clients: 4472,
        },
        mounts: [
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 874,
              non_entity_tokens: 0,
              non_entity_clients: 601,
              clients: 1475,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 921,
              non_entity_tokens: 0,
              non_entity_clients: 428,
              clients: 1349,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 885,
              non_entity_tokens: 0,
              non_entity_clients: 204,
              clients: 1089,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 396,
              non_entity_tokens: 0,
              non_entity_clients: 163,
              clients: 559,
            },
          },
        ],
      },
      {
        namespace_id: 'sJRLj',
        namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2298,
          non_entity_tokens: 0,
          non_entity_clients: 1632,
          clients: 3930,
        },
        mounts: [
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 858,
              non_entity_tokens: 0,
              non_entity_clients: 663,
              clients: 1521,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 669,
              non_entity_tokens: 0,
              non_entity_clients: 272,
              clients: 941,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 183,
              non_entity_tokens: 0,
              non_entity_clients: 567,
              clients: 750,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 588,
              non_entity_tokens: 0,
              non_entity_clients: 130,
              clients: 718,
            },
          },
        ],
      },
      {
        namespace_id: 'PU6JB',
        namespace_path: 'test-ns-2/',
        counts: {
          distinct_entities: 0,
          entity_clients: 1792,
          non_entity_tokens: 0,
          non_entity_clients: 842,
          clients: 2634,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 611,
              non_entity_tokens: 0,
              non_entity_clients: 215,
              clients: 826,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 365,
              non_entity_tokens: 0,
              non_entity_clients: 368,
              clients: 733,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 469,
              non_entity_tokens: 0,
              non_entity_clients: 244,
              clients: 713,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 347,
              non_entity_tokens: 0,
              non_entity_clients: 15,
              clients: 362,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 0,
        entity_clients: 5855,
        non_entity_tokens: 0,
        non_entity_clients: 4729,
        clients: 10584,
      },
      namespaces: [
        {
          namespace_id: 'opmJ1',
          namespace_path: 'test-ns-1/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1409,
            non_entity_tokens: 0,
            non_entity_clients: 1429,
            clients: 2838,
          },
          mounts: [
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 532,
                non_entity_tokens: 0,
                non_entity_clients: 563,
                clients: 1095,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 655,
                non_entity_tokens: 0,
                non_entity_clients: 179,
                clients: 834,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 34,
                non_entity_tokens: 0,
                non_entity_clients: 658,
                clients: 692,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 188,
                non_entity_tokens: 0,
                non_entity_clients: 29,
                clients: 217,
              },
            },
          ],
        },
        {
          namespace_id: '3lq5r',
          namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1869,
            non_entity_tokens: 0,
            non_entity_clients: 592,
            clients: 2461,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 745,
                non_entity_tokens: 0,
                non_entity_clients: 239,
                clients: 984,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 539,
                non_entity_tokens: 0,
                non_entity_clients: 132,
                clients: 671,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 294,
                non_entity_tokens: 0,
                non_entity_clients: 110,
                clients: 404,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 291,
                non_entity_tokens: 0,
                non_entity_clients: 111,
                clients: 402,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 0,
            entity_clients: 838,
            non_entity_tokens: 0,
            non_entity_clients: 1486,
            clients: 2324,
          },
          mounts: [
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 629,
                non_entity_tokens: 0,
                non_entity_clients: 742,
                clients: 1371,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 166,
                non_entity_tokens: 0,
                non_entity_clients: 410,
                clients: 576,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 12,
                non_entity_tokens: 0,
                non_entity_clients: 279,
                clients: 291,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 31,
                non_entity_tokens: 0,
                non_entity_clients: 55,
                clients: 86,
              },
            },
          ],
        },
        {
          namespace_id: 'sJRLj',
          namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 996,
            non_entity_tokens: 0,
            non_entity_clients: 805,
            clients: 1801,
          },
          mounts: [
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 484,
                non_entity_tokens: 0,
                non_entity_clients: 145,
                clients: 629,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 396,
                non_entity_tokens: 0,
                non_entity_clients: 156,
                clients: 552,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 18,
                non_entity_tokens: 0,
                non_entity_clients: 401,
                clients: 419,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 98,
                non_entity_tokens: 0,
                non_entity_clients: 103,
                clients: 201,
              },
            },
          ],
        },
        {
          namespace_id: 'PU6JB',
          namespace_path: 'test-ns-2/',
          counts: {
            distinct_entities: 0,
            entity_clients: 743,
            non_entity_tokens: 0,
            non_entity_clients: 417,
            clients: 1160,
          },
          mounts: [
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 188,
                non_entity_tokens: 0,
                non_entity_clients: 168,
                clients: 356,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 196,
                non_entity_tokens: 0,
                non_entity_clients: 115,
                clients: 311,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 291,
                non_entity_tokens: 0,
                non_entity_clients: 3,
                clients: 294,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 68,
                non_entity_tokens: 0,
                non_entity_clients: 131,
                clients: 199,
              },
            },
          ],
        },
      ],
    },
  },
  {
    timestamp: formatISO(addMonths(UPGRADE_DATE, 3)),
    counts: {
      distinct_entities: 0,
      entity_clients: 10342,
      non_entity_tokens: 0,
      non_entity_clients: 13170,
      clients: 23512,
    },
    namespaces: [
      {
        namespace_id: 'PU6JB',
        namespace_path: 'test-ns-2/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2816,
          non_entity_tokens: 0,
          non_entity_clients: 3098,
          clients: 5914,
        },
        mounts: [
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 726,
              non_entity_tokens: 0,
              non_entity_clients: 995,
              clients: 1721,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 737,
              non_entity_tokens: 0,
              non_entity_clients: 850,
              clients: 1587,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 754,
              non_entity_tokens: 0,
              non_entity_clients: 617,
              clients: 1371,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 599,
              non_entity_tokens: 0,
              non_entity_clients: 636,
              clients: 1235,
            },
          },
        ],
      },
      {
        namespace_id: 'sJRLj',
        namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 2253,
          non_entity_tokens: 0,
          non_entity_clients: 2404,
          clients: 4657,
        },
        mounts: [
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 775,
              non_entity_tokens: 0,
              non_entity_clients: 689,
              clients: 1464,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 699,
              non_entity_tokens: 0,
              non_entity_clients: 652,
              clients: 1351,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 566,
              non_entity_tokens: 0,
              non_entity_clients: 487,
              clients: 1053,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 213,
              non_entity_tokens: 0,
              non_entity_clients: 576,
              clients: 789,
            },
          },
        ],
      },
      {
        namespace_id: 'opmJ1',
        namespace_path: 'test-ns-1/',
        counts: {
          distinct_entities: 0,
          entity_clients: 1725,
          non_entity_tokens: 0,
          non_entity_clients: 2927,
          clients: 4652,
        },
        mounts: [
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 811,
              non_entity_tokens: 0,
              non_entity_clients: 417,
              clients: 1228,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 294,
              non_entity_tokens: 0,
              non_entity_clients: 900,
              clients: 1194,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 503,
              non_entity_tokens: 0,
              non_entity_clients: 620,
              clients: 1123,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 117,
              non_entity_tokens: 0,
              non_entity_clients: 990,
              clients: 1107,
            },
          },
        ],
      },
      {
        namespace_id: 'root',
        namespace_path: '',
        counts: {
          distinct_entities: 0,
          entity_clients: 1678,
          non_entity_tokens: 0,
          non_entity_clients: 2775,
          clients: 4453,
        },
        mounts: [
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 972,
              non_entity_tokens: 0,
              non_entity_clients: 608,
              clients: 1580,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 172,
              non_entity_tokens: 0,
              non_entity_clients: 957,
              clients: 1129,
            },
          },
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 220,
              non_entity_tokens: 0,
              non_entity_clients: 756,
              clients: 976,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 314,
              non_entity_tokens: 0,
              non_entity_clients: 454,
              clients: 768,
            },
          },
        ],
      },
      {
        namespace_id: '3lq5r',
        namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
        counts: {
          distinct_entities: 0,
          entity_clients: 1870,
          non_entity_tokens: 0,
          non_entity_clients: 1966,
          clients: 3836,
        },
        mounts: [
          {
            mount_path: 'path-4-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 839,
              non_entity_tokens: 0,
              non_entity_clients: 762,
              clients: 1601,
            },
          },
          {
            mount_path: 'path-3-with-over-18-characters',
            counts: {
              distinct_entities: 0,
              entity_clients: 447,
              non_entity_tokens: 0,
              non_entity_clients: 583,
              clients: 1030,
            },
          },
          {
            mount_path: 'path-2',
            counts: {
              distinct_entities: 0,
              entity_clients: 382,
              non_entity_tokens: 0,
              non_entity_clients: 375,
              clients: 757,
            },
          },
          {
            mount_path: 'path-1',
            counts: {
              distinct_entities: 0,
              entity_clients: 202,
              non_entity_tokens: 0,
              non_entity_clients: 246,
              clients: 448,
            },
          },
        ],
      },
    ],
    new_clients: {
      counts: {
        distinct_entities: 0,
        entity_clients: 5959,
        non_entity_tokens: 0,
        non_entity_clients: 6985,
        clients: 12944,
      },
      namespaces: [
        {
          namespace_id: 'opmJ1',
          namespace_path: 'test-ns-1/',
          counts: {
            distinct_entities: 0,
            entity_clients: 873,
            non_entity_tokens: 0,
            non_entity_clients: 2355,
            clients: 3228,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 196,
                non_entity_tokens: 0,
                non_entity_clients: 811,
                clients: 1007,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 38,
                non_entity_tokens: 0,
                non_entity_clients: 931,
                clients: 969,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 148,
                non_entity_tokens: 0,
                non_entity_clients: 608,
                clients: 756,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 491,
                non_entity_tokens: 0,
                non_entity_clients: 5,
                clients: 496,
              },
            },
          ],
        },
        {
          namespace_id: 'sJRLj',
          namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1352,
            non_entity_tokens: 0,
            non_entity_clients: 1506,
            clients: 2858,
          },
          mounts: [
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 245,
                non_entity_tokens: 0,
                non_entity_clients: 560,
                clients: 805,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 465,
                non_entity_tokens: 0,
                non_entity_clients: 332,
                clients: 797,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 529,
                non_entity_tokens: 0,
                non_entity_clients: 117,
                clients: 646,
              },
            },
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 113,
                non_entity_tokens: 0,
                non_entity_clients: 497,
                clients: 610,
              },
            },
          ],
        },
        {
          namespace_id: '3lq5r',
          namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1355,
            non_entity_tokens: 0,
            non_entity_clients: 1353,
            clients: 2708,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 557,
                non_entity_tokens: 0,
                non_entity_clients: 538,
                clients: 1095,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 410,
                non_entity_tokens: 0,
                non_entity_clients: 496,
                clients: 906,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 146,
                non_entity_tokens: 0,
                non_entity_clients: 237,
                clients: 383,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 242,
                non_entity_tokens: 0,
                non_entity_clients: 82,
                clients: 324,
              },
            },
          ],
        },
        {
          namespace_id: 'PU6JB',
          namespace_path: 'test-ns-2/',
          counts: {
            distinct_entities: 0,
            entity_clients: 1514,
            non_entity_tokens: 0,
            non_entity_clients: 578,
            clients: 2092,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 602,
                non_entity_tokens: 0,
                non_entity_clients: 147,
                clients: 749,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 259,
                non_entity_tokens: 0,
                non_entity_clients: 344,
                clients: 603,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 349,
                non_entity_tokens: 0,
                non_entity_clients: 43,
                clients: 392,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 304,
                non_entity_tokens: 0,
                non_entity_clients: 44,
                clients: 348,
              },
            },
          ],
        },
        {
          namespace_id: 'root',
          namespace_path: '',
          counts: {
            distinct_entities: 0,
            entity_clients: 865,
            non_entity_tokens: 0,
            non_entity_clients: 1193,
            clients: 2058,
          },
          mounts: [
            {
              mount_path: 'path-4-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 10,
                non_entity_tokens: 0,
                non_entity_clients: 722,
                clients: 732,
              },
            },
            {
              mount_path: 'path-1',
              counts: {
                distinct_entities: 0,
                entity_clients: 643,
                non_entity_tokens: 0,
                non_entity_clients: 4,
                clients: 647,
              },
            },
            {
              mount_path: 'path-2',
              counts: {
                distinct_entities: 0,
                entity_clients: 93,
                non_entity_tokens: 0,
                non_entity_clients: 379,
                clients: 472,
              },
            },
            {
              mount_path: 'path-3-with-over-18-characters',
              counts: {
                distinct_entities: 0,
                entity_clients: 119,
                non_entity_tokens: 0,
                non_entity_clients: 88,
                clients: 207,
              },
            },
          ],
        },
      ],
    },
  },
];

function generateNullMonths(startDate, endDate) {
  let numberOfMonths = differenceInCalendarMonths(endDate, startDate);
  let months = [];
  for (let i = 0; i < numberOfMonths; i++) {
    months.push({
      timestamp: formatRFC3339(startOfMonth(addMonths(startDate, i))),
      counts: null,
      namespace: null,
      new_clients: null,
    });
    continue;
  }
  return months;
}

const handleMockQuery = (queryStartTimestamp, queryEndTimestamp, monthlyData) => {
  const queryStartDate = startOfMonth(parseAPITimestamp(queryStartTimestamp));
  const queryEndDate = parseAPITimestamp(queryEndTimestamp);
  // monthlyData is oldest to newest
  const dataEarliestMonth = parseAPITimestamp(monthlyData[0].timestamp);
  const dataLatestMonth = parseAPITimestamp(monthlyData[monthlyData.length - 1].timestamp);
  let transformedMonthlyArray = [...monthlyData];
  // If query end is before last month in array, return only through end query
  if (isBefore(queryEndDate, dataLatestMonth)) {
    let indexQueryStart = monthlyData.findIndex((e) =>
      isSameMonth(queryStartDate, parseAPITimestamp(e.timestamp))
    );
    let indexQueryEnd = monthlyData.findIndex((e) =>
      isSameMonth(queryEndDate, parseAPITimestamp(e.timestamp))
    );
    return transformedMonthlyArray.slice(indexQueryStart, indexQueryEnd + 1);
  }
  // If query wants months previous to the data we have, generate months without data prior
  if (isBefore(queryStartDate, dataEarliestMonth)) {
    return [...generateNullMonths(queryStartDate, dataEarliestMonth), ...transformedMonthlyArray];
  }
  // If query is after earliest month in array, return latest to month that matches query
  if (isAfter(queryStartDate, dataEarliestMonth)) {
    let index = monthlyData.findIndex((e) => isSameMonth(queryStartDate, parseAPITimestamp(e.timestamp)));
    return transformedMonthlyArray.slice(index);
  }
  return transformedMonthlyArray;
};

export default function (server) {
  server.get('sys/license/status', function () {
    return {
      request_id: 'my-license-request-id',
      data: {
        autoloaded: {
          license_id: 'my-license-id',
          start_time: formatRFC3339(LICENSE_START),
          expiration_time: formatRFC3339(LICENSE_END),
        },
      },
    };
  });

  server.get('sys/internal/counters/config', function () {
    return {
      request_id: 'some-config-id',
      data: {
        default_report_months: 12,
        enabled: 'default-enable',
        queries_available: true,
        retention_months: 24,
      },
    };
  });

  server.get('/sys/internal/counters/activity', (schema, req) => {
    const { start_time, end_time } = req.queryParams;

    return {
      request_id: '25f55fbb-f253-9c46-c6f0-3cdd3ada91ab',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        by_namespace: [
          {
            namespace_id: 'PU6JB',
            namespace_path: 'test-ns-2/',
            counts: {
              distinct_entities: 23326,
              entity_clients: 23326,
              non_entity_tokens: 17826,
              non_entity_clients: 17826,
              clients: 41152,
            },
            mounts: [
              {
                mount_path: 'path-2',
                counts: {
                  distinct_entities: 6508,
                  entity_clients: 6508,
                  non_entity_tokens: 3634,
                  non_entity_clients: 3634,
                  clients: 10142,
                },
              },
              {
                mount_path: 'path-4-with-over-18-characters',
                counts: {
                  distinct_entities: 5118,
                  entity_clients: 5118,
                  non_entity_tokens: 4942,
                  non_entity_clients: 4942,
                  clients: 10060,
                },
              },
              {
                mount_path: 'path-3-with-over-18-characters',
                counts: {
                  distinct_entities: 5931,
                  entity_clients: 5931,
                  non_entity_tokens: 4057,
                  non_entity_clients: 4057,
                  clients: 9988,
                },
              },
              {
                mount_path: 'path-1',
                counts: {
                  distinct_entities: 4962,
                  entity_clients: 4962,
                  non_entity_tokens: 3739,
                  non_entity_clients: 3739,
                  clients: 8701,
                },
              },
            ],
          },
          {
            namespace_id: '3lq5r',
            namespace_path: 'test-ns-2-with-namespace-length-over-18-characters/',
            counts: {
              distinct_entities: 19842,
              entity_clients: 19842,
              non_entity_tokens: 20799,
              non_entity_clients: 20799,
              clients: 40641,
            },
            mounts: [
              {
                mount_path: 'path-4-with-over-18-characters',
                counts: {
                  distinct_entities: 4695,
                  entity_clients: 4695,
                  non_entity_tokens: 6620,
                  non_entity_clients: 6620,
                  clients: 11315,
                },
              },
              {
                mount_path: 'path-3-with-over-18-characters',
                counts: {
                  distinct_entities: 5762,
                  entity_clients: 5762,
                  non_entity_tokens: 4112,
                  non_entity_clients: 4112,
                  clients: 9874,
                },
              },
              {
                mount_path: 'path-2',
                counts: {
                  distinct_entities: 5303,
                  entity_clients: 5303,
                  non_entity_tokens: 4538,
                  non_entity_clients: 4538,
                  clients: 9841,
                },
              },
              {
                mount_path: 'path-1',
                counts: {
                  distinct_entities: 3501,
                  entity_clients: 3501,
                  non_entity_tokens: 4974,
                  non_entity_clients: 4974,
                  clients: 8475,
                },
              },
            ],
          },
          {
            namespace_id: 'sJRLj',
            namespace_path: 'test-ns-1-with-namespace-length-over-18-characters/',
            counts: {
              distinct_entities: 20389,
              entity_clients: 20389,
              non_entity_tokens: 19445,
              non_entity_clients: 19445,
              clients: 39834,
            },
            mounts: [
              {
                mount_path: 'path-3-with-over-18-characters',
                counts: {
                  distinct_entities: 5356,
                  entity_clients: 5356,
                  non_entity_tokens: 5075,
                  non_entity_clients: 5075,
                  clients: 10431,
                },
              },
              {
                mount_path: 'path-4-with-over-18-characters',
                counts: {
                  distinct_entities: 4639,
                  entity_clients: 4639,
                  non_entity_tokens: 5242,
                  non_entity_clients: 5242,
                  clients: 9881,
                },
              },
              {
                mount_path: 'path-2',
                counts: {
                  distinct_entities: 4926,
                  entity_clients: 4926,
                  non_entity_tokens: 4163,
                  non_entity_clients: 4163,
                  clients: 9089,
                },
              },
              {
                mount_path: 'path-1',
                counts: {
                  distinct_entities: 4437,
                  entity_clients: 4437,
                  non_entity_tokens: 4201,
                  non_entity_clients: 4201,
                  clients: 8638,
                },
              },
            ],
          },
          {
            namespace_id: 'opmJ1',
            namespace_path: 'test-ns-1/',
            counts: {
              distinct_entities: 19316,
              entity_clients: 19316,
              non_entity_tokens: 18450,
              non_entity_clients: 18450,
              clients: 37766,
            },
            mounts: [
              {
                mount_path: 'path-3-with-over-18-characters',
                counts: {
                  distinct_entities: 4952,
                  entity_clients: 4952,
                  non_entity_tokens: 5080,
                  non_entity_clients: 5080,
                  clients: 10032,
                },
              },
              {
                mount_path: 'path-4-with-over-18-characters',
                counts: {
                  distinct_entities: 5198,
                  entity_clients: 5198,
                  non_entity_tokens: 3825,
                  non_entity_clients: 3825,
                  clients: 9023,
                },
              },
              {
                mount_path: 'path-1',
                counts: {
                  distinct_entities: 3827,
                  entity_clients: 3827,
                  non_entity_tokens: 5156,
                  non_entity_clients: 5156,
                  clients: 8983,
                },
              },
              {
                mount_path: 'path-2',
                counts: {
                  distinct_entities: 3981,
                  entity_clients: 3981,
                  non_entity_tokens: 3661,
                  non_entity_clients: 3661,
                  clients: 7642,
                },
              },
            ],
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 15416,
              entity_clients: 15416,
              non_entity_tokens: 19892,
              non_entity_clients: 19892,
              clients: 35308,
            },
            mounts: [
              {
                mount_path: 'path-2',
                counts: {
                  distinct_entities: 3936,
                  entity_clients: 3936,
                  non_entity_tokens: 5428,
                  non_entity_clients: 5428,
                  clients: 9364,
                },
              },
              {
                mount_path: 'path-1',
                counts: {
                  distinct_entities: 4021,
                  entity_clients: 4021,
                  non_entity_tokens: 4530,
                  non_entity_clients: 4530,
                  clients: 8551,
                },
              },
              {
                mount_path: 'path-4-with-over-18-characters',
                counts: {
                  distinct_entities: 2934,
                  entity_clients: 2934,
                  non_entity_tokens: 5357,
                  non_entity_clients: 5357,
                  clients: 8291,
                },
              },
              {
                mount_path: 'path-3-with-over-18-characters',
                counts: {
                  distinct_entities: 3938,
                  entity_clients: 3938,
                  non_entity_tokens: 3932,
                  non_entity_clients: 3932,
                  clients: 7870,
                },
              },
            ],
          },
        ],
        end_time: end_time || formatISO(endOfMonth(sub(NEW_DATE, { months: 1 }))),
        months: handleMockQuery(start_time, end_time, MOCK_MONTHLY_DATA),
        start_time: isBefore(new Date(start_time), COUNTS_START) ? formatRFC3339(COUNTS_START) : start_time,
        total: {
          distinct_entities: 98289,
          entity_clients: 98289,
          non_entity_tokens: 96412,
          non_entity_clients: 96412,
          clients: 194701,
        },
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });

  server.get('/sys/internal/counters/activity/monthly', function () {
    const timestamp = NEW_DATE;
    return {
      request_id: '26be5ab9-dcac-9237-ec12-269a8ca64742',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        by_namespace: [
          {
            namespace_id: '0lHBL',
            namespace_path: 'ns1/',
            counts: {
              distinct_entities: 85,
              non_entity_tokens: 15,
              clients: 100,
            },
            mounts: [
              {
                mount_path: 'auth/method/uMGBU',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                mount_path: 'auth/method/woiej',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
            ],
          },
          {
            namespace_id: 'RxD81',
            namespace_path: 'ns2/',
            counts: {
              distinct_entities: 35,
              non_entity_tokens: 20,
              clients: 55,
            },
            mounts: [
              {
                mount_path: 'auth/method/ABCD1',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                mount_path: 'auth/method/ABCD2',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
            ],
          },
          {
            namespace_id: 'root',
            namespace_path: '',
            counts: {
              distinct_entities: 12,
              non_entity_tokens: 8,
              clients: 20,
            },
            mounts: [
              {
                mount_path: 'auth/method/XYZZ2',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                mount_path: 'auth/method/XYZZ1',
                counts: {
                  clients: 35,
                  entity_clients: 20,
                  non_entity_clients: 15,
                },
              },
              {
                mount_path: 'auth_userpass_3158c012',
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
            timestamp: startOfMonth(timestamp).toISOString(),
            counts: {
              distinct_entities: 0,
              entity_clients: 4,
              non_entity_tokens: 0,
              non_entity_clients: 0,
              clients: 4,
            },
            namespaces: [
              {
                namespace_id: 'lHmap',
                namespace_path: 'education/',
                counts: {
                  distinct_entities: 0,
                  entity_clients: 2,
                  non_entity_tokens: 0,
                  non_entity_clients: 0,
                  clients: 2,
                },
                mounts: [
                  {
                    mount_path: 'auth_userpass_a36c8125',
                    counts: {
                      distinct_entities: 0,
                      entity_clients: 2,
                      non_entity_tokens: 0,
                      non_entity_clients: 0,
                      clients: 2,
                    },
                  },
                ],
              },
              {
                namespace_id: 'root',
                namespace_path: '',
                counts: {
                  distinct_entities: 0,
                  entity_clients: 2,
                  non_entity_tokens: 0,
                  non_entity_clients: 0,
                  clients: 2,
                },
                mounts: [
                  {
                    mount_path: 'auth_userpass_3158c012',
                    counts: {
                      distinct_entities: 0,
                      entity_clients: 2,
                      non_entity_tokens: 0,
                      non_entity_clients: 0,
                      clients: 2,
                    },
                  },
                ],
              },
            ],
            new_clients: {
              counts: {
                distinct_entities: 0,
                entity_clients: 4,
                non_entity_tokens: 0,
                non_entity_clients: 0,
                clients: 4,
              },
              namespaces: [
                {
                  namespace_id: 'root',
                  namespace_path: '',
                  counts: {
                    distinct_entities: 0,
                    entity_clients: 2,
                    non_entity_tokens: 0,
                    non_entity_clients: 0,
                    clients: 2,
                  },
                  mounts: [
                    {
                      mount_path: 'auth_userpass_3158c012',
                      counts: {
                        distinct_entities: 0,
                        entity_clients: 2,
                        non_entity_tokens: 0,
                        non_entity_clients: 0,
                        clients: 2,
                      },
                    },
                  ],
                },
                {
                  namespace_id: 'lHmap',
                  namespace_path: 'education/',
                  counts: {
                    distinct_entities: 0,
                    entity_clients: 2,
                    non_entity_tokens: 0,
                    non_entity_clients: 0,
                    clients: 2,
                  },
                  mounts: [
                    {
                      mount_path: 'auth_userpass_a36c8125',
                      counts: {
                        distinct_entities: 0,
                        entity_clients: 2,
                        non_entity_tokens: 0,
                        non_entity_clients: 0,
                        clients: 2,
                      },
                    },
                  ],
                },
              ],
            },
          },
        ],
        distinct_entities: 132,
        entity_clients: 132,
        non_entity_tokens: 43,
        non_entity_clients: 43,
        clients: 175,
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });
}
