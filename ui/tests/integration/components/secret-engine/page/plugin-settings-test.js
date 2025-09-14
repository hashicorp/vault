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

module('Integration | Component | SecretEngine::Page::PluginSettings', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.models = {
      keymgmt: {
        secretsEngine: {
          type: 'keymgmt',
        },
        config: null,
      },
      aws: {
        secretsEngine: {
          type: 'aws',
        },
        config: {
          region: 'us-west-2',
          access_key: '123-key',
          iam_endpoint: 'iam-endpoint',
          sts_endpoint: 'sts-endpoint',
          max_retries: 1,
        },
      },
      azure: {
        secretsEngine: {
          type: 'azure',
        },
        config: {
          client_secret: 'client-secret',
          subscription_id: 'subscription-id',
          tenant_id: 'tenant-id',
          client_id: 'client-id',
          root_password_ttl: '1800000s',
          environment: 'AZUREPUBLICCLOUD',
        },
      },
      gcp: {
        secretsEngine: {
          type: 'gcp',
        },
        config: {
          credentials: '{"some-key":"some-value"}',
          ttl: '100s',
          max_ttl: '101s',
        },
      },
      ssh: {
        secretsEngine: {
          type: 'ssh',
        },
        config: {
          public_key: 'public-key',
          generate_signing_key: true,
        },
      },
    };
  });

  test('it shows empty state when the engine is not configurable', async function (assert) {
    assert.expect(2);
    this.model = this.models['keymgmt'];
    await render(hbs`
      <SecretEngine::Page::PluginSettings @model={{this.model}} />
    `);
    assert.dom(GENERAL.emptyStateTitle).hasText(`No configuration details available`);
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        `Key Management does not have any plugin specific configuration. All configurable parameters for this engine are under 'General Settings'.`
      );
  });

  for (const type of ALL_ENGINES.filter((engine) => engine.isConfigurable ?? false).map(
    (engine) => engine.type
  )) {
    test(`${type}: it shows config details if configModel(s) are passed in`, async function (assert) {
      this.model = this.models[type];

      await render(hbs`<SecretEngine::Page::PluginSettings @model={{this.model}} />`);

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
