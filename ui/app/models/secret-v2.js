import DS from 'ember-data';
import { match } from '@ember/object/computed';

const { attr, hasMany, belongsTo, Model } = DS;

export default Model.extend({
  engine: belongsTo('secret-engine'),
  versions: hasMany('secret-v2-version', { async: false, inverse: null }),
  selectedVersion: belongsTo('secret-v2-version', { inverse: 'secret' }),
  createdTime: attr(),
  updatedTime: attr(),
  currentVersion: attr('number'),
  oldestVersion: attr('number'),
  maxVersions: attr('number'),
  casRequired: attr('boolean'),
  isFolder: match('id', /\/$/),
});
