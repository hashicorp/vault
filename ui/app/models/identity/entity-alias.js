import IdentityModel from './_base';
import DS from 'ember-data';
import Ember from 'ember';
import { queryRecord } from 'ember-computed-query';
const { attr, belongsTo } = DS;
const { computed } = Ember;

export default IdentityModel.extend({
  parentType: 'entity',
  formFields: ['name', 'mountAccessor', 'metadata'],
  entity: belongsTo('identity/entity', { readOnly: true, async: false }),

  name: attr('string'),
  canonicalId: attr('string'),
  mountAccessor: attr('string', {
    label: 'Auth Backend',
    editType: 'mountAccessor',
  }),
  metadata: attr({
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

  updatePath: queryRecord(
    'capabilities',
    context => {
      const { identityType, id } = context.getProperties('identityType', 'id');
      //identity/entity-alias/id/efb8b562-77fd-335f-a754-740373a778e6
      return {
        id: `identity/${identityType}/id/${id}`,
      };
    },
    'id',
    'identityType'
  ),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),
});
