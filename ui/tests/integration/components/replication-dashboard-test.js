import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const DR = {
  primaryClusterAddr: 'https://127.0.0.1:8201',
};

module('Integration | Enterprise | Component | replication-dashboard', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('dr', DR);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ReplicationDashboard
    @dr={{dr}}
    />`);

    assert.dom('[data-test-replication-dashboard]').exists();
  });

  test('it renders with primary cluster address when set, and documentation link', async function(assert) {
    await render(hbs`<ReplicationDashboard
    @dr={{dr}}
    />`);

    assert
      .dom('[data-test-primary-cluster-address]')
      .includesText(DR.primaryClusterAddr, `shows the correct primary cluster address value`);

    assert.dom('[data-test-replication-doc-link]').exists();
  });
});
