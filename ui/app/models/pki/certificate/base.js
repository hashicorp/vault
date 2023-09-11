/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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
 * This model also displays leaf certs and their parsed attributes (which exist as an object in
 * the attribute `parsedCertificate`)
 */

// also displays parsedCertificate values in the template
const certDisplayFields = ['certificate', 'commonName', 'revocationTime', 'serialNumber'];

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

  // The attributes parsed from parse-pki-cert util live here
  @attr parsedCertificate;

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
      'The requested Subject Alternative Names; if email protection is enabled for the role, this may contain email addresses.',
    editType: 'stringArray',
  })
  altNames;

  // SANs below are editType: stringArray from openApi
  @attr('string', {
    label: 'IP Subject Alternative Names (IP SANs)',
    subText: 'Only valid if the role allows IP SANs (which is the default).',
  })
  ipSans;

  @attr('string', {
    label: 'URI Subject Alternative Names (URI SANs)',
    subText: 'If any requested URIs do not match role policy, the entire request will be denied.',
  })
  uriSans;

  @attr('string', {
    subText: 'Requested other SANs with the format <oid>;UTF8:<utf8 string value> for each entry.',
  })
  otherSans;

  // Attrs that come back from API POST request
  @attr({ label: 'CA Chain', isCertificate: true }) caChain;
  @attr('string', { isCertificate: true }) certificate;
  @attr('number') expiration;
  @attr('string', { label: 'Issuing CA', isCertificate: true }) issuingCa;
  @attr('string', { isCertificate: true }) privateKey; // only returned for type=exported and /issue
  @attr('string') privateKeyType; // only returned for type=exported and /issue
  @attr('number', { formatDate: true }) revocationTime;
  @attr('string') serialNumber;

  @lazyCapabilities(apiPath`${'backend'}/revoke`, 'backend') revokePath;
  get canRevoke() {
    return this.revokePath.get('isLoading') || this.revokePath.get('canCreate') !== false;
  }
}
