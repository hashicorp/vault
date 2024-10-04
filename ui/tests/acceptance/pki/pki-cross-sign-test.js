/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { visit, click, fillIn, find } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCmd } from 'vault/tests/helpers/commands';
import { verifyCertificates } from 'vault/utils/parse-pki-cert';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import {
  PKI_CONFIGURE_CREATE,
  PKI_CROSS_SIGN,
  PKI_ISSUER_DETAILS,
} from 'vault/tests/helpers/pki/pki-selectors';

module('Acceptance | pki/pki cross sign', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await authPage.login();
    this.parentMountPath = `parent-mount-${uuidv4()}`;
    this.oldParentIssuerName = 'old-parent-issuer'; // old parent issuer we're transferring from
    this.parentIssuerName = 'new-parent-issuer'; // issuer where cross-signing action will begin
    this.intMountPath = `intermediate-mount-${uuidv4()}`; // first input box in cross-signing page
    this.intIssuerName = 'my-intermediate-issuer'; // second input box in cross-signing page
    this.newlySignedIssuer = 'my-newly-signed-int'; // third input
    await enablePage.enable('pki', this.parentMountPath);
    await enablePage.enable('pki', this.intMountPath);

    await runCmd([
      `write "${this.parentMountPath}/root/generate/internal" common_name="Long-Lived Root X1" ttl=8960h issuer_name="${this.oldParentIssuerName}"`,
      `write "${this.parentMountPath}/root/generate/internal" common_name="Long-Lived Root X2" ttl=8960h issuer_name="${this.parentIssuerName}"`,
      `write "${this.parentMountPath}/config/issuers" default="${this.parentIssuerName}"`,
    ]);
  });

  hooks.afterEach(async function () {
    // Cleanup engine
    await runCmd([`delete sys/mounts/${this.intMountPath}`]);
    await runCmd([`delete sys/mounts/${this.parentMountPath}`]);
  });

  test('it cross-signs an issuer', async function (assert) {
    // configure parent and intermediate mounts to make them cross-signable
    await visit(`/vault/secrets/${this.intMountPath}/pki/configuration/create`);
    await click(PKI_CONFIGURE_CREATE.optionByKey('generate-csr'));
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    await fillIn(GENERAL.inputByAttr('commonName'), 'Short-Lived Int R1');
    await click('[data-test-save]');
    const csr = find(PKI_CROSS_SIGN.copyButton('CSR')).getAttribute('data-test-copy-button');
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.oldParentIssuerName}/sign`);
    await fillIn(GENERAL.inputByAttr('csr'), csr);
    await fillIn(GENERAL.inputByAttr('format'), 'pem_bundle');
    await click('[data-test-pki-sign-intermediate-save]');
    const pemBundle = find(PKI_CROSS_SIGN.copyButton('CA Chain'))
      .getAttribute('data-test-copy-button')
      .replace(/,/, '\n');
    await visit(`vault/secrets/${this.intMountPath}/pki/configuration/create`);
    await click(PKI_CONFIGURE_CREATE.optionByKey('import'));
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', pemBundle);
    await click(PKI_CONFIGURE_CREATE.importSubmit);
    await visit(`vault/secrets/${this.intMountPath}/pki/issuers`);
    await click('[data-test-is-default]');
    // name default issuer of intermediate
    const oldIntIssuerId = find(PKI_CROSS_SIGN.rowValue('Issuer ID')).innerText;
    const oldIntCert = find(PKI_CROSS_SIGN.copyButton('Certificate')).getAttribute('data-test-copy-button');
    await click(PKI_ISSUER_DETAILS.configure);
    await fillIn(GENERAL.inputByAttr('issuerName'), this.intIssuerName);
    await click('[data-test-save]');

    // perform cross-sign
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.parentIssuerName}/cross-sign`);
    await fillIn(PKI_CROSS_SIGN.objectListInput('intermediateMount'), this.intMountPath);
    await fillIn(PKI_CROSS_SIGN.objectListInput('intermediateIssuer'), this.intIssuerName);
    await fillIn(PKI_CROSS_SIGN.objectListInput('newCrossSignedIssuer'), this.newlySignedIssuer);
    await click(GENERAL.saveButton);
    assert
      .dom(`${PKI_CROSS_SIGN.signedIssuerCol('intermediateMount')} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.intMountPath}/pki/overview`);
    assert
      .dom(`${PKI_CROSS_SIGN.signedIssuerCol('intermediateIssuer')} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.intMountPath}/pki/issuers/${oldIntIssuerId}/details`);

    // get certificate data of newly signed issuer
    await click(`${PKI_CROSS_SIGN.signedIssuerCol('newCrossSignedIssuer')} a`);
    const newIntCert = find(PKI_CROSS_SIGN.copyButton('Certificate')).getAttribute('data-test-copy-button');

    // verify cross-sign was accurate by creating a role to issue a leaf certificate
    const myRole = 'some-role';
    await runCmd([
      `write ${this.intMountPath}/roles/${myRole} \
    issuer_ref=${this.newlySignedIssuer}\
    allow_any_name=true \
    max_ttl="720h"`,
    ]);
    await visit(`vault/secrets/${this.intMountPath}/pki/roles/${myRole}/generate`);
    await fillIn(GENERAL.inputByAttr('commonName'), 'my-leaf');
    await fillIn('[data-test-ttl-value="TTL"]', '3600');
    await click(GENERAL.saveButton);
    const myLeafCert = find(PKI_CROSS_SIGN.copyButton('Certificate')).getAttribute('data-test-copy-button');

    // see comments in utils/parse-pki-cert.js for step-by-step explanation of of verifyCertificates method
    assert.true(
      await verifyCertificates(oldIntCert, newIntCert, myLeafCert),
      'the leaf certificate validates against both intermediate certificates'
    );
  });
});
