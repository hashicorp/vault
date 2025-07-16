/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { belongsTo } from '@ember-data/model';

export default Model.extend({
  backend: belongsTo('auth-method', {
    inverse: 'authConfigs',
    readOnly: true,
    async: false,
    as: 'auth-config',
  }),
});
