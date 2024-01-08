/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';

const SELECTORS = {
  rowLabel: (attr) => `[data-test-row-label="${attr}"]`,
  rowValue: (attr) => `[data-test-value-div="${attr}"]`,
  rowIcon: (attr, icon) => `[data-test-row-value="${attr}"] [data-test-icon="${icon}"]`,
};

module('Integration | Component | Page::PkiConfigurationDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';

    this.store = this.owner.lookup('service:store');
    this.cluster = this.store.createRecord('pki/config/cluster', {
      id: 'pki-test',
      path: 'https://pr-a.vault.example.com/v1/ns1/pki-root',
    });
    this.urls = this.store.createRecord('pki/config/urls', {
      id: 'pki-test',
      issuingCertificates: 'example.com',
    });
    this.crl = this.store.createRecord('pki/config/crl', {
      id: 'pki-test',
      expiry: '20h',
      disable: false,
      autoRebuild: true,
      autoRebuildGracePeriod: '13h',
      enableDelta: true,
      deltaRebuildInterval: '15m',
      ocspExpiry: '77h',
      ocspDisable: false,
      crossClusterRevocation: true,
      unifiedCrl: true,
      unifiedCrlOnExistingPaths: true,
    });
    this.mountConfig = {
      id: 'pki-test',
      engineType: 'pki',
      path: '/pki-test',
      accessor: 'pki_33345b0d',
      local: false,
      sealWrap: true,
      config: this.store.createRecord('mount-config', {
        defaultLease: '12h',
        maxLeaseTtl: '400h',
        allowedManagedKeys: true,
      }),
    };
  });

  test('shows the correct information on cluster config', async function (assert) {
    await render(hbs`<Page::PkiConfigurationDetails @cluster={{this.cluster}} @hasConfig={{true}} />,`, {
      owner: this.engine,
    });

    assert
      .dom(SELECTORS.rowValue("Mount's API path"))
      .hasText('https://pr-a.vault.example.com/v1/ns1/pki-root', 'mount API path row renders');
    assert.dom(SELECTORS.rowValue('AIA path')).hasText('None', "renders 'None' when no data");
  });

  test('shows the correct information on global urls section', async function (assert) {
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

    assert
      .dom(SELECTORS.rowLabel('Issuing certificates'))
      .hasText('Issuing certificates', 'issuing certificate row label renders');
    assert
      .dom(SELECTORS.rowValue('Issuing certificates'))
      .hasText('example.com', 'issuing certificate value renders');
    this.urls.issuingCertificates = null;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert
      .dom(SELECTORS.rowValue('Issuing certificates'))
      .hasText('None', 'issuing certificate value renders None if none is configured');
    assert
      .dom(SELECTORS.rowLabel('CRL distribution points'))
      .hasText('CRL distribution points', 'crl distribution points row label renders');
    assert
      .dom(SELECTORS.rowValue('CRL distribution points'))
      .hasText('None', 'crl distribution points value renders None if none is configured');
  });

  test('shows the correct information on crl section', async function (assert) {
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.rowLabel('CRL building')).hasText('CRL building', 'crl expiry row label renders');
    assert.dom(SELECTORS.rowValue('CRL building')).hasText('Enabled', 'enabled renders');
    assert.dom(SELECTORS.rowValue('Expiry')).hasText('20h', 'expiry value renders');
    assert.dom(SELECTORS.rowLabel('Auto-rebuild')).hasText('Auto-rebuild', 'auto rebuild label renders');
    assert.dom(SELECTORS.rowValue('Auto-rebuild')).hasText('On', 'it renders truthy auto build');
    assert.dom(SELECTORS.rowIcon('Auto-rebuild', 'check-circle'));
    assert
      .dom(SELECTORS.rowValue('Auto-rebuild grace period'))
      .hasText('13h', 'it renders auto build grace period');
    assert.dom(SELECTORS.rowValue('Delta CRL building')).hasText('On', 'it renders truthy delta crl build');
    assert.dom(SELECTORS.rowIcon('Delta CRL building', 'check-circle'));
    assert
      .dom(SELECTORS.rowValue('Delta rebuild interval'))
      .hasText('15m', 'it renders delta build duration');
    assert
      .dom(SELECTORS.rowValue('Responder APIs'))
      .hasText('Enabled', 'responder apis value renders Enabled if ocsp_disable=false');
    assert.dom(SELECTORS.rowValue('Interval')).hasText('77h', 'interval value renders');
    // check falsy aut_rebuild and _enable_delta hides duration values
    this.crl.autoRebuild = false;
    this.crl.enableDelta = false;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.rowValue('Auto-rebuild')).hasText('Off', 'it renders falsy auto build');
    assert.dom(SELECTORS.rowIcon('Auto-rebuild', 'x-square'));
    assert
      .dom(SELECTORS.rowValue('Auto-rebuild grace period'))
      .doesNotExist('does not render auto-rebuild grace period');
    assert.dom(SELECTORS.rowValue('Delta CRL building')).hasText('Off', 'it renders falsy delta cr build');
    assert.dom(SELECTORS.rowIcon('Delta CRL building', 'x-square'));
    assert
      .dom(SELECTORS.rowValue('Delta rebuild interval'))
      .doesNotExist('does not render delta rebuild duration');

    // check falsy disable and ocsp_disable hides duration values and other params
    this.crl.autoRebuild = true;
    this.crl.enableDelta = true;
    this.crl.disable = true;
    this.crl.ocspDisable = true;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.rowValue('CRL building')).hasText('Disabled', 'disabled renders');
    assert.dom(SELECTORS.rowValue('Expiry')).doesNotExist();
    assert
      .dom(SELECTORS.rowValue('Responder APIs'))
      .hasText('Disabled', 'responder apis value renders Disabled');
    assert.dom(SELECTORS.rowValue('Interval')).doesNotExist();
    assert.dom(SELECTORS.rowValue('Auto-rebuild')).doesNotExist();
    assert.dom(SELECTORS.rowValue('Auto-rebuild grace period')).doesNotExist();
    assert.dom(SELECTORS.rowValue('Delta CRL building')).doesNotExist();
    assert.dom(SELECTORS.rowValue('Delta rebuild interval')).doesNotExist();
  });

  test('it renders enterprise params in crl section', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1+ent';
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.rowValue('Cross-cluster revocation')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Cross-cluster revocation', 'check-circle'));
    assert.dom(SELECTORS.rowValue('Unified CRL')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Unified CRL', 'check-circle'));
    assert.dom(SELECTORS.rowValue('Unified CRL on existing paths')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Unified CRL on existing paths', 'check-circle'));
  });

  test('it does not render enterprise params in crl section', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.13.1';
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.rowValue('Cross-cluster revocation')).doesNotExist();
    assert.dom(SELECTORS.rowValue('Unified CRL')).doesNotExist();
    assert.dom(SELECTORS.rowValue('Unified CRL on existing paths')).doesNotExist();
  });

  test('shows the correct information on mount configuration section', async function (assert) {
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

    assert
      .dom(SELECTORS.rowLabel('Secret engine type'))
      .hasText('Secret engine type', 'engine type row label renders');
    assert.dom(SELECTORS.rowValue('Secret engine type')).hasText('pki', 'engine type row value renders');
    assert.dom(SELECTORS.rowLabel('Path')).hasText('Path', 'path row label renders');
    assert.dom(SELECTORS.rowValue('Path')).hasText('/pki-test', 'path row value renders');
    assert.dom(SELECTORS.rowLabel('Accessor')).hasText('Accessor', 'accessor row label renders');
    assert.dom(SELECTORS.rowValue('Accessor')).hasText('pki_33345b0d', 'accessor row value renders');
    assert.dom(SELECTORS.rowLabel('Local')).hasText('Local', 'local row label renders');
    assert.dom(SELECTORS.rowValue('Local')).hasText('No', 'local row value renders');
    assert.dom(SELECTORS.rowLabel('Seal wrap')).hasText('Seal wrap', 'seal wrap row label renders');
    assert
      .dom(SELECTORS.rowValue('Seal wrap'))
      .hasText('Yes', 'seal wrap row value renders Yes if sealWrap is true');
    assert.dom(SELECTORS.rowLabel('Max lease TTL')).hasText('Max lease TTL', 'max lease label renders');
    assert.dom(SELECTORS.rowValue('Max lease TTL')).hasText('400h', 'max lease value renders');
    assert
      .dom(SELECTORS.rowLabel('Allowed managed keys'))
      .hasText('Allowed managed keys', 'allowed managed keys label renders');
    assert
      .dom(SELECTORS.rowValue('Allowed managed keys'))
      .hasText('Yes', 'allowed managed keys value renders');
  });

  test('shows mount configuration when hasConfig is false', async function (assert) {
    this.urls = 403;
    this.crl = 403;

    await render(
      hbs`<Page::PkiConfigurationDetails @mountConfig={{this.mountConfig}} @hasConfig={{false}} />,`,
      { owner: this.engine }
    );

    assert
      .dom(SELECTORS.rowLabel('Secret engine type'))
      .hasText('Secret engine type', 'engine type row label renders');
    assert.dom(SELECTORS.rowValue('Secret engine type')).hasText('pki', 'engine type row value renders');
    assert.dom(SELECTORS.rowLabel('Path')).hasText('Path', 'path row label renders');
    assert.dom(SELECTORS.rowValue('Path')).hasText('/pki-test', 'path row value renders');
    assert.dom(SELECTORS.rowLabel('Accessor')).hasText('Accessor', 'accessor row label renders');
    assert.dom(SELECTORS.rowValue('Accessor')).hasText('pki_33345b0d', 'accessor row value renders');
    assert.dom(SELECTORS.rowLabel('Local')).hasText('Local', 'local row label renders');
    assert.dom(SELECTORS.rowValue('Local')).hasText('No', 'local row value renders');
    assert.dom(SELECTORS.rowLabel('Seal wrap')).hasText('Seal wrap', 'seal wrap row label renders');
    assert
      .dom(SELECTORS.rowValue('Seal wrap'))
      .hasText('Yes', 'seal wrap row value renders Yes if sealWrap is true');
    assert.dom(SELECTORS.rowLabel('Max lease TTL')).hasText('Max lease TTL', 'max lease label renders');
    assert.dom(SELECTORS.rowValue('Max lease TTL')).hasText('400h', 'max lease value renders');
    assert
      .dom(SELECTORS.rowLabel('Allowed managed keys'))
      .hasText('Allowed managed keys', 'allowed managed keys label renders');
    assert
      .dom(SELECTORS.rowValue('Allowed managed keys'))
      .hasText('Yes', 'allowed managed keys value renders');
  });
});
