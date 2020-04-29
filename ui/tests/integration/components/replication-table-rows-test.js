import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const DATA = {
  clusterId: 'b829d963-6835-33eb-a903-57376024b97a',
  mode: 'primary',
  merkleRoot: 'c21c8428a0a06135cef6ae25bf8e0267ff1592a6',
};

module('Integration | Component | replication-table-rows', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('data', DATA);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationTableRows @data={{data}}/>`);

    assert.dom('[data-test-table-rows]').exists();
  });

  test('it renders with merkle root, mode, replication set', async function(assert) {
    await render(hbs`<ReplicationTableRows @data={{data}}/>`);

    assert.dom('.empty-state').doesNotExist('does not show empty state when data is found');

    assert
      .dom('[data-test-row-value="Merkle root index"]')
      .includesText(DATA.merkleRoot, `shows the correct Merkle Root`);
    assert.dom('[data-test-row-value="Mode"]').includesText(DATA.mode, `shows the correct Merkle Root`);
    assert
      .dom('[data-test-row-value="Replication set"]')
      .includesText(DATA.clusterId, `shows the correct Merkle Root`);
  });

  test('it renders unknown if values cannot be found', async function(assert) {
    const noAttrs = {
      clusterId: null,
      mode: null,
      merkleRoot: null,
    };
    this.set('data', noAttrs);
    await render(hbs`<ReplicationTableRows @data={{data}}/>`);

    assert.dom('[data-test-table-rows]').includesText('unknown');
  });
});
