/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['action', 'roleType'],
  action: '',
  roleType: '',
  reset() {
    this.set('action', '');
    this.set('roleType', '');
  },
});
