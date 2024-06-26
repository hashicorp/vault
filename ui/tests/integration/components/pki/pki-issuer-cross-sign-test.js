/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';
import { PKI_CROSS_SIGN } from 'vault/tests/helpers/pki/pki-selectors';

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
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const store = this.owner.lookup('service:store');
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

    this.testInputs = {
      intermediateMount: this.intMountPath,
      intermediateIssuer: this.intIssuerData.issuer_name,
      newCrossSignedIssuer: this.newIssuerData.issuer_name,
    };

    store.pushPayload('pki/issuer', { modelName: 'pki/issuer', data: this.parentIssuerData });
    this.parentIssuerModel = store.peekRecord('pki/issuer', this.parentIssuerData.issuer_id);
  });

  test('it makes requests to the correct endpoints', async function (assert) {
    assert.expect(18);
    this.server.get(`/${this.intMountPath}/issuer/${this.intIssuerData.issuer_name}`, () => {
      assert.ok(true, 'Step 1. GET request is made to fetch existing issuer data');
      return { data: this.intIssuerData };
    });
    this.server.post(`/${this.intMountPath}/intermediate/generate/existing`, (schema, req) => {
      assert.ok(true, 'Step 2. POST request is made to generate new CSR');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          common_name: newCSR.common_name,
          country: null,
          exclude_cn_from_sans: false,
          format: 'pem',
          locality: null,
          organization: null,
          ou: null,
          province: null,
          key_ref: this.intIssuerData.key_id,
        },
        'payload contains correct key ref'
      );
      return {
        data: { csr: newCSR.csr, key_id: this.intIssuerData.key_id },
        request_id: '1234',
      };
    });
    this.server.post(
      `/${this.backend}/issuer/${this.parentIssuerData.issuer_name}/sign-intermediate`,
      (schema, req) => {
        assert.ok(true, 'Step 3. POST request is made to sign CSR with new parent issuer');
        assert.propEqual(JSON.parse(req.requestBody), newCSR, 'payload has common name and csr');
        return {
          data: { ca_chain: [newlySignedCert, parentIssuerCert] },
          request_id: '1234',
        };
      }
    );
    this.server.post(`/${this.intMountPath}/issuers/import/bundle`, (schema, req) => {
      assert.ok(true, 'Step 4. POST request made to import issuer');
      assert.propEqual(
        JSON.parse(req.requestBody),
        { pem_bundle: [newlySignedCert, parentIssuerCert].join('\n') },
        'payload contains pem bundle'
      );
      return {
        request_id: '1234',
        data: {
          imported_issuers: null,
          imported_keys: null,
          mapping: { [this.newIssuerData.issuer_id]: this.intIssuerData.key_id },
        },
      };
    });
    this.server.get(`/${this.intMountPath}/issuer/${this.newIssuerData.issuer_id}`, () => {
      assert.ok(true, 'Step 5. GET request is made to newly imported issuer');
      return { data: this.newIssuerData };
    });

    this.server.post(`/${this.intMountPath}/issuer/${this.newIssuerData.issuer_id}`, (schema, req) => {
      assert.ok(true, 'Step 6. POST request is made to update issuer name');
      assert.propEqual(
        JSON.parse(req.requestBody),
        {
          issuer_name: 'newly-cross-signed-cert',
          leaf_not_after_behavior: 'err',
          usage: 'crl-signing,issuing-certificates,ocsp-signing,read-only',
        },
        'payload has correct data '
      );
      return { data: this.newIssuerData };
    });

    await render(hbs`<PkiIssuerCrossSign @parentIssuer={{this.parentIssuerModel}} /> `, {
      owner: this.engine,
    });
    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.saveButton);

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
    assert.expect(13);
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

    this.server.get(`/${this.intMountPath}/issuer/${this.intIssuerData.issuer_name}`, () => {
      assert.ok(true, 'request is made to sign first cert');
      return { data: this.intIssuerData };
    });

    this.server.get(`/${this.intMountPath}/issuer/${nonexistentIssuer.intermediateIssuer}`, () => {
      assert.ok(true, 'request is made to second cert');
      return new Response(
        500,
        { 'Content-Type': 'application/json' },
        JSON.stringify({
          errors: [
            `1 error occurred:\n\t* unable to find PKI issuer for reference: ${nonexistentIssuer.intermediateIssuer}\n\n`,
          ],
        })
      );
    });

    this.server.get(`/${this.intMountPath}/issuer/${unsupportedCert.intermediateIssuer}`, () => {
      assert.ok(true, 'request is made to third cert');
      return { data: { isser_name: unsupportedCert.intermediateIssuer, certificate: unsupportedOids } };
    });

    await render(hbs`<PkiIssuerCrossSign @parentIssuer={{this.parentIssuerModel}} /> `, {
      owner: this.engine,
    });

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
    await click(GENERAL.saveButton);

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
    this.server.get(`/${this.intMountPath}/issuer/${this.intIssuerData.issuer_name}`, () => {
      return new Response(
        500,
        { 'Content-Type': 'application/json' },
        JSON.stringify({
          errors: ['1 error occurred:\n\t* unable to find PKI issuer for reference: nonexistent-mount\n\n'],
        })
      );
    });

    await render(hbs`<PkiIssuerCrossSign @parentIssuer={{this.parentIssuerModel}} /> `, {
      owner: this.engine,
    });

    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.saveButton);

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
    this.server.get(`/${this.intMountPath}/issuer/${this.intIssuerData.issuer_name}`, () => {
      return { data: unsupportedIssuerCert };
    });

    await render(hbs`<PkiIssuerCrossSign @parentIssuer={{this.parentIssuerModel}} /> `, {
      owner: this.engine,
    });
    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.saveButton);
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
        'certificate contains unsupported subject OIDs: 1.2.840.113549.1.9.1, certificate contains unsupported extension OIDs: 2.5.29.37'
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
    this.server.get(`/${this.backend}/issuer/${this.parentIssuerData.issuer_name}`, () => {
      return { data: this.parentIssuerData };
    });

    await render(hbs`<PkiIssuerCrossSign @parentIssuer={{this.parentIssuerModel}} /> `, {
      owner: this.engine,
    });
    // fill out form and submit
    for (const field of FIELDS) {
      await fillIn(PKI_CROSS_SIGN.objectListInput(field.key), this.testInputs[field.key]);
    }
    await click(GENERAL.saveButton);
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
