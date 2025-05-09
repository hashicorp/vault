/**
 * Copyright (c) HashiCorp, Inc.
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
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('shows the correct information on cluster config', async function (assert) {
    await render(hbs`<Page::PkiConfigurationDetails @cluster={{this.cluster}} @hasConfig={{true}} />,`, {
      owner: this.engine,
    });
    assert
      .dom(GENERAL.infoRowValue("Mount's API path"))
      .hasText('https://pr-a.vault.example.com/v1/ns1/pki-root', 'mount API path row renders');
    assert.dom(GENERAL.infoRowValue('AIA path')).hasText('None', "renders 'None' when no data");
  });

  test('shows the correct information on global urls section', async function (assert) {
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

    assert
      .dom(GENERAL.infoRowLabel('Issuing certificates'))
      .hasText('Issuing certificates', 'issuing certificate row label renders');
    assert
      .dom(GENERAL.infoRowValue('Issuing certificates'))
      .hasText('example.com', 'issuing certificate value renders');
    this.urls.issuingCertificates = null;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
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
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );

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
    this.crl.autoRebuild = false;
    this.crl.enableDelta = false;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
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
    this.crl.autoRebuild = true;
    this.crl.enableDelta = true;
    this.crl.disable = true;
    this.crl.ocspDisable = true;
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
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
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert.dom(GENERAL.infoRowValue('Cross-cluster revocation')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Cross-cluster revocation', 'check-circle'));
    assert.dom(GENERAL.infoRowValue('Unified CRL')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Unified CRL', 'check-circle'));
    assert.dom(GENERAL.infoRowValue('Unified CRL on existing paths')).hasText('Yes');
    assert.dom(SELECTORS.rowIcon('Unified CRL on existing paths', 'check-circle'));
  });

  test('it does not render enterprise params in crl section', async function (assert) {
    this.version = this.owner.lookup('service:version');
    this.version.type = 'community';
    await render(
      hbs`<Page::PkiConfigurationDetails @urls={{this.urls}} @crl={{this.crl}} @hasConfig={{true}} />,`,
      { owner: this.engine }
    );
    assert.dom(GENERAL.infoRowValue('Cross-cluster revocation')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Unified CRL')).doesNotExist();
    assert.dom(GENERAL.infoRowValue('Unified CRL on existing paths')).doesNotExist();
  });
});
