import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MODEL = {
  replicationMode: 'dr',
};

module('Integration | Enterprise | Component | replication-page', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('model', MODEL);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationPage @model={{model}} />`);

    assert.dom('[data-test-replication-page]').exists();
  });
});
