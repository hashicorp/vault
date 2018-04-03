import IdentityModel from './_base';
import DS from 'ember-data';
const { attr, belongsTo } = DS;

export default IdentityModel.extend({
  formFields: ['name', 'mountAccessor', 'metadata'],
  entity: belongsTo('identity/entity', { readOnly: true, async: false }),

  name: attr('string'),
  canonicalId: attr('string'),
  mountAccessor: attr('string', {
    label: 'Auth Backend',
    editType: 'mountAccessor',
  }),
  metadata: attr('object', {
    editType: 'kv',
  }),
  mountPath: attr('string', {
    readOnly: true,
  }),
  mountType: attr('string', {
    readOnly: true,
  }),
  creationTime: attr('string', {
    readOnly: true,
  }),
  lastUpdateTime: attr('string', {
    readOnly: true,
  }),
  mergedFromCanonicalIds: attr(),
});
