/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';
import {
  createSecretsEngine,
  createConfig,
  configUrl,
  expectedConfigKeys,
  expectedValueOfConfigKeys,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

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
      .hasText(
        'We are unable to access the mount information for this engine. Ask you administrator if you think you should have access to this secret engine.'
      );
  });

  test('it shows prompt message if no config is returned', async function (assert) {
    assert.expect(CONFIGURABLE_SECRET_ENGINES.length * 2);
    for (const type of CONFIGURABLE_SECRET_ENGINES) {
      const title = type.toUpperCase();
      const backend = `test-404-${type}`;
      this.model = createSecretsEngine(this.store, type, backend);
      this.server.get(configUrl(type, backend), () => {
        return overrideResponse(404);
      });

      await render(hbs`<SecretEngine::ConfigurableSecretEngineDetails @model={{this.model}}/>`);
      assert.dom(GENERAL.emptyStateTitle).hasText(`${title} not configured`);
      assert.dom(GENERAL.emptyStateMessage).hasText(`Get started by configuring your ${title} engine.`);
    }
  });

  test('it shows API error', async function (assert) {
    assert.expect(CONFIGURABLE_SECRET_ENGINES.length * 2);
    for (const type of CONFIGURABLE_SECRET_ENGINES) {
      const backend = `test-400-${type}`;
      this.model = createSecretsEngine(this.store, type, backend);
      this.server.get(configUrl(type, backend), () => {
        return overrideResponse(400, { errors: ['error'] });
      });

      await render(hbs`<SecretEngine::ConfigurableSecretEngineDetails @model={{this.model}}/>`);
      assert.dom(GENERAL.emptyStateTitle).hasText(`Something went wrong`);
      assert.dom(GENERAL.emptyStateMessage).hasText(`error`);
    }
  });

  test('it shows config details if config data is returned', async function (assert) {
    assert.expect(14);
    for (const type of CONFIGURABLE_SECRET_ENGINES) {
      const backend = `test-${type}`;
      this.model = createSecretsEngine(this.store, type, backend);
      createConfig(this.store, backend, type);
      this.server.get(configUrl(type, backend), () => {
        return overrideResponse(200);
      });

      await render(hbs`<SecretEngine::ConfigurableSecretEngineDetails @model={{this.model}}/>`);
      for (const key of expectedConfigKeys(type)) {
        assert.dom(GENERAL.infoRowLabel(key)).exists(`${key} on the ${type} config details exists.`);
        const responseKeyAndValue = expectedValueOfConfigKeys(type, key);
        assert
          .dom(GENERAL.infoRowValue(key))
          .hasText(responseKeyAndValue, `${key} value for the ${type} config details exists.`);
      }
    }
  });
});
