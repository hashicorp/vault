import Ember from 'ember';
import IdentityModel from './_base';
import DS from 'ember-data';
import { queryRecord } from 'ember-computed-query';
const { computed } = Ember;

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

  updatePath: queryRecord(
    'capabilities',
    context => {
      const { identityType, id } = context.getProperties('identityType', 'id');
      //identity/entity/id/efb8b562-77fd-335f-a754-740373a778e6
      return {
        id: `identity/${identityType}/id/${id}`,
      };
    },
    'id',
    'identityType'
  ),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),
  aliasPath: queryRecord('capabilities', () => {
    //identity/entity-alias
    return {
      id: `identity/entity-alias`,
    };
  }),

  canAddAlias: computed.alias('aliasPath.canCreate'),
});
