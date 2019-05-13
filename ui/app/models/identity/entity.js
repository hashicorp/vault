import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import IdentityModel from './_base';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import identityCapabilities from 'vault/macros/identity-capabilities';

const { attr, hasMany } = DS;

export default IdentityModel.extend({
  formFields: computed(function() {
    return ['name', 'disabled', 'policies', 'metadata'];
  }),
  name: attr('string'),
  disabled: attr('boolean', {
    defaultValue: false,
    label: 'Disable entity',
    helpText: 'All associated tokens cannot be used, but are not revoked.',
  }),
  mergedEntityIds: attr(),
  metadata: attr({
    editType: 'kv',
  }),
  policies: attr({
    label: 'Policies',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['policy/acl', 'policy/rgp'],
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

  updatePath: identityCapabilities(),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),

  aliasPath: lazyCapabilities(apiPath`identity/entity-alias`),
  canAddAlias: alias('aliasPath.canCreate'),
});
