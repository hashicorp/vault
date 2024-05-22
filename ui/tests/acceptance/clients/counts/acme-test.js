/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit, click, currentURL } from '@ember/test-helpers';
import { getUnixTime } from 'date-fns';
import sinon from 'sinon';
import timestamp from 'core/utils/timestamp';
import authPage from 'vault/tests/pages/auth';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CHARTS, CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { ACTIVITY_RESPONSE_STUB, assertBarChart } from 'vault/tests/helpers/clients/client-count-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';

const { searchSelect } = GENERAL;

// integration test handle general display assertions, acceptance handles nav + filtering
module('Acceptance | clients | counts | acme', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    this.server.get('sys/internal/counters/activity', () => {
      return {
        request_id: 'some-activity-id',
        data: ACTIVITY_RESPONSE_STUB,
      };
    });
    // store serialized activity data for value comparison
    const { byMonth, byNamespace } = await this.owner
      .lookup('service:store')
      .queryRecord('clients/activity', {
        start_time: { timestamp: getUnixTime(LICENSE_START) },
        end_time: { timestamp: getUnixTime(STATIC_NOW) },
      });
    this.nsPath = 'ns1';
    this.mountPath = 'pki-engine-0';
    this.expectedValues = {
      nsTotals: byNamespace.find((ns) => ns.label === this.nsPath),
      nsMonthlyUsage: byMonth.map((m) => m?.namespaces_by_key[this.nsPath]).filter((d) => !!d),
      nsMonthActivity: byMonth.find(({ month }) => month === '9/23').namespaces_by_key[this.nsPath],
    };

    await authPage.login();
    return visit('/vault');
  });

  test('it navigates to acme tab', async function (assert) {
    assert.expect(3);
    await click(GENERAL.navLink('Client Count'));
    await click(GENERAL.tab('acme'));
    assert.strictEqual(currentURL(), '/vault/clients/counts/acme', 'it navigates to acme tab');
    assert.dom(GENERAL.tab('acme')).hasClass('active');
    await click(GENERAL.navLink('Back to main navigation'));
    assert.strictEqual(currentURL(), '/vault/dashboard', 'it navigates back to dashboard');
  });

  test('it filters by namespace data and renders charts', async function (assert) {
    const { nsTotals, nsMonthlyUsage, nsMonthActivity } = this.expectedValues;
    const nsMonthlyNew = nsMonthlyUsage.map((m) => m?.new_clients);
    assert.expect(7 + nsMonthlyUsage.length + nsMonthlyNew.length);

    await visit('/vault/clients/counts/acme');
    await click(searchSelect.trigger('namespace-search-select'));
    await click(searchSelect.option(searchSelect.optionIndex(this.nsPath)));

    // each chart assertion count is data array length + 2
    assertBarChart(assert, 'ACME usage', nsMonthlyUsage);
    assertBarChart(assert, 'Monthly new', nsMonthlyNew);
    assert.strictEqual(
      currentURL(),
      `/vault/clients/counts/acme?ns=${this.nsPath}`,
      'namespace filter updates URL query param'
    );
    assert
      .dom(CLIENT_COUNT.statText('Total ACME clients'))
      .hasTextContaining(
        `${formatNumber([nsTotals.acme_clients])}`,
        'renders total acme clients for namespace'
      );
    // there is only one month in the stubbed data, so in this case the average is the same as the total new clients
    assert
      .dom(CLIENT_COUNT.statText('Average new ACME clients per month'))
      .hasTextContaining(
        `${formatNumber([nsMonthActivity.new_clients.acme_clients])}`,
        'renders average acme clients for namespace'
      );
  });

  test('it filters by mount data and renders charts', async function (assert) {
    const { nsTotals, nsMonthlyUsage, nsMonthActivity } = this.expectedValues;
    const mountTotals = nsTotals.mounts.find((m) => m.label === this.mountPath);
    const mountMonthlyUsage = nsMonthlyUsage.map((ns) => ns.mounts_by_key[this.mountPath]).filter((d) => !!d);
    const mountMonthlyNew = mountMonthlyUsage.map((m) => m?.new_clients);
    assert.expect(7 + mountMonthlyUsage.length + mountMonthlyNew.length);

    await visit('/vault/clients/counts/acme');
    await click(searchSelect.trigger('namespace-search-select'));
    await click(searchSelect.option(searchSelect.optionIndex(this.nsPath)));
    await click(searchSelect.trigger('mounts-search-select'));
    await click(searchSelect.option(searchSelect.optionIndex(this.mountPath)));

    // each chart assertion count is data array length + 2
    assertBarChart(assert, 'ACME usage', mountMonthlyUsage);
    assertBarChart(assert, 'Monthly new', mountMonthlyNew);
    assert.strictEqual(
      currentURL(),
      `/vault/clients/counts/acme?mountPath=${this.mountPath}&ns=${this.nsPath}`,
      'mount filter updates URL query param'
    );
    assert
      .dom(CLIENT_COUNT.statText('Total ACME clients'))
      .hasTextContaining(
        `${formatNumber([mountTotals.acme_clients])}`,
        'renders total acme clients for mount'
      );
    // there is only one month in the stubbed data, so in this case the average is the same as the total new clients
    const mountMonthActivity = nsMonthActivity.mounts_by_key[this.mountPath];
    assert
      .dom(CLIENT_COUNT.statText('Average new ACME clients per month'))
      .hasTextContaining(
        `${formatNumber([mountMonthActivity.new_clients.acme_clients])}`,
        'renders average acme clients for mount'
      );
  });

  test('it renders empty chart for no mount data ', async function (assert) {
    assert.expect(3);
    await visit('/vault/clients/counts/acme');
    await click(searchSelect.trigger('namespace-search-select'));
    await click(searchSelect.option(searchSelect.optionIndex(this.nsPath)));
    await click(searchSelect.trigger('mounts-search-select'));
    // no data because this is an auth mount (acme_clients come from pki mounts)
    await click(searchSelect.option(searchSelect.optionIndex('auth/authid/0')));
    assert.dom(CLIENT_COUNT.statText('Total ACME clients')).hasTextContaining('0');
    assert.dom(`${CHARTS.chart('ACME usage')} ${CHARTS.verticalBar}`).isNotVisible();
    assert.dom(CHARTS.container('Monthly new')).doesNotExist();
  });
});
