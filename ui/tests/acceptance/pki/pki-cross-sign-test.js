/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, click, fillIn, find } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-issuer-cross-sign';
import { verifyCertificates } from 'vault/utils/parse-pki-cert';
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

    await runCommands([
      `write "${this.parentMountPath}/root/generate/internal" common_name="Long-Lived Root X1" ttl=8960h issuer_name="${this.oldParentIssuerName}"`,
      `write "${this.parentMountPath}/root/generate/internal" common_name="Long-Lived Root X2" ttl=8960h issuer_name="${this.parentIssuerName}"`,
      `write "${this.parentMountPath}/config/issuers" default="${this.parentIssuerName}"`,
    ]);
  });

  hooks.afterEach(async function () {
    // Cleanup engine
    await runCommands([`delete sys/mounts/${this.intMountPath}`]);
    await runCommands([`delete sys/mounts/${this.parentMountPath}`]);
  });

  test('it cross-signs an issuer', async function (assert) {
    // configure parent and intermediate mounts to make them cross-signable
    await visit(`/vault/secrets/${this.intMountPath}/pki/configuration/create`);
    await click(SELECTORS.configure.optionByKey('generate-csr'));
    await fillIn(SELECTORS.inputByName('type'), 'internal');
    await fillIn(SELECTORS.inputByName('commonName'), 'Short-Lived Int R1');
    await click('[data-test-save]');
    const csr = find(SELECTORS.copyButton('CSR')).getAttribute('data-clipboard-text');
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.oldParentIssuerName}/sign`);
    await fillIn(SELECTORS.inputByName('csr'), csr);
    await fillIn(SELECTORS.inputByName('format'), 'pem_bundle');
    await click('[data-test-pki-sign-intermediate-save]');
    const pemBundle = find(SELECTORS.copyButton('CA Chain'))
      .getAttribute('data-clipboard-text')
      .replace(/,/, '\n');
    await visit(`vault/secrets/${this.intMountPath}/pki/configuration/create`);
    await click(SELECTORS.configure.optionByKey('import'));
    await click('[data-test-text-toggle]');
    await fillIn('[data-test-text-file-textarea]', pemBundle);
    await click(SELECTORS.configure.importSubmit);
    await visit(`vault/secrets/${this.intMountPath}/pki/issuers`);
    await click('[data-test-is-default]');
    // name default issuer of intermediate
    const oldIntIssuerId = find(SELECTORS.rowValue('Issuer ID')).innerText;
    const oldIntCert = find(SELECTORS.copyButton('Certificate')).getAttribute('data-clipboard-text');
    await click(SELECTORS.details.configure);
    await fillIn(SELECTORS.inputByName('issuerName'), this.intIssuerName);
    await click('[data-test-save]');

    // perform cross-sign
    await visit(`vault/secrets/${this.parentMountPath}/pki/issuers/${this.parentIssuerName}/cross-sign`);
    await fillIn(SELECTORS.objectListInput('intermediateMount'), this.intMountPath);
    await fillIn(SELECTORS.objectListInput('intermediateIssuer'), this.intIssuerName);
    await fillIn(SELECTORS.objectListInput('newCrossSignedIssuer'), this.newlySignedIssuer);
    await click(SELECTORS.submitButton);
    assert
      .dom(`${SELECTORS.signedIssuerCol('intermediateMount')} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.intMountPath}/pki/overview`);
    assert
      .dom(`${SELECTORS.signedIssuerCol('intermediateIssuer')} a`)
      .hasAttribute('href', `/ui/vault/secrets/${this.intMountPath}/pki/issuers/${oldIntIssuerId}/details`);

    // get certificate data of newly signed issuer
    await click(`${SELECTORS.signedIssuerCol('newCrossSignedIssuer')} a`);
    const newIntCert = find(SELECTORS.copyButton('Certificate')).getAttribute('data-clipboard-text');

    // verify cross-sign was accurate by creating a role to issue a leaf certificate
    const myRole = 'some-role';
    await runCommands([
      `write ${this.intMountPath}/roles/${myRole} \
    issuer_ref=${this.newlySignedIssuer}\
    allow_any_name=true \
    max_ttl="720h"`,
    ]);
    await visit(`vault/secrets/${this.intMountPath}/pki/roles/${myRole}/generate`);
    await fillIn(SELECTORS.inputByName('commonName'), 'my-leaf');
    await fillIn('[data-test-ttl-value="TTL"]', '3600');
    await click('[data-test-pki-generate-button]');
    const myLeafCert = find(SELECTORS.copyButton('Certificate')).getAttribute('data-clipboard-text');

    // see comments in utils/parse-pki-cert.js for step-by-step explanation of of verifyCertificates method
    assert.true(
      await verifyCertificates(oldIntCert, newIntCert, myLeafCert),
      'the leaf certificate validates against both intermediate certificates'
    );
  });
});
