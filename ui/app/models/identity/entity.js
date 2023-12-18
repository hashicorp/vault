/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hasMany, attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import IdentityModel from './_base';
import apiPath from 'vault/utils/api-path';
import lazyCapabilities from 'vault/macros/lazy-capabilities';

export default IdentityModel.extend({
  formFields: computed(function () {
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
    editType: 'yield',
    isSectionHeader: true,
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
  updatePath: lazyCapabilities(apiPath`identity/entity/id/${'id'}`, 'id'),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),
  aliasPath: lazyCapabilities(apiPath`identity/entity-alias`),
  canAddAlias: alias('aliasPath.canCreate'),
  policyPath: lazyCapabilities(apiPath`sys/policies`),
  canCreatePolicies: alias('policyPath.canCreate'),
});
