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

module('Integration | Enterprise | Component | replication-dashboard', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('clusterMode', 'secondary');
    this.set('isSecondary', true);
  });

  test('it renders', async function(assert) {
    this.set('replicationDetails', REPLICATION_DETAILS);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{replicationDetails}}
      @clusterMode={{clusterMode}}
      @isSecondary={{isSecondary}}
    />`);

    assert.dom('[data-test-replication-dashboard]').exists();
    assert.dom('[data-test-table-rows').exists();
    assert.dom('[data-test-selectable-card-container="secondary"]').exists();
    assert.dom('[data-test-replication-doc-link]').exists();
    assert.dom('[data-test-flash-message]').doesNotExist('no flash message is displayed on render');
  });

  test('it renders the primary selectable-card-container when the cluster is a primary', async function(assert) {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('isSecondary', false);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{replicationDetails}}
      @clusterMode={{clusterMode}}
      @isSecondary={{isSecondary}}
    />`);

    assert.dom('[data-test-selectable-card-container="primary"]').exists();
  });

  test('it renders an alert banner if the dashboard is syncing', async function(assert) {
    this.set('replicationDetails', IS_SYNCING);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{replicationDetails}}
      @clusterMode={{clusterMode}}
      @isSecondary={{isSecondary}}
      @componentToRender='replication-secondary-card'
    />`);

    assert.dom('[data-test-isSyncing]').exists();
    assert
      .dom('[data-test-isReindexing]')
      .doesNotExist('does not show reindexing banner if cluster is cluster is not reindexing');
  });

  test('it shows an alert banner if the cluster is reindexing', async function(assert) {
    this.set('replicationDetails', IS_REINDEXING);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{replicationDetails}}
      @clusterMode={{clusterMode}}
      @isSecondary={{isSecondary}}
      @componentToRender='replication-secondary-card'
    />`);

    assert.dom('[data-test-isReindexing]').exists();
    assert.dom('.message-title').includesText('Building', 'shows reindexing stage if there is one');
    assert
      .dom('.message-title>.progress')
      .hasValue(
        IS_REINDEXING.reindex_building_progress,
        'shows the reindexing progress inside the alert banner'
      );
  });
});
