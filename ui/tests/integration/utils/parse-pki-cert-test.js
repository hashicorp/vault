/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { parseCertificate, parseExtensions, parseSubject, formatValues } from 'vault/utils/parse-pki-cert';
import * as asn1js from 'asn1js';
import { fromBase64, stringToArrayBuffer } from 'pvutils';
import { Certificate } from 'pkijs';
import { addHours, fromUnixTime, isSameDay } from 'date-fns';
import errorMessage from 'vault/utils/error-message';
import { OTHER_OIDs, SAN_TYPES } from 'vault/utils/parse-pki-cert-oids';
import { verifyCertificates, jsonToCertObject, verifySignature } from 'vault/utils/parse-pki-cert';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';

const {
  certWithoutCN,
  loadedCert,
  pssTrueCert,
  skeletonCert,
  unsupportedOids,
  unsupportedSignatureRoot,
  unsupportedSignatureInt,
} = CERTIFICATES;

module('Integration | Util | parse pki certificate', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.getErrorMessages = (certErrors) => certErrors.map((error) => errorMessage(error));
    this.certSchema = (cert) => {
      const cert_base64 = cert.replace(/(-----(BEGIN|END) CERTIFICATE-----|\n)/g, '');
      const cert_der = fromBase64(cert_base64);
      const cert_asn1 = asn1js.fromBER(stringToArrayBuffer(cert_der));
      return new Certificate({ schema: cert_asn1.result });
    };
    this.parsableLoadedCert = this.certSchema(loadedCert);
    this.parsableUnsupportedCert = this.certSchema(unsupportedOids);
  });

  test('it parses a certificate with supported values', async function (assert) {
    assert.expect(2);
    // certificate contains all allowable params
    const parsedCert = parseCertificate(loadedCert);
    assert.propEqual(
      parsedCert,
      {
        alt_names: 'altname1, altname2',
        can_parse: true,
        common_name: 'common-name.com',
        country: 'France',
        other_sans: '1.3.1.4.1.5.9.2.6;UTF8:some-utf-string',
        exclude_cn_from_sans: true,
        ip_sans: '192.158.1.38, 1234:0fd2:5621:0001:0089:0000:0000:4500',
        key_usage: 'CertSign, CRLSign',
        locality: 'Paris',
        max_path_length: 17,
        not_valid_after: 1678210083,
        not_valid_before: 1675445253,
        organization: 'Widget',
        ou: 'Finance',
        parsing_errors: [],
        permitted_dns_domains: 'dnsname1.com, dsnname2.com',
        postal_code: '123456',
        province: 'Champagne',
        subject_serial_number: 'cereal1292',
        signature_bits: '256',
        street_address: '234 sesame',
        ttl: '768h',
        uri_sans: 'testuri1, testuri2',
        use_pss: false,
      },
      'it contains expected attrs, cn is excluded from alt_names (exclude_cn_from_sans: true) and ipV6 is compressed correctly'
    );
    assert.ok(
      isSameDay(
        addHours(fromUnixTime(parsedCert.not_valid_before), Number(parsedCert.ttl.split('h')[0])),
        fromUnixTime(parsedCert.not_valid_after),
        'ttl value is correct'
      )
    );
  });

  test('it parses a certificate with use_pass=true and exclude_cn_from_sans=false', async function (assert) {
    assert.expect(2);
    const parsedPssCert = parseCertificate(pssTrueCert);
    assert.propContains(
      parsedPssCert,
      { signature_bits: '256', ttl: '768h', use_pss: true },
      'returns signature_bits value and use_pss is true'
    );
    assert.propContains(
      parsedPssCert,
      {
        alt_names: 'common-name.com',
        can_parse: true,
        common_name: 'common-name.com',
        exclude_cn_from_sans: false,
      },
      'common name is included in alt_names'
    );
  });

  test('it returns parsing_errors when certificate has unsupported values', async function (assert) {
    assert.expect(2);
    const parsedCert = parseCertificate(unsupportedOids); // contains unsupported subject and extension OIDs
    const parsingErrors = this.getErrorMessages(parsedCert.parsing_errors);
    assert.propContains(
      parsedCert,
      {
        alt_names: 'dns-NameSupported',
        common_name: 'fancy-cert-unsupported-subj-and-ext-oids',
        ip_sans: '192.158.1.38',
        parsing_errors: [{}, {}],
        uri_sans: 'uriSupported',
      },
      'supported values are present when unsupported values exist'
    );
    assert.propEqual(
      parsingErrors,
      [
        'certificate contains unsupported subject OIDs: 1.2.840.113549.1.9.1',
        'certificate contains unsupported extension OIDs: 2.5.29.37',
      ],
      'it contains expected error messages'
    );
  });

  test('it returns attr with a null value if nonexistent', async function (assert) {
    assert.expect(1);
    const onlyHasCommonName = parseCertificate(skeletonCert);
    assert.propContains(
      onlyHasCommonName,
      {
        alt_names: 'common-name.com',
        common_name: 'common-name.com',
        country: null,
        ip_sans: null,
        locality: null,
        max_path_length: undefined,
        organization: null,
        ou: null,
        postal_code: null,
        province: null,
        subject_serial_number: null,
        street_address: null,
        uri_sans: null,
      },
      'it contains expected attrs'
    );
  });

  test('the helper parseSubject returns object with correct key/value pairs', async function (assert) {
    assert.expect(3);
    const supportedSubj = parseSubject(this.parsableLoadedCert.subject.typesAndValues);
    assert.propEqual(
      supportedSubj,
      {
        subjErrors: [],
        subjValues: {
          common_name: 'common-name.com',
          country: 'France',
          locality: 'Paris',
          organization: 'Widget',
          ou: 'Finance',
          postal_code: '123456',
          province: 'Champagne',
          subject_serial_number: 'cereal1292',
          street_address: '234 sesame',
        },
      },
      'it returns supported subject values'
    );

    const unsupportedSubj = parseSubject(this.parsableUnsupportedCert.subject.typesAndValues);
    assert.propEqual(
      this.getErrorMessages(unsupportedSubj.subjErrors),
      ['certificate contains unsupported subject OIDs: 1.2.840.113549.1.9.1'],
      'it returns subject errors'
    );
    assert.ok(
      unsupportedSubj.subjErrors.every((e) => e instanceof Error),
      'subjErrors contain error objects'
    );
  });

  test('the helper parseExtensions returns object with correct key/value pairs', async function (assert) {
    assert.expect(11);
    // assert supported extensions return correct type
    const supportedExtensions = parseExtensions(this.parsableLoadedCert.extensions);
    let { extValues, extErrors } = supportedExtensions;
    for (const keyName in SAN_TYPES) {
      assert.ok(Array.isArray(extValues[keyName]), `${keyName} is an array`);
    }
    assert.ok(Array.isArray(extValues.permitted_dns_domains), 'permitted_dns_domains is an array');
    assert.ok(Number.isInteger(extValues.max_path_length), 'max_path_length is an integer');
    assert.propEqual(extValues.key_usage, ['CertSign', 'CRLSign'], 'parses key_usage');
    assert.strictEqual(extErrors.length, 0, 'no extension errors');

    // assert unsupported extensions return errors
    const unsupportedExt = parseExtensions(this.parsableUnsupportedCert.extensions);
    ({ extValues, extErrors } = unsupportedExt);
    assert.propEqual(
      this.getErrorMessages(extErrors),
      ['certificate contains unsupported extension OIDs: 2.5.29.37'],
      'it returns extension errors'
    );
    assert.ok(
      extErrors.every((e) => e instanceof Error),
      'subjErrors contain error objects'
    );
    assert.ok(Number.isInteger(extValues.max_path_length), 'max_path_length is an integer');
  });

  test('the helper formatValues returns object with correct types', async function (assert) {
    assert.expect(1);
    const supportedSubj = parseSubject(this.parsableLoadedCert.subject.typesAndValues);
    const supportedExtensions = parseExtensions(this.parsableLoadedCert.extensions);
    assert.propContains(
      formatValues(supportedSubj, supportedExtensions),
      {
        alt_names: 'altname1, altname2',
        ip_sans: '192.158.1.38, 1234:0fd2:5621:0001:0089:0000:0000:4500',
        permitted_dns_domains: 'dnsname1.com, dsnname2.com',
        uri_sans: 'testuri1, testuri2',
        parsing_errors: [],
        exclude_cn_from_sans: true,
      },
      `values for ${Object.keys(SAN_TYPES).join(', ')} are comma separated strings (and no longer arrays)`
    );
  });

  test('the helper verifyCertificates catches errors', async function (assert) {
    assert.expect(5);
    const verifiedRoot = await verifyCertificates(unsupportedSignatureRoot, unsupportedSignatureRoot);
    assert.true(verifiedRoot, 'returns true for root certificate');
    const verifiedInt = await verifyCertificates(unsupportedSignatureInt, unsupportedSignatureInt);
    assert.false(verifiedInt, 'returns false for intermediate cert');

    const filterExtensions = (list, oid) => list.filter((ext) => ext.extnID !== oid);
    const { subject_key_identifier, authority_key_identifier } = OTHER_OIDs;
    const testCert = jsonToCertObject(unsupportedSignatureRoot);
    const certWithoutSKID = testCert;
    certWithoutSKID.extensions = filterExtensions(testCert.extensions, subject_key_identifier);
    assert.false(
      await verifySignature(certWithoutSKID, certWithoutSKID),
      'returns false if no subject key ID'
    );

    const certWithoutAKID = testCert;
    certWithoutAKID.extensions = filterExtensions(testCert.extensions, authority_key_identifier);
    assert.false(await verifySignature(certWithoutAKID, certWithoutAKID), 'returns false if no AKID');

    const certWithoutKeyID = testCert;
    certWithoutAKID.extensions = [];
    assert.false(
      await verifySignature(certWithoutKeyID, certWithoutKeyID),
      'returns false if neither SKID or AKID'
    );
  });

  test('it fails silently when passed null', async function (assert) {
    assert.expect(3);
    const parsedCert = parseCertificate(certWithoutCN);
    assert.propEqual(
      parsedCert,
      {
        can_parse: true,
        common_name: null,
        country: null,
        exclude_cn_from_sans: false,
        key_usage: null,
        locality: null,
        max_path_length: 10,
        not_valid_after: 1989876490,
        not_valid_before: 1674516490,
        organization: null,
        ou: null,
        parsing_errors: [{}, {}],
        postal_code: null,
        province: null,
        subject_serial_number: null,
        signature_bits: '256',
        street_address: null,
        ttl: '87600h',
        use_pss: false,
      },
      'it parses a cert without CN'
    );
    const parsingErrors = this.getErrorMessages(parsedCert.parsing_errors);
    assert.propEqual(
      parsingErrors,
      [
        'certificate contains unsupported subject OIDs: 1.2.840.113549.1.9.1',
        'certificate contains unsupported extension OIDs: 2.5.29.37',
      ],
      'it returns correct errors'
    );
    assert.propEqual(
      formatValues(null, null),
      { parsing_errors: [Error('error parsing certificate')] },
      'it returns error if unable to format values'
    );
  });
});
