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
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-configuration-edit';
import sinon from 'sinon';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

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

    assert.dom(SELECTORS.configEditSection).exists('renders config section');
    assert.dom(SELECTORS.urlsEditSection).exists('renders urls section');
    assert.dom(SELECTORS.crlEditSection).exists('renders crl section');
    assert.dom(SELECTORS.cancelButton).exists();
    this.urls.eachAttribute((name) => {
      assert.dom(SELECTORS.urlFieldInput(name)).exists(`renders ${name} input`);
    });
    assert.dom(SELECTORS.urlFieldInput('issuingCertificates')).hasValue('hashicorp.com');
    assert.dom(SELECTORS.urlFieldInput('crlDistributionPoints')).hasValue('some-crl-distribution.com');
    assert.dom(SELECTORS.urlFieldInput('ocspServers')).hasValue('ocsp-stuff.com');

    // cluster config
    await fillIn(SELECTORS.configInput('path'), 'https://pr-a.vault.example.com/v1/ns1/pki-root');
    await fillIn(SELECTORS.configInput('aiaPath'), 'http://another-path.com');

    // acme config;
    await click(SELECTORS.configInput('enabled'));
    await fillIn(SELECTORS.stringListInput('allowedRoles'), 'my-role');
    await fillIn(SELECTORS.stringListInput('allowedIssuers'), '*');
    await fillIn(SELECTORS.configInput('eabPolicy'), 'new-account-required');
    await fillIn(SELECTORS.configInput('dnsResolver'), 'some-dns');

    // urls
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

    assert.dom(SELECTORS.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(SELECTORS.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(SELECTORS.groupHeader('Unified Revocation')).exists();
    await click(SELECTORS.checkboxInput('crossClusterRevocation'));
    await click(SELECTORS.checkboxInput('unifiedCrl'));
    await click(SELECTORS.checkboxInput('unifiedCrlOnExistingPaths'));
    await click(SELECTORS.saveButton);
  });

  test('it does not render enterprise only params for OSS', async function (assert) {
    assert.expect(9);
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1';
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

    assert.dom(SELECTORS.checkboxInput('crossClusterRevocation')).doesNotExist();
    assert.dom(SELECTORS.checkboxInput('unifiedCrl')).doesNotExist();
    assert.dom(SELECTORS.checkboxInput('unifiedCrlOnExistingPaths')).doesNotExist();
    assert.dom(SELECTORS.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(SELECTORS.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(SELECTORS.groupHeader('Unified Revocation')).doesNotExist();
    await click(SELECTORS.saveButton);
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
      .dom(`${SELECTORS.configEditSection} [data-test-component="empty-state"]`)
      .hasText(
        "You do not have permission to set this mount's the cluster config Ask your administrator if you think you should have access to: POST /pki-engine/config/cluster"
      );
    assert
      .dom(`${SELECTORS.acmeEditSection} [data-test-component="empty-state"]`)
      .hasText(
        "You do not have permission to set this mount's ACME config Ask your administrator if you think you should have access to: POST /pki-engine/config/acme"
      );
    assert
      .dom(`${SELECTORS.urlsEditSection} [data-test-component="empty-state"]`)
      .hasText(
        "You do not have permission to set this mount's URLs Ask your administrator if you think you should have access to: POST /pki-engine/config/urls"
      );
    assert
      .dom(`${SELECTORS.crlEditSection} [data-test-component="empty-state"]`)
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

    await click(SELECTORS.saveButton);
    assert
      .dom(SELECTORS.errorBanner)
      .hasText(
        'Error POST config/cluster: something wrong with cluster POST config/acme: something wrong with acme POST config/urls: something wrong with urls POST config/crl: something wrong with crl'
      );
    assert.dom(`${SELECTORS.errorBanner} ul`).hasClass('bullet');

    // change 3 out of 4 requests to be successful to assert single error renders correctly
    this.server.post(`/${this.backend}/config/acme`, () => new Response(200));
    this.server.post(`/${this.backend}/config/cluster`, () => new Response(200));
    this.server.post(`/${this.backend}/config/crl`, () => new Response(200));

    await click(SELECTORS.saveButton);
    assert.dom(SELECTORS.errorBanner).hasText('Error POST config/urls: something wrong with urls');
    assert.dom(`${SELECTORS.errorBanner} ul`).doesNotHaveClass('bullet');
  });
});
