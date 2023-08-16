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
    assert.dom('[data-test-dr-perf-replication] [data-test-dr-title]').hasText('DR primary');
    assert
      .dom('[data-test-dr-perf-replication] [data-test-stat-text="dr state"] .stat-value')
      .hasText('running');
    assert
      .dom('[data-test-dr-perf-replication] [data-test-stat-text="dr state"] [data-test-icon="check-circle"]')
      .exists();
    assert.dom('[data-test-dr-perf-replication] [data-test-performance-title]').hasText('Perf primary');
    assert
      .dom('[data-test-dr-perf-replication] [data-test-stat-text="performance state"] .stat-value')
      .hasText('running');
    assert
      .dom(
        '[data-test-dr-perf-replication] [data-test-stat-text="performance state"] [data-test-icon="check-circle"]'
      )
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
    assert.dom('[data-test-dr-perf-replication] [data-test-dr-title]').hasText('DR primary');
    assert
      .dom('[data-test-dr-perf-replication] [data-test-stat-text="dr state"] .stat-value')
      .hasText('running');
    assert
      .dom('[data-test-dr-perf-replication] [data-test-stat-text="dr state"] [data-test-icon="check-circle"]')
      .exists();

    assert.dom('[data-test-dr-perf-replication] [data-test-performance-title]').hasText('Perf primary');

    assert
      .dom('[data-test-dr-perf-replication] [data-test-stat-text="performance state"] .stat-value')
      .hasText('not set up');
    assert
      .dom(
        '[data-test-dr-perf-replication] [data-test-stat-text="performance state"] [data-test-icon="x-circle-fill"]'
      )
      .exists();
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
    assert.dom('[data-test-stat-text="dr state"] .stat-label').hasText('state');
    assert
      .dom('[data-test-stat-text="dr state"] .stat-text')
      .hasText('The current operating state of the cluster.');
    assert.dom('[data-test-stat-text="dr state"] .stat-value').hasText('running');
    assert.dom('[data-test-stat-text="dr state"] .stat-value [data-test-icon="check-circle"]').exists();
    assert.dom('[data-test-stat-text="known secondaries"] .stat-label').hasText('known secondaries');
    assert
      .dom('[data-test-stat-text="known secondaries"] .stat-text')
      .hasText('Number of secondaries connected to this primary.');
    assert.dom('[data-test-stat-text="known secondaries"] .stat-value').hasText('1');
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
    assert.dom('[data-test-component="empty-state"]').exists();
    assert.dom('[data-test-component="empty-state"] .empty-state-title').hasText('Replication not set up');
    assert
      .dom('[data-test-component="empty-state"] .empty-state-message')
      .hasText('Data will be listed here. Enable a primary replication cluster to get started.');
    assert.dom('[data-test-component="empty-state"] .empty-state-actions').hasText('Details');
  });
});
