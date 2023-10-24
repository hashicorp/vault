import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | open-api-explorer | swagger-ui', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'open-api-explorer');
  setupMirage(hooks);
  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it renders', async function (assert) {
    assert.expect(2);
    const openApiResponse = this.server.create('open-api-explorer');
    this.server.get('sys/internal/specs/openapi', () => {
      assert.ok(true, 'request made to correct endpoint.');
      return openApiResponse;
    });

    await render(hbs`<SwaggerUi/>`, {
      owner: this.engine,
    });

    assert.dom('[data-test-swagger-ui]').exists('renders component');
  });
});
