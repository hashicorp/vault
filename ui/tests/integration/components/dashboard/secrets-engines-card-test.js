/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { DASHBOARD } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

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

    // verify overflow style exists on secret engine text
    assert
      .dom(SES.secretPath('kubernetes-test/'))
      .hasClass('overflow-wrap', 'secret engine name has overflow class ');

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
      assert.dom(DASHBOARD.cardHeader('Secrets engines')).hasText('Secrets engines');
      assert.dom(DASHBOARD.tableRow('Secrets engines')).exists({ count: 5 });
      assert.dom('[data-test-secrets-engine-total-help-text]').exists();
      assert
        .dom('[data-test-secrets-engine-total-help-text]')
        .hasText(
          `Showing 5 out of ${this.secretsEngines.length} secret engines. Navigate to details to view more.`
        );
    });

    test('it should display the secrets engines accessor and path', async function (assert) {
      await this.renderComponent();
      assert.dom(DASHBOARD.cardHeader('Secrets engines')).hasText('Secrets engines');
      assert.dom(DASHBOARD.tableRow('Secrets engines')).exists({ count: 5 });

      this.secretsEngines.slice(0, 5).forEach((engine) => {
        assert.dom(DASHBOARD.secretsEnginesCard.secretEngineAccessorRow(engine.id)).hasText(engine.accessor);
        if (engine.description) {
          assert
            .dom(DASHBOARD.secretsEnginesCard.secretEngineDescription(engine.id))
            .hasText(engine.description);
        } else {
          assert
            .dom(DASHBOARD.secretsEnginesCard.secretEngineDescription(engine.id))
            .doesNotExist(engine.description);
        }
      });
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      await this.renderComponent();
      assert
        .dom(SES.secretPath('secrets-test/'))
        .hasClass('has-text-black', 'does not apply disabled class to supported secret engine');
      assert
        .dom(SES.secretPath('nomad/'))
        .hasClass('has-text-grey', 'nomad is not a supported secret engine and has disabled class');
    });
  });

  module('favorites functionality', function (hooks) {
    hooks.beforeEach(function () {
      // Clear localStorage before each test
      window.localStorage.clear();

      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'kv_123',
          path: 'secret/',
          type: 'kv',
          id: 'kv_123',
        },
      });

      this.store.pushPayload('secret-engine', {
        modelName: 'secret-engine',
        data: {
          accessor: 'pki_456',
          path: 'pki/',
          type: 'pki',
          id: 'pki_456',
        },
      });

      this.secretsEngines = this.store.peekAll('secret-engine', {});
    });

    hooks.afterEach(function () {
      // Clean up localStorage after each test
      window.localStorage.clear();
    });

    test('it shows star icon for favorite engines', async function (assert) {
      // Pre-populate localStorage with favorites
      localStorage.setItem('vault-favorite-engines', JSON.stringify(['kv_123']));

      await render(hbs`<Dashboard::SecretsEnginesCard @secretsEngines={{this.secretsEngines}} />`);

      // Check that favorite engine shows star icon
      assert
        .dom('[data-test-secrets-engines-row="kv_123"] .hds-icon-star-fill')
        .exists('favorite engine shows star icon');

      // Check that non-favorite engine does not show star icon
      assert
        .dom('[data-test-secrets-engines-row="pki_456"] .hds-icon-star-fill')
        .doesNotExist('non-favorite engine does not show star icon');
    });

    test('it sorts favorite engines to the top', async function (assert) {
      // Pre-populate localStorage with pki_456 as favorite
      localStorage.setItem('vault-favorite-engines', JSON.stringify(['pki_456']));

      await render(hbs`<Dashboard::SecretsEnginesCard @secretsEngines={{this.secretsEngines}} />`);

      // Get all table rows
      const rows = this.element.querySelectorAll('[data-test-dashboard-table="Secrets engines"] tbody tr');

      // First row should be the favorite (pki_456)
      assert
        .dom(rows[0])
        .hasAttribute('data-test-secrets-engines-row', 'pki_456', 'favorite engine appears first');

      // Second row should be the non-favorite (kv_123)
      assert
        .dom(rows[1])
        .hasAttribute('data-test-secrets-engines-row', 'kv_123', 'non-favorite engine appears second');
    });

    test('it applies favorite styling to favorite engines', async function (assert) {
      // Pre-populate localStorage with favorites
      localStorage.setItem('vault-favorite-engines', JSON.stringify(['kv_123']));

      await render(hbs`<Dashboard::SecretsEnginesCard @secretsEngines={{this.secretsEngines}} />`);

      // Check that favorite row has appropriate class
      assert
        .dom('[data-test-secrets-engines-row="kv_123"]')
        .hasClass('is-favorite-row', 'favorite engine row has favorite class');

      // Check that non-favorite row does not have the class
      assert
        .dom('[data-test-secrets-engines-row="pki_456"]')
        .doesNotHaveClass('is-favorite-row', 'non-favorite engine row does not have favorite class');
    });

    test('it loads favorites from localStorage on component initialization', async function (assert) {
      // Pre-populate localStorage
      localStorage.setItem('vault-favorite-engines', JSON.stringify(['kv_123', 'pki_456']));

      await render(hbs`<Dashboard::SecretsEnginesCard @secretsEngines={{this.secretsEngines}} />`);

      // Both engines should show as favorites
      assert
        .dom('[data-test-secrets-engines-row="kv_123"] .hds-icon-star-fill')
        .exists('first favorite engine shows star icon');
      assert
        .dom('[data-test-secrets-engines-row="pki_456"] .hds-icon-star-fill')
        .exists('second favorite engine shows star icon');
    });

    test('it handles corrupted localStorage gracefully', async function (assert) {
      // Set corrupted data in localStorage
      localStorage.setItem('vault-favorite-engines', 'invalid-json');

      await render(hbs`<Dashboard::SecretsEnginesCard @secretsEngines={{this.secretsEngines}} />`);

      // Component should still render without errors
      assert
        .dom('[data-test-dashboard-card-header="Secrets engines"]')
        .exists('component renders successfully with corrupted localStorage');

      // No engines should show as favorites
      assert
        .dom('.hds-icon-star-fill')
        .doesNotExist('no engines show as favorites with corrupted localStorage');
    });
  });
});
