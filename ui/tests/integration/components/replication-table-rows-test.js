import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const REPLICATION_ATTRS = {
  clusterId: 'b829d963-6835-33eb-a903-57376024b97a',
  mode: 'primary',
  merkleRoot: 'c21c8428a0a06135cef6ae25bf8e0267ff1592a6',
};

module('Integration | Component | replication-table-rows', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('data', REPLICATION_ATTRS);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationTableRows @data={{replicationAttrs}}/>`);

    assert.dom('.replication-table-rows').exists();
  });

  // it renders with merkle root, mode, replication set
  // test('', async function(assert) {
  //   await render(hbs`<ReplicationTableRows @data={{replicationAttrs}}/>`);

  //   Object.keys(REPLICATION_ATTRS).forEach(attr => {
  //     let expected = REPLICATION_ATTRS[attr];
  //     let found = this.element.querySelector(`[data-test-attr-${expected}]`);

  //     debugger;

  //     assert.equal(found.textContent.trim(), expected);
  //   });
  // });

  // it renders unknown if values cannot be found
});
