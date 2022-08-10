import { module, test } from 'qunit';
import { visit, currentURL, settled, click } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import { overrideResponse, SELECTORS } from '../helpers/clients';

const searchSelect = create(ss);

module('Acceptance | clients current tab', function (hooks) {
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

  test('shows empty state when config disabled, no data', async function (assert) {
    assert.expect(3);
    this.server.get('sys/internal/counters/config', () => {
      return {
        request_id: 'some-config-id',
        data: {
          default_report_months: 12,
          enabled: 'default-disable',
          retention_months: 24,
        },
      };
    });
    this.server.get('sys/internal/counters/activity/monthly', () => overrideResponse(204));
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('Tracking is disabled');
  });

  test('shows empty state when config enabled, no data', async function (assert) {
    assert.expect(3);
    this.server.get('sys/internal/counters/activity/monthly', () => {
      return {
        request_id: 'some-monthly-id',
        data: {
          by_namespace: [],
          clients: 0,
          distinct_entities: 0,
          entity_clients: 0,
          months: [],
          non_entity_clients: 0,
          non_entity_tokens: 0,
        },
      };
    });
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.emptyStateTitle).hasText('No data received');
  });

  test('filters correctly on current with full data', async function (assert) {
    assert.expect(27);
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.usageStats).exists('usage stats block exists');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');

    // TODO update with dynamic counts
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('175');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('132');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('43');
    assert.dom('[data-test-clients-attribution]').exists('Shows attribution area');
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
    await settled();

    // FILTER BY NAMESPACE
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('100');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('85');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('15');
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
    await settled();

    // FILTER BY AUTH METHOD
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();

    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('35');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('20');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('15');
    assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution block');

    // Delete auth filter goes back to filtered only on namespace
    await click('#auth-method-search-select [data-test-selected-list-button="delete"]');

    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('100');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('85');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('15');
    await settled();
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
    assert.dom(SELECTORS.attributionBlock).exists('Still shows attribution block');
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();
    // Delete namespace filter with auth filter on
    await click('#namespace-search-select-monthly [data-test-selected-list-button="delete"]');
    // Goes back to no filters
    // TODO update with dynamic counts
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('175');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('132');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('43');
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
    assert.dom('[data-test-chart-container="new-clients"] [data-test-empty-state-subtext]').doesNotExist();
  });

  test('filters correctly on current with no auth mounts', async function (assert) {
    assert.expect(16);
    await visit('/vault/clients/current');
    assert.equal(currentURL(), '/vault/clients/current');
    assert.dom(SELECTORS.currentMonthActiveTab).hasText('Current month', 'current month tab is active');
    assert.dom(SELECTORS.usageStats).exists('usage stats block exists');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    // TODO CMB update with dynamic data
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('175');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('132');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('43');
    assert.dom('[data-test-clients-attribution]').exists('Shows attribution area');
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();

    // Filter by namespace
    await clickTrigger();
    await searchSelect.options.objectAt(0).click();
    await settled();

    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('100');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('85');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('15');

    // TODO add month data without mounts
    // assert.dom(SELECTORS.attributionBlock).doesNotExist('Does not show attribution');
    // assert.dom('#auth-method-search-select').doesNotExist('Auth method filter is not shown');

    // Remove namespace filter
    await click('#namespace-search-select-monthly [data-test-selected-list-button="delete"]');

    // Goes back to no filters
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('175');
    assert.dom('[data-test-stat-text="entity-clients"] .stat-value').hasText('132');
    assert.dom('[data-test-stat-text="non-entity-clients"] .stat-value').hasText('43');
    assert.dom('[data-test-chart-container="new-clients"]').doesNotExist();
  });

  test('shows correct empty state when config off but no read on config', async function (assert) {
    assert.expect(2);
    this.server.get('sys/internal/counters/activity/monthly', () => overrideResponse(204));
    this.server.get('sys/internal/counters/config', () => overrideResponse(403));
    await visit('/vault/clients/current');
    assert.dom(SELECTORS.filterBar).doesNotExist('Filter bar is not shown');
    assert.dom(SELECTORS.emptyStateTitle).containsText('No data available', 'Shows no data empty state');
  });
});
