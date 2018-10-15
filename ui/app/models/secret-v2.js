import DS from 'ember-data';
import { computed } from '@ember/object';
import { match } from '@ember/object/computed';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr, hasMany, belongsTo, Model } = DS;

export default Model.extend({
  engine: belongsTo('secret-engine'),
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
  isFolder: match('id', /\/$/),
  fields: computed(function() {
    return expandAttributeMeta(this, ['maxVersions', 'casRequired']);
  }),
});
