import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiRolesCreateRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    return this.pathHelp.getNewModel('pki/pki-role-engine', 'pki');
  }

  model() {
    return this.store.createRecord('pki/pki-role-engine', {
      backend: this.secretMountPath.currentPath,
    });
  }
}
