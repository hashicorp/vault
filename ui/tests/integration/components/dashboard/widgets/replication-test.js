/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import timestamp from 'core/utils/timestamp';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | dashboard/widgets/replication', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.replication = {
      dr: {
        clusterId: '123',
        state: 'running',
        mode: 'primary',
      },
      performance: {
        clusterId: 'abc-1',
        state: 'running',
        mode: 'primary',
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
        <Dashboard::Widgets::Replication
          @replication={{this.replication}}
          @version={{this.version}}
          @updatedAt={{this.updatedAt}}
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(GENERAL.tableData('0', 'DR')).hasText('Disaster recovery primary Known secondaries');
    assert.dom(GENERAL.tableData('0', 'DR value')).hasText('Running 0');
    assert.dom(GENERAL.badge('DR')).hasClass('hds-badge--color-success');
    assert.dom(GENERAL.tableData('1', 'Performance')).hasText('Performance replication primary');
    assert.dom(GENERAL.tableData('1', 'Performance value')).hasText('Running');
    assert.dom(GENERAL.badge('Performance')).hasClass('hds-badge--color-success');
  });

  test('it should display replication information if both dr and performance replication are enabled as features and only dr is setup', async function (assert) {
    this.replication.performance = { mode: 'disabled' };
    await render(
      hbs`
        <Dashboard::Widgets::Replication
          @replication={{this.replication}}
          @version={{this.version}}
          @updatedAt={{this.updatedAt}}
          @refresh={{this.refresh}} />
          `
    );

    assert.dom(GENERAL.tableData('0', 'DR')).hasText('Disaster recovery primary Known secondaries');
    assert.dom(GENERAL.tableData('0', 'DR value')).hasText('Running 0');
    assert.dom(GENERAL.badge('DR')).hasClass('hds-badge--color-success');
    assert.dom(GENERAL.tableData('1', 'Performance')).hasText('Performance replication');
    assert.dom(GENERAL.tableData('1', 'Performance value')).hasText('Not set up');
    assert.dom(GENERAL.badge('Performance')).hasClass('hds-badge--color-warning');
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
        <Dashboard::Widgets::Replication
          @replication={{this.replication}}
          @version={{this.version}}
          @updatedAt={{this.updatedAt}}
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(GENERAL.tableData('0', 'DR')).hasText('Disaster recovery Known secondaries');
    assert.dom(GENERAL.tableData('0', 'DR value')).hasText('Running 1');
    assert.dom(GENERAL.badge('DR')).hasClass('hds-badge--color-success');
    assert.dom(GENERAL.tableData('1', 'Performance')).doesNotExist();
  });

  test('it should show correct icons if dr and performance replication is idle or shutdown states', async function (assert) {
    this.replication = {
      dr: {
        clusterId: 'abc',
        state: 'idle',
        mode: 'primary',
      },
      performance: {
        clusterId: 'def',
        state: 'shutdown',
        mode: 'primary',
      },
    };
    await render(
      hbs`
        <Dashboard::Widgets::Replication
          @replication={{this.replication}}
          @version={{this.version}}
          @updatedAt={{this.updatedAt}}
          @refresh={{this.refresh}} />
          `
    );

    assert.dom(GENERAL.tableData('0', 'DR')).hasText('Disaster recovery primary Known secondaries');
    assert.dom(GENERAL.tableData('0', 'DR value')).hasText('Idle 0');
    assert.dom(GENERAL.badge('DR')).hasClass('hds-badge--color-critical');

    assert.dom(GENERAL.tableData('1', 'Performance')).hasText('Performance replication primary');
    assert.dom(GENERAL.tableData('1', 'Performance value')).hasText('Shutdown');
    assert.dom(GENERAL.badge('Performance')).hasClass('hds-badge--color-critical');
  });

  test('it should show correct performance titles if primary vs secondary', async function (assert) {
    this.replication = {
      dr: {
        clusterId: 'abc',
        state: 'running',
        mode: 'primary',
      },
      performance: {
        clusterId: 'def',
        mode: 'primary',
      },
    };
    await render(
      hbs`
        <Dashboard::Widgets::Replication
          @replication={{this.replication}}
          @version={{this.version}}
          @updatedAt={{this.updatedAt}}
          @refresh={{this.refresh}} />
          `
    );
    assert.dom(GENERAL.tableData('0', 'DR')).hasText('Disaster recovery primary Known secondaries');
    assert.dom(GENERAL.tableData('1', 'Performance')).hasText('Performance replication primary');

    this.replication.performance.mode = 'secondary';
    await render(
      hbs`
          <Dashboard::Widgets::Replication
            @replication={{this.replication}}
            @version={{this.version}}
            @updatedAt={{this.updatedAt}}
            @refresh={{this.refresh}} />
            `
    );
    assert.dom(GENERAL.tableData('1', 'Performance')).hasText('Performance replication secondary');
  });

  test('it should show replication card empty table and Enable replication button', async function (assert) {
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
        <Dashboard::Widgets::Replication
          @replication={{this.replication}}
          @version={{this.version}}
          @updatedAt={{this.updatedAt}}
          @refresh={{this.refresh}} />
          `
    );

    assert.dom(GENERAL.tableData('0', 'DR')).hasText('Disaster recovery Known secondaries');
    assert.dom(GENERAL.tableData('0', 'DR value')).hasText('Not set up 0');
    assert.dom(GENERAL.badge('DR')).hasClass('hds-badge--color-warning');
    assert.dom(GENERAL.tableData('1', 'Performance')).hasText('Performance replication');
    assert.dom(GENERAL.tableData('1', 'Performance value')).hasText('Not set up');
    assert.dom(GENERAL.badge('Performance')).hasClass('hds-badge--color-warning');
  });
});
