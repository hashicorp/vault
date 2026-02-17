/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  rowIcon: (attr, icon) => `${GENERAL.infoRowValue(attr)} ${GENERAL.icon(icon)}`,
};

module('Integration | Component | Page::PkiConfigurationDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.cluster = {
      path: 'https://pr-a.vault.example.com/v1/ns1/pki-root',
    };
    this.urls = {
      issuing_certificates: 'example.com',
    };
    this.crl = {
      expiry: '20h',
      disable: false,
      auto_rebuild: true,
      auto_rebuild_grace_period: '13h',
      enable_delta: true,
      delta_rebuild_interval: '15m',
      ocsp_expiry: '77h',
      ocsp_disable: false,
      cross_cluster_revocation: true,
      unified_crl: true,
      unified_crl_on_existing_paths: true,
    };
    this.acme = {
      enabled: true,
      default_directory_policy: 'foo',
      allowed_roles: 'test',
      allow_role_ext_key_usage: false,
      allowed_issuers: 'bar',
      eab_policy: 'not-required',
      dns_resolver: 'resolver',
      max_ttl: '72h',
    };
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });

    this.renderComponent = () =>
      render(
        hbs`
        <Page::PkiConfigurationDetails
          @acme={{this.acme}}
          @cluster={{this.cluster}}
          @urls={{this.urls}}
          @crl={{this.crl}}
          @backend="pki-test"
          @canDeleteAllIssuers={{this.canDeleteAllIssuers}}
        />
      `,
        { owner: this.engine }
      );
  });

  test('shows the correct information on cluster config', async function (assert) {
    await this.renderComponent();

    assert
      .dom(GENERAL.infoRowValue("Mount's API path"))
      .hasText('https://pr-a.vault.example.com/v1/ns1/pki-root', 'mount API path row renders');
    assert.dom(GENERAL.infoRowValue('AIA path')).hasText('None', "renders 'None' when no data");
  });

  test('shows the correct information on acme config', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('ACME enabled')).hasText('Yes', 'enabled value renders');
    assert
      .dom(GENERAL.infoRowValue('Default directory policy'))
      .hasText('foo', 'default_directory_policy value renders');
    assert.dom(GENERAL.infoRowValue('Allowed roles')).hasText('test', 'allowed_roles value renders');
    assert
      .dom(GENERAL.infoRowValue('Allow role ExtKeyUsage'))
      .hasText('None', 'allow_role_ext_key_usage value renders');
    assert.dom(GENERAL.infoRowValue('Allowed issuers')).hasText('bar', 'allowed_issuers value renders');
    assert.dom(GENERAL.infoRowValue('EAB policy')).hasText('not-required', 'eab_policy value renders');
    assert.dom(GENERAL.infoRowValue('DNS resolver')).hasText('resolver', 'dns_resolver value renders');
    assert.dom(GENERAL.infoRowValue('Max TTL')).hasText('3 days', 'max ttl value renders');
  });

  test('shows the correct information on global urls section', async function (assert) {
    await this.renderComponent();

    assert
      .dom(GENERAL.infoRowLabel('Issuing certificates'))
      .hasText('Issuing certificates', 'issuing certificate row label renders');
    assert
      .dom(GENERAL.infoRowValue('Issuing certificates'))
      .hasText('example.com', 'issuing certificate value renders');

    this.urls.issuing_certificates = null;
    await this.renderComponent();

    assert
      .dom(GENERAL.infoRowValue('Issuing certificates'))
      .hasText('None', 'issuing certificate value renders None if none is configured');
    assert
      .dom(GENERAL.infoRowLabel('CRL distribution points'))
      .hasText('CRL distribution points', 'crl distribution points row label renders');
    assert
      .dom(GENERAL.infoRowValue('CRL distribution points'))
      .hasText('None', 'crl distribution points value renders None if none is configured');
  });

  test('shows the correct information on crl section', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.infoRowLabel('CRL building')).hasText('CRL building', 'crl expiry row label renders');
    assert.dom(GENERAL.infoRowValue('CRL building')).hasText('Enabled', 'enabled renders');
    assert.dom(GENERAL.infoRowValue('Expiry')).hasText('20h', 'expiry value renders');
    assert.dom(GENERAL.infoRowLabel('Auto-rebuild')).hasText('Auto-rebuild', 'auto rebuild label renders');
    assert.dom(GENERAL.infoRowValue('Auto-rebuild')).hasText('On', 'it renders truthy auto build');
    assert.dom(SELECTORS.rowIcon('Auto-rebuild', 'check-circle'));
    assert
      .dom(GENERAL.infoRowValue('Auto-rebuild grace period'))
      .hasText('13h', 'it renders auto build grace period');
    assert.dom(GENERAL.infoRowValue('Delta CRL building')).hasText('On', 'it renders truthy delta crl build');
    assert.dom(SELECTORS.rowIcon('Delta CRL building', 'check-circle'));
    assert
      .dom(GENERAL.infoRowValue('Delta rebuild interval'))
      .hasText('15m', 'it renders delta build duration');
    assert
      .dom(GENERAL.infoRowValue('Responder APIs'))
      .hasText('Enabled', 'responder apis value renders Enabled if ocsp_disable=false');
    assert.dom(GENERAL.infoRowValue('Interval')).hasText('77h', 'interval value renders');
    // check falsy aut_rebuild and _enable_delta hides duration values
    this.crl.auto_rebuild = false;
    this.crl.enable_delta = false;
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('Auto-rebuild')).hasText('Off', 'it renders falsy auto build');
    assert.dom(SELECTORS.rowIcon('Auto-rebuild', 'x-square'));
    assert
      .dom(GENERAL.infoRowValue('Auto-rebuild grace period'))
      .doesNotExist('does not render auto-rebuild grace period');
    assert.dom(GENERAL.infoRowValue('Delta CRL building')).hasText('Off', 'it renders falsy delta cr build');
    assert.dom(SELECTORS.rowIcon('Delta CRL building', 'x-square'));
    assert
      .dom(GENERAL.infoRowValue('Delta rebuild interval'))
      .doesNotExist('does not render delta rebuild duration');

    // check falsy disable and ocsp_disable hides duration values and other params
    this.crl.auto_rebuild = true;
    this.crl.enable_delta = true;
    this.crl.disable = true;
    this.crl.ocsp_disable = true;
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('CRL building')).hasText('Disabled', 'disabled renders');
    assert.dom(GENERAL.infoRowValue('Expiry')).doesNotExist();
    assert
      .dom(GENERAL.infoRowValue('Responder APIs'))
      .hasText('Disabled', 'responder apis value renders Disabled');
    assert.dom(GENERAL.infoRowValue('Interval')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Auto-rebuild')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Auto-rebuild grace period')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Delta CRL building')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Delta rebuild interval')).doesNotExist();
  });

  test('it renders enterprise params in crl section', async function (assert) {
    this.owner.lookup('service:version').type = 'enterprise';
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('Cross-cluster revocation')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Cross-cluster revocation', 'check-circle'));
    assert.dom(GENERAL.infoRowValue('Unified CRL')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Unified CRL', 'check-circle'));
    assert.dom(GENERAL.infoRowValue('Unified CRL on existing paths')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Unified CRL on existing paths', 'check-circle'));
  });

  test('it does not render enterprise params in crl section', async function (assert) {
    this.owner.lookup('service:version').type = 'community';
    await this.renderComponent();

    assert.dom(GENERAL.infoRowValue('Cross-cluster revocation')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Unified CRL')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Unified CRL on existing paths')).doesNotExist();
  });
});
