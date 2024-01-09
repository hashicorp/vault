/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { belongsTo, attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import IdentityModel from './_base';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import identityCapabilities from 'vault/macros/identity-capabilities';

export default IdentityModel.extend({
  formFields: computed('type', function () {
    const fields = ['name', 'type', 'policies', 'metadata'];
    if (this.type === 'internal') {
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
    editType: 'yield',
    isSectionHeader: true,
  }),
  memberGroupIds: attr({
    label: 'Member Group IDs',
    editType: 'searchSelect',
    isSectionHeader: true,
    fallbackComponent: 'string-list',
    models: ['identity/group'],
  }),
  parentGroupIds: attr({
    label: 'Parent Group IDs',
    editType: 'searchSelect',
    isSectionHeader: true,
    fallbackComponent: 'string-list',
    models: ['identity/group'],
  }),
  memberEntityIds: attr({
    label: 'Member Entity IDs',
    editType: 'searchSelect',
    isSectionHeader: true,
    fallbackComponent: 'string-list',
    models: ['identity/entity'],
  }),
  hasMembers: computed(
    'memberEntityIds',
    'memberEntityIds.[]',
    'memberGroupIds',
    'memberGroupIds.[]',
    function () {
      const { memberEntityIds, memberGroupIds } = this;
      const numEntities = (memberEntityIds && memberEntityIds.length) || 0;
      const numGroups = (memberGroupIds && memberGroupIds.length) || 0;
      return numEntities + numGroups > 0;
    }
  ),
  policyPath: lazyCapabilities(apiPath`sys/policies`),
  canCreatePolicies: alias('policyPath.canCreate'),
  alias: belongsTo('identity/group-alias', { async: false, readOnly: true, inverse: 'group' }),
  updatePath: identityCapabilities(),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),

  aliasPath: lazyCapabilities(apiPath`identity/group-alias`),
  canAddAlias: computed('aliasPath.canCreate', 'type', 'alias', function () {
    const type = this.type;
    const alias = this.alias;
    // internal groups can't have aliases, and external groups can only have one
    if (type === 'internal' || alias) {
      return false;
    }
    return this.aliasPath.canCreate;
  }),
});
