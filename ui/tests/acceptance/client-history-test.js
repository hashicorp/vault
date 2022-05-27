import { module, test } from 'qunit';
// import { visit, currentURL, click, findAll, find } from '@ember/test-helpers';
import { visit, currentURL, click, findAll } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import { addMonths, format, formatRFC3339, startOfMonth, subMonths } from 'date-fns';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { SELECTORS, sendResponse, overrideResponse } from '../helpers/clients';
// import endOfMonth from 'date-fns/endOfMonth';
// import { create } from 'ember-cli-page-object';
// import { clickTrigger } from 'ember-power-select/test-support/helpers';
// import ss from 'vault/tests/pages/components/search-select';
// const searchSelect = create(ss);

const NEW_DATE = new Date();
const LICENSE_START = startOfMonth(subMonths(NEW_DATE, 6));
// const LICENSE_END = endOfMonth(addMonths(NEW_DATE, 6));
const lastMonth = startOfMonth(subMonths(NEW_DATE, 1));

// upgrade happened 1 month after license start
// const UPGRADE_DATE = addMonths(LICENSE_START, 1);

module('Acceptance | clients history tab', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'clients';
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  hooks.beforeEach(function () {
    return authPage.login();
  });

  // hooks.afterEach(function () {
  //   this.server.shutdown();
  // });

  test('shows warning when config off, no data, queries_available=true', async function (assert) {
    assert.expect(6);
    this.server.get('sys/internal/counters/activity', () => sendResponse(null, 204));
    this.server.get('sys/internal/counters/config', () => {
      return {
        request_id: 'some-config-id',
        data: {
          default_report_months: 12,
          enabled: 'default-disable',
          queries_available: true,
          retention_months: 24,
        },
      };
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    assert.dom('[data-test-tracking-disabled] .message-title').hasText('Tracking is disabled');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No data received');
    assert.dom(SELECTORS.filterBar).doesNotExist('Shows filter bar to search previous dates');
    assert.dom(SELECTORS.usageStats).doesNotExist('No usage stats');
  });

  test('shows warning when config off, no data, queries_available=false', async function (assert) {
    assert.expect(4);
    this.server.get('sys/internal/counters/activity', () => sendResponse(null, 204));
    this.server.get('sys/internal/counters/config', () => {
      return {
        request_id: 'some-config-id',
        data: {
          default_report_months: 12,
          enabled: 'default-disable',
          queries_available: false,
          retention_months: 24,
        },
      };
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('Data tracking is disabled');
    assert.dom(SELECTORS.filterBar).doesNotExist('Filter bar is hidden when no data available');
  });

  test('shows empty state when config enabled and queries_available=false', async function (assert) {
    assert.expect(4);
    this.server.get('sys/internal/counters/activity', () => sendResponse(null, 204));
    this.server.get('sys/internal/counters/config', () => {
      return {
        request_id: 'some-config-id',
        data: {
          default_report_months: 12,
          enabled: 'default-enable',
          queries_available: false,
          retention_months: 24,
        },
      };
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');

    assert.dom(SELECTORS.emptyStateTitle).hasText('No monthly history');
    assert.dom(SELECTORS.filterBar).doesNotExist('Does not show filter bar');
  });

  test('visiting history tab config on and data with mounts', async function (assert) {
    assert.expect(7);
    // TODO CMB: wire up dynamic generateActivity to mirage handler
    // const activity = generateActivityResponse(5, LICENSE_START, lastMonth);
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');

    assert
      .dom(SELECTORS.dateDisplay)
      .hasText(format(LICENSE_START, 'MMMM yyyy'), 'billing start month is correctly parsed from license');
    assert
      .dom(SELECTORS.rangeDropdown)
      .hasText(
        `${format(LICENSE_START, 'MMMM yyyy')} - ${format(lastMonth, 'MMMM yyyy')}`,
        'Date range shows dates correctly parsed activity response'
      );
    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert.dom(SELECTORS.monthlyUsageBlock).exists('Shows monthly usage block');
    assert
      .dom(SELECTORS.runningTotalMonthlyCharts)
      .exists('Shows running totals with monthly breakdown charts');
    // TODO CMB update when generate monthly data dynamically
    assert.equal(findAll('[data-test-line-chart="plot-point"]').length, 5, `5 plot points show`);
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
  //   const config = overrideConfigResponse();
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
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    assert.dom('[data-test-alert-banner="alert"]').includesText('Vault was upgraded');
  });

  test('Shows empty if license start date is current month', async function (assert) {
    const licenseStart = NEW_DATE;
    const licenseEnd = addMonths(NEW_DATE, 12);
    this.server.get('sys/license/status', function () {
      return {
        request_id: 'my-license-request-id',
        data: {
          autoloaded: {
            license_id: 'my-license-id',
            start_time: formatRFC3339(licenseStart),
            expiration_time: formatRFC3339(licenseEnd),
          },
        },
      };
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
    this.server.get('/sys/license/status', () => overrideResponse(403));
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    // Message changes depending on ent or OSS
    assert.dom(SELECTORS.emptyStateTitle).exists('Empty state exists');
    assert.dom(SELECTORS.monthDropdown).exists('Dropdown exists to select month');
    assert.dom(SELECTORS.yearDropdown).exists('Dropdown exists to select year');
  });

  test('shows error template if permissions denied querying activity response with no data', async function (assert) {
    this.server.get('sys/license/status', () => overrideResponse(403));
    this.server.get('sys/version-history', () => overrideResponse(403));
    this.server.get('sys/internal/counters/config', () => overrideResponse(403));
    this.server.get('sys/internal/counters/activity', () => overrideResponse(403));

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
