/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
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
    // both models only use findRecord. API parameters for pki/crl
    // are set by default backend values when the engine is mounted
    this.store.pushPayload('pki/config/crl', {
      modelName: 'pki/config/crl',
      id: this.backend,
      auto_rebuild: false,
      auto_rebuild_grace_period: '12h',
      delta_rebuild_interval: '15m',
      disable: false,
      enable_delta: false,
      expiry: '72h',
      ocsp_disable: false,
      ocsp_expiry: '12h',
    });
    this.store.pushPayload('pki/config/urls', {
      modelName: 'pki/config/urls',
      id: this.backend,
      issuing_certificates: ['hashicorp.com'],
      crl_distribution_points: ['some-crl-distribution.com'],
      ocsp_servers: ['ocsp-stuff.com'],
    });
    this.urls = this.store.peekRecord('pki/config/urls', this.backend);
    this.crl = this.store.peekRecord('pki/config/crl', this.backend);
  });

  test('it renders with config data and updates config', async function (assert) {
    assert.expect(27);
    this.server.post(`/${this.backend}/config/crl`, (schema, req) => {
      assert.ok(true, 'request made to save crl config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          auto_rebuild: true,
          auto_rebuild_grace_period: '24h',
          delta_rebuild_interval: '45m',
          disable: false,
          enable_delta: true,
          expiry: '1152h',
          ocsp_disable: false,
          ocsp_expiry: '24h',
        },
        'it updates crl model attributes'
      );
    });
    this.server.post(`/${this.backend}/config/urls`, (schema, req) => {
      assert.ok(true, 'request made to save urls config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          crl_distribution_points: ['test-crl.com'],
          issuing_certificates: ['update-hashicorp.com'],
          ocsp_servers: ['ocsp.com'],
        },
        'it updates url model attributes'
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

    await fillIn(SELECTORS.urlFieldInput('issuingCertificates'), 'update-hashicorp.com');
    await fillIn(SELECTORS.urlFieldInput('crlDistributionPoints'), 'test-crl.com');
    await fillIn(SELECTORS.urlFieldInput('ocspServers'), 'ocsp.com');

    // confirm default toggle state and text
    this.crl.eachAttribute((name, { options }) => {
      if (['expiry', 'ocspExpiry'].includes(name)) {
        assert.dom(SELECTORS.crlToggleInput(name)).isChecked(`${name} defaults to toggled on`);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.label);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.helperTextEnabled);
      }
      if (['autoRebuildGracePeriod', 'deltaRebuildInterval'].includes(name)) {
        assert.dom(SELECTORS.crlToggleInput(name)).isNotChecked(`${name} defaults off`);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.labelDisabled);
        assert.dom(SELECTORS.crlFieldLabel(name)).hasTextContaining(options.helperTextDisabled);
      }
    });

    // toggle everything on
    await click(SELECTORS.crlToggleInput('autoRebuildGracePeriod'));
    assert
      .dom(SELECTORS.crlFieldLabel('autoRebuildGracePeriod'))
      .hasTextContaining(
        'Auto-rebuild on Vault will rebuild the CRL in the below grace period before expiration',
        'it renders auto rebuild toggled on text'
      );
    await click(SELECTORS.crlToggleInput('deltaRebuildInterval'));
    assert
      .dom(SELECTORS.crlFieldLabel('deltaRebuildInterval'))
      .hasTextContaining(
        'Delta CRL building on Vault will rebuild the delta CRL at the interval below:',
        'it renders delta crl build toggled on text'
      );

    // assert ttl values update model attributes
    await fillIn(SELECTORS.crlTtlInput('Expiry'), '48');
    await fillIn(SELECTORS.crlTtlInput('Auto-rebuild on'), '24');
    await fillIn(SELECTORS.crlTtlInput('Delta CRL building on'), '45');
    await fillIn(SELECTORS.crlTtlInput('OCSP responder APIs enabled'), '24');
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
          delta_rebuild_interval: '15m',
          disable: true,
          enable_delta: false,
          expiry: '72h',
          ocsp_disable: true,
          ocsp_expiry: '12h',
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
    await click(SELECTORS.crlToggleInput('expiry'));
    assert.dom(SELECTORS.crlFieldLabel('expiry')).hasText('No expiry The CRL will not be built.');
    assert
      .dom(SELECTORS.crlToggleInput('autoRebuildGracePeriod'))
      .doesNotExist('expiry off hides the auto rebuild toggle');
    assert
      .dom(SELECTORS.crlToggleInput('deltaRebuildInterval'))
      .doesNotExist('expiry off hides delta crl toggle');
    await click(SELECTORS.crlToggleInput('ocspExpiry'));
    assert
      .dom(SELECTORS.crlFieldLabel('ocspExpiry'))
      .hasTextContaining(
        'OCSP responder APIs disabled Requests cannot be made to check if an individual certificate is valid.',
        'it renders correct toggled off text'
      );

    await click(SELECTORS.saveButton);
  });

  test('it renders enterprise only params', async function (assert) {
    assert.expect(6);
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    this.server.post(`/${this.backend}/config/crl`, (schema, req) => {
      assert.ok(true, 'request made to save crl config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          auto_rebuild: false,
          auto_rebuild_grace_period: '12h',
          delta_rebuild_interval: '15m',
          disable: false,
          enable_delta: false,
          expiry: '72h',
          ocsp_disable: false,
          ocsp_expiry: '12h',
          cross_cluster_revocation: true,
          unified_crl: true,
          unified_crl_on_existing_paths: true,
        },
        'crl payload includes enterprise params'
      );
    });
    this.server.post(`/${this.backend}/config/urls`, () => {
      assert.ok(true, 'request made to save urls config');
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
    assert.dom(SELECTORS.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(SELECTORS.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(SELECTORS.groupHeader('Unified Revocation')).exists();
    await click(SELECTORS.checkboxInput('crossClusterRevocation'));
    await click(SELECTORS.checkboxInput('unifiedCrl'));
    await click(SELECTORS.checkboxInput('unifiedCrlOnExistingPaths'));
    await click(SELECTORS.saveButton);
  });

  test('it renders does not render enterprise only params for OSS', async function (assert) {
    assert.expect(9);
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1';
    this.server.post(`/${this.backend}/config/crl`, (schema, req) => {
      assert.ok(true, 'request made to save crl config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          auto_rebuild: false,
          auto_rebuild_grace_period: '12h',
          delta_rebuild_interval: '15m',
          disable: false,
          enable_delta: false,
          expiry: '72h',
          ocsp_disable: false,
          ocsp_expiry: '12h',
        },
        'crl payload does not include enterprise params'
      );
    });
    this.server.post(`/${this.backend}/config/urls`, () => {
      assert.ok(true, 'request made to save urls config');
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

    assert.dom(SELECTORS.checkboxInput('crossClusterRevocation')).doesNotExist();
    assert.dom(SELECTORS.checkboxInput('unifiedCrl')).doesNotExist();
    assert.dom(SELECTORS.checkboxInput('unifiedCrlOnExistingPaths')).doesNotExist();
    assert.dom(SELECTORS.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(SELECTORS.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(SELECTORS.groupHeader('Unified Revocation')).doesNotExist();
    await click(SELECTORS.saveButton);
  });
});
