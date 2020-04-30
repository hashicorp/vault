import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const REPLICATION_DETAILS = {
  state: 'stream-wals',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

const REPLICATION_DETAILS_SYNCING = {
  state: 'merkle-diff',
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

module('Integration | Enterprise | Component | replication-dashboard', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('replicationDetailsSyncing', REPLICATION_DETAILS_SYNCING);
    this.set('componentToRender', 'replication-secondary-card');
    this.set('clusterMode', 'secondary');
    this.set('isSecondary', true);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationDashboard 
    @replicationDetails={{replicationDetails}} 
    />`);

    assert.dom('[data-test-replication-dashboard]').exists();
  });

  test('it renders table rows', async function(assert) {
    await render(hbs`<ReplicationDashboard @replicationDetails={{replicationDetails}}/>`);
    assert.dom('[data-test-table-rows').exists();
  });

  test('it renders with primary cluster address when set, and documentation link', async function(assert) {
    await render(hbs`<ReplicationDashboard 
    @replicationDetails={{replicationDetails}} 
    @clusterMode={{clusterMode}}
    @isSecondary={{isSecondary}}
    />`);

    assert
      .dom('[data-test-primary-cluster-address]')
      .includesText(
        REPLICATION_DETAILS.primaryClusterAddr,
        `shows the correct primary cluster address value`
      );

    assert.dom('[data-test-replication-doc-link]').exists();
  });

  test('it renders alert banner if state is merkle-diff and isSecondary', async function(assert) {
    await render(hbs`<ReplicationDashboard 
    @replicationDetails={{replicationDetailsSyncing}} 
    @clusterMode={{clusterMode}}
    @isSecondary={{true}}
    @componentToRender={{componentToRender}}
    />`);

    assert.dom('[data-test-isSyncing]').exists();
  });
});
