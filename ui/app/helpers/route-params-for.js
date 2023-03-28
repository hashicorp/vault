/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Helper from '@ember/component/helper';
import { inject as service } from '@ember/service';

export default Helper.extend({
  permissions: service(),
  compute([navItem]) {
    const permissions = this.permissions;
    return permissions.navPathParams(navItem);
  },
});
