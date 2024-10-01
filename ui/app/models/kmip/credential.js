/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { computed } from '@ember/object';
import apiPath from 'vault/utils/api-path';
import lazyCapabilities from 'vault/macros/lazy-capabilities';

export default Model.extend({
  backend: attr({ readOnly: true }),
  scope: attr({ readOnly: true }),
  role: attr({ readOnly: true }),
  certificate: attr('string', { readOnly: true }),
  caChain: attr({ readOnly: true }),
  privateKey: attr('string', {
    readOnly: true,
    sensitive: true,
  }),
  format: attr('string', {
    possibleValues: ['pem', 'der', 'pem_bundle'],
    defaultValue: 'pem',
    label: 'Certificate format',
  }),
  fieldGroups: computed(function () {
    const groups = [
      {
        default: ['format'],
      },
    ];

    return fieldToAttrs(this, groups);
  }),
  deletePath: lazyCapabilities(
    apiPath`${'backend'}/scope/${'scope'}/role/${'role'}/credentials/revoke`,
    'backend',
    'scope',
    'role'
  ),
});
