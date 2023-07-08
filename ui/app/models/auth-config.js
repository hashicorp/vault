/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { belongsTo } from '@ember-data/model';

export default Model.extend({
  backend: belongsTo('auth-method', { inverse: 'authConfigs', readOnly: true, async: false }),
  getHelpUrl: function (backend) {
    return `/v1/auth/${backend}/config?help=1`;
  },
});
