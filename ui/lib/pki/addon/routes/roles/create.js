import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'pki/decorators/confirm-leave';
import PkiRolesIndexRoute from '.';

@withConfirmLeave
export default class PkiRolesCreateRoute extends PkiRolesIndexRoute {
  @service store;
  @service secretMountPath;

  model() {
    return this.store.createRecord('pki/role', {
      backend: this.secretMountPath.currentPath,
    });
  }
}
