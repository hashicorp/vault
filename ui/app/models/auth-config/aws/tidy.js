/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import AuthConfig from '../../auth-config';

export default AuthConfig.extend({
  safetyBuffer: attr({
    defaultValue: '72h',
    editType: 'ttl',
  }),

  disablePeriodicTidy: attr('boolean', {
    defaultValue: false,
  }),

  attrs: computed(function () {
    return expandAttributeMeta(this, ['safetyBuffer', 'disablePeriodicTidy']);
  }),
});
