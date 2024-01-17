/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { formatRFC3339 } from 'date-fns';
import { findAll } from '@ember/test-helpers';
import { calculateAverage } from 'vault/utils/chart-helpers';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | clients/running-total', function (hooks) {
  setupRenderingTest(hooks);
  const MONTHLY_ACTIVITY = [
    {
      month: '8/21',
      timestamp: '2021-08-01T00:00:00Z',
      counts: null,
      namespaces: [],
      new_clients: {
        month: '8/21',
        namespaces: [],
      },
      namespaces_by_key: {},
    },
    {
      month: '9/21',
      clients: 19251,
      entity_clients: 10713,
      non_entity_clients: 8538,
      namespaces: [
        {
          label: 'root',
          clients: 4852,
          entity_clients: 3108,
          non_entity_clients: 1744,
          mounts: [
            {
              label: 'path-3-with-over-18-characters',
              clients: 1598,
              entity_clients: 687,
              non_entity_clients: 911,
            },
            {
              label: 'path-1',
              clients: 1429,
              entity_clients: 981,
              non_entity_clients: 448,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 965,
              entity_clients: 720,
              non_entity_clients: 245,
            },
            {
              label: 'path-2',
              clients: 860,
              entity_clients: 720,
              non_entity_clients: 140,
            },
          ],
        },
        {
          label: 'test-ns-2/',
          clients: 4702,
          entity_clients: 3057,
          non_entity_clients: 1645,
          mounts: [
            {
              label: 'path-3-with-over-18-characters',
              clients: 1686,
              entity_clients: 926,
              non_entity_clients: 760,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 1525,
              entity_clients: 789,
              non_entity_clients: 736,
            },
            {
              label: 'path-2',
              clients: 905,
              entity_clients: 849,
              non_entity_clients: 56,
            },
            {
              label: 'path-1',
              clients: 586,
              entity_clients: 493,
              non_entity_clients: 93,
            },
          ],
        },
        {
          label: 'test-ns-1/',
          clients: 4569,
          entity_clients: 1871,
          non_entity_clients: 2698,
          mounts: [
            {
              label: 'path-4-with-over-18-characters',
              clients: 1534,
              entity_clients: 619,
              non_entity_clients: 915,
            },
            {
              label: 'path-3-with-over-18-characters',
              clients: 1528,
              entity_clients: 589,
              non_entity_clients: 939,
            },
            {
              label: 'path-1',
              clients: 828,
              entity_clients: 612,
              non_entity_clients: 216,
            },
            {
              label: 'path-2',
              clients: 679,
              entity_clients: 51,
              non_entity_clients: 628,
            },
          ],
        },
        {
          label: 'test-ns-2-with-namespace-length-over-18-characters/',
          clients: 3771,
          entity_clients: 2029,
          non_entity_clients: 1742,
          mounts: [
            {
              label: 'path-3-with-over-18-characters',
              clients: 1249,
              entity_clients: 793,
              non_entity_clients: 456,
            },
            {
              label: 'path-1',
              clients: 1046,
              entity_clients: 444,
              non_entity_clients: 602,
            },
            {
              label: 'path-2',
              clients: 930,
              entity_clients: 277,
              non_entity_clients: 653,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 546,
              entity_clients: 515,
              non_entity_clients: 31,
            },
          ],
        },
        {
          label: 'test-ns-1-with-namespace-length-over-18-characters/',
          clients: 1357,
          entity_clients: 648,
          non_entity_clients: 709,
          mounts: [
            {
              label: 'path-1',
              clients: 613,
              entity_clients: 23,
              non_entity_clients: 590,
            },
            {
              label: 'path-3-with-over-18-characters',
              clients: 543,
              entity_clients: 465,
              non_entity_clients: 78,
            },
            {
              label: 'path-2',
              clients: 146,
              entity_clients: 141,
              non_entity_clients: 5,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 55,
              entity_clients: 19,
              non_entity_clients: 36,
            },
          ],
        },
      ],
      namespaces_by_key: {
        root: {
          month: '9/21',
          clients: 4852,
          entity_clients: 3108,
          non_entity_clients: 1744,
          new_clients: {
            month: '9/21',
            label: 'root',
            clients: 2525,
            entity_clients: 1315,
            non_entity_clients: 1210,
          },
          mounts_by_key: {
            'path-3-with-over-18-characters': {
              month: '9/21',
              label: 'path-3-with-over-18-characters',
              clients: 1598,
              entity_clients: 687,
              non_entity_clients: 911,
              new_clients: {
                month: '9/21',
                label: 'path-3-with-over-18-characters',
                clients: 1055,
                entity_clients: 257,
                non_entity_clients: 798,
              },
            },
            'path-1': {
              month: '9/21',
              label: 'path-1',
              clients: 1429,
              entity_clients: 981,
              non_entity_clients: 448,
              new_clients: {
                month: '9/21',
                label: 'path-1',
                clients: 543,
                entity_clients: 340,
                non_entity_clients: 203,
              },
            },
            'path-4-with-over-18-characters': {
              month: '9/21',
              label: 'path-4-with-over-18-characters',
              clients: 965,
              entity_clients: 720,
              non_entity_clients: 245,
              new_clients: {
                month: '9/21',
                label: 'path-4-with-over-18-characters',
                clients: 136,
                entity_clients: 7,
                non_entity_clients: 129,
              },
            },
            'path-2': {
              month: '9/21',
              label: 'path-2',
              clients: 860,
              entity_clients: 720,
              non_entity_clients: 140,
              new_clients: {
                month: '9/21',
                label: 'path-2',
                clients: 791,
                entity_clients: 711,
                non_entity_clients: 80,
              },
            },
          },
        },
        'test-ns-2/': {
          month: '9/21',
          clients: 4702,
          entity_clients: 3057,
          non_entity_clients: 1645,
          new_clients: {
            month: '9/21',
            label: 'test-ns-2/',
            clients: 1537,
            entity_clients: 662,
            non_entity_clients: 875,
          },
          mounts_by_key: {
            'path-3-with-over-18-characters': {
              month: '9/21',
              label: 'path-3-with-over-18-characters',
              clients: 1686,
              entity_clients: 926,
              non_entity_clients: 760,
              new_clients: {
                month: '9/21',
                label: 'path-3-with-over-18-characters',
                clients: 520,
                entity_clients: 13,
                non_entity_clients: 507,
              },
            },
            'path-4-with-over-18-characters': {
              month: '9/21',
              label: 'path-4-with-over-18-characters',
              clients: 1525,
              entity_clients: 789,
              non_entity_clients: 736,
              new_clients: {
                month: '9/21',
                label: 'path-4-with-over-18-characters',
                clients: 499,
                entity_clients: 197,
                non_entity_clients: 302,
              },
            },
            'path-2': {
              month: '9/21',
              label: 'path-2',
              clients: 905,
              entity_clients: 849,
              non_entity_clients: 56,
              new_clients: {
                month: '9/21',
                label: 'path-2',
                clients: 398,
                entity_clients: 370,
                non_entity_clients: 28,
              },
            },
            'path-1': {
              month: '9/21',
              label: 'path-1',
              clients: 586,
              entity_clients: 493,
              non_entity_clients: 93,
              new_clients: {
                month: '9/21',
                label: 'path-1',
                clients: 120,
                entity_clients: 82,
                non_entity_clients: 38,
              },
            },
          },
        },
        'test-ns-1/': {
          month: '9/21',
          clients: 4569,
          entity_clients: 1871,
          non_entity_clients: 2698,
          new_clients: {
            month: '9/21',
            label: 'test-ns-1/',
            clients: 2712,
            entity_clients: 879,
            non_entity_clients: 1833,
          },
          mounts_by_key: {
            'path-4-with-over-18-characters': {
              month: '9/21',
              label: 'path-4-with-over-18-characters',
              clients: 1534,
              entity_clients: 619,
              non_entity_clients: 915,
              new_clients: {
                month: '9/21',
                label: 'path-4-with-over-18-characters',
                clients: 740,
                entity_clients: 39,
                non_entity_clients: 701,
              },
            },
            'path-3-with-over-18-characters': {
              month: '9/21',
              label: 'path-3-with-over-18-characters',
              clients: 1528,
              entity_clients: 589,
              non_entity_clients: 939,
              new_clients: {
                month: '9/21',
                label: 'path-3-with-over-18-characters',
                clients: 1250,
                entity_clients: 536,
                non_entity_clients: 714,
              },
            },
            'path-1': {
              month: '9/21',
              label: 'path-1',
              clients: 828,
              entity_clients: 612,
              non_entity_clients: 216,
              new_clients: {
                month: '9/21',
                label: 'path-1',
                clients: 463,
                entity_clients: 283,
                non_entity_clients: 180,
              },
            },
            'path-2': {
              month: '9/21',
              label: 'path-2',
              clients: 679,
              entity_clients: 51,
              non_entity_clients: 628,
              new_clients: {
                month: '9/21',
                label: 'path-2',
                clients: 259,
                entity_clients: 21,
                non_entity_clients: 238,
              },
            },
          },
        },
        'test-ns-2-with-namespace-length-over-18-characters/': {
          month: '9/21',
          clients: 3771,
          entity_clients: 2029,
          non_entity_clients: 1742,
          new_clients: {
            month: '9/21',
            label: 'test-ns-2-with-namespace-length-over-18-characters/',
            clients: 2087,
            entity_clients: 902,
            non_entity_clients: 1185,
          },
          mounts_by_key: {
            'path-3-with-over-18-characters': {
              month: '9/21',
              label: 'path-3-with-over-18-characters',
              clients: 1249,
              entity_clients: 793,
              non_entity_clients: 456,
              new_clients: {
                month: '9/21',
                label: 'path-3-with-over-18-characters',
                clients: 472,
                entity_clients: 260,
                non_entity_clients: 212,
              },
            },
            'path-1': {
              month: '9/21',
              label: 'path-1',
              clients: 1046,
              entity_clients: 444,
              non_entity_clients: 602,
              new_clients: {
                month: '9/21',
                label: 'path-1',
                clients: 775,
                entity_clients: 349,
                non_entity_clients: 426,
              },
            },
            'path-2': {
              month: '9/21',
              label: 'path-2',
              clients: 930,
              entity_clients: 277,
              non_entity_clients: 653,
              new_clients: {
                month: '9/21',
                label: 'path-2',
                clients: 632,
                entity_clients: 90,
                non_entity_clients: 542,
              },
            },
            'path-4-with-over-18-characters': {
              month: '9/21',
              label: 'path-4-with-over-18-characters',
              clients: 546,
              entity_clients: 515,
              non_entity_clients: 31,
              new_clients: {
                month: '9/21',
                label: 'path-4-with-over-18-characters',
                clients: 208,
                entity_clients: 203,
                non_entity_clients: 5,
              },
            },
          },
        },
        'test-ns-1-with-namespace-length-over-18-characters/': {
          month: '9/21',
          clients: 1357,
          entity_clients: 648,
          non_entity_clients: 709,
          new_clients: {
            month: '9/21',
            label: 'test-ns-1-with-namespace-length-over-18-characters/',
            clients: 560,
            entity_clients: 189,
            non_entity_clients: 371,
          },
          mounts_by_key: {
            'path-1': {
              month: '9/21',
              label: 'path-1',
              clients: 613,
              entity_clients: 23,
              non_entity_clients: 590,
              new_clients: {
                month: '9/21',
                label: 'path-1',
                clients: 318,
                entity_clients: 12,
                non_entity_clients: 306,
              },
            },
            'path-3-with-over-18-characters': {
              month: '9/21',
              label: 'path-3-with-over-18-characters',
              clients: 543,
              entity_clients: 465,
              non_entity_clients: 78,
              new_clients: {
                month: '9/21',
                label: 'path-3-with-over-18-characters',
                clients: 126,
                entity_clients: 89,
                non_entity_clients: 37,
              },
            },
            'path-2': {
              month: '9/21',
              label: 'path-2',
              clients: 146,
              entity_clients: 141,
              non_entity_clients: 5,
              new_clients: {
                month: '9/21',
                label: 'path-2',
                clients: 76,
                entity_clients: 75,
                non_entity_clients: 1,
              },
            },
            'path-4-with-over-18-characters': {
              month: '9/21',
              label: 'path-4-with-over-18-characters',
              clients: 55,
              entity_clients: 19,
              non_entity_clients: 36,
              new_clients: {
                month: '9/21',
                label: 'path-4-with-over-18-characters',
                clients: 40,
                entity_clients: 13,
                non_entity_clients: 27,
              },
            },
          },
        },
      },
      new_clients: {
        month: '9/21',
        clients: 9421,
        entity_clients: 3947,
        non_entity_clients: 5474,
        namespaces: [
          {
            label: 'test-ns-1/',
            clients: 2712,
            entity_clients: 879,
            non_entity_clients: 1833,
            mounts: [
              {
                label: 'path-3-with-over-18-characters',
                clients: 1250,
                entity_clients: 536,
                non_entity_clients: 714,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 740,
                entity_clients: 39,
                non_entity_clients: 701,
              },
              {
                label: 'path-1',
                clients: 463,
                entity_clients: 283,
                non_entity_clients: 180,
              },
              {
                label: 'path-2',
                clients: 259,
                entity_clients: 21,
                non_entity_clients: 238,
              },
            ],
          },
          {
            label: 'root',
            clients: 2525,
            entity_clients: 1315,
            non_entity_clients: 1210,
            mounts: [
              {
                label: 'path-3-with-over-18-characters',
                clients: 1055,
                entity_clients: 257,
                non_entity_clients: 798,
              },
              {
                label: 'path-2',
                clients: 791,
                entity_clients: 711,
                non_entity_clients: 80,
              },
              {
                label: 'path-1',
                clients: 543,
                entity_clients: 340,
                non_entity_clients: 203,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 136,
                entity_clients: 7,
                non_entity_clients: 129,
              },
            ],
          },
          {
            label: 'test-ns-2-with-namespace-length-over-18-characters/',
            clients: 2087,
            entity_clients: 902,
            non_entity_clients: 1185,
            mounts: [
              {
                label: 'path-1',
                clients: 775,
                entity_clients: 349,
                non_entity_clients: 426,
              },
              {
                label: 'path-2',
                clients: 632,
                entity_clients: 90,
                non_entity_clients: 542,
              },
              {
                label: 'path-3-with-over-18-characters',
                clients: 472,
                entity_clients: 260,
                non_entity_clients: 212,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 208,
                entity_clients: 203,
                non_entity_clients: 5,
              },
            ],
          },
          {
            label: 'test-ns-2/',
            clients: 1537,
            entity_clients: 662,
            non_entity_clients: 875,
            mounts: [
              {
                label: 'path-3-with-over-18-characters',
                clients: 520,
                entity_clients: 13,
                non_entity_clients: 507,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 499,
                entity_clients: 197,
                non_entity_clients: 302,
              },
              {
                label: 'path-2',
                clients: 398,
                entity_clients: 370,
                non_entity_clients: 28,
              },
              {
                label: 'path-1',
                clients: 120,
                entity_clients: 82,
                non_entity_clients: 38,
              },
            ],
          },
          {
            label: 'test-ns-1-with-namespace-length-over-18-characters/',
            clients: 560,
            entity_clients: 189,
            non_entity_clients: 371,
            mounts: [
              {
                label: 'path-1',
                clients: 318,
                entity_clients: 12,
                non_entity_clients: 306,
              },
              {
                label: 'path-3-with-over-18-characters',
                clients: 126,
                entity_clients: 89,
                non_entity_clients: 37,
              },
              {
                label: 'path-2',
                clients: 76,
                entity_clients: 75,
                non_entity_clients: 1,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 40,
                entity_clients: 13,
                non_entity_clients: 27,
              },
            ],
          },
        ],
      },
    },
    {
      month: '10/21',
      clients: 19417,
      entity_clients: 10105,
      non_entity_clients: 9312,
      namespaces: [
        {
          label: 'root',
          clients: 4835,
          entity_clients: 2364,
          non_entity_clients: 2471,
          mounts: [
            {
              label: 'path-3-with-over-18-characters',
              clients: 1797,
              entity_clients: 883,
              non_entity_clients: 914,
            },
            {
              label: 'path-1',
              clients: 1501,
              entity_clients: 663,
              non_entity_clients: 838,
            },
            {
              label: 'path-2',
              clients: 1461,
              entity_clients: 800,
              non_entity_clients: 661,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 76,
              entity_clients: 18,
              non_entity_clients: 58,
            },
          ],
        },
        {
          label: 'test-ns-2/',
          clients: 4027,
          entity_clients: 1692,
          non_entity_clients: 2335,
          mounts: [
            {
              label: 'path-4-with-over-18-characters',
              clients: 1223,
              entity_clients: 820,
              non_entity_clients: 403,
            },
            {
              label: 'path-3-with-over-18-characters',
              clients: 1110,
              entity_clients: 111,
              non_entity_clients: 999,
            },
            {
              label: 'path-1',
              clients: 1034,
              entity_clients: 462,
              non_entity_clients: 572,
            },
            {
              label: 'path-2',
              clients: 660,
              entity_clients: 299,
              non_entity_clients: 361,
            },
          ],
        },
        {
          label: 'test-ns-2-with-namespace-length-over-18-characters/',
          clients: 3924,
          entity_clients: 2132,
          non_entity_clients: 1792,
          mounts: [
            {
              label: 'path-3-with-over-18-characters',
              clients: 1411,
              entity_clients: 765,
              non_entity_clients: 646,
            },
            {
              label: 'path-2',
              clients: 1205,
              entity_clients: 382,
              non_entity_clients: 823,
            },
            {
              label: 'path-1',
              clients: 884,
              entity_clients: 850,
              non_entity_clients: 34,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 424,
              entity_clients: 135,
              non_entity_clients: 289,
            },
          ],
        },
        {
          label: 'test-ns-1-with-namespace-length-over-18-characters/',
          clients: 3639,
          entity_clients: 2314,
          non_entity_clients: 1325,
          mounts: [
            {
              label: 'path-1',
              clients: 1062,
              entity_clients: 781,
              non_entity_clients: 281,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 1021,
              entity_clients: 609,
              non_entity_clients: 412,
            },
            {
              label: 'path-2',
              clients: 849,
              entity_clients: 426,
              non_entity_clients: 423,
            },
            {
              label: 'path-3-with-over-18-characters',
              clients: 707,
              entity_clients: 498,
              non_entity_clients: 209,
            },
          ],
        },
        {
          label: 'test-ns-1/',
          clients: 2992,
          entity_clients: 1603,
          non_entity_clients: 1389,
          mounts: [
            {
              label: 'path-1',
              clients: 1140,
              entity_clients: 480,
              non_entity_clients: 660,
            },
            {
              label: 'path-4-with-over-18-characters',
              clients: 1058,
              entity_clients: 651,
              non_entity_clients: 407,
            },
            {
              label: 'path-2',
              clients: 575,
              entity_clients: 416,
              non_entity_clients: 159,
            },
            {
              label: 'path-3-with-over-18-characters',
              clients: 219,
              entity_clients: 56,
              non_entity_clients: 163,
            },
          ],
        },
      ],
      namespaces_by_key: {
        root: {
          month: '10/21',
          clients: 4835,
          entity_clients: 2364,
          non_entity_clients: 2471,
          new_clients: {
            month: '10/21',
            label: 'root',
            clients: 1732,
            entity_clients: 586,
            non_entity_clients: 1146,
          },
          mounts_by_key: {
            'path-3-with-over-18-characters': {
              month: '10/21',
              label: 'path-3-with-over-18-characters',
              clients: 1797,
              entity_clients: 883,
              non_entity_clients: 914,
              new_clients: {
                month: '10/21',
                label: 'path-3-with-over-18-characters',
                clients: 907,
                entity_clients: 192,
                non_entity_clients: 715,
              },
            },
            'path-1': {
              month: '10/21',
              label: 'path-1',
              clients: 1501,
              entity_clients: 663,
              non_entity_clients: 838,
              new_clients: {
                month: '10/21',
                label: 'path-1',
                clients: 276,
                entity_clients: 202,
                non_entity_clients: 74,
              },
            },
            'path-2': {
              month: '10/21',
              label: 'path-2',
              clients: 1461,
              entity_clients: 800,
              non_entity_clients: 661,
              new_clients: {
                month: '10/21',
                label: 'path-2',
                clients: 502,
                entity_clients: 189,
                non_entity_clients: 313,
              },
            },
            'path-4-with-over-18-characters': {
              month: '10/21',
              label: 'path-4-with-over-18-characters',
              clients: 76,
              entity_clients: 18,
              non_entity_clients: 58,
              new_clients: {
                month: '10/21',
                label: 'path-4-with-over-18-characters',
                clients: 47,
                entity_clients: 3,
                non_entity_clients: 44,
              },
            },
          },
        },
        'test-ns-2/': {
          month: '10/21',
          clients: 4027,
          entity_clients: 1692,
          non_entity_clients: 2335,
          new_clients: {
            month: '10/21',
            label: 'test-ns-2/',
            clients: 2301,
            entity_clients: 678,
            non_entity_clients: 1623,
          },
          mounts_by_key: {
            'path-4-with-over-18-characters': {
              month: '10/21',
              label: 'path-4-with-over-18-characters',
              clients: 1223,
              entity_clients: 820,
              non_entity_clients: 403,
              new_clients: {
                month: '10/21',
                label: 'path-4-with-over-18-characters',
                clients: 602,
                entity_clients: 212,
                non_entity_clients: 390,
              },
            },
            'path-3-with-over-18-characters': {
              month: '10/21',
              label: 'path-3-with-over-18-characters',
              clients: 1110,
              entity_clients: 111,
              non_entity_clients: 999,
              new_clients: {
                month: '10/21',
                label: 'path-3-with-over-18-characters',
                clients: 440,
                entity_clients: 7,
                non_entity_clients: 433,
              },
            },
            'path-1': {
              month: '10/21',
              label: 'path-1',
              clients: 1034,
              entity_clients: 462,
              non_entity_clients: 572,
              new_clients: {
                month: '10/21',
                label: 'path-1',
                clients: 980,
                entity_clients: 454,
                non_entity_clients: 526,
              },
            },
            'path-2': {
              month: '10/21',
              label: 'path-2',
              clients: 660,
              entity_clients: 299,
              non_entity_clients: 361,
              new_clients: {
                month: '10/21',
                label: 'path-2',
                clients: 279,
                entity_clients: 5,
                non_entity_clients: 274,
              },
            },
          },
        },
        'test-ns-2-with-namespace-length-over-18-characters/': {
          month: '10/21',
          clients: 3924,
          entity_clients: 2132,
          non_entity_clients: 1792,
          new_clients: {
            month: '10/21',
            label: 'test-ns-2-with-namespace-length-over-18-characters/',
            clients: 1561,
            entity_clients: 1225,
            non_entity_clients: 336,
          },
          mounts_by_key: {
            'path-3-with-over-18-characters': {
              month: '10/21',
              label: 'path-3-with-over-18-characters',
              clients: 1411,
              entity_clients: 765,
              non_entity_clients: 646,
              new_clients: {
                month: '10/21',
                label: 'path-3-with-over-18-characters',
                clients: 948,
                entity_clients: 660,
                non_entity_clients: 288,
              },
            },
            'path-2': {
              month: '10/21',
              label: 'path-2',
              clients: 1205,
              entity_clients: 382,
              non_entity_clients: 823,
              new_clients: {
                month: '10/21',
                label: 'path-2',
                clients: 305,
                entity_clients: 289,
                non_entity_clients: 16,
              },
            },
            'path-1': {
              month: '10/21',
              label: 'path-1',
              clients: 884,
              entity_clients: 850,
              non_entity_clients: 34,
              new_clients: {
                month: '10/21',
                label: 'path-1',
                clients: 230,
                entity_clients: 207,
                non_entity_clients: 23,
              },
            },
            'path-4-with-over-18-characters': {
              month: '10/21',
              label: 'path-4-with-over-18-characters',
              clients: 424,
              entity_clients: 135,
              non_entity_clients: 289,
              new_clients: {
                month: '10/21',
                label: 'path-4-with-over-18-characters',
                clients: 78,
                entity_clients: 69,
                non_entity_clients: 9,
              },
            },
          },
        },
        'test-ns-1-with-namespace-length-over-18-characters/': {
          month: '10/21',
          clients: 3639,
          entity_clients: 2314,
          non_entity_clients: 1325,
          new_clients: {
            month: '10/21',
            label: 'test-ns-1-with-namespace-length-over-18-characters/',
            clients: 1245,
            entity_clients: 710,
            non_entity_clients: 535,
          },
          mounts_by_key: {
            'path-1': {
              month: '10/21',
              label: 'path-1',
              clients: 1062,
              entity_clients: 781,
              non_entity_clients: 281,
              new_clients: {
                month: '10/21',
                label: 'path-1',
                clients: 288,
                entity_clients: 63,
                non_entity_clients: 225,
              },
            },
            'path-4-with-over-18-characters': {
              month: '10/21',
              label: 'path-4-with-over-18-characters',
              clients: 1021,
              entity_clients: 609,
              non_entity_clients: 412,
              new_clients: {
                month: '10/21',
                label: 'path-4-with-over-18-characters',
                clients: 440,
                entity_clients: 323,
                non_entity_clients: 117,
              },
            },
            'path-2': {
              month: '10/21',
              label: 'path-2',
              clients: 849,
              entity_clients: 426,
              non_entity_clients: 423,
              new_clients: {
                month: '10/21',
                label: 'path-2',
                clients: 339,
                entity_clients: 308,
                non_entity_clients: 31,
              },
            },
            'path-3-with-over-18-characters': {
              month: '10/21',
              label: 'path-3-with-over-18-characters',
              clients: 707,
              entity_clients: 498,
              non_entity_clients: 209,
              new_clients: {
                month: '10/21',
                label: 'path-3-with-over-18-characters',
                clients: 178,
                entity_clients: 16,
                non_entity_clients: 162,
              },
            },
          },
        },
        'test-ns-1/': {
          month: '10/21',
          clients: 2992,
          entity_clients: 1603,
          non_entity_clients: 1389,
          new_clients: {
            month: '10/21',
            label: 'test-ns-1/',
            clients: 820,
            entity_clients: 356,
            non_entity_clients: 464,
          },
          mounts_by_key: {
            'path-1': {
              month: '10/21',
              label: 'path-1',
              clients: 1140,
              entity_clients: 480,
              non_entity_clients: 660,
              new_clients: {
                month: '10/21',
                label: 'path-1',
                clients: 239,
                entity_clients: 30,
                non_entity_clients: 209,
              },
            },
            'path-4-with-over-18-characters': {
              month: '10/21',
              label: 'path-4-with-over-18-characters',
              clients: 1058,
              entity_clients: 651,
              non_entity_clients: 407,
              new_clients: {
                month: '10/21',
                label: 'path-4-with-over-18-characters',
                clients: 256,
                entity_clients: 63,
                non_entity_clients: 193,
              },
            },
            'path-2': {
              month: '10/21',
              label: 'path-2',
              clients: 575,
              entity_clients: 416,
              non_entity_clients: 159,
              new_clients: {
                month: '10/21',
                label: 'path-2',
                clients: 259,
                entity_clients: 245,
                non_entity_clients: 14,
              },
            },
            'path-3-with-over-18-characters': {
              month: '10/21',
              label: 'path-3-with-over-18-characters',
              clients: 219,
              entity_clients: 56,
              non_entity_clients: 163,
              new_clients: {
                month: '10/21',
                label: 'path-3-with-over-18-characters',
                clients: 66,
                entity_clients: 18,
                non_entity_clients: 48,
              },
            },
          },
        },
      },
      new_clients: {
        month: '10/21',
        clients: 7659,
        entity_clients: 3555,
        non_entity_clients: 4104,
        namespaces: [
          {
            label: 'test-ns-2/',
            clients: 2301,
            entity_clients: 678,
            non_entity_clients: 1623,
            mounts: [
              {
                label: 'path-1',
                clients: 980,
                entity_clients: 454,
                non_entity_clients: 526,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 602,
                entity_clients: 212,
                non_entity_clients: 390,
              },
              {
                label: 'path-3-with-over-18-characters',
                clients: 440,
                entity_clients: 7,
                non_entity_clients: 433,
              },
              {
                label: 'path-2',
                clients: 279,
                entity_clients: 5,
                non_entity_clients: 274,
              },
            ],
          },
          {
            label: 'root',
            clients: 1732,
            entity_clients: 586,
            non_entity_clients: 1146,
            mounts: [
              {
                label: 'path-3-with-over-18-characters',
                clients: 907,
                entity_clients: 192,
                non_entity_clients: 715,
              },
              {
                label: 'path-2',
                clients: 502,
                entity_clients: 189,
                non_entity_clients: 313,
              },
              {
                label: 'path-1',
                clients: 276,
                entity_clients: 202,
                non_entity_clients: 74,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 47,
                entity_clients: 3,
                non_entity_clients: 44,
              },
            ],
          },
          {
            label: 'test-ns-2-with-namespace-length-over-18-characters/',
            clients: 1561,
            entity_clients: 1225,
            non_entity_clients: 336,
            mounts: [
              {
                label: 'path-3-with-over-18-characters',
                clients: 948,
                entity_clients: 660,
                non_entity_clients: 288,
              },
              {
                label: 'path-2',
                clients: 305,
                entity_clients: 289,
                non_entity_clients: 16,
              },
              {
                label: 'path-1',
                clients: 230,
                entity_clients: 207,
                non_entity_clients: 23,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 78,
                entity_clients: 69,
                non_entity_clients: 9,
              },
            ],
          },
          {
            label: 'test-ns-1-with-namespace-length-over-18-characters/',
            clients: 1245,
            entity_clients: 710,
            non_entity_clients: 535,
            mounts: [
              {
                label: 'path-4-with-over-18-characters',
                clients: 440,
                entity_clients: 323,
                non_entity_clients: 117,
              },
              {
                label: 'path-2',
                clients: 339,
                entity_clients: 308,
                non_entity_clients: 31,
              },
              {
                label: 'path-1',
                clients: 288,
                entity_clients: 63,
                non_entity_clients: 225,
              },
              {
                label: 'path-3-with-over-18-characters',
                clients: 178,
                entity_clients: 16,
                non_entity_clients: 162,
              },
            ],
          },
          {
            label: 'test-ns-1/',
            clients: 820,
            entity_clients: 356,
            non_entity_clients: 464,
            mounts: [
              {
                label: 'path-2',
                clients: 259,
                entity_clients: 245,
                non_entity_clients: 14,
              },
              {
                label: 'path-4-with-over-18-characters',
                clients: 256,
                entity_clients: 63,
                non_entity_clients: 193,
              },
              {
                label: 'path-1',
                clients: 239,
                entity_clients: 30,
                non_entity_clients: 209,
              },
              {
                label: 'path-3-with-over-18-characters',
                clients: 66,
                entity_clients: 18,
                non_entity_clients: 48,
              },
            ],
          },
        ],
      },
    },
  ];
  const NEW_ACTIVITY = MONTHLY_ACTIVITY.map((d) => d.new_clients);
  const TOTAL_USAGE_COUNTS = {
    clients: 38668,
    entity_clients: 20818,
    non_entity_clients: 17850,
  };
  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => new Date('2018-04-03T14:15:30'));
  });
  hooks.beforeEach(function () {
    this.set('timestamp', formatRFC3339(timestamp.now()));
    this.set('chartLegend', [
      { label: 'entity clients', key: 'entity_clients' },
      { label: 'non-entity clients', key: 'non_entity_clients' },
    ]);
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });
  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it renders with full monthly activity data', async function (assert) {
    this.set('byMonthActivityData', MONTHLY_ACTIVITY);
    this.set('totalUsageCounts', TOTAL_USAGE_COUNTS);
    const expectedTotalEntity = formatNumber([TOTAL_USAGE_COUNTS.entity_clients]);
    const expectedTotalNonEntity = formatNumber([TOTAL_USAGE_COUNTS.non_entity_clients]);
    const expectedNewEntity = formatNumber([calculateAverage(NEW_ACTIVITY, 'entity_clients')]);
    const expectedNewNonEntity = formatNumber([calculateAverage(NEW_ACTIVITY, 'non_entity_clients')]);

    await render(hbs`
            <Clients::RunningTotal
      @chartLegend={{this.chartLegend}}
      @selectedAuthMethod={{this.selectedAuthMethod}}
      @byMonthActivityData={{this.byMonthActivityData}}
      @runningTotals={{this.totalUsageCounts}}
      @upgradeData={{this.upgradeDuringActivity}}
      @responseTimestamp={{this.timestamp}}
      @isHistoricalMonth={{false}}
    />
    `);

    assert.dom('[data-test-running-total]').exists('running total component renders');
    assert.dom('[data-test-line-chart]').exists('line chart renders');
    assert.dom('[data-test-vertical-bar-chart]').exists('vertical bar chart renders');
    assert.dom('[data-test-running-total-legend]').exists('legend renders');
    assert.dom('[data-test-running-total-timestamp]').exists('renders timestamp');
    assert
      .dom('[data-test-running-total-entity] p.data-details')
      .hasText(`${expectedTotalEntity}`, `renders correct total average ${expectedTotalEntity}`);
    assert
      .dom('[data-test-running-total-nonentity] p.data-details')
      .hasText(`${expectedTotalNonEntity}`, `renders correct new average ${expectedTotalNonEntity}`);
    assert
      .dom('[data-test-running-new-entity] p.data-details')
      .hasText(`${expectedNewEntity}`, `renders correct total average ${expectedNewEntity}`);
    assert
      .dom('[data-test-running-new-nonentity] p.data-details')
      .hasText(`${expectedNewNonEntity}`, `renders correct new average ${expectedNewNonEntity}`);

    // assert line chart is correct
    findAll('[data-test-line-chart="x-axis-labels"] text').forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${MONTHLY_ACTIVITY[i].month}`,
          `renders x-axis labels for line chart: ${MONTHLY_ACTIVITY[i].month}`
        );
    });
    assert
      .dom('[data-test-line-chart="plot-point"]')
      .exists(
        { count: MONTHLY_ACTIVITY.filter((m) => m.counts !== null).length },
        'renders correct number of plot points'
      );

    // assert bar chart is correct
    findAll('[data-test-vertical-chart="x-axis-labels"] text').forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${MONTHLY_ACTIVITY[i].month}`,
          `renders x-axis labels for bar chart: ${MONTHLY_ACTIVITY[i].month}`
        );
    });
    assert
      .dom('[data-test-vertical-chart="data-bar"]')
      .exists(
        { count: MONTHLY_ACTIVITY.filter((m) => m.counts !== null).length * 2 },
        'renders correct number of data bars'
      );
  });

  test('it renders with no new monthly data', async function (assert) {
    const monthlyWithoutNew = MONTHLY_ACTIVITY.map((d) => ({ ...d, new_clients: { month: d.month } }));
    this.set('byMonthActivityData', monthlyWithoutNew);
    this.set('totalUsageCounts', TOTAL_USAGE_COUNTS);
    const expectedTotalEntity = formatNumber([TOTAL_USAGE_COUNTS.entity_clients]);
    const expectedTotalNonEntity = formatNumber([TOTAL_USAGE_COUNTS.non_entity_clients]);

    await render(hbs`
            <Clients::RunningTotal
      @chartLegend={{this.chartLegend}}
      @selectedAuthMethod={{this.selectedAuthMethod}}
      @byMonthActivityData={{this.byMonthActivityData}}
      @runningTotals={{this.totalUsageCounts}}
      @responseTimestamp={{this.timestamp}}
      @isHistoricalMonth={{false}}
    />
    `);
    assert.dom('[data-test-running-total]').exists('running total component renders');
    assert.dom('[data-test-line-chart]').exists('line chart renders');
    assert.dom('[data-test-vertical-bar-chart]').doesNotExist('vertical bar chart does not render');
    assert.dom('[data-test-running-total-legend]').doesNotExist('legend does not render');
    assert.dom('[data-test-component="empty-state"]').exists('renders empty state');
    assert.dom('[data-test-empty-state-title]').hasText('No new clients');
    assert.dom('[data-test-running-total-timestamp]').exists('renders timestamp');
    assert
      .dom('[data-test-running-total-entity] p.data-details')
      .hasText(`${expectedTotalEntity}`, `renders correct total average ${expectedTotalEntity}`);
    assert
      .dom('[data-test-running-total-nonentity] p.data-details')
      .hasText(`${expectedTotalNonEntity}`, `renders correct new average ${expectedTotalNonEntity}`);
    assert
      .dom('[data-test-running-new-entity] p.data-details')
      .doesNotExist('new client counts does not exist');
    assert
      .dom('[data-test-running-new-nonentity] p.data-details')
      .doesNotExist('average new client counts does not exist');
  });

  test('it renders with single historical month data', async function (assert) {
    const singleMonth = MONTHLY_ACTIVITY[MONTHLY_ACTIVITY.length - 1];
    const singleMonthNew = NEW_ACTIVITY[NEW_ACTIVITY.length - 1];
    this.set('singleMonth', [singleMonth]);
    const expectedTotalClients = formatNumber([singleMonth.clients]);
    const expectedTotalEntity = formatNumber([singleMonth.entity_clients]);
    const expectedTotalNonEntity = formatNumber([singleMonth.non_entity_clients]);
    const expectedNewClients = formatNumber([singleMonthNew.clients]);
    const expectedNewEntity = formatNumber([singleMonthNew.entity_clients]);
    const expectedNewNonEntity = formatNumber([singleMonthNew.non_entity_clients]);

    await render(hbs`
            <Clients::RunningTotal
      @chartLegend={{this.chartLegend}}
      @selectedAuthMethod={{this.selectedAuthMethod}}
      @byMonthActivityData={{this.singleMonth}}
      @runningTotals={{this.totalUsageCounts}}
      @responseTimestamp={{this.timestamp}}
      @isHistoricalMonth={{true}}
    />
    `);
    assert.dom('[data-test-running-total]').exists('running total component renders');
    assert.dom('[data-test-line-chart]').doesNotExist('line chart does not render');
    assert.dom('[data-test-vertical-bar-chart]').doesNotExist('vertical bar chart does not render');
    assert.dom('[data-test-running-total-legend]').doesNotExist('legend does not render');
    assert.dom('[data-test-running-total-timestamp]').doesNotExist('renders timestamp');
    assert.dom('[data-test-stat-text-container]').exists({ count: 6 }, 'renders stat text containers');
    assert
      .dom('[data-test-new] [data-test-stat-text-container="New clients"] div.stat-value')
      .hasText(`${expectedNewClients}`, `renders correct total new clients: ${expectedNewClients}`);
    assert
      .dom('[data-test-new] [data-test-stat-text-container="Entity clients"] div.stat-value')
      .hasText(`${expectedNewEntity}`, `renders correct total new entity: ${expectedNewEntity}`);
    assert
      .dom('[data-test-new] [data-test-stat-text-container="Non-entity clients"] div.stat-value')
      .hasText(`${expectedNewNonEntity}`, `renders correct total new non-entity: ${expectedNewNonEntity}`);
    assert
      .dom('[data-test-total] [data-test-stat-text-container="Total monthly clients"] div.stat-value')
      .hasText(`${expectedTotalClients}`, `renders correct total clients: ${expectedTotalClients}`);
    assert
      .dom('[data-test-total] [data-test-stat-text-container="Entity clients"] div.stat-value')
      .hasText(`${expectedTotalEntity}`, `renders correct total entity: ${expectedTotalEntity}`);
    assert
      .dom('[data-test-total] [data-test-stat-text-container="Non-entity clients"] div.stat-value')
      .hasText(`${expectedTotalNonEntity}`, `renders correct total non-entity: ${expectedTotalNonEntity}`);
  });
});
