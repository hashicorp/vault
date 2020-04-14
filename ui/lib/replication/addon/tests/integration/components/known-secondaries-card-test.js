import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import hbs from 'htmlbars-inline-precompile';

const CLUSTER = {
  canAddSecondary: true,
  replicationMode: 'dr',
};

const REPLICATION_ATTRS = {
  knownSecondaries: [5, 4343, 3432],
};

module('Integration | Component | replication known-secondaries-card', function(hooks) {
  setupRenderingTest(hooks, { resolver: engineResolverFor('replication') });

  hooks.beforeEach(function() {
    this.set('cluster', CLUSTER);
    this.set('replicationAttrs', REPLICATION_ATTRS);
  });

  test('it renders with secondary ids', async function(assert) {
    await render(hbs`<KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />`);

    assert.dom('.secondaries').exists();
  });

  test('it renders an empty state if there are no known secondaries', async function(assert) {
    const noSecondaries = {
      knownSecondaries: null,
    };
    await render(hbs`<KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />`);

    assert.dom('.empty-state').exists();
  });
});
