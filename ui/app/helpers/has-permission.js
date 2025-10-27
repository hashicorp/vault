/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-observers */
import Helper from '@ember/component/helper';
import { service } from '@ember/service';
import { observer } from '@ember/object';

export default Helper.extend({
  permissions: service(),
  namespace: service(),

  // Recompute when either ACL OR namespace path changes
  onPermissionsChange: observer(
    'permissions.exactPaths',
    'permissions.globPaths',
    'permissions.canViewAll',
    'permissions.chrootNamespace',
    'namespace.path',
    function () {
      this.recompute();
    }
  ),

  compute([route], params) {
    const { routeParams, requireAll } = params || {};
    return this.permissions.hasNavPermission(route, routeParams, requireAll);
  },
});
