import { formatRFC3339 } from 'date-fns';

/** Scenarios
 * Config off, no data
 * * queries available (hist only)
 * * queries unavailable (hist only)
 * Config on, no data
 * Config on, with data
 * Filtering (data with mounts)
 * Filtering (data without mounts)
 * -- HISTORY ONLY --
 * No permissions for license
 * Version (hist only)
 * License start date this month
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

export function generateConfigResponse(overrides = {}) {
  return {
    request_id: 'some-config-id',
    data: {
      default_report_months: 12,
      enabled: 'default-enable',
      queries_available: true,
      retention_months: 24,
      ...overrides,
    },
  };
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
        months: [],
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
      months: [],
    },
  };
}

export function generateLicenseResponse(startDate, endDate) {
  return {
    request_id: 'my-license-request-id',
    data: {
      autoloaded: {
        license_id: 'my-license-id',
        start_time: formatRFC3339(startDate),
        expiration_time: formatRFC3339(endDate),
      },
    },
  };
}

export function generateCurrentMonthResponse(namespaceCount, skipMounts = false) {
  if (!namespaceCount) {
    return {
      request_id: 'monthly-response-id',
      data: {
        by_namespace: [],
        clients: 0,
        entity_clients: 0,
        non_entity_clients: 0,
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
    },
  };
}
