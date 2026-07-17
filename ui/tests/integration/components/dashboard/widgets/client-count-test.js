/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import sinon from 'sinon';
import { STATIC_NOW } from 'vault/mirage/handlers/clients';
import timestamp from 'core/utils/timestamp';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { dateFormat } from 'core/helpers/date-format';

module('Integration | Component | dashboard/widgets/client-count', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('it should display client count information', async function (assert) {
    assert.expect(5);
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW)); // 1/25/24
    this.server.get('sys/internal/counters/activity', () => {
      // this assertion should be hit twice, once initially and then again clicking 'refresh'
      assert.true(true, 'makes request to sys/internal/counters/activity');
      return { data: ACTIVITY_RESPONSE_STUB };
    });

    await render(hbs`<Dashboard::Widgets::ClientCount />`);
    assert.dom(GENERAL.textDisplay('Client count')).hasText('Client count');
    assert
      .dom(GENERAL.textBody('Client count total'))
      .hasText('Total: The number of clients in this billing period (Jun 2023 - Sep 2023).');
    assert
      .dom(GENERAL.textBody('Client count new'))
      .hasText('New: The number of clients new to Vault in the current month.');
    assert.dom(GENERAL.textBody('Client count updated at')).hasTextContaining(
      `Updated ${dateFormat([STATIC_NOW.toISOString(), 'MMMM d, yyyy, h:mm:ss aaa'], {
        withTimeZone: true,
      })}`
    );
  });

  test('it shows no data subtext if no start or end timestamp', async function (assert) {
    assert.expect(4);
    // as far as I know, responses will always have a start/end time
    // stubbing this unrealistic response just to test component subtext logic
    this.server.get('sys/internal/counters/activity', () => {
      return {
        data: { by_namespace: [], months: [], total: {} },
      };
    });

    await render(hbs`<Dashboard::Widgets::ClientCount />`);
    assert.dom(GENERAL.textBody('Client count total')).hasText('Total: No total client data available.');
    assert.dom(GENERAL.textBody('Client count new')).hasText('New: No new client data available.');
    assert.dom(GENERAL.tableData(0, 'total value')).hasText('-');
    assert.dom(GENERAL.tableData(1, 'new value')).hasText('-');
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
    await render(hbs`<Dashboard::Widgets::ClientCount />`);
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
    await render(hbs`<Dashboard::Widgets::ClientCount />`);
    assert.dom(GENERAL.emptyStateTitle).hasText('Activity configuration data is unavailable');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Reporting status is unknown and could be enabled or disabled. Check the Vault logs for more information.'
      );
  });
});
