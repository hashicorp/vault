import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiRolesCreateRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    return this.pathHelp.getNewModel('pki/role', 'pki');
  }

  model() {
    return this.store.createRecord('pki/role', {
      backend: this.secretMountPath.currentPath,
    });
  }
}
