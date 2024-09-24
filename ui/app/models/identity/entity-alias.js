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
  parentType: 'entity',
  formFields: computed(function () {
    return ['name', 'mountAccessor'];
  }),
  entity: belongsTo('identity/entity', { readOnly: true, async: false, inverse: 'aliases' }),

  name: attr('string'),
  canonicalId: attr('string'),
  mountAccessor: attr('string', {
    label: 'Auth Backend',
    editType: 'mountAccessor',
  }),
  metadata: attr({
    editType: 'kv',
    isSectionHeader: true,
  }),
  mountPath: attr('string', {
    readOnly: true,
  }),
  mountType: attr('string', {
    readOnly: true,
  }),
  creationTime: attr('string', {
    readOnly: true,
  }),
  lastUpdateTime: attr('string', {
    readOnly: true,
  }),
  mergedFromCanonicalIds: attr(),

  updatePath: identityCapabilities(),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
});
