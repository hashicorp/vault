/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import apiPath from 'vault/utils/api-path';

import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import lazyCapabilities from 'vault/macros/lazy-capabilities';

export default Model.extend({
  name: attr('string'),
  backend: attr({ readOnly: true }),
  attrs: computed(function () {
    return expandAttributeMeta(this, ['name']);
  }),
  updatePath: lazyCapabilities(apiPath`${'backend'}/scope/${'id'}`, 'backend', 'id'),
});
