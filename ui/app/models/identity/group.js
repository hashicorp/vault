import Ember from 'ember';
import IdentityModel from './_base';
import DS from 'ember-data';
import { queryRecord } from 'ember-computed-query';

const { computed } = Ember;
const { attr, belongsTo } = DS;

export default IdentityModel.extend({
  formFields: computed('type', function() {
    let fields = ['name', 'type', 'policies', 'metadata'];
    if (this.get('type') === 'internal') {
      return fields.concat(['memberGroupIds', 'memberEntityIds']);
    }
    return fields;
  }),
  name: attr('string'),
  type: attr('string', {
    defaultValue: 'internal',
    possibleValues: ['internal', 'external'],
  }),
  creationTime: attr('string', {
    readOnly: true,
  }),
  lastUpdateTime: attr('string', {
    readOnly: true,
  }),
  metadata: attr('object', {
    editType: 'kv',
  }),
  policies: attr({
    editType: 'stringArray',
  }),
  memberGroupIds: attr({
    label: 'Member Group IDs',
    editType: 'stringArray',
  }),
  memberEntityIds: attr({
    label: 'Member Entity IDs',
    editType: 'stringArray',
  }),
  hasMembers: computed(
    'memberEntityIds',
    'memberEntityIds.[]',
    'memberGroupIds',
    'memberGroupIds.[]',
    function() {
      let { memberEntityIds, memberGroupIds } = this.getProperties('memberEntityIds', 'memberGroupIds');
      let numEntities = (memberEntityIds && memberEntityIds.get('length')) || 0;
      let numGroups = (memberGroupIds && memberGroupIds.get('length')) || 0;
      return numEntities + numGroups > 0;
    }
  ),

  alias: belongsTo('identity/group-alias', { async: false, readOnly: true }),
  updatePath: queryRecord(
    'capabilities',
    context => {
      const { identityType, id } = context.getProperties('identityType', 'id');
      //identity/group/id/efb8b562-77fd-335f-a754-740373a778e6
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
      id: `identity/group-alias`,
    };
  }),

  canAddAlias: computed('aliasPath.canCreate', 'type', 'alias', function() {
    let type = this.get('type');
    let alias = this.get('alias');
    // internal groups can't have aliases, and external groups can only have one
    if (type === 'internal' || alias) {
      return false;
    }
    return this.get('aliasPath.canCreate');
  }),
});
