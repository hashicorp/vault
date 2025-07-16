/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import {
  expectedConfigKeys,
  expectedValueOfConfigKeys,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';
import engineDisplayData from 'vault/helpers/engines-display-data';

module('Integration | Component | SecretEngine::ConfigurationDetails', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.configs = {
      aws: {
        region: 'us-west-2',
        accessKey: '123-key',
        iamEndpoint: 'iam-endpoint',
        stsEndpoint: 'sts-endpoint',
        maxRetries: 1,
      },
      azure: {
        clientSecret: 'client-secret',
        subscriptionId: 'subscription-id',
        tenantId: 'tenant-id',
        clientId: 'client-id',
        rootPasswordTtl: '1800000s',
        environment: 'AZUREPUBLICCLOUD',
      },
      gcp: {
        credentials: '{"some-key":"some-value"}',
        ttl: '100s',
        maxTtl: '101s',
      },
      ssh: {
        publicKey: 'public-key',
        generateSigningKey: true,
      },
    };
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

  for (const type of ALL_ENGINES.filter((engine) => engine.isConfigurable ?? false).map(
    (engine) => engine.type
  )) {
    test(`${type}: it shows config details if configModel(s) are passed in`, async function (assert) {
      this.config = this.configs[type];
      this.typeDisplay = engineDisplayData(type).displayName;

      await render(
        hbs`<SecretEngine::ConfigurationDetails @config={{this.config}} @typeDisplay={{this.typeDisplay}}/>`
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
