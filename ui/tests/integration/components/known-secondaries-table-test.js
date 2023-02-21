/* eslint qunit/no-conditional-assertions: "warn" */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';

const SECONDARIES = [
  { node_id: 'secondary-1', api_address: 'https://127.0.0.1:52304', connection_status: 'connected' },
  { node_id: '2nd', connection_status: 'disconnected' },
  { node_id: '_three_', api_address: 'http://127.0.0.1:8202', connection_status: 'connected' },
];

module('Integration | Component | replication known-secondaries-table', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'replication');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    this.set('secondaries', SECONDARIES);
  });

  test('it renders a table of known secondaries', async function (assert) {
    await render(hbs`<KnownSecondariesTable @secondaries={{this.secondaries}} />`, this.context);

    assert.dom('[data-test-known-secondaries-table]').exists();
  });

  test('it shows the secondary URL and connection_status', async function (assert) {
    assert.expect(9);
    await render(hbs`<KnownSecondariesTable @secondaries={{this.secondaries}} />`, this.context);

    SECONDARIES.forEach((secondary) => {
      assert.strictEqual(
        this.element.querySelector(`[data-test-secondaries=row-for-${secondary.node_id}]`).innerHTML.trim(),
        secondary.node_id,
        'shows a table row and ID for each known secondary'
      );

      if (secondary.api_address) {
        const expectedUrl = `${secondary.api_address}/ui/`;

        assert.strictEqual(
          this.element.querySelector(`[data-test-secondaries=api-address-for-${secondary.node_id}]`).href,
          expectedUrl,
          'renders a URL to the secondary UI'
        );
      } else {
        assert.notOk(
          this.element.querySelector(`[data-test-secondaries=api-address-for-${secondary.node_id}]`)
        );
      }

      assert.strictEqual(
        this.element
          .querySelector(`[data-test-secondaries=connection-status-for-${secondary.node_id}]`)
          .innerHTML.trim(),
        secondary.connection_status,
        'shows the connection status'
      );
    });
  });
});
