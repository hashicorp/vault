import Model, { belongsTo, hasMany, attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withKeyMixin } from 'vault/decorators/key-mixin';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  maxVersions: [
    { type: 'number', options: { asString: true }, message: 'Maximum versions must be a number.' },
    { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
  ],
};

@withKeyMixin()
@withModelValidations(validations)
export default class SecretV2Model extends Model {
  @attr('boolean') failedServerRead;
  @belongsTo('secret-engine', { async: false }) engine;
  @attr('string') engineId;
  @hasMany('secret-v2-version', { async: false, inverse: null }) versions;
  @belongsTo('secret-v2-version', { async: false, inverse: 'secret' }) selectedVersion;
  @attr createdTime;
  @attr updatedTime;
  @attr('number') currentVersion;
  @attr('number') oldestVersion;
  @attr('object', {
    editType: 'kv',
    subText: 'An optional set of informational key-value pairs that will be stored with all secret versions.',
  })
  customMetadata;
  @attr('number', {
    defaultValue: 10,
    label: 'Maximum number of versions',
    subText:
      'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted.',
  })
  maxVersions;
  @attr('boolean', {
    defaultValue: false,
    label: 'Require Check and Set',
    subText:
      'Writes will only be allowed if the key’s current version matches the version specified in the cas parameter.',
  })
  casRequired;
  @attr({
    defaultValue: 0,
    editType: 'ttl',
    label: 'Automate secret deletion',
    helperTextDisabled: 'A secret’s version must be manually deleted.',
    helperTextEnabled: 'Delete all new versions of this secret after',
  })
  deleteVersionAfter;

  // since getters are triggered each time they are accessed this will fire repeatedly on re-render
  // this causes problems with inputs losing focus when continually calling expandAttributeMeta
  // cache result on first get and return that instead
  get fields() {
    if (!this._fieldsCache) {
      this._fieldsCache = expandAttributeMeta(this, [
        'customMetadata',
        'maxVersions',
        'casRequired',
        'deleteVersionAfter',
      ]);
    }
    return this._fieldsCache;
  }

  secretDataPath = lazyCapabilities(apiPath`${'engineId'}/data/${'id'}`, 'engineId', 'id');
  secretMetadataPath = lazyCapabilities(apiPath`${'engineId'}/metadata/${'id'}`, 'engineId', 'id');

  get canListMetadata() {
    return this.secretMetadataPath?.canList || false;
  }
  get canReadMetadata() {
    return this.secretMetadataPath?.canRead || false;
  }
  get canUpdateMetadata() {
    return this.secretMetadataPath?.canUpdate || false;
  }

  get canReadSecretData() {
    return this.secretDataPath?.canRead || false;
  }
  get canEditSecretData() {
    return this.secretDataPath?.canUpdate || false;
  }
  get canDeleteSecretData() {
    return this.secretDataPath?.canDelete || false;
  }
}
