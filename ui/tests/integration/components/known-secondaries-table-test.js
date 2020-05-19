import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import hbs from 'htmlbars-inline-precompile';
const resolver = engineResolverFor('replication');

const REPLICATION_ATTRS = {
  secondaries: [
    { id: 'secondary-1', api_address: 'https://stuff.com/', connection_status: 'connected' },
    { id: '2nd', api_address: 'https://10.0.0.2:1234/', connection_status: 'disconnected' },
    { id: '_three_', api_address: 'https://10.0.0.2:1000/', connection_status: 'connected' },
  ],
};

module('Integration | Component | replication known-secondaries-table', function(hooks) {
  setupRenderingTest(hooks, { resolver });

  hooks.beforeEach(function() {
    this.set('replicationAttrs', REPLICATION_ATTRS);
  });

  test('it renders a table of known secondaries', async function(assert) {
    await render(hbs`<KnownSecondariesTable @replicationAttrs={{replicationAttrs}} />`);

    assert.dom('[data-test-known-secondaries-table]').exists();
  });

  test('it shows the secondary URL and connection_status', async function(assert) {
    await render(hbs`<KnownSecondariesTable @replicationAttrs={{replicationAttrs}} />`);

    REPLICATION_ATTRS.secondaries.forEach(secondary => {
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
