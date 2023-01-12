import Model, { attr } from '@ember-data/model';
import { inject as service } from '@ember/service';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withFormFields } from 'vault/decorators/model-form-fields';

@withFormFields()
export default class PkiConfigModel extends Model {
  @service secretMountPath;

  /* formType import */
  @attr('string') pemBundle;

  /* formType generate-root */
  @attr('string', {
    possibleValues: ['exported', 'internal', 'existing', 'kms'],
  })
  type;

  @attr('string') issuerName; // REQUIRED, cannot be "default"

  @attr('string') keyName; // cannot be "default"

  @attr('string') keyRef; // search-select? only for type=existing

  @attr('string') commonName; // REQUIRED

  @attr('string', {
    label: 'DNS/Email Subject Alternative Names (SANs)',
  })
  altNames; // comma sep strings

  @attr('string', {
    label: 'IP Subject Alternative Names (SANs)',
  })
  ipSans;

  @attr('string') uriSans;

  @attr('string') otherSans;

  @attr('string', {
    defaultValue: 'pem',
    possibleValues: ['pem', 'der', 'pem_bundle'],
  })
  format;

  @attr('string', {
    defaultValue: 'der',
    possibleValues: ['der', 'pkcs8'],
  })
  privateKeyFormat;

  @attr('string', {
    defaultValue: 'rsa',
    possibleValues: ['rsa', 'ed25519', 'ec'],
  })
  keyType;

  @attr('number', {
    defaultValue: 0,
    // options management happens in pki-key-parameters
  })
  keyBits;

  @attr('number', {
    defaultValue: -1,
  })
  maxPathLength;

  @attr('boolean', {
    defaultValue: false,
  })
  excludeCnFromSans;

  @attr('string', {
    label: 'Permitted DNS domains',
  })
  permittedDnsDomains;

  @attr('string') ou;
  @attr('string') organization;
  @attr('string') country;
  @attr('string') locality;
  @attr('string') province;
  @attr('string') streetAddress;
  @attr('string') postalCode;

  @attr('string') serialNumber;

  @attr({
    label: 'Backdate validity',
    detailsLabel: 'Issued certificate backdating',
    helperTextDisabled: 'Vault will use the default value, 30s',
    helperTextEnabled:
      'Also called the not_before_duration property. Allows certificates to be valid for a certain time period before now. This is useful to correct clock misalignment on various systems when setting up your CA.',
    editType: 'ttl',
    defaultValue: '30s',
  })
  notBeforeDuration;

  @attr('string') managedKeyName;
  @attr('string') managedKeyId;

  @attr({
    label: 'Not valid after',
    detailsLabel: 'Issued certificates expire after',
    subText:
      'The time after which this certificate will no longer be valid. This can be a TTL (a range of time from now) or a specific date.',
    editType: 'yield',
  })
  customTtl;
  @attr('string') ttl;
  @attr('date') notAfter;

  get backend() {
    return this.secretMountPath.currentPath;
  }

  @lazyCapabilities(apiPath`${'backend'}/issuers/import/bundle`, 'backend') importBundlePath;
  @lazyCapabilities(apiPath`${'backend'}/issuers/generate/root/${'type'}`, 'backend', 'type')
  generateIssuerRootPath;
  @lazyCapabilities(apiPath`${'backend'}/issuers/generate/intermediate/${'type'}`, 'backend', 'type')
  generateIssuerCsrPath;

  get canImportBundle() {
    return this.importBundlePath.get('canCreate') !== false;
  }
  get canGenerateIssuerRoot() {
    return this.generateIssuerRootPath.get('canCreate') !== false;
  }
  get canGenerateIssuerIntermediate() {
    return this.generateIssuerCsrPath.get('canCreate') !== false;
  }
  get canCrossSign() {
    return this.crossSignPath.get('canCreate') !== false;
  }
}
