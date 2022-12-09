import FetchConfigRoute from '../../fetch-config';
import { hash } from 'ember-concurrency';
export default class KubernetesRoleCredentialsRoute extends FetchConfigRoute {
  model() {
    const roleModel = this.modelFor('roles.role');

    return hash({
      roleModel,
      backend: this.secretMountPath.get(),
    });
  }
}
