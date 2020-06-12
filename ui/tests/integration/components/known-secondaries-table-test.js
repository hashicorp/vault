import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import hbs from 'htmlbars-inline-precompile';
const resolver = engineResolverFor('replication');

const SECONDARIES = [
  { node_id: 'secondary-1', api_address: 'https://127.0.0.1:52304', connection_status: 'connected' },
  { node_id: '2nd', connection_status: 'disconnected' },
  { node_id: '_three_', api_address: 'http://127.0.0.1:8202', connection_status: 'connected' },
];

module('Integration | Component | replication known-secondaries-table', function(hooks) {
  setupRenderingTest(hooks, { resolver });

  hooks.beforeEach(function() {
    this.set('replicationAttrs', SECONDARIES);
  });

  test('it renders a table of known secondaries', async function(assert) {
    await render(hbs`<KnownSecondariesTable @replicationAttrs={{replicationAttrs}} />`);

    assert.dom('[data-test-known-secondaries-table]').exists();
  });

  test('it shows the secondary URL and connection_status', async function(assert) {
    await render(hbs`<KnownSecondariesTable @replicationAttrs={{replicationAttrs}} />`);

    SECONDARIES.forEach(secondary => {
      assert.equal(
        this.element.querySelector(`[data-test-secondaries=row-for-${secondary.id}]`).innerHTML.trim(),
        secondary.id,
        'shows a table row and ID for each known secondary'
      );

      assert.equal(
        this.element.querySelector(`[data-test-secondaries=api-address-for-${secondary.id}]`).href,
        secondary.api_address,
        'renders a URL to the secondary UI'
      );

      assert.equal(
        this.element
          .querySelector(`[data-test-secondaries=connection-status-for-${secondary.id}]`)
          .innerHTML.trim(),
        secondary.connection_status,
        'shows the connection status'
      );
    });
  });
});
