import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const REPLICATION_DETAILS = {
  state: 'stream-wals',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

const IS_SYNCING = {
  state: 'merkle-diff',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

const IS_REINDEXING = {
  reindex_building_progress: 26838,
  reindex_building_total: 305443,
  reindex_in_progress: true,
  reindexing_stage: 'building',
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
    assert.dom('[data-test-replication-doc-link]').exists();
    assert.dom('[data-test-flash-message]').doesNotExist('no flash message is displayed on render');
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
    assert.dom('[data-test-isReindexing]').doesNotExist('does not show reindexing banner if cluster is ');
  });

  test('it shows an alert banner if the cluster is reindexing', async function(assert) {
    this.set('replicationDetails', IS_REINDEXING);

    await render(hbs`<ReplicationDashboard
      @replicationDetails={{replicationDetails}}
      @clusterMode={{clusterMode}}
      @isSecondary={{isSecondary}}
      @componentToRender='replication-secondary-card'
    />`);

    assert
      .dom('[data-test-isSyncing]')
      .doesNotExist('does not show syncing alert banner if cluster is not syncing');
    assert.dom('[data-test-isReindexing]').exists();
  });
});
