/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';

module('Integration | Component | kubernetes | Page::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.secretsEngine = new SecretsEngineResource({
      accessor: 'kubernetes_f3400dee',
      path: 'kubernetes-test/',
      type: 'kubernetes',
    });

    this.config = null;

    this.setConfig = (disableLocal) => {
      this.config = this.server.create(
        'kubernetes-config',
        !disableLocal ? { disable_local_ca_jwt: false } : null
      );
    };

    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretsEngine.id },
    ];

    this.renderComponent = () => {
      return render(
        hbs`<Page::Configuration @config={{this.config}} @secretsEngine={{this.secretsEngine}} @breadcrumbs={{this.breadcrumbs}} />`,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it should render tab page header', async function (assert) {
    await this.renderComponent();
    assert
      .dom(GENERAL.icon('kubernetes-color'))
      .hasClass('hds-icon-kubernetes-color', 'Kubernetes icon renders in title');
    assert
      .dom(GENERAL.hdsPageHeaderTitle)
      .hasText('kubernetes-test configuration', 'Mount path renders in title');
    assert.dom(SES.configure).doesNotExist('Toolbar action does not render when engine is not configured');
  });

  test('it should render message for inferred configuration', async function (assert) {
    this.setConfig(false);
    await this.renderComponent();
    assert
      .dom('[data-test-inferred-message] svg')
      .hasClass('hds-icon-check-circle-fill', 'Inferred message icon renders');
    const message =
      'These details were successfully inferred from Vault’s kubernetes environment and were not explicity set in this config.';
    assert.dom('[data-test-inferred-message]').hasText(message, 'Inferred message renders');
    assert.dom(SES.configure).exists().hasText('Edit configuration', 'Toolbar action has correct text');
  });

  test('it should render host and certificate info', async function (assert) {
    this.setConfig(true);
    await this.renderComponent();
    assert.dom('[data-test-row-label="Kubernetes host"]').exists('Kubernetes host label renders');
    assert
      .dom('[data-test-row-value="Kubernetes host"]')
      .hasText(this.config.kubernetes_host, 'Kubernetes host value renders');

    assert.dom('[data-test-row-label="Certificate"]').exists('Certificate label renders');
    assert.dom('[data-test-certificate-card]').exists('Certificate card component renders');
    assert.dom('[data-test-certificate-icon]').hasClass('hds-icon-certificate', 'Certificate icon renders');
    assert.dom(GENERAL.copyButton).exists('Certificate copy button renders');
    assert.dom('[data-test-certificate-label]').hasText('PEM Format', 'Certificate label renders');
    assert
      .dom('[data-test-certificate-value]')
      .hasText(this.config.kubernetes_ca_cert, 'Certificate value renders');
  });
});
