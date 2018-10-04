import Secret from './secret';
import DS from 'ember-data';

const { attr, hasMany, belongsTo, Model } = DS;

export default Model.extend({
  engine: belongsTo('secret-engine'),
  versions: hasMany('secret-v2-version', { async: false }),
  createdTime: attr(),
  updatedTime: attr(),
  currentVersion: attr('number'),
  oldestVersion: attr('number'),
  maxVersions: attr('number'),
  casRequired: attr('boolean'),
});
