import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | raft-storage-overview', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    let model = [
      { address: '127.0.0.1:8200', voter: true },
      { address: '127.0.0.1:8200', voter: true, leader: true },
    ];
    this.set('model', model);
    await render(hbs`<RaftStorageOverview @model={{this.model}} />`);
    assert.dom('[data-raft-row]').exists({ count: 2 });
  });
});
