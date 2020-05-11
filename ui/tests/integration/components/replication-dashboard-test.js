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

const REPLICATION_DETAILS_SYNCING = {
  state: 'merkle-diff',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

const REPLICATION_DETAILS_REINDEXING = {
  state: 'stream-wals',
  reindex_in_progress: true,
  primaryClusterAddr: 'https://127.0.0.1:8201',
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

  test('it renders an alert banner if the dashboard is syncing', async function(assert) {
    this.set('replicationDetailsSyncing', REPLICATION_DETAILS_SYNCING);

    await render(hbs`<ReplicationDashboard 
    @replicationDetails={{replicationDetailsSyncing}} 
    @clusterMode={{clusterMode}}
    @isSecondary={{isSecondary}}
    @componentToRender='replication-secondary-card'
    />`);

    assert.dom('[data-test-isSyncing]').exists();
    assert.dom('[data-test-isReindexing]').doesNotExist();
  });

  test('it renders an alert banner if the dashboard is reIndexing', async function(assert) {
    this.set('replicationDetailsReindexing', REPLICATION_DETAILS_REINDEXING);

    await render(hbs`<ReplicationDashboard 
    @replicationDetails={{replicationDetailsReindexing}} 
    @clusterMode={{clusterMode}}
    @isSecondary={{isSecondary}}
    @componentToRender='replication-primary-card'
    />`);

    assert.dom('[data-test-isSyncing]').doesNotExist();
    assert.dom('[data-test-isReindexing]').exists();
  });
});
