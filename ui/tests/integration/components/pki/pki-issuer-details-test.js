/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-issuer-details';

module('Integration | Component | page/pki-issuer-details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.context = { owner: this.engine };
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.issuer = this.store.createRecord('pki/issuer', { issuerId: 'abcd-efgh' });
  });

  test('it renders with correct toolbar by default', async function (assert) {
    await render(
      hbs`
      <Page::PkiIssuerDetails @issuer={{this.issuer}} />
            `,
      this.context
    );

    assert.dom(SELECTORS.rotateRoot).doesNotExist();
    assert.dom(SELECTORS.crossSign).doesNotExist();
    assert.dom(SELECTORS.signIntermediate).doesNotExist();
    assert.dom(SELECTORS.download).hasText('Download');
    assert.dom(SELECTORS.configure).doesNotExist();
    assert.dom(SELECTORS.parsingAlertBanner).doesNotExist();
  });

  test('it renders toolbar actions depending on passed capabilities', async function (assert) {
    this.set('isRotatable', true);
    this.set('canRotate', true);
    this.set('canCrossSign', true);
    this.set('canSignIntermediate', true);
    this.set('canConfigure', true);

    await render(
      hbs`
      <Page::PkiIssuerDetails
        @issuer={{this.issuer}}
        @isRotatable={{this.isRotatable}}
        @canRotate={{this.canRotate}}
        @canCrossSign={{this.canCrossSign}}
        @canSignIntermediate={{this.canSignIntermediate}}
        @canConfigure={{this.canConfigure}}
      />
            `,
      this.context
    );

    assert.dom(SELECTORS.parsingAlertBanner).doesNotExist();
    assert.dom(SELECTORS.rotateRoot).hasText('Rotate this root');
    assert.dom(SELECTORS.crossSign).hasText('Cross-sign issuers');
    assert.dom(SELECTORS.signIntermediate).hasText('Sign Intermediate');
    assert.dom(SELECTORS.download).hasText('Download');
    assert.dom(SELECTORS.configure).hasText('Configure');

    this.set('canRotate', false);
    this.set('canCrossSign', false);
    this.set('canSignIntermediate', false);
    this.set('canConfigure', false);
    await settled();

    assert.dom(SELECTORS.rotateRoot).doesNotExist();
    assert.dom(SELECTORS.crossSign).doesNotExist();
    assert.dom(SELECTORS.signIntermediate).doesNotExist();
    assert.dom(SELECTORS.download).hasText('Download');
    assert.dom(SELECTORS.configure).doesNotExist();
  });

  test('it renders correct details by default', async function (assert) {
    await render(
      hbs`
        <Page::PkiIssuerDetails @issuer={{this.issuer}} />
                `,
      this.context
    );

    // Default group details:
    assert.dom(SELECTORS.defaultGroup).exists('Default group of details exists');
    assert.dom(SELECTORS.valueByName('Certificate')).exists('Certificate detail exists');
    assert.dom(SELECTORS.copyButtonByName('Certificate')).exists('Certificate is copyable');
    assert.dom(SELECTORS.valueByName('CA Chain')).exists('CA Chain detail exists');
    assert.dom(SELECTORS.copyButtonByName('CA Chain')).exists('CA Chain is copyable');
    assert.dom(SELECTORS.valueByName('Common name')).exists('Common name detail exists');
    assert.dom(SELECTORS.valueByName('Issuer name')).exists('Issuer name detail exists');
    assert.dom(SELECTORS.valueByName('Issuer ID')).exists('Issuer ID detail exists');
    assert.dom(SELECTORS.copyButtonByName('Issuer ID')).exists('Issuer ID is copyable');
    assert.dom(SELECTORS.valueByName('Default key ID')).exists('Default key ID detail exists');

    // Issuer URLs group details:
    assert.dom(SELECTORS.urlsGroup).exists('Issuer URLs group of details exists');
    assert.dom(SELECTORS.valueByName('Issuing certificates')).exists('Issuing certificates detail exists');
    assert
      .dom(SELECTORS.valueByName('CRL distribution points'))
      .exists('CRL distribution points detail exists');
    assert.dom(SELECTORS.valueByName('OCSP servers')).exists('OCSP servers detail exists');
  });

  test('it renders parsing error banner if issuer certificate contains unsupported OIDs', async function (assert) {
    this.issuer.parsedCertificate = {
      common_name: 'fancy-cert-unsupported-subj-and-ext-oids',
      subject_serial_number: null,
      ou: null,
      organization: 'Acme, Inc',
      country: 'US',
      locality: 'Topeka',
      province: 'Kansas',
      street_address: null,
      parsing_errors: [new Error('certificate contains stuff we cannot parse')],
      can_parse: true,
    };
    await render(
      hbs`
      <Page::PkiIssuerDetails @issuer={{this.issuer}} />
            `,
      this.context
    );

    assert.dom(SELECTORS.parsingAlertBanner).exists();
    assert
      .dom(SELECTORS.parsingAlertBanner)
      .hasText(
        "There was an error parsing certificate metadata Vault cannot display unparsed values, but this will not interfere with the certificate's functionality. However, if you wish to cross-sign this issuer it must be done manually using the CLI. Parsing error(s): certificate contains stuff we cannot parse"
      );
  });

  test('it renders parsing error banner if can_parse=false but no parsing_errors', async function (assert) {
    this.issuer.parsedCertificate = {
      common_name: 'fancy-cert-unsupported-subj-and-ext-oids',
      subject_serial_number: null,
      ou: null,
      organization: 'Acme, Inc',
      country: 'US',
      locality: 'Topeka',
      province: 'Kansas',
      street_address: null,
      parsing_errors: [],
      can_parse: false,
    };
    await render(
      hbs`
      <Page::PkiIssuerDetails @issuer={{this.issuer}} />
            `,
      this.context
    );

    assert.dom(SELECTORS.parsingAlertBanner).exists();
    assert
      .dom(SELECTORS.parsingAlertBanner)
      .hasText(
        "There was an error parsing certificate metadata Vault cannot display unparsed values, but this will not interfere with the certificate's functionality. However, if you wish to cross-sign this issuer it must be done manually using the CLI."
      );
  });

  test('it renders parsing error banner if no key for parsing_errors', async function (assert) {
    this.issuer.parsedCertificate = {
      common_name: 'fancy-cert-unsupported-subj-and-ext-oids',
      subject_serial_number: null,
      ou: null,
      organization: 'Acme, Inc',
      country: 'US',
      locality: 'Topeka',
      province: 'Kansas',
      street_address: null,
      can_parse: false,
    };

    await render(
      hbs`
      <Page::PkiIssuerDetails @issuer={{this.issuer}} />
            `,
      this.context
    );

    assert.dom(SELECTORS.parsingAlertBanner).exists();
    assert
      .dom(SELECTORS.parsingAlertBanner)
      .hasText(
        "There was an error parsing certificate metadata Vault cannot display unparsed values, but this will not interfere with the certificate's functionality. However, if you wish to cross-sign this issuer it must be done manually using the CLI."
      );
  });
});
