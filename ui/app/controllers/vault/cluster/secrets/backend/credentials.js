/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['roleType'],
  roleType: '',
  reset() {
    this.set('roleType', '');
  },
});
