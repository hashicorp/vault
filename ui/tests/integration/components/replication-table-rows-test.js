import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const REPLICATION_DETAILS = {
  clusterId: 'b829d963-6835-33eb-a903-57376024b97a',
  merkleRoot: 'c21c8428a0a06135cef6ae25bf8e0267ff1592a6',
};

const CLUSTER_MODE = 'primary';

module('Integration | Component | replication-table-rows', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('clusterMode', CLUSTER_MODE);
  });

  test('it renders', async function(assert) {
    await render(
      hbs`<ReplicationTableRows @replicationDetails={{replicationDetails}} @clusterMode={{clusterMode}}/>`
    );
    assert.dom('[data-test-table-rows]').exists();
  });

  test('it renders with merkle root, mode, replication set', async function(assert) {
    await render(
      hbs`<ReplicationTableRows @replicationDetails={{replicationDetails}} @clusterMode={{clusterMode}}/>`
    );
    assert.dom('.empty-state').doesNotExist('does not show empty state when data is found');

    assert
      .dom('[data-test-row-value="Merkle root index"]')
      .includesText(REPLICATION_DETAILS.merkleRoot, `shows the correct Merkle Root`);
    assert.dom('[data-test-row-value="Mode"]').includesText(CLUSTER_MODE, `shows the correct Mode`);
    assert
      .dom('[data-test-row-value="Replication set"]')
      .includesText(REPLICATION_DETAILS.clusterId, `shows the correct Cluster ID`);
  });

  test('it renders unknown if values cannot be found', async function(assert) {
    const noAttrs = {
      clusterId: null,
      merkleRoot: null,
    };
    const clusterMode = null;
    this.set('replicationDetails', noAttrs);
    this.set('clusterMode', clusterMode);
    await render(
      hbs`<ReplicationTableRows @replicationDetails={{replicationDetails}} @clusterMode={{clusterMode}}/>`
    );

    assert.dom('[data-test-table-rows]').includesText('unknown');
  });
});
