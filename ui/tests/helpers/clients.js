import { addMonths, differenceInCalendarMonths, formatRFC3339, startOfMonth } from 'date-fns';
import { Response } from 'miragejs';

/** Scenarios
  Config off, no data
  Config on, no data
  Config on, with data
  Filtering (data with mounts)
  Filtering (data without mounts)
  Filtering (data without mounts)

 * -- HISTORY ONLY --
  No permissions for license
  Version
  queries available
  queries unavailable
  License start date this month
 */

// TODO
/*
Filtering different date ranges (hist only)
Upgrade warning 

*/
export const SELECTORS = {
  currentMonthActiveTab: '.active[data-test-current-month]',
  historyActiveTab: '.active[data-test-history]',
  emptyStateTitle: '[data-test-empty-state-title]',
  usageStats: '[data-test-usage-stats]',
  dateDisplay: '[data-test-date-display]',
  attributionBlock: '[data-test-clients-attribution]',
  filterBar: '[data-test-clients-filter-bar]',
  rangeDropdown: '[data-test-popup-menu-trigger]',
  monthDropdown: '[data-test-popup-menu-trigger="month"]',
  yearDropdown: '[data-test-popup-menu-trigger="year"]',
  dateDropdownSubmit: '[data-test-date-dropdown-submit]',
  runningTotalMonthStats: '[data-test-running-total="single-month-stats"]',
  runningTotalMonthlyCharts: '[data-test-running-total="monthly-charts"]',
  monthlyUsageBlock: '[data-test-monthly-usage]',
};

export const CHART_ELEMENTS = {
  entityClientDataBars: '[data-test-group="entity_clients"]',
  nonEntityDataBars: '[data-test-group="non_entity_clients"]',
  yLabels: '[data-test-group="y-labels"]',
  actionBars: '[data-test-group="action-bars"]',
  labelActionBars: '[data-test-group="label-action-bars"]',
  totalValues: '[data-test-group="total-values"]',
};

export function sendResponse(data, httpStatus = 200) {
  if (httpStatus === 403) {
    return [
      httpStatus,
      { 'Content-Type': 'application/json' },
      JSON.stringify({ errors: ['permission denied'] }),
    ];
  }
  if (httpStatus === 204) {
    // /activity endpoint returns 204 when no data, while
    // /activity/monthly returns 200 with zero values on data
    return [httpStatus, { 'Content-Type': 'application/json' }];
  }
  return [httpStatus, { 'Content-Type': 'application/json' }, JSON.stringify(data)];
}

export function overrideResponse(httpStatus, data) {
  if (httpStatus === 403) {
    return new Response(
      403,
      { 'Content-Type': 'application/json' },
      JSON.stringify({ errors: ['permission denied'] })
    );
  }
  // /activity endpoint returns 204 when no data, while
  // /activity/monthly returns 200 with zero values on data
  if (httpStatus === 204) {
    return new Response(204, { 'Content-Type': 'application/json' });
  }
  return new Response(200, { 'Content-Type': 'application/json' }, JSON.stringify(data));
}

function generateNamespaceBlock(idx = 0, skipMounts = false) {
  let mountCount = 1;
  const nsBlock = {
    namespace_id: `${idx}UUID`,
    namespace_path: `${idx}/namespace`,
    counts: {
      clients: mountCount * 15,
      entity_clients: mountCount * 5,
      non_entity_clients: mountCount * 10,
      distinct_entities: mountCount * 5,
      non_entity_tokens: mountCount * 10,
    },
  };
  if (!skipMounts) {
    mountCount = Math.floor((Math.random() + idx) * 20);
    let mounts = [];
    Array.from(Array(mountCount)).forEach((v, index) => {
      mounts.push({
        mount_path: `auth/authid${index}`,
        counts: {
          clients: 5,
          entity_clients: 3,
          non_entity_clients: 2,
          distinct_entities: 3,
          non_entity_tokens: 2,
        },
      });
    });
    nsBlock.mounts = mounts;
  }
  return nsBlock;
}

function generateCounts(max, arrayLength) {
  function randomBetween(min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min);
  }
  var result = [];
  var sum = 0;
  for (var i = 0; i < arrayLength - 1; i++) {
    result[i] = randomBetween(1, max - (arrayLength - i - 1) - sum);
    sum += result[i];
  }
  result[arrayLength - 1] = max - sum;
  return result.sort((a, b) => b - a);
}

function generateMonths(startDate, endDate, hasNoData = false) {
  let numberOfMonths = differenceInCalendarMonths(endDate, startDate) + 1;
  let months = [];

  for (let i = 0; i < numberOfMonths; i++) {
    if (hasNoData) {
      months.push({
        timestamp: formatRFC3339(startOfMonth(addMonths(startDate, i))),
        counts: null,
        namespace: null,
        new_clients: null,
      });
      continue;
    }
    const namespaces = Array.from(Array(5)).map((v, idx) => {
      return generateNamespaceBlock(idx);
    });
    const clients = numberOfMonths * 5 + i * 5;
    const [entity_clients, non_entity_clients] = generateCounts(clients, 2);
    const counts = {
      clients,
      entity_clients,
      non_entity_clients,
      distinct_entities: entity_clients,
      non_entity_tokens: non_entity_clients,
    };
    const new_counts = 5 + i;
    const [new_entity, new_non_entity] = generateCounts(new_counts, 2);
    months.push({
      timestamp: formatRFC3339(startOfMonth(addMonths(startDate, i))),
      counts,
      namespaces,
      new_clients: {
        counts: {
          distinct_entities: new_entity,
          entity_clients: new_entity,
          non_entity_tokens: new_non_entity,
          non_entity_clients: new_non_entity,
          clients: new_counts,
        },
        namespaces,
      },
    });
  }
  return months;
}

export function generateActivityResponse(nsCount = 1, startDate, endDate) {
  if (nsCount === 0) {
    return {
      request_id: 'some-activity-id',
      data: {
        start_time: formatRFC3339(startDate),
        end_time: formatRFC3339(endDate),
        total: {
          clients: 0,
          entity_clients: 0,
          non_entity_clients: 0,
        },
        by_namespace: [
          {
            namespace_id: `root`,
            namespace_path: '',
            counts: {
              entity_clients: 0,
              non_entity_clients: 0,
              clients: 0,
            },
          },
        ],
        months: generateMonths(startDate, endDate, false),
      },
    };
  }
  let namespaces = Array.from(Array(nsCount)).map((v, idx) => {
    return generateNamespaceBlock(idx);
  });
  return {
    request_id: 'some-activity-id',
    data: {
      start_time: formatRFC3339(startDate),
      end_time: formatRFC3339(endDate),
      total: {
        clients: 999,
        entity_clients: 666,
        non_entity_clients: 333,
      },
      by_namespace: namespaces,
      months: generateMonths(startDate, endDate),
    },
  };
}

export function generateCurrentMonthResponse(namespaceCount, skipMounts = false, configEnabled = true) {
  if (!configEnabled) {
    return {
      data: { id: 'no-data' },
    };
  }
  if (!namespaceCount) {
    return {
      request_id: 'monthly-response-id',
      data: {
        by_namespace: [],
        clients: 0,
        distinct_entities: 0,
        entity_clients: 0,
        non_entity_clients: 0,
        non_entity_tokens: 0,
        months: [],
      },
    };
  }
  // generate by_namespace data
  const by_namespace = Array.from(Array(namespaceCount)).map((ns, idx) =>
    generateNamespaceBlock(idx, skipMounts)
  );
  const counts = by_namespace.reduce(
    (prev, curr) => {
      return {
        clients: prev.clients + curr.counts.clients,
        entity_clients: prev.entity_clients + curr.counts.entity_clients,
        non_entity_clients: prev.non_entity_clients + curr.counts.non_entity_clients,
      };
    },
    { clients: 0, entity_clients: 0, non_entity_clients: 0 }
  );
  return {
    request_id: 'monthly-response-id',
    data: {
      by_namespace,
      ...counts,
      months: [],
    },
  };
}
