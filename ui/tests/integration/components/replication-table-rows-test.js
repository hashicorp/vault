/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import ReplicationAttributesResource from 'vault/resources/replication-attributes';

const REPLICATION_DETAILS = {
  clusterId: 'b829d963-6835-33eb-a903-57376024b97a',
  merkleRoot: 'c21c8428a0a06135cef6ae25bf8e0267ff1592a6',
};

const CLUSTER_MODE = 'primary';

module('Integration | Component | replication-table-rows', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('replicationDetails', REPLICATION_DETAILS);
    this.set('clusterMode', CLUSTER_MODE);
  });

  test('it renders', async function (assert) {
    await render(
      hbs`<ReplicationTableRows @replicationDetails={{this.replicationDetails}} @clusterMode={{this.clusterMode}}/>`
    );
    assert.dom('[data-test-table-rows]').exists();
  });

  test('it renders with merkle root, mode, replication set', async function (assert) {
    await render(
      hbs`<ReplicationTableRows @replicationDetails={{this.replicationDetails}} @clusterMode={{this.clusterMode}}/>`
    );
    assert
      .dom('[data-test-row-value="Merkle root index"]')
      .includesText(REPLICATION_DETAILS.merkleRoot, `shows the correct Merkle Root`);
    assert.dom('[data-test-row-value="Mode"]').includesText(CLUSTER_MODE, `shows the correct Mode`);
    assert
      .dom('[data-test-row-value="Replication set"]')
      .includesText(REPLICATION_DETAILS.clusterId, `shows the correct Cluster ID`);
  });

  test('it renders unknown if values cannot be found', async function (assert) {
    const noAttrs = {
      clusterId: null,
      merkleRoot: null,
    };
    const clusterMode = null;
    this.set('replicationDetails', noAttrs);
    this.set('clusterMode', clusterMode);
    await render(
      hbs`<ReplicationTableRows @replicationDetails={{this.replicationDetails}} @clusterMode={{this.clusterMode}}/>`
    );

    assert.dom('[data-test-table-rows]').includesText('unknown');
  });

  test('it updates when the replication resource is refreshed in place', async function (assert) {
    const details = new ReplicationAttributesResource({
      cluster_id: 'b829d963-6835-33eb-a903-57376024b97a',
      mode: 'primary',
      merkle_root: 'root-one',
      primary_cluster_addr: 'https://127.0.0.1:8201',
    });
    this.set('replicationDetails', details);
    await render(
      hbs`<ReplicationTableRows @replicationDetails={{this.replicationDetails}} @clusterMode={{this.clusterMode}}/>`
    );
    assert.dom('[data-test-row-value="Merkle root index"]').includesText('root-one');

    details.update({ merkle_root: 'root-two', primary_cluster_addr: 'https://127.0.0.1:9201' });
    await settled();
    assert
      .dom('[data-test-row-value="Merkle root index"]')
      .includesText('root-two', 'the merkle root is updated');
    assert
      .dom('[data-test-row-value="primary_cluster_addr"]')
      .includesText('https://127.0.0.1:9201', 'the primary address is updated');
  });
});
