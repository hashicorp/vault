/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  approved: attr('boolean'),
  requestPath: attr('string'),
  // plain objects: requestEntity is { id, name } or null, authorizations is [{ id, name }]
  requestEntity: attr(),
  authorizations: attr(),

  authorizePath: lazyCapabilities(apiPath`sys/control-group/authorize`),
  canAuthorize: alias('authorizePath.canUpdate'),
  configurePath: lazyCapabilities(apiPath`sys/config/control-group`),
  canConfigure: alias('configurePath.canUpdate'),
});
