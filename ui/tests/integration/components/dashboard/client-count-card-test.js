/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import timestamp from 'core/utils/timestamp';
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients/client-count-helpers';
import { formatNumber } from 'core/helpers/format-number';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

module('Integration | Component | dashboard/client-count-card', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => STATIC_NOW);
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it should display client count information', async function (assert) {
    assert.expect(6);
    const { months, total } = ACTIVITY_RESPONSE_STUB;
    const [latestMonth] = months.slice(-1);
    this.server.get('sys/internal/counters/activity', () => {
      // this assertion should be hit twice, once initially and then again clicking 'refresh'
      assert.true(true, 'makes request to sys/internal/counters/activity');
      return {
        request_id: 'some-activity-id',
        data: ACTIVITY_RESPONSE_STUB,
      };
    });
    this.server.get('sys/internal/counters/config', function () {
      assert.true(true, 'sys/internal/counters/config');
      return {
        request_id: 'some-config-id',
        data: {
          billing_start_timestamp: LICENSE_START.toISOString(),
        },
      };
    });

    await render(hbs`<Dashboard::ClientCountCard @isEnterprise={{true}} />`);
    assert.dom('[data-test-client-count-title]').hasText('Client count');
    assert
      .dom(CLIENT_COUNT.statText('Total'))
      .hasText(
        `Total The number of clients in this billing period (Jul 2023 - Jan 2024). ${formatNumber([
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

    // fires second request to /activity
    await click('[data-test-refresh]');
  });

  test('it does not query activity for community edition', async function (assert) {
    assert.expect(4);
    this.server.get('sys/internal/counters/activity', () => {
      // this assertion should NOT be hit in this test
      assert.true(true, 'uh oh! makes request to sys/internal/counters/activity');
      return {
        request_id: 'some-activity-id',
        data: ACTIVITY_RESPONSE_STUB,
      };
    });
    this.server.get('sys/internal/counters/config', function () {
      assert.true(true, 'sys/internal/counters/config');
      return {
        request_id: 'some-config-id',
        data: {
          billing_start_timestamp: '0001-01-01T00:00:00Z',
        },
      };
    });

    await render(hbs`<Dashboard::ClientCountCard @isEnterprise={{false}} />`);
    assert.dom(CLIENT_COUNT.statText('Total')).hasText('Total No total client data available. -');
    assert.dom(CLIENT_COUNT.statText('New')).hasText('New No new client data available. -');

    // attempt second request to /activity but component task should return instead of hitting endpoint
    await click('[data-test-refresh]');
  });
});
