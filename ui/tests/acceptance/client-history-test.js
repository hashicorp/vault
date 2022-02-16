import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import authPage from 'vault/tests/pages/auth';
import { addMonths, format, startOfMonth, subMonths } from 'date-fns';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import {
  generateActivityResponse,
  generateConfigResponse,
  generateLicenseResponse,
  SELECTORS,
  sendResponse,
} from '../helpers/clients';

const searchSelect = create(ss);

module('Acceptance | clients history', function (hooks) {
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
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(activity));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.activeTab).hasText('History', 'history tab is active');

    assert.dom('[data-test-tracking-disabled] .message-title').hasText('Tracking is disabled');
    // TODO: still allows query by previous dates
    // TODO: 0's on stat text
    await this.pauseTest();
  });

  test('shows warning when config off, no data, queries unavailable', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse({ enabled: 'default-disable', queries_available: false });
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/license/status', () => sendResponse(license));
      this.get('/v1/sys/internal/counters/activity', () => sendResponse(activity));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
    });
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.activeTab).hasText('History', 'history tab is active');
    await this.pauseTest();
    assert.dom(SELECTORS.emptyStateTitle).hasText('Data tracking is disabled');
    // TODO: filter bar hidden
    assert.dom(SELECTORS.filterBar).doesNotExist('Filter bar is not hidden when no data available');
    // Hide billing start month?
    await this.pauseTest();
  });

  test('shows empty state and warning when no queries', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    const config = generateConfigResponse({ queries_available: false });
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
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
    // History Tab
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom(SELECTORS.activeTab).hasText('History', 'history tab is active');

    assert.dom(SELECTORS.emptyStateTitle).hasText('No monthly history');
    await this.pauseTest();
  });
  test('visiting history tab with no data and config on', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const config = generateConfigResponse();
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
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
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Attribution block is not shown when no data');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    // TODO: Filters correct
    // TODO: don't show namespace filter if none exist
  });
  test('filters correctly on history with full data', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const config = generateConfigResponse();
    const activity = generateActivityResponse(5, licenseStart, licenseEnd);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
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
