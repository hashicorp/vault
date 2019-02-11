import Helper from '@ember/component/helper';
import { inject as service } from '@ember/service';

export default Helper.extend({
  permissions: service(),
  compute([route], { routeParams, capability }) {
    let permissions = this.permissions;
    return permissions.hasNavPermission(route, routeParams, capability);
  },
});
