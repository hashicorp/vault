import Helper from '@ember/component/helper';
import { inject as service } from '@ember/service';

export default Helper.extend({
  permissions: service(),
  compute([navItem]) {
    let permissions = this.permissions;
    return permissions.navPathParams(navItem);
  },
});
