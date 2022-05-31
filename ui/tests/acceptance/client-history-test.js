import { module, test } from 'qunit';
import { visit, currentURL, click } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import authPage from 'vault/tests/pages/auth';
import { addMonths, format, formatRFC3339, startOfMonth, subMonths } from 'date-fns';
import {
  generateActivityResponse,
  generateConfigResponse,
  generateLicenseResponse,
  SELECTORS,
  sendResponse,
} from '../helpers/clients';
// import { create } from 'ember-cli-page-object';
// import { clickTrigger } from 'ember-power-select/test-support/helpers';
// import ss from 'vault/tests/pages/components/search-select';

// const searchSelect = create(ss);

module('Acceptance | clients history tab', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('shows warning when config off, no data, queries available', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse({ enabled: 'default-disable' });
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(null, 204));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');

    assert.dom('[data-test-tracking-disabled] .message-title').hasText('Tracking is disabled');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No data received');
    assert.dom(SELECTORS.filterBar).doesNotExist('Shows filter bar to search previous dates');
    assert.dom(SELECTORS.usageStats).doesNotExist('No usage stats');
  });

  test('shows warning when config off, no data, queries unavailable', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse({ enabled: 'default-disable', queries_available: false });
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(null, 204));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('Data tracking is disabled');
    assert.dom(SELECTORS.filterBar).doesNotExist('Filter bar is hidden when no data available');
  });

  test('shows empty state when config on and no queries', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse({ queries_available: false });
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(null, 204));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    // History Tab
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');

    assert.dom(SELECTORS.emptyStateTitle).hasText('No monthly history');
    assert.dom(SELECTORS.filterBar).doesNotExist('Does not show filter bar');
  });

  test('visiting history tab config on and data with mounts', async function (assert) {
    assert.expect(9);
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const lastMonth = addMonths(new Date(), -1);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse();
    const activity = generateActivityResponse(5, licenseStart, lastMonth);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(activity));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
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
    assert
      .dom(SELECTORS.rangeDropdown)
      .hasText(
        `${format(licenseStart, 'MMMM yyyy')} - ${format(lastMonth, 'MMMM yyyy')}`,
        'Date range shows dates correctly parsed activity response'
      );
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    const { clients, entity_clients, non_entity_clients } = activity.data.total;
    assert
      .dom('[data-test-stat-text="total-clients"] .stat-value')
      .hasText(clients.toString(), 'total clients stat is correct');
    assert
      .dom('[data-test-stat-text="entity-clients"] .stat-value')
      .hasText(entity_clients.toString(), 'entity clients stat is correct');
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText(non_entity_clients.toString(), 'non-entity clients stat is correct');
    assert.dom('[data-test-clients-attribution]').exists('Shows attribution area');
    assert.dom('[data-test-top-attribution]').includesText('Top namespace');

    // TODO CMB - add assertion so double charts show for single historical month
    // TODO and check for empty state there
    // assert
    // .dom('[data-test-chart-container="new-clients"] [data-test-empty-state-subtext]')
    // .includesText(
    //   'There are no new clients for this namespace during this time period.',
    //   'Shows empty state if no new client counts'
    // );
  });

  // flaky test -- does not consistently run the same number of assertions
  // refactor before using assert.expect
  // eslint-disable-next-line qunit/no-commented-tests
  // test('filters correctly on history with full data', async function (assert) {
  //   /* eslint qunit/require-expect: "warn" */
  //   // assert.expect(44);
  //   const licenseStart = startOfMonth(subMonths(new Date(), 6));
  //   const licenseEnd = addMonths(new Date(), 6);
  //   const lastMonth = addMonths(new Date(), -1);
  //   const config = generateConfigResponse();
  //   const activity = generateActivityResponse(5, licenseStart, lastMonth);
  //   const license = generateLicenseResponse(licenseStart, licenseEnd);
  //   this.server = new Pretender(function () {
  //     this.get('/v1/sys/license/status', () => sendResponse(license));
  //     this.get('/v1/sys/internal/counters/activity', () => sendResponse(activity));
  //     this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
  //     this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
  //     this.get('/v1/sys/health', this.passthrough);
  //     this.get('/v1/sys/seal-status', this.passthrough);
  //     this.post('/v1/sys/capabilities-self', this.passthrough);
  //     this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
  //   });
  //   await visit('/vault/clients/history');
  //   assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
  //   assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
  //   assert.dom(SELECTORS.usageStats).exists('usage stats block exists');
  //   assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
  //   const { total } = activity.data;

  //   // FILTER BY NAMESPACE
  //   await clickTrigger();
  //   await searchSelect.options.objectAt(0).click();
  //   await settled();
  //   assert.ok(true, 'Filter by first namespace');
  //   assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('15');
  //   assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('5');
  //   assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('10');
  //   assert.dom('[data-test-top-attribution]').includesText('Top auth method');

  //   // FILTER BY AUTH METHOD
  //   await clickTrigger();
  //   await searchSelect.options.objectAt(0).click();
  //   await settled();
  //   assert.ok(true, 'Filter by first auth method');
  //   assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('5');
  //   assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('3');
  //   assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('2');
  //   assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution block');

  //   await click('#namespace-search-select [data-test-selected-list-button="delete"]');
  //   assert.ok(true, 'Remove namespace filter without first removing auth method filter');
  //   assert.dom('[data-test-top-attribution]').includesText('Top namespace');
  //   assert
  //     .dom('[data-test-stat-text="total-clients"] .stat-value')
  //     .hasText(total.clients.toString(), 'total clients stat is back to unfiltered value');
  // });

  test('shows warning if upgrade happened within license period', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const lastMonth = addMonths(new Date(), -1);
    const config = generateConfigResponse();
    const activity = generateActivityResponse(5, licenseStart, lastMonth);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(activity));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () =>
        sendResponse({
          data: {
            keys: ['1.9.0'],
            key_info: {
              '1.9.0': {
                previous_version: null,
                timestamp_installed: formatRFC3339(addMonths(new Date(), -2)),
              },
            },
          },
        })
      );
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    assert.dom('[data-test-alert-banner="alert"]').includesText('Vault was upgraded');
  });

  test('Shows empty if license start date is current month', async function (assert) {
    const licenseStart = new Date();
    const licenseEnd = addMonths(new Date(), 12);
    const config = generateConfigResponse();
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => this.passthrough);
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () =>
        sendResponse({
          keys: [],
        })
      );
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No data for this billing period');
    assert
      .dom(SELECTORS.dateDisplay)
      .hasText(format(licenseStart, 'MMMM yyyy'), 'Shows license date, gives ability to edit');
    assert.dom(SELECTORS.monthDropdown).exists('Dropdown exists to select month');
    assert.dom(SELECTORS.yearDropdown).exists('Dropdown exists to select year');
  });

  test('shows correct interface if no permissions on license', async function (assert) {
    const config = generateConfigResponse();
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(null, 403));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    // Message changes depending on ent or OSS
    assert.dom(SELECTORS.emptyStateTitle).exists('Empty state exists');
    assert.dom(SELECTORS.monthDropdown).exists('Dropdown exists to select month');
    assert.dom(SELECTORS.yearDropdown).exists('Dropdown exists to select year');
  });

  test('shows error template if permissions denied querying activity response with no data', async function (assert) {
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(null, 403));
      this.get('/v1/sys/version-history', () => sendResponse(null, 403));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(null, 403));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(null, 403));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert
      .dom(SELECTORS.emptyStateTitle)
      .includesText('start date found', 'Empty state shows no billing start date');
    await click(SELECTORS.monthDropdown);
    await click(this.element.querySelector('[data-test-month-list] button:not([disabled])'));
    await click(SELECTORS.yearDropdown);
    await click(this.element.querySelector('[data-test-year-list] button:not([disabled])'));
    await click(SELECTORS.dateDropdownSubmit);
    assert
      .dom(SELECTORS.emptyStateTitle)
      .hasText('You are not authorized', 'Empty state displays not authorized message');
  });
});
