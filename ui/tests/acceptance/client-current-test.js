import { module, test } from 'qunit';
import { visit, currentURL, settled, click } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import authPage from 'vault/tests/pages/auth';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import {
  generateConfigResponse,
  generateCurrentMonthResponse,
  SELECTORS,
  sendResponse,
} from '../helpers/clients';

const searchSelect = create(ss);

module('Acceptance | clients current', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('shows empty state when config disabled, no data', async function (assert) {
    const config = generateConfigResponse({ enabled: 'default-disable' });
    const monthly = generateCurrentMonthResponse({ configEnabled: false });
    this.server = new Pretender(function () {
      this.get('/v1/sys/internal/counters/activity/monthly', () => sendResponse(monthly));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({}));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
    });
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('Tracking is disabled');
  });

  test('shows empty state when config enabled, no data', async function (assert) {
    const config = generateConfigResponse();
    const monthly = generateCurrentMonthResponse();
    this.server = new Pretender(function () {
      this.get('/v1/sys/internal/counters/activity/monthly', () => sendResponse(monthly));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No data received');
  });
  // flaky test -- assertion count is not consistent
  // eslint-disable-next-line
  // test('filters correctly on current with full data', async function (assert) {
  //   // uncomment once assertion count is consistent
  //   // assert.expect(65);
  //   const config = generateConfigResponse();
  //   const monthly = generateCurrentMonthResponse(3);
  //   this.server = new Pretender(function () {
  //     this.get('/v1/sys/internal/counters/activity/monthly', () => sendResponse(monthly));
  //     this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
  //     this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
  //     this.get('/v1/sys/health', this.passthrough);
  //     this.get('/v1/sys/seal-status', this.passthrough);
  //     this.post('/v1/sys/capabilities-self', this.passthrough);
  //     this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
  //   });
  //   await visit('/vault/clients/current');
  //   assert.equal(currentURL(), '/vault/clients/current');
  //   assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
  //   assert.dom(SELECTORS.usageStats).exists('usage stats block exists');
  //   assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
  //   const { clients, entity_clients, non_entity_clients } = monthly.data;
  //   assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText(clients.toString());
  //   assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText(entity_clients.toString());
  //   assert
  //     .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
  //     .hasText(non_entity_clients.toString());
  //   assert.dom('[data-test-clients-attribution]').exists('Shows attribution area');
  //   assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
  //   await settled();

  //   // FILTER BY NAMESPACE
  //   await clickTrigger();
  //   await searchSelect.options.objectAt(0).click();
  //   await settled();
  //   assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('15');
  //   assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('5');
  //   assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('10');
  //   assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
  //   await settled();
  //   // FILTER BY AUTH METHOD
  //   await clickTrigger();
  //   await searchSelect.options.objectAt(0).click();
  //   await waitUntil(() => find('#auth-method-search-select'));
  //   assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('5');
  //   assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('3');
  //   assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('2');
  //   assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution block');
  //   // Delete auth filter goes back to filtered only on namespace
  //   await click('#auth-method-search-select [data-test-selected-list-button="delete"]');
  //   assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('15');
  //   assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('5');
  //   assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('10');
  //   await settled();
  //   assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
  //   assert.dom(SELECTORS.attributionBlock).exists('Still shows attribution block');
  //   await clickTrigger();
  //   await searchSelect.options.objectAt(0).click();
  //   await settled();
  //   // Delete namespace filter with auth filter on
  //   await click('#namespace-search-select-monthly [data-test-selected-list-button="delete"]');
  //   // Goes back to no filters
  //   assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText(clients.toString());
  //   assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText(entity_clients.toString());
  //   assert
  //     .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
  //     .hasText(non_entity_clients.toString());
  //   assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
  //   assert.dom('[data-test-chart-container="new-clients"] [data-test-empty-state-subtext]').doesNotExist();
  // });

  test('filters correctly on current with no auth mounts', async function (assert) {
    const config = generateConfigResponse();
    const monthly = generateCurrentMonthResponse(3, true /* skip mounts */);
    this.server = new Pretender(function () {
      this.get('/v1/sys/internal/counters/activity/monthly', () => sendResponse(monthly));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(config));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.usageStats).exists('usage stats block exists');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    const { clients, entity_clients, non_entity_clients } = monthly.data;
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText(clients.toString());
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText(entity_clients.toString());
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText(non_entity_clients.toString());
    assert.dom('[data-test-clients-attribution]').exists('Shows attribution area');
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();

    // Filter by namespace
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('15');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('5');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('10');
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution');
    assert.dom('#auth-method-search-select').doesNotExist('Auth method filter is not shown');
    // Remove namespace filter
    await click('#namespace-search-select-monthly [data-test-selected-list-button="delete"]');
    // Goes back to no filters
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText(clients.toString());
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText(entity_clients.toString());
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText(non_entity_clients.toString());
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
  });

  test('shows correct empty state when config off but no read on config', async function (assert) {
    this.server = new Pretender(function () {
      // Monthly responds with 204 when config off
      this.get('/v1/sys/internal/counters/activity/monthly', () => sendResponse(null, 204));
      this.get('/v1/sys/internal/counters/config', () => sendResponse(null, 403));
      this.get('/v1/sys/version-history', () => sendResponse({ keys: [] }));
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/internal/ui/feature-flags', this.passthrough);
    });
    await visit('/vault/clients/current');
    assert.dom(SELECTORS.filterBar).doesNotExist('Filter bar is not shown');
    assert.dom(SELECTORS.emptyStateTitle).containsText('No data available', 'Shows no data empty state');
  });
});
