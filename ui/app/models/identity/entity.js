import Ember from 'ember';
import IdentityModel from './_base';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import identityCapabilities from 'vault/macros/identity-capabilities';

const { computed } = Ember;

const { attr, hasMany } = DS;

export default IdentityModel.extend({
  formFields: ['name', 'disabled', 'policies', 'metadata'],
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

  updatePath: identityCapabilities(),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),

  aliasPath: lazyCapabilities(apiPath`identity/entity-alias`),
  canAddAlias: computed.alias('aliasPath.canCreate'),
});
