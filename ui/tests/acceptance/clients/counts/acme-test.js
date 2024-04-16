/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { visit, click, currentURL } from '@ember/test-helpers';
import { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import { getUnixTime } from 'date-fns';
import sinon from 'sinon';
import timestamp from 'core/utils/timestamp';
import authPage from 'vault/tests/pages/auth';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { ACTIVITY_RESPONSE_STUB, assertChart } from 'vault/tests/helpers/clients/client-count-helpers';
import { formatNumber } from 'core/helpers/format-number';

const { searchSelect } = GENERAL;

// integration test handle general display assertions, acceptance handles nav + filtering
module('Acceptance | clients | counts | acme', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => STATIC_NOW);
  });

  hooks.beforeEach(async function () {
    this.server.get('sys/internal/counters/activity', () => {
      return {
        request_id: 'some-activity-id',
        data: ACTIVITY_RESPONSE_STUB,
      };
    });
    const query = {
      start_time: { timestamp: getUnixTime(LICENSE_START) },
      end_time: { timestamp: getUnixTime(STATIC_NOW) },
    };
    // store serialized activity data for value comparison
    this.activity = await this.owner.lookup('service:store').queryRecord('clients/activity', query);
    await authPage.login();
    return visit('/vault');
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it navigates to acme tab', async function (assert) {
    assert.expect(3);
    await click(GENERAL.navLink('Client Count'));
    await click(GENERAL.tab('acme'));
    assert.strictEqual(currentURL(), '/vault/clients/counts/acme');
    assert.dom(GENERAL.tab('acme')).hasClass('active');
    await click(GENERAL.navLink('Back to main navigation'));
    assert.strictEqual(currentURL(), '/vault');
  });

  test('it filters data and renders charts', async function (assert) {
    assert.expect(20);
    await visit('/vault/clients/counts/acme');
    const nsPath = 'ns1';
    const mountPath = 'pki-engine-0';
    const nsData = this.activity.byNamespace.find((ns) => ns.label === nsPath);
    const mountData = this.activity.byMonth.find(({ month }) => month === '9/23').namespaces_by_key[nsPath]
      .mounts_by_key[mountPath];

    // filter by namespace
    await click(searchSelect.trigger('namespace-search-select'));
    await click(searchSelect.option(searchSelect.optionIndex(nsPath)));
    assert.strictEqual(
      currentURL(),
      `/vault/clients/counts/acme?ns=${nsPath}`,
      'namespace filter updates URL query param'
    );
    assert
      .dom(CLIENT_COUNT.statText('Total ACME clients'))
      .hasTextContaining(
        `${formatNumber([nsData.acme_clients])}`,
        'renders total acme clients for namespace'
      );
    const monthlyNsData = this.activity.byMonth.map((m) => m?.namespaces_by_key[nsPath]).filter((d) => !!d);
    const monthlyNewNsData = monthlyNsData.map((m) => m?.new_clients);
    // each chart assertion count is data array length + 2
    assertChart(assert, 'ACME usage', monthlyNsData);
    assertChart(assert, 'Monthly new', monthlyNewNsData);

    // filter by mount
    await click(searchSelect.trigger('mounts-search-select'));
    await click(searchSelect.option(searchSelect.optionIndex(mountPath)));
    assert.strictEqual(
      currentURL(),
      `/vault/clients/counts/acme?mountPath=${mountPath}&ns=${nsPath}`,
      'mount filter updates URL query param'
    );
    assert
      .dom(CLIENT_COUNT.statText('Total ACME clients'))
      .hasTextContaining(`${formatNumber([mountData.acme_clients])}`, 'renders total acme clients for mount');

    const monthlyMountData = monthlyNsData.map((ns) => ns.mounts_by_key[mountPath]).filter((d) => !!d);
    const monthlyNewMountData = monthlyMountData.map((m) => m?.new_clients);
    // there is only one month in the chart, so in this case the average is the same as the total new clients
    assert
      .dom(CLIENT_COUNT.statText('Average new ACME clients per month'))
      .hasTextContaining(
        `${formatNumber([mountData.new_clients.acme_clients])}`,
        'renders average acme clients for mount'
      );
    // each chart assertion count is data array length + 2
    assertChart(assert, 'ACME usage', monthlyMountData);
    assertChart(assert, 'Monthly new', monthlyNewMountData);
    await click(searchSelect.removeSelected);
    await click(searchSelect.trigger('namespace-search-select'));
    await click(searchSelect.option(searchSelect.optionIndex(nsPath)));
    await click(searchSelect.trigger('mounts-search-select'));
    // no data because this is an auth mount (acme_clients come from pki mounts)
    await click(searchSelect.option(searchSelect.optionIndex('auth/authid0')));
    assert.dom(CLIENT_COUNT.statText('Total ACME clients')).hasTextContaining('0');
    assert.dom(`${CLIENT_COUNT.charts.chart('ACME usage')} ${CLIENT_COUNT.charts.dataBar}`).isNotVisible();
    assert.dom(CLIENT_COUNT.charts.chart('Monthly new')).doesNotExist();
  });
});
