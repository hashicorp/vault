import Helper from '@ember/component/helper';
import { inject as service } from '@ember/service';
import { observer } from '@ember/object';

export default Helper.extend({
  permissions: service(),
  onPermissionsChange: observer(
    'permissions.exactPaths',
    'permissions.globPaths',
    'permissions.canViewAll',
    function() {
      this.recompute();
    }
  ),

  compute([route], { routeParams, capability }) {
    let permissions = this.permissions;
    return permissions.hasNavPermission(route, routeParams, capability);
  },
});
