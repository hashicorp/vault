/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  add,
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
2. Find "ADD NEW CLIENT TYPES HERE" comment below and add type to destructuring array
3. Add generateMounts() for that client type to the mounts array
*/
export const LICENSE_START = new Date('2023-07-02T00:00:00Z');
export const STATIC_NOW = new Date('2024-01-25T23:59:59Z');
const COUNTS_START = subMonths(STATIC_NOW, 12); // user started Vault cluster on 2023-01-25
// upgrade happened 2 month after license start
export const UPGRADE_DATE = addMonths(LICENSE_START, 2); // monthly attribution added

// exported so that tests not using this scenario can use the same response
export const CONFIG_RESPONSE = {
  request_id: 'some-config-id',
  data: {
    billing_start_timestamp: formatRFC3339(LICENSE_START),
    default_report_months: 12,
    enabled: 'default-enabled',
    minimum_retention_months: 48,
    queries_available: false,
    reporting_enabled: true,
    retention_months: 48,
  },
};

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

function generateMountBlock(path, counts) {
  const baseObject = CLIENT_TYPES.reduce((obj, key) => {
    obj[key] = 0;
    return obj;
  }, {});
  return {
    mount_path: path,
    counts: {
      ...baseObject,
      distinct_entities: 0,
      non_entity_tokens: 0,
      // object contains keys for which 0-values of base object to overwrite
      ...counts,
    },
  };
}

function generateNamespaceBlock(idx = 0, isLowerCounts = false, ns, skipCounts = false) {
  const min = isLowerCounts ? 10 : 50;
  const max = isLowerCounts ? 100 : 1000;
  const nsBlock = {
    namespace_id: ns?.namespace_id || (idx === 0 ? 'root' : Math.random().toString(36).slice(2, 7) + idx),
    namespace_path: ns?.namespace_path || (idx === 0 ? '' : `ns${idx}`),
    counts: {},
    mounts: {},
  };

  if (skipCounts) return nsBlock; // skip counts to generate empty ns block with namespace ids and paths

  // generates one mount per client type
  const mountsArray = (idx) => {
    // * ADD NEW CLIENT TYPES HERE and pass to a new generateMountBlock() function below
    const [acme_clients, entity_clients, non_entity_clients, secret_syncs] = CLIENT_TYPES.map(() =>
      randomBetween(min, max)
    );

    // each mount type generates a different type of client
    return [
      generateMountBlock(`auth/authid/${idx}`, {
        clients: non_entity_clients + entity_clients,
        non_entity_clients,
        entity_clients,
      }),
      generateMountBlock(`kvv2-engine-${idx}`, { clients: secret_syncs, secret_syncs }),
      generateMountBlock(`pki-engine-${idx}`, { clients: acme_clients, acme_clients }),
    ];
  };

  // two mounts per client type for more varied mock data
  const mounts = [...mountsArray(0), ...mountsArray(1)];

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

  // only generate monthly block if queried dates span or follow upgrade to 1.10
  const upgradeWithin = isWithinInterval(UPGRADE_DATE, { start: startDateObject, end: endDateObject });
  const upgradeAfter = isAfter(startDateObject, UPGRADE_DATE);
  if (upgradeWithin || upgradeAfter) {
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

function generateActivityResponse(startDate, endDate) {
  let namespaces = Array.from(Array(12)).map((v, idx) => generateNamespaceBlock(idx, null, null, true));
  const months = generateMonths(startDate, endDate, namespaces);
  if (months.length) {
    const monthlyCounts = months.filter((m) => m.counts);
    // sum namespace counts from monthly data
    namespaces.forEach((ns) => {
      const nsData = monthlyCounts.map((d) =>
        d.namespaces.find((n) => n.namespace_path === ns.namespace_path)
      );
      const mountCounts = nsData.flatMap((d) => d.mounts);
      const paths = nsData[0].mounts.map(({ mount_path }) => mount_path);
      ns.mounts = paths.map((path) => {
        const counts = getTotalCounts(mountCounts.filter((m) => m.mount_path === path));
        return { mount_path: path, counts };
      });
      ns.counts = getTotalCounts(nsData);
    });
  } else {
    // if no monthly data due to upgrade stuff, generate counts
    namespaces = Array.from(Array(12)).map((v, idx) => generateNamespaceBlock(idx));
  }
  namespaces.sort((a, b) => b.counts.clients - a.counts.clients);
  return {
    start_time: isAfter(new Date(startDate), COUNTS_START) ? startDate : formatRFC3339(COUNTS_START),
    end_time: endDate,
    by_namespace: namespaces,
    months,
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
    return CONFIG_RESPONSE;
  });

  server.get('/sys/internal/counters/activity', (schema, req) => {
    let { start_time, end_time } = req.queryParams;
    if (req.queryParams.current_billing_period) {
      // { current_billing_period: true } automatically queries the activity log
      // from the builtin license start timestamp to the current month
      start_time = LICENSE_START.toISOString();
      end_time = STATIC_NOW.toISOString();
    }
    // backend returns a timestamp if given unix time, so first convert to timestamp string here
    if (!start_time.includes('T')) start_time = fromUnixTime(start_time).toISOString();
    if (!end_time.includes('T')) end_time = fromUnixTime(end_time).toISOString();
    return {
      request_id: 'some-activity-id',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: generateActivityResponse(start_time, end_time),
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
        keys: ['1.9.0', '1.9.1', '1.10.1', '1.10.3', '1.14.4', '1.16.0', '1.17.0'],
        key_info: {
          // entity/non-entity breakdown added
          '1.9.0': {
            // we don't currently use build_date, including for accuracy. it's only tracked in versions >= 1.11.0
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
            timestamp_installed: addMonths(LICENSE_START, 2).toISOString(), // same as UPGRADE_DATE
          },
          '1.10.3': {
            build_date: null,
            previous_version: '1.10.1',
            timestamp_installed: add(LICENSE_START, { months: 2, weeks: 3 }).toISOString(),
          },
          // no notable UI changes
          '1.14.4': {
            build_date: addMonths(LICENSE_START, 3).toISOString(),
            previous_version: '1.10.3',
            timestamp_installed: addMonths(LICENSE_START, 3).toISOString(),
          },
          // sync clients added
          '1.16.0': {
            build_date: addMonths(LICENSE_START, 4).toISOString(),
            previous_version: '1.14.4',
            timestamp_installed: addMonths(LICENSE_START, 4).toISOString(),
          },
          // acme_clients separated from non-entity clients
          '1.17.0': {
            build_date: addMonths(LICENSE_START, 5).toISOString(),
            previous_version: '1.16.0',
            timestamp_installed: addMonths(LICENSE_START, 5).toISOString(),
          },
        },
      },
    };
  });
}
