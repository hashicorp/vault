/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { assert } from '@ember/debug';
import { service } from '@ember/service';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

/**
 * There are many actions that involve certificates in PKI world.
 * The base certificate model contains shared attributes that make up a certificate's content.
 * Other models under pki/certificate will extend this model and include additional attributes
 * and associated adapter methods for performing various generation and signing actions.
 * This model also displays leaf certs and their parsed attributes (parsed parameters only
 * render if included in certDisplayFields below).
 */

const certDisplayFields = [
  'certificate',
  'commonName',
  'revocationTime',
  'serialNumber',
  'notValidBefore',
  'notValidAfter',
];

@withFormFields(certDisplayFields)
export default class PkiCertificateBaseModel extends Model {
  @service secretMountPath;

  get useOpenAPI() {
    return true;
  }
  get backend() {
    return this.secretMountPath.currentPath;
  }
  getHelpUrl() {
    assert('You must provide a helpUrl for OpenAPI', true);
  }

  @attr('string') commonName;

  @attr({
    label: 'Not valid after',
    detailsLabel: 'Issued certificates expire after',
    subText:
      'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
    editType: 'yield',
  })
  customTtl; // sets ttl and notAfter via one input <PkiNotValidAfterForm>

  @attr('boolean', {
    label: 'Exclude common name from SANs',
    subText:
      'If checked, the common name will not be included in DNS or Email Subject Alternate Names. This is useful if the CN is a human-readable identifier, not a hostname or email address.',
    defaultValue: false,
  })
  excludeCnFromSans;

  @attr('string', {
    label: 'Subject Alternative Names (SANs)',
    subText:
      'The requested Subject Alternative Names; if email protection is enabled for the role, this may contain email addresses. Add one per row.',
    editType: 'stringArray',
  })
  altNames;

  // SANs below are editType: stringArray from openApi
  @attr('string', {
    label: 'IP Subject Alternative Names (IP SANs)',
    subText: 'Only valid if the role allows IP SANs (which is the default). Add one per row.',
  })
  ipSans;

  @attr('string', {
    label: 'URI Subject Alternative Names (URI SANs)',
    subText:
      'If any requested URIs do not match role policy, the entire request will be denied. Add one per row.',
  })
  uriSans;

  @attr('string', {
    subText:
      'Requested other SANs with the format <oid>;UTF8:<utf8 string value> for each entry. Add one per row.',
  })
  otherSans;

  // Attrs that come back from API POST request
  @attr({ label: 'CA Chain', masked: true }) caChain;
  @attr('string', { masked: true }) certificate;
  @attr('number') expiration;
  @attr('string', { label: 'Issuing CA', masked: true }) issuingCa;
  @attr('string') privateKey; // only returned for type=exported
  @attr('string') privateKeyType; // only returned for type=exported
  @attr('number', { formatDate: true }) revocationTime;
  @attr('string') serialNumber;

  // read only attrs parsed from certificate contents in serializer on GET requests (see parse-pki-cert.js)
  @attr('number', { formatDate: true }) notValidAfter; // set by ttl or notAfter (customTtL above)
  @attr('number', { formatDate: true }) notValidBefore; // date certificate was issued
  @attr('string') signatureBits;

  @lazyCapabilities(apiPath`${'backend'}/revoke`, 'backend') revokePath;
  get canRevoke() {
    return this.revokePath.get('isLoading') || this.revokePath.get('canCreate') !== false;
  }
}
