import Ember from 'ember';
import IdentityModel from './_base';
import DS from 'ember-data';

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
});
