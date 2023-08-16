import { module, test } from 'qunit';
import { visit, currentURL, click, findAll, find, settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import { addMonths, format, formatRFC3339, startOfMonth, subMonths } from 'date-fns';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { SELECTORS, overrideResponse } from '../helpers/clients';
import { create } from 'ember-cli-page-object';
import ss from 'vault/tests/pages/components/search-select';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';

const searchSelect = create(ss);

const NEW_DATE = new Date();
const LAST_MONTH = startOfMonth(subMonths(NEW_DATE, 1));
const COUNTS_START = subMonths(NEW_DATE, 12); // pretend vault user started cluster 1 year ago

// for testing, we're in the middle of a license/billing period
const LICENSE_START = startOfMonth(subMonths(NEW_DATE, 6));

// upgrade happened 1 month after license start
const UPGRADE_DATE = addMonths(LICENSE_START, 1);

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

  test('shows warning when config off, no data, queries_available=true', async function (assert) {
    assert.expect(6);
    this.server.get('sys/internal/counters/activity', () => overrideResponse(204));
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
    this.server.get('sys/internal/counters/activity', () => overrideResponse(204));
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
    this.server.get('sys/internal/counters/activity', () => overrideResponse(204));
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
    assert.expect(8);
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');

    assert
      .dom(SELECTORS.dateDisplay)
      .hasText(format(LICENSE_START, 'MMMM yyyy'), 'billing start month is correctly parsed from license');
    assert
      .dom(SELECTORS.rangeDropdown)
      .hasText(
        `${format(LICENSE_START, 'MMMM yyyy')} - ${format(LAST_MONTH, 'MMMM yyyy')}`,
        'Date range shows dates correctly parsed activity response'
      );
    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert.dom(SELECTORS.monthlyUsageBlock).exists('Shows monthly usage block');
    assert
      .dom(SELECTORS.runningTotalMonthlyCharts)
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(find('[data-test-line-chart="x-axis-labels"] g.tick text'))
      .hasText(`${format(LICENSE_START, 'M/yy')}`, 'x-axis labels start with billing start date');
    assert.equal(
      findAll('[data-test-line-chart="plot-point"]').length,
      5,
      `line chart plots 5 points to match query`
    );
  });

  test('updates correctly when querying date ranges', async function (assert) {
    assert.expect(26);
    // TODO CMB: wire up dynamically generated activity to mirage clients handler
    // const activity = generateActivityResponse(5, LICENSE_START, LAST_MONTH);
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');

    // query for single, historical month with no new counts
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-show-calendar]');
    if (parseInt(find('[data-test-display-year]').innerText) > LICENSE_START.getFullYear()) {
      await click('[data-test-previous-year]');
    }
    await click(find(`[data-test-calendar-month=${ARRAY_OF_MONTHS[LICENSE_START.getMonth()]}]`));

    assert.dom('[data-test-usage-stats]').exists('total usage stats show');
    assert
      .dom(SELECTORS.runningTotalMonthStats)
      .doesNotExist('running total single month stat boxes do not show');
    assert
      .dom(SELECTORS.runningTotalMonthlyCharts)
      .doesNotExist('running total month over month charts do not show');
    assert.dom(SELECTORS.monthlyUsageBlock).doesNotExist('does not show monthly usage block');
    assert.dom(SELECTORS.attributionBlock).exists('attribution area shows');
    assert
      .dom('[data-test-chart-container="new-clients"] [data-test-component="empty-state"]')
      .exists('new client attribution has empty state');
    assert
      .dom('[data-test-empty-state-subtext]')
      .hasText('There are no new clients for this namespace during this time period.    ');
    assert.dom('[data-test-chart-container="total-clients"]').exists('total client attribution chart shows');

    // reset to billing period
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-current-billing-period]');

    // change billing start to month/year of first upgrade
    await click('[data-test-start-date-editor] button');
    await click(SELECTORS.monthDropdown);
    await click(find(`[data-test-date-modal-month="${ARRAY_OF_MONTHS[UPGRADE_DATE.getMonth()]}"]`));
    await click(SELECTORS.yearDropdown);
    await click(find(`[data-test-date-modal-year="${UPGRADE_DATE.getFullYear()}`));
    await click('[data-test-modal-save]');

    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert.dom(SELECTORS.monthlyUsageBlock).exists('Shows monthly usage block');
    assert
      .dom(SELECTORS.runningTotalMonthlyCharts)
      .exists('Shows running totals with monthly breakdown charts');
    assert
      .dom(find('[data-test-line-chart="x-axis-labels"] g.tick text'))
      .hasText(`${format(UPGRADE_DATE, 'M/yy')}`, 'x-axis labels start with updated billing start month');
    assert.equal(
      findAll('[data-test-line-chart="plot-point"]').length,
      5,
      `line chart plots 5 points to match query`
    );

    // query custom end month
    const customEndDate = subMonths(NEW_DATE, 3);
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-show-calendar]');
    if (parseInt(find('[data-test-display-year]').innerText) !== customEndDate.getFullYear()) {
      await click('[data-test-previous-year]');
    }
    await click(find(`[data-test-calendar-month=${ARRAY_OF_MONTHS[customEndDate.getMonth()]}]`));

    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert.dom(SELECTORS.monthlyUsageBlock).exists('Shows monthly usage block');
    assert
      .dom(SELECTORS.runningTotalMonthlyCharts)
      .exists('Shows running totals with monthly breakdown charts');
    assert.equal(
      findAll('[data-test-line-chart="plot-point"]').length,
      3,
      `line chart plots 3 points to match query`
    );
    let xAxisLabels = findAll('[data-test-line-chart="x-axis-labels"] g.tick text');
    assert
      .dom(xAxisLabels[xAxisLabels.length - 1])
      .hasText(`${format(subMonths(LAST_MONTH, 2), 'M/yy')}`, 'x-axis labels end with queried end month');

    // query for single, historical month
    await click(SELECTORS.rangeDropdown);
    await click('[data-test-show-calendar]');
    if (parseInt(find('[data-test-display-year]').innerText) !== UPGRADE_DATE.getFullYear()) {
      await click('[data-test-previous-year]');
    }
    await click(find(`[data-test-calendar-month=${ARRAY_OF_MONTHS[UPGRADE_DATE.getMonth()]}]`));

    assert.dom(SELECTORS.runningTotalMonthStats).exists('running total single month stat boxes show');
    assert
      .dom(SELECTORS.runningTotalMonthlyCharts)
      .doesNotExist('running total month over month charts do not show');
    assert.dom(SELECTORS.monthlyUsageBlock).doesNotExist('Does not show monthly usage block');
    assert.dom(SELECTORS.attributionBlock).exists('attribution area shows');
    assert.dom('[data-test-chart-container="new-clients"]').exists('new client attribution chart shows');
    assert.dom('[data-test-chart-container="total-clients"]').exists('total client attribution chart shows');

    // reset to billing period
    await click('[data-test-popup-menu-trigger]');
    await click('[data-test-current-billing-period]');

    // query month older than count start date
    await click('[data-test-start-date-editor] button');
    await click(SELECTORS.yearDropdown);
    await click(find(`[data-test-date-modal-year="${LICENSE_START.getFullYear() - 3}`));
    await click('[data-test-modal-save]');

    assert
      .dom('[data-test-alert-banner="alert"]')
      .hasTextContaining(
        `We only have data from ${format(COUNTS_START, 'MMMM yyyy')}`,
        'warning banner displays that date queried was prior to count start date'
      );
  });

  test('filters correctly on history with full data', async function (assert) {
    assert.expect(19);
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    assert
      .dom(SELECTORS.runningTotalMonthlyCharts)
      .exists('Shows running totals with monthly breakdown charts');
    assert.dom(SELECTORS.attributionBlock).exists('Shows attribution area');
    assert.dom(SELECTORS.monthlyUsageBlock).exists('Shows monthly usage block');

    // FILTER BY NAMESPACE
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();

    await settled();
    assert.ok(true, 'Filter by first namespace');
    assert.dom('[data-test-top-attribution]').includesText('Top auth method');
    assert.dom('[data-test-running-total-entity]').includesText('23,326', 'total entity clients is accurate');
    assert
      .dom('[data-test-running-total-nonentity]')
      .includesText('17,826', 'total non-entity clients is accurate');
    assert.dom('[data-test-attribution-clients]').includesText('10,142', 'top attribution clients accurate');

    // FILTER BY AUTH METHOD
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    assert.ok(true, 'Filter by first auth method');
    assert.dom('[data-test-running-total-entity]').includesText('6,508', 'total entity clients is accurate');
    assert
      .dom('[data-test-running-total-nonentity]')
      .includesText('3,634', 'total non-entity clients is accurate');
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution block');

    await click('#namespace-search-select [data-test-selected-list-button="delete"]');
    assert.ok(true, 'Remove namespace filter without first removing auth method filter');
    assert.dom('[data-test-top-attribution]').includesText('Top namespace');
    assert
      .dom('[data-test-attribution-clients]')
      .hasTextContaining('41,152', 'top attribution clients back to unfiltered value');
    assert
      .dom('[data-test-running-total-entity]')
      .hasTextContaining('98,289', 'total entity clients is back to unfiltered value');
    assert
      .dom('[data-test-running-total-nonentity]')
      .hasTextContaining('96,412', 'total non-entity clients is back to unfiltered value');
  });

  test('shows warning if upgrade happened within license period', async function (assert) {
    assert.expect(3);
    this.server.get('sys/version-history', function () {
      return {
        data: {
          keys: ['1.9.0', '1.9.1', '1.9.2', '1.10.1'],
          key_info: {
            '1.9.0': {
              previous_version: null,
              timestamp_installed: formatRFC3339(subMonths(UPGRADE_DATE, 4)),
            },
            '1.9.1': {
              previous_version: '1.9.0',
              timestamp_installed: formatRFC3339(subMonths(UPGRADE_DATE, 3)),
            },
            '1.9.2': {
              previous_version: '1.9.1',
              timestamp_installed: formatRFC3339(subMonths(UPGRADE_DATE, 2)),
            },
            '1.10.1': {
              previous_version: '1.9.2',
              timestamp_installed: formatRFC3339(UPGRADE_DATE),
            },
          },
        },
      };
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history', 'clients/history URL is correct');
    assert.dom(SELECTORS.historyActiveTab).hasText('History', 'history tab is active');
    assert
      .dom('[data-test-alert-banner="alert"]')
      .hasTextContaining(
        `Warning Vault was upgraded to 1.10.1 on ${format(
          UPGRADE_DATE,
          'MMM d, yyyy'
        )}. We added monthly breakdowns and mount level attribution starting in 1.10, so keep that in mind when looking at the data below. Learn more here.`
      );
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
