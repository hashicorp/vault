/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-observers */
import Helper from '@ember/component/helper';
import { service } from '@ember/service';
import { observer } from '@ember/object';

export default Helper.extend({
  permissions: service(),
  onPermissionsChange: observer(
    'permissions.exactPaths',
    'permissions.globPaths',
    'permissions.canViewAll',
    function () {
      this.recompute();
    }
  ),

  compute([route], params) {
    const { routeParams, requireAll } = params;
    const permissions = this.permissions;

    return permissions.hasNavPermission(route, routeParams, requireAll);
  },
});
