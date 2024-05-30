/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const REPLICATION_DETAILS = {
  state: 'stream-wals',
  primaryClusterAddr: 'https://127.0.0.1:8201',
  merkleRoot: '352f6e58ba2e8ec3935e05da1d142653dc76fe17',
  clusterId: '68999e13-a09d-b5e4-d66c-b35da566a177',
};

const IS_SYNCING = {
  state: 'merkle-diff',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

const IS_REINDEXING = {
  reindex_building_progress: 26838,
  reindex_building_total: 305443,
  reindex_in_progress: true,
  reindex_stage: 'building',
  state: 'running',
};

module('Integration | Component | replication-dashboard', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('clusterMode', 'secondary');
    this.set('isSecondary', true);
  });

  test('it renders', async function (assert) {
    await render(hbs`<ReplicationDashboard
      @replicationDetails={{this.replicationDetails}}
      @clusterMode={{this.clusterMode}}
      @isSecondary={{this.isSecondary}}
    />`);

    assert.dom('[data-test-replication-dashboard]').exists();
    assert.dom('[data-test-table-rows]').exists();
    assert.dom('[data-test-selectable-card-container="secondary"]').exists();
    assert.dom('[data-test-replication-doc-link]').exists();
    assert.dom('[data-test-flash-message]').doesNotExist('no flash message is displayed on render');
  });

  test('it updates the dashboard when the replication mode has changed', async function (assert) {
    await render(hbs`<ReplicationDashboard
      @replicationDetails={{this.replicationDetails}}
      @clusterMode={{this.clusterMode}}
      @isSecondary={{this.isSecondary}}
    />`);

    assert.dom('[data-test-selectable-card-container="secondary"]').exists();
    assert.dom('[data-test-selectable-card-container="primary"]').doesNotExist();

    this.set('clusterMode', 'primary');
    this.set('isSecondary', false);

    assert.dom('[data-test-selectable-card-container="primary"]').exists();
    assert.dom('[data-test-selectable-card-container="secondary"]').doesNotExist();
  });

  test('it renders the primary selectable-card-container when the cluster is a primary', async function (assert) {
    this.set('isSecondary', false);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{this.replicationDetails}}
      @clusterMode={{this.clusterMode}}
      @isSecondary={{this.isSecondary}}
    />`);

    assert.dom('[data-test-selectable-card-container="primary"]').exists();
    assert.dom('[data-test-selectable-card-container="secondary"]').doesNotExist();
  });

  test('it renders an alert banner if the dashboard is syncing', async function (assert) {
    this.set('replicationDetails', IS_SYNCING);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{this.replicationDetails}}
      @clusterMode={{this.clusterMode}}
      @isSecondary={{this.isSecondary}}
      @componentToRender='replication-secondary-card'
    />`);
    assert.dom('[data-test-isSyncing]').exists();
    assert
      .dom('[data-test-isReindexing]')
      .doesNotExist('does not show reindexing banner if cluster is cluster is not reindexing');
  });

  test('it shows an alert banner if the cluster is reindexing', async function (assert) {
    this.set('replicationDetails', IS_REINDEXING);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{this.replicationDetails}}
      @clusterMode={{this.clusterMode}}
      @isSecondary={{this.isSecondary}}
      @componentToRender='replication-secondary-card'
    />`);
    assert.dom('[data-test-isReindexing]').exists();
    assert
      .dom('[data-test-reindexing-title]')
      .includesText('Building', 'shows reindexing stage if there is one');
    assert
      .dom('[data-test-reindexing-progress]')
      .hasValue(
        IS_REINDEXING.reindex_building_progress,
        'shows the reindexing progress inside the alert banner'
      );
    const reindexingInProgress = { ...IS_REINDEXING, reindex_building_progress: 152721 };
    this.set('replicationDetails', reindexingInProgress);
    assert
      .dom('[data-test-reindexing-progress]')
      .hasValue(reindexingInProgress.reindex_building_progress, 'updates the reindexing progress');
  });

  test('it renders replication-summary-card when isSummaryDashboard', async function (assert) {
    const replicationDetailsSummary = {
      dr: {
        state: 'running',
        lastWAL: 10,
        knownSecondaries: ['https://127.0.0.1:8201', 'https://127.0.0.1:8202'],
      },
      performance: {
        state: 'running',
        lastWAL: 20,
        knownSecondaries: ['https://127.0.0.1:8201'],
      },
    };
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('replicationDetailsSummary', replicationDetailsSummary);
    this.set('isSummaryDashboard', true);
    this.set('clusterMode', 'primary');
    this.set('isSecondary', false);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{this.replicationDetails}}
      @replicationDetailsSummary={{this.replicationDetailsSummary}}
      @isSummaryDashboard={{this.isSummaryDashboard}}
      @clusterMode={{this.clusterMode}}
      @isSecondary={{this.isSecondary}}
      @componentToRender='replication-summary-card'
    />`);

    assert.dom('.summary-state').exists('it renders the summary dashboard');
    assert
      .dom('[data-test-summary-state]')
      .includesText(replicationDetailsSummary.dr.state, 'shows the correct state value');
    assert.dom('[data-test-icon]').exists('shows an icon if state is ok');
    assert
      .dom('[data-test-selectable-card-container-summary]')
      .exists('it renders with the correct selectable card container');
    assert.dom('[data-test-selectable-card-container-primary]').doesNotExist();
  });

  test('it renders replication-summary-card with an error message when the state is not OK', async function (assert) {
    const replicationDetailsSummary = {
      dr: {
        state: 'shutdown',
        lastWAL: 10,
        knownSecondaries: ['https://127.0.0.1:8201', 'https://127.0.0.1:8202'],
      },
      performance: {
        state: 'shutdown',
        lastWAL: 20,
        knownSecondaries: ['https://127.0.0.1:8201'],
      },
    };
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('replicationDetailsSummary', replicationDetailsSummary);
    this.set('isSummaryDashboard', true);
    this.set('clusterMode', 'primary');
    this.set('isSecondary', false);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{this.replicationDetails}}
      @replicationDetailsSummary={{this.replicationDetailsSummary}}
      @isSummaryDashboard={{this.isSummaryDashboard}}
      @clusterMode={{this.clusterMode}}
      @isSecondary={{this.isSecondary}}
      @componentToRender='replication-summary-card'
    />`);

    assert.dom('[data-test-error]').includesText('state', 'show correct error title');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('The cluster is shutdown. Please check your server logs.', 'show correct error message');
  });
});
