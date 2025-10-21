/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['roleType'],
  // used for database credentials
  roleType: '',
  reset() {
    this.set('roleType', '');
  },
});
