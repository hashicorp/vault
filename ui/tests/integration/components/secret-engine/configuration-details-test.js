/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';
import {
  createConfig,
  expectedConfigKeys,
  expectedValueOfConfigKeys,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

const allEnginesArray = allEngines(); // saving as const so we don't invoke the method multiple times via the for loop

module('Integration | Component | SecretEngine/ConfigurationDetails', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.configModels = [];
  });

  test('it shows prompt message if no config models are passed in', async function (assert) {
    assert.expect(2);
    await render(hbs`
      <SecretEngine::ConfigurationDetails @typeDisplay="Display Name" />
    `);
    assert.dom(GENERAL.emptyStateTitle).hasText(`Display Name not configured`);
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(`Get started by configuring your Display Name secrets engine.`);
  });

  for (const type of CONFIGURABLE_SECRET_ENGINES) {
    test(`${type}: it shows config details if configModel(s) are passed in`, async function (assert) {
      const backend = `test-${type}`;
      this.configModels = createConfig(this.store, backend, type);
      this.typeDisplay = allEnginesArray.find((engine) => engine.type === type).displayName;

      await render(
        hbs`<SecretEngine::ConfigurationDetails @configModels={{array this.configModels}} @typeDisplay={{this.typeDisplay}}/>`
      );

      for (const key of expectedConfigKeys(type)) {
        if (
          key === 'Secret key' ||
          key === 'Client secret' ||
          key === 'Private key' ||
          key === 'Credentials'
        ) {
          // these keys are not returned by the API and should not show on the details page
          assert
            .dom(GENERAL.infoRowLabel(key))
            .doesNotExist(`${key} on the ${type} config details does NOT exists.`);
        } else {
          // check the label appears
          assert.dom(GENERAL.infoRowLabel(key)).exists(`${key} on the ${type} config details exists.`);
          const responseKeyAndValue = expectedValueOfConfigKeys(type, key);
          // check the value appears
          assert
            .dom(GENERAL.infoRowValue(key))
            .hasText(responseKeyAndValue, `${key} value for the ${type} config details exists.`);
          // make sure the values that should be masked are masked, and others are not.
          if (key === 'Public Key') {
            assert.dom(GENERAL.infoRowValue(key)).hasClass('masked-input', `${key} is masked`);
          } else {
            assert.dom(GENERAL.infoRowValue(key)).doesNotHaveClass('masked-input', `${key} is not masked`);
          }
        }
      }
    });
  }
});
