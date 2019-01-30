import DS from 'ember-data';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import KeyMixin from 'vault/mixins/key-mixin';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const { attr, hasMany, belongsTo, Model } = DS;

export default Model.extend(KeyMixin, {
  engine: belongsTo('secret-engine', { async: false }),
  engineId: attr('string'),
  versions: hasMany('secret-v2-version', { async: false, inverse: null }),
  selectedVersion: belongsTo('secret-v2-version', { async: false, inverse: 'secret' }),
  createdTime: attr(),
  updatedTime: attr(),
  currentVersion: attr('number'),
  oldestVersion: attr('number'),
  maxVersions: attr('number', {
    defaultValue: 10,
    label: 'Maximum Number of Versions',
  }),
  casRequired: attr('boolean', {
    defaultValue: false,
    label: 'Require Check and Set',
    helpText:
      'Writes will only be allowed if the key’s current version matches the version specified in the cas parameter',
  }),
  fields: computed(function() {
    return expandAttributeMeta(this, ['maxVersions', 'casRequired']);
  }),
  versionPath: lazyCapabilities(apiPath`${'engineId'}/data/${'id'}`, 'engineId', 'id'),
  secretPath: lazyCapabilities(apiPath`${'engineId'}/metadata/${'id'}`, 'engineId', 'id'),

  canEdit: alias('versionPath.canUpdate'),
  canDelete: alias('secretPath.canDelete'),
  canRead: alias('secretPath.canRead'),
});
