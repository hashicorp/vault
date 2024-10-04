/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | kubernetes | Page::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kubernetes');
  setupMirage(hooks);

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
    this.backend = this.store.peekRecord('secret-engine', 'kubernetes-test');
    this.config = null;

    this.setConfig = (disableLocal) => {
      const data = this.server.create(
        'kubernetes-config',
        !disableLocal ? { disable_local_ca_jwt: false } : null
      );
      this.store.pushPayload('kubernetes/config', {
        modelName: 'kubernetes/config',
        backend: 'kubernetes-test',
        ...data,
      });
      this.config = this.store.peekRecord('kubernetes/config', 'kubernetes-test');
    };

    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend.id },
    ];

    this.renderComponent = () => {
      return render(
        hbs`<Page::Configuration @backend={{this.backend}} @config={{this.config}} @breadcrumbs={{this.breadcrumbs}} />`,
        {
          owner: this.engine,
        }
      );
    };
  });

  test('it should render tab page header, config cta and mount config', async function (assert) {
    await this.renderComponent();
    assert.dom('.title svg').hasClass('flight-icon-kubernetes-color', 'Kubernetes icon renders in title');
    assert.dom('.title').hasText('kubernetes-test', 'Mount path renders in title');
    assert
      .dom('[data-test-toolbar-config-action]')
      .hasText('Configure Kubernetes', 'Toolbar action has correct text');
    assert.dom('[data-test-config-cta]').exists('Config cta renders');
    assert.dom('[data-test-mount-config]').exists('Mount config renders');
  });

  test('it should render message for inferred configuration', async function (assert) {
    this.setConfig(false);
    await this.renderComponent();
    assert
      .dom('[data-test-inferred-message] svg')
      .hasClass('flight-icon-check-circle-fill', 'Inferred message icon renders');
    const message =
      'These details were successfully inferred from Vault’s kubernetes environment and were not explicity set in this config.';
    assert.dom('[data-test-inferred-message]').hasText(message, 'Inferred message renders');
    assert
      .dom('[data-test-toolbar-config-action]')
      .hasText('Edit configuration', 'Toolbar action has correct text');
  });

  test('it should render host and certificate info', async function (assert) {
    this.setConfig(true);
    await this.renderComponent();
    assert.dom('[data-test-row-label="Kubernetes host"]').exists('Kubernetes host label renders');
    assert
      .dom('[data-test-row-value="Kubernetes host"]')
      .hasText(this.config.kubernetesHost, 'Kubernetes host value renders');

    assert.dom('[data-test-row-label="Certificate"]').exists('Certificate label renders');
    assert.dom('[data-test-certificate-card]').exists('Certificate card component renders');
    assert
      .dom('[data-test-certificate-icon]')
      .hasClass('flight-icon-certificate', 'Certificate icon renders');
    assert
      .dom('[data-test-certificate-card] [data-test-copy-button]')
      .exists('Certificate copy button renders');
    assert.dom('[data-test-certificate-label]').hasText('PEM Format', 'Certificate label renders');
    assert
      .dom('[data-test-certificate-value]')
      .hasText(this.config.kubernetesCaCert, 'Certificate value renders');
  });
});
