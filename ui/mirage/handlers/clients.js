/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  addMonths,
  differenceInCalendarMonths,
  endOfMonth,
  formatRFC3339,
  fromUnixTime,
  isAfter,
  isBefore,
  isSameMonth,
  isWithinInterval,
  startOfMonth,
  subMonths,
} from 'date-fns';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { CLIENT_TYPES } from 'core/utils/client-count-utils';

/*
HOW TO ADD NEW TYPES:
1. add key to CLIENT_TYPES 
2. Find "ADD NEW CLIENT TYPES HERE" comment below and generate mock counts for that key
3. Add generateMounts() for that client type to the mounts array
*/
export const LICENSE_START = new Date('2023-07-02T00:00:00Z');
export const STATIC_NOW = new Date('2024-01-25T23:59:59Z');
const COUNTS_START = subMonths(STATIC_NOW, 12); // user started Vault cluster on 2023-01-25
// upgrade happened 2 month after license start
export const UPGRADE_DATE = addMonths(LICENSE_START, 2); // monthly attribution added

function getSum(array, key) {
  return array.reduce((sum, { counts }) => sum + counts[key], 0);
}

function getTotalCounts(array) {
  const counts = CLIENT_TYPES.reduce((obj, key) => {
    obj[key] = getSum(array, key);
    return obj;
  }, {});

  // add deprecated keys
  return {
    ...counts,
    distinct_entities: counts.entity_clients,
    non_entity_tokens: counts.non_entity_clients,
  };
}

function randomBetween(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

// generates array of counts that add up to max
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

function generateMounts(pathPrefix, counts) {
  const baseObject = CLIENT_TYPES.reduce((obj, key) => {
    obj[key] = 0;
    return obj;
  }, {});
  return Array.from(Array(5)).map((mount, index) => {
    return {
      mount_path: `${pathPrefix}${index}`,
      counts: {
        ...baseObject,
        distinct_entities: 0,
        non_entity_tokens: 0,
        // object contains keys for which 0-values of base object to overwrite
        ...counts,
      },
    };
  });
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

  // * ADD NEW CLIENT TYPES HERE and spread to the mounts array below
  const authClients = randomBetween(min, max);
  const [non_entity_clients, entity_clients] = arrayOfCounts(authClients, 2);
  const secret_syncs = randomBetween(min, max);
  const acme_clients = randomBetween(min, max);

  // each mount type generates a different type of client
  const mounts = [
    ...generateMounts('auth/authid/', { clients: authClients, non_entity_clients, entity_clients }),
    ...generateMounts('kvv2-engine-', { clients: secret_syncs, secret_syncs }),
    ...generateMounts('pki-engine-', { clients: acme_clients, acme_clients }),
  ];

  mounts.sort((a, b) => b.counts.clients - a.counts.clients);
  nsBlock.mounts = mounts;
  nsBlock.counts = getTotalCounts(mounts);
  return nsBlock;
}

function generateMonths(startDate, endDate, namespaces) {
  const startDateObject = parseAPITimestamp(startDate);
  const endDateObject = parseAPITimestamp(endDate);
  const numberOfMonths = differenceInCalendarMonths(endDateObject, startDateObject) + 1;
  const months = [];

  // only generate monthly block if queried dates span an upgrade
  if (isWithinInterval(UPGRADE_DATE, { start: startDateObject, end: endDateObject })) {
    for (let i = 0; i < numberOfMonths; i++) {
      const month = addMonths(startOfMonth(startDateObject), i);
      const hasNoData = isBefore(month, UPGRADE_DATE) && !isSameMonth(month, UPGRADE_DATE);
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
          expiration_time: formatRFC3339(endOfMonth(addMonths(STATIC_NOW, 6))),
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
            timestamp_installed: LICENSE_START.toISOString(),
          },
          '1.9.1': {
            build_date: null,
            previous_version: '1.9.0',
            timestamp_installed: addMonths(LICENSE_START, 1).toISOString(),
          },
          // auth mount attribution added in 1.10.0
          '1.10.1': {
            build_date: null,
            previous_version: '1.9.1',
            timestamp_installed: UPGRADE_DATE.toISOString(),
          },
          // no notable UI changes
          '1.14.4': {
            build_date: addMonths(LICENSE_START, 3).toISOString(),
            previous_version: '1.10.1',
            timestamp_installed: addMonths(LICENSE_START, 3).toISOString(),
          },
          // sync clients added
          '1.16.0': {
            build_date: addMonths(LICENSE_START, 4).toISOString(),
            previous_version: '1.14.4',
            timestamp_installed: addMonths(LICENSE_START, 4).toISOString(),
          },
        },
      },
    };
  });
}
