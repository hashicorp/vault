/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import KeyMixin from 'vault/mixins/key-mixin';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { alias } from '@ember/object/computed';

export default Model.extend(KeyMixin, {
  failedServerRead: attr('boolean'),
  auth: attr('string'),
  lease_duration: attr('number'),
  lease_id: attr('string'),
  renewable: attr('boolean'),

  code: attr('string'),

  codePath: lazyCapabilities(apiPath`${'backend'}/codes/${'id'}`, 'backend', 'id'),
  canRead: alias('codePath.canRead'),
});
