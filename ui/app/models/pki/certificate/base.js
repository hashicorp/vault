import Model, { attr } from '@ember-data/model';
import { assert } from '@ember/debug';
import { service } from '@ember/service';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

/**
 * There are many ways to generate a cert, but we want to display them in a consistent way.
 * This base certificate model will set the attributes we want to display, and other
 * models under pki/certificate will extend this model and have their own required
 * attributes and adapter methods.
 */

const certDisplayFields = [
  'certificate',
  'commonName',
  'revocationTime',
  'issueDate',
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

  // Required input for all certificates
  @attr('string') commonName;

  // Attrs that come back from API POST request
  @attr() caChain;
  @attr('string', { masked: true }) certificate;
  @attr('number') expiration;
  @attr('number', { formatDate: true }) revocationTime;
  @attr('string') issuingCa;
  @attr('string') privateKey;
  @attr('string') privateKeyType;
  @attr('string') serialNumber;

  // Parsed from cert in serializer
  @attr('date', { formatDate: true }) issueDate;
  @attr('date', { formatDate: true }) notValidAfter;
  @attr('date', { formatDate: true }) notValidBefore;

  // For importing
  @attr('string') pemBundle;
  @attr importedIssuers;
  @attr importedKeys;

  @lazyCapabilities(apiPath`${'backend'}/revoke`, 'backend') revokePath;
  get canRevoke() {
    return this.revokePath.get('isLoading') || this.revokePath.get('canCreate') !== false;
  }
}
