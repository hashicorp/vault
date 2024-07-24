import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { overrideResponse } from 'vault/tests/helpers/stubs';

import { createSecretsEngine, configUrl } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

// const CONFIGURABLE_SECRET_ENGINES = ['aws'];

module('Integration | Component | SecretEngine::configurable-secret-engine-details', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
  });

  test('it shows error message if no model is passed in', async function (assert) {
    await render(hbs`<SecretEngine::ConfigurableSecretEngineDetails @model={{this.model}}/>`);

    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText('You may not have permissions to configure this engine. Reach out to an admin for help.');
  });

  test('it shows prompt message if no config is returned', async function () {
    // for (const type of CONFIGURABLE_SECRET_ENGINES) {
    const type = 'aws';
    const backend = `test-${type}`;
    this.model = createSecretsEngine(this.store, type, backend);
    this.server.get(configUrl(type, backend), () => {
      return overrideResponse(404);
    });

    await render(hbs`<SecretEngine::ConfigurableSecretEngineDetails @model={{this.model}}/>`);
    // ARG TODO stopped here
    // }
  });
});
