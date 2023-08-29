/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import SELECTORS from 'vault/tests/helpers/components/dashboard/secrets-engines-card';

module('Integration | Component | dashboard/secrets-engines-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it should hide show all button', async function (assert) {
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kubernetes_f3400dee',
        path: 'kubernetes-test/',
        type: 'kubernetes',
      },
    });

    this.secretsEngines = this.store.peekAll('secret-engine', {});

    await render(hbs`<Dashboard::SecretsEnginesCard @secretsEngines={{this.secretsEngines}} />`);

    assert.dom('[data-test-secrets-engines-card-show-all]').doesNotExist();
  });

  module('secrets engines with 5 or more enabled', function (hooks) {
    hooks.beforeEach(function () {
      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'kubernetes_f3400dee',
          path: 'kubernetes-test/',
          type: 'kubernetes',
        },
      });
      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'pki_i1234dd',
          path: 'pki-test/',
          type: 'pki',
        },
      });
      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'secrets_j2350ii',
          path: 'secrets-test/',
          type: 'kv',
        },
      });
      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'nomad_123hh',
          path: 'nomad/',
          type: 'nomad',
        },
      });
      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'pki_f3400dee',
          path: 'pki-0-test/',
          type: 'pki',
        },
      });
      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'pki_i1234dd',
          path: 'pki-1-test/',
          description: 'pki-1-path-description',
          type: 'pki',
        },
      });
      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'secrets_j2350ii',
          path: 'secrets-1-test/',
          type: 'kv',
        },
      });

      this.secretsEngines = this.store.peekAll('secret-engine', {});

      this.renderComponent = () => {
        return render(hbs`<Dashboard::SecretsEnginesCard @secretsEngines={{this.secretsEngines}} />`);
      };
    });

    test('it should display only five secrets engines and show help text for more than 5 engines', async function (assert) {
      await this.renderComponent();
      assert.dom(SELECTORS.cardTitle).hasText('Secrets engines');
      assert.dom(SELECTORS.secretEnginesTableRows).exists({ count: 5 });
      assert.dom('[data-test-secrets-engine-total-help-text]').exists();
      assert
        .dom('[data-test-secrets-engine-total-help-text]')
        .hasText(
          `Showing 5 out of ${this.secretsEngines.length} secret engines. Navigate to details to view more.`
        );
    });

    test('it should display the secrets engines accessor and path', async function (assert) {
      await this.renderComponent();
      assert.dom(SELECTORS.cardTitle).hasText('Secrets engines');
      assert.dom(SELECTORS.secretEnginesTableRows).exists({ count: 5 });

      this.secretsEngines.slice(0, 5).forEach((engine) => {
        assert.dom(SELECTORS.getSecretEngineAccessor(engine.id)).hasText(engine.accessor);
        if (engine.description) {
          assert.dom(SELECTORS.getSecretEngineDescription(engine.id)).hasText(engine.description);
        } else {
          assert.dom(SELECTORS.getSecretEngineDescription(engine.id)).doesNotExist(engine.description);
        }
      });
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      await this.renderComponent();
      assert.dom('[data-test-secrets-engines-row="nomad"] [data-test-view]').doesNotExist();
      assert.dom('[data-test-icon="nomad"]').hasClass('has-text-grey');
    });
  });
});
