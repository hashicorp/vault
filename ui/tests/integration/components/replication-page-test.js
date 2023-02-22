import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const MODEL = {
  replicationMode: 'dr',
  replicationAttrs: {
    mode: 'secondary',
    clusterId: '12ab',
    replicationDisabled: false,
  },
};

module('Integration | Component | replication-page', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('model', MODEL);
  });

  test('it renders', async function (assert) {
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-replication-page]').exists();
    assert.dom('[data-test-layout-loading]').doesNotExist();
  });

  test('it renders loader when either clusterId is unknown or mode is bootstrapping', async function (assert) {
    this.set('model.replicationAttrs.clusterId', '');
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-layout-loading]').exists();

    this.set('model.replicationAttrs.clusterId', '123456');
    this.set('model.replicationAttrs.mode', 'bootstrapping');
    await render(hbs`<ReplicationPage @model={{this.model}} />`);
    assert.dom('[data-test-layout-loading]').exists();
  });
});
