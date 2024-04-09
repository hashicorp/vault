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
import { ACTIVITY_RESPONSE_STUB } from 'vault/tests/helpers/clients';

module('Integration | Component | dashboard/client-count-card', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => STATIC_NOW);
  });

  hooks.beforeEach(function () {
    this.license = {
      startTime: LICENSE_START.toISOString(),
    };
  });

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it should display client count information', async function (assert) {
    assert.expect(9);
    this.server.get('sys/internal/counters/activity', () => {
      // this assertion should be hit twice, once initially and then again clicking 'refresh'
      assert.true(true, 'makes request to sys/internal/counters/activity');
      return {
        request_id: 'some-activity-id',
        data: ACTIVITY_RESPONSE_STUB,
      };
    });

    await render(hbs`<Dashboard::ClientCountCard @license={{this.license}} />`);
    assert.dom('[data-test-client-count-title]').hasText('Client count');
    assert.dom('[data-test-stat-text="total-clients"] .stat-label').hasText('Total');
    assert
      .dom('[data-test-stat-text="total-clients"] .stat-text')
      .hasText('The number of clients in this billing period (Jul 2023 - Jan 2024).');
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('7,805');
    assert.dom('[data-test-stat-text="new-clients"] .stat-label').hasText('New');
    assert
      .dom('[data-test-stat-text="new-clients"] .stat-text')
      .hasText('The number of clients new to Vault in the current month.');
    assert.dom('[data-test-stat-text="new-clients"] .stat-value').hasText('336');

    // fires second request to /activity
    await click('[data-test-refresh]');
  });
});
