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
import { SELECTORS } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';

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
        isPrimary: true,
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
    assert.dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert.dom(SELECTORS.replicationCard.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('running');
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle'))
      .exists();
    assert
      .dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'Performance primary'))
      .hasText('Performance primary');
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipTitle('dr-perf', 'Performance primary'))
      .hasText('running');
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'Performance primary', 'check-circle'))
      .exists();
  });
  test('it should display replication information if both dr and performance replication are enabled as features and only dr is setup', async function (assert) {
    this.replication = {
      dr: {
        clusterId: '123',
        state: 'running',
      },
      performance: {
        clusterId: '',
        isPrimary: true,
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
    assert.dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert.dom(SELECTORS.replicationCard.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('running');
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle'))
      .exists();
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle'))
      .hasClass('has-text-success');

    assert
      .dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'Performance primary'))
      .hasText('Performance primary');

    assert
      .dom(SELECTORS.replicationCard.getStateTooltipTitle('dr-perf', 'Performance primary'))
      .hasText('not set up');
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'Performance primary', 'x-circle'))
      .exists();
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'Performance primary', 'x-circle'))
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
    assert.dom(SELECTORS.replicationCard.getReplicationTitle('dr', 'state')).hasText('state');
    assert
      .dom(SELECTORS.replicationCard.drOnlyStateSubText)
      .hasText('The current operating state of the cluster.');
    assert.dom(SELECTORS.replicationCard.getStateTooltipTitle('dr', 'state')).hasText('running');
    assert.dom(SELECTORS.replicationCard.getStateTooltipIcon('dr', 'state', 'check-circle')).exists();
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr', 'state', 'check-circle'))
      .hasClass('has-text-success');
    assert.dom(SELECTORS.replicationCard.knownSecondariesLabel).hasText('known secondaries');
    assert
      .dom(SELECTORS.replicationCard.knownSecondariesSubtext)
      .hasText('Number of secondaries connected to this primary.');
    assert.dom(SELECTORS.replicationCard.knownSecondariesValue).hasText('1');
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
        isPrimary: true,
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
    assert.dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert.dom(SELECTORS.replicationCard.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('idle');
    assert.dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'DR primary', 'x-square')).exists();
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'DR primary', 'x-square'))
      .hasClass('has-text-danger');

    assert
      .dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'Performance primary'))
      .hasText('Performance primary');
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipTitle('dr-perf', 'Performance primary'))
      .hasText('shutdown');
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'Performance primary', 'x-circle'))
      .exists();
    assert
      .dom(SELECTORS.replicationCard.getStateTooltipIcon('dr-perf', 'Performance primary', 'x-circle'))
      .hasClass('has-text-danger');
  });

  test('it should show correct performance titles if primary vs secondary', async function (assert) {
    this.replication = {
      dr: {
        clusterId: 'abc',
        state: 'running',
      },
      performance: {
        clusterId: 'def',
        isPrimary: true,
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
    assert.dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'DR primary')).hasText('DR primary');
    assert
      .dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'Performance primary'))
      .hasText('Performance primary');

    this.replication = {
      dr: {
        clusterId: 'abc',
        state: 'running',
      },
      performance: {
        clusterId: 'def',
        isPrimary: false,
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
    assert
      .dom(SELECTORS.replicationCard.getReplicationTitle('dr-perf', 'Performance secondary'))
      .hasText('Performance secondary');
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
    assert.dom(SELECTORS.emptyState('replication')).exists();
    assert.dom(SELECTORS.emptyStateTitle('replication')).hasText('Replication not set up');
    assert
      .dom(SELECTORS.emptyStateMessage('replication'))
      .hasText('Data will be listed here. Enable a primary replication cluster to get started.');
    assert.dom(SELECTORS.emptyStateActions('replication')).hasText('Enable replication');
  });
});
