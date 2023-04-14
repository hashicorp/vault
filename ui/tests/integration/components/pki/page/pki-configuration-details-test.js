/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-configuration-details';

module('Integration | Component | Page::PkiConfigurationDetails', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';

    this.store = this.owner.lookup('service:store');
    this.urls = this.store.createRecord('pki/urls', { id: 'pki-test', issuingCertificates: 'example.com' });
    this.crl = this.store.createRecord('pki/crl', {
      id: 'pki-test',
      expiry: '20h',
      autoRebuild: false,
      ocspExpiry: '77h',
      oscpDisable: true,
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

  test('shows the correct information on global urls section', async function (assert) {
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

    assert
      .dom(SELECTORS.issuingCertificatesLabel)
      .hasText('Issuing certificates', 'issuing certificate row label renders');
    assert
      .dom(SELECTORS.issuingCertificatesRowVal)
      .hasText('example.com', 'issuing certificate value renders');
    this.urls.issuingCertificates = null;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert
      .dom(SELECTORS.issuingCertificatesRowVal)
      .hasText('None', 'issuing certificate value renders None if none is configured');
    assert
      .dom(SELECTORS.crlDistributionPointsLabel)
      .hasText('CRL distribution points', 'crl distribution points row label renders');
    assert
      .dom(SELECTORS.crlDistributionPointsRowVal)
      .hasText('None', 'crl distribution points value renders None if none is configured');
  });

  test('shows the correct information on crl section', async function (assert) {
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.expiryLabel).hasText('Expiry', 'crl expiry row label renders');
    assert.dom(SELECTORS.expiryRowVal).hasText('20h', 'expiry value renders');
    assert.dom(SELECTORS.rebuildLabel).hasText('Auto-rebuild', 'auto rebuild label renders');
    assert
      .dom(SELECTORS.rebuildRowVal)
      .hasText('Off', 'auto-rebuild value renders off if auto rebuild is false');
    this.crl.autoRebuild = true;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert
      .dom(SELECTORS.rebuildRowVal)
      .hasText('On', 'auto-rebuild value renders on if auto rebuild is true');
    assert.dom(SELECTORS.responderApiLabel).hasText('Responder APIs', 'responder apis row label renders');
    assert
      .dom(SELECTORS.responderApiRowVal)
      .hasText('Enabled', 'responder apis value renders Enabled if oscpDisable is true');
    assert.dom(SELECTORS.intervalLabel).hasText('Interval', 'interval row label renders');
    assert.dom(SELECTORS.intervalRowVal).hasText('77h', 'interval value renders');
  });

  test('shows the correct information on mount configuration section', async function (assert) {
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @mountConfig={{this.mountConfig}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.engineTypeLabel).hasText('Secret engine type', 'engine type row label renders');
    assert.dom(SELECTORS.engineTypeRowVal).hasText('pki', 'engine type row value renders');
    assert.dom(SELECTORS.pathLabel).hasText('Path', 'path row label renders');
    assert.dom(SELECTORS.pathRowVal).hasText('/pki-test', 'path row value renders');
    assert.dom(SELECTORS.accessorLabel).hasText('Accessor', 'accessor row label renders');
    assert.dom(SELECTORS.accessorRowVal).hasText('pki_33345b0d', 'accessor row value renders');
    assert.dom(SELECTORS.localLabel).hasText('Local', 'local row label renders');
    assert.dom(SELECTORS.localRowVal).hasText('No', 'local row value renders');
    assert.dom(SELECTORS.sealWrapLabel).hasText('Seal wrap', 'seal wrap row label renders');
    assert
      .dom(SELECTORS.sealWrapRowVal)
      .hasText('Yes', 'seal wrap row value renders Yes if sealWrap is true');
    assert.dom(SELECTORS.maxLeaseTtlLabel).hasText('Max lease TTL', 'max lease label renders');
    assert.dom(SELECTORS.maxLeaseTtlRowVal).hasText('400h', 'max lease value renders');
    assert
      .dom(SELECTORS.allowedManagedKeysLabel)
      .hasText('Allowed managed keys', 'allowed managed keys label renders');
    assert.dom(SELECTORS.allowedManagedKeysRowVal).hasText('Yes', 'allowed managed keys value renders');
  });

  test('shows mount configuration when hasConfig is false', async function (assert) {
    this.urls = 403;
    this.crl = 403;

    await render(
      hbs`<Page::PkiConfigurationDetails @mountConfig={{this.mountConfig}} @hasConfig={{false}} />,`,
      { owner: this.engine }
    );

    assert.dom(SELECTORS.engineTypeLabel).hasText('Secret engine type', 'engine type row label renders');
    assert.dom(SELECTORS.engineTypeRowVal).hasText('pki', 'engine type row value renders');
    assert.dom(SELECTORS.pathLabel).hasText('Path', 'path row label renders');
    assert.dom(SELECTORS.pathRowVal).hasText('/pki-test', 'path row value renders');
    assert.dom(SELECTORS.accessorLabel).hasText('Accessor', 'accessor row label renders');
    assert.dom(SELECTORS.accessorRowVal).hasText('pki_33345b0d', 'accessor row value renders');
    assert.dom(SELECTORS.localLabel).hasText('Local', 'local row label renders');
    assert.dom(SELECTORS.localRowVal).hasText('No', 'local row value renders');
    assert.dom(SELECTORS.sealWrapLabel).hasText('Seal wrap', 'seal wrap row label renders');
    assert
      .dom(SELECTORS.sealWrapRowVal)
      .hasText('Yes', 'seal wrap row value renders Yes if sealWrap is true');
    assert.dom(SELECTORS.maxLeaseTtlLabel).hasText('Max lease TTL', 'max lease label renders');
    assert.dom(SELECTORS.maxLeaseTtlRowVal).hasText('400h', 'max lease value renders');
    assert
      .dom(SELECTORS.allowedManagedKeysLabel)
      .hasText('Allowed managed keys', 'allowed managed keys label renders');
    assert.dom(SELECTORS.allowedManagedKeysRowVal).hasText('Yes', 'allowed managed keys value renders');
  });
});
