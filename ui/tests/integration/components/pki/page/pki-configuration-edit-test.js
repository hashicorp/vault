/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-configuration-edit';
import sinon from 'sinon';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | page/pki-configuration-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    this.store = this.owner.lookup('service:store');
    this.cancelSpy = sinon.spy();
    this.backend = 'pki-engine';
    this.store.pushPayload('pki/crl', {
      modelName: 'pki/crl',
      id: this.backend,
      auto_rebuild: false,
      auto_rebuild_grace_period: '12h',
      delta_rebuild_interval: '3d',
      disable: false,
      enable_delta: false,
      expiry: '24h',
      ocsp_disable: false,
      ocsp_expiry: '18m',
    });
    this.store.pushPayload('pki/urls', {
      modelName: 'pki/urls',
      id: this.backend,
      issuing_certificates: ['hashicorp.com'],
      crl_distribution_points: ['some-crl-distribution.com'],
      ocsp_servers: ['ocsp-stuff.com'],
    });
    this.urls = this.store.peekRecord('pki/urls', this.backend);
    this.crl = this.store.peekRecord('pki/crl', this.backend);
  });

  test('it renders with config data and updates config', async function (assert) {
    assert.expect(27);
    this.server.post(`/${this.backend}/config/crl`, (schema, req) => {
      assert.ok(true, 'request made to save crl config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          auto_rebuild: true,
          auto_rebuild_grace_period: '12h',
          delta_rebuild_interval: '72h',
          disable: false,
          enable_delta: true,
          expiry: '24h',
          ocsp_disable: false,
          ocsp_expiry: '18m',
        },
        'crl payload has correct data'
      );
    });
    this.server.post(`/${this.backend}/config/urls`, (schema, req) => {
      assert.ok(true, 'request made to save urls config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          crl_distribution_points: ['some-crl-distribution.com'],
          issuing_certificates: ['hashicorp.com'],
          ocsp_servers: ['ocsp-stuff.com'],
        },
        'url payload has correct data'
      );
    });
    await render(
      hbs`
      <Page::PkiConfigurationEdit
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    assert.dom(SELECTORS.urlsEditSection).exists('renders urls section');
    assert.dom(SELECTORS.crlEditSection).exists('renders crl section');
    assert.dom(SELECTORS.cancelButton).exists();
    this.urls.eachAttribute((name) => {
      assert.dom(SELECTORS.urlFieldInput(name)).exists(`renders ${name} input`);
    });
    assert.dom(SELECTORS.urlFieldInput('issuingCertificates')).hasValue('hashicorp.com');
    assert.dom(SELECTORS.urlFieldInput('crlDistributionPoints')).hasValue('some-crl-distribution.com');
    assert.dom(SELECTORS.urlFieldInput('ocspServers')).hasValue('ocsp-stuff.com');

    // confirm default toggle state and text
    this.crl.eachAttribute((name, { options }) => {
      if (['crlExpiryData', 'ocspExpiryData'].includes(name)) {
        assert.dom(SELECTORS.crlFieldInput(name)).isChecked(`${name} defaults to toggled on`);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.label);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.helperTextEnabled);
      }
      if (['autoRebuildData', 'deltaCrlBuildingData'].includes(name)) {
        assert.dom(SELECTORS.crlFieldInput(name)).isNotChecked(`${name} defaults off`);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.labelDisabled);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.helperTextDisabled);
      }
    });

    // toggle everything on
    await click(SELECTORS.crlFieldInput('autoRebuildData'));
    assert
      .dom(SELECTORS.crlFieldLabel('autoRebuildData'))
      .hasTextContaining(
        'Auto-rebuild on Vault will rebuild the CRL in the below grace period before expiration',
        'it renders auto rebuild toggled on text'
      );
    await click(SELECTORS.crlFieldInput('deltaCrlBuildingData'));
    assert
      .dom(SELECTORS.crlFieldLabel('deltaCrlBuildingData'))
      .hasTextContaining(
        'Delta CRL building on Vault will rebuild the delta CRL at the interval below:',
        'it renders delta crl build toggled on text'
      );
    await click(SELECTORS.saveButton);
  });

  test('it removes urls and sends false crl values', async function (assert) {
    assert.expect(8);
    this.server.post(`/${this.backend}/config/crl`, (schema, req) => {
      assert.ok(true, 'request made to save crl config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          auto_rebuild: false,
          auto_rebuild_grace_period: '12h',
          delta_rebuild_interval: '3d',
          disable: true,
          enable_delta: false,
          expiry: '24h',
          ocsp_disable: true,
          ocsp_expiry: '18m',
        },
        'crl payload has correct data'
      );
    });
    this.server.post(`/${this.backend}/config/urls`, (schema, req) => {
      assert.ok(true, 'request made to save urls config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          crl_distribution_points: [],
          issuing_certificates: [],
          ocsp_servers: [],
        },
        'url payload has empty arrays'
      );
    });
    await render(
      hbs`
      <Page::PkiConfigurationEdit
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    await click(SELECTORS.deleteButton('issuingCertificates'));
    await click(SELECTORS.deleteButton('crlDistributionPoints'));
    await click(SELECTORS.deleteButton('ocspServers'));

    // toggle everything off
    await click(SELECTORS.crlFieldInput('crlExpiryData'));
    assert.dom(SELECTORS.crlFieldLabel('crlExpiryData')).hasText('No expiry The CRL will not be built.');
    assert
      .dom(SELECTORS.crlFieldInput('autoRebuildData'))
      .doesNotExist('expiry off hides the auto rebuild toggle');
    assert
      .dom(SELECTORS.crlFieldInput('deltaCrlBuildingData'))
      .doesNotExist('expiry off hides delta crl toggle');
    await click(SELECTORS.crlFieldInput('ocspExpiryData'));
    assert
      .dom(SELECTORS.crlFieldLabel('ocspExpiryData'))
      .hasTextContaining(
        'OCSP responder APIs disabled Requests cannot be made to check if an individual certificate is valid.',
        'it renders correct toggled off text'
      );

    await click(SELECTORS.saveButton);
  });
});
