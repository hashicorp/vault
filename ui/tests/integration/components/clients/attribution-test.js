/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { endOfMonth, formatRFC3339 } from 'date-fns';
import { click } from '@ember/test-helpers';
import subMonths from 'date-fns/subMonths';
import timestamp from 'core/utils/timestamp';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SERIALIZED_ACTIVITY_RESPONSE } from 'vault/tests/helpers/clients/client-count-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { capabilitiesStub, overrideResponse } from 'vault/tests/helpers/stubs';

module('Integration | Component | clients/attribution', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    this.timestampStub = sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
  });

  hooks.beforeEach(function () {
    const { total, by_namespace } = SERIALIZED_ACTIVITY_RESPONSE;
    this.csvDownloadStub = sinon.stub(this.owner.lookup('service:download'), 'csv');
    const mockNow = this.timestampStub();
    this.mockNow = mockNow;
    this.startTimestamp = formatRFC3339(subMonths(mockNow, 6));
    this.timestamp = formatRFC3339(mockNow);
    this.selectedNamespace = null;
    this.totalUsageCounts = total;
    this.totalClientAttribution = [...by_namespace];
    this.namespaceMountsData = by_namespace.find((ns) => ns.label === 'ns1').mounts;
  });

  hooks.after(function () {
    this.csvDownloadStub.restore();
  });

  test('it renders empty state with no data', async function (assert) {
    await render(hbs`
      <Clients::Attribution />
    `);

    assert.dom('[data-test-component="empty-state"]').exists();
    assert.dom('[data-test-empty-state-title]').hasText('No data found');
    assert.dom('[data-test-attribution-description]').hasText('There is a problem gathering data');
    assert.dom('[data-test-attribution-export-button]').doesNotExist();
    assert.dom('[data-test-attribution-timestamp]').doesNotHaveTextContaining('Updated');
  });

  test('it renders with data for namespaces', async function (assert) {
    await render(hbs`
      <Clients::Attribution
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
    assert.dom('[data-test-top-attribution]').includesText('namespace').includesText('ns1');
    assert.dom('[data-test-attribution-clients]').includesText('namespace').includesText('18,903');
  });

  test('it renders two charts and correct text for single, historical month', async function (assert) {
    this.start = formatRFC3339(subMonths(this.mockNow, 1));
    this.end = formatRFC3339(subMonths(endOfMonth(this.mockNow), 1));
    await render(hbs`
      <Clients::Attribution
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
    this.set('selectedNamespace', 'ns1');

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
      <Clients::Attribution
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
      <Clients::Attribution
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
    this.set('selectedNamespace', 'ns1');
    await render(hbs`
      <Clients::Attribution
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
    assert.dom('[data-test-top-attribution]').includesText('auth method').includesText('auth/authid/0');
    assert.dom('[data-test-attribution-clients]').includesText('auth method').includesText('8,394');
  });

  module('export modal', function () {
    test('it renders modal', async function (assert) {
      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.namespaceMountsData}}
          @responseTimestamp={{this.timestamp}}
          @startTimestamp="2022-06-01T23:00:11.050Z"
          @endTimestamp="2022-12-01T23:00:11.050Z"
          />
      `);
      await click('[data-test-attribution-export-button]');
      assert
        .dom('[data-test-export-modal-title]')
        .hasText('Export activity data', 'modal appears to export data');
      assert.dom('[data-test-export-date-range]').includesText('June 2022 - December 2022');
    });

    test('it downloads csv data for date range', async function (assert) {
      assert.expect(2);
      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
          start_time: '2022-06-01T23:00:11.050Z',
          end_time: '2022-12-01T23:00:11.050Z',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          @startTimestamp="2022-06-01T23:00:11.050Z"
          @endTimestamp="2022-12-01T23:00:11.050Z"
          />
      `);
      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const [filename] = this.csvDownloadStub.lastCall.args;
      assert.strictEqual(
        filename,
        'clients_by_namespace_June 2022-December 2022',
        'csv has expected filename'
      );
    });

    test('it downloads csv data for a single month', async function (assert) {
      assert.expect(2);
      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
          start_time: '2022-06-01T23:00:11.050Z',
          end_time: '2022-06-21T23:00:11.050Z',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });
      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          @startTimestamp="2022-06-01T23:00:11.050Z"
          @endTimestamp="2022-06-21T23:00:11.050Z"
          />
      `);
      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const [filename] = this.csvDownloadStub.lastCall.args;
      assert.strictEqual(filename, 'clients_by_namespace_June 2022', 'csv has single month in filename');
    });
    test('csv filename omits date if no start/end timestamp', async function (assert) {
      assert.expect(2);

      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.queryParams, {
          format: 'csv',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          />
      `);

      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      const [filename] = this.csvDownloadStub.lastCall.args;
      assert.strictEqual(filename, 'clients_by_namespace');
    });

    test('shows the API error on the modal', async function (assert) {
      this.server.get('/sys/internal/counters/activity/export', function () {
        return new Response(
          403,
          { 'Content-Type': 'application/json' },
          { errors: ['this is an error from the API'] }
        );
      });

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          />
      `);

      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
      assert.dom('[data-test-export-error]').hasText('this is an error from the API');
    });

    test('hides the export button if user does not have SUDO capabilities', async function (assert) {
      this.server.post('/sys/capabilities-self', () =>
        capabilitiesStub('sys/internal/counters/activity/export', ['read'])
      );

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          />
      `);
      assert.dom('[data-test-attribution-export-button]').doesNotExist();
    });

    test('shows the export button if user has SUDO capabilities', async function (assert) {
      this.server.post('/sys/capabilities-self', () =>
        capabilitiesStub('sys/internal/counters/activity/export', ['sudo'])
      );

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          />
      `);
      assert.dom('[data-test-attribution-export-button]').exists();
    });

    test('defaults to show the export button if capabilities cannot be read', async function (assert) {
      this.server.post('/sys/capabilities-self', () => overrideResponse(403));

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          />
      `);
      assert.dom('[data-test-attribution-export-button]').exists();
    });

    test('it sends the current namespace in export request', async function (assert) {
      assert.expect(2);
      const namespaceSvc = this.owner.lookup('service:namespace');
      namespaceSvc.path = 'foo';
      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.requestHeaders, {
          'X-Vault-Namespace': 'foo',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          />
      `);
      assert.dom('[data-test-attribution-export-button]').exists();
      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
    });
    test('it sends the selected namespace in export request', async function (assert) {
      assert.expect(2);
      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.requestHeaders, {
          'X-Vault-Namespace': 'foobar',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });
      this.selectedNamespace = 'foobar/';

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          @selectedNamespace={{this.selectedNamespace}}
          />
      `);
      assert.dom('[data-test-attribution-export-button]').exists();
      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
    });

    test('it sends the current + selected namespace in export request', async function (assert) {
      assert.expect(2);
      const namespaceSvc = this.owner.lookup('service:namespace');
      namespaceSvc.path = 'foo';
      this.server.get('/sys/internal/counters/activity/export', function (_, req) {
        assert.deepEqual(req.requestHeaders, {
          'X-Vault-Namespace': 'foo/bar',
        });
        return new Response(200, { 'Content-Type': 'text/csv' }, '');
      });
      this.selectedNamespace = 'bar/';

      await render(hbs`
        <Clients::Attribution
          @totalClientAttribution={{this.totalClientAttribution}}
          @responseTimestamp={{this.timestamp}}
          @selectedNamespace={{this.selectedNamespace}}
          />
      `);
      assert.dom('[data-test-attribution-export-button]').exists();
      await click('[data-test-attribution-export-button]');
      await click(GENERAL.confirmButton);
    });
  });
});
