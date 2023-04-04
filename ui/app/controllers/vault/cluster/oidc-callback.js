/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['state', 'code'],
  code: null,
  state: null,
});
