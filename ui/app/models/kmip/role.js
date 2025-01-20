/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import apiPath from 'vault/utils/api-path';
import lazyCapabilities from 'vault/macros/lazy-capabilities';
import { removeManyFromArray } from 'vault/helpers/remove-from-array';

export const COMPUTEDS = {
  operationFields: computed('newFields', function () {
    return this.newFields.filter((key) => key.startsWith('operation'));
  }),

  operationFieldsWithoutSpecial: computed('operationFields', function () {
    return removeManyFromArray(this.operationFields, ['operationAll', 'operationNone']);
  }),

  tlsFields: computed(function () {
    return ['tlsClientKeyBits', 'tlsClientKeyType', 'tlsClientTtl'];
  }),

  // For rendering on the create/edit pages
  defaultFields: computed('newFields', 'operationFields', 'tlsFields', function () {
    const excludeFields = ['role'].concat(this.operationFields, this.tlsFields);
    return removeManyFromArray(this.newFields, excludeFields);
  }),

  // For adapter/serializer
  nonOperationFields: computed('newFields', 'operationFields', function () {
    return removeManyFromArray(this.newFields, this.operationFields);
  }),
};

export default Model.extend(COMPUTEDS, {
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  name: attr({ readOnly: true }),

  fieldGroups: computed('fields', 'defaultFields.length', 'tlsFields', function () {
    const groups = [{ TLS: this.tlsFields }];
    if (this.defaultFields.length) {
      groups.unshift({ default: this.defaultFields });
    }
    const ret = fieldToAttrs(this, groups);
    return ret;
  }),

  operationFormFields: computed('operationFieldsWithoutSpecial', function () {
    const objects = [
      'operationCreate',
      'operationActivate',
      'operationGet',
      'operationLocate',
      'operationRekey',
      'operationRevoke',
      'operationDestroy',
    ];

    const attributes = ['operationAddAttribute', 'operationGetAttributes'];
    const server = ['operationDiscoverVersions'];
    const others = removeManyFromArray(this.operationFieldsWithoutSpecial, [
      ...objects,
      ...attributes,
      ...server,
    ]);
    const groups = [
      { 'Managed Cryptographic Objects': objects },
      { 'Object Attributes': attributes },
      { Server: server },
    ];
    if (others.length) {
      groups.push({
        Other: others,
      });
    }
    return fieldToAttrs(this, groups);
  }),
  tlsFormFields: computed('tlsFields', function () {
    return expandAttributeMeta(this, this.tlsFields);
  }),
  fields: computed('defaultFields', function () {
    return expandAttributeMeta(this, this.defaultFields);
  }),

  updatePath: lazyCapabilities(apiPath`${'backend'}/scope/${'scope'}/role/${'id'}`, 'backend', 'scope', 'id'),
});
