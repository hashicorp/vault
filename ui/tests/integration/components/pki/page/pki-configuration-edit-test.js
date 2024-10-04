/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import sinon from 'sinon';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { PKI_CONFIG_EDIT } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | page/pki-configuration-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    // test context setup
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    this.store = this.owner.lookup('service:store');
    this.router = this.owner.lookup('service:router');
    sinon.stub(this.router, 'transitionTo');

    // component data setup
    this.backend = 'pki-engine';
    // both models only use findRecord. API parameters for pki/crl
    // are set by default backend values when the engine is mounted
    this.store.pushPayload('pki/config/cluster', {
      modelName: 'pki/config/cluster',
      id: this.backend,
    });
    this.store.pushPayload('pki/config/acme', {
      modelName: 'pki/config/acme',
      id: this.backend,
    });
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
    this.acme = this.store.peekRecord('pki/config/acme', this.backend);
    this.cluster = this.store.peekRecord('pki/config/cluster', this.backend);
    this.crl = this.store.peekRecord('pki/config/crl', this.backend);
    this.urls = this.store.peekRecord('pki/config/urls', this.backend);
  });

  hooks.afterEach(function () {
    this.router.transitionTo.restore();
  });

  test('it renders with config data and updates config', async function (assert) {
    assert.expect(32);
    this.server.post(`/${this.backend}/config/acme`, (schema, req) => {
      assert.ok(true, 'request made to save acme config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          allowed_issuers: ['*'],
          allowed_roles: ['my-role'],
          dns_resolver: 'some-dns',
          eab_policy: 'new-account-required',
          enabled: true,
        },
        'it updates acme config model attributes'
      );
    });
    this.server.post(`/${this.backend}/config/cluster`, (schema, req) => {
      assert.ok(true, 'request made to save cluster config');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          path: 'https://pr-a.vault.example.com/v1/ns1/pki-root',
          aia_path: 'http://another-path.com',
        },
        'it updates cluster config model attributes'
      );
    });
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
        'it updates crl config model attributes'
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
        'it updates url config model attributes'
      );
    });
    await render(
      hbs`
      <Page::PkiConfigurationEdit
        @acme={{this.acme}}
        @cluster={{this.cluster}}
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    assert.dom(PKI_CONFIG_EDIT.configEditSection).exists('renders config section');
    assert.dom(PKI_CONFIG_EDIT.urlsEditSection).exists('renders urls section');
    assert.dom(PKI_CONFIG_EDIT.crlEditSection).exists('renders crl section');
    assert.dom(PKI_CONFIG_EDIT.cancelButton).exists();
    this.urls.eachAttribute((name) => {
      assert.dom(PKI_CONFIG_EDIT.urlFieldInput(name)).exists(`renders ${name} input`);
    });
    assert.dom(PKI_CONFIG_EDIT.urlFieldInput('issuingCertificates')).hasValue('hashicorp.com');
    assert.dom(PKI_CONFIG_EDIT.urlFieldInput('crlDistributionPoints')).hasValue('some-crl-distribution.com');
    assert.dom(PKI_CONFIG_EDIT.urlFieldInput('ocspServers')).hasValue('ocsp-stuff.com');

    // cluster config
    await fillIn(PKI_CONFIG_EDIT.configInput('path'), 'https://pr-a.vault.example.com/v1/ns1/pki-root');
    await fillIn(PKI_CONFIG_EDIT.configInput('aiaPath'), 'http://another-path.com');

    // acme config;
    await click(PKI_CONFIG_EDIT.configInput('enabled'));
    await fillIn(PKI_CONFIG_EDIT.stringListInput('allowedRoles'), 'my-role');
    await fillIn(PKI_CONFIG_EDIT.stringListInput('allowedIssuers'), '*');
    await fillIn(PKI_CONFIG_EDIT.configInput('eabPolicy'), 'new-account-required');
    await fillIn(PKI_CONFIG_EDIT.configInput('dnsResolver'), 'some-dns');

    // urls
    await fillIn(PKI_CONFIG_EDIT.urlFieldInput('issuingCertificates'), 'update-hashicorp.com');
    await fillIn(PKI_CONFIG_EDIT.urlFieldInput('crlDistributionPoints'), 'test-crl.com');
    await fillIn(PKI_CONFIG_EDIT.urlFieldInput('ocspServers'), 'ocsp.com');

    // confirm default toggle state and text
    this.crl.eachAttribute((name, { options }) => {
      if (['expiry', 'ocspExpiry'].includes(name)) {
        assert.dom(PKI_CONFIG_EDIT.crlToggleInput(name)).isChecked(`${name} defaults to toggled on`);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.label);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.helperTextEnabled);
      }
      if (['autoRebuildGracePeriod', 'deltaRebuildInterval'].includes(name)) {
        assert.dom(PKI_CONFIG_EDIT.crlToggleInput(name)).isNotChecked(`${name} defaults off`);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.labelDisabled);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.helperTextDisabled);
      }
    });

    // toggle everything on
    await click(PKI_CONFIG_EDIT.crlToggleInput('autoRebuildGracePeriod'));
    assert
      .dom(PKI_CONFIG_EDIT.crlFieldLabel('autoRebuildGracePeriod'))
      .hasTextContaining(
        'Auto-rebuild on Vault will rebuild the CRL in the below grace period before expiration',
        'it renders auto rebuild toggled on text'
      );
    await click(PKI_CONFIG_EDIT.crlToggleInput('deltaRebuildInterval'));
    assert
      .dom(PKI_CONFIG_EDIT.crlFieldLabel('deltaRebuildInterval'))
      .hasTextContaining(
        'Delta CRL building on Vault will rebuild the delta CRL at the interval below:',
        'it renders delta crl build toggled on text'
      );

    // assert ttl values update model attributes
    await fillIn(PKI_CONFIG_EDIT.crlTtlInput('Expiry'), '48');
    await fillIn(PKI_CONFIG_EDIT.crlTtlInput('Auto-rebuild on'), '24');
    await fillIn(PKI_CONFIG_EDIT.crlTtlInput('Delta CRL building on'), '45');
    await fillIn(PKI_CONFIG_EDIT.crlTtlInput('OCSP responder APIs enabled'), '24');

    await click(PKI_CONFIG_EDIT.saveButton);
  });

  test('it removes urls and sends false crl values', async function (assert) {
    assert.expect(8);
    this.server.post(`/${this.backend}/config/acme`, () => {});
    this.server.post(`/${this.backend}/config/cluster`, () => {});
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
        @acme={{this.acme}}
        @cluster={{this.cluster}}
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    await click(PKI_CONFIG_EDIT.deleteButton('issuingCertificates'));
    await click(PKI_CONFIG_EDIT.deleteButton('crlDistributionPoints'));
    await click(PKI_CONFIG_EDIT.deleteButton('ocspServers'));

    // toggle everything off
    await click(PKI_CONFIG_EDIT.crlToggleInput('expiry'));
    assert.dom(PKI_CONFIG_EDIT.crlFieldLabel('expiry')).hasText('No expiry The CRL will not be built.');
    assert
      .dom(PKI_CONFIG_EDIT.crlToggleInput('autoRebuildGracePeriod'))
      .doesNotExist('expiry off hides the auto rebuild toggle');
    assert
      .dom(PKI_CONFIG_EDIT.crlToggleInput('deltaRebuildInterval'))
      .doesNotExist('expiry off hides delta crl toggle');
    await click(PKI_CONFIG_EDIT.crlToggleInput('ocspExpiry'));
    assert
      .dom(PKI_CONFIG_EDIT.crlFieldLabel('ocspExpiry'))
      .hasTextContaining(
        'OCSP responder APIs disabled Requests cannot be made to check if an individual certificate is valid.',
        'it renders correct toggled off text'
      );

    await click(PKI_CONFIG_EDIT.saveButton);
  });

  test('it renders enterprise only params', async function (assert) {
    assert.expect(6);
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    this.server.post(`/${this.backend}/config/acme`, () => {});
    this.server.post(`/${this.backend}/config/cluster`, () => {});
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
        @acme={{this.acme}}
        @cluster={{this.cluster}}
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    assert.dom(PKI_CONFIG_EDIT.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Unified Revocation')).exists();
    await click(PKI_CONFIG_EDIT.checkboxInput('crossClusterRevocation'));
    await click(PKI_CONFIG_EDIT.checkboxInput('unifiedCrl'));
    await click(PKI_CONFIG_EDIT.checkboxInput('unifiedCrlOnExistingPaths'));
    await click(PKI_CONFIG_EDIT.saveButton);
  });

  test('it does not render enterprise only params for OSS', async function (assert) {
    assert.expect(9);
    this.version = this.owner.lookup('service:version');
    this.version.type = 'community';
    this.server.post(`/${this.backend}/config/acme`, () => {});
    this.server.post(`/${this.backend}/config/cluster`, () => {});
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
        @acme={{this.acme}}
        @cluster={{this.cluster}}
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    assert.dom(PKI_CONFIG_EDIT.checkboxInput('crossClusterRevocation')).doesNotExist();
    assert.dom(PKI_CONFIG_EDIT.checkboxInput('unifiedCrl')).doesNotExist();
    assert.dom(PKI_CONFIG_EDIT.checkboxInput('unifiedCrlOnExistingPaths')).doesNotExist();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Unified Revocation')).doesNotExist();
    await click(PKI_CONFIG_EDIT.saveButton);
  });

  test('it renders empty states if no update capabilities', async function (assert) {
    assert.expect(4);
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub(['read']));

    await render(
      hbs`
      <Page::PkiConfigurationEdit
        @acme={{this.acme}}
        @cluster={{this.cluster}}
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    assert
      .dom(`${PKI_CONFIG_EDIT.configEditSection} [data-test-component="empty-state"]`)
      .hasText(
        "You do not have permission to set this mount's the cluster config Ask your administrator if you think you should have access to: POST /pki-engine/config/cluster"
      );
    assert
      .dom(`${PKI_CONFIG_EDIT.acmeEditSection} [data-test-component="empty-state"]`)
      .hasText(
        "You do not have permission to set this mount's ACME config Ask your administrator if you think you should have access to: POST /pki-engine/config/acme"
      );
    assert
      .dom(`${PKI_CONFIG_EDIT.urlsEditSection} [data-test-component="empty-state"]`)
      .hasText(
        "You do not have permission to set this mount's URLs Ask your administrator if you think you should have access to: POST /pki-engine/config/urls"
      );
    assert
      .dom(`${PKI_CONFIG_EDIT.crlEditSection} [data-test-component="empty-state"]`)
      .hasText(
        "You do not have permission to set this mount's revocation configuration Ask your administrator if you think you should have access to: POST /pki-engine/config/crl"
      );
  });

  test('it renders alert banner and endpoint respective error', async function (assert) {
    assert.expect(4);
    this.server.post(`/${this.backend}/config/acme`, () => {
      return new Response(500, {}, { errors: ['something wrong with acme'] });
    });
    this.server.post(`/${this.backend}/config/cluster`, () => {
      return new Response(500, {}, { errors: ['something wrong with cluster'] });
    });
    this.server.post(`/${this.backend}/config/crl`, () => {
      return new Response(500, {}, { errors: ['something wrong with crl'] });
    });
    this.server.post(`/${this.backend}/config/urls`, () => {
      return new Response(500, {}, { errors: ['something wrong with urls'] });
    });
    await render(
      hbs`
      <Page::PkiConfigurationEdit
        @acme={{this.acme}}
        @cluster={{this.cluster}}
        @urls={{this.urls}}
        @crl={{this.crl}}
        @backend={{this.backend}}
      />
    `,
      this.context
    );

    await click(PKI_CONFIG_EDIT.saveButton);
    assert
      .dom(PKI_CONFIG_EDIT.errorBanner)
      .hasText(
        'Error POST config/cluster: something wrong with cluster POST config/acme: something wrong with acme POST config/urls: something wrong with urls POST config/crl: something wrong with crl'
      );
    assert.dom(`${PKI_CONFIG_EDIT.errorBanner} ul`).hasClass('bullet');

    // change 3 out of 4 requests to be successful to assert single error renders correctly
    this.server.post(`/${this.backend}/config/acme`, () => new Response(200));
    this.server.post(`/${this.backend}/config/cluster`, () => new Response(200));
    this.server.post(`/${this.backend}/config/crl`, () => new Response(200));

    await click(PKI_CONFIG_EDIT.saveButton);
    assert.dom(PKI_CONFIG_EDIT.errorBanner).hasText('Error POST config/urls: something wrong with urls');
    assert.dom(`${PKI_CONFIG_EDIT.errorBanner} ul`).doesNotHaveClass('bullet');
  });
});
