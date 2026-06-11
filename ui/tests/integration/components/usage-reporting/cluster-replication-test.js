/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | usage-reporting/cluster-replication', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders DR and Performance rows', async function (assert) {
    await render(hbs`
      <UsageReporting::ClusterReplication
        @disasterRecoveryState="primary"
        @performanceState="secondary"
        @isVaultDedicated={{false}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-cluster-replication-dr-row]')
      .exists('DR row renders')
      .includesText('Disaster recovery', 'DR row label is correct');
    assert
      .dom('[data-test-vault-reporting-cluster-replication-perf-row]')
      .exists('Performance row renders')
      .includesText('Performance', 'Performance row label is correct');
  });

  test('it renders success badge for enabled states', async function (assert) {
    await render(hbs`
      <UsageReporting::ClusterReplication
        @disasterRecoveryState="primary"
        @performanceState="secondary"
        @isVaultDedicated={{false}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-cluster-replication-dr-badge]')
      .hasText('primary', 'DR badge shows state text');
    assert
      .dom('[data-test-vault-reporting-cluster-replication-perf-badge]')
      .hasText('secondary', 'Performance badge shows state text');
  });

  test('it renders disabled badge when state is "disabled"', async function (assert) {
    await render(hbs`
      <UsageReporting::ClusterReplication
        @disasterRecoveryState="disabled"
        @performanceState="disabled"
        @isVaultDedicated={{false}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-cluster-replication-dr-badge]')
      .hasText('disabled', 'DR badge shows disabled');
    assert
      .dom('[data-test-vault-reporting-cluster-replication-perf-badge]')
      .hasText('disabled', 'Performance badge shows disabled');
  });

  test('it renders the description link when both states are disabled', async function (assert) {
    await render(hbs`
      <UsageReporting::ClusterReplication
        @disasterRecoveryState="disabled"
        @performanceState="disabled"
        @isVaultDedicated={{false}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-cluster-replication-description-link]')
      .exists('replication docs link appears in empty description');
  });

  test('it renders the card title link when not vault dedicated', async function (assert) {
    await render(hbs`
      <UsageReporting::ClusterReplication
        @disasterRecoveryState="primary"
        @performanceState="disabled"
        @isVaultDedicated={{false}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-dashboard-card-title-link]')
      .exists('title row link renders for non-dedicated cluster');
  });

  test('it does not render the card title link when vault dedicated', async function (assert) {
    await render(hbs`
      <UsageReporting::ClusterReplication
        @disasterRecoveryState="primary"
        @performanceState="primary"
        @isVaultDedicated={{true}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-dashboard-card-title-link]')
      .doesNotExist('title row link is hidden for vault dedicated');
  });
});
