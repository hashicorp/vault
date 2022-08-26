import Model, { belongsTo, hasMany, attr } from '@ember-data/model';
import { computed } from '@ember/object'; // eslint-disable-line
import { alias } from '@ember/object/computed'; // eslint-disable-line
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import KeyMixin from 'vault/mixins/key-mixin';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  maxVersions: [
    { type: 'number', message: 'Maximum versions must be a number.' },
    { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
  ],
};

@withModelValidations(validations)
class SecretV2Model extends Model {}
export default SecretV2Model.extend(KeyMixin, {
  failedServerRead: attr('boolean'),
  engine: belongsTo('secret-engine', { async: false }),
  engineId: attr('string'),
  versions: hasMany('secret-v2-version', { async: false, inverse: null }),
  selectedVersion: belongsTo('secret-v2-version', { async: false, inverse: 'secret' }),
  createdTime: attr(),
  updatedTime: attr(),
  currentVersion: attr('number'),
  oldestVersion: attr('number'),
  customMetadata: attr('object', {
    editType: 'kv',
    subText: 'An optional set of informational key-value pairs that will be stored with all secret versions.',
  }),
  maxVersions: attr('number', {
    defaultValue: 10,
    label: 'Maximum number of versions',
    subText:
      'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted.',
  }),
  casRequired: attr('boolean', {
    defaultValue: false,
    label: 'Require Check and Set',
    subText:
      'Writes will only be allowed if the key’s current version matches the version specified in the cas parameter.',
  }),
  deleteVersionAfter: attr({
    defaultValue: 0,
    editType: 'ttl',
    label: 'Automate secret deletion',
    helperTextDisabled: 'A secret’s version must be manually deleted.',
    helperTextEnabled: 'Delete all new versions of this secret after',
  }),
  fields: computed(function () {
    return expandAttributeMeta(this, ['customMetadata', 'maxVersions', 'casRequired', 'deleteVersionAfter']);
  }),
  secretDataPath: lazyCapabilities(apiPath`${'engineId'}/data/${'id'}`, 'engineId', 'id'),
  secretMetadataPath: lazyCapabilities(apiPath`${'engineId'}/metadata/${'id'}`, 'engineId', 'id'),

  canListMetadata: alias('secretMetadataPath.canList'),
  canReadMetadata: alias('secretMetadataPath.canRead'),
  canUpdateMetadata: alias('secretMetadataPath.canUpdate'),

  canReadSecretData: alias('secretDataPath.canRead'),
  canEditSecretData: alias('secretDataPath.canUpdate'),
  canDeleteSecretData: alias('secretDataPath.canDelete'),
});
