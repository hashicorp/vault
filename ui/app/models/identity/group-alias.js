import IdentityModel from './_base';
import DS from 'ember-data';
const { attr, belongsTo } = DS;

export default IdentityModel.extend({
  parentType: 'group',
  formFields: ['name', 'mountAccessor'],
  group: belongsTo('identity/group', { readOnly: true, async: false }),

  name: attr('string'),
  canonicalId: attr('string'),

  mountPath: attr('string', {
    readOnly: true,
  }),
  mountType: attr('string', {
    readOnly: true,
  }),
  mountAccessor: attr('string', {
    label: 'Auth Backend',
    editType: 'mountAccessor',
  }),

  creationTime: attr('string', {
    readOnly: true,
  }),
  lastUpdateTime: attr('string', {
    readOnly: true,
  }),
});
