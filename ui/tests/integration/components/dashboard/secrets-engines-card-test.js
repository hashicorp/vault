import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const SELECTORS = {
  cardTitle: '[data-test-dashboard-secrets-engines-header] h3',
  secretEnginesTableRows: '[data-test-dashboard-secrets-engines-table] tr',
};

module('Integration | Component | dashboard/secrets-engines-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
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

  test('it should display only five secrets engines', async function (assert) {
    await this.renderComponent();
    assert.dom(SELECTORS.cardTitle).hasText('Secrets Engines');
    assert.dom(SELECTORS.secretEnginesTableRows).exists({ count: 5 });
  });

  test('it should display the secrets engines accessor and path', async function (assert) {
    await this.renderComponent();
    assert.dom(SELECTORS.cardTitle).hasText('Secrets Engines');
    assert.dom(SELECTORS.secretEnginesTableRows).exists({ count: 5 });

    this.secretsEngines.slice(0, 5).forEach((engine) => {
      assert
        .dom(`[data-test-secrets-engines-row=${engine.id}] [data-test-accessor]`)
        .hasText(engine.accessor);
      if (engine.description) {
        assert
          .dom(`[data-test-secrets-engines-row=${engine.id}] [data-test-description]`)
          .hasText(engine.description);
      } else {
        assert
          .dom(`[data-test-secrets-engines-row=${engine.id}] [data-test-description]`)
          .doesNotExist(engine.description);
      }
    });
  });
});
