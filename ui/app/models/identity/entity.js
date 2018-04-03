import IdentityModel from './_base';
import DS from 'ember-data';
const { attr, hasMany } = DS;

export default IdentityModel.extend({
  formFields: ['name', 'policies', 'metadata'],
  name: attr('string'),
  mergedEntityIds: attr(),
  metadata: attr('object', {
    editType: 'kv',
  }),
  policies: attr({
    editType: 'stringArray',
  }),
  creationTime: attr('string', {
    readOnly: true,
  }),
  lastUpdateTime: attr('string', {
    readOnly: true,
  }),
  aliases: hasMany('identity/entity-alias', { async: false, readOnly: true }),
  groupIds: attr({
    readOnly: true,
  }),
  directGroupIds: attr({
    readOnly: true,
  }),
  inheritedGroupIds: attr({
    readOnly: true,
  }),
});
