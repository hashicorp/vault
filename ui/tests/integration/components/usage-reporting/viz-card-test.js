/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | usage-reporting/viz-card', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders empty state when all values are zero', async function (assert) {
    this.set('data', [
      { label: 'kubernetes', value: 0 },
      { label: 'token', value: 0 },
    ]);

    await render(hbs`
      <UsageReporting::VizCard
        @data={{this.data}}
        @title="Authentication methods"
        @description="Enabled authentication methods for this cluster."
      />
    `);

    assert
      .dom('[data-test-vault-reporting-horizontal-bar-chart-empty-state]')
      .exists('renders empty state when all values are filtered out');
    assert
      .dom('[data-test-vault-reporting-horizontal-bar-chart-carbon]')
      .doesNotExist('does not render chart when all values are zero');
  });
});
