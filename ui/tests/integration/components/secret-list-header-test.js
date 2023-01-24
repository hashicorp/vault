import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { setupMirage } from 'ember-cli-mirage/test-support';
import mirageToModels from 'vault/tests/helpers/mirage-to-models';

module('Integration | Component | secret-list-header', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('it should render version badge for kv and generic engine types', async function (assert) {
    const backends = supportedSecretBackends();
    assert.expect(backends.length);

    this.server.post('/sys/capabilities-self', () => {});

    for (const type of backends) {
      const data = this.server.create('secret-engine', 2, { type });
      this.model = mirageToModels(data);
      await render(hbs`
        <SecretListHeader
          @model={{this.model}}
        />
        <div id="modal-wormhole"></div>
      `);
      const selector = '[data-test-kv-version-badge]';

      if (['kv', 'generic'].includes(type)) {
        assert
          .dom(selector)
          .hasText(
            `Version ${this.model.version || 1}`,
            `Badge renders with correct version for ${type} engine type`
          );
      } else {
        assert.dom(selector).doesNotExist(`Version badge does not render for ${type} engine type`);
      }
    }
  });

  test('it should render new pki beta button and remain the same for other engines', async function (assert) {
    const backends = supportedSecretBackends();
    const numExpects = backends.length + 1;
    assert.expect(numExpects);

    this.server.post('/sys/capabilities-self', () => {});

    for (const type of backends) {
      const data = this.server.create('secret-engine', 2, { type });
      this.model = mirageToModels(data);
      await render(hbs`
        <SecretListHeader
          @model={{this.model}}
        />
        <div id="modal-wormhole"></div>
      `);
      const oldPkiBetaButtonSelector = '[data-test-old-pki-beta-button]';
      const oldPkiBetaModalSelector = '[data-test-modal-background="New PKI Beta"]';

      if (type === 'pki') {
        assert.dom(oldPkiBetaButtonSelector).hasText('New PKI UI available');
        await click(oldPkiBetaButtonSelector);
        assert.dom(oldPkiBetaModalSelector).exists();
      } else {
        assert
          .dom(oldPkiBetaButtonSelector)
          .doesNotExist(`Version badge does not render for ${type} engine type`);
      }
    }
  });

  test('it should render return to old pki from new pki', async function (assert) {
    const backends = supportedSecretBackends();
    assert.expect(backends.length);

    this.server.post('/sys/capabilities-self', () => {});

    for (const type of backends) {
      const data = this.server.create('secret-engine', 2, { type });
      this.model = mirageToModels(data);
      await render(hbs`
        <SecretListHeader
          @model={{this.model}}
          @isEngine={{true}}
        />
        <div id="modal-wormhole"></div>
      `);
      const newPkiButtonSelector = '[data-test-new-pki-beta-button]';

      if (type === 'pki') {
        assert.dom(newPkiButtonSelector).hasText('Return to old PKI');
      } else {
        assert.dom(newPkiButtonSelector).doesNotExist(`No return to old pki exists`);
      }
    }
  });

  test('it should show the pki modal when New PKI UI available button is clicked', async function (assert) {
    const backends = supportedSecretBackends();
    const numExpects = backends.length + 1;
    assert.expect(numExpects);

    this.server.post('/sys/capabilities-self', () => {});

    for (const type of backends) {
      const data = this.server.create('secret-engine', 2, { type });
      this.model = mirageToModels(data);
      await render(hbs`
        <SecretListHeader
          @model={{this.model}}
        />
        <div id="modal-wormhole"></div>
      `);
      const oldPkiButtonSelector = '[data-test-old-pki-beta-button]';
      const cancelPkiBetaModal = '[data-test-cancel-pki-beta-modal]';

      if (type === 'pki') {
        await click(oldPkiButtonSelector);
        assert.dom('.modal.is-active').exists('Pki beta modal is open');
        await click(cancelPkiBetaModal);
        assert.dom('.modal').exists('Pki beta modal is closed');
      } else {
        assert.dom(oldPkiButtonSelector).doesNotExist(`No return to old pki exists`);
      }
    }
  });
});
