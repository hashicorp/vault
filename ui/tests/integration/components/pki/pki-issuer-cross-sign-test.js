/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { PKI_CROSS_SIGN } from 'vault/tests/helpers/pki/pki-selectors';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import sinon from 'sinon';

const FIELDS = [
  {
    label: 'Mount path',
    key: 'intermediateMount',
    placeholder: 'Mount path',
    helpText: 'The mount in which your new certificate can be found.',
  },
  {
    label: "Issuer's current name",
    key: 'intermediateIssuer',
    placeholder: 'Current issuer name',
    helpText: 'The API name of the previous intermediate which was cross-signed.',
  },
  {
    label: 'New issuer name',
    key: 'newCrossSignedIssuer',
    placeholder: 'Enter a new issuer name',
    helpText: `This is your new issuer’s name in the API.`,
  },
];
const { intIssuerCert, newCSR, newlySignedCert, oldParentIssuerCert, parentIssuerCert, unsupportedOids } =
  CERTIFICATES;
module('Integration | Component | pki issuer cross sign', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.backend = 'my-parent-issuer-mount';
    this.intMountPath = 'int-mount';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    // parent issuer
    this.parentIssuerData = {
      ca_chain: [parentIssuerCert],
      certificate: parentIssuerCert,
      crl_distribution_points: [],
      issuer_id: '0c983955-6426-22b2-1b3f-c0bdca40fd15',
      issuer_name: 'my-parent-issuer-name',
      issuing_certificates: [],
      key_id: '8b8d0017-a067-ac50-c5cf-475876f9aac5',
      leaf_not_after_behavior: 'err',
      manual_chain: null,
      ocsp_servers: [],
      revocation_signature_algorithm: 'SHA256WithRSA',
      revoked: false,
      usage: 'crl-signing,issuing-certificates,ocsp-signing,read-only',
    };
    // intermediate issuer
    this.intIssuerData = {
      ca_chain: [intIssuerCert, oldParentIssuerCert],
      certificate: intIssuerCert,
      crl_distribution_points: [],
      issuer_id: '6c286455-7904-5698-bf86-8aba81e680e6',
      issuer_name: 'source-int-name',
      issuing_certificates: [],
      key_id: '2e2b8baf-4dac-c46f-cee4-8afbc7f2d8b2',
      leaf_not_after_behavior: 'err',
      manual_chain: null,
      ocsp_servers: [],
      revocation_signature_algorithm: '',
      revoked: false,
      usage: 'crl-signing,issuing-certificates,ocsp-signing,read-only',
    };
    // newly cross signed issuer
    this.newIssuerData = {
      ca_chain: [newlySignedCert, parentIssuerCert],
      certificate: newlySignedCert,
      crl_distribution_points: [],
      issuer_id: 'bc159ba8-930c-c894-e871-2f3e889e8e02',
      issuer_name: 'newly-cross-signed-cert',
      issuing_certificates: [],
      key_id: '2e2b8baf-4dac-c46f-cee4-8afbc7f2d8b2',
      leaf_not_after_behavior: 'err',
      manual_chain: null,
      ocsp_servers: [],
      revocation_signature_algorithm: '',
      revoked: false,
      usage: 'crl-signing,issuing-certificates,ocsp-signing,read-only',
    };

    this.certData = {
      can_parse: true,
      common_name: 'Short-Lived Int R1',
      country: null,
      exclude_cn_from_sans: false,
      key_usage: 'CertSign, CRLSign',
      locality: null,
      max_path_length: undefined,
      not_valid_after: 1677371103,
      not_valid_before: 1674606273,
      organization: null,
      ou: null,
      parsing_errors: [],
      postal_code: null,
      province: null,
      signature_bits: '256',
      street_address: null,
      serial_number: null,
      ttl: '768h',
      use_pss: false,
    };

    this.testInputs = {
      intermediateMount: this.intMountPath,
      intermediateIssuer: this.intIssuerData.issuer_name,
      newCrossSignedIssuer: this.newIssuerData.issuer_name,
    };

    const api = this.owner.lookup('service:api');
    this.listStub = sinon.stub(api.secrets, 'pkiListIssuers').resolves({
      keys: [this.parentIssuerData.issuer_id],
      key_info: { [this.parentIssuerData.issuer_id]: this.parentIssuerData },
    });
    this.readStub = sinon.stub(api.secrets, 'pkiReadIssuer').resolves(this.intIssuerData);
    this.generateStub = sinon
      .stub(api.secrets, 'pkiGenerateIntermediate')
      .resolves({ csr: newCSR.csr, key_id: this.intIssuerData.key_id });
    this.signStub = sinon
      .stub(api.secrets, 'pkiIssuerSignIntermediate')
      .resolves({ ca_chain: [newlySignedCert, parentIssuerCert] });
    this.importStub = sinon
      .stub(api.secrets, 'pkiIssuersImportBundle')
      .resolves({ mapping: { [this.newIssuerData.issuer_id]: this.intIssuerData.key_id } });
    this.writeStub = sinon.stub(api.secrets, 'pkiWriteIssuer').resolves(this.newIssuerData);

    this.renderComponent = () =>
      render(hbs`<PkiIssuerCrossSign @parentIssuer={{this.parentIssuerData}} /> `, {
        owner: this.engine,
      });
  });

  test('it makes requests to the correct endpoints', async function (assert) {
    assert.expect(15);

    this.readStub.onSecondCall().resolves(this.newIssuerData);

    await this.renderComponent();
    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.submitButton);

    assert.true(
      this.listStub.calledWith(this.intMountPath),
      'Step 0. GET request is made to list issuers in intermediate mount for name validation'
    );

    assert.true(
      this.readStub.calledWith(this.intIssuerData.issuer_name, this.intMountPath),
      'Step 1. GET request is made to fetch existing issuer data'
    );

    assert.true(
      this.generateStub.calledWith('existing', this.intMountPath, {
        ...this.certData,
        key_ref: this.intIssuerData.key_id,
      }),
      'Step 2. POST request is made to generate new CSR'
    );

    assert.true(
      this.signStub.calledWith(this.parentIssuerData.issuer_name, this.backend, {
        ...newCSR,
        ...this.certData,
      }),
      'Step 3. POST request is made to sign CSR with new parent issuer'
    );

    assert.true(
      this.importStub.calledWith(this.intMountPath, {
        pem_bundle: [newlySignedCert, parentIssuerCert].join('\n'),
      }),
      'Step 4. POST request is made to import issuer'
    );

    assert.true(
      this.readStub.calledWith(this.newIssuerData.issuer_id, this.intMountPath),
      'Step 5. GET request is made to newly imported issuer'
    );

    assert.true(
      this.writeStub.calledWith(this.newIssuerData.issuer_id, this.intMountPath, {
        ...this.newIssuerData,
        issuer_name: this.newIssuerData.issuer_name,
      }),
      'Step 6. POST request is made to update issuer name'
    );

    assert.dom(PKI_CROSS_SIGN.statusCount).hasText('Cross-signing complete (1 successful, 0 errors)');
    assert
      .dom(`${PKI_CROSS_SIGN.signedIssuerRow()} [data-test-icon="check-circle"]`)
      .exists('row has success icon');
    for (const field of FIELDS) {
      assert
        .dom(`${PKI_CROSS_SIGN.signedIssuerCol(field.key)}`)
        .hasText(this.testInputs[field.key], `${field.key} displays correct value`);
      assert.dom(`${PKI_CROSS_SIGN.signedIssuerCol(field.key)} a`).hasTagName('a');
    }
  });

  test('it cross-signs multiple certs', async function (assert) {
    assert.expect(10);

    const nonexistentIssuer = {
      intermediateMount: this.intMountPath,
      intermediateIssuer: 'some-fake-issuer',
      newCrossSignedIssuer: 'failed-cert-1',
    };
    const unsupportedCert = {
      intermediateMount: this.intMountPath,
      intermediateIssuer: 'some-fancy-issuer',
      newCrossSignedIssuer: 'failed-cert-2',
    };

    const error = getErrorResponse(
      {
        errors: [
          `1 error occurred:\n\t* unable to find PKI issuer for reference: ${nonexistentIssuer.intermediateIssuer}\n\n`,
        ],
      },
      500
    );
    this.readStub.onCall(1).rejects(error);
    this.readStub
      .onCall(2)
      .resolves({ issuer_name: unsupportedCert.intermediateIssuer, certificate: unsupportedOids });
    this.readStub.onCall(3).rejects(error);

    await this.renderComponent();

    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(PKI_CROSS_SIGN.addRow);
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key, 1), nonexistentIssuer[field.key]);
    }
    await click(PKI_CROSS_SIGN.addRow);
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key, 2), unsupportedCert[field.key]);
    }
    await click(GENERAL.submitButton);

    assert.dom(PKI_CROSS_SIGN.statusCount).hasText('Cross-signing complete (0 successful, 3 errors)');
    for (const field of FIELDS) {
      assert
        .dom(`${PKI_CROSS_SIGN.signedIssuerRow()} ${PKI_CROSS_SIGN.signedIssuerCol(field.key)}`)
        .hasText(this.testInputs[field.key], `first row has correct values`);
    }
    for (const field of FIELDS) {
      assert
        .dom(`${PKI_CROSS_SIGN.signedIssuerRow(1)} ${PKI_CROSS_SIGN.signedIssuerCol(field.key)}`)
        .hasText(nonexistentIssuer[field.key], `second row has correct values`);
    }
    for (const field of FIELDS) {
      assert
        .dom(`${PKI_CROSS_SIGN.signedIssuerRow(2)} ${PKI_CROSS_SIGN.signedIssuerCol(field.key)}`)
        .hasText(unsupportedCert[field.key], `third row has correct values`);
    }
  });

  test('it returns API errors when a request fails', async function (assert) {
    assert.expect(7);

    this.readStub.rejects(
      getErrorResponse(
        { errors: ['1 error occurred:\n\t* unable to find PKI issuer for reference: nonexistent-mount\n\n'] },
        500
      )
    );

    await this.renderComponent();

    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.submitButton);

    assert.dom(PKI_CROSS_SIGN.statusCount).hasText('Cross-signing complete (0 successful, 1 error)');
    assert
      .dom(`${PKI_CROSS_SIGN.signedIssuerRow()} [data-test-icon="alert-circle-fill"]`)
      .exists('row has failure icon');

    assert.dom('[data-test-cross-sign-alert-title]').hasText('Cross-sign failed');
    assert
      .dom('[data-test-cross-sign-alert-message]')
      .hasText('1 error occurred: * unable to find PKI issuer for reference: nonexistent-mount');

    for (const field of FIELDS) {
      assert
        .dom(`${PKI_CROSS_SIGN.signedIssuerCol(field.key)}`)
        .hasText(this.testInputs[field.key], `${field.key} displays correct value`);
    }
  });

  test('it returns an error when a certificate contains unsupported values', async function (assert) {
    assert.expect(7);

    const unsupportedIssuerCert = { ...this.intIssuerData, certificate: unsupportedOids };
    this.readStub.resolves(unsupportedIssuerCert);

    await this.renderComponent();

    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.submitButton);
    assert.dom(PKI_CROSS_SIGN.statusCount).hasText('Cross-signing complete (0 successful, 1 error)');
    assert
      .dom(`${PKI_CROSS_SIGN.signedIssuerRow()} [data-test-icon="alert-circle-fill"]`)
      .exists('row has failure icon');
    assert
      .dom('[data-test-cross-sign-alert-title]')
      .hasText('Certificate must be manually cross-signed using the CLI.');
    assert
      .dom('[data-test-cross-sign-alert-message]')
      .hasText(
        'certificate contains unsupported subject OIDs: 1.2.840.113549.1.9.1, certificate contains unsupported extension OIDs: 2.5.29.37, unsupported key usage value on issuer certificate: DigitalSignature, KeyEncipherment'
      );

    for (const field of FIELDS) {
      assert
        .dom(`${PKI_CROSS_SIGN.signedIssuerCol(field.key)}`)
        .hasText(this.testInputs[field.key], `${field.key} displays correct value`);
    }
  });

  test('it returns an error when attempting to self-cross-sign', async function (assert) {
    assert.expect(7);

    this.testInputs = {
      intermediateMount: this.backend,
      intermediateIssuer: this.parentIssuerData.issuer_name,
      newCrossSignedIssuer: this.newIssuerData.issuer_name,
    };

    this.readStub.resolves(this.parentIssuerData);

    await this.renderComponent();

    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.submitButton);
    assert.dom(PKI_CROSS_SIGN.statusCount).hasText('Cross-signing complete (0 successful, 1 error)');
    assert
      .dom(`${PKI_CROSS_SIGN.signedIssuerRow()} [data-test-icon="alert-circle-fill"]`)
      .exists('row has failure icon');
    assert.dom('[data-test-cross-sign-alert-title]').hasText('Cross-sign failed');
    assert
      .dom('[data-test-cross-sign-alert-message]')
      .hasText('Cross-signing a root issuer with itself must be performed manually using the CLI.');

    for (const field of FIELDS) {
      assert
        .dom(`${PKI_CROSS_SIGN.signedIssuerCol(field.key)}`)
        .hasText(this.testInputs[field.key], `${field.key} displays correct value`);
    }
  });
});
