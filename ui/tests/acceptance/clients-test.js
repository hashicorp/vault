import { module, test, skip } from 'qunit';
import { visit, currentURL, settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import authPage from 'vault/tests/pages/auth';
import { addMonths, format, formatRFC3339, startOfMonth, subMonths } from 'date-fns';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';

const searchSelect = create(ss);

function generateNamespaceBlock(idx = 0, skipMounts = false) {
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

function generateConfigResponse(overrides = {}) {
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
function generateActivityResponse(nsCount = 1, startDate, endDate) {
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
        // months: [],
      },
    };
  }
  let namespaces = Array.from(Array(nsCount)).map((v, idx) => {
    return generateNamespaceBlock(idx);
  });
  console.log({ namespaces });
  return {
    request_id: 'some-activity-id',
    data: {
      start_time: formatRFC3339(startDate),
      end_time: formatRFC3339(endDate),
      total: {
        clients: 3637,
        entity_clients: 1643,
        non_entity_clients: 1994,
      },
      by_namespace: namespaces,
      // months: [],
    },
  };
}
function generateCurrentMonthResponse(namespaceCount) {
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
  const by_namespace = Array.from(Array(namespaceCount)).map((ns, idx) => generateNamespaceBlock(idx));
  const counts = by_namespace.reduce(
    (prev, curr) => {
      console.log(prev, curr, 'combine namespaces');
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
function generateLicenseResponse(startDate, endDate) {
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
function send(data, httpStatus = 200) {
  return [httpStatus, { 'Content-Type': 'application/json' }, JSON.stringify(data)];
}

const SELECTORS = {
  activeTab: '.nav-tab-link.is-active',
  emptyStateTitle: '[data-test-empty-state-title]',
  usageStats: '[data-test-usage-stats]',
  dateDisplay: '[data-test-date-display]',
  attributionBlock: '[data-test-clients-attribution]',
};

module('Acceptance | clients route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('shows empty state and warning when config disabled, queries available, no data', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse({ enabled: 'default-disable' });
    const monthly = generateCurrentMonthResponse();
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => send(license));
      this.get('/v1/sys/internal/counters/activity', () => send(activity));
      this.get('/v1/sys/internal/counters/activity/monthly', () => send(monthly));
      this.get('/v1/sys/internal/counters/config', () => send(config));
      this.get('/v1/sys/version-history', () => send({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
    });
    // Current Tab
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.activeTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('Tracking is disabled');
    // History Tab
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.activeTab).hasText('History', 'history tab is active');

    assert.dom('[data-test-tracking-disabled] .message-title').hasText('Tracking is disabled');
    // TODO: still allows query by previous dates
  });

  test('shows empty state and warning when no queries', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse({ queries_available: false });
    const monthly = generateCurrentMonthResponse();
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => send(license));
      this.get('/v1/sys/internal/counters/activity', () => send(activity));
      this.get('/v1/sys/internal/counters/activity/monthly', () => send(monthly));
      this.get('/v1/sys/internal/counters/config', () => send(config));
      this.get('/v1/sys/version-history', () => send({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    // Current tab
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.activeTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No data received');
    // History Tab
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.activeTab).hasText('History', 'history tab is active');

    assert.dom(SELECTORS.emptyStateTitle).hasText('No monthly history');
  });
  test('visiting history tab with no data and config on', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const config = generateConfigResponse();
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => send(license));
      this.get('/v1/sys/internal/counters/activity', () => send(activity));
      this.get('/v1/sys/internal/counters/config', () => send(config));
      this.get('/v1/sys/version-history', () => send({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert
      .dom(SELECTORS.dateDisplay)
      .hasText(format(licenseStart, 'MMMM yyyy'), 'billing start month is correctly parsed from license');
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Attribution block is not shown when no data');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    // TODO: Filters correct
    // TODO: don't show namespace filter if none exist
  });
  test('filters correctly on current with full data', async function (assert) {
    const config = generateConfigResponse();
    const monthly = generateCurrentMonthResponse(3);
    this.server = new Pretender(function () {
      this.get('/v1/sys/internal/counters/activity/monthly', () => send(monthly));
      this.get('/v1/sys/internal/counters/config', () => send(config));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.activeTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.usageStats).exists('usage stats block exists');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    const { clients, entity_clients, non_entity_clients } = monthly.data;
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText(clients.toString());
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText(entity_clients.toString());
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText(non_entity_clients.toString());
    assert.dom('[data-test-clients-attribution]').exists('Shows attribution area');
    assert.dom('[data-test-horizontal-bar-chart]').exists('Shows attribution bar chart');
    assert.dom('[data-test-top-attribution]').hasText('Top namespace');
    // Filter by namespace
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('15');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('5');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('10');
    assert.dom('[data-test-horizontal-bar-chart]').exists('Still shows attribution bar chart');
    assert.dom('[data-test-top-attribution]').hasText('Top auth method');
    // Filter by auth method
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('5');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('3');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('2');
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution block');
  });
  test('filters correctly on history with full data', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const config = generateConfigResponse();
    const activity = generateActivityResponse(5, licenseStart, licenseEnd);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => send(license));
      this.get('/v1/sys/internal/counters/activity', () => send(activity));
      this.get('/v1/sys/internal/counters/config', () => send(config));
      this.get('/v1/sys/version-history', () => send({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    console.log({ activity });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.activeTab).hasText('History', 'history tab is active');
    assert.dom(SELECTORS.usageStats).exists('usage stats block exists');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    const { clients, entity_clients, non_entity_clients } = activity.data.total;
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText(clients.toString());
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText(entity_clients.toString());
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText(non_entity_clients.toString());
    await this.pauseTest();
    assert.dom('[data-test-clients-attribution]').exists('Shows attribution area');
    assert.dom('[data-test-horizontal-bar-chart]').exists('Shows attribution bar chart');
    assert.dom('[data-test-top-attribution]').hasText('Top namespace');
    // Filter by namespace
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('15');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('5');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('10');
    assert.dom('[data-test-horizontal-bar-chart]').exists('Still shows attribution bar chart');
    assert.dom('[data-test-top-attribution]').hasText('Top auth method');
    // Filter by auth method
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('5');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('2');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('3');
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution block');
    await this.pauseTest();
    await click('#allowed_roles [data-test-selected-list-button="delete"]');
  });
});
