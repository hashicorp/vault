/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { PKI_CONFIG_EDIT } from 'vault/tests/helpers/pki/pki-selectors';
import { configCapabilities } from 'vault/tests/helpers/pki/pki-helpers';
import PkiConfigClusterForm from 'vault/forms/secrets/pki/config/cluster';
import PkiConfigAcmeForm from 'vault/forms/secrets/pki/config/acme';
import PkiConfigCrlForm from 'vault/forms/secrets/pki/config/crl';
import PkiConfigUrlsForm from 'vault/forms/secrets/pki/config/urls';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | page/pki-configuration-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    // test context setup
    this.store = this.owner.lookup('service:store');
    this.router = this.owner.lookup('service:router');
    sinon.stub(this.router, 'transitionTo');

    // component data setup
    this.backend = 'pki-engine';
    this.clusterForm = new PkiConfigClusterForm();
    this.acmeForm = new PkiConfigAcmeForm();
    this.crlForm = new PkiConfigCrlForm({
      auto_rebuild: false,
      auto_rebuild_grace_period: '12h',
      delta_rebuild_interval: '15m',
      disable: false,
      enable_delta: false,
      expiry: '72h',
      ocsp_disable: false,
      ocsp_expiry: '12h',
    });
    this.urlsForm = new PkiConfigUrlsForm({
      issuing_certificates: ['hashicorp.com'],
      crl_distribution_points: ['some-crl-distribution.com'],
      ocsp_servers: ['ocsp-stuff.com'],
    });
    this.capabilities = configCapabilities;

    // api stubs
    const { secrets } = this.owner.lookup('service:api');
    this.clusterStub = sinon.stub(secrets, 'pkiConfigureCluster').resolves();
    this.acmeStub = sinon.stub(secrets, 'pkiConfigureAcme').resolves();
    this.urlsStub = sinon.stub(secrets, 'pkiConfigureUrls').resolves();
    this.crlStub = sinon.stub(secrets, 'pkiConfigureCrl').resolves();

    this.renderComponent = () =>
      render(
        hbs`
        <Page::PkiConfigurationEdit
          @acmeForm={{this.acmeForm}}
          @clusterForm={{this.clusterForm}}
          @urlsForm={{this.urlsForm}}
          @crlForm={{this.crlForm}}
          @backend={{this.backend}}
          @capabilities={{this.capabilities}}
        />
      `,
        { owner: this.engine }
      );
  });

  test('it renders with config data and updates config', async function (assert) {
    assert.expect(30);

    await this.renderComponent();

    assert.dom(PKI_CONFIG_EDIT.configEditSection).exists('renders config section');
    assert.dom(PKI_CONFIG_EDIT.urlsEditSection).exists('renders urls section');
    assert.dom(PKI_CONFIG_EDIT.crlEditSection).exists('renders crl section');
    assert.dom(PKI_CONFIG_EDIT.cancelButton).exists();
    this.urlsForm.formFields.forEach(({ name }) => {
      const isTextarea = name !== 'enable_templating';
      assert.dom(PKI_CONFIG_EDIT.urlFieldInput(name, isTextarea)).exists(`renders ${name} input`);
    });
    assert.dom(PKI_CONFIG_EDIT.urlFieldInput('issuing_certificates')).hasValue('hashicorp.com');
    assert
      .dom(PKI_CONFIG_EDIT.urlFieldInput('crl_distribution_points'))
      .hasValue('some-crl-distribution.com');
    assert.dom(PKI_CONFIG_EDIT.urlFieldInput('ocsp_servers')).hasValue('ocsp-stuff.com');

    // cluster config
    await fillIn(PKI_CONFIG_EDIT.configInput('path'), 'https://pr-a.vault.example.com/v1/ns1/pki-root');
    await fillIn(PKI_CONFIG_EDIT.configInput('aia_path'), 'http://another-path.com');

    // acme config;
    await click(PKI_CONFIG_EDIT.configInput('enabled'));
    await fillIn(PKI_CONFIG_EDIT.stringListInput('allowed_roles'), 'my-role');
    await fillIn(PKI_CONFIG_EDIT.stringListInput('allowed_issuers'), '*');
    await fillIn(PKI_CONFIG_EDIT.configInput('eab_policy'), 'new-account-required');
    await fillIn(PKI_CONFIG_EDIT.configInput('dns_resolver'), 'some-dns');

    // urls
    await fillIn(PKI_CONFIG_EDIT.urlFieldInput('issuing_certificates'), 'update-hashicorp.com');
    await fillIn(PKI_CONFIG_EDIT.urlFieldInput('crl_distribution_points'), 'test-crl.com');
    await fillIn(PKI_CONFIG_EDIT.urlFieldInput('ocsp_servers'), 'ocsp.com');

    // confirm default toggle state and text
    this.crlForm.formFields.forEach(({ name, options }) => {
      if (['expiry', 'ocsp_expiry'].includes(name)) {
        assert.dom(PKI_CONFIG_EDIT.crlToggleInput(name)).isChecked(`${name} defaults to toggled on`);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.label);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.helperTextEnabled);
      }
      if (['auto_rebuild_grace_period', 'delta_rebuild_interval'].includes(name)) {
        assert.dom(PKI_CONFIG_EDIT.crlToggleInput(name)).isNotChecked(`${name} defaults off`);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.labelDisabled);
        assert.dom(PKI_CONFIG_EDIT.crlFieldLabel(name)).hasTextContaining(options.helperTextDisabled);
      }
    });

    // toggle everything on
    await click(PKI_CONFIG_EDIT.crlToggleInput('auto_rebuild_grace_period'));
    assert
      .dom(PKI_CONFIG_EDIT.crlFieldLabel('auto_rebuild_grace_period'))
      .hasTextContaining(
        'Auto-rebuild on Vault will rebuild the CRL in the below grace period before expiration',
        'it renders auto rebuild toggled on text'
      );
    await click(PKI_CONFIG_EDIT.crlToggleInput('delta_rebuild_interval'));
    assert
      .dom(PKI_CONFIG_EDIT.crlFieldLabel('delta_rebuild_interval'))
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

    assert.true(
      this.acmeStub.calledWith(this.backend, {
        allowed_issuers: ['*'],
        allowed_roles: ['my-role'],
        dns_resolver: 'some-dns',
        eab_policy: 'new-account-required',
        enabled: true,
      }),
      'request made to update acme config'
    );
    assert.true(
      this.clusterStub.calledWith(this.backend, {
        path: 'https://pr-a.vault.example.com/v1/ns1/pki-root',
        aia_path: 'http://another-path.com',
      }),
      'request made to update cluster config'
    );
    assert.true(
      this.urlsStub.calledWith(this.backend, {
        crl_distribution_points: ['test-crl.com'],
        issuing_certificates: ['update-hashicorp.com'],
        ocsp_servers: ['ocsp.com'],
      }),
      'request made to update urls config'
    );
    assert.true(
      this.crlStub.calledWith(this.backend, {
        auto_rebuild: true,
        auto_rebuild_grace_period: '24h',
        delta_rebuild_interval: '45m',
        disable: false,
        enable_delta: true,
        expiry: '1152h',
        ocsp_disable: false,
        ocsp_expiry: '24h',
      }),
      'request made to update crl config'
    );
  });

  test('it removes urls and sends false crl values', async function (assert) {
    assert.expect(6);

    await this.renderComponent();

    await click(PKI_CONFIG_EDIT.deleteButton('issuing_certificates'));
    await click(PKI_CONFIG_EDIT.deleteButton('crl_distribution_points'));
    await click(PKI_CONFIG_EDIT.deleteButton('ocsp_servers'));

    // toggle everything off
    await click(PKI_CONFIG_EDIT.crlToggleInput('expiry'));
    assert.dom(PKI_CONFIG_EDIT.crlFieldLabel('expiry')).hasText('No expiry The CRL will not be built.');
    assert
      .dom(PKI_CONFIG_EDIT.crlToggleInput('auto_rebuild_grace_period'))
      .doesNotExist('expiry off hides the auto rebuild toggle');
    assert
      .dom(PKI_CONFIG_EDIT.crlToggleInput('delta_rebuild_interval'))
      .doesNotExist('expiry off hides delta crl toggle');
    await click(PKI_CONFIG_EDIT.crlToggleInput('ocsp_expiry'));
    assert
      .dom(PKI_CONFIG_EDIT.crlFieldLabel('ocsp_expiry'))
      .hasTextContaining(
        'OCSP responder APIs disabled Requests cannot be made to check if an individual certificate is valid.',
        'it renders correct toggled off text'
      );

    await click(PKI_CONFIG_EDIT.saveButton);

    assert.true(
      this.crlStub.calledWith(this.backend, {
        auto_rebuild: false,
        auto_rebuild_grace_period: '12h',
        delta_rebuild_interval: '15m',
        disable: true,
        enable_delta: false,
        expiry: '72h',
        ocsp_disable: true,
        ocsp_expiry: '12h',
      }),
      'request made to update crl config with false values'
    );
    assert.true(
      this.urlsStub.calledWith(this.backend, {
        crl_distribution_points: [],
        issuing_certificates: [],
        ocsp_servers: [],
      }),
      'request made to update urls config with empty arrays'
    );
  });

  test('it renders enterprise only params', async function (assert) {
    assert.expect(4);

    this.owner.lookup('service:version').type = 'enterprise';

    await this.renderComponent();

    assert.dom(PKI_CONFIG_EDIT.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Unified Revocation')).exists();
    await click(PKI_CONFIG_EDIT.checkboxInput('cross_cluster_revocation'));
    await click(PKI_CONFIG_EDIT.checkboxInput('unified_crl'));
    await click(PKI_CONFIG_EDIT.checkboxInput('unified_crl_on_existing_paths'));
    await click(PKI_CONFIG_EDIT.saveButton);

    assert.true(
      this.crlStub.calledWith(this.backend, {
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
      }),
      'request made to save crl config with enterprise params'
    );
  });

  test('it does not render enterprise only params for OSS', async function (assert) {
    assert.expect(7);

    this.owner.lookup('service:version').type = 'community';

    await this.renderComponent();

    assert.dom(PKI_CONFIG_EDIT.checkboxInput('cross_cluster_revocation')).doesNotExist();
    assert.dom(PKI_CONFIG_EDIT.checkboxInput('unified_crl')).doesNotExist();
    assert.dom(PKI_CONFIG_EDIT.checkboxInput('unified_crl_on_existing_paths')).doesNotExist();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Certificate Revocation List (CRL)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Online Certificate Status Protocol (OCSP)')).exists();
    assert.dom(PKI_CONFIG_EDIT.groupHeader('Unified Revocation')).doesNotExist();
    await click(PKI_CONFIG_EDIT.saveButton);

    assert.true(
      this.crlStub.calledWith(this.backend, {
        auto_rebuild: false,
        auto_rebuild_grace_period: '12h',
        delta_rebuild_interval: '15m',
        disable: false,
        enable_delta: false,
        expiry: '72h',
        ocsp_disable: false,
        ocsp_expiry: '12h',
      }),
      'request made to save crl config without enterprise params'
    );
  });

  test('it renders empty states if no update capabilities', async function (assert) {
    assert.expect(4);

    this.capabilities = {};

    await this.renderComponent();

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

    this.acmeStub.rejects(getErrorResponse({ errors: ['something wrong with acme'] }, 500));
    this.clusterStub.rejects(getErrorResponse({ errors: ['something wrong with cluster'] }, 500));
    this.crlStub.rejects(getErrorResponse({ errors: ['something wrong with crl'] }, 500));
    this.urlsStub.rejects(getErrorResponse({ errors: ['something wrong with urls'] }, 500));

    await this.renderComponent();

    await click(PKI_CONFIG_EDIT.saveButton);
    assert
      .dom(PKI_CONFIG_EDIT.errorBanner)
      .hasText(
        'Error POST config/cluster: something wrong with cluster POST config/acme: something wrong with acme POST config/urls: something wrong with urls POST config/crl: something wrong with crl'
      );
    assert.dom(`${PKI_CONFIG_EDIT.errorBanner} ul`).hasClass('bullet');

    // change 3 out of 4 requests to be successful to assert single error renders correctly
    this.acmeStub.resolves();
    this.clusterStub.resolves();
    this.crlStub.resolves();

    await click(PKI_CONFIG_EDIT.saveButton);
    assert.dom(PKI_CONFIG_EDIT.errorBanner).hasText('Error POST config/urls: something wrong with urls');
    assert.dom(`${PKI_CONFIG_EDIT.errorBanner} ul`).doesNotHaveClass('bullet');
  });
});
