import Ember from 'ember';
import IdentityModel from './_base';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import identityCapabilities from 'vault/macros/identity-capabilities';

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
  numMemberEntities: attr('number', {
    readOnly: true,
  }),
  numParentGroups: attr('number', {
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
  parentGroupIds: attr({
    label: 'Parent Group IDs',
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
  updatePath: identityCapabilities(),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),

  aliasPath: lazyCapabilities(apiPath`identity/group-alias`),
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
