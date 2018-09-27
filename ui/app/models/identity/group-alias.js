import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import IdentityModel from './_base';
import DS from 'ember-data';
import identityCapabilities from 'vault/macros/identity-capabilities';

const { attr, belongsTo } = DS;

export default IdentityModel.extend({
  parentType: 'group',
  formFields: computed(function() {
    return ['name', 'mountAccessor'];
  }),
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

  updatePath: identityCapabilities(),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
});
