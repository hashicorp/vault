/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { hasMany, belongsTo, attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  approved: attr('boolean'),
  requestPath: attr('string'),
  requestEntity: belongsTo('identity/entity', { async: false, inverse: null }),
  authorizations: hasMany('identity/entity', { async: false, inverse: null }),

  authorizePath: lazyCapabilities(apiPath`sys/control-group/authorize`),
  canAuthorize: alias('authorizePath.canUpdate'),
  configurePath: lazyCapabilities(apiPath`sys/config/control-group`),
  canConfigure: alias('configurePath.canUpdate'),
});
