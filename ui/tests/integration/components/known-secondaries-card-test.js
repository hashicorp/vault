import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import hbs from 'htmlbars-inline-precompile';
const resolver = engineResolverFor('replication');

const CLUSTER = {
  canAddSecondary: true,
  replicationMode: 'dr',
};

const REPLICATION_ATTRS = {
  secondaries: [
    { id: 'secondary-1', api_address: 'https://stuff.com/', connection_status: 'connected' },
    { id: '2nd', api_address: 'https://10.0.0.2:1234/', connection_status: 'disconnected' },
    { id: '_three_', api_address: 'https://10.0.0.2:1000/', connection_status: 'connected' },
  ],
};

module('Integration | Component | replication known-secondaries-card', function(hooks) {
  setupRenderingTest(hooks, { resolver });

  hooks.beforeEach(function() {
    this.set('cluster', CLUSTER);
    this.set('replicationAttrs', REPLICATION_ATTRS);
  });

  test('it renders with a table of known secondaries', async function(assert) {
    await render(hbs`<KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />`);

    assert
      .dom('[data-test-known-secondaries-table]')
      .exists('shows known secondaries table when there are known secondaries');
    assert.dom('[data-test-manage-link]').exists('shows manage link');
  });

  // TODO: this test can be deleted once the 'secondaries' array replaces 'knownSecondaries'
  // in the backend
  test('it renders a known secondaries table even if api address and connection_status are missing', async function(assert) {
    const missingSecondariesInfo = {
      knownSecondaries: ['secondary-1', '2nd', '_three_'],
    };
    this.set('replicationAttrs', missingSecondariesInfo);

    await render(hbs`<KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />`);

    assert
      .dom('[data-test-known-secondaries-table]')
      .exists('shows known secondaries table when there are known secondaries');
  });

  test('it renders an empty state if there are no known secondaries', async function(assert) {
    const noSecondaries = {
      secondaries: [],
    };
    this.set('replicationAttrs', noSecondaries);
    await render(hbs`<KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />`);

    assert
      .dom('[data-test-known-secondaries-table]')
      .doesNotExist('does not show the known secondaries table');
    assert
      .dom('.empty-state')
      .includesText('No known dr secondary clusters', 'has a message with the replication mode');
  });

  test('it renders an Add secondary link if user has capabilites', async function(assert) {
    await render(hbs`<KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />`);

    assert.dom('.add-secondaries').exists();
  });

  test('it does not render an Add secondary link if user does not have capabilites', async function(assert) {
    const noCapabilities = {
      canAddSecondary: false,
    };
    this.set('cluster', noCapabilities);
    await render(hbs`<KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />`);

    assert.dom('.add-secondaries').doesNotExist();
  });
});
