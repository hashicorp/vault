/** Scenarios
 * Config off, no data, queries unavailable/available (hist only)
 * Config on, no data
 * Config on, with data
 * Version (hist only)
 * Filtering (data with mounts)
 * Filtering (data without mounts)
 *
 */

export const SELECTORS = {
  activeTab: '.nav-tab-link.is-active',
  emptyStateTitle: '[data-test-empty-state-title]',
  usageStats: '[data-test-usage-stats]',
  dateDisplay: '[data-test-date-display]',
  attributionBlock: '[data-test-clients-attribution]',
};

export function sendResponse(data, httpStatus = 200) {
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

export function generateNamespaceBlock(idx = 0, skipMounts = false) {
  let mountCount = 1;
  const nsBlock = {
    namespace_id: `${idx}UUID`,
    namespace_path: `my-namespace-${idx}/`,
    counts: {
      entity_clients: mountCount * 5,
      non_entity_clients: mountCount * 10,
      clients: mountCount * 15,
    },
  };
  if (!skipMounts) {
    mountCount = Math.floor((Math.random() + idx) * 20);
    let mounts = [];
    if (!skipMounts) {
      Array.from(Array(mountCount)).forEach((v, index) => {
        mounts.push({
          id: index,
          path: `auth/method/authid${index}`,
          counts: {
            clients: 5,
            entity_clients: 3,
            non_entity_clients: 2,
          },
        });
      });
    }
    nsBlock.mounts = mounts;
  }
  return nsBlock;
}
