import Model, { belongsTo, hasMany, attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import KeyMixin from 'vault/mixins/key-mixin';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { validator, buildValidations } from 'ember-cp-validations';

const Validations = buildValidations({
  // ARG TODO later:
  //The new custom_metadata field will compete with the version history for storage space in the key metadata entry. To attempt to prevent contention, Vault will impose limits on both the keys and values within the custom_metadata field. The keys and values will be limited to 128 and 512 bytes, respectively. Vault will also impose a limit of 64 total custom_metadata keys.
  customMetadata: validator('format', {
    regex: /^[^\/]+$/,
    message: 'Custom Values cannot contain a forward slash.',
  }),
  maxVersions: [
    validator('number', {
      allowString: false,
      integer: true,
      message: 'Maximum versions must be a number.',
    }),
    validator('length', {
      min: 1,
      max: 16,
      message: 'You cannot go over 16 characters.',
    }),
  ],
});

export default Model.extend(KeyMixin, Validations, {
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
    label: 'Maximum Number of Versions',
    subText:
      'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted.',
  }),
  casRequired: attr('boolean', {
    defaultValue: false,
    label: 'Require Check and Set',
    subText:
      'Writes will only be allowed if the key’s current version matches the version specified in the cas parameter',
  }),
  deleteVersionAfter: attr({
    defaultValue: 0,
    editType: 'ttl',
    label: 'Automate secret deletion',
    helperTextDisabled: 'A secret’s version must be manually deleted.',
    helperTextEnabled: 'Delete all new versions of this secret after',
  }),
  fields: computed(function() {
    return expandAttributeMeta(this, ['customMetadata', 'maxVersions', 'casRequired', 'deleteVersionAfter']);
  }),
  versionPath: lazyCapabilities(apiPath`${'engineId'}/data/${'id'}`, 'engineId', 'id'),
  secretPath: lazyCapabilities(apiPath`${'engineId'}/metadata/${'id'}`, 'engineId', 'id'),

  canEdit: alias('versionPath.canUpdate'),
  canDelete: alias('secretPath.canDelete'),
  canRead: alias('secretPath.canRead'),
});
