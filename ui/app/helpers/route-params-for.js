/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Helper from '@ember/component/helper';
import { service } from '@ember/service';

export default Helper.extend({
  permissions: service(),
  compute([navItem]) {
    const permissions = this.permissions;
    return permissions.navPathParams(navItem);
  },
});
