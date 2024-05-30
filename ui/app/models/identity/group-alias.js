/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { belongsTo, attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import IdentityModel from './_base';
import identityCapabilities from 'vault/macros/identity-capabilities';

export default IdentityModel.extend({
  parentType: 'group',
  formFields: computed(function () {
    return ['name', 'mountAccessor'];
  }),
  group: belongsTo('identity/group', { readOnly: true, async: false, inverse: 'alias' }),

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
