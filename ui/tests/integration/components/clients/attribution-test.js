import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/attribution', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('totalClientsData', [{ label: 'first', total: 5, entity_clients: 3, non_entity_clients: 2 }]);
    this.set('chartLegend', [{ label: 'first', key: 'one' }]);
  });

  test('it does not render export button if no data', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <Clients::Attribution @chartLegend={{chartLegend}} />
    `);
    assert.dom('[data-test-export-attribution-data]').doesNotExist('Export button not rendered');
    // Shows "problem gathering data" messages
    await this.pauseTest();
  });
});
