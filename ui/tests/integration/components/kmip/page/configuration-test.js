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
import { KMIP_SELECTORS } from 'vault/tests/helpers/kmip/selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | kmip | Page::Configuration', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:secret-mount-path').update('kmip-test');

    this.config = {
      default_tls_client_key_bits: 256,
      default_tls_client_key_type: 'ec',
      default_tls_client_ttl: 1209600,
      listen_addrs: ['127.0.0.1:5696'],
      server_hostnames: ['localhost'],
      server_ips: ['127.0.0.1', '0.0.0.0'],
      tls_ca_key_bits: 256,
      tls_ca_key_type: 'ec',
      tls_min_version: 'tls12',
      ca_pem: '-----BEGIN CERTIFICATE-----',
    };

    this.renderComponent = () =>
      render(hbs`<Page::Configuration @config={{this.config}} />`, { owner: this.engine });
  });

  test('it should render header and toolbar actions', async function (assert) {
    await this.renderComponent();
    assert.dom(KMIP_SELECTORS.tabs.scope).exists('Page header renders');
    assert.dom(KMIP_SELECTORS.toolbar.download).exists('Download cert button renders');
    assert.dom(KMIP_SELECTORS.toolbar.config).exists('Configure button renders');

    this.config.ca_pem = undefined;
    await this.renderComponent();
    assert
      .dom(KMIP_SELECTORS.toolbar.download)
      .doesNotExist('Download button is hidden when no ca_pem is present');
  });

  test('it should render empty state when config is not provided', async function (assert) {
    this.config = undefined;
    await this.renderComponent();

    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No configuration for this secrets engine', 'renders empty state title');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        "We'll need to configure a few things before getting started.",
        'renders empty state description'
      );
  });

  test('it should render config details', async function (assert) {
    await this.renderComponent();
    assert
      .dom(GENERAL.infoRowValue('Listen addresses'))
      .hasText('127.0.0.1:5696', 'renders listen addresses');
    assert
      .dom(GENERAL.infoRowValue('Default TLS client key bits'))
      .hasText('256', 'renders default tls client key bits');
    assert
      .dom(GENERAL.infoRowValue('Default TLS client key type'))
      .hasText('ec', 'renders default tls client key type');
    assert
      .dom(GENERAL.infoRowValue('Default TLS client TTL'))
      .hasText('14 days', 'renders default tls client ttl');
    assert.dom(GENERAL.infoRowValue('Server hostnames')).hasText('localhost', 'renders server hostnames');
    assert.dom(GENERAL.infoRowValue('Server IPs')).hasText('127.0.0.1,0.0.0.0', 'renders server IPs');
    assert.dom(GENERAL.infoRowValue('TLS CA key bits')).hasText('256', 'renders tls ca key bits');
    assert.dom(GENERAL.infoRowValue('TLS CA key type')).hasText('ec', 'renders tls ca key type');
    assert.dom(GENERAL.infoRowValue('Minimum TLS version')).hasText('tls12', 'renders tls minimum version');
    assert.dom(`${GENERAL.infoRowValue('CA PEM')} ${GENERAL.maskedInput}`).exists('renders ca pem');
  });
});
