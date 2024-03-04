/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  isBefore,
  startOfMonth,
  endOfMonth,
  addMonths,
  subMonths,
  differenceInCalendarMonths,
  fromUnixTime,
  isAfter,
  formatRFC3339,
} from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';

export const STATIC_START = new Date('2023-10-01T00:00:00Z');
export const STATIC_NOW = new Date('2024-01-25T23:59:59Z');
// Matches mocked date in client-dashboard-test file
const CURRENT_DATE = new Date('2023-01-13T14:15:00');
const COUNTS_START = subMonths(CURRENT_DATE, 12); // pretend vault user started cluster 6 months ago
// for testing, we're in the middle of a license/billing period
const LICENSE_START = startOfMonth(subMonths(CURRENT_DATE, 6));
// upgrade happened 1 month after license start
const UPGRADE_DATE = addMonths(LICENSE_START, 1);

function getSum(array, key) {
  return array.reduce((sum, { counts }) => sum + counts[key], 0);
}

function getTotalCounts(array) {
  return {
    distinct_entities: getSum(array, 'entity_clients'),
    entity_clients: getSum(array, 'entity_clients'),
    non_entity_tokens: getSum(array, 'non_entity_clients'),
    non_entity_clients: getSum(array, 'non_entity_clients'),
    secret_syncs: getSum(array, 'secret_syncs'),
    clients: getSum(array, 'clients'),
  };
}

function randomBetween(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

function arrayOfCounts(max, arrayLength) {
  var result = [];
  var sum = 0;
  for (var i = 0; i < arrayLength - 1; i++) {
    result[i] = randomBetween(1, max - (arrayLength - i - 1) - sum);
    sum += result[i];
  }
  result[arrayLength - 1] = max - sum;
  return result.sort((a, b) => b - a);
}

function generateNamespaceBlock(idx = 0, isLowerCounts = false, ns) {
  const min = isLowerCounts ? 10 : 50;
  const max = isLowerCounts ? 100 : 5000;
  const nsBlock = {
    namespace_id: ns?.namespace_id || (idx === 0 ? 'root' : Math.random().toString(36).slice(2, 7) + idx),
    namespace_path: ns?.namespace_path || (idx === 0 ? '' : `ns/${idx}`),
    counts: {},
    mounts: {},
  };
  const mounts = [];

  Array.from(Array(5)).forEach((mount, index) => {
    const [secretSyncs] = arrayOfCounts(randomBetween(min, max), 1);
    mounts.push({
      mount_path: `kvv2-engine-${index}`,
      counts: {
        clients: secretSyncs,
        // TODO test with live backend to confirm entity keys are present (and 0) for kv mounts
        entity_clients: 0,
        non_entity_clients: 0,
        distinct_entities: 0,
        non_entity_tokens: 0,
        secret_syncs: secretSyncs,
      },
    });
  });

  // generate auth mounts array
  Array.from(Array(10)).forEach((mount, index) => {
    const mountClients = randomBetween(min, max);
    const [nonEntity, entity] = arrayOfCounts(mountClients, 2);
    mounts.push({
      mount_path: `auth/authid${index}`,
      counts: {
        clients: mountClients,
        entity_clients: entity,
        non_entity_clients: nonEntity,
        distinct_entities: entity,
        non_entity_tokens: nonEntity,
        // TODO test with live backend to confirm this key is present (and 0) for auth mounts (non-kv mounts)
        secret_syncs: 0,
      },
    });
  });
  mounts.sort((a, b) => b.counts.clients - a.counts.clients);
  nsBlock.mounts = mounts;
  nsBlock.counts = getTotalCounts(mounts);
  return nsBlock;
}

function generateMonths(startDate, endDate, namespaces) {
  const startDateObject = startOfMonth(parseAPITimestamp(startDate));
  const endDateObject = startOfMonth(parseAPITimestamp(endDate));
  const numberOfMonths = differenceInCalendarMonths(endDateObject, startDateObject) + 1;
  const months = [];
  if (isBefore(startDateObject, UPGRADE_DATE) && isBefore(endDateObject, UPGRADE_DATE)) {
    // months block is empty if dates do not span an upgrade
    return [];
  }
  for (let i = 0; i < numberOfMonths; i++) {
    const month = addMonths(startDateObject, i);
    const hasNoData = isBefore(month, UPGRADE_DATE);
    if (hasNoData) {
      months.push({
        timestamp: formatRFC3339(month),
        counts: null,
        namespaces: null,
        new_clients: null,
      });
      continue;
    }

    const monthNs = namespaces.map((ns, idx) => generateNamespaceBlock(idx, false, ns));
    const newClients = namespaces.map((ns, idx) => generateNamespaceBlock(idx, true, ns));
    months.push({
      timestamp: formatRFC3339(month),
      counts: getTotalCounts(monthNs),
      namespaces: monthNs.sort((a, b) => b.counts.clients - a.counts.clients),
      new_clients: {
        counts: getTotalCounts(newClients),
        namespaces: newClients.sort((a, b) => b.counts.clients - a.counts.clients),
      },
    });
  }
  return months;
}

function generateActivityResponse(namespaces, startDate, endDate) {
  return {
    start_time: isAfter(new Date(startDate), COUNTS_START) ? startDate : formatRFC3339(COUNTS_START),
    end_time: endDate,
    by_namespace: namespaces.sort((a, b) => b.counts.clients - a.counts.clients),
    months: generateMonths(startDate, endDate, namespaces),
    total: getTotalCounts(namespaces),
  };
}

export default function (server) {
  server.get('sys/license/status', function () {
    return {
      request_id: 'my-license-request-id',
      data: {
        autoloaded: {
          license_id: 'my-license-id',
          start_time: formatRFC3339(LICENSE_START),
          expiration_time: formatRFC3339(endOfMonth(addMonths(CURRENT_DATE, 6))),
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
        billing_start_timestamp: formatRFC3339(LICENSE_START),
      },
    };
  });

  server.get('/sys/internal/counters/activity', (schema, req) => {
    let { start_time, end_time } = req.queryParams;
    // backend returns a timestamp if given unix time, so first convert to timestamp string here
    if (!start_time.includes('T')) start_time = fromUnixTime(start_time).toISOString();
    if (!end_time.includes('T')) end_time = fromUnixTime(end_time).toISOString();
    const namespaces = Array.from(Array(12)).map((v, idx) => generateNamespaceBlock(idx));
    return {
      request_id: 'some-activity-id',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: generateActivityResponse(namespaces, start_time, end_time),
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });

  // client counting has changed in different ways since 1.9 see link below for details
  // https://developer.hashicorp.com/vault/docs/concepts/client-count/faq#client-count-faq
  server.get('sys/version-history', function () {
    return {
      request_id: 'version-history-request-id',
      data: {
        keys: ['1.9.0', '1.9.1', '1.10.1', '1.14.4', '1.16.0'],
        key_info: {
          // entity/non-entity breakdown added
          '1.9.0': {
            // we don't currently use build_date, including for accuracy. it's only tracked in versions >= 1.110
            build_date: null,
            previous_version: null,
            timestamp_installed: addMonths(STATIC_START, 1).toISOString(),
          },
          '1.9.1': {
            build_date: null,
            previous_version: '1.9.0',
            timestamp_installed: addMonths(STATIC_START, 2).toISOString(),
          },
          // auth mount attribution added in 1.10.0
          '1.10.1': {
            build_date: null,
            previous_version: '1.9.1',
            timestamp_installed: addMonths(STATIC_START, 3).toISOString(),
          },
          // no notable UI changes
          '1.14.4': {
            build_date: addMonths(STATIC_START, 3).toISOString(),
            previous_version: '1.10.1',
            timestamp_installed: addMonths(STATIC_START, 4).toISOString(),
          },
          // sync clients added
          '1.16.0': {
            build_date: addMonths(STATIC_START, 4).toISOString(),
            previous_version: '1.14.4',
            timestamp_installed: addMonths(STATIC_START, 5).toISOString(),
          },
        },
      },
    };
  });
}
