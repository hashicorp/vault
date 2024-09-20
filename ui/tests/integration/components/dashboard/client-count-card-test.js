/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import sinon from 'sinon';
import { STATIC_NOW } from 'vault/mirage/handlers/clients';
import timestamp from 'core/utils/timestamp';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | dashboard/client-count-card', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('it should display client count information', async function (assert) {
    assert.expect(6);
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW)); // 1/25/24
    const { months, total } = ACTIVITY_RESPONSE_STUB;
    const [latestMonth] = months.slice(-1);
    this.server.get('sys/internal/counters/activity', () => {
      // this assertion should be hit twice, once initially and then again clicking 'refresh'
      assert.true(true, 'makes request to sys/internal/counters/activity');
      return { data: ACTIVITY_RESPONSE_STUB };
    });

    await render(hbs`<Dashboard::ClientCountCard />`);
    assert.dom('[data-test-client-count-title]').hasText('Client count');
    assert
      .dom(CLIENT_COUNT.statText('Total'))
      .hasText(
        `Total The number of clients in this billing period (Jun 2023 - Sep 2023). ${formatNumber([
          total.clients,
        ])}`
      );
    assert
      .dom(CLIENT_COUNT.statText('New'))
      .hasText(
        `New The number of clients new to Vault in the current month. ${formatNumber([
          latestMonth.new_clients.counts.clients,
        ])}`
      );
    assert.dom('[data-test-updated-timestamp]').hasTextContaining('Updated Jan 25 2024');

    // fires second request to /activity
    await click('[data-test-refresh]');
  });

  test('it shows no data subtext if no start or end timestamp', async function (assert) {
    assert.expect(2);
    // as far as I know, responses will always have a start/end time
    // stubbing this unrealistic response just to test component subtext logic
    this.server.get('sys/internal/counters/activity', () => {
      return {
        data: { by_namespace: [], months: [], total: {} },
      };
    });

    await render(hbs`<Dashboard::ClientCountCard />`);
    assert.dom(CLIENT_COUNT.statText('Total')).hasText('Total No total client data available. -');
    assert.dom(CLIENT_COUNT.statText('New')).hasText('New No new client data available. -');
  });

  test('it shows empty state if no activity data and reporting is enabled', async function (assert) {
    // the activity response has changed and now should ALWAYS return something
    // but adding this test until we update the adapter to reflect that
    assert.expect(4);
    this.server.get('sys/internal/counters/activity', () => {
      assert.true(true, 'makes request to sys/internal/counters/activity');
      return { data: {} };
    });
    this.server.get('sys/internal/counters/config', () => {
      assert.true(true, 'makes request to sys/internal/counters/config');
      return {
        request_id: '25a94b99-b49a-c4ac-cb7b-5ba0eb390a25',
        data: { reporting_enabled: true },
      };
    });
    await render(hbs`<Dashboard::ClientCountCard />`);
    assert.dom(GENERAL.emptyStateTitle).hasText('No data received');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('Tracking is turned on and Vault is gathering data. It should appear here within 30 minutes.');
  });

  test('it shows empty state if no activity data and config data is unavailable', async function (assert) {
    assert.expect(4);
    this.server.get('sys/internal/counters/activity', () => {
      assert.true(true, 'makes request to sys/internal/counters/activity');
      return { data: {} };
    });
    this.server.get('sys/internal/counters/config', () => {
      assert.true(true, 'makes request to sys/internal/counters/config');
      return new Response(
        403,
        { 'Content-Type': 'application/json' },
        JSON.stringify({ errors: ['permission denied'] })
      );
    });
    await render(hbs`<Dashboard::ClientCountCard />`);
    assert.dom(GENERAL.emptyStateTitle).hasText('Activity configuration data is unavailable');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Reporting status is unknown and could be enabled or disabled. Check the Vault logs for more information.'
      );
  });
});
