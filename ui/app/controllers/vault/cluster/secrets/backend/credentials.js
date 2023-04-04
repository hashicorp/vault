/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
