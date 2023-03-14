import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { endOfMonth, formatRFC3339 } from 'date-fns';
import { click } from '@ember/test-helpers';
import subMonths from 'date-fns/subMonths';
import timestamp from 'core/utils/timestamp';

module('Integration | Component | clients/attribution', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const mockNow = timestamp.now();
    this.mockNow = mockNow;
    this.set('startTimestamp', formatRFC3339(subMonths(mockNow, 6)));
    this.set('timestamp', formatRFC3339(mockNow));
    this.set('selectedNamespace', null);
    this.set('chartLegend', [
      { label: 'entity clients', key: 'entity_clients' },
      { label: 'non-entity clients', key: 'non_entity_clients' },
    ]);
    this.set('totalUsageCounts', { clients: 15, entity_clients: 10, non_entity_clients: 5 });
    this.set('totalClientAttribution', [
      { label: 'second', clients: 10, entity_clients: 7, non_entity_clients: 3 },
      { label: 'first', clients: 5, entity_clients: 3, non_entity_clients: 2 },
    ]);
    this.set('totalMountsData', { clients: 5, entity_clients: 3, non_entity_clients: 2 });
    this.set('namespaceMountsData', [
      { label: 'auth1/', clients: 3, entity_clients: 2, non_entity_clients: 1 },
      { label: 'auth2/', clients: 2, entity_clients: 1, non_entity_clients: 1 },
    ]);
  });

  test('it renders empty state with no data', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution @chartLegend={{this.chartLegend}} />
    `);

    assert.dom('[data-test-component="empty-state"]').exists();
    assert.dom('[data-test-empty-state-title]').hasText('No data found');
    assert.dom('[data-test-attribution-description]').hasText('There is a problem gathering data');
    assert.dom('[data-test-attribution-export-button]').doesNotExist();
    assert.dom('[data-test-attribution-timestamp]').doesNotHaveTextContaining('Updated');
  });

  test('it renders with data for namespaces', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution
        @chartLegend={{this.chartLegend}}
        @totalClientAttribution={{this.totalClientAttribution}}
        @totalUsageCounts={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.timestamp}}
        @selectedNamespace={{this.selectedNamespace}}
        @isHistoricalMonth={{false}}
        />
    `);

    assert.dom('[data-test-component="empty-state"]').doesNotExist();
    assert.dom('[data-test-horizontal-bar-chart]').exists('chart displays');
    assert.dom('[data-test-attribution-export-button]').exists();
    assert
      .dom('[data-test-attribution-description]')
      .hasText(
        'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.'
      );
    assert
      .dom('[data-test-attribution-subtext]')
      .hasText(
        'The total clients in the namespace for this date range. This number is useful for identifying overall usage volume.'
      );
    assert.dom('[data-test-top-attribution]').includesText('namespace').includesText('second');
    assert.dom('[data-test-attribution-clients]').includesText('namespace').includesText('10');
  });

  test('it renders two charts and correct text for single, historical month', async function (assert) {
    this.start = formatRFC3339(subMonths(this.mockNow, 1));
    this.end = formatRFC3339(subMonths(endOfMonth(this.mockNow), 1));
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution
        @chartLegend={{this.chartLegend}}
        @totalClientAttribution={{this.totalClientAttribution}}
        @totalUsageCounts={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
        @startTimestamp={{this.start}}
        @endTimestamp={{this.end}}
        @selectedNamespace={{this.selectedNamespace}}
        @isHistoricalMonth={{true}}
        />
    `);
    assert
      .dom('[data-test-attribution-description]')
      .includesText(
        'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.',
        'renders correct auth attribution description'
      );
    assert
      .dom('[data-test-chart-container="total-clients"] .chart-description')
      .includesText(
        'The total clients in the namespace for this month. This number is useful for identifying overall usage volume.',
        'renders total monthly namespace text'
      );
    assert
      .dom('[data-test-chart-container="new-clients"] .chart-description')
      .includesText(
        'The new clients in the namespace for this month. This aids in understanding which namespaces create and use new clients.',
        'renders new monthly namespace text'
      );
    this.set('selectedNamespace', 'second');

    assert
      .dom('[data-test-attribution-description]')
      .includesText(
        'This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where clients are originating. Authentication methods are organized by path.',
        'renders correct auth attribution description'
      );
    assert
      .dom('[data-test-chart-container="total-clients"] .chart-description')
      .includesText(
        'The total clients used by the auth method for this month. This number is useful for identifying overall usage volume.',
        'renders total monthly auth method text'
      );
    assert
      .dom('[data-test-chart-container="new-clients"] .chart-description')
      .includesText(
        'The new clients used by the auth method for this month. This aids in understanding which auth methods create and use new clients.',
        'renders new monthly auth method text'
      );
  });

  test('it renders single chart for current month', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution
        @chartLegend={{this.chartLegend}}
        @totalClientAttribution={{this.totalClientAttribution}}
        @totalUsageCounts={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
        @startTimestamp={{this.timestamp}}
        @endTimestamp={{this.timestamp}}
        @selectedNamespace={{this.selectedNamespace}}
        @isHistoricalMonth={{false}}
        />
    `);
    assert
      .dom('[data-test-chart-container="single-chart"]')
      .exists('renders single chart with total clients');
    assert
      .dom('[data-test-attribution-subtext]')
      .hasTextContaining('this month', 'renders total monthly namespace text');
  });

  test('it renders single chart and correct text for for date range', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution
        @chartLegend={{this.chartLegend}}
        @totalClientAttribution={{this.totalClientAttribution}}
        @totalUsageCounts={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.timestamp}}
        @selectedNamespace={{this.selectedNamespace}}
        @isHistoricalMonth={{false}}
        />
    `);

    assert
      .dom('[data-test-chart-container="single-chart"]')
      .exists('renders single chart with total clients');
    assert
      .dom('[data-test-attribution-subtext]')
      .hasTextContaining('date range', 'renders total monthly namespace text');
  });

  test('it renders with data for selected namespace auth methods for a date range', async function (assert) {
    this.set('selectedNamespace', 'second');
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution
        @chartLegend={{this.chartLegend}}
        @totalClientAttribution={{this.namespaceMountsData}}
        @totalUsageCounts={{this.totalUsageCounts}}
        @responseTimestamp={{this.timestamp}}
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.timestamp}}
        @selectedNamespace={{this.selectedNamespace}}
        @isHistoricalMonth={{this.isHistoricalMonth}}
        />
    `);

    assert.dom('[data-test-component="empty-state"]').doesNotExist();
    assert.dom('[data-test-horizontal-bar-chart]').exists('chart displays');
    assert.dom('[data-test-attribution-export-button]').exists();
    assert
      .dom('[data-test-attribution-description]')
      .hasText(
        'This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where clients are originating. Authentication methods are organized by path.'
      );
    assert
      .dom('[data-test-attribution-subtext]')
      .hasText(
        'The total clients used by the auth method for this date range. This number is useful for identifying overall usage volume.'
      );
    assert.dom('[data-test-top-attribution]').includesText('auth method').includesText('auth1/');
    assert.dom('[data-test-attribution-clients]').includesText('auth method').includesText('3');
  });

  test('it renders modal', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution
        @chartLegend={{this.chartLegend}}
        @totalClientAttribution={{this.namespaceMountsData}}
        @responseTimestamp={{this.timestamp}}
        @startTimestamp="2022-06-01T23:00:11.050Z"
        @endTimestamp="2022-12-01T23:00:11.050Z"
        />
    `);
    await click('[data-test-attribution-export-button]');
    assert.dom('.modal.is-active .title').hasText('Export attribution data', 'modal appears to export csv');
    assert.dom('.modal.is-active').includesText('June 2022 - December 2022');
  });
});
