/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import timestamp from 'core/utils/timestamp';

const SELECTORS = {
  getReplicationTitle: (type, name) => `[data-test-${type}-replication] [data-test-title="${name}"]`,
  getStateTooltipTitle: (type, name) => `[data-test-${type}-replication] [data-test-tooltip-title="${name}"]`,
  getStateTooltipIcon: (type, name, icon) =>
    `[data-test-${type}-replication] [data-test-tooltip-title="${name}"] [data-test-icon="${icon}"]`,
  drOnlyStateSubText: '[data-test-dr-replication] [data-test-subtext="state"]',
  knownSecondariesLabel: '[data-test-stat-text="known secondaries"] .stat-label',
  knownSecondariesSubtext: '[data-test-stat-text="known secondaries"] .stat-text',
  knownSecondariesValue: '[data-test-stat-text="known secondaries"] .stat-value',
  replicationEmptyState: '[data-test-component="empty-state"]',
  replicationEmptyStateTitle: '[data-test-component="empty-state"] .empty-state-title',
  replicationEmptyStateMessage: '[data-test-component="empty-state"] .empty-state-message',
  replicationEmptyStateActions: '[data-test-component="empty-state"] .empty-state-actions',
};

module('Integration | Component | dashboard/replication-card', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.replication = {
      dr: {
        clusterId: '123',
        state: 'running',
      },
      performance: {
        clusterId: 'abc-1',
        state: 'running',
      },
    };
    this.version = {
      hasPerfReplication: true,
      hasDRReplication: true,
    };
    this.updatedAt = timestamp.now().toISOString();
    this.refresh = () => {};
  });

  test('it should display replication information if both dr and performance replication are enabled as features', async function (assert) {
    await render(
      hbs`
        <Dashboard::ReplicationCard 
          @replication={{this.replication}} 
          @version={{this.version}} 
          @updatedAt={{this.updatedAt}} 
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(SELECTORS.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('running');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle')).exists();
    assert.dom(SELECTORS.getReplicationTitle('dr-perf', 'Perf primary')).hasText('Perf primary');
    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'Perf primary')).hasText('running');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'Perf primary', 'check-circle')).exists();
  });
  test('it should display replication information if both dr and performance replication are enabled as features and only dr is setup', async function (assert) {
    this.replication = {
      dr: {
        clusterId: '123',
        state: 'running',
      },
      performance: {
        clusterId: '',
      },
    };
    await render(
      hbs`
        <Dashboard::ReplicationCard 
          @replication={{this.replication}} 
          @version={{this.version}} 
          @updatedAt={{this.updatedAt}} 
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(SELECTORS.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('running');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle')).exists();
    assert
      .dom(SELECTORS.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle'))
      .hasClass('has-text-success');

    assert.dom(SELECTORS.getReplicationTitle('dr-perf', 'Perf primary')).hasText('Perf primary');

    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'Perf primary')).hasText('not set up');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'Perf primary', 'x-circle')).exists();
    assert
      .dom(SELECTORS.getStateTooltipIcon('dr-perf', 'Perf primary', 'x-circle'))
      .hasClass('has-text-danger');
  });

  test('it should display only dr replication information if vault version only has hasDRReplication', async function (assert) {
    this.version = {
      hasPerfReplication: false,
      hasDRReplication: true,
    };
    this.replication = {
      dr: {
        clusterId: '123',
        state: 'running',
        knownSecondaries: [{ id: 1 }],
      },
    };
    await render(
      hbs`
        <Dashboard::ReplicationCard 
          @replication={{this.replication}} 
          @version={{this.version}} 
          @updatedAt={{this.updatedAt}} 
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(SELECTORS.getReplicationTitle('dr', 'state')).hasText('state');
    assert.dom(SELECTORS.drOnlyStateSubText).hasText('The current operating state of the cluster.');
    assert.dom(SELECTORS.getStateTooltipTitle('dr', 'state')).hasText('running');
    assert.dom(SELECTORS.getStateTooltipIcon('dr', 'state', 'check-circle')).exists();
    assert.dom(SELECTORS.getStateTooltipIcon('dr', 'state', 'check-circle')).hasClass('has-text-success');
    assert.dom(SELECTORS.knownSecondariesLabel).hasText('known secondaries');
    assert.dom(SELECTORS.knownSecondariesSubtext).hasText('Number of secondaries connected to this primary.');
    assert.dom(SELECTORS.knownSecondariesValue).hasText('1');
  });

  test('it should show correct icons if dr and performance replication is idle or shutdown states', async function (assert) {
    this.replication = {
      dr: {
        clusterId: 'abc',
        state: 'idle',
      },
      performance: {
        clusterId: 'def',
        state: 'shutdown',
      },
    };
    await render(
      hbs`
        <Dashboard::ReplicationCard 
          @replication={{this.replication}} 
          @version={{this.version}} 
          @updatedAt={{this.updatedAt}} 
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(SELECTORS.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('idle');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'DR primary', 'x-square')).exists();
    assert
      .dom(SELECTORS.getStateTooltipIcon('dr-perf', 'DR primary', 'x-square'))
      .hasClass('has-text-danger');

    assert.dom(SELECTORS.getReplicationTitle('dr-perf', 'Perf primary')).hasText('Perf primary');
    assert.dom(SELECTORS.getStateTooltipTitle('dr-perf', 'Perf primary')).hasText('shutdown');
    assert.dom(SELECTORS.getStateTooltipIcon('dr-perf', 'Perf primary', 'x-circle')).exists();
    assert
      .dom(SELECTORS.getStateTooltipIcon('dr-perf', 'Perf primary', 'x-circle'))
      .hasClass('has-text-danger');
  });

  test('it should show empty state', async function (assert) {
    this.replication = {
      dr: {
        clusterId: '',
      },
      performance: {
        clusterId: '',
      },
    };
    await render(
      hbs`
        <Dashboard::ReplicationCard 
          @replication={{this.replication}} 
          @version={{this.version}} 
          @updatedAt={{this.updatedAt}} 
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(SELECTORS.replicationEmptyState).exists();
    assert.dom(SELECTORS.replicationEmptyStateTitle).hasText('Replication not set up');
    assert
      .dom(SELECTORS.replicationEmptyStateMessage)
      .hasText('Data will be listed here. Enable a primary replication cluster to get started.');
    assert.dom(SELECTORS.replicationEmptyStateActions).hasText('Enable replication');
  });
});
